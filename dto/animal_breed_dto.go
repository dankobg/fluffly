package dto

import (
	api "github.com/dankobg/fluffly/api/gen"
	"github.com/dankobg/fluffly/db/gen/models"
	"github.com/dankobg/fluffly/persistence/dbtype"
)

func AnimalBreedToResponse(data models.AnimalBreed) api.AnimalBreed {
	resp := api.AnimalBreed{
		// @TODO: too lazy, stub for now
		Primary:   data.Primary,
		CreatedAt: data.CreatedAt,
		UpdatedAt: data.UpdatedAt,
	}

	return resp
}

func AnimalBreedWithJoinDataToResponse(data dbtype.AnimalBreedWithJoinData) api.AnimalBreed {
	resp := api.AnimalBreed{
		ID:             data.ID,
		AnimalSpecieID: data.AnimalSpecieID,
		Name:           data.Name,
		Primary:        data.Primary,
		CreatedAt:      data.CreatedAt,
		UpdatedAt:      data.UpdatedAt,
	}

	return resp
}
