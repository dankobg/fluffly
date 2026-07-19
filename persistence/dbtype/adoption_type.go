package dbtype

import (
	api "github.com/dankobg/fluffly/api/gen"
	"github.com/dankobg/fluffly/db/gen/models"
)

type ListAdoptionsFilters struct {
	api.ListAdoptionsParams
}

type GetAdoptionByIDFilters struct {
	api.GetAdoptionParams
}

type AdoptionWithJoinData struct {
	models.Adoption
	Organization                      *models.Organization
	OrganizationWorkHour              *models.OrganizationWorkHour
	OrganizationContact               *models.OrganizationContact
	OrganizationContactAddress        *models.Address
	OrganizationContactAddressCountry *models.Country
	User                              *models.User
}
