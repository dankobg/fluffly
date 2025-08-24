package dto

import (
	"encoding/json"
	"fmt"

	api "github.com/dankobg/fluffly/api/gen"
	"github.com/dankobg/fluffly/db/dbmodel"
	"github.com/dankobg/fluffly/db/queries"
	"github.com/dankobg/fluffly/ptr"
)

func GetOrganizationByIdRowToResponse(row queries.GetOrganizationByIdRow) api.Organization {
	latDec := row.OrganizationContactAddressLat.Ptr()
	lngDec := row.OrganizationContactAddressLNG.Ptr()
	var lat *api.Latitude
	var lng *api.Longitude
	if latDec != nil {
		val, _ := latDec.Float64()
		lat = ptr.Of(float32(val))
	}
	if lngDec != nil {
		val, _ := lngDec.Float64()
		lat = ptr.Of(float32(val))
	}
	photos := []api.OrganizationPhoto{}
	socials := []api.OrganizationSocial{}
	if err := json.Unmarshal(row.Photos.Val, &photos); err != nil {
		fmt.Println("could not unmarshal photos: ", err)
	}
	if err := json.Unmarshal(row.Socials.Val, &socials); err != nil {
		fmt.Println("could not unmarshal socials: ", err)
	}

	org := api.Organization{
		ID:               row.ID,
		Name:             row.Name,
		AdoptionPolicy:   row.AdoptionPolicy.Ptr(),
		AdoptionURL:      row.AdoptionURL.Ptr(),
		Distance:         row.Distance.Ptr(),
		MissionStatement: row.MissionStatement.Ptr(),
		Website:          row.Website.Ptr(),
		Contact: api.Contact{
			Address: api.Address{
				City: row.OrganizationContactAddressCity.MustGet(),
				Country: api.Country{
					ID:         row.OrganizationContactAddressCountryID.MustGet(),
					Name:       row.OrganizationContactAddressCountryName.MustGet(),
					IsoAlpha2:  row.OrganizationContactAddressCountryIsoAlpha2.MustGet(),
					IsoAlpha3:  row.OrganizationContactAddressCountryIsoAlpha3.MustGet(),
					IsoNumeric: row.OrganizationContactAddressCountryIsoNumeric.MustGet(),
					CreatedAt:  row.OrganizationContactAddressCountryCreatedAt.MustGet(),
					UpdatedAt:  row.OrganizationContactAddressCountryUpdatedAt.MustGet(),
				},
				ID:            row.OrganizationContactAddressID.MustGet(),
				Lat:           lat,
				Lng:           lng,
				Note:          row.OrganizationContactAddressNote.Ptr(),
				PostalCode:    row.OrganizationContactAddressPostalCode.Ptr(),
				Region:        row.OrganizationContactAddressRegion.Ptr(),
				StreetAddress: row.OrganizationContactAddressStreetAddress.MustGet(),
				StreetNumber:  row.OrganizationContactAddressStreetNumber.Ptr(),
				UnitNumber:    row.OrganizationContactAddressUnitNumber.Ptr(),
				CreatedAt:     row.OrganizationContactAddressCreatedAt.MustGet(),
				UpdatedAt:     row.OrganizationContactAddressUpdatedAt.MustGet(),
			},
			Email:     row.OrganizationContactEmail.MustGet(),
			Phone:     row.OrganizationContactPhone.MustGet(),
			CreatedAt: row.OrganizationContactCreatedAt.Ptr(),
			UpdatedAt: row.OrganizationContactUpdatedAt.Ptr(),
		},
		WorkHour: &api.OrganizationWorkHour{
			ID:             row.OrganizationWorkHourID.MustGet(),
			OrganizationID: row.OrganizationWorkHourOrganizationID.MustGet(),
			Monday:         row.OrganizationWorkHourSaturday.Ptr(),
			Tuesday:        row.OrganizationWorkHourTuesday.Ptr(),
			Wednesday:      row.OrganizationWorkHourWednesday.Ptr(),
			Thursday:       row.OrganizationWorkHourThursday.Ptr(),
			Friday:         row.OrganizationWorkHourFriday.Ptr(),
			Saturday:       row.OrganizationWorkHourSaturday.Ptr(),
			Sunday:         row.OrganizationWorkHourSunday.Ptr(),
			CreatedAt:      row.OrganizationWorkHourCreatedAt.Ptr(),
			UpdatedAt:      row.OrganizationWorkHourUpdatedAt.Ptr(),
		},
		Photos:    photos,
		Socials:   socials,
		CreatedAt: row.CreatedAt,
		UpdatedAt: row.UpdatedAt,
	}
	return org
}

func ListOrganizationsRowToResponse(row queries.ListOrganizationsRow_) api.Organization {
	latDec := row.OrganizationContact.Address.Lat.Ptr()
	lngDec := row.OrganizationContact.Address.LNG.Ptr()
	var lat *api.Latitude
	var lng *api.Longitude
	if latDec != nil {
		val, _ := latDec.Float64()
		lat = ptr.Of(float32(val))
	}
	if lngDec != nil {
		val, _ := lngDec.Float64()
		lat = ptr.Of(float32(val))
	}
	photos := []api.OrganizationPhoto{}
	socials := []api.OrganizationSocial{}
	if err := json.Unmarshal(row.Photos.Val, &photos); err != nil {
		fmt.Println("could not unmarshal photos: ", err)
	}
	if err := json.Unmarshal(row.Socials.Val, &socials); err != nil {
		fmt.Println("could not unmarshal socials: ", err)
	}

	org := api.Organization{
		ID:               row.ID,
		Name:             row.Name,
		AdoptionPolicy:   row.AdoptionPolicy.Ptr(),
		AdoptionURL:      row.AdoptionURL.Ptr(),
		Distance:         row.Distance.Ptr(),
		MissionStatement: row.MissionStatement.Ptr(),
		Website:          row.Website.Ptr(),
		Contact: api.Contact{
			Address: api.Address{
				City: row.OrganizationContact.Address.City.MustGet(),
				Country: api.Country{
					ID:         row.OrganizationContact.Address.Country.ID.MustGet(),
					Name:       row.OrganizationContact.Address.Country.Name.MustGet(),
					IsoAlpha2:  row.OrganizationContact.Address.Country.IsoAlpha2.MustGet(),
					IsoAlpha3:  row.OrganizationContact.Address.Country.IsoAlpha3.MustGet(),
					IsoNumeric: row.OrganizationContact.Address.Country.IsoNumeric.MustGet(),
					CreatedAt:  row.OrganizationContact.Address.Country.CreatedAt.MustGet(),
					UpdatedAt:  row.OrganizationContact.Address.Country.UpdatedAt.MustGet(),
				},
				ID:            row.OrganizationContact.Address.ID.MustGet(),
				Lat:           lat,
				Lng:           lng,
				Note:          row.OrganizationContact.Address.Note.Ptr(),
				PostalCode:    row.OrganizationContact.Address.PostalCode.Ptr(),
				Region:        row.OrganizationContact.Address.Region.Ptr(),
				StreetAddress: row.OrganizationContact.Address.StreetAddress.MustGet(),
				StreetNumber:  row.OrganizationContact.Address.StreetNumber.Ptr(),
				UnitNumber:    row.OrganizationContact.Address.UnitNumber.Ptr(),
				CreatedAt:     row.OrganizationContact.Address.CreatedAt.MustGet(),
				UpdatedAt:     row.OrganizationContact.Address.UpdatedAt.MustGet(),
			},
			Email:     row.OrganizationContact.Email.MustGet(),
			Phone:     row.OrganizationContact.Phone.MustGet(),
			CreatedAt: row.OrganizationContact.CreatedAt.Ptr(),
			UpdatedAt: row.OrganizationContact.UpdatedAt.Ptr(),
		},
		WorkHour: &api.OrganizationWorkHour{
			ID:             row.OrganizationWorkHour.ID.MustGet(),
			OrganizationID: row.OrganizationWorkHour.OrganizationID.MustGet(),
			Monday:         row.OrganizationWorkHour.Saturday.Ptr(),
			Tuesday:        row.OrganizationWorkHour.Tuesday.Ptr(),
			Wednesday:      row.OrganizationWorkHour.Wednesday.Ptr(),
			Thursday:       row.OrganizationWorkHour.Thursday.Ptr(),
			Friday:         row.OrganizationWorkHour.Friday.Ptr(),
			Saturday:       row.OrganizationWorkHour.Saturday.Ptr(),
			Sunday:         row.OrganizationWorkHour.Sunday.Ptr(),
			CreatedAt:      row.OrganizationWorkHour.CreatedAt.Ptr(),
			UpdatedAt:      row.OrganizationWorkHour.UpdatedAt.Ptr(),
		},
		Photos:    photos,
		Socials:   socials,
		CreatedAt: row.CreatedAt,
		UpdatedAt: row.UpdatedAt,
	}
	return org
}

func OrganizationToResponse(organization *dbmodel.Organization) api.Organization {
	latDec := organization.R.OrganizationContact.R.Address.Lat.Ptr()
	lngDec := organization.R.OrganizationContact.R.Address.LNG.Ptr()
	var lat *api.Latitude
	var lng *api.Longitude
	if latDec != nil {
		val, _ := latDec.Float64()
		lat = ptr.Of(float32(val))
	}
	if lngDec != nil {
		val, _ := lngDec.Float64()
		lat = ptr.Of(float32(val))
	}
	resp := api.Organization{
		ID:               organization.ID,
		Name:             organization.Name,
		AdoptionPolicy:   organization.AdoptionPolicy.Ptr(),
		AdoptionURL:      organization.AdoptionURL.Ptr(),
		MissionStatement: organization.MissionStatement.Ptr(),
		Distance:         organization.Distance.Ptr(),
		Website:          organization.Website.Ptr(),
		CreatedAt:        organization.CreatedAt,
		UpdatedAt:        organization.UpdatedAt,
	}
	if organization.R.OrganizationContact != nil {
		resp.Contact = api.Contact{
			Email:     organization.R.OrganizationContact.Email,
			Phone:     organization.R.OrganizationContact.Phone,
			CreatedAt: &organization.R.OrganizationContact.CreatedAt,
			UpdatedAt: &organization.R.OrganizationContact.UpdatedAt,
		}
	}
	if organization.R.OrganizationContact.R.Address != nil {
		resp.Contact.Address = api.Address{
			City:          organization.R.OrganizationContact.R.Address.City,
			ID:            organization.R.OrganizationContact.R.Address.ID,
			Lat:           lat,
			Lng:           lng,
			Note:          organization.R.OrganizationContact.R.Address.Note.Ptr(),
			PostalCode:    organization.R.OrganizationContact.R.Address.PostalCode.Ptr(),
			Region:        organization.R.OrganizationContact.R.Address.Region.Ptr(),
			StreetAddress: organization.R.OrganizationContact.R.Address.StreetNumber.MustGet(),
			StreetNumber:  organization.R.OrganizationContact.R.Address.StreetNumber.Ptr(),
			UnitNumber:    organization.R.OrganizationContact.R.Address.UnitNumber.Ptr(),
			CreatedAt:     organization.R.OrganizationContact.R.Address.CreatedAt,
			UpdatedAt:     organization.R.OrganizationContact.R.Address.UpdatedAt,
		}
	}
	if organization.R.OrganizationContact.R.Address.R.Country != nil {
		resp.Contact.Address.Country = api.Country{
			ID:         organization.R.OrganizationContact.R.Address.R.Country.ID,
			IsoAlpha2:  organization.R.OrganizationContact.R.Address.R.Country.IsoAlpha2,
			IsoAlpha3:  organization.R.OrganizationContact.R.Address.R.Country.IsoAlpha3,
			IsoNumeric: organization.R.OrganizationContact.R.Address.R.Country.IsoNumeric,
			Name:       organization.R.OrganizationContact.R.Address.R.Country.Name,
			CreatedAt:  organization.R.OrganizationContact.R.Address.R.Country.CreatedAt,
			UpdatedAt:  organization.R.OrganizationContact.R.Address.R.Country.UpdatedAt,
		}
	}
	if organization.R.OrganizationWorkHour != nil {
		resp.WorkHour = &api.OrganizationWorkHour{
			ID:             organization.R.OrganizationWorkHour.ID,
			OrganizationID: organization.R.OrganizationWorkHour.OrganizationID.MustGet(),
			Monday:         organization.R.OrganizationWorkHour.Monday.Ptr(),
			Tuesday:        organization.R.OrganizationWorkHour.Tuesday.Ptr(),
			Wednesday:      organization.R.OrganizationWorkHour.Wednesday.Ptr(),
			Thursday:       organization.R.OrganizationWorkHour.Thursday.Ptr(),
			Friday:         organization.R.OrganizationWorkHour.Friday.Ptr(),
			Saturday:       organization.R.OrganizationWorkHour.Saturday.Ptr(),
			Sunday:         organization.R.OrganizationWorkHour.Sunday.Ptr(),
			CreatedAt:      &organization.R.OrganizationWorkHour.CreatedAt,
			UpdatedAt:      &organization.R.OrganizationWorkHour.UpdatedAt,
		}
	}
	photos := []api.OrganizationPhoto{}
	for _, photo := range organization.R.OrganizationPhotos {
		photos = append(photos, api.OrganizationPhoto{
			ID:             photo.ID,
			OrganizationID: photo.OrganizationID.MustGet(),
			Small:          photo.Small.Ptr(),
			Medium:         photo.Medium.Ptr(),
			Large:          photo.Large.Ptr(),
			Full:           photo.Full.Ptr(),
			CreatedAt:      &photo.CreatedAt,
			UpdatedAt:      &photo.UpdatedAt,
		})
	}
	resp.Photos = photos
	socials := []api.OrganizationSocial{}
	for _, social := range organization.R.OrganizationSocials {
		socials = append(socials, api.OrganizationSocial{
			Platform: social.Platform,
			URL:      social.URL,
		})
	}
	resp.Socials = socials
	return resp
}
