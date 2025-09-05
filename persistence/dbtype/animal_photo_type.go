package dbtype

import (
	"time"

	"github.com/dankobg/fluffly/db/gen/test/public/model"
	t "github.com/dankobg/fluffly/db/gen/test/public/table"
	"github.com/dankobg/fluffly/ptr"
	p "github.com/go-jet/jet/v2/postgres"
	"github.com/oapi-codegen/nullable"
)

type AnimalPhotoSetter struct {
	AnimalID  nullable.Nullable[int64]     `json:"animalid"`
	Small     nullable.Nullable[string]    `json:"small"`
	Medium    nullable.Nullable[string]    `json:"medium"`
	Large     nullable.Nullable[string]    `json:"large"`
	Full      nullable.Nullable[string]    `json:"full"`
	CreatedAt nullable.Nullable[time.Time] `json:"created_at"`
	UpdatedAt nullable.Nullable[time.Time] `json:"updated_at"`
}

func (s AnimalPhotoSetter) ToModel(isPatch ...bool) (p.ColumnList, model.AnimalPhoto) {
	var cols p.ColumnList
	var m model.AnimalPhoto

	if len(isPatch) > 0 {
		if isPatch[0] {
			cols = append(cols, t.AnimalPhoto.UpdatedAt)
			m.UpdatedAt = time.Now()
		}
	}

	if s.AnimalID.IsSpecified() {
		cols = append(cols, t.AnimalPhoto.AnimalID)
		if !s.AnimalID.IsNull() {
			m.AnimalID = ptr.Of(s.AnimalID.MustGet())
		} else {
			m.AnimalID = nil
		}
	}
	if s.Small.IsSpecified() {
		cols = append(cols, t.AnimalPhoto.Small)
		if !s.Small.IsNull() {
			m.Small = ptr.Of(s.Small.MustGet())
		} else {
			m.Small = nil
		}
	}
	if s.Medium.IsSpecified() {
		cols = append(cols, t.AnimalPhoto.Medium)
		if !s.Medium.IsNull() {
			m.Medium = ptr.Of(s.Medium.MustGet())
		} else {
			m.Medium = nil
		}
	}
	if s.Large.IsSpecified() {
		cols = append(cols, t.AnimalPhoto.Large)
		if !s.Large.IsNull() {
			m.Large = ptr.Of(s.Large.MustGet())
		} else {
			m.Large = nil
		}
	}
	if s.Full.IsSpecified() {
		cols = append(cols, t.AnimalPhoto.Full)
		if !s.Full.IsNull() {
			m.Full = ptr.Of(s.Full.MustGet())
		} else {
			m.Full = nil
		}
	}
	if s.CreatedAt.IsSpecified() {
		cols = append(cols, t.AnimalPhoto.CreatedAt)
		if !s.CreatedAt.IsNull() {
			m.CreatedAt = s.CreatedAt.MustGet()
		}
	}
	if s.UpdatedAt.IsSpecified() {
		cols = append(cols, t.AnimalPhoto.UpdatedAt)
		if !s.UpdatedAt.IsNull() {
			m.UpdatedAt = s.UpdatedAt.MustGet()
		}
	}

	return cols, m
}
