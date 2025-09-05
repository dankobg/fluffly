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

type MicrochipFilters struct {
	Pagination *api.PaginationParams
}

type MicrochipSetter struct {
	AnimalID    nullable.Nullable[int64]     `json:"animal_id"`
	Number      nullable.Nullable[string]    `json:"number"`
	Brand       nullable.Nullable[string]    `json:"brand"`
	Description nullable.Nullable[string]    `json:"description"`
	Location    nullable.Nullable[string]    `json:"location"`
	CreatedAt   nullable.Nullable[time.Time] `json:"created_at"`
	UpdatedAt   nullable.Nullable[time.Time] `json:"updated_at"`
}

func (s MicrochipSetter) ToModel(isPatch ...bool) (p.ColumnList, model.Microchip) {
	var cols p.ColumnList
	var m model.Microchip

	if len(isPatch) > 0 {
		if isPatch[0] {
			cols = append(cols, t.Microchip.UpdatedAt)
			m.UpdatedAt = time.Now()
		}
	}

	if s.AnimalID.IsSpecified() {
		cols = append(cols, t.Microchip.AnimalID)
		if !s.AnimalID.IsNull() {
			m.AnimalID = ptr.Of(s.AnimalID.MustGet())
		} else {
			m.AnimalID = nil
		}
	}
	if s.Number.IsSpecified() {
		cols = append(cols, t.Microchip.Number)
		if !s.Number.IsNull() {
			m.Number = s.Number.MustGet()
		}
	}
	if s.Brand.IsSpecified() {
		cols = append(cols, t.Microchip.Brand)
		if !s.Brand.IsNull() {
			m.Brand = ptr.Of(s.Brand.MustGet())
		} else {
			m.Brand = nil
		}
	}
	if s.Description.IsSpecified() {
		cols = append(cols, t.Microchip.Description)
		if !s.Description.IsNull() {
			m.Description = ptr.Of(s.Description.MustGet())
		} else {
			m.Description = nil
		}
	}
	if s.Location.IsSpecified() {
		cols = append(cols, t.Microchip.Location)
		if !s.Location.IsNull() {
			m.Location = ptr.Of(s.Location.MustGet())
		} else {
			m.Location = nil
		}
	}
	if s.CreatedAt.IsSpecified() {
		cols = append(cols, t.Microchip.CreatedAt)
		if !s.CreatedAt.IsNull() {
			m.CreatedAt = s.CreatedAt.MustGet()
		}
	}
	if s.UpdatedAt.IsSpecified() {
		cols = append(cols, t.Microchip.UpdatedAt)
		if !s.UpdatedAt.IsNull() {
			m.UpdatedAt = s.UpdatedAt.MustGet()
		}
	}

	return cols, m
}
