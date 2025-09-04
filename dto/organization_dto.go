package dto

import (
	api "github.com/dankobg/fluffly/api/gen"
	"github.com/dankobg/fluffly/db/gen/test/public/model"
	"github.com/dankobg/fluffly/persistence"
)

func OrganizationToResponse(data model.Organization) api.Organization {
	return api.Organization{
		ID:               data.ID,
		Name:             data.Name,
		AdoptionPolicy:   data.AdoptionPolicy,
		AdoptionURL:      data.AdoptionURL,
		Distance:         data.Distance,
		MissionStatement: data.MissionStatement,
		Website:          data.Website,
		CreatedAt:        data.CreatedAt,
		UpdatedAt:        data.UpdatedAt,
	}
}

func OrganizationWithJoinDataToResponse(data persistence.OrganizationWithJoinData) api.Organization {
	resp := OrganizationToResponse(data.Organization)
	resp.Contact = ContactToResponse(data.Contact.OrganizationContact, data.Contact.Address.Address, data.Contact.Address.Country)
	workHour := WorkHourToResponse(data.WorkHour)
	resp.WorkHour = &workHour
	resp.Photos = make([]api.OrganizationPhoto, len(data.Photos))
	for i, photo := range data.Photos {
		resp.Photos[i] = OrganizationPhotoToResp(photo)
	}
	for i, video := range data.Videos {
		resp.Videos[i] = OrganizationVideoToResp(video)
	}
	resp.Socials = make([]api.OrganizationSocial, len(data.Socials))
	for i, social := range data.Socials {
		resp.Socials[i] = OrganizationSocialToResp(social)
	}
	return resp
}
