package dto

import (
	api "github.com/dankobg/fluffly/api/gen"
	"github.com/dankobg/fluffly/db/gen/test/public/model"
)

func AnimalBreedToResp(data model.Breed) api.Breed {
	resp := api.Breed{
		ID:           data.ID,
		AnimalTypeID: data.AnimalSpeciesID,
		Name:         data.Name,
		Primary:      false,
		CreatedAt:    data.CreatedAt,
		UpdatedAt:    data.UpdatedAt,
	}
	return resp

}
