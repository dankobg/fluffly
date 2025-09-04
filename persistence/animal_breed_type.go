package persistence

import (
	"time"

	"github.com/dankobg/fluffly/db/gen/test/public/model"
	t "github.com/dankobg/fluffly/db/gen/test/public/table"
	p "github.com/go-jet/jet/v2/postgres"
	"github.com/oapi-codegen/nullable"
)

type AnimalBreedWithJoinData struct {
	model.Breed
	Primary bool
}

type AnimalBreedSetter struct {
	AnimalID  nullable.Nullable[int64]     `json:"animal_id"`
	BreedID   nullable.Nullable[int64]     `json:"breed_id"`
	Primary   nullable.Nullable[bool]      `json:"primary"`
	Name      nullable.Nullable[string]    `json:"name"`
	CreatedAt nullable.Nullable[time.Time] `json:"created_at"`
	UpdatedAt nullable.Nullable[time.Time] `json:"updated_at"`
}

func (s AnimalBreedSetter) ToModel(isPatch ...bool) (p.ColumnList, model.AnimalBreed) {
	var cols p.ColumnList
	var m model.AnimalBreed

	if len(isPatch) > 0 {
		if isPatch[0] {
			cols = append(cols, t.AnimalBreed.UpdatedAt)
			m.UpdatedAt = time.Now()
		}
	}

	if s.AnimalID.IsSpecified() {
		cols = append(cols, t.AnimalBreed.AnimalID)
		if !s.AnimalID.IsNull() {
			m.AnimalID = s.AnimalID.MustGet()
		}
	}
	if s.BreedID.IsSpecified() {
		cols = append(cols, t.AnimalBreed.BreedID)
		if !s.BreedID.IsNull() {
			m.BreedID = s.BreedID.MustGet()
		}
	}
	if s.Primary.IsSpecified() {
		cols = append(cols, t.AnimalBreed.Primary)
		if !s.Primary.IsNull() {
			m.Primary = s.Primary.MustGet()
		}
	}
	if s.CreatedAt.IsSpecified() {
		cols = append(cols, t.AnimalBreed.CreatedAt)
		if !s.CreatedAt.IsNull() {
			m.CreatedAt = s.CreatedAt.MustGet()
		}
	}
	if s.UpdatedAt.IsSpecified() {
		cols = append(cols, t.AnimalBreed.UpdatedAt)
		if !s.UpdatedAt.IsNull() {
			m.UpdatedAt = s.UpdatedAt.MustGet()
		}
	}

	return cols, m
}
