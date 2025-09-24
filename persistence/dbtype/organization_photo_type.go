package dbtype

import (
	"time"

	"github.com/dankobg/fluffly/db/gen/test/public/model"
	t "github.com/dankobg/fluffly/db/gen/test/public/table"
	"github.com/dankobg/fluffly/ptr"
	p "github.com/go-jet/jet/v2/postgres"
	"github.com/oapi-codegen/nullable"
)

type OrganizationPhotoSetter struct {
	OrganizationID  nullable.Nullable[int64]     `json:"organizationid"`
	ObjectKind      nullable.Nullable[string]    `json:"object_kind"`
	ObjectRefSmall  nullable.Nullable[string]    `json:"object_ref_small"`
	ObjectRefMedium nullable.Nullable[string]    `json:"object_ref_medium"`
	ObjectRefLarge  nullable.Nullable[string]    `json:"object_ref_large"`
	ObjectRefFull   nullable.Nullable[string]    `json:"object_ref_full"`
	CreatedAt       nullable.Nullable[time.Time] `json:"created_at"`
	UpdatedAt       nullable.Nullable[time.Time] `json:"updated_at"`
}

func (s OrganizationPhotoSetter) ToModel(isPatch ...bool) (p.ColumnList, model.OrganizationPhoto) {
	var cols p.ColumnList
	var m model.OrganizationPhoto

	if len(isPatch) > 0 {
		if isPatch[0] {
			cols = append(cols, t.OrganizationPhoto.UpdatedAt)
			m.UpdatedAt = time.Now()
		}
	}

	if s.OrganizationID.IsSpecified() {
		cols = append(cols, t.OrganizationPhoto.OrganizationID)
		if !s.OrganizationID.IsNull() {
			m.OrganizationID = ptr.Of(s.OrganizationID.MustGet())
		} else {
			m.OrganizationID = nil
		}
	}
	if s.ObjectKind.IsSpecified() {
		cols = append(cols, t.OrganizationPhoto.ObjectKind)
		if !s.ObjectKind.IsNull() {
			m.ObjectKind = s.ObjectKind.MustGet()
		}
	}
	if s.ObjectRefSmall.IsSpecified() {
		cols = append(cols, t.OrganizationPhoto.ObjectRefSmall)
		if !s.ObjectRefSmall.IsNull() {
			m.ObjectRefSmall = ptr.Of(s.ObjectRefSmall.MustGet())
		} else {
			m.ObjectRefSmall = nil
		}
	}
	if s.ObjectRefMedium.IsSpecified() {
		cols = append(cols, t.OrganizationPhoto.ObjectRefMedium)
		if !s.ObjectRefMedium.IsNull() {
			m.ObjectRefMedium = ptr.Of(s.ObjectRefMedium.MustGet())
		} else {
			m.ObjectRefMedium = nil
		}
	}
	if s.ObjectRefLarge.IsSpecified() {
		cols = append(cols, t.OrganizationPhoto.ObjectRefLarge)
		if !s.ObjectRefLarge.IsNull() {
			m.ObjectRefLarge = ptr.Of(s.ObjectRefLarge.MustGet())
		} else {
			m.ObjectRefLarge = nil
		}
	}
	if s.ObjectRefFull.IsSpecified() {
		cols = append(cols, t.OrganizationPhoto.ObjectRefFull)
		if !s.ObjectRefFull.IsNull() {
			m.ObjectRefFull = ptr.Of(s.ObjectRefFull.MustGet())
		} else {
			m.ObjectRefFull = nil
		}
	}
	if s.CreatedAt.IsSpecified() {
		cols = append(cols, t.OrganizationPhoto.CreatedAt)
		if !s.CreatedAt.IsNull() {
			m.CreatedAt = s.CreatedAt.MustGet()
		}
	}
	if s.UpdatedAt.IsSpecified() {
		cols = append(cols, t.OrganizationPhoto.UpdatedAt)
		if !s.UpdatedAt.IsNull() {
			m.UpdatedAt = s.UpdatedAt.MustGet()
		}
	}

	return cols, m
}
