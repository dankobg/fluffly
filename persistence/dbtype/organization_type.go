package dbtype

import (
	"time"

	api "github.com/dankobg/fluffly/api/gen"
	"github.com/dankobg/fluffly/db/gen/test/public/model"
	t "github.com/dankobg/fluffly/db/gen/test/public/table"
	"github.com/dankobg/fluffly/ptr"
	p "github.com/go-jet/jet/v2/postgres"
	"github.com/oapi-codegen/nullable"
)

type OrganizationFilters struct {
	Pagination *api.PaginationParams
}

type OrganizationWithJoinData struct {
	model.Organization
	WorkHour model.OrganizationWorkHour
	Contact  struct {
		model.OrganizationContact
		Address struct {
			model.Address
			Country model.Country
		}
	}
	Photos  []model.OrganizationPhoto  `json_column:"photos"`
	Videos  []model.OrganizationVideo  `json_column:"videos"`
	Socials []model.OrganizationSocial `json_column:"socials"`
}

type OrganizationCreateSetter struct {
	Organization OrganizationSetter
	Contact      OrganizationContactSetter
	Address      AddressSetter
	WorkHour     nullable.Nullable[OrganizationWorkHourSetter]
	Photos       nullable.Nullable[[]OrganizationPhotoSetter]
	Videos       nullable.Nullable[[]OrganizationVideoSetter]
	Socials      nullable.Nullable[[]OrganizationSocialSetter]
}

type OrganizationSetter struct {
	Name             nullable.Nullable[string]    `json:"name"`
	Website          nullable.Nullable[string]    `json:"website"`
	MissionStatement nullable.Nullable[string]    `json:"mission_statement"`
	AdoptionPolicy   nullable.Nullable[string]    `json:"adoption_policy"`
	AdoptionURL      nullable.Nullable[string]    `json:"adoption_url"`
	Distance         nullable.Nullable[string]    `json:"distance"`
	CreatedAt        nullable.Nullable[time.Time] `json:"created_at"`
	UpdatedAt        nullable.Nullable[time.Time] `json:"updated_at"`
}

func (s OrganizationSetter) ToModel(isPatch ...bool) (p.ColumnList, model.Organization) {
	var cols p.ColumnList
	var m model.Organization

	if len(isPatch) > 0 {
		if isPatch[0] {
			cols = append(cols, t.Organization.UpdatedAt)
			m.UpdatedAt = time.Now()
		}
	}

	if s.Name.IsSpecified() {
		cols = append(cols, t.Organization.Name)
		if !s.Name.IsNull() {
			m.Name = s.Name.MustGet()
		}
	}
	if s.Website.IsSpecified() {
		cols = append(cols, t.Organization.Website)
		if !s.Website.IsNull() {
			m.Website = ptr.Of(s.Website.MustGet())
		} else {
			m.Website = nil
		}
	}
	if s.MissionStatement.IsSpecified() {
		cols = append(cols, t.Organization.MissionStatement)
		if !s.MissionStatement.IsNull() {
			m.MissionStatement = ptr.Of(s.MissionStatement.MustGet())
		} else {
			m.MissionStatement = nil
		}
	}
	if s.AdoptionPolicy.IsSpecified() {
		cols = append(cols, t.Organization.AdoptionPolicy)
		if !s.AdoptionPolicy.IsNull() {
			m.AdoptionPolicy = ptr.Of(s.AdoptionPolicy.MustGet())
		} else {
			m.AdoptionPolicy = nil
		}
	}
	if s.AdoptionURL.IsSpecified() {
		cols = append(cols, t.Organization.AdoptionURL)
		if !s.AdoptionURL.IsNull() {
			m.AdoptionURL = ptr.Of(s.AdoptionURL.MustGet())
		} else {
			m.AdoptionURL = nil
		}
	}
	if s.Distance.IsSpecified() {
		cols = append(cols, t.Organization.Distance)
		if !s.Distance.IsNull() {
			m.Distance = ptr.Of(s.Distance.MustGet())
		} else {
			m.Distance = nil
		}
	}
	if s.CreatedAt.IsSpecified() {
		cols = append(cols, t.Organization.CreatedAt)
		if !s.CreatedAt.IsNull() {
			m.CreatedAt = s.CreatedAt.MustGet()
		}
	}
	if s.UpdatedAt.IsSpecified() {
		cols = append(cols, t.Organization.UpdatedAt)
		if !s.UpdatedAt.IsNull() {
			m.UpdatedAt = s.UpdatedAt.MustGet()
		}
	}

	return cols, m
}
