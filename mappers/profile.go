package mappers

import (
	"github.com/google/uuid"
	"n_users/entity"
)

func ToProfile(i entity.CreateProfileRequest) entity.Profile {
	p := entity.Profile{}

	p.FullName = i.FullName
	p.EmailID = i.EmailID
	p.Mobile = i.Mobile
	p.Gender = i.Gender
	p.BirthDate = i.BirthDate
	p.CityID = i.CityID
	p.CountryID = i.CountryID
	p.Address = i.Address
	p.Latitude = i.Latitude
	p.Longitude = i.Longitude

	p.Active = true
	p.TenantID = "default"
	p.ProfileID = uuid.New().String()

	return p
}
