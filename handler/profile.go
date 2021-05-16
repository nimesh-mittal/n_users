package handler

import (
	"encoding/json"
	"log"
	"net/http"
	"os"
	"strconv"

	"n_users/controller"
	"n_users/entity"
	"n_users/mappers"

	"n_users/gateway/s3store"
	"n_users/repo"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/go-chi/chi/v5"
)

const maxUploadFileSize = int64(2 * 1024000)
const S3ImagePrefix = "https://s3.ap-south-1.amazonaws.com/images.repo.bucket1/"

// ProfileHandler handles profile endpoints
type ProfileHandler interface {
	CreateProfile(w http.ResponseWriter, r *http.Request)
	DeleteProfile(w http.ResponseWriter, r *http.Request)
	SearchProfile(w http.ResponseWriter, r *http.Request)
	UpdateProfile(w http.ResponseWriter, r *http.Request)
	UploadProfileImage(w http.ResponseWriter, r *http.Request)
	NewProfileRouter() http.Handler
}

type profileHandler struct {
	ProfileService controller.ProfileService
	AWSSession     *session.Session
}

// NewProfileHandler creates ProfileHandler
func NewProfileHandler() ProfileHandler {
	url := os.Getenv("POSTGRE_URL_VALUE")
	pr, err := repo.New("postgres", url)

	if err != nil {
		log.Fatal("error init profile handler", err)
	}

	// Create an AWS session
	AWSRegion := os.Getenv("AWS_REGION")
	AWSSecretID := os.Getenv("AWS_SECRET_ID")
	AWSSecret := os.Getenv("AWS_SECRET")
	s, err := session.NewSession(&aws.Config{
		Region:      aws.String(AWSRegion),
		Credentials: credentials.NewStaticCredentials(AWSSecretID, AWSSecret, ""),
	})

	if err != nil {
		// TODO: check if regular session refresh is required?
		log.Fatal("error creating AWS session", err)
	}

	return &profileHandler{ProfileService: controller.New(pr), AWSSession: s}
}

// NewProfileRouter returns new router for profile endpoints
func (h *profileHandler) NewProfileRouter() http.Handler {
	r := chi.NewRouter()

	r.Post("/", h.CreateProfile)
	r.Delete("/{ProfileID}", h.DeleteProfile)
	r.Put("/{ProfileID}", h.UpdateProfile)
	r.Post("/_search", h.SearchProfile)
	r.Put("/{ProfileID}/_upload", h.UploadProfileImage)

	return r
}

func (h *profileHandler) CreateProfile(w http.ResponseWriter, r *http.Request) {
	decoder := json.NewDecoder(r.Body)
	defer r.Body.Close()

	var createProfileRequest entity.CreateProfileRequest
	err := decoder.Decode(&createProfileRequest)

	if err != nil {
		e := entity.NewError("invalid create profile request " + err.Error())
		res, _ := json.Marshal(e)
		w.Write(res)
		return
	}

	tenant := r.Header.Get("ntenant")
	if len(tenant) == 0 {
		tenant = "default"
	}

	p := mappers.ToProfile(createProfileRequest)
	p.TenantID = tenant

	id, err := h.ProfileService.Create(p)

	if err != nil {
		e := entity.NewError("error processing create profile request " + err.Error())
		res, _ := json.Marshal(e)
		w.Write(res)
		return
	}

	e := entity.CreateProfileResponse{ProfileID: id, TenantID: tenant}
	res, _ := json.Marshal(e)
	w.Write(res)
}

func (h *profileHandler) DeleteProfile(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "ProfileID")

	tenant := r.Header.Get("ntenant")
	if len(tenant) == 0 {
		tenant = "default"
	}

	status, err := h.ProfileService.Delete(id, tenant)

	if err != nil {
		e := entity.NewError("error processing delete profile request " + err.Error())
		res, _ := json.Marshal(e)
		w.Write(res)
		return
	}

	e := entity.SuccessResponse{Status: strconv.FormatBool(status)}
	res, _ := json.Marshal(e)
	w.Write(res)
}

func (h *profileHandler) UpdateProfile(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "ProfileID")

	decoder := json.NewDecoder(r.Body)
	defer r.Body.Close()

	var updateProfileRequest entity.UpdateProfileRequest
	err := decoder.Decode(&updateProfileRequest)

	if err != nil {
		e := entity.NewError("invalid update profile request " + err.Error())
		res, _ := json.Marshal(e)
		w.Write(res)
		return
	}

	tenant := r.Header.Get("ntenant")
	if len(tenant) == 0 {
		tenant = "default"
	}

	filter := map[string]interface{}{"profile_id": id, "tenant_id": tenant}

	fieldsToUpdate := map[string]interface{}{
		"full_name":  updateProfileRequest.FullName,
		"gender":     updateProfileRequest.Gender,
		"email_id":   updateProfileRequest.EmailID,
		"mobile":     updateProfileRequest.Mobile,
		"birth_date": updateProfileRequest.BirthDate,
		"address":    updateProfileRequest.Address,
	}

	fieldsToUpdate = entity.RemoveEmptyValues(fieldsToUpdate)
	status, err := h.ProfileService.Update(filter, fieldsToUpdate)
	if err != nil {
		e := entity.NewError("error processing update profile request")
		res, _ := json.Marshal(e)
		w.Write(res)
		return
	}

	e := entity.SuccessResponse{Status: strconv.FormatBool(status)}
	res, _ := json.Marshal(e)
	w.Write(res)
}

func (h *profileHandler) SearchProfile(w http.ResponseWriter, r *http.Request) {
	decoder := json.NewDecoder(r.Body)
	defer r.Body.Close()

	var searchProfileRequest entity.SearchProfileRequest
	err := decoder.Decode(&searchProfileRequest)

	if err != nil {
		e := entity.NewError("invalid search profile request " + err.Error())
		res, _ := json.Marshal(e)
		w.Write(res)
		return
	}

	tenant := r.Header.Get("ntenant")
	if len(tenant) == 0 {
		tenant = "default"
	}

	query := searchProfileRequest.Query
	sortBy := searchProfileRequest.SortBy
	profiles, err := h.ProfileService.Search(query,
		int(searchProfileRequest.Limit),
		int(searchProfileRequest.Offset),
		sortBy,
		tenant)

	if err != nil {
		e := entity.NewError("Error processing search profile request. " + err.Error())
		res, _ := json.Marshal(e)
		w.Write(res)
		return
	}

	res, _ := json.Marshal(profiles)
	w.Write(res)
}

func (h *profileHandler) UploadProfileImage(w http.ResponseWriter, r *http.Request) {
	// allow only 2MB of file size
	err := r.ParseMultipartForm(maxUploadFileSize)
	if err != nil {
		res, _ := entity.NewErrorJSON("Image too large, max file size allowed is 2MB. " + err.Error())
		w.Write(res)
		return
	}

	file, fileHeader, err := r.FormFile("profile_image")
	if err != nil {
		res, _ := entity.NewErrorJSON("Error fetching uploaded file. " + err.Error())
		w.Write(res)
		return
	}
	defer file.Close()

	fileName, err := s3store.UploadFileToS3(h.AWSSession, file, fileHeader)
	if err != nil {
		res, _ := entity.NewErrorJSON("Unable to upload file. " + err.Error())
		w.Write(res)
		return
	}

	profileImageURL := S3ImagePrefix + fileName

	// update profile in database with image url
	id := chi.URLParam(r, "ProfileID")
	tenant := r.Header.Get("ntenant")
	if len(tenant) == 0 {
		tenant = "default"
	}

	filter := map[string]interface{}{"profile_id": id, "tenant_id": tenant}
	fieldsToUpdate := map[string]interface{}{"profile_image_url": profileImageURL}
	_, err = h.ProfileService.Update(filter, fieldsToUpdate)
	if err != nil {
		res, _ := entity.NewErrorJSON("Unable to update profile with image url. " + err.Error())
		w.Write(res)
		return
	}

	// prepare response object and return
	e := entity.SuccessResponse{Status: "File upload successful. File uploaded at " + profileImageURL}
	res, _ := json.Marshal(e)
	w.Write(res)
}
