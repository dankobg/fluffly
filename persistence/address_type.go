package persistence

import (
	"time"

	"github.com/dankobg/fluffly/db/gen/test/public/model"
	t "github.com/dankobg/fluffly/db/gen/test/public/table"
	"github.com/dankobg/fluffly/ptr"
	p "github.com/go-jet/jet/v2/postgres"
	"github.com/oapi-codegen/nullable"
)

type AddressSetter struct {
	CountryID     nullable.Nullable[int64]     `json:"country_id"`
	UnitNumber    nullable.Nullable[string]    `json:"unit_number"`
	StreetNumber  nullable.Nullable[string]    `json:"street_number"`
	StreetAddress nullable.Nullable[string]    `json:"street_address"`
	City          nullable.Nullable[string]    `json:"city"`
	Region        nullable.Nullable[string]    `json:"region"`
	PostalCode    nullable.Nullable[string]    `json:"postal_code"`
	Lat           nullable.Nullable[float64]   `json:"lat"`
	Lng           nullable.Nullable[float64]   `json:"lng"`
	Note          nullable.Nullable[string]    `json:"note"`
	CreatedAt     nullable.Nullable[time.Time] `json:"created_at"`
	UpdatedAt     nullable.Nullable[time.Time] `json:"updated_at"`
}

func (s AddressSetter) ToModel(isPatch ...bool) (p.ColumnList, model.Address) {
	var cols p.ColumnList
	var m model.Address

	if len(isPatch) > 0 {
		if isPatch[0] {
			cols = append(cols, t.Address.UpdatedAt)
			m.UpdatedAt = time.Now()
		}
	}

	if s.CountryID.IsSpecified() {
		cols = append(cols, t.Address.CountryID)
		if !s.CountryID.IsNull() {
			m.CountryID = s.CountryID.MustGet()
		}
	}
	if s.UnitNumber.IsSpecified() {
		cols = append(cols, t.Address.UnitNumber)
		if !s.UnitNumber.IsNull() {
			m.UnitNumber = ptr.Of(s.UnitNumber.MustGet())
		} else {
			m.UnitNumber = nil
		}
	}
	if s.StreetNumber.IsSpecified() {
		cols = append(cols, t.Address.StreetNumber)
		if !s.StreetNumber.IsNull() {
			m.StreetNumber = ptr.Of(s.StreetNumber.MustGet())
		} else {
			m.StreetNumber = nil
		}
	}
	if s.StreetAddress.IsSpecified() {
		cols = append(cols, t.Address.StreetAddress)
		if !s.StreetAddress.IsNull() {
			m.StreetAddress = s.StreetAddress.MustGet()
		}
	}
	if s.City.IsSpecified() {
		cols = append(cols, t.Address.City)
		if !s.City.IsNull() {
			m.City = s.City.MustGet()
		}
	}
	if s.Region.IsSpecified() {
		cols = append(cols, t.Address.Region)
		if !s.Region.IsNull() {
			m.Region = ptr.Of(s.Region.MustGet())
		} else {
			m.Region = nil
		}
	}
	if s.PostalCode.IsSpecified() {
		cols = append(cols, t.Address.PostalCode)
		if !s.PostalCode.IsNull() {
			m.PostalCode = ptr.Of(s.PostalCode.MustGet())
		} else {
			m.PostalCode = nil
		}
	}
	if s.Lat.IsSpecified() {
		cols = append(cols, t.Address.Lat)
		if !s.Lat.IsNull() {
			m.Lat = ptr.Of(s.Lat.MustGet())
		} else {
			m.Lat = nil
		}
	}
	if s.Lng.IsSpecified() {
		cols = append(cols, t.Address.Lng)
		if !s.Lng.IsNull() {
			m.Lng = ptr.Of(s.Lng.MustGet())
		} else {
			m.Lng = nil
		}
	}
	if s.Note.IsSpecified() {
		cols = append(cols, t.Address.Note)
		if !s.Note.IsNull() {
			m.Note = ptr.Of(s.Note.MustGet())
		} else {
			m.Note = nil
		}
	}
	if s.CreatedAt.IsSpecified() {
		cols = append(cols, t.Address.CreatedAt)
		if !s.CreatedAt.IsNull() {
			m.CreatedAt = s.CreatedAt.MustGet()
		}
	}
	if s.UpdatedAt.IsSpecified() {
		cols = append(cols, t.Address.UpdatedAt)
		if !s.UpdatedAt.IsNull() {
			m.UpdatedAt = s.UpdatedAt.MustGet()
		}
	}

	return cols, m
}
