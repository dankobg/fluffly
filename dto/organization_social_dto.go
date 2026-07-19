package dto

import (
	api "github.com/dankobg/fluffly/api/gen"
	"github.com/dankobg/fluffly/db/gen/models"
)

func OrganizationSocialToResponse(data models.OrganizationSocial) api.OrganizationSocial {
	return api.OrganizationSocial{
		ID:        data.ID,
		Platform:  data.Platform,
		URL:       data.URL,
		CreatedAt: data.CreatedAt,
		UpdatedAt: data.UpdatedAt,
	}
}
