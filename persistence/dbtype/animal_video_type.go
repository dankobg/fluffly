package dbtype

import (
	"time"

	"github.com/dankobg/fluffly/db/gen/test/public/model"
	t "github.com/dankobg/fluffly/db/gen/test/public/table"
	"github.com/dankobg/fluffly/ptr"
	p "github.com/go-jet/jet/v2/postgres"
	"github.com/oapi-codegen/nullable"
)

type AnimalVideoSetter struct {
	AnimalID  nullable.Nullable[int64]     `json:"animalid"`
	URL       nullable.Nullable[string]    `json:"url"`
	CreatedAt nullable.Nullable[time.Time] `json:"created_at"`
	UpdatedAt nullable.Nullable[time.Time] `json:"updated_at"`
}

func (s AnimalVideoSetter) ToModel(isPatch ...bool) (p.ColumnList, model.AnimalVideo) {
	var cols p.ColumnList
	var m model.AnimalVideo

	if len(isPatch) > 0 {
		if isPatch[0] {
			cols = append(cols, t.AnimalVideo.UpdatedAt)
			m.UpdatedAt = time.Now()
		}
	}

	if s.AnimalID.IsSpecified() {
		cols = append(cols, t.AnimalVideo.AnimalID)
		if !s.AnimalID.IsNull() {
			m.AnimalID = ptr.Of(s.AnimalID.MustGet())
		} else {
			m.AnimalID = nil
		}
	}
	if s.URL.IsSpecified() {
		cols = append(cols, t.AnimalVideo.URL)
		if !s.URL.IsNull() {
			m.URL = s.URL.MustGet()
		}
	}
	if s.CreatedAt.IsSpecified() {
		cols = append(cols, t.AnimalVideo.CreatedAt)
		if !s.CreatedAt.IsNull() {
			m.CreatedAt = s.CreatedAt.MustGet()
		}
	}
	if s.UpdatedAt.IsSpecified() {
		cols = append(cols, t.AnimalVideo.UpdatedAt)
		if !s.UpdatedAt.IsNull() {
			m.UpdatedAt = s.UpdatedAt.MustGet()
		}
	}

	return cols, m
}
