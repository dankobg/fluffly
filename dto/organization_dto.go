package dto

import (
	api "github.com/dankobg/fluffly/api/gen"
	"github.com/dankobg/fluffly/db/gen/models"
	"github.com/dankobg/fluffly/media"
	"github.com/dankobg/fluffly/persistence/dbtype"
)

func OrganizationToResponse(data models.Organization) api.Organization {
	return api.Organization{
		ID:               data.ID,
		Name:             data.Name,
		AdoptionPolicy:   data.AdoptionPolicy.Ptr(),
		AdoptionURL:      data.AdoptionURL.Ptr(),
		Distance:         data.Distance.Ptr(),
		MissionStatement: data.MissionStatement.Ptr(),
		Website:          data.Website.Ptr(),
		Status:           api.OrganizationStatus(data.Status),
		CreatedAt:        data.CreatedAt,
		UpdatedAt:        data.UpdatedAt,
	}
}

func OrganizationWithJoinDataToResponse(data dbtype.OrganizationWithJoinData, upl media.Uploader) api.Organization {
	resp := OrganizationToResponse(data.Organization)

	if data.Contact != nil && data.Contact.ID != 0 {
		resp.Contact = new(ContactToResponse(*data.Contact, *data.ContactAddress, *data.ContactAddressCountry))
	}

	if data.WorkHour != nil && data.WorkHour.ID != 0 {
		resp.WorkHour = new(WorkHourToResponse(*data.WorkHour))
	}

	if data.Photos.Val != nil {
		photos := make([]api.OrganizationPhoto, len(*data.Photos.Val))
		for i, photo := range *data.Photos.Val {
			photos[i] = OrganizationPhotoToResponse(photo, upl)
		}

		resp.Photos = &photos
	}

	if data.Videos.Val != nil {
		videos := make([]api.OrganizationVideo, len(*data.Videos.Val))
		for i, video := range *data.Videos.Val {
			videos[i] = OrganizationVideoToResponse(video, upl)
		}

		resp.Videos = &videos
	}

	if data.Socials.Val != nil {
		socials := make([]api.OrganizationSocial, len(*data.Socials.Val))
		for i, social := range *data.Socials.Val {
			socials[i] = OrganizationSocialToResponse(social)
		}

		resp.Socials = &socials
	}

	return resp
}
