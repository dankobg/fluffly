package postgres

import (
	"context"
	"fmt"
	"slices"
	"strings"

	api "github.com/dankobg/fluffly/api/gen"
	"github.com/dankobg/fluffly/db/gen/models"
	"github.com/dankobg/fluffly/persistence"
	"github.com/dankobg/fluffly/persistence/dbtype"
	"github.com/google/uuid"
	"github.com/stephenafamo/bob"
	"github.com/stephenafamo/bob/dialect/psql"
	"github.com/stephenafamo/bob/dialect/psql/dm"
	"github.com/stephenafamo/bob/dialect/psql/im"
	"github.com/stephenafamo/bob/dialect/psql/sm"
	"github.com/stephenafamo/bob/dialect/psql/um"
	"github.com/stephenafamo/bob/orm"
	"github.com/stephenafamo/scan"
)

var _ persistence.UserPersistor = (*PgUserPersistor)(nil)

type PgUserPersistor struct {
	*PgPersistor
}

func NewPgUserPersistor(ps *PgPersistor) *PgUserPersistor {
	return &PgUserPersistor{
		PgPersistor: ps,
	}
}

func (pst *PgUserPersistor) ListUsers(ctx context.Context, filters dbtype.ListUsersFilters) (dbtype.PagedResult[models.User], error) {
	q := psql.Select(
		sm.Columns(models.Users.Columns),
		sm.From(models.Users.Name()),
		sm.GroupBy(models.Users.Columns.ID),
	)
	addOrderBy(&q, filters.Sort, models.Users.Columns.Names())
	addPagination(&q, filters.Page, filters.PageSize)

	type ListUsersRow struct {
		models.User
		TotalCount int64
	}

	countries, err := bob.All(ctx, pst.exec, q, scan.StructMapper[ListUsersRow]())
	if err != nil {
		return dbtype.PagedResult[models.User]{}, fmt.Errorf("query users")
	}

	result := dbtype.PagedResult[models.User]{
		Data: make([]models.User, len(countries)),
	}
	for i, row := range countries {
		result.Data[i] = row.User
	}

	if len(countries) > 0 {
		result.TotalCount = countries[0].TotalCount
	}

	return result, nil
}

func (pst *PgUserPersistor) GetUserByID(ctx context.Context, userID uuid.UUID) (models.User, error) {
	q := psql.Select(
		sm.Columns(models.Users.Columns),
		sm.From(models.Users.Name()),
		sm.Where(models.Users.Columns.ID.EQ(psql.Arg(userID))),
	)

	user, err := bob.One(ctx, pst.exec, q, scan.StructMapper[models.User]())
	if err != nil {
		return models.User{}, fmt.Errorf("query user")
	}

	return user, nil
}

func (pst *PgUserPersistor) DeleteUserByID(ctx context.Context, userID uuid.UUID) (uuid.UUID, error) {
	q := models.Users.Delete(dm.Where(models.Users.Columns.ID.EQ(psql.Arg(userID))))
	if _, err := bob.Exec(ctx, pst.exec, q); err != nil {
		return uuid.Nil, fmt.Errorf("delete country: %w", err)
	}

	return userID, nil
}

func (pst *PgUserPersistor) CreateUser(ctx context.Context, in models.UserSetter) (models.User, error) {
	q := models.Users.Insert(&in, im.Returning(models.Users.Columns))

	user, err := bob.One(ctx, pst.exec, q, scan.StructMapper[models.User]())
	if err != nil {
		return models.User{}, fmt.Errorf("insert user")
	}

	return user, nil
}

func (pst *PgUserPersistor) UpdateUser(ctx context.Context, userID uuid.UUID, in models.UserSetter) (models.User, error) {
	q := models.Users.Update(
		in.UpdateMod(),
		um.Where(models.Users.Columns.ID.EQ(psql.Arg(userID))),
		um.Returning(models.Users.Columns),
	)

	user, err := bob.One(ctx, pst.exec, q, scan.StructMapper[models.User]())
	if err != nil {
		return models.User{}, fmt.Errorf("update user")
	}

	return user, nil
}

func (pst *PgUserPersistor) ListMyAnimals(ctx context.Context, userID uuid.UUID, filters dbtype.ListMyAnimalsFilters) (dbtype.PagedResult[dbtype.AnimalWithJoinData], error) {
	q := psql.Select(
		sm.Columns(
			models.Animals.Columns,
			models.AnimalTypes.Columns.WithPrefix("type."),
			models.AnimalSpecies.Columns.WithPrefix("specie."),
			psql.Raw(`COUNT("user_animal_like"."animal_id") as "likes"`),
			psql.Exists(psql.Select(
				sm.Columns(1),
				sm.From(models.UserAnimalLikes.Name()),
				sm.Where(models.UserAnimalLikes.Columns.AnimalID.EQ(models.Animals.Columns.ID).
					And(models.UserAnimalLikes.Columns.UserID.EQ(psql.Arg(userID))),
				),
			)).As("liked"),
			psql.Select(
				sm.Columns(models.Adoptions.Columns.ID.As("adoption_id")),
				sm.From(models.Adoptions.Name()),
				sm.Where(models.Adoptions.Columns.AnimalID.EQ(models.Animals.Columns.ID)),
			),
		),
		sm.From(models.Animals.Name()),
		sm.Where(models.Animals.Columns.UserID.EQ(psql.Arg(userID))),
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

	if hasAnyLogicFilters(&filters.ListMyAnimalsParams) {
		if filters.Embed != nil {
			for _, embed := range *filters.Embed {
				switch embed {
				case api.ListMyAnimalsParamsEmbedBreeds:
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
				case api.ListMyAnimalsParamsEmbedMicrochip:
					q.Apply(sm.Columns(
						models.Microchips.Columns.WithPrefix("microchip."),
					))
					q.Apply(
						sm.LeftJoin(models.Microchips.Name()).
							On(models.Animals.Columns.ID.EQ(models.Microchips.Columns.AnimalID)),
					)
					q.Apply(sm.GroupBy(models.Microchips.Columns.ID))
				case api.ListMyAnimalsParamsEmbedTags:
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

				case api.ListMyAnimalsParamsEmbedPhotos:
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

				case api.ListMyAnimalsParamsEmbedVideos:
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

func (pst *PgUserPersistor) ListMyFavoriteAnimals(ctx context.Context, userID uuid.UUID, filters dbtype.ListMyFavoriteAnimalsFilters) (dbtype.PagedResult[dbtype.AnimalWithJoinData], error) {
	q := psql.Select(
		sm.Columns(
			models.Animals.Columns,
			models.AnimalTypes.Columns.WithPrefix("type."),
			models.AnimalSpecies.Columns.WithPrefix("specie."),
			psql.Raw(`COUNT("user_animal_like"."animal_id") as "likes"`),
			psql.Exists(psql.Select(
				sm.Columns(1),
				sm.From(models.UserAnimalLikes.Name()),
				sm.Where(models.UserAnimalLikes.Columns.AnimalID.EQ(models.Animals.Columns.ID).
					And(models.UserAnimalLikes.Columns.UserID.EQ(psql.Arg(userID))),
				),
			)).As("liked"),
			psql.Select(
				sm.Columns(models.Adoptions.Columns.ID.As("adoption_id")),
				sm.From(models.Adoptions.Name()),
				sm.Where(models.Adoptions.Columns.AnimalID.EQ(models.Animals.Columns.ID)),
			),
		),
		sm.From(models.Animals.Name()),
		sm.Where(
			models.Animals.Columns.ID.In(
				psql.Select(
					sm.Columns(models.UserAnimalLikes.Columns.AnimalID),
					sm.From(models.UserAnimalLikes.Name()),
					sm.Where(models.UserAnimalLikes.Columns.UserID.EQ(psql.Arg(userID))),
				),
			),
		),
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

	if filters.Sort != nil {
		var likedAtSort string

		*filters.Sort = slices.DeleteFunc(*filters.Sort, func(x string) bool {
			if x == "liked_at" || x == "-liked_at" {
				likedAtSort = x
				return true
			}

			return false
		})

		if len(likedAtSort) > 0 {
			if strings.HasPrefix(likedAtSort, "-") {
				q.Apply(sm.OrderBy("MAX(liked_at) DESC"))
			} else {
				q.Apply(sm.OrderBy("MAX(liked_at) ASC"))
			}
		}
	}

	addOrderBy(&q, filters.Sort, append(models.Animals.Columns.Except(
		models.Animals.Columns.Properties.String(),
	).Names(), models.UserAnimalLikes.Columns.LikedAt.String()))
	addPagination(&q, filters.Page, filters.PageSize)

	if hasAnyLogicFilters(&filters.ListMyFavoriteAnimalsParams) {
		if filters.Embed != nil {
			for _, embed := range *filters.Embed {
				switch embed {
				case api.ListMyFavoriteAnimalsParamsEmbedBreeds:
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
				case api.ListMyFavoriteAnimalsParamsEmbedMicrochip:
					q.Apply(sm.Columns(
						models.Microchips.Columns.WithPrefix("microchip."),
					))
					q.Apply(
						sm.LeftJoin(models.Microchips.Name()).
							On(models.Animals.Columns.ID.EQ(models.Microchips.Columns.AnimalID)),
					)
					q.Apply(sm.GroupBy(models.Microchips.Columns.ID))
				case api.ListMyFavoriteAnimalsParamsEmbedOrganization:
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
				case api.ListMyFavoriteAnimalsParamsEmbedTags:
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

				case api.ListMyFavoriteAnimalsParamsEmbedPhotos:
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

				case api.ListMyFavoriteAnimalsParamsEmbedVideos:
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

func (pst *PgUserPersistor) ListMyAdoptions(ctx context.Context, userID uuid.UUID, filters dbtype.ListMyAdoptionsFilters) (dbtype.PagedResult[dbtype.AnimalWithJoinData], error) {
	q := psql.Select(
		sm.Columns(
			models.Animals.Columns,
			models.AnimalTypes.Columns.WithPrefix("type."),
			models.AnimalSpecies.Columns.WithPrefix("specie."),
			psql.Raw(`COUNT("user_animal_like"."animal_id") as "likes"`),
			psql.Exists(psql.Select(
				sm.Columns(1),
				sm.From(models.UserAnimalLikes.Name()),
				sm.Where(models.UserAnimalLikes.Columns.AnimalID.EQ(models.Animals.Columns.ID).
					And(models.UserAnimalLikes.Columns.UserID.EQ(psql.Arg(userID))),
				),
			)).As("liked"),
			psql.Select(
				sm.Columns(models.Adoptions.Columns.ID.As("adoption_id")),
				sm.From(models.Adoptions.Name()),
				sm.Where(models.Adoptions.Columns.AnimalID.EQ(models.Animals.Columns.ID)),
			),
		),
		sm.From(models.Animals.Name()),
		sm.Where(models.Adoptions.Columns.UserID.EQ(psql.Arg(userID))),
		sm.LeftJoin(models.AnimalTypes.Name()).
			On(models.Animals.Columns.TypeID.EQ(models.AnimalTypes.Columns.ID)),
		sm.LeftJoin(models.AnimalSpecies.Name()).
			On(models.Animals.Columns.SpecieID.EQ(models.AnimalSpecies.Columns.ID)),
		sm.LeftJoin(models.UserAnimalLikes.Name()).
			On(models.UserAnimalLikes.Columns.AnimalID.EQ(models.Animals.Columns.ID)),
		sm.LeftJoin(models.Adoptions.Name()).
			On(models.Adoptions.Columns.AnimalID.EQ(models.Animals.Columns.ID)),
		sm.GroupBy(models.Animals.Columns.ID),
		sm.GroupBy(models.AnimalTypes.Columns.ID),
		sm.GroupBy(models.AnimalSpecies.Columns.ID),
	)
	addOrderBy(&q, filters.Sort, models.Animals.Columns.Except(
		models.Animals.Columns.Properties.String(),
	).Names())
	addPagination(&q, filters.Page, filters.PageSize)

	if hasAnyLogicFilters(&filters.ListMyAdoptionsParams) {
		if filters.Embed != nil {
			for _, embed := range *filters.Embed {
				switch embed {
				case api.ListMyAdoptionsParamsEmbedBreeds:
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
				case api.ListMyAdoptionsParamsEmbedMicrochip:
					q.Apply(sm.Columns(
						models.Microchips.Columns.WithPrefix("microchip."),
					))
					q.Apply(
						sm.LeftJoin(models.Microchips.Name()).
							On(models.Animals.Columns.ID.EQ(models.Microchips.Columns.AnimalID)),
					)
					q.Apply(sm.GroupBy(models.Microchips.Columns.ID))
				case api.ListMyAdoptionsParamsEmbedTags:
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

				case api.ListMyAdoptionsParamsEmbedPhotos:
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

				case api.ListMyAdoptionsParamsEmbedVideos:
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

func (pst *PgUserPersistor) ListMyOrganizations(ctx context.Context, userID uuid.UUID, filters dbtype.ListMyOrganizationsFilters) (dbtype.PagedResult[dbtype.OrganizationWithJoinData], error) {
	q := psql.Select(
		sm.Columns(models.Organizations.Columns),
		sm.From(models.Organizations.Name()),
		sm.LeftJoin(models.OrganizationMemberships.Name()).
			On(models.Organizations.Columns.ID.EQ(models.OrganizationMemberships.Columns.OrganizationID)),
		sm.Where(models.OrganizationMemberships.Columns.UserID.EQ(psql.Arg(userID))),
		sm.GroupBy(models.Organizations.Columns.ID),
	)
	addOrderBy(&q, filters.Sort, models.Organizations.Columns.Names())
	addPagination(&q, filters.Page, filters.PageSize)

	if hasAnyLogicFilters(&filters.ListMyOrganizationsParams) {
		if filters.Embed != nil {
			for _, embed := range *filters.Embed {
				switch embed {
				case api.ListMyOrganizationsParamsEmbedContact:
					q.Apply(sm.Columns(
						models.OrganizationContacts.Columns.WithPrefix("contact."),
						models.Addresses.Columns.Except("coords").WithPrefix("contact_address."),
						psql.Raw(`ST_AsBinary("address"."coords") AS "contact_address.coords"`),
						models.Countries.Columns.WithPrefix("contact_address_country."),
					))
					q.Apply(
						sm.LeftJoin(models.OrganizationContacts.Name()).
							On(models.Organizations.Columns.ID.EQ(models.OrganizationContacts.Columns.OrganizationID)),
						sm.LeftJoin(models.Addresses.Name()).
							On(models.Addresses.Columns.ID.EQ(models.OrganizationContacts.Columns.AddressID)),
						sm.LeftJoin(models.Countries.Name()).
							On(models.Countries.Columns.ID.EQ(models.Addresses.Columns.CountryID)),
					)
					q.Apply(sm.GroupBy(models.OrganizationContacts.Columns.ID))
					q.Apply(sm.GroupBy(models.Addresses.Columns.ID))
					q.Apply(sm.GroupBy(models.Countries.Columns.ID))

				case api.ListMyOrganizationsParamsEmbedWorkHour:
					q.Apply(sm.Columns(
						models.OrganizationWorkHours.Columns.WithPrefix("work_hour."),
					))
					q.Apply(
						sm.LeftJoin(models.OrganizationWorkHours.Name()).
							On(models.Organizations.Columns.ID.EQ(models.OrganizationWorkHours.Columns.OrganizationID)),
					)
					q.Apply(sm.GroupBy(models.OrganizationWorkHours.Columns.ID))

				case api.ListMyOrganizationsParamsEmbedPhotos:
					q.Apply(sm.Columns(
						psql.Raw(`(
  SELECT
    json_agg(row_to_json(photos_row))
  FROM
    (
      SELECT
        organization_photo.id AS "id",
        organization_photo.organization_id AS "organizationID",
        organization_photo.object_kind AS "objectKind",
        organization_photo.object_ref_small AS "objectRefSmall",
        organization_photo.object_ref_medium AS "objectRefMedium",
        organization_photo.object_ref_large AS "objectRefLarge",
        organization_photo.object_ref_full AS "objectRefFull",
        organization_photo.created_at AS "createdAt",
        organization_photo.updated_at AS "updatedAt"
      FROM
        public.organization_photo
      WHERE
        organization_photo.organization_id = organization.id
    ) AS photos_row
) AS "photos"`),
					))

				case api.ListMyOrganizationsParamsEmbedVideos:
					q.Apply(sm.Columns(
						psql.Raw(`(
  SELECT
    json_agg(row_to_json(videos_row))
  FROM
    (
      SELECT
        organization_video.id AS "id",
        organization_video.organization_id AS "organizationID",
        organization_video.object_kind AS "objectKind",
        organization_video.object_ref AS "objectRef",
        organization_video.created_at AS "createdAt",
        organization_video.updated_at AS "updatedAt"
      FROM
        public.organization_video
      WHERE
        organization_video.organization_id = organization.id
    ) AS videos_row
) AS "videos"`),
					))

				case api.ListMyOrganizationsParamsEmbedSocials:
					q.Apply(sm.Columns(
						psql.Raw(`(
  SELECT
    json_agg(row_to_json(socials_row))
  FROM
    (
      SELECT
        organization_social.id AS "id",
        organization_social.organization_id AS "organizationID",
        organization_social.platform AS "platform",
        organization_social.url AS "url",
        organization_social.created_at AS "createdAt",
        organization_social.updated_at AS "updatedAt"
      FROM
        public.organization_social
      WHERE
        organization_social.organization_id = organization.id
    ) AS socials_row
) AS "socials"`),
					))

				}
			}
		}

		if filters.ID != nil {
			ids := make([]any, len(*filters.ID))
			for i, id := range *filters.ID {
				ids[i] = id
			}

			q.Apply(sm.Where(models.Organizations.Columns.ID.In(psql.Arg(ids...))))
		}

		if filters.Name != nil {
			q.Apply(sm.Where(models.Organizations.Columns.Name.ILike(psql.Arg("%" + *filters.Name + "%"))))
		}

		if filters.Status != nil {
			statuses := make([]any, len(*filters.Status))
			for i, id := range *filters.Status {
				statuses[i] = id
			}

			q.Apply(sm.Where(models.Organizations.Columns.Status.In(psql.Arg(statuses...))))
		}
	}

	type ListOrganizationsRow struct {
		dbtype.OrganizationWithJoinData
		TotalCount int64
	}

	organizations, err := bob.All(ctx, pst.exec, q, scan.StructMapper[ListOrganizationsRow](scan.WithTypeConverter(orm.NullTypeConverter{})))
	if err != nil {
		return dbtype.PagedResult[dbtype.OrganizationWithJoinData]{}, fmt.Errorf("query organizations")
	}

	result := dbtype.PagedResult[dbtype.OrganizationWithJoinData]{
		Data: make([]dbtype.OrganizationWithJoinData, len(organizations)),
	}
	for i, row := range organizations {
		result.Data[i] = row.OrganizationWithJoinData
	}

	if len(organizations) > 0 {
		result.TotalCount = organizations[0].TotalCount
	}

	return result, nil
}
