package dto

import (
	api "github.com/dankobg/fluffly/api/gen"
	"github.com/dankobg/fluffly/db/gen/models"
	"github.com/dankobg/fluffly/persistence/dbcustom"
)

func AnimalSpeciesToResponse(data models.AnimalSpecie) api.AnimalSpecie {
	var propertiesSchema *dbcustom.PropertiesSchema

	jsonPropertiesSchema := data.PropertiesSchema.Ptr()
	if jsonPropertiesSchema != nil {
		propertiesSchema = &jsonPropertiesSchema.Val
	}

	resp := api.AnimalSpecie{
		ID:               data.ID,
		AnimalTypeID:     data.AnimalTypeID,
		Name:             data.Name,
		PropertiesSchema: propertiesSchema,
		CreatedAt:        data.CreatedAt,
		UpdatedAt:        data.UpdatedAt,
	}

	return resp
}

func AnimalSpeciesMinToResponse(data models.AnimalSpecie) api.AnimalSpecieMin {
	resp := api.AnimalSpecieMin{
		ID:           data.ID,
		AnimalTypeID: data.AnimalTypeID,
		Name:         data.Name,
		CreatedAt:    data.CreatedAt,
		UpdatedAt:    data.UpdatedAt,
	}

	return resp
}
