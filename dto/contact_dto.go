package dto

import (
	api "github.com/dankobg/fluffly/api/gen"
	"github.com/dankobg/fluffly/db/gen/models"
)

func CountryToResponse(data models.Country) api.Country {
	return api.Country{
		ID:         data.ID,
		Name:       data.Name,
		IsoAlpha2:  data.IsoAlpha2,
		IsoAlpha3:  data.IsoAlpha3,
		IsoNumeric: data.IsoNumeric,
		CreatedAt:  data.CreatedAt,
		UpdatedAt:  data.UpdatedAt,
	}
}

func AddressToResponse(data models.Address, country models.Country) api.Address {
	var lat, lon *float64

	if data.Coords.IsValue() {
		flatCoords := data.Coords.MustGet().FlatCoords()
		if len(flatCoords) == 2 {
			lon = &flatCoords[0]
			lat = &flatCoords[1]
		}
	}

	resp := api.Address{
		ID:            data.ID,
		City:          data.City,
		Country:       CountryToResponse(country),
		Note:          data.Note.Ptr(),
		PostalCode:    data.PostalCode.Ptr(),
		Region:        data.Region.Ptr(),
		StreetAddress: data.StreetAddress,
		StreetNumber:  data.StreetNumber.Ptr(),
		UnitNumber:    data.UnitNumber.Ptr(),
		Lat:           lat,
		Lon:           lon,
		CreatedAt:     data.CreatedAt,
		UpdatedAt:     data.UpdatedAt,
	}

	return resp
}

func ContactToResponse(data models.OrganizationContact, addr models.Address, country models.Country) api.Contact {
	resp := api.Contact{
		Address:   AddressToResponse(addr, country),
		Email:     data.Email,
		Phone:     data.Phone,
		CreatedAt: &data.CreatedAt,
		UpdatedAt: &data.UpdatedAt,
	}

	return resp
}
