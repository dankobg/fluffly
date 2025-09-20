package dto

import (
	api "github.com/dankobg/fluffly/api/gen"
	"github.com/dankobg/fluffly/db/gen/test/public/model"
	"github.com/dankobg/fluffly/media"
)

func OrganizationVideoToResp(data model.OrganizationVideo, upl media.Uploader) api.OrganizationVideo {
	getURL := func(name, kind string, upl media.Uploader) *string {
		u, err := upl.URL(name, kind)
		if err != nil {
			return nil
		}
		return &u
	}

	resp := api.OrganizationVideo{
		ID:        data.ID,
		URL:       getURL(data.ObjectRef, data.ObjectKind, upl),
		CreatedAt: data.CreatedAt,
		UpdatedAt: data.UpdatedAt,
	}
	if data.OrganizationID != nil {
		resp.OrganizationID = *data.OrganizationID
	}
	return resp

}
