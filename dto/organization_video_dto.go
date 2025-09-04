package dto

import (
	api "github.com/dankobg/fluffly/api/gen"
	"github.com/dankobg/fluffly/db/gen/test/public/model"
)

func OrganizationVideoToResp(data model.OrganizationVideo) api.OrganizationVideo {
	resp := api.OrganizationVideo{
		ID:        data.ID,
		URL:       &data.URL,
		CreatedAt: data.CreatedAt,
		UpdatedAt: data.UpdatedAt,
	}
	if data.OrganizationID != nil {
		resp.OrganizationID = *data.OrganizationID
	}
	return resp

}
