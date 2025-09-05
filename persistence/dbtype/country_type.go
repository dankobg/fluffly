package dbtype

import (
	"time"

	api "github.com/dankobg/fluffly/api/gen"
	"github.com/dankobg/fluffly/db/gen/test/public/model"
	t "github.com/dankobg/fluffly/db/gen/test/public/table"
	p "github.com/go-jet/jet/v2/postgres"
	"github.com/oapi-codegen/nullable"
)

type CountryFilters struct {
	Pagination *api.PaginationParams
}

type CountrySetter struct {
	Name       nullable.Nullable[string]    `json:"name"`
	IsoAlpha2  nullable.Nullable[string]    `json:"iso_alpha2"`
	IsoAlpha3  nullable.Nullable[string]    `json:"iso_alpha3"`
	IsoNumeric nullable.Nullable[string]    `json:"iso_numeric"`
	CreatedAt  nullable.Nullable[time.Time] `json:"created_at"`
	UpdatedAt  nullable.Nullable[time.Time] `json:"updated_at"`
}

func (s CountrySetter) ToModel(isPatch ...bool) (p.ColumnList, model.Country) {
	var cols p.ColumnList
	var m model.Country

	if len(isPatch) > 0 {
		if isPatch[0] {
			cols = append(cols, t.Country.UpdatedAt)
			m.UpdatedAt = time.Now()
		}
	}

	if s.Name.IsSpecified() {
		cols = append(cols, t.Country.Name)
		if !s.Name.IsNull() {
			m.Name = s.Name.MustGet()
		}
	}
	if s.IsoAlpha2.IsSpecified() {
		cols = append(cols, t.Country.IsoAlpha2)
		if !s.IsoAlpha2.IsNull() {
			m.IsoAlpha2 = s.IsoAlpha2.MustGet()
		}
	}
	if s.IsoAlpha3.IsSpecified() {
		cols = append(cols, t.Country.IsoAlpha3)
		if !s.IsoAlpha3.IsNull() {
			m.IsoAlpha3 = s.IsoAlpha3.MustGet()
		}
	}
	if s.IsoNumeric.IsSpecified() {
		cols = append(cols, t.Country.IsoNumeric)
		if !s.IsoNumeric.IsNull() {
			m.IsoNumeric = s.IsoNumeric.MustGet()
		}
	}
	if s.CreatedAt.IsSpecified() {
		cols = append(cols, t.Country.CreatedAt)
		if !s.CreatedAt.IsNull() {
			m.CreatedAt = s.CreatedAt.MustGet()
		}
	}
	if s.UpdatedAt.IsSpecified() {
		cols = append(cols, t.Country.UpdatedAt)
		if !s.UpdatedAt.IsNull() {
			m.UpdatedAt = s.UpdatedAt.MustGet()
		}
	}

	return cols, m
}
