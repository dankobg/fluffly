package dto

import (
	api "github.com/dankobg/fluffly/api/gen"
	"github.com/dankobg/fluffly/db/gen/test/public/model"
	"github.com/dankobg/fluffly/media"
)

func AnimalVideoToResp(data model.AnimalVideo, upl media.Uploader) api.AnimalVideo {
	getURL := func(name, kind string, upl media.Uploader) *string {
		u, err := upl.URL(name, kind)
		if err != nil {
			return nil
		}
		return &u
	}

	resp := api.AnimalVideo{
		ID:        data.ID,
		URL:       getURL(data.ObjectRef, data.ObjectKind, upl),
		CreatedAt: data.CreatedAt,
		UpdatedAt: data.UpdatedAt,
	}
	if data.AnimalID != nil {
		resp.AnimalID = *data.AnimalID
	}
	return resp

}
