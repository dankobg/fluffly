package dto

import (
	api "github.com/dankobg/fluffly/api/gen"
	"github.com/dankobg/fluffly/db/model"
)

func OrganizationToResponse(organization model.Organization) api.Organization {
	return api.Organization{
		ID:               organization.ID,
		ContactID:        organization.ContactID,
		Name:             organization.Name,
		AdoptionPolicy:   organization.AdoptionPolicy.Ptr(),
		AdoptionURL:      organization.AdoptionURL.Ptr(),
		MissionStatement: organization.MissionStatement.Ptr(),
		Distance:         organization.Distance.Ptr(),
		Website:          organization.Website.Ptr(),
		Facebook:         organization.Facebook.Ptr(),
		Instagram:        organization.Instagram.Ptr(),
		Pinterest:        organization.Pinterest.Ptr(),
		Twitter:          organization.Twitter.Ptr(),
		Youtube:          organization.Youtube.Ptr(),
		CreatedAt:        organization.CreatedAt,
		UpdatedAt:        organization.UpdatedAt,
	}
}

// func OrganizationToResponse(organization kratos.Organization) (api.Organization, error) {
// 	dispatches := make([]api.OrganizationDispatch, 0, len(organization.Dispatches))
// 	for _, d := range organization.Dispatches {
// 		id, err := uuid.Parse(organization.Id)
// 		if err != nil {
// 			return api.Organization{}, fmt.Errorf("failed to parse dispatch uuid: %w", err)
// 		}
// 		organizationID, err := uuid.Parse(organization.Id)
// 		if err != nil {
// 			return api.Organization{}, fmt.Errorf("failed to parse dispatch organizationid uuid: %w", err)
// 		}
// 		dispatches = append(dispatches, api.OrganizationDispatch{
// 			ID:        id,
// 			OrganizationID: organizationID,
// 			Status:    api.OrganizationDispatchStatus(d.Status),
// 			Error:     &d.Error,
// 			CreatedAt: d.CreatedAt,
// 			UpdatedAt: d.UpdatedAt,
// 		})
// 	}
// 	id, err := uuid.Parse(organization.Id)
// 	if err != nil {
// 		return api.Organization{}, fmt.Errorf("failed to parse organization uuid: %w", err)
// 	}
// 	resp := api.Organization{
// 		ID:           id,
// 		Body:         organization.Body,
// 		Subject:      organization.Subject,
// 		Channel:      organization.Channel,
// 		Recipient:    organization.Recipient,
// 		Status:       api.CourierOrganizationStatus(organization.Status),
// 		TemplateType: api.CourierOrganizationTemplateType(organization.TemplateType),
// 		Type:         api.CourierOrganizationType(organization.Type),
// 		SendCount:    organization.SendCount,
// 		Dispatches:   &dispatches,
// 		CreatedAt:    organization.CreatedAt,
// 		UpdatedAt:    organization.UpdatedAt,
// 	}
// 	return resp, nil
// }
