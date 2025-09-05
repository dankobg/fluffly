package postgres

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"github.com/dankobg/fluffly/db/gen/test/public/model"
	t "github.com/dankobg/fluffly/db/gen/test/public/table"
	"github.com/dankobg/fluffly/persistence"
	"github.com/dankobg/fluffly/persistence/dbtype"
	p "github.com/go-jet/jet/v2/postgres"
	"github.com/go-jet/jet/v2/qrm"
	"github.com/jackc/pgerrcode"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/oapi-codegen/nullable"
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

type ErrAnimalIntegrityViolation struct{ errIntegrityViolation }
type ErrAnimalUniqueViolation struct{ errUniqueViolation }
type ErrAnimalForeignKeyViolation struct{ errForeignKeyViolation }
type ErrAnimalCheckViolation struct{ errCheckViolation }

var (
	ErrAnimalNotFound  = errors.New("animal not found")
	errAnimalIntegrity = ErrAnimalIntegrityViolation{}

	errAnimalForeignKeyUserID           = ErrAnimalForeignKeyViolation{errForeignKeyViolation: errForeignKeyViolation{Name: "user_id"}}
	errAnimalForeignKeyOrganizationID   = ErrAnimalForeignKeyViolation{errForeignKeyViolation: errForeignKeyViolation{Name: "organization_id"}}
	errAnimalCheckAgeValid              = ErrAnimalCheckViolation{errCheckViolation: errCheckViolation{Name: "age valid"}}
	errAnimalCheckSizeValid             = ErrAnimalCheckViolation{errCheckViolation: errCheckViolation{Name: "size valid"}}
	errAnimalCheckUserIdOrOrgIDProvided = ErrAnimalCheckViolation{errCheckViolation: errCheckViolation{Name: "user_id_or_organization_id provided"}}
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
		}
		return errAnimalIntegrity
	}
	return pgErr
}

func (po *PgAnimalPersistor) ListAnimals(ctx context.Context, filters dbtype.AnimalFilters) (dbtype.PagedResult[dbtype.AnimalWithJoinData], error) {
	q := p.SELECT(
		t.Animal.AllColumns,
		t.AnimalType.AllColumns,
		t.AnimalSpecies.AllColumns,
		t.Microchip.AllColumns,
		t.Organization.AllColumns,
		t.OrganizationContact.AllColumns,
		t.Address.AllColumns,
		t.Country.AllColumns,
		t.OrganizationWorkHour.AllColumns,
		p.SELECT_JSON_ARR(t.AnimalBreed.AllColumns, t.Breed.AllColumns).
			FROM(
				t.AnimalBreed.
					LEFT_JOIN(t.Breed, t.AnimalBreed.BreedID.EQ(t.Breed.ID)),
			).
			WHERE(t.AnimalBreed.AnimalID.EQ(t.Animal.ID)).
			AS("breeds"),
		p.SELECT_JSON_ARR(t.Tag.AllColumns).
			FROM(t.Tag).
			WHERE(t.Tag.AnimalID.EQ(t.Animal.ID)).
			AS("tags"),
		p.SELECT_JSON_ARR(t.AnimalPhoto.AllColumns).
			FROM(t.AnimalPhoto).
			WHERE(t.AnimalPhoto.AnimalID.EQ(t.Animal.ID)).
			AS("photos"),
		p.SELECT_JSON_ARR(t.AnimalVideo.AllColumns).
			FROM(t.AnimalVideo).
			WHERE(t.AnimalVideo.AnimalID.EQ(t.Animal.ID)).
			AS("videos"),
		getSelectTotalCount(filters.Pagination),
	).
		FROM(
			t.Animal.
				LEFT_JOIN(t.AnimalType, t.Animal.ID.EQ(t.AnimalType.ID)).
				LEFT_JOIN(t.AnimalSpecies, t.Animal.ID.EQ(t.AnimalSpecies.ID)).
				LEFT_JOIN(t.Microchip, t.Animal.ID.EQ(t.Microchip.AnimalID)).
				LEFT_JOIN(t.Organization, t.Animal.OrganizationID.EQ(t.Organization.ID)).
				LEFT_JOIN(t.OrganizationContact, t.Organization.ID.EQ(t.OrganizationContact.OrganizationID)).
				LEFT_JOIN(t.Address, t.Address.ID.EQ(t.OrganizationContact.AddressID)).
				LEFT_JOIN(t.Country, t.Country.ID.EQ(t.Address.CountryID)).
				LEFT_JOIN(t.OrganizationWorkHour, t.Organization.ID.EQ(t.OrganizationWorkHour.OrganizationID)),
		).
		GROUP_BY(
			t.Animal.ID,
			t.AnimalType.ID,
			t.AnimalSpecies.ID,
			t.Microchip.ID,
			t.Organization.ID,
			t.OrganizationContact.ID,
			t.Address.ID,
			t.Country.ID,
			t.OrganizationWorkHour.ID,
		)
	q = getLimitOffset(q, filters.Pagination)

	var dest []struct {
		dbtype.AnimalWithJoinData
		TotalCount int64 `db:"total_count"`
	}
	if err := q.QueryContext(ctx, po.db, &dest); err != nil {
		return dbtype.PagedResult[dbtype.AnimalWithJoinData]{}, err
	}
	result := dbtype.PagedResult[dbtype.AnimalWithJoinData]{
		Data: make([]dbtype.AnimalWithJoinData, len(dest)),
	}
	for i, row := range dest {
		result.Data[i] = row.AnimalWithJoinData
	}
	if len(dest) > 0 {
		result.TotalCount = dest[0].TotalCount
	}
	return result, nil
}

func (po *PgAnimalPersistor) GetAnimalByID(ctx context.Context, animalID int64) (dbtype.AnimalWithJoinData, error) {
	q := p.SELECT(
		t.Animal.AllColumns,
		t.AnimalType.AllColumns,
		t.AnimalSpecies.AllColumns,
		t.Microchip.AllColumns,
		t.Organization.AllColumns,
		t.OrganizationContact.AllColumns,
		t.Address.AllColumns,
		t.Country.AllColumns,
		t.OrganizationWorkHour.AllColumns,
		p.SELECT_JSON_ARR(t.AnimalBreed.AllColumns, t.Breed.AllColumns).
			FROM(
				t.AnimalBreed.
					LEFT_JOIN(t.Breed, t.AnimalBreed.BreedID.EQ(t.Breed.ID)),
			).
			WHERE(t.AnimalBreed.AnimalID.EQ(t.Animal.ID)).
			AS("breeds"),
		p.SELECT_JSON_ARR(t.Tag.AllColumns).
			FROM(t.Tag).
			WHERE(t.Tag.AnimalID.EQ(t.Animal.ID)).
			AS("tags"),
		p.SELECT_JSON_ARR(t.AnimalPhoto.AllColumns).
			FROM(t.AnimalPhoto).
			WHERE(t.AnimalPhoto.AnimalID.EQ(t.Animal.ID)).
			AS("photos"),
		p.SELECT_JSON_ARR(t.AnimalVideo.AllColumns).
			FROM(t.AnimalVideo).
			WHERE(t.AnimalVideo.AnimalID.EQ(t.Animal.ID)).
			AS("videos"),
	).
		FROM(
			t.Animal.
				LEFT_JOIN(t.AnimalType, t.Animal.ID.EQ(t.AnimalType.ID)).
				LEFT_JOIN(t.AnimalSpecies, t.Animal.ID.EQ(t.AnimalSpecies.ID)).
				LEFT_JOIN(t.Microchip, t.Animal.ID.EQ(t.Microchip.AnimalID)).
				LEFT_JOIN(t.Organization, t.Animal.OrganizationID.EQ(t.Organization.ID)).
				LEFT_JOIN(t.OrganizationContact, t.Organization.ID.EQ(t.OrganizationContact.OrganizationID)).
				LEFT_JOIN(t.Address, t.Address.ID.EQ(t.OrganizationContact.AddressID)).
				LEFT_JOIN(t.Country, t.Country.ID.EQ(t.Address.CountryID)).
				LEFT_JOIN(t.OrganizationWorkHour, t.Organization.ID.EQ(t.OrganizationWorkHour.OrganizationID)),
		).
		WHERE(t.Animal.ID.EQ(p.Int64(animalID))).
		GROUP_BY(
			t.Animal.ID,
			t.AnimalType.ID,
			t.AnimalSpecies.ID,
			t.Microchip.ID,
			t.Organization.ID,
			t.OrganizationContact.ID,
			t.Address.ID,
			t.Country.ID,
			t.OrganizationWorkHour.ID,
		)
	var dest dbtype.AnimalWithJoinData
	if err := q.QueryContext(ctx, po.db, &dest); err != nil {
		if errors.Is(err, qrm.ErrNoRows) {
			return dest, ErrAnimalNotFound
		}
		return dest, err
	}
	return dest, nil
}

func (po *PgAnimalPersistor) CreateAnimal(ctx context.Context, in dbtype.AnimalCreateSetter) (model.Animal, error) {
	var insertedAnimal model.Animal

	txErr := po.WithTx(ctx, func(tx *sql.Tx) error {
		animalCols, org := in.Animal.ToModel()
		q1 := t.Animal.INSERT(animalCols).
			MODEL(org).
			RETURNING(t.Animal.AllColumns)
		if err := q1.QueryContext(ctx, tx, &insertedAnimal); err != nil {
			return fmt.Errorf("failed to insert an animal: %w", err)
		}

		if in.Microchip.IsSpecified() && !in.Microchip.IsNull() {
			var insertedMicrochip model.Microchip
			microchipInput := in.Microchip.MustGet()
			microchipInput.AnimalID = nullable.NewNullableWithValue(insertedAnimal.ID)
			microchipCols, microchip := microchipInput.ToModel()
			q1 := t.Microchip.INSERT(microchipCols).
				MODEL(microchip).
				RETURNING(t.Microchip.AllColumns)
			if err := q1.QueryContext(ctx, tx, &insertedMicrochip); err != nil {
				return fmt.Errorf("failed to insert animal microchip: %w", err)
			}
		}

		if in.Breeds.IsSpecified() && !in.Breeds.IsNull() {
			animalBreedsInput := in.Breeds.MustGet()
			if len(animalBreedsInput) > 0 {
				var insertedAnimalBreeds []model.AnimalBreed
				breeds := make([]model.AnimalBreed, len(animalBreedsInput))
				var animalBreedCols p.ColumnList
				for i, abreed := range animalBreedsInput {
					abreed.AnimalID = nullable.NewNullableWithValue(insertedAnimal.ID)
					cols, m := abreed.ToModel()
					if i == 0 {
						animalBreedCols = cols
					}
					breeds[i] = m
				}
				q3 := t.AnimalBreed.INSERT(animalBreedCols).
					MODELS(breeds).
					RETURNING(t.AnimalBreed.AllColumns)
				if err := q3.QueryContext(ctx, tx, &insertedAnimalBreeds); err != nil {
					return fmt.Errorf("failed to insert animal breeds: %w", err)
				}
			}
		}

		if in.Tags.IsSpecified() && !in.Tags.IsNull() {
			tagsInput := in.Tags.MustGet()
			if len(tagsInput) > 0 {
				var insertedTags []model.Tag
				photos := make([]model.Tag, len(tagsInput))
				var tagCols p.ColumnList
				for i, tag := range tagsInput {
					tag.AnimalID = nullable.NewNullableWithValue(insertedAnimal.ID)
					cols, m := tag.ToModel()
					if i == 0 {
						tagCols = cols
					}
					photos[i] = m
				}
				q3 := t.Tag.INSERT(tagCols).
					MODELS(photos).
					RETURNING(t.Tag.AllColumns)
				if err := q3.QueryContext(ctx, tx, &insertedTags); err != nil {
					return fmt.Errorf("failed to insert animal tags: %w", err)
				}
			}
		}

		if in.Photos.IsSpecified() && !in.Photos.IsNull() {
			photosInput := in.Photos.MustGet()
			if len(photosInput) > 0 {
				var insertedPhotos []model.AnimalPhoto
				photos := make([]model.AnimalPhoto, len(photosInput))
				var photoCols p.ColumnList
				for i, photo := range photosInput {
					photo.AnimalID = nullable.NewNullableWithValue(insertedAnimal.ID)
					cols, m := photo.ToModel()
					if i == 0 {
						photoCols = cols
					}
					photos[i] = m
				}
				q4 := t.AnimalPhoto.INSERT(photoCols).
					MODELS(photos).
					RETURNING(t.AnimalPhoto.AllColumns)
				if err := q4.QueryContext(ctx, tx, &insertedPhotos); err != nil {
					return fmt.Errorf("failed to insert animal photos: %w", err)
				}
			}
		}

		if in.Videos.IsSpecified() && !in.Videos.IsNull() {
			videosInput := in.Videos.MustGet()
			if len(videosInput) > 0 {
				var insertedVideos []model.AnimalVideo
				videos := make([]model.AnimalVideo, len(videosInput))
				var videoCols p.ColumnList
				for i, video := range videosInput {
					video.AnimalID = nullable.NewNullableWithValue(insertedAnimal.ID)
					cols, m := video.ToModel()
					if i == 0 {
						videoCols = cols
					}
					videos[i] = m
				}
				q5 := t.AnimalVideo.INSERT(videoCols).
					MODELS(videos).
					RETURNING(t.AnimalVideo.AllColumns)
				if err := q5.QueryContext(ctx, tx, &insertedVideos); err != nil {
					return fmt.Errorf("failed to insert animal videos: %w", err)
				}
			}
		}

		return nil
	})
	if txErr != nil {
		var pgErr *pgconn.PgError
		if errors.As(txErr, &pgErr) {
			return model.Animal{}, convertAnimalPgError(pgErr)
		}
		return model.Animal{}, txErr
	}
	return insertedAnimal, nil
}

func (po *PgAnimalPersistor) UpdateAnimal(ctx context.Context, animalID int64, in dbtype.AnimalSetter) (model.Animal, error) {
	cols, m := in.ToModel(true)
	q := t.Animal.UPDATE(cols).
		MODEL(m).
		WHERE(t.Animal.ID.EQ(p.Int64(animalID))).
		RETURNING(t.Animal.AllColumns)

	var dest model.Animal
	if err := q.QueryContext(ctx, po.db, &dest); err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) {
			return dest, convertAnimalPgError(pgErr)
		}
		return dest, fmt.Errorf("failed to update an animal: %w", err)
	}
	return dest, nil
}

func (po *PgAnimalPersistor) DeleteAnimalByID(ctx context.Context, animalID int64) (int64, error) {
	q := t.Animal.DELETE().WHERE(t.Animal.ID.EQ(p.Int64(animalID)))
	if _, err := q.ExecContext(ctx, po.db); err != nil {
		return 0, fmt.Errorf("failed to delete an animal: %w", err)
	}
	return animalID, nil
}
