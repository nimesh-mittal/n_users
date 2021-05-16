package controller

import (
	"n_users/entity"
	"n_users/repo"

	"go.uber.org/zap"
)

// ProfileService represents interface to manage profile
type ProfileService interface {
	Create(profile entity.Profile) (string, error)
	Delete(profileID string, tenantID string) (bool, error)
	Search(query string, limit int, offset int, sortBy string, tenantID string) ([]entity.Profile, error)
	Update(filters map[string]interface{}, fieldsToUpdate map[string]interface{}) (bool, error)
	UploadProfileImage(profileID string, image []byte) (bool, error)
}

type service struct {
	Repo repo.ProfileRepo
}

// New creates new object of ProfileService
func New(repo repo.ProfileRepo) ProfileService {
	return &service{Repo: repo}
}

func (s *service) Create(profile entity.Profile) (string, error) {
	zap.L().Info("receive create profile request",
		zap.String("profile_id", profile.ProfileID),
		zap.String("tenant_id", profile.TenantID))

	id, err := s.Repo.Create(profile)

	if err != nil {
		zap.L().Error("error processing created profile request", zap.Error(err))
		return "", err
	}

	return id, nil
}

func (s *service) Delete(profileID string, tenantID string) (bool, error) {
	zap.L().Info("receive delete profile request",
		zap.String("profile_id", profileID),
		zap.String("tenant_id", tenantID))

	status, err := s.Repo.Delete(profileID, tenantID)

	if err != nil {
		zap.L().Error("error processing created profile request", zap.Error(err))
		return false, err
	}

	return status, nil
}

func (s *service) Search(query string, limit int, offset int, sortBy string, tenantID string) ([]entity.Profile, error) {
	zap.L().Info("receive search profile request",
		zap.String("query", query),
		zap.String("tenant_id", tenantID))

	profiles, err := s.Repo.Search(query, limit, offset, sortBy, tenantID)

	if err != nil {
		zap.L().Error("error processing created profile request", zap.Error(err))
		return nil, err
	}

	return profiles, nil
}

func (s *service) Update(filters map[string]interface{}, fieldsToUpdate map[string]interface{}) (bool, error) {
	zap.L().Info("receive update profile request")

	status, err := s.Repo.Update(filters, fieldsToUpdate)

	if err != nil {
		zap.L().Error("error processing update profile request", zap.Error(err))
		return false, err
	}

	return status, nil
}

func (s *service) UploadProfileImage(profileID string, image []byte) (bool, error) {
	zap.L().Info("receive upload profile image request")

	status, err := s.Repo.UploadProfileImage(profileID, image)

	if err != nil {
		zap.L().Error("error processing update profile request", zap.Error(err))
		return false, err
	}

	return status, nil
}
