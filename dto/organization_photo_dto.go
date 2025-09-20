package dto

import (
	api "github.com/dankobg/fluffly/api/gen"
	"github.com/dankobg/fluffly/db/gen/test/public/model"
	"github.com/dankobg/fluffly/media"
)

func OrganizationPhotoToResp(data model.OrganizationPhoto, upl media.Uploader) api.OrganizationPhoto {
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
		ID:        data.ID,
		SmallURL:  getURL(data.ObjectRefSmall, data.ObjectKind, upl),
		MediumURL: getURL(data.ObjectRefMedium, data.ObjectKind, upl),
		LargeURL:  getURL(data.ObjectRefLarge, data.ObjectKind, upl),
		FullURL:   getURL(data.ObjectRefFull, data.ObjectKind, upl),
		CreatedAt: data.CreatedAt,
		UpdatedAt: data.UpdatedAt,
	}
	if data.OrganizationID != nil {
		resp.OrganizationID = *data.OrganizationID
	}
	return resp
}
