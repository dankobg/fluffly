package postgres

import (
	"context"
	"errors"
	"fmt"

	"github.com/dankobg/fluffly/db/gen/models"
	"github.com/dankobg/fluffly/persistence/dbtype"
	"github.com/stephenafamo/bob"
	"github.com/stephenafamo/bob/dialect/psql"
	"github.com/stephenafamo/bob/dialect/psql/dm"
	"github.com/stephenafamo/bob/dialect/psql/im"
	"github.com/stephenafamo/bob/dialect/psql/sm"
	"github.com/stephenafamo/bob/dialect/psql/um"
	"github.com/stephenafamo/scan"
)

type (
	ErrBreedIntegrityViolation  struct{ errIntegrityViolation }
	ErrBreedUniqueViolation     struct{ errUniqueViolation }
	ErrBreedForeignKeyViolation struct{ errForeignKeyViolation }
	ErrBreedCheckViolation      struct{ errCheckViolation }
)

var ErrBreedNotFound = errors.New("breed not found")

func (pst *PgAnimalPersistor) ListBreeds(ctx context.Context, filters dbtype.ListBreedsFilters) (dbtype.PagedResult[models.Breed], error) {
	q := psql.Select(
		sm.Columns(models.Breeds.Columns),
		sm.From(models.Breeds.Name()),
		sm.GroupBy(models.Breeds.Columns.ID),
	)
	addOrderBy(&q, filters.Sort, models.Breeds.Columns.Names())
	addPagination(&q, filters.Page, filters.PageSize)

	if hasAnyLogicFilters(&filters.ListBreedsParams) {
		if filters.ID != nil {
			ids := make([]any, len(*filters.ID))
			for i, id := range *filters.ID {
				ids[i] = id
			}

			q.Apply(sm.Where(models.Breeds.Columns.ID.In(psql.Arg(ids...))))
		}

		if filters.Name != nil {
			q.Apply(sm.Where(models.Breeds.Columns.Name.ILike(psql.Arg("%" + *filters.Name + "%"))))
		}
	}

	type ListBreedsRow struct {
		models.Breed
		TotalCount int64
	}

	breeds, err := bob.All(ctx, pst.exec, q, scan.StructMapper[ListBreedsRow]())
	if err != nil {
		return dbtype.PagedResult[models.Breed]{}, fmt.Errorf("query breeds")
	}

	result := dbtype.PagedResult[models.Breed]{
		Data: make([]models.Breed, len(breeds)),
	}
	for i, row := range breeds {
		result.Data[i] = row.Breed
	}

	if len(breeds) > 0 {
		result.TotalCount = breeds[0].TotalCount
	}

	return result, nil
}

func (pst *PgAnimalPersistor) GetBreedByID(ctx context.Context, breedID int64) (models.Breed, error) {
	q := psql.Select(
		sm.Columns(models.Breeds.Columns),
		sm.From(models.Breeds.Name()),
		sm.Where(models.Breeds.Columns.ID.EQ(psql.Arg(breedID))),
	)

	breed, err := bob.One(ctx, pst.exec, q, scan.StructMapper[models.Breed]())
	if err != nil {
		return models.Breed{}, fmt.Errorf("query breed")
	}

	return breed, nil
}

func (pst *PgAnimalPersistor) CreateBreed(ctx context.Context, in models.BreedSetter) (models.Breed, error) {
	q := models.Breeds.Insert(&in, im.Returning(models.Breeds.Columns))

	breed, err := bob.One(ctx, pst.exec, q, scan.StructMapper[models.Breed]())
	if err != nil {
		return models.Breed{}, fmt.Errorf("insert breed")
	}

	return breed, nil
}

func (pst *PgAnimalPersistor) UpdateBreed(ctx context.Context, breedID int64, in models.BreedSetter) (models.Breed, error) {
	q := models.Breeds.Update(
		in.UpdateMod(),
		um.Where(models.Breeds.Columns.ID.EQ(psql.Arg(breedID))),
		um.Returning(models.Breeds.Columns),
	)

	breed, err := bob.One(ctx, pst.exec, q, scan.StructMapper[models.Breed]())
	if err != nil {
		return models.Breed{}, fmt.Errorf("update breed")
	}

	return breed, nil
}

func (pst *PgAnimalPersistor) DeleteBreedByID(ctx context.Context, breedID int64) (int64, error) {
	q := models.Breeds.Delete(dm.Where(models.Breeds.Columns.ID.EQ(psql.Arg(breedID))))
	if _, err := bob.Exec(ctx, pst.exec, q); err != nil {
		return 0, fmt.Errorf("delete breed: %w", err)
	}

	return breedID, nil
}

func (pst *PgAnimalPersistor) DeleteBreeds(ctx context.Context, ids []int64) error {
	breedIDs := make([]any, len(ids))
	for i, id := range ids {
		breedIDs[i] = id
	}

	q := models.Breeds.Delete(dm.Where(models.Breeds.Columns.ID.In(psql.Arg(breedIDs...))))
	if _, err := bob.Exec(ctx, pst.exec, q); err != nil {
		return fmt.Errorf("delete breeds: %w", err)
	}

	return nil
}
