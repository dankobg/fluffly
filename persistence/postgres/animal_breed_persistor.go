package postgres

import (
	"context"
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

func (pst *PgAnimalPersistor) ListAnimalBreeds(ctx context.Context, animalID int64, filters dbtype.ListAnimalBreedsFilters) (dbtype.PagedResult[dbtype.AnimalBreedWithJoinData], error) {
	q := psql.Select(
		sm.Columns(
			models.Breeds.Columns,
			models.AnimalBreeds.Columns.Primary,
		),
		sm.From(models.Breeds.Name()),
		sm.InnerJoin(models.AnimalBreeds.Name()).
			On(models.Breeds.Columns.ID.EQ(models.AnimalBreeds.Columns.BreedID)),
		sm.Where(models.AnimalBreeds.Columns.AnimalID.EQ(psql.Arg(animalID))),
	)
	addOrderBy(&q, filters.Sort, models.Countries.Columns.Names())
	addPagination(&q, filters.Page, filters.PageSize)

	if hasAnyLogicFilters(&filters.ListAnimalBreedsParams) {
		if filters.ID != nil {
			ids := make([]any, len(*filters.ID))
			for i, id := range *filters.ID {
				ids[i] = id
			}

			q.Apply(sm.Where(models.AnimalBreeds.Columns.BreedID.In(psql.Arg(ids...))))
		}
	}

	type ListAnimalTagsRow struct {
		dbtype.AnimalBreedWithJoinData
		TotalCount int64
	}

	animalTags, err := bob.All(ctx, pst.exec, q, scan.StructMapper[ListAnimalTagsRow]())
	if err != nil {
		return dbtype.PagedResult[dbtype.AnimalBreedWithJoinData]{}, fmt.Errorf("query animal tags")
	}

	result := dbtype.PagedResult[dbtype.AnimalBreedWithJoinData]{
		Data: make([]dbtype.AnimalBreedWithJoinData, len(animalTags)),
	}
	for i, row := range animalTags {
		result.Data[i] = row.AnimalBreedWithJoinData
	}

	if len(animalTags) > 0 {
		result.TotalCount = animalTags[0].TotalCount
	}

	return result, nil
}

func (pst *PgAnimalPersistor) GetAnimalBreed(ctx context.Context, animalID, breedID int64) (dbtype.AnimalBreedWithJoinData, error) {
	q := psql.Select(
		sm.Columns(
			models.Breeds.Columns,
			models.AnimalBreeds.Columns.Primary,
		),
		sm.From(models.Breeds.Name()),
		sm.InnerJoin(models.AnimalBreeds.Name()).
			On(models.Breeds.Columns.ID.EQ(models.AnimalBreeds.Columns.BreedID)),
		sm.Where(models.AnimalBreeds.Columns.AnimalID.EQ(psql.Arg(animalID)).
			And(models.AnimalBreeds.Columns.BreedID.EQ(psql.Arg(breedID))),
		),
	)

	animalBreed, err := bob.One(ctx, pst.exec, q, scan.StructMapper[dbtype.AnimalBreedWithJoinData]())
	if err != nil {
		return dbtype.AnimalBreedWithJoinData{}, fmt.Errorf("query animal breed")
	}

	return animalBreed, nil
}

func (pst *PgAnimalPersistor) CreateAnimalBreeds(ctx context.Context, animalID int64, in []models.AnimalBreedSetter) ([]models.AnimalBreed, error) {
	setters := make([]*models.AnimalBreedSetter, len(in))
	for i, x := range in {
		setters[i] = &x
	}

	q := models.AnimalBreeds.Insert(bob.ToMods(setters...), im.Returning(models.AnimalBreeds.Columns))

	animalBreeds, err := bob.All(ctx, pst.exec, q, scan.StructMapper[models.AnimalBreed]())
	if err != nil {
		return nil, fmt.Errorf("insert animal breeds")
	}

	return animalBreeds, nil
}

func (pst *PgAnimalPersistor) DeleteAnimalBreeds(ctx context.Context, animalID int64, ids []int64) error {
	breedIDs := make([]any, len(ids))
	for i, id := range ids {
		breedIDs[i] = id
	}

	q := models.AnimalBreeds.Delete(
		dm.Where(models.AnimalBreeds.Columns.AnimalID.EQ(psql.Arg(animalID)).And(
			models.AnimalBreeds.Columns.BreedID.In(psql.Arg(breedIDs...)),
		)),
	)
	if _, err := bob.Exec(ctx, pst.exec, q); err != nil {
		return fmt.Errorf("delete animal breeds: %w", err)
	}

	return nil
}

func (pst *PgAnimalPersistor) DeleteAnimalBreed(ctx context.Context, animalID, breedID int64) (int64, error) {
	q := models.AnimalBreeds.Delete(
		dm.Where(models.AnimalBreeds.Columns.AnimalID.EQ(psql.Arg(animalID)).And(
			models.AnimalBreeds.Columns.BreedID.EQ(psql.Arg(breedID)),
		)),
	)
	if _, err := bob.Exec(ctx, pst.exec, q); err != nil {
		return 0, fmt.Errorf("delete animal breed: %w", err)
	}

	return breedID, nil
}

func (pst *PgAnimalPersistor) UpdateAnimalBreed(ctx context.Context, animalID, breedID int64, in models.AnimalBreedSetter) (models.AnimalBreed, error) {
	q := models.AnimalBreeds.Update(
		in.UpdateMod(),
		um.Where(models.AnimalBreeds.Columns.AnimalID.EQ(psql.Arg(animalID)).And(
			models.AnimalBreeds.Columns.BreedID.EQ(psql.Arg(breedID)),
		)),
		um.Returning(models.AnimalBreeds.Columns),
	)

	animalBreed, err := bob.One(ctx, pst.exec, q, scan.StructMapper[models.AnimalBreed]())
	if err != nil {
		return models.AnimalBreed{}, fmt.Errorf("update animal breed")
	}

	return animalBreed, nil
}
