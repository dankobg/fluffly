package dbtype

import (
	"time"

	"github.com/dankobg/fluffly/db/gen/test/public/model"
	t "github.com/dankobg/fluffly/db/gen/test/public/table"
	"github.com/dankobg/fluffly/ptr"
	p "github.com/go-jet/jet/v2/postgres"
	"github.com/oapi-codegen/nullable"
)

type OrganizationWorkHourSetter struct {
	OrganizationID nullable.Nullable[int64]     `json:"organizationID"`
	Monday         nullable.Nullable[string]    `json:"monday"`
	Tuesday        nullable.Nullable[string]    `json:"tuesday"`
	Wednesday      nullable.Nullable[string]    `json:"wednesday"`
	Thursday       nullable.Nullable[string]    `json:"thursday"`
	Friday         nullable.Nullable[string]    `json:"friday"`
	Saturday       nullable.Nullable[string]    `json:"saturday"`
	Sunday         nullable.Nullable[string]    `json:"sunday"`
	CreatedAt      nullable.Nullable[time.Time] `json:"created_at"`
	UpdatedAt      nullable.Nullable[time.Time] `json:"updated_at"`
}

func (s OrganizationWorkHourSetter) ToModel(isPatch ...bool) (p.ColumnList, model.OrganizationWorkHour) {
	var cols p.ColumnList
	var m model.OrganizationWorkHour

	if len(isPatch) > 0 {
		if isPatch[0] {
			cols = append(cols, t.OrganizationWorkHour.UpdatedAt)
			m.UpdatedAt = time.Now()
		}
	}

	if s.OrganizationID.IsSpecified() {
		cols = append(cols, t.OrganizationWorkHour.OrganizationID)
		if !s.OrganizationID.IsNull() {
			m.OrganizationID = ptr.Of(s.OrganizationID.MustGet())
		} else {
			m.OrganizationID = nil
		}
	}
	if s.Monday.IsSpecified() {
		cols = append(cols, t.OrganizationWorkHour.Monday)
		if !s.Monday.IsNull() {
			m.Monday = ptr.Of(s.Monday.MustGet())
		} else {
			m.Monday = nil
		}
	}
	if s.Tuesday.IsSpecified() {
		cols = append(cols, t.OrganizationWorkHour.Tuesday)
		if !s.Tuesday.IsNull() {
			m.Tuesday = ptr.Of(s.Tuesday.MustGet())
		} else {
			m.Tuesday = nil
		}
	}
	if s.Wednesday.IsSpecified() {
		cols = append(cols, t.OrganizationWorkHour.Wednesday)
		if !s.Wednesday.IsNull() {
			m.Wednesday = ptr.Of(s.Wednesday.MustGet())
		} else {
			m.Wednesday = nil
		}
	}
	if s.Thursday.IsSpecified() {
		cols = append(cols, t.OrganizationWorkHour.Thursday)
		if !s.Thursday.IsNull() {
			m.Thursday = ptr.Of(s.Thursday.MustGet())
		} else {
			m.Thursday = nil
		}
	}
	if s.Friday.IsSpecified() {
		cols = append(cols, t.OrganizationWorkHour.Friday)
		if !s.Friday.IsNull() {
			m.Friday = ptr.Of(s.Friday.MustGet())
		} else {
			m.Friday = nil
		}
	}
	if s.Saturday.IsSpecified() {
		cols = append(cols, t.OrganizationWorkHour.Saturday)
		if !s.Saturday.IsNull() {
			m.Saturday = ptr.Of(s.Saturday.MustGet())
		} else {
			m.Saturday = nil
		}
	}
	if s.Sunday.IsSpecified() {
		cols = append(cols, t.OrganizationWorkHour.Sunday)
		if !s.Sunday.IsNull() {
			m.Sunday = ptr.Of(s.Sunday.MustGet())
		} else {
			m.Sunday = nil
		}
	}
	if s.CreatedAt.IsSpecified() {
		cols = append(cols, t.OrganizationWorkHour.CreatedAt)
		if !s.CreatedAt.IsNull() {
			m.CreatedAt = s.CreatedAt.MustGet()
		}
	}
	if s.UpdatedAt.IsSpecified() {
		cols = append(cols, t.OrganizationWorkHour.UpdatedAt)
		if !s.UpdatedAt.IsNull() {
			m.UpdatedAt = s.UpdatedAt.MustGet()
		}
	}

	return cols, m
}
