package dto

import (
	api "github.com/dankobg/fluffly/api/gen"
	"github.com/dankobg/fluffly/db/gen/models"
	"github.com/dankobg/fluffly/media"
)

func OrganizationPhotoToResponse(data models.OrganizationPhoto, upl media.Uploader) api.OrganizationPhoto {
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

	resp := api.OrganizationPhoto{
		ID:             data.ID,
		OrganizationID: data.OrganizationID.GetOrZero(),
		SmallURL:       getURL(data.ObjectRefSmall.Ptr(), data.ObjectKind, upl),
		MediumURL:      getURL(data.ObjectRefMedium.Ptr(), data.ObjectKind, upl),
		LargeURL:       getURL(data.ObjectRefLarge.Ptr(), data.ObjectKind, upl),
		FullURL:        getURL(data.ObjectRefFull.Ptr(), data.ObjectKind, upl),
		CreatedAt:      data.CreatedAt,
		UpdatedAt:      data.UpdatedAt,
	}

	return resp
}
