package dbtype

import (
	"time"

	api "github.com/dankobg/fluffly/api/gen"
	"github.com/dankobg/fluffly/db/gen/test/public/model"
	t "github.com/dankobg/fluffly/db/gen/test/public/table"
	"github.com/dankobg/fluffly/persistence/dbcustom"
	"github.com/dankobg/fluffly/ptr"
	p "github.com/go-jet/jet/v2/postgres"
	"github.com/google/uuid"
	"github.com/oapi-codegen/nullable"
)

type AnimalFilters struct {
	Pagination *api.PaginationParams
}

type AnimalWithJoinData struct {
	model.Animal
	Type         model.AnimalType
	Species      model.AnimalSpecies
	Microchip    model.Microchip
	Breeds       []AnimalBreedWithJoinData `json_column:"breeds"`
	Tags         []model.Tag               `json_column:"tags"`
	Photos       []model.AnimalPhoto       `json_column:"photos"`
	Videos       []model.AnimalVideo       `json_column:"videos"`
	Organization struct {
		model.Organization
		WorkHour model.OrganizationWorkHour
		Contact  struct {
			model.OrganizationContact
			Address struct {
				model.Address
				Country model.Country
			}
		}
	}
	User  model.User
	Likes int64
}

type AnimalCreateSetter struct {
	Animal    AnimalSetter
	Microchip nullable.Nullable[MicrochipSetter]
	Breeds    nullable.Nullable[[]AnimalBreedSetter]
	Tags      nullable.Nullable[[]AnimalTagSetter]
	Photos    nullable.Nullable[[]AnimalPhotoSetter]
	Videos    nullable.Nullable[[]AnimalVideoSetter]
}

type AnimalSetter struct {
	UserID               nullable.Nullable[uuid.UUID]         `json:"user_id"`
	OrganizationID       nullable.Nullable[int64]             `json:"organization_id"`
	TypeID               nullable.Nullable[int64]             `json:"type_id"`
	SpeciesID            nullable.Nullable[int64]             `json:"species_id"`
	Name                 nullable.Nullable[string]            `json:"name"`
	Gender               nullable.Nullable[model.Gender]      `json:"gender"`
	Hermaphrodite        nullable.Nullable[bool]              `json:"hermaphrodite"`
	Age                  nullable.Nullable[string]            `json:"age"`
	Size                 nullable.Nullable[string]            `json:"size"`
	ImageObjectKind      nullable.Nullable[string]            `json:"image_object_kind"`
	ImageObjectRefSmall  nullable.Nullable[string]            `json:"image_object_ref_small"`
	ImageObjectRefMedium nullable.Nullable[string]            `json:"image_object_ref_medium"`
	ImageObjectRefLarge  nullable.Nullable[string]            `json:"image_object_ref_large"`
	ImageObjectRefFull   nullable.Nullable[string]            `json:"image_object_ref_full"`
	Description          nullable.Nullable[string]            `json:"description"`
	Distance             nullable.Nullable[string]            `json:"distance"`
	Properties           nullable.Nullable[dbcustom.JsonType] `json:"properties"`
	Status               nullable.Nullable[string]            `json:"status"`
	StatusChangedAt      nullable.Nullable[time.Time]         `json:"status_changed_at"`
	AdoptedAt            nullable.Nullable[time.Time]         `json:"adopted_at"`
	CreatedAt            nullable.Nullable[time.Time]         `json:"created_at"`
	UpdatedAt            nullable.Nullable[time.Time]         `json:"updated_at"`
}

func (s AnimalSetter) ToModel(isPatch ...bool) (p.ColumnList, model.Animal) {
	var cols p.ColumnList
	var m model.Animal

	if len(isPatch) > 0 {
		if isPatch[0] {
			cols = append(cols, t.Animal.UpdatedAt)
			m.UpdatedAt = time.Now()
		}
	}

	if s.UserID.IsSpecified() {
		cols = append(cols, t.Animal.UserID)
		if !s.UserID.IsNull() {
			m.UserID = ptr.Of(s.UserID.MustGet())
		} else {
			m.UserID = nil
		}
	}
	if s.OrganizationID.IsSpecified() {
		cols = append(cols, t.Animal.OrganizationID)
		if !s.OrganizationID.IsNull() {
			m.OrganizationID = ptr.Of(s.OrganizationID.MustGet())
		} else {
			m.OrganizationID = nil
		}
	}
	if s.TypeID.IsSpecified() {
		cols = append(cols, t.Animal.TypeID)
		if !s.TypeID.IsNull() {
			m.TypeID = s.TypeID.MustGet()
		}
	}
	if s.SpeciesID.IsSpecified() {
		cols = append(cols, t.Animal.SpeciesID)
		if !s.SpeciesID.IsNull() {
			m.SpeciesID = s.SpeciesID.MustGet()
		}
	}
	if s.Name.IsSpecified() {
		cols = append(cols, t.Animal.Name)
		if !s.Name.IsNull() {
			m.Name = s.Name.MustGet()
		}
	}
	if s.Gender.IsSpecified() {
		cols = append(cols, t.Animal.Gender)
		if !s.Gender.IsNull() {
			m.Gender = ptr.Of(s.Gender.MustGet())
		} else {
			m.Gender = nil
		}
	}
	if s.Hermaphrodite.IsSpecified() {
		cols = append(cols, t.Animal.Hermaphrodite)
		if !s.Hermaphrodite.IsNull() {
			m.Hermaphrodite = s.Hermaphrodite.MustGet()
		}
	}
	if s.Age.IsSpecified() {
		cols = append(cols, t.Animal.Age)
		if !s.Age.IsNull() {
			m.Age = s.Age.MustGet()
		}
	}
	if s.Size.IsSpecified() {
		cols = append(cols, t.Animal.Size)
		if !s.Size.IsNull() {
			m.Size = s.Size.MustGet()
		}
	}
	if s.ImageObjectKind.IsSpecified() {
		cols = append(cols, t.Animal.ImageObjectKind)
		if !s.ImageObjectKind.IsNull() {
			m.ImageObjectKind = s.ImageObjectKind.MustGet()
		}
	}
	if s.ImageObjectRefSmall.IsSpecified() {
		cols = append(cols, t.Animal.ImageObjectRefSmall)
		if !s.ImageObjectRefSmall.IsNull() {
			m.ImageObjectRefSmall = s.ImageObjectRefSmall.MustGet()
		}
	}
	if s.ImageObjectRefMedium.IsSpecified() {
		cols = append(cols, t.Animal.ImageObjectRefMedium)
		if !s.ImageObjectRefMedium.IsNull() {
			m.ImageObjectRefMedium = s.ImageObjectRefMedium.MustGet()
		}
	}
	if s.ImageObjectRefLarge.IsSpecified() {
		cols = append(cols, t.Animal.ImageObjectRefLarge)
		if !s.ImageObjectRefLarge.IsNull() {
			m.ImageObjectRefLarge = s.ImageObjectRefLarge.MustGet()
		}
	}
	if s.ImageObjectRefFull.IsSpecified() {
		cols = append(cols, t.Animal.ImageObjectRefFull)
		if !s.ImageObjectRefFull.IsNull() {
			m.ImageObjectRefFull = s.ImageObjectRefFull.MustGet()
		}
	}
	if s.Description.IsSpecified() {
		cols = append(cols, t.Animal.Description)
		if !s.Description.IsNull() {
			m.Description = ptr.Of(s.Description.MustGet())
		} else {
			m.Description = nil
		}
	}
	if s.Distance.IsSpecified() {
		cols = append(cols, t.Animal.Distance)
		if !s.Distance.IsNull() {
			m.Distance = ptr.Of(s.Distance.MustGet())
		} else {
			m.Distance = nil
		}
	}
	if s.Properties.IsSpecified() {
		cols = append(cols, t.Animal.Properties)
		if !s.Properties.IsNull() {
			m.Properties = s.Properties.MustGet()
		}
	}
	if s.Status.IsSpecified() {
		cols = append(cols, t.Animal.Status)
		if !s.Status.IsNull() {
			m.Status = ptr.Of(s.Status.MustGet())
		} else {
			m.Status = nil
		}
	}
	if s.StatusChangedAt.IsSpecified() {
		cols = append(cols, t.Animal.StatusChangedAt)
		if !s.StatusChangedAt.IsNull() {
			m.StatusChangedAt = ptr.Of(s.StatusChangedAt.MustGet())
		} else {
			m.StatusChangedAt = nil
		}
	}
	if s.AdoptedAt.IsSpecified() {
		cols = append(cols, t.Animal.AdoptedAt)
		if !s.AdoptedAt.IsNull() {
			m.AdoptedAt = ptr.Of(s.AdoptedAt.MustGet())
		} else {
			m.AdoptedAt = nil
		}
	}
	if s.CreatedAt.IsSpecified() {
		cols = append(cols, t.Animal.CreatedAt)
		if !s.CreatedAt.IsNull() {
			m.CreatedAt = s.CreatedAt.MustGet()
		}
	}
	if s.UpdatedAt.IsSpecified() {
		cols = append(cols, t.Animal.UpdatedAt)
		if !s.UpdatedAt.IsNull() {
			m.UpdatedAt = s.UpdatedAt.MustGet()
		}
	}

	return cols, m
}
