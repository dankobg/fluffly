package dto

import (
	api "github.com/dankobg/fluffly/api/gen"
	"github.com/dankobg/fluffly/db/gen/models"
	"github.com/dankobg/fluffly/media"
)

func AnimalPhotoToResponse(data models.AnimalPhoto, upl media.Uploader) api.AnimalPhoto {
	getURL := func(name *string, kind string, upl media.Uploader) *string {
		if name == nil {
			return nil
		}

		u, err := upl.URL(*name, kind)
		if err != nil {
			return nil
		}

		return &u
	}

	resp := api.AnimalPhoto{
		ID:        data.ID,
		AnimalID:  data.AnimalID.GetOrZero(),
		SmallURL:  getURL(data.ObjectRefSmall.Ptr(), data.ObjectKind, upl),
		MediumURL: getURL(data.ObjectRefMedium.Ptr(), data.ObjectKind, upl),
		LargeURL:  getURL(data.ObjectRefLarge.Ptr(), data.ObjectKind, upl),
		FullURL:   getURL(data.ObjectRefFull.Ptr(), data.ObjectKind, upl),
		CreatedAt: data.CreatedAt,
		UpdatedAt: data.UpdatedAt,
	}

	return resp
}
