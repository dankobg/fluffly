package dbtype

import (
	"time"

	"github.com/dankobg/fluffly/db/gen/test/public/model"
	t "github.com/dankobg/fluffly/db/gen/test/public/table"
	"github.com/dankobg/fluffly/ptr"
	p "github.com/go-jet/jet/v2/postgres"
	"github.com/oapi-codegen/nullable"
)

type OrganizationVideoSetter struct {
	OrganizationID nullable.Nullable[int64]     `json:"organizationid"`
	ObjectKind     nullable.Nullable[string]    `json:"object_kind"`
	ObjectRef      nullable.Nullable[string]    `json:"object_ref"`
	CreatedAt      nullable.Nullable[time.Time] `json:"created_at"`
	UpdatedAt      nullable.Nullable[time.Time] `json:"updated_at"`
}

func (s OrganizationVideoSetter) ToModel(isPatch ...bool) (p.ColumnList, model.OrganizationVideo) {
	var cols p.ColumnList
	var m model.OrganizationVideo

	if len(isPatch) > 0 {
		if isPatch[0] {
			cols = append(cols, t.OrganizationVideo.UpdatedAt)
			m.UpdatedAt = time.Now()
		}
	}

	if s.OrganizationID.IsSpecified() {
		cols = append(cols, t.OrganizationVideo.OrganizationID)
		if !s.OrganizationID.IsNull() {
			m.OrganizationID = ptr.Of(s.OrganizationID.MustGet())
		} else {
			m.OrganizationID = nil
		}
	}
	if s.ObjectKind.IsSpecified() {
		cols = append(cols, t.OrganizationVideo.ObjectKind)
		if !s.ObjectKind.IsNull() {
			m.ObjectKind = s.ObjectKind.MustGet()
		}
	}
	if s.ObjectRef.IsSpecified() {
		cols = append(cols, t.OrganizationVideo.ObjectRef)
		if !s.ObjectRef.IsNull() {
			m.ObjectRef = s.ObjectRef.MustGet()
		}
	}
	if s.CreatedAt.IsSpecified() {
		cols = append(cols, t.OrganizationVideo.CreatedAt)
		if !s.CreatedAt.IsNull() {
			m.CreatedAt = s.CreatedAt.MustGet()
		}
	}
	if s.UpdatedAt.IsSpecified() {
		cols = append(cols, t.OrganizationVideo.UpdatedAt)
		if !s.UpdatedAt.IsNull() {
			m.UpdatedAt = s.UpdatedAt.MustGet()
		}
	}

	return cols, m
}
