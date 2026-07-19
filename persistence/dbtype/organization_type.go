package dbtype

import (
	"github.com/aarondl/opt/omitnull"
	api "github.com/dankobg/fluffly/api/gen"
	"github.com/dankobg/fluffly/db/gen/models"
	"github.com/stephenafamo/bob/types"
)

type ListMyOrganizationsFilters struct {
	api.ListMyOrganizationsParams
}

type ListOrganizationsFilters struct {
	api.ListOrganizationsParams
}

type GetOrganizationByIDFilters struct {
	api.GetOrganizationParams
}

type OrganizationWithJoinData struct {
	models.Organization
	WorkHour              *models.OrganizationWorkHour
	Contact               *models.OrganizationContact
	ContactAddress        *models.Address
	ContactAddressCountry *models.Country
	Photos                types.JSON[*[]models.OrganizationPhoto]
	Videos                types.JSON[*[]models.OrganizationVideo]
	Socials               types.JSON[*[]models.OrganizationSocial]
}

type OrganizationApplyForSetter struct {
	Organization models.OrganizationSetter
	Contact      models.OrganizationContactSetter
	Address      models.AddressSetter
	WorkHour     omitnull.Val[models.OrganizationWorkHourSetter]
	Photos       omitnull.Val[[]models.OrganizationPhotoSetter]
	Videos       omitnull.Val[[]models.OrganizationVideoSetter]
	Socials      omitnull.Val[[]models.OrganizationSocialSetter]
}

type OrganizationCreateSetter struct {
	Organization models.OrganizationSetter
	Contact      models.OrganizationContactSetter
	Address      models.AddressSetter
	WorkHour     omitnull.Val[models.OrganizationWorkHourSetter]
	Photos       omitnull.Val[[]models.OrganizationPhotoSetter]
	Videos       omitnull.Val[[]models.OrganizationVideoSetter]
	Socials      omitnull.Val[[]models.OrganizationSocialSetter]
}

type OrganizationUpdateSetter struct {
	Organization omitnull.Val[models.OrganizationSetter]
	Contact      omitnull.Val[models.OrganizationContactSetter]
	Address      omitnull.Val[models.AddressSetter]
	WorkHour     omitnull.Val[models.OrganizationWorkHourSetter]
	Photos       omitnull.Val[[]models.OrganizationPhotoSetter]
	Videos       omitnull.Val[[]models.OrganizationVideoSetter]
	Socials      omitnull.Val[[]models.OrganizationSocialSetter]
}
