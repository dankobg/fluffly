package dto

import (
	api "github.com/dankobg/fluffly/api/gen"
	"github.com/dankobg/fluffly/db/gen/models"
	"github.com/dankobg/fluffly/media"
	"github.com/dankobg/fluffly/persistence/dbtype"
)

func AdoptionToResponse(data models.Adoption) api.Adoption {
	return api.Adoption{
		ID:             data.ID,
		AnimalID:       data.AnimalID,
		UserID:         data.UserID,
		OrganizationID: data.OrganizationID.Ptr(),
		IsPermanent:    data.IsPermanent,
		Status:         api.AdoptionStatus(data.Status),
		Note:           data.Note.Ptr(),
		ReturnedAt:     data.ReturnedAt.Ptr(),
		AdoptedAt:      &data.AdoptedAt,
		CreatedAt:      data.CreatedAt,
		UpdatedAt:      data.UpdatedAt,
	}
}

func AdoptionWithJoinDataToResponse(data dbtype.AdoptionWithJoinData, upl media.Uploader) api.Adoption {
	resp := AdoptionToResponse(data.Adoption)

	// if data.Organization != nil && data.Organization.ID != 0 {
	// 	resp.Organization = new(OrganizationToResponse(*data.Organization))
	// }

	if data.Organization != nil && data.Organization.ID != 0 {
		resp.Organization = new(OrganizationWithJoinDataToResponse(dbtype.OrganizationWithJoinData{
			Organization:          *data.Organization,
			WorkHour:              data.OrganizationWorkHour,
			Contact:               data.OrganizationContact,
			ContactAddress:        data.OrganizationContactAddress,
			ContactAddressCountry: data.OrganizationContactAddressCountry,
		}, upl))
	}

	return resp
}
