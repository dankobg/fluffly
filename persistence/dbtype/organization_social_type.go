package dbtype

import (
	"time"

	"github.com/dankobg/fluffly/db/gen/test/public/model"
	t "github.com/dankobg/fluffly/db/gen/test/public/table"
	p "github.com/go-jet/jet/v2/postgres"
	"github.com/oapi-codegen/nullable"
)

type OrganizationSocialSetter struct {
	OrganizationID nullable.Nullable[int64]     `json:"organizationid"`
	Platform       nullable.Nullable[string]    `json:"platform"`
	URL            nullable.Nullable[string]    `json:"url"`
	CreatedAt      nullable.Nullable[time.Time] `json:"created_at"`
	UpdatedAt      nullable.Nullable[time.Time] `json:"updated_at"`
}

func (s OrganizationSocialSetter) ToModel(isPatch ...bool) (p.ColumnList, model.OrganizationSocial) {
	var cols p.ColumnList
	var m model.OrganizationSocial

	if len(isPatch) > 0 {
		if isPatch[0] {
			cols = append(cols, t.OrganizationSocial.UpdatedAt)
			m.UpdatedAt = time.Now()
		}
	}

	if s.OrganizationID.IsSpecified() {
		cols = append(cols, t.OrganizationSocial.OrganizationID)
		if !s.OrganizationID.IsNull() {
			m.OrganizationID = s.OrganizationID.MustGet()
		}
	}
	if s.Platform.IsSpecified() {
		cols = append(cols, t.OrganizationSocial.Platform)
		if !s.Platform.IsNull() {
			m.Platform = s.Platform.MustGet()
		}
	}
	if s.URL.IsSpecified() {
		cols = append(cols, t.OrganizationSocial.URL)
		if !s.URL.IsNull() {
			m.URL = s.URL.MustGet()
		}
	}
	if s.CreatedAt.IsSpecified() {
		cols = append(cols, t.OrganizationSocial.CreatedAt)
		if !s.CreatedAt.IsNull() {
			m.CreatedAt = s.CreatedAt.MustGet()
		}
	}
	if s.UpdatedAt.IsSpecified() {
		cols = append(cols, t.OrganizationSocial.UpdatedAt)
		if !s.UpdatedAt.IsNull() {
			m.UpdatedAt = s.UpdatedAt.MustGet()
		}
	}

	return cols, m
}
