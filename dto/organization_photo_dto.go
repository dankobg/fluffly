package dto

import (
	api "github.com/dankobg/fluffly/api/gen"
	"github.com/dankobg/fluffly/db/gen/test/public/model"
)

func OrganizationPhotoToResp(data model.OrganizationPhoto) api.OrganizationPhoto {
	resp := api.OrganizationPhoto{
		ID:        data.ID,
		Small:     data.Small,
		Medium:    data.Medium,
		Large:     data.Large,
		Full:      data.Full,
		CreatedAt: data.CreatedAt,
		UpdatedAt: data.UpdatedAt,
	}
	if data.OrganizationID != nil {
		resp.OrganizationID = *data.OrganizationID
	}
	return resp

}
