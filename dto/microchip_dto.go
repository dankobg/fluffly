package dto

import (
	api "github.com/dankobg/fluffly/api/gen"
	"github.com/dankobg/fluffly/db/gen/models"
)

func MicrochipToResponse(data models.Microchip) api.Microchip {
	resp := api.Microchip{
		ID:          data.ID,
		Number:      data.Number,
		Brand:       data.Brand.Ptr(),
		Description: data.Description.Ptr(),
		Location:    data.Location.Ptr(),
		CreatedAt:   data.CreatedAt,
		UpdatedAt:   data.UpdatedAt,
	}

	return resp
}
