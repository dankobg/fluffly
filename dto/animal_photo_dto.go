package dto

import (
	api "github.com/dankobg/fluffly/api/gen"
	"github.com/dankobg/fluffly/db/gen/test/public/model"
)

func AnimalPhotoToResp(data model.AnimalPhoto) api.AnimalPhoto {
	resp := api.AnimalPhoto{
		ID:        data.ID,
		Small:     data.Small,
		Medium:    data.Medium,
		Large:     data.Large,
		Full:      data.Full,
		CreatedAt: data.CreatedAt,
		UpdatedAt: data.UpdatedAt,
	}
	if data.AnimalID != nil {
		resp.AnimalID = *data.AnimalID
	}
	return resp

}
