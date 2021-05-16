package repo

import (
	"errors"
	"n_users/entity"

	"go.uber.org/zap"

	"github.com/jinzhu/gorm"
	// required for postgre
	_ "github.com/lib/pq"
)

// ProfileRepo represent interface to perform CRUD on database
type ProfileRepo interface {
	Create(profile entity.Profile) (string, error)
	Delete(profileID string, tenantID string) (bool, error)
	Search(query string, limit int, offset int, sortBy string, tenantID string) ([]entity.Profile, error)
	Update(filters map[string]interface{}, fieldsToUpdate map[string]interface{}) (bool, error)
	UploadProfileImage(profileID string, image []byte) (bool, error)
	SafeClose()
}

type profileRepo struct {
	DB *gorm.DB
}

// New creates new object of ProfileRepo
func New(dialect string, dbName string) (ProfileRepo, error) {
	db, err := gorm.Open(dialect, dbName)

	if err != nil {
		zap.L().Fatal("failed to connect database", zap.Error(err))
		return nil, err
	}

	db.DB().SetMaxIdleConns(5)
	db.DB().SetMaxOpenConns(100)
	db.LogMode(true)
	//db.SetLogger(zap.L()) TODO: fix this

	// Migrate the schema
	db.AutoMigrate(&entity.Profile{})

	defer zap.L().Info("sql database setup completed")
	return &profileRepo{DB: db}, nil
}

func (pr *profileRepo) SafeClose() {
	pr.DB.Close()
}

func (pr *profileRepo) Create(profile entity.Profile) (string, error) {
	res := pr.DB.Create(profile)
	if res.Error != nil {
		zap.L().Error(res.Error.Error())
		return "", res.Error
	}

	return profile.ProfileID, nil
}

func (pr *profileRepo) Delete(profileID string, tenantID string) (bool, error) {
	var profile entity.Profile
	res := pr.DB.Where("profile_id = ? AND tenant_id = ?", profileID, tenantID).Delete(&profile)

	if res.Error != nil {
		zap.L().Error(res.Error.Error())
		return false, res.Error
	}

	return res.RowsAffected > 0, nil
}

func (pr *profileRepo) Search(query string, limit int, offset int, sortBy string, tenantID string) ([]entity.Profile, error) {
	var profiles []entity.Profile

	res := pr.DB.Where(query).
		Where("tenant_id = ?", tenantID).
		Limit(limit).
		Offset(offset).
		Order(sortBy).
		Find(&profiles)

	if res.Error != nil {
		zap.L().Error(res.Error.Error())
		return nil, res.Error
	}

	return profiles, nil
}

func (pr *profileRepo) Update(filters map[string]interface{}, fieldsToUpdate map[string]interface{}) (bool, error) {

	profile := entity.Profile{}

	if value, ok := filters["profile_id"]; ok {
		profile.ProfileID = value.(string)
	}

	if value, ok := filters["tenant_id"]; ok {
		profile.TenantID = value.(string)
	}

	res := pr.DB.
		Model(&profile).
		Where("profile_id = ? and tenant_id = ?", profile.ProfileID, profile.TenantID).
		Updates(fieldsToUpdate)

	if res.Error != nil {
		zap.L().Error(res.Error.Error())
		return false, res.Error
	}

	return res.RowsAffected > 0, nil
}

func (pr *profileRepo) UploadProfileImage(profileID string, image []byte) (bool, error) {
	return false, errors.New("not implemented")
}
