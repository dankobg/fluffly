package dto

import (
	api "github.com/dankobg/fluffly/api/gen"
	"github.com/dankobg/fluffly/db/gen/test/public/model"
)

func WorkHourToResponse(data model.OrganizationWorkHour) api.OrganizationWorkHour {
	resp := api.OrganizationWorkHour{
		ID:        data.ID,
		Monday:    data.Monday,
		Tuesday:   data.Tuesday,
		Wednesday: data.Wednesday,
		Thursday:  data.Thursday,
		Friday:    data.Friday,
		Saturday:  data.Saturday,
		Sunday:    data.Sunday,
		UpdatedAt: data.CreatedAt,
		CreatedAt: data.UpdatedAt,
	}
	if data.OrganizationID != nil {
		resp.OrganizationID = *data.OrganizationID
	}
	return resp
}
