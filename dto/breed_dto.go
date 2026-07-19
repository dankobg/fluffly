package dto

import (
	api "github.com/dankobg/fluffly/api/gen"
	"github.com/dankobg/fluffly/db/gen/models"
)

func BreedToResponse(data models.Breed) api.Breed {
	resp := api.Breed{
		ID:             data.ID,
		AnimalSpecieID: data.AnimalSpecieID,
		Name:           data.Name,
		CreatedAt:      data.CreatedAt,
		UpdatedAt:      data.UpdatedAt,
	}

	return resp
}
