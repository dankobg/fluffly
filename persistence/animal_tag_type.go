package persistence

import (
	"time"

	"github.com/dankobg/fluffly/db/gen/test/public/model"
	t "github.com/dankobg/fluffly/db/gen/test/public/table"
	"github.com/dankobg/fluffly/ptr"
	p "github.com/go-jet/jet/v2/postgres"
	"github.com/oapi-codegen/nullable"
)

type AnimalTagSetter struct {
	AnimalID  nullable.Nullable[int64]     `json:"animalid"`
	Name      nullable.Nullable[string]    `json:"name"`
	CreatedAt nullable.Nullable[time.Time] `json:"created_at"`
	UpdatedAt nullable.Nullable[time.Time] `json:"updated_at"`
}

func (s AnimalTagSetter) ToModel(isPatch ...bool) (p.ColumnList, model.Tag) {
	var cols p.ColumnList
	var m model.Tag

	if len(isPatch) > 0 {
		if isPatch[0] {
			cols = append(cols, t.Tag.UpdatedAt)
			m.UpdatedAt = time.Now()
		}
	}

	if s.AnimalID.IsSpecified() {
		cols = append(cols, t.Tag.AnimalID)
		if !s.AnimalID.IsNull() {
			m.AnimalID = ptr.Of(s.AnimalID.MustGet())
		} else {
			m.AnimalID = nil
		}
	}
	if s.Name.IsSpecified() {
		cols = append(cols, t.Tag.Name)
		if !s.Name.IsNull() {
			m.Name = s.Name.MustGet()
		}
	}
	if s.CreatedAt.IsSpecified() {
		cols = append(cols, t.Tag.CreatedAt)
		if !s.CreatedAt.IsNull() {
			m.CreatedAt = s.CreatedAt.MustGet()
		}
	}
	if s.UpdatedAt.IsSpecified() {
		cols = append(cols, t.Tag.UpdatedAt)
		if !s.UpdatedAt.IsNull() {
			m.UpdatedAt = s.UpdatedAt.MustGet()
		}
	}

	return cols, m
}
