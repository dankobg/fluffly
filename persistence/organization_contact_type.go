package persistence

import (
	"time"

	"github.com/dankobg/fluffly/db/gen/test/public/model"
	t "github.com/dankobg/fluffly/db/gen/test/public/table"
	p "github.com/go-jet/jet/v2/postgres"
	"github.com/oapi-codegen/nullable"
)

type OrganizationContactSetter struct {
	OrganizationID nullable.Nullable[int64]     `json:"organizationid"`
	AddressID      nullable.Nullable[int64]     `json:"addressid"`
	Phone          nullable.Nullable[string]    `json:"phone"`
	Email          nullable.Nullable[string]    `json:"email"`
	CreatedAt      nullable.Nullable[time.Time] `json:"created_at"`
	UpdatedAt      nullable.Nullable[time.Time] `json:"updated_at"`
}

func (s OrganizationContactSetter) ToModel(isPatch ...bool) (p.ColumnList, model.OrganizationContact) {
	var cols p.ColumnList
	var m model.OrganizationContact

	if len(isPatch) > 0 {
		if isPatch[0] {
			cols = append(cols, t.OrganizationContact.UpdatedAt)
			m.UpdatedAt = time.Now()
		}
	}

	if s.OrganizationID.IsSpecified() {
		cols = append(cols, t.OrganizationContact.OrganizationID)
		if !s.OrganizationID.IsNull() {
			m.OrganizationID = s.OrganizationID.MustGet()
		}
	}
	if s.AddressID.IsSpecified() {
		cols = append(cols, t.OrganizationContact.AddressID)
		if !s.AddressID.IsNull() {
			m.AddressID = s.AddressID.MustGet()
		}
	}
	if s.Phone.IsSpecified() {
		cols = append(cols, t.OrganizationContact.Phone)
		if !s.Phone.IsNull() {
			m.Phone = s.Phone.MustGet()
		}
	}
	if s.Email.IsSpecified() {
		cols = append(cols, t.OrganizationContact.Email)
		if !s.Email.IsNull() {
			m.Email = s.Email.MustGet()
		}
	}
	if s.CreatedAt.IsSpecified() {
		cols = append(cols, t.OrganizationContact.CreatedAt)
		if !s.CreatedAt.IsNull() {
			m.CreatedAt = s.CreatedAt.MustGet()
		}
	}
	if s.UpdatedAt.IsSpecified() {
		cols = append(cols, t.OrganizationContact.UpdatedAt)
		if !s.UpdatedAt.IsNull() {
			m.UpdatedAt = s.UpdatedAt.MustGet()
		}
	}

	return cols, m
}
