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

func (pst *PgAnimalPersistor) ListAnimalTypes(ctx context.Context, filters dbtype.ListAnimalTypesFilters) (dbtype.PagedResult[models.AnimalType], error) {
	q := psql.Select(
		sm.Columns(models.AnimalTypes.Columns),
		sm.From(models.AnimalTypes.Name()),
		sm.GroupBy(models.AnimalTypes.Columns.ID),
	)
	addOrderBy(&q, filters.Sort, models.AnimalTypes.Columns.Names())
	addPagination(&q, filters.Page, filters.PageSize)

	if hasAnyLogicFilters(&filters.ListAnimalTypesParams) {
		if filters.ID != nil {
			ids := make([]any, len(*filters.ID))
			for i, id := range *filters.ID {
				ids[i] = id
			}

			q.Apply(sm.Where(models.AnimalTypes.Columns.ID.In(psql.Arg(ids...))))
		}

		if filters.Name != nil {
			q.Apply(sm.Where(models.AnimalTypes.Columns.Name.ILike(psql.Arg("%" + *filters.Name + "%"))))
		}
	}

	type ListAnimalTypesRow struct {
		models.AnimalType
		TotalCount int64
	}

	animalTypes, err := bob.All(ctx, pst.exec, q, scan.StructMapper[ListAnimalTypesRow]())
	if err != nil {
		return dbtype.PagedResult[models.AnimalType]{}, fmt.Errorf("query animalTypes")
	}

	result := dbtype.PagedResult[models.AnimalType]{
		Data: make([]models.AnimalType, len(animalTypes)),
	}
	for i, row := range animalTypes {
		result.Data[i] = row.AnimalType
	}

	if len(animalTypes) > 0 {
		result.TotalCount = animalTypes[0].TotalCount
	}

	return result, nil
}

func (pst *PgAnimalPersistor) GetAnimalTypeByID(ctx context.Context, animalTypeID int64) (models.AnimalType, error) {
	q := psql.Select(
		sm.Columns(models.AnimalTypes.Columns),
		sm.From(models.AnimalTypes.Name()),
		sm.Where(models.AnimalTypes.Columns.ID.EQ(psql.Arg(animalTypeID))),
	)

	animalType, err := bob.One(ctx, pst.exec, q, scan.StructMapper[models.AnimalType]())
	if err != nil {
		return models.AnimalType{}, fmt.Errorf("query animal specie")
	}

	return animalType, nil
}

func (pst *PgAnimalPersistor) CreateAnimalType(ctx context.Context, in models.AnimalTypeSetter) (models.AnimalType, error) {
	q := models.AnimalTypes.Insert(&in, im.Returning(models.AnimalTypes.Columns))

	animalType, err := bob.One(ctx, pst.exec, q, scan.StructMapper[models.AnimalType]())
	if err != nil {
		return models.AnimalType{}, fmt.Errorf("insert animal specie")
	}

	return animalType, nil
}

func (pst *PgAnimalPersistor) UpdateAnimalType(ctx context.Context, animalTypeID int64, in models.AnimalTypeSetter) (models.AnimalType, error) {
	q := models.AnimalTypes.Update(
		in.UpdateMod(),
		um.Where(models.AnimalTypes.Columns.ID.EQ(psql.Arg(animalTypeID))),
		um.Returning(models.AnimalTypes.Columns),
	)

	animalTypes, err := bob.One(ctx, pst.exec, q, scan.StructMapper[models.AnimalType]())
	if err != nil {
		return models.AnimalType{}, fmt.Errorf("update animal specie")
	}

	return animalTypes, nil
}

func (pst *PgAnimalPersistor) DeleteAnimalTypeByID(ctx context.Context, animalTypeID int64) (int64, error) {
	q := models.AnimalTypes.Delete(dm.Where(models.AnimalTypes.Columns.ID.EQ(psql.Arg(animalTypeID))))
	if _, err := bob.Exec(ctx, pst.exec, q); err != nil {
		return 0, fmt.Errorf("delete animal specie: %w", err)
	}

	return animalTypeID, nil
}

func (pst *PgAnimalPersistor) DeleteAnimalTypes(ctx context.Context, ids []int64) error {
	specieIDs := make([]any, len(ids))
	for i, id := range ids {
		specieIDs[i] = id
	}

	q := models.AnimalTypes.Delete(dm.Where(models.AnimalTypes.Columns.ID.In(psql.Arg(specieIDs...))))
	if _, err := bob.Exec(ctx, pst.exec, q); err != nil {
		return fmt.Errorf("delete animal species: %w", err)
	}

	return nil
}
