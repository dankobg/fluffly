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

func (pst *PgAnimalPersistor) ListAnimalTags(ctx context.Context, animalID int64, filters dbtype.ListAnimalTagsFilters) (dbtype.PagedResult[models.AnimalTag], error) {
	q := psql.Select(
		sm.Columns(models.AnimalTags.Columns),
		sm.From(models.AnimalTags.Name()),
		sm.GroupBy(models.AnimalTags.Columns.ID),
	)
	addOrderBy(&q, filters.Sort, models.AnimalTags.Columns.Names())
	addPagination(&q, filters.Page, filters.PageSize)

	if hasAnyLogicFilters(&filters.ListAnimalTagsParams) {
		if filters.ID != nil {
			ids := make([]any, len(*filters.ID))
			for i, id := range *filters.ID {
				ids[i] = id
			}

			q.Apply(sm.Where(models.AnimalTags.Columns.ID.In(psql.Arg(ids...))))
		}
	}

	type ListAnimalTagsRow struct {
		models.AnimalTag
		TotalCount int64
	}

	animalTags, err := bob.All(ctx, pst.exec, q, scan.StructMapper[ListAnimalTagsRow]())
	if err != nil {
		return dbtype.PagedResult[models.AnimalTag]{}, fmt.Errorf("query animal tags")
	}

	result := dbtype.PagedResult[models.AnimalTag]{
		Data: make([]models.AnimalTag, len(animalTags)),
	}
	for i, row := range animalTags {
		result.Data[i] = row.AnimalTag
	}

	if len(animalTags) > 0 {
		result.TotalCount = animalTags[0].TotalCount
	}

	return result, nil
}

func (pst *PgAnimalPersistor) GetAnimalTag(ctx context.Context, animalID, tagID int64) (models.AnimalTag, error) {
	q := psql.Select(
		sm.Columns(models.AnimalTags.Columns),
		sm.From(models.AnimalTags.Name()),
		sm.Where(models.AnimalTags.Columns.ID.EQ(psql.Arg(tagID))),
	)

	animalTag, err := bob.One(ctx, pst.exec, q, scan.StructMapper[models.AnimalTag]())
	if err != nil {
		return models.AnimalTag{}, fmt.Errorf("query animal tag")
	}

	return animalTag, nil
}

func (pst *PgAnimalPersistor) CreateAnimalTags(ctx context.Context, animalID int64, in []models.AnimalTagSetter) ([]models.AnimalTag, error) {
	setters := make([]*models.AnimalTagSetter, len(in))
	for i, x := range in {
		setters[i] = &x
	}

	q := models.AnimalTags.Insert(bob.ToMods(setters...), im.Returning(models.AnimalTags.Columns))

	animalTags, err := bob.All(ctx, pst.exec, q, scan.StructMapper[models.AnimalTag]())
	if err != nil {
		return nil, fmt.Errorf("insert animal tags")
	}

	return animalTags, nil
}

func (pst *PgAnimalPersistor) DeleteAnimalTags(ctx context.Context, animalID int64, ids []int64) error {
	tagIDs := make([]any, len(ids))
	for i, id := range ids {
		tagIDs[i] = id
	}

	q := models.AnimalTags.Delete(
		dm.Where(models.AnimalTags.Columns.AnimalID.EQ(psql.Arg(animalID)).And(
			models.AnimalTags.Columns.ID.In(psql.Arg(tagIDs...)),
		)),
	)
	if _, err := bob.Exec(ctx, pst.exec, q); err != nil {
		return fmt.Errorf("delete animal tags: %w", err)
	}

	return nil
}

func (pst *PgAnimalPersistor) UpdateAnimalTag(ctx context.Context, animalID, tagID int64, in models.AnimalTagSetter) (models.AnimalTag, error) {
	q := models.AnimalTags.Update(
		in.UpdateMod(),
		um.Where(models.AnimalTags.Columns.AnimalID.EQ(psql.Arg(animalID)).And(
			models.AnimalTags.Columns.ID.EQ(psql.Arg(tagID)),
		)),
		um.Returning(models.AnimalTags.Columns),
	)

	animalTag, err := bob.One(ctx, pst.exec, q, scan.StructMapper[models.AnimalTag]())
	if err != nil {
		return models.AnimalTag{}, fmt.Errorf("update animal tag")
	}

	return animalTag, nil
}

func (pst *PgAnimalPersistor) DeleteAnimalTag(ctx context.Context, animalID, tagID int64) (int64, error) {
	q := models.AnimalTags.Delete(
		dm.Where(models.AnimalTags.Columns.AnimalID.EQ(psql.Arg(animalID)).And(
			models.AnimalTags.Columns.ID.EQ(psql.Arg(tagID)),
		)),
	)
	if _, err := bob.Exec(ctx, pst.exec, q); err != nil {
		return 0, fmt.Errorf("delete animal tag: %w", err)
	}

	return tagID, nil
}
