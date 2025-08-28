package persistence

import (
	"time"

	"github.com/dankobg/fluffly/db/gen/test/public/model"
	t "github.com/dankobg/fluffly/db/gen/test/public/table"
	p "github.com/go-jet/jet/v2/postgres"
	"github.com/google/uuid"
	"github.com/oapi-codegen/nullable"
)

type UserSetter struct {
	ID        nullable.Nullable[uuid.UUID] `json:"id"`
	CreatedAt nullable.Nullable[time.Time] `json:"created_at"`
	UpdatedAt nullable.Nullable[time.Time] `json:"updated_at"`
}

func (s UserSetter) ToModel(isPatch ...bool) (p.ColumnList, model.User) {
	var cols p.ColumnList
	var m model.User

	if len(isPatch) > 0 {
		if isPatch[0] {
			cols = append(cols, t.User.UpdatedAt)
			m.UpdatedAt = time.Now()
		}
	}

	if s.ID.IsSpecified() {
		cols = append(cols, t.User.ID)
		if !s.ID.IsNull() {
			m.ID = s.ID.MustGet()
		}
	}
	if s.CreatedAt.IsSpecified() {
		cols = append(cols, t.User.CreatedAt)
		if !s.CreatedAt.IsNull() {
			m.CreatedAt = s.CreatedAt.MustGet()
		}
	}
	if s.UpdatedAt.IsSpecified() {
		cols = append(cols, t.User.UpdatedAt)
		if !s.UpdatedAt.IsNull() {
			m.UpdatedAt = s.UpdatedAt.MustGet()
		}
	}

	return cols, m
}
