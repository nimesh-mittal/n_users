package handler

import (
	"encoding/json"
	"errors"
	"n_users/controller"
	"n_users/entity"
	"n_users/mocks"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/golang/mock/gomock"
)

func GetCreateProfileRequest() *http.Request {
	data := entity.CreateProfileRequest{
		FullName: "Nimesh",
		EmailID:  "nimesh@gmail.com",
		Mobile:   "8888800000",
	}
	b, _ := json.Marshal(data)
	payload := strings.NewReader(string(b))
	req, _ := http.NewRequest(http.MethodPost, "http://localhost:8085/", payload)
	return req
}

func GetMockCreateProfileHandler(t *testing.T, triggerError bool) ProfileHandler {
	mockCtrl := gomock.NewController(t)

	mockProfileRepo := mocks.NewMockProfileRepo(mockCtrl)

	var err error
	if triggerError == true {
		err = errors.New("error")
	}

	mockProfileRepo.EXPECT().Create(gomock.Any()).Return("101", err).Times(1)

	return &profileHandler{ProfileService: controller.New(mockProfileRepo)}
}

func TestCreateProfile(t *testing.T) {
	w := httptest.NewRecorder()

	GetMockCreateProfileHandler(t, false).NewProfileRouter().ServeHTTP(w, GetCreateProfileRequest())
	resp := w.Result()

	if resp.StatusCode != http.StatusOK {
		t.Errorf("create profile didn’t respond 200 OK: %s", resp.Status)
	}

	var sr entity.CreateProfileResponse
	if err := json.NewDecoder(resp.Body).Decode(&sr); err != nil {
		t.Errorf("create profile response parsing error %s", err)
	}

	if sr.ProfileID != "101" {
		t.Errorf("create profile response id is %s but expected 101", sr.ProfileID)
	}
}

func GetDeleteProfileRequest() *http.Request {
	req, _ := http.NewRequest(http.MethodDelete, "http://localhost:8085/201", nil)
	return req
}

func GetMockDeleteProfileHandler(t *testing.T, triggerError bool) ProfileHandler {
	mockCtrl := gomock.NewController(t)

	mockProfileRepo := mocks.NewMockProfileRepo(mockCtrl)

	var err error
	if triggerError == true {
		err = errors.New("error")
	}

	mockProfileRepo.EXPECT().Delete("201", gomock.Any()).Return(true, err).Times(1)

	return &profileHandler{ProfileService: controller.New(mockProfileRepo)}
}

func TestDeleteProfile(t *testing.T) {
	w := httptest.NewRecorder()

	GetMockDeleteProfileHandler(t, false).NewProfileRouter().ServeHTTP(w, GetDeleteProfileRequest())
	resp := w.Result()

	if resp.StatusCode != http.StatusOK {
		t.Errorf("delete profile didn’t respond 200 OK: %s", resp.Status)
	}

	var sr entity.SuccessResponse
	if err := json.NewDecoder(resp.Body).Decode(&sr); err != nil {
		t.Errorf("delete profile response parsing error %s", err)
	}

	if sr.Status != "true" {
		t.Errorf("delete profile status is %s but expected true", sr.Status)
	}
}

func GetUpdateProfileRequest() *http.Request {
	data := entity.UpdateProfileRequest{
		FullName: "Nimesh",
		EmailID:  "nimesh@gmail.com",
		Mobile:   "8888800000",
	}
	b, _ := json.Marshal(data)
	payload := strings.NewReader(string(b))
	req, _ := http.NewRequest(http.MethodPut, "http://localhost:8085/301", payload)
	return req
}

func GetMockUpdateProfileHandler(t *testing.T, triggerError bool) ProfileHandler {
	mockCtrl := gomock.NewController(t)

	mockProfileRepo := mocks.NewMockProfileRepo(mockCtrl)

	var err error
	if triggerError == true {
		err = errors.New("error")
	}

	mockProfileRepo.EXPECT().Update(gomock.Any(), gomock.Any()).Return(true, err).Times(1)

	return &profileHandler{ProfileService: controller.New(mockProfileRepo)}
}

func TestUpdateProfile(t *testing.T) {
	w := httptest.NewRecorder()

	GetMockUpdateProfileHandler(t, false).NewProfileRouter().ServeHTTP(w, GetUpdateProfileRequest())
	resp := w.Result()

	if resp.StatusCode != http.StatusOK {
		t.Errorf("update profile didn’t respond 200 OK: %s", resp.Status)
	}

	var sr entity.SuccessResponse
	if err := json.NewDecoder(resp.Body).Decode(&sr); err != nil {
		t.Errorf("update profile response parsing error %s", err)
	}

	if sr.Status != "true" {
		t.Errorf("update profile status is %s but expected true", sr.Status)
	}
}
