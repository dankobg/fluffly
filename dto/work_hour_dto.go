package dto

import (
	api "github.com/dankobg/fluffly/api/gen"
	"github.com/dankobg/fluffly/db/gen/models"
)

func WorkHourToResponse(data models.OrganizationWorkHour) api.OrganizationWorkHour {
	resp := api.OrganizationWorkHour{
		ID:             data.ID,
		OrganizationID: data.OrganizationID.GetOrZero(),
		Monday:         data.Monday.Ptr(),
		Tuesday:        data.Tuesday.Ptr(),
		Wednesday:      data.Wednesday.Ptr(),
		Thursday:       data.Thursday.Ptr(),
		Friday:         data.Friday.Ptr(),
		Saturday:       data.Saturday.Ptr(),
		Sunday:         data.Sunday.Ptr(),
		UpdatedAt:      data.CreatedAt,
		CreatedAt:      data.UpdatedAt,
	}

	return resp
}
