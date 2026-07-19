package postgres

import (
	"context"
	"errors"
	"fmt"
	"slices"

	"github.com/aarondl/opt/omit"
	"github.com/aarondl/opt/omitnull"
	api "github.com/dankobg/fluffly/api/gen"
	"github.com/dankobg/fluffly/db/gen/models"
	"github.com/dankobg/fluffly/persistence"
	"github.com/dankobg/fluffly/persistence/dbcustom"
	"github.com/dankobg/fluffly/persistence/dbtype"
	"github.com/google/uuid"
	"github.com/jackc/pgerrcode"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/stephenafamo/bob"
	"github.com/stephenafamo/bob/dialect/psql"
	"github.com/stephenafamo/bob/dialect/psql/dm"
	"github.com/stephenafamo/bob/dialect/psql/im"
	"github.com/stephenafamo/bob/dialect/psql/sm"
	"github.com/stephenafamo/bob/dialect/psql/um"
	bobpgx "github.com/stephenafamo/bob/drivers/pgx"
	"github.com/stephenafamo/bob/orm"
	"github.com/stephenafamo/scan"
)

var _ persistence.AnimalPersistor = (*PgAnimalPersistor)(nil)

type PgAnimalPersistor struct {
	*PgPersistor
}

func NewPgAnimalPersistor(ps *PgPersistor) *PgAnimalPersistor {
	return &PgAnimalPersistor{
		PgPersistor: ps,
	}
}

type (
	ErrAnimalIntegrityViolation  struct{ errIntegrityViolation }
	ErrAnimalUniqueViolation     struct{ errUniqueViolation }
	ErrAnimalForeignKeyViolation struct{ errForeignKeyViolation }
	ErrAnimalCheckViolation      struct{ errCheckViolation }
)

var (
	ErrAnimalBreedNotFound   = errors.New("animal breed not found")
	ErrAnimalTagNotFound     = errors.New("animal tag not found")
	ErrAnimalPhotoNotFound   = errors.New("animal photo not found")
	ErrAnimalVideoNotFound   = errors.New("animal video not found")
	ErrAnimalNotFound        = errors.New("animal not found")
	ErrAnimalTypeNotFound    = errors.New("animal type not found")
	ErrAnimalSpeciesNotFound = errors.New("animal species not found")
	errAnimalIntegrity       = ErrAnimalIntegrityViolation{}
	ErrAnimalAlreadyAdopted  = errors.New("animal already adopted")

	errAnimalForeignKeyUserID           = ErrAnimalForeignKeyViolation{errForeignKeyViolation: errForeignKeyViolation{Name: "user_id"}}
	errAnimalForeignKeyOrganizationID   = ErrAnimalForeignKeyViolation{errForeignKeyViolation: errForeignKeyViolation{Name: "organization_id"}}
	errAnimalCheckAgeValid              = ErrAnimalCheckViolation{errCheckViolation: errCheckViolation{Name: "age valid"}}
	errAnimalCheckSizeValid             = ErrAnimalCheckViolation{errCheckViolation: errCheckViolation{Name: "size valid"}}
	errAnimalCheckUserIdOrOrgIDProvided = ErrAnimalCheckViolation{errCheckViolation: errCheckViolation{Name: "user_id or organization_id not provided"}}
	errAnimalLikeUnique                 = ErrAnimalUniqueViolation{errUniqueViolation: errUniqueViolation{Name: "like"}}

	errAnimalTypeUniqueName = ErrAnimalUniqueViolation{errUniqueViolation: errUniqueViolation{Name: "animal_type.name"}}
)

func convertAnimalPgError(pgErr *pgconn.PgError) error {
	if pgerrcode.IsIntegrityConstraintViolation(pgErr.Code) {
		switch pgErr.ConstraintName {
		case "fk_animal_user_id":
			return errAnimalForeignKeyUserID
		case "fk_animal_organization_id":
			return errAnimalForeignKeyOrganizationID
		case "ck_animal_age_valid":
			return errAnimalCheckAgeValid
		case "ck_animal_size_valid":
			return errAnimalCheckSizeValid
		case "ck_animal_user_or_organization_provided":
			return errAnimalCheckUserIdOrOrgIDProvided
		case "pk_user_animal_like_user_id_animal_id":
			return errAnimalLikeUnique
		}

		return errAnimalIntegrity
	}

	return pgErr
}

func (pst *PgAnimalPersistor) GetAnimalMinimalByID(ctx context.Context, animalID int64) (models.Animal, error) {
	q := psql.Select(
		sm.Columns(models.Animals.Columns),
		sm.From(models.Animals.Name()),
		sm.Where(models.Animals.Columns.ID.EQ(psql.Arg(animalID))),
	)

	animal, err := bob.One(ctx, pst.exec, q, scan.StructMapper[models.Animal]())
	if err != nil {
		return models.Animal{}, fmt.Errorf("query animal minimal")
	}

	return animal, nil
}

func (pst *PgAnimalPersistor) ListAnimals(ctx context.Context, filters dbtype.ListAnimalsFilters) (dbtype.PagedResult[dbtype.AnimalWithJoinData], error) {
	q := psql.Select(
		sm.Columns(
			models.Animals.Columns,
			models.AnimalTypes.Columns.WithPrefix("type."),
			models.AnimalSpecies.Columns.WithPrefix("specie."),
			psql.Raw(`COUNT("user_animal_like"."animal_id") as "likes"`),
			psql.Select(
				sm.Columns(models.Adoptions.Columns.ID.As("adoption_id")),
				sm.From(models.Adoptions.Name()),
				sm.Where(models.Adoptions.Columns.AnimalID.EQ(models.Animals.Columns.ID)),
			),
		),
		sm.From(models.Animals.Name()),
		sm.LeftJoin(models.AnimalTypes.Name()).
			On(models.Animals.Columns.TypeID.EQ(models.AnimalTypes.Columns.ID)),
		sm.LeftJoin(models.AnimalSpecies.Name()).
			On(models.Animals.Columns.SpecieID.EQ(models.AnimalSpecies.Columns.ID)),
		sm.LeftJoin(models.UserAnimalLikes.Name()).
			On(models.UserAnimalLikes.Columns.AnimalID.EQ(models.Animals.Columns.ID)),
		sm.GroupBy(models.Animals.Columns.ID),
		sm.GroupBy(models.AnimalTypes.Columns.ID),
		sm.GroupBy(models.AnimalSpecies.Columns.ID),
	)
	addOrderBy(&q, filters.Sort, models.Animals.Columns.Except(
		models.Animals.Columns.Properties.String(),
	).Names())
	addPagination(&q, filters.Page, filters.PageSize)

	if filters.UserID != nil {
		q.Apply(sm.Columns(psql.Exists(psql.Select(
			sm.Columns(1),
			sm.From(models.UserAnimalLikes.Name()),
			sm.Where(models.UserAnimalLikes.Columns.AnimalID.EQ(models.Animals.Columns.ID).
				And(models.UserAnimalLikes.Columns.UserID.EQ(psql.Arg(*filters.UserID))),
			),
		)).As("liked")))
	}

	if hasAnyLogicFilters(&filters.ListAnimalsParams) {
		if filters.Embed != nil {
			for _, embed := range *filters.Embed {
				switch embed {
				case api.ListAnimalsParamsEmbedBreeds:
					q.Apply(sm.Columns(
						psql.Raw(`(
	  SELECT
	    json_agg(row_to_json(breeds_row))
	  FROM
	    (
	      SELECT
	     	  breed.id AS "id",
	        breed.name AS "name",
	        animal_breed.primary AS "primary",
	        animal_breed.created_at AS "createdAt",
	        animal_breed.updated_at AS "updatedAt"
	      FROM
	        public.animal_breed
	        INNER JOIN public.breed ON animal_breed.breed_id = breed.id
	      WHERE
	        animal_breed.animal_id = animal.id
	    ) AS breeds_row
	) AS "breeds"`),
					))
				case api.ListAnimalsParamsEmbedMicrochip:
					q.Apply(sm.Columns(
						models.Microchips.Columns.WithPrefix("microchip."),
					))
					q.Apply(
						sm.LeftJoin(models.Microchips.Name()).
							On(models.Animals.Columns.ID.EQ(models.Microchips.Columns.AnimalID)),
					)
					q.Apply(sm.GroupBy(models.Microchips.Columns.ID))
				case api.ListAnimalsParamsEmbedOrganization:
					q.Apply(sm.Columns(
						models.Organizations.Columns.WithPrefix("organization."),
						models.OrganizationContacts.Columns.WithPrefix("organization_contact."),
						models.OrganizationWorkHours.Columns.WithPrefix("organization_work_hour."),
						models.Addresses.Columns.Except("coords").WithPrefix("organization_contact_address."),
						psql.Raw(`ST_AsBinary("address"."coords") AS "organization_contact_address.coords"`),
						models.Countries.Columns.WithPrefix("organization_contact_address_country."),
					))
					q.Apply(
						sm.LeftJoin(models.Organizations.Name()).
							On(models.Organizations.Columns.ID.EQ(models.Animals.Columns.OrganizationID)),
						sm.LeftJoin(models.OrganizationContacts.Name()).
							On(models.Organizations.Columns.ID.EQ(models.OrganizationContacts.Columns.OrganizationID)),
						sm.LeftJoin(models.Addresses.Name()).
							On(models.Addresses.Columns.ID.EQ(models.OrganizationContacts.Columns.AddressID)),
						sm.LeftJoin(models.Countries.Name()).
							On(models.Countries.Columns.ID.EQ(models.Addresses.Columns.CountryID)),
						sm.LeftJoin(models.OrganizationWorkHours.Name()).
							On(models.Organizations.Columns.ID.EQ(models.OrganizationWorkHours.Columns.OrganizationID)),
					)
					q.Apply(sm.GroupBy(models.Organizations.Columns.ID))
					q.Apply(sm.GroupBy(models.OrganizationContacts.Columns.ID))
					q.Apply(sm.GroupBy(models.Addresses.Columns.ID))
					q.Apply(sm.GroupBy(models.Countries.Columns.ID))
					q.Apply(sm.GroupBy(models.OrganizationWorkHours.Columns.ID))
				case api.ListAnimalsParamsEmbedTags:
					q.Apply(sm.Columns(
						psql.Raw(`(
	  SELECT
	    json_agg(row_to_json(tags_row))
	  FROM
	    (
	      SELECT
	        animal_tag.id AS "id",
	        animal_tag.name AS "name",
	        animal_tag.created_at AS "createdAt",
	        animal_tag.updated_at AS "updatedAt"
	      FROM
	        public.animal_tag
	      WHERE
	        animal_tag.animal_id = animal.id
	    ) AS "tags_row"
	) AS "tags"`),
					))

				case api.ListAnimalsParamsEmbedPhotos:
					q.Apply(sm.Columns(
						psql.Raw(`(
	  SELECT
	    json_agg(row_to_json(photos_row))
	  FROM
	    (
	      SELECT
	        animal_photo.id AS "id",
	        animal_photo.animal_id AS "animalID",
	        animal_photo.object_kind AS "objectKind",
	        animal_photo.object_ref_small AS "objectRefSmall",
	        animal_photo.object_ref_medium AS "objectRefMedium",
	        animal_photo.object_ref_large AS "objectRefLarge",
	        animal_photo.object_ref_full AS "objectRefFull",
	        animal_photo.created_at AS "createdAt",
	        animal_photo.updated_at AS "updatedAt"
	      FROM
	        public.animal_photo
	      WHERE
	        animal_photo.animal_id = animal.id
	    ) AS "photos_row"
	) AS "photos"`),
					))

				case api.ListAnimalsParamsEmbedVideos:
					q.Apply(sm.Columns(
						psql.Raw(`(
	  SELECT
	    json_agg(row_to_json(videos_row))
	  FROM
	    (
	      SELECT
	        animal_video.id AS "id",
	        animal_video.animal_id AS "animalID",
	        animal_video.object_kind AS "objectKind",
	        animal_video.object_ref AS "objectRef",
	        animal_video.created_at AS "createdAt",
	        animal_video.updated_at AS "updatedAt"
	      FROM
	        public.animal_video
	      WHERE
	        animal_video.animal_id = animal.id
	    ) AS "videos_row"
	) AS "videos"`),
					))
				}
			}
		}

		if filters.AnimalTypeID != nil {
			typeIDs := make([]any, len(*filters.AnimalTypeID))
			for i, id := range *filters.AnimalTypeID {
				typeIDs[i] = id
			}

			q.Apply(sm.Where(models.AnimalTypes.Columns.ID.In(psql.Arg(typeIDs...))))
		}

		if filters.AnimalSpecieID != nil {
			specieIDs := make([]any, len(*filters.AnimalSpecieID))
			for i, id := range *filters.AnimalSpecieID {
				specieIDs[i] = id
			}

			q.Apply(sm.Where(models.AnimalSpecies.Columns.ID.In(psql.Arg(specieIDs...))))
		}

		if filters.AnimalBreedID != nil {
			breedIDs := make([]any, len(*filters.AnimalBreedID))
			for i, id := range *filters.AnimalBreedID {
				breedIDs[i] = id
			}

			q.Apply(sm.Where(psql.Raw(`
					EXISTS (
					select 1 from
					animal_breed
					where animal_breed.animal_id = animal.id and
						    animal_breed.breed_id in ?
					)
					`, psql.ArgGroup(breedIDs...))))
		}

		if filters.Name != nil {
			q.Apply(sm.Where(models.Animals.Columns.Name.ILike(psql.Arg("%" + *filters.Name + "%"))))
		}

		if filters.Age != nil {
			ages := make([]any, len(*filters.Age))
			for i, age := range *filters.Age {
				ages[i] = age
			}

			q.Apply(sm.Where(models.Animals.Columns.Age.In(psql.Arg(ages...))))
		}

		if filters.Size != nil {
			sizes := make([]any, len(*filters.Size))
			for i, age := range *filters.Size {
				sizes[i] = age
			}

			q.Apply(sm.Where(models.Animals.Columns.Size.In(psql.Arg(sizes...))))
		}

		if filters.Gender != nil {
			genders := make([]any, len(*filters.Gender))
			for i, age := range *filters.Gender {
				genders[i] = age
			}

			q.Apply(sm.Where(models.Animals.Columns.Gender.In(psql.Arg(genders...))))
		}

		if filters.Hermaphrodite != nil {
			q.Apply(sm.Where(models.Animals.Columns.Hermaphrodite.EQ(psql.Arg(*filters.Hermaphrodite))))
		}

		if filters.Microchip != nil {
			if filters.Embed == nil || (filters.Embed != nil && !slices.Contains(*filters.Embed, api.ListAnimalsParamsEmbedMicrochip)) {
				q.Apply(sm.LeftJoin(models.Microchips.Name()).
					On(models.Animals.Columns.ID.EQ(models.Microchips.Columns.AnimalID)))
			}

			if *filters.Microchip {
				q.Apply(sm.Where(models.Microchips.Columns.ID.IsNotNull()))
			} else {
				q.Apply(sm.Where(models.Microchips.Columns.ID.IsNull()))
			}
		}

		if filters.Tag != nil {
			// could probably use `any array` or `unnest` but whatever...
			orTagLike := []bob.Expression{psql.Raw("false")}
			for _, tag := range *filters.Tag {
				orTagLike = append(orTagLike, models.AnimalTags.Columns.Name.ILike(psql.Arg("%"+tag+"%")))
			}

			q.Apply(sm.Where(psql.Exists(
				psql.Select(
					sm.Columns(1),
					sm.From(models.AnimalTags.Name()),
					sm.Where(models.AnimalTags.Columns.AnimalID.EQ(models.Animals.Columns.ID).And(psql.Or(orTagLike...))),
				),
			)))
		}

		if filters.DaysLt != nil {
			q.Apply(sm.Where(models.Animals.Columns.CreatedAt.GTE(psql.Raw("now() - make_interval(days => ?)", *filters.DaysLt))))
		}

		if filters.DaysGt != nil {
			q.Apply(sm.Where(models.Animals.Columns.CreatedAt.LTE(psql.Raw("now() - make_interval(days => ?)", *filters.DaysGt))))
		}

		if filters.Lat != nil && filters.Lon != nil && filters.RadiusM != nil && *filters.RadiusM != 0 {
			if filters.Embed == nil || (filters.Embed != nil && !slices.Contains(*filters.Embed, api.ListAnimalsParamsEmbedOrganization)) {
				q.Apply(
					sm.LeftJoin(models.Organizations.Name()).
						On(models.Organizations.Columns.ID.EQ(models.Animals.Columns.OrganizationID)),
					sm.LeftJoin(models.OrganizationContacts.Name()).
						On(models.Organizations.Columns.ID.EQ(models.OrganizationContacts.Columns.OrganizationID)),
					sm.LeftJoin(models.Addresses.Name()).
						On(models.Addresses.Columns.ID.EQ(models.OrganizationContacts.Columns.AddressID)),
					sm.LeftJoin(models.Countries.Name()).
						On(models.Countries.Columns.ID.EQ(models.Addresses.Columns.CountryID)),
					sm.LeftJoin(models.OrganizationWorkHours.Name()).
						On(models.Organizations.Columns.ID.EQ(models.OrganizationWorkHours.Columns.OrganizationID)),
				)
			}

			q.Apply(sm.Where(psql.Raw(`ST_DWithin("address"."coords", ST_SetSRID(ST_MakePoint(?, ?), ?)::geography, ?)`,
				*filters.Lon, *filters.Lat, dbcustom.SRID, *filters.RadiusM)))
		}
	}

	if filters.Status != nil {
		statuses := make([]any, len(*filters.Status))
		for i, age := range *filters.Status {
			statuses[i] = age
		}

		q.Apply(sm.Where(models.Animals.Columns.Status.In(psql.Arg(statuses...))))
	}

	if filters.Properties != nil && len(filters.PropertiesFilters) > 0 {
		for _, f := range filters.PropertiesFilters {
			switch f.Type {
			case dbtype.FilterTypeBool:
				boolVals, ok := f.Value.([]bool)
				if !ok {
					continue
				}

				boolArgs := make([]any, len(boolVals))
				for i, x := range boolVals {
					boolArgs[i] = x
				}

				raw := fmt.Sprintf(`("animal"."properties"->>'%s')::boolean IN ?`, f.Name)
				q.Apply(sm.Where(psql.Raw(raw, psql.ArgGroup(boolArgs...))))

			case dbtype.FilterTypeString:
				strVals, ok := f.Value.([]string)
				if !ok {
					continue
				}

				strArgs := make([]any, len(strVals))
				for i, x := range strVals {
					strArgs[i] = x
				}

				raw := fmt.Sprintf(`"animal"."properties"->>'%s' IN ?`, f.Name)
				q.Apply(sm.Where(psql.Raw(raw, psql.ArgGroup(strArgs...))))
			}
		}
	}

	type ListAnimalsRow struct {
		dbtype.AnimalWithJoinData
		TotalCount int64
	}

	animals, err := bob.All(ctx, pst.exec, q, scan.StructMapper[ListAnimalsRow](scan.WithTypeConverter(orm.NullTypeConverter{})))
	if err != nil {
		return dbtype.PagedResult[dbtype.AnimalWithJoinData]{}, fmt.Errorf("query animals")
	}

	result := dbtype.PagedResult[dbtype.AnimalWithJoinData]{
		Data: make([]dbtype.AnimalWithJoinData, len(animals)),
	}
	for i, row := range animals {
		result.Data[i] = row.AnimalWithJoinData
	}

	if len(animals) > 0 {
		result.TotalCount = animals[0].TotalCount
	}

	return result, nil
}

func (pst *PgAnimalPersistor) GetAnimalByID(ctx context.Context, animalID int64, filters dbtype.GetAnimalByIDFilters) (dbtype.AnimalWithJoinData, error) {
	q := psql.Select(
		sm.Columns(
			models.Animals.Columns,
			models.AnimalTypes.Columns.WithPrefix("type."),
			models.AnimalSpecies.Columns.WithPrefix("specie."),
			psql.Raw(`COUNT("user_animal_like"."animal_id") as "likes"`),
			psql.Select(
				sm.Columns(models.Adoptions.Columns.ID.As("adoption_id")),
				sm.From(models.Adoptions.Name()),
				sm.Where(models.Adoptions.Columns.AnimalID.EQ(models.Animals.Columns.ID)),
			),
		),
		sm.From(models.Animals.Name()),
		sm.LeftJoin(models.AnimalTypes.Name()).
			On(models.Animals.Columns.TypeID.EQ(models.AnimalTypes.Columns.ID)),
		sm.LeftJoin(models.AnimalSpecies.Name()).
			On(models.Animals.Columns.SpecieID.EQ(models.AnimalSpecies.Columns.ID)),
		sm.LeftJoin(models.UserAnimalLikes.Name()).
			On(models.UserAnimalLikes.Columns.AnimalID.EQ(models.Animals.Columns.ID)),
		sm.Where(models.Animals.Columns.ID.EQ(psql.Arg(animalID))),
		sm.GroupBy(models.Animals.Columns.ID),
		sm.GroupBy(models.AnimalTypes.Columns.ID),
		sm.GroupBy(models.AnimalSpecies.Columns.ID),
	)

	if filters.UserID != nil {
		q.Apply(sm.Columns(psql.Exists(psql.Select(
			sm.Columns(1),
			sm.From(models.UserAnimalLikes.Name()),
			sm.Where(models.UserAnimalLikes.Columns.AnimalID.EQ(models.Animals.Columns.ID).
				And(models.UserAnimalLikes.Columns.UserID.EQ(psql.Arg(*filters.UserID))),
			),
		)).As("liked")))
	}

	if hasAnyLogicFilters(&filters.GetAnimalParams) {
		if filters.Embed != nil {
			for _, embed := range *filters.Embed {
				switch embed {
				case api.GetAnimalParamsEmbedBreeds:
					q.Apply(sm.Columns(
						psql.Raw(`(
  SELECT
    json_agg(row_to_json(breeds_row))
  FROM
    (
      SELECT
     	  breed.id AS "id",
        breed.name AS "name",
        animal_breed.primary AS "primary",
        animal_breed.created_at AS "createdAt",
        animal_breed.updated_at AS "updatedAt"
      FROM
        public.animal_breed
        INNER JOIN public.breed ON animal_breed.breed_id = breed.id
      WHERE
        animal_breed.animal_id = animal.id
    ) AS breeds_row
) AS "breeds"`),
					))
				case api.GetAnimalParamsEmbedMicrochip:
					q.Apply(sm.Columns(
						models.Microchips.Columns.WithPrefix("microchip."),
					))
					q.Apply(
						sm.LeftJoin(models.Microchips.Name()).
							On(models.Animals.Columns.ID.EQ(models.Microchips.Columns.AnimalID)),
					)
					q.Apply(sm.GroupBy(models.Microchips.Columns.ID))

				case api.GetAnimalParamsEmbedOrganization:
					q.Apply(sm.Columns(
						models.Organizations.Columns.WithPrefix("organization."),
						models.OrganizationContacts.Columns.WithPrefix("organization_contact."),
						models.OrganizationWorkHours.Columns.WithPrefix("organization_work_hour."),
						models.Addresses.Columns.Except("coords").WithPrefix("organization_contact_address."),
						psql.Raw(`ST_AsBinary("address"."coords") AS "organization_contact_address.coords"`),
						models.Countries.Columns.WithPrefix("organization_contact_address_country."),
					))
					q.Apply(
						sm.LeftJoin(models.Organizations.Name()).
							On(models.Organizations.Columns.ID.EQ(models.Animals.Columns.OrganizationID)),
						sm.LeftJoin(models.OrganizationContacts.Name()).
							On(models.Organizations.Columns.ID.EQ(models.OrganizationContacts.Columns.OrganizationID)),
						sm.LeftJoin(models.Addresses.Name()).
							On(models.Addresses.Columns.ID.EQ(models.OrganizationContacts.Columns.AddressID)),
						sm.LeftJoin(models.Countries.Name()).
							On(models.Countries.Columns.ID.EQ(models.Addresses.Columns.CountryID)),
						sm.LeftJoin(models.OrganizationWorkHours.Name()).
							On(models.Organizations.Columns.ID.EQ(models.OrganizationWorkHours.Columns.OrganizationID)),
					)
					q.Apply(sm.GroupBy(models.Organizations.Columns.ID))
					q.Apply(sm.GroupBy(models.OrganizationContacts.Columns.ID))
					q.Apply(sm.GroupBy(models.Addresses.Columns.ID))
					q.Apply(sm.GroupBy(models.Countries.Columns.ID))
					q.Apply(sm.GroupBy(models.OrganizationWorkHours.Columns.ID))
				case api.GetAnimalParamsEmbedTags:
					q.Apply(sm.Columns(
						psql.Raw(`(
 SELECT
   json_agg(row_to_json(tags_row))
 FROM
   (
     SELECT
       animal_tag.id AS "id",
       animal_tag.name AS "name",
       animal_tag.created_at AS "createdAt",
       animal_tag.updated_at AS "updatedAt"
     FROM
       public.animal_tag
     WHERE
       animal_tag.animal_id = animal.id
   ) AS "tags_row"
) AS "tags"`),
					))

				case api.GetAnimalParamsEmbedPhotos:
					q.Apply(sm.Columns(
						psql.Raw(`(
 SELECT
   json_agg(row_to_json(photos_row))
 FROM
   (
     SELECT
       animal_photo.id AS "id",
       animal_photo.animal_id AS "animalID",
       animal_photo.object_kind AS "objectKind",
       animal_photo.object_ref_small AS "objectRefSmall",
       animal_photo.object_ref_medium AS "objectRefMedium",
       animal_photo.object_ref_large AS "objectRefLarge",
       animal_photo.object_ref_full AS "objectRefFull",
       animal_photo.created_at AS "createdAt",
       animal_photo.updated_at AS "updatedAt"
     FROM
       public.animal_photo
     WHERE
       animal_photo.animal_id = animal.id
   ) AS "photos_row"
) AS "photos"`),
					))

				case api.GetAnimalParamsEmbedVideos:
					q.Apply(sm.Columns(
						psql.Raw(`(
 SELECT
   json_agg(row_to_json(videos_row))
 FROM
   (
     SELECT
       animal_video.id AS "id",
       animal_video.animal_id AS "animalID",
       animal_video.object_kind AS "objectKind",
       animal_video.object_ref AS "objectRef",
       animal_video.created_at AS "createdAt",
       animal_video.updated_at AS "updatedAt"
     FROM
       public.animal_video
     WHERE
       animal_video.animal_id = animal.id
   ) AS "videos_row"
) AS "videos"`),
					))
				}
			}
		}
	}

	animal, err := bob.One(ctx, pst.exec, q, scan.StructMapper[dbtype.AnimalWithJoinData](scan.WithTypeConverter(orm.NullTypeConverter{})))
	if err != nil {
		return dbtype.AnimalWithJoinData{}, fmt.Errorf("query animal")
	}

	return animal, nil
}

func (pst *PgAnimalPersistor) CreateAnimal(ctx context.Context, in dbtype.AnimalCreateSetter) (models.Animal, error) {
	var insertedAnimal models.Animal

	txErr := pst.WithTx(ctx, func(tx bobpgx.Tx) error {
		q1 := models.Animals.Insert(&in.Animal, im.Returning(models.Animals.Columns))

		animal, err := bob.One(ctx, tx, q1, scan.StructMapper[models.Animal]())
		if err != nil {
			return fmt.Errorf("insert animal: %w", err)
		}

		insertedAnimal = animal

		if !in.Microchip.IsUnset() && !in.Microchip.IsNull() {
			inMicrochip := in.Microchip.MustGet()
			inMicrochip.AnimalID.Set(insertedAnimal.ID)

			q2 := models.Microchips.Insert(&inMicrochip, im.Returning(models.Microchips.Columns))
			if _, err := bob.One(ctx, tx, q2, scan.StructMapper[models.Microchip]()); err != nil {
				return fmt.Errorf("insert microchip: %w", err)
			}
		}

		if !in.Breeds.IsUnset() && !in.Breeds.IsNull() {
			inBreeds := in.Breeds.MustGet()
			if len(inBreeds) > 0 {
				setters := make([]*models.AnimalBreedSetter, len(inBreeds))
				for i, x := range inBreeds {
					x.AnimalID.Set(insertedAnimal.ID)
					setters[i] = &x
				}

				q3 := models.AnimalBreeds.Insert(bob.ToMods(setters...), im.Returning(models.AnimalBreeds.Columns))
				if _, err := bob.Exec(ctx, tx, q3); err != nil {
					return fmt.Errorf("insert animal breeds: %w", err)
				}
			}
		}

		if !in.Tags.IsUnset() && !in.Tags.IsNull() {
			inTags := in.Tags.MustGet()
			if len(inTags) > 0 {
				setters := make([]*models.AnimalTagSetter, len(inTags))
				for i, x := range inTags {
					x.AnimalID.Set(insertedAnimal.ID)
					setters[i] = &x
				}

				q4 := models.AnimalTags.Insert(bob.ToMods(setters...), im.Returning(models.AnimalTags.Columns))
				if _, err := bob.Exec(ctx, tx, q4); err != nil {
					return fmt.Errorf("insert animal tags: %w", err)
				}
			}
		}

		if !in.Photos.IsUnset() && !in.Photos.IsNull() {
			inPhotos := in.Photos.MustGet()
			if len(inPhotos) > 0 {
				setters := make([]*models.AnimalPhotoSetter, len(inPhotos))
				for i, x := range inPhotos {
					x.AnimalID.Set(insertedAnimal.ID)
					setters[i] = &x
				}

				q5 := models.AnimalPhotos.Insert(bob.ToMods(setters...), im.Returning(models.AnimalPhotos.Columns))
				if _, err := bob.Exec(ctx, tx, q5); err != nil {
					return fmt.Errorf("insert animal photos: %w", err)
				}
			}
		}

		if !in.Videos.IsUnset() && !in.Videos.IsNull() {
			inVideos := in.Videos.MustGet()
			if len(inVideos) > 0 {
				setters := make([]*models.AnimalVideoSetter, len(inVideos))
				for i, x := range inVideos {
					x.AnimalID.Set(insertedAnimal.ID)
					setters[i] = &x
				}

				q6 := models.AnimalVideos.Insert(bob.ToMods(setters...), im.Returning(models.AnimalVideos.Columns))
				if _, err := bob.Exec(ctx, tx, q6); err != nil {
					return fmt.Errorf("insert animal videos: %w", err)
				}
			}
		}

		return nil
	})

	return insertedAnimal, txErr
}

func (pst *PgAnimalPersistor) UpdateAnimal(ctx context.Context, animalID int64, in dbtype.AnimalUpdateSetter) (models.Animal, error) {
	q := psql.Select(
		sm.Columns(models.Microchips.Columns.ID.As("microchip_id")),
		sm.From(models.Microchips.Name()),
		sm.Where(models.Microchips.Columns.AnimalID.EQ(psql.Arg(animalID))),
	)

	type idMicrochipResult struct {
		MicrochipID *int64
	}

	idMicrochip, err := bob.One(ctx, pst.exec, q, scan.StructMapper[idMicrochipResult]())
	if err != nil && !errors.Is(pgx.ErrNoRows, err) {
		return models.Animal{}, fmt.Errorf("query microchip id: %w", err)
	}

	var updatedAnimal models.Animal

	txErr := pst.WithTx(ctx, func(tx bobpgx.Tx) error {
		if !in.Animal.IsUnset() && !in.Animal.IsNull() {
			inAnimal := in.Animal.MustGet()
			q1 := models.Animals.Update(
				inAnimal.UpdateMod(),
				um.Where(models.Animals.Columns.ID.EQ(psql.Arg(animalID))),
				um.Returning(models.Animals.Columns),
			)

			animal, err := bob.One(ctx, tx, q1, scan.StructMapper[models.Animal]())
			if err != nil {
				return fmt.Errorf("update animal: %w", err)
			}

			updatedAnimal = animal
		}

		if !in.Microchip.IsUnset() {
			if in.Microchip.IsNull() {
				if idMicrochip.MicrochipID != nil {
					q2 := models.Microchips.Delete(dm.Where(models.Microchips.Columns.ID.EQ(psql.Arg(*idMicrochip.MicrochipID))))
					if _, err := bob.Exec(ctx, tx, q2); err != nil {
						return fmt.Errorf("delete animal microchip: %w", err)
					}
				}
			} else {
				inMicrochip := in.Microchip.MustGet()
				inMicrochip.AnimalID.Set(animalID)

				var isPatch bool
				if idMicrochip.MicrochipID != nil {
					isPatch = true
				}

				if isPatch {
					q3 := models.Microchips.Update(
						inMicrochip.UpdateMod(),
						um.Where(models.Microchips.Columns.ID.EQ(psql.Arg(*idMicrochip.MicrochipID))),
					)
					if _, err := bob.Exec(ctx, tx, q3); err != nil {
						return fmt.Errorf("update animal microchip: %w", err)
					}
				} else {
					q4 := models.Microchips.Insert(&inMicrochip)
					if _, err := bob.Exec(ctx, tx, q4); err != nil {
						return fmt.Errorf("insert animal microchip: %w", err)
					}
				}
			}
		}

		if !in.Photos.IsUnset() && !in.Photos.IsNull() {
			inPhotos := in.Photos.MustGet()
			if len(inPhotos) > 0 {
				setters := make([]*models.AnimalPhotoSetter, len(inPhotos))
				for i, x := range inPhotos {
					x.AnimalID.Set(updatedAnimal.ID)
					setters[i] = &x
				}

				q5 := models.AnimalPhotos.Insert(bob.ToMods(setters...), im.Returning(models.AnimalPhotos.Columns))
				if _, err := bob.Exec(ctx, tx, q5); err != nil {
					return fmt.Errorf("insert animal photos: %w", err)
				}
			}
		}

		if !in.Videos.IsUnset() && !in.Videos.IsNull() {
			inVideos := in.Videos.MustGet()
			if len(inVideos) > 0 {
				setters := make([]*models.AnimalVideoSetter, len(inVideos))
				for i, x := range inVideos {
					x.AnimalID.Set(updatedAnimal.ID)
					setters[i] = &x
				}

				q6 := models.AnimalVideos.Insert(bob.ToMods(setters...), im.Returning(models.AnimalVideos.Columns))
				if _, err := bob.Exec(ctx, tx, q6); err != nil {
					return fmt.Errorf("insert animal videos: %w", err)
				}
			}
		}

		if !in.Breeds.IsUnset() {
			deleteOld := in.Breeds.IsNull() || (!in.Breeds.IsNull() && len(in.Breeds.MustGet()) > 0)
			if deleteOld {
				q7 := models.AnimalBreeds.Delete(dm.Where(models.AnimalBreeds.Columns.AnimalID.EQ(psql.Arg(animalID))))
				if _, err := bob.Exec(ctx, tx, q7); err != nil {
					return fmt.Errorf("delete old animal breeds: %w", err)
				}
			}

			if !in.Breeds.IsUnset() && len(in.Breeds.MustGet()) > 0 {
				inBreeds := in.Breeds.MustGet()
				if len(inBreeds) > 0 {
					setters := make([]*models.AnimalBreedSetter, len(inBreeds))
					for i, x := range inBreeds {
						x.AnimalID.Set(updatedAnimal.ID)
						setters[i] = &x
					}
					// @TODO: should probably do upsert but i am lazy,
					// i just delete old and add new, no need to fetch old, diff, upsert, delete etc...
					q8 := models.AnimalBreeds.Insert(bob.ToMods(setters...), im.Returning(models.AnimalBreeds.Columns))
					if _, err := bob.Exec(ctx, tx, q8); err != nil {
						return fmt.Errorf("insert animal breeds: %w", err)
					}
				}
			}
		}

		if !in.Tags.IsUnset() {
			deleteOld := in.Tags.IsNull() || (!in.Tags.IsNull() && len(in.Tags.MustGet()) > 0)
			if deleteOld {
				q9 := models.AnimalTags.Delete(dm.Where(models.AnimalTags.Columns.AnimalID.EQ(psql.Arg(animalID))))
				if _, err := bob.Exec(ctx, tx, q9); err != nil {
					return fmt.Errorf("delete old animal tags: %w", err)
				}
			}

			if !in.Tags.IsUnset() && len(in.Tags.MustGet()) > 0 {
				inTags := in.Tags.MustGet()
				if len(inTags) > 0 {
					setters := make([]*models.AnimalTagSetter, len(inTags))
					for i, x := range inTags {
						x.AnimalID.Set(updatedAnimal.ID)
						setters[i] = &x
					}
					// @TODO: should probably do upsert but i am lazy,
					// i just delete old and add new, no need to fetch old, diff, upsert, delete etc...
					q8 := models.AnimalTags.Insert(bob.ToMods(setters...), im.Returning(models.AnimalTags.Columns))
					if _, err := bob.Exec(ctx, tx, q8); err != nil {
						return fmt.Errorf("insert animal tags: %w", err)
					}
				}
			}
		}

		return nil
	})

	return updatedAnimal, txErr
}

func (pst *PgAnimalPersistor) DeleteAnimalByID(ctx context.Context, animalID int64) (int64, error) {
	q := models.Animals.Delete(dm.Where(models.Animals.Columns.ID.EQ(psql.Arg(animalID))))
	if _, err := bob.Exec(ctx, pst.exec, q); err != nil {
		return 0, fmt.Errorf("delete animal: %w", err)
	}

	return animalID, nil
}

func (pst *PgAnimalPersistor) DeleteAnimals(ctx context.Context, ids []int64) error {
	animalIDs := make([]any, len(ids))
	for i, id := range ids {
		animalIDs[i] = id
	}

	q := models.Animals.Delete(dm.Where(models.Animals.Columns.ID.In(psql.Arg(animalIDs...))))
	if _, err := bob.Exec(ctx, pst.exec, q); err != nil {
		return fmt.Errorf("delete animals: %w", err)
	}

	return nil
}

func (pst *PgAnimalPersistor) LikeAnimal(ctx context.Context, userID uuid.UUID, animalID int64) error {
	q := models.UserAnimalLikes.Insert(&models.UserAnimalLikeSetter{UserID: omit.From(userID), AnimalID: omit.From(animalID)})
	if _, err := bob.Exec(ctx, pst.exec, q); err != nil {
		return fmt.Errorf("insert animal like: %w", err)
	}

	return nil
}

func (pst *PgAnimalPersistor) UnlikeAnimal(ctx context.Context, userID uuid.UUID, animalID int64) error {
	q := models.UserAnimalLikes.Delete(dm.Where(
		models.UserAnimalLikes.Columns.UserID.EQ(psql.Arg(userID)).
			And(models.UserAnimalLikes.Columns.AnimalID.EQ(psql.Arg(animalID))),
	))
	if _, err := bob.Exec(ctx, pst.exec, q); err != nil {
		return fmt.Errorf("delete animal like: %w", err)
	}

	return nil
}

func (pst *PgAnimalPersistor) ApplyForAdoption(ctx context.Context, animalID int64, userID uuid.UUID, organizationID *int64) (models.Adoption, error) {
	adoptionSetter := &models.AdoptionSetter{
		AnimalID:       omit.From(animalID),
		UserID:         omit.From(userID),
		OrganizationID: omitnull.FromPtr(organizationID),
		Status:         omit.From("pending"),
	}

	animalSetter := &models.AnimalSetter{
		Status: omit.From("reserved"),
	}

	var appliedAdoption models.Adoption

	txErr := pst.WithTx(ctx, func(tx bobpgx.Tx) error {
		q1 := models.Adoptions.Insert(adoptionSetter, im.Returning(models.Adoptions.Columns))
		adoption, err := bob.One(ctx, pst.exec, q1, scan.StructMapper[models.Adoption]())
		if err != nil {
			return fmt.Errorf("insert animal pending adoption: %w", err)
		}
		appliedAdoption = adoption

		q2 := models.Animals.Update(
			animalSetter.UpdateMod(),
			um.Where(models.Animals.Columns.ID.EQ(psql.Arg(animalID)).And(models.Animals.Columns.Status.EQ(psql.Arg("adoptable")))),
		)
		if _, err := bob.Exec(ctx, pst.exec, q2); err != nil {
			return fmt.Errorf("update animal adoption reserved status: %w", err)
		}

		return nil
	})

	return appliedAdoption, txErr
}

func (pst *PgAnimalPersistor) ApproveAdoption(ctx context.Context, adoptionID int64) error {
	animalSetter := &models.AnimalSetter{
		Status: omit.From("adopted"),
	}

	adoptionSetter := &models.AdoptionSetter{
		Status: omit.From("approved"),
	}

	txErr := pst.WithTx(ctx, func(tx bobpgx.Tx) error {
		q1 := models.Animals.Update(
			animalSetter.UpdateMod(),
			um.From(models.Adoptions.Name()),
			um.Where(models.Adoptions.Columns.ID.EQ(psql.Arg(adoptionID)).
				And(models.Adoptions.Columns.AnimalID.EQ(models.Animals.Columns.ID).
					And(models.Animals.Columns.Status.EQ(psql.Arg("reserved"))),
				),
			),
		)
		if _, err := bob.Exec(ctx, pst.exec, q1); err != nil {
			return fmt.Errorf("update animal adoption approved status: %w", err)
		}

		q2 := models.Adoptions.Update(
			adoptionSetter.UpdateMod(),
			um.Where(models.Adoptions.Columns.ID.EQ(psql.Arg(adoptionID))),
		)
		if _, err := bob.Exec(ctx, pst.exec, q2); err != nil {
			return fmt.Errorf("update adoption approved status: %w", err)
		}

		return nil
	})

	return txErr
}

func (pst *PgAnimalPersistor) RejectAdoption(ctx context.Context, adoptionID int64) error {
	animalSetter := &models.AnimalSetter{
		Status: omit.From("adoptable"),
	}

	adoptionSetter := &models.AdoptionSetter{
		Status: omit.From("rejected"),
	}

	txErr := pst.WithTx(ctx, func(tx bobpgx.Tx) error {
		q1 := models.Animals.Update(
			animalSetter.UpdateMod(),
			um.From(models.Adoptions.Name()),
			um.Where(models.Adoptions.Columns.ID.EQ(psql.Arg(adoptionID)).
				And(models.Adoptions.Columns.AnimalID.EQ(models.Animals.Columns.ID).
					And(models.Animals.Columns.Status.EQ(psql.Arg("reserved"))),
				),
			),
		)
		if _, err := bob.Exec(ctx, pst.exec, q1); err != nil {
			return fmt.Errorf("update animal adoption rejected status: %w", err)
		}

		q2 := models.Adoptions.Update(
			adoptionSetter.UpdateMod(),
			um.Where(models.Adoptions.Columns.ID.EQ(psql.Arg(adoptionID))),
		)
		if _, err := bob.Exec(ctx, pst.exec, q2); err != nil {
			return fmt.Errorf("update adoption rejected status: %w", err)
		}

		return nil
	})

	return txErr
}

func (pst *PgAnimalPersistor) ApproveAnimal(ctx context.Context, animalID int64) error {
	q := models.Animals.Update(
		models.AnimalSetter{Status: omit.From("adoptable")}.UpdateMod(),
		um.Where(models.Animals.Columns.ID.EQ(psql.Arg(animalID))),
		um.Returning(models.Animals.Columns),
	)
	if _, err := bob.Exec(ctx, pst.exec, q); err != nil {
		return fmt.Errorf("update animal status to adoptable: %w", err)
	}

	return nil
}

func (pst *PgAnimalPersistor) RejectAnimal(ctx context.Context, animalID int64) error {
	q := models.Animals.Update(
		models.AnimalSetter{Status: omit.From("rejected")}.UpdateMod(),
		um.Where(models.Animals.Columns.ID.EQ(psql.Arg(animalID))),
		um.Returning(models.Animals.Columns),
	)
	if _, err := bob.Exec(ctx, pst.exec, q); err != nil {
		return fmt.Errorf("update animal status to rejected: %w", err)
	}

	return nil
}
