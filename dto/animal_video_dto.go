package dto

import (
	api "github.com/dankobg/fluffly/api/gen"
	"github.com/dankobg/fluffly/db/gen/test/public/model"
)

func AnimalVideoToResp(data model.AnimalVideo) api.AnimalVideo {
	resp := api.AnimalVideo{
		ID:        data.ID,
		URL:       &data.URL,
		CreatedAt: data.CreatedAt,
		UpdatedAt: data.UpdatedAt,
	}
	if data.AnimalID != nil {
		resp.AnimalID = *data.AnimalID
	}
	return resp

}
