package dto

import (
	api "github.com/dankobg/fluffly/api/gen"
	"github.com/dankobg/fluffly/persistence"
)

func AnimalBreedWithJounDataToResp(data persistence.AnimalBreedWithJoinData) api.Breed {
	resp := api.Breed{
		ID:           data.ID,
		AnimalTypeID: data.AnimalSpeciesID,
		Name:         data.Name,
		Primary:      data.Primary,
		CreatedAt:    data.CreatedAt,
		UpdatedAt:    data.UpdatedAt,
	}
	return resp

}
