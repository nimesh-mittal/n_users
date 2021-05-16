package entity

import "time"

// Profile represents user profile object
type Profile struct {
	TenantID        string `json:"tenant_id" gorm:"primaryKey" validate:"required"`
	ProfileID       string `json:"profile_id" gorm:"primaryKey" validate:"required"`
	FullName        string `json:"full_name" validate:"required"`
	Gender          string
	EmailID         string `json:"email_id" gorm:"unique" validate:"required"`
	Mobile          string `json:"mobile" gorm:"unique" validate:"required"`
	BirthDate       time.Time
	CityID          string
	CountryID       string
	Address         string
	Latitude        float64
	Longitude       float64
	ProfileImageURL string
	// who columns
	Active    bool
	CreatedBy string
	CreatedAt time.Time
	UpdatedBy string
	UpdatedAt time.Time
	DeletedBy string
	DeletedAt *time.Time
}

// CreateProfileRequest represent create profile request
type CreateProfileRequest struct {
	FullName  string `json:"full_name" validate:"required"`
	Gender    string
	EmailID   string `json:"email_id" validate:"required"`
	Mobile    string
	BirthDate time.Time `json:"birth_date" validate:"required"`
	CityID    string    `json:"city_id" validate:"required"`
	CountryID string    `json:"country_id" validate:"required"`
	Address   string
	Latitude  float64
	Longitude float64
}

// CreateProfileResponse represent create profile response
type CreateProfileResponse struct {
	TenantID  string
	ProfileID string
}

// UpdateProfileRequest represent update profile request
type UpdateProfileRequest struct {
	FullName  string `json:"full_name"`
	Gender    string
	EmailID   string `json:"email_id"`
	Mobile    string
	BirthDate time.Time `json:"birth_date"`
	Address   string
}

// SearchProfileRequest represent search profile request
type SearchProfileRequest struct {
	Query  string
	SortBy string `json:"sort_by"`
	Limit  int64
	Offset int64
}
