package dto

import (
	api "github.com/dankobg/fluffly/api/gen"
	"github.com/dankobg/fluffly/db/gen/models"
	"github.com/dankobg/fluffly/media"
)

func AnimalVideoToResponse(data models.AnimalVideo, upl media.Uploader) api.AnimalVideo {
	getURL := func(name, kind string, upl media.Uploader) *string {
		u, err := upl.URL(name, kind)
		if err != nil {
			return nil
		}

		return &u
	}

	resp := api.AnimalVideo{
		ID:        data.ID,
		AnimalID:  data.AnimalID.GetOrZero(),
		URL:       getURL(data.ObjectRef, data.ObjectKind, upl),
		CreatedAt: data.CreatedAt,
		UpdatedAt: data.UpdatedAt,
	}

	return resp
}
