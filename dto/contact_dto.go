package dto

import (
	api "github.com/dankobg/fluffly/api/gen"
	"github.com/dankobg/fluffly/db/gen/test/public/model"
	"github.com/dankobg/fluffly/ptr"
)

func CountryToResponse(data model.Country) api.Country {
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

func AddressToResponse(data model.Address, country model.Country) api.Address {
	resp := api.Address{
		ID:            data.ID,
		City:          data.City,
		Country:       CountryToResponse(country),
		Note:          data.Note,
		PostalCode:    data.PostalCode,
		Region:        data.Region,
		StreetAddress: data.StreetAddress,
		StreetNumber:  data.StreetNumber,
		UnitNumber:    data.UnitNumber,
		CreatedAt:     data.CreatedAt,
		UpdatedAt:     data.UpdatedAt,
	}
	if data.Lat != nil {
		resp.Lat = ptr.Of(float32(*data.Lat))
	}
	if data.Lng != nil {
		resp.Lng = ptr.Of(float32(*data.Lng))
	}
	return resp
}

func ContactToResponse(data model.OrganizationContact, addr model.Address, country model.Country) api.Contact {
	resp := api.Contact{
		Address:   AddressToResponse(addr, country),
		Email:     data.Email,
		Phone:     data.Phone,
		CreatedAt: &data.CreatedAt,
		UpdatedAt: &data.UpdatedAt,
	}
	return resp
}
