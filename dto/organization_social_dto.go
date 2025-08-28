package dto

import (
	api "github.com/dankobg/fluffly/api/gen"
	"github.com/dankobg/fluffly/db/gen/test/public/model"
)

func OrganizationSocialToResp(data model.OrganizationSocial) api.OrganizationSocial {
	return api.OrganizationSocial{
		Platform:  data.Platform,
		URL:       data.URL,
		CreatedAt: data.CreatedAt,
		UpdatedAt: data.UpdatedAt,
	}
}
