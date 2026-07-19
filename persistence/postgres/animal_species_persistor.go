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

func (pst *PgAnimalPersistor) ListAnimalSpecies(ctx context.Context, filters dbtype.ListAnimalSpeciesFilters) (dbtype.PagedResult[models.AnimalSpecie], error) {
	q := psql.Select(
		sm.Columns(models.AnimalSpecies.Columns),
		sm.From(models.AnimalSpecies.Name()),
		sm.GroupBy(models.AnimalSpecies.Columns.ID),
	)
	addOrderBy(&q, filters.Sort, models.AnimalSpecies.Columns.Names())
	addPagination(&q, filters.Page, filters.PageSize)

	if hasAnyLogicFilters(&filters.ListAnimalSpeciesParams) {
		if filters.ID != nil {
			ids := make([]any, len(*filters.ID))
			for i, id := range *filters.ID {
				ids[i] = id
			}

			q.Apply(sm.Where(models.AnimalSpecies.Columns.ID.In(psql.Arg(ids...))))
		}

		if filters.AnimalTypeID != nil {
			typeIDs := make([]any, len(*filters.AnimalTypeID))
			for i, id := range *filters.AnimalTypeID {
				typeIDs[i] = id
			}

			q.Apply(sm.Where(models.AnimalSpecies.Columns.AnimalTypeID.In(psql.Arg(typeIDs...))))
		}

		if filters.Name != nil {
			q.Apply(sm.Where(models.AnimalSpecies.Columns.Name.ILike(psql.Arg("%" + *filters.Name + "%"))))
		}
	}

	type ListAnimalSpeciesRow struct {
		models.AnimalSpecie
		TotalCount int64
	}

	animalSpecies, err := bob.All(ctx, pst.exec, q, scan.StructMapper[ListAnimalSpeciesRow]())
	if err != nil {
		return dbtype.PagedResult[models.AnimalSpecie]{}, fmt.Errorf("query animalSpecies")
	}

	result := dbtype.PagedResult[models.AnimalSpecie]{
		Data: make([]models.AnimalSpecie, len(animalSpecies)),
	}
	for i, row := range animalSpecies {
		result.Data[i] = row.AnimalSpecie
	}

	if len(animalSpecies) > 0 {
		result.TotalCount = animalSpecies[0].TotalCount
	}

	return result, nil
}

func (pst *PgAnimalPersistor) GetAnimalSpecieByID(ctx context.Context, animalSpecieID int64) (models.AnimalSpecie, error) {
	q := psql.Select(
		sm.Columns(models.AnimalSpecies.Columns),
		sm.From(models.AnimalSpecies.Name()),
		sm.Where(models.AnimalSpecies.Columns.ID.EQ(psql.Arg(animalSpecieID))),
	)

	animalSpecie, err := bob.One(ctx, pst.exec, q, scan.StructMapper[models.AnimalSpecie]())
	if err != nil {
		return models.AnimalSpecie{}, fmt.Errorf("query animal specie")
	}

	return animalSpecie, nil
}

func (pst *PgAnimalPersistor) CreateAnimalSpecie(ctx context.Context, in models.AnimalSpecieSetter) (models.AnimalSpecie, error) {
	q := models.AnimalSpecies.Insert(&in, im.Returning(models.AnimalSpecies.Columns))

	animalSpecie, err := bob.One(ctx, pst.exec, q, scan.StructMapper[models.AnimalSpecie]())
	if err != nil {
		return models.AnimalSpecie{}, fmt.Errorf("insert animal specie")
	}

	return animalSpecie, nil
}

func (pst *PgAnimalPersistor) UpdateAnimalSpecie(ctx context.Context, animalSpecieID int64, in models.AnimalSpecieSetter) (models.AnimalSpecie, error) {
	q := models.AnimalSpecies.Update(
		in.UpdateMod(),
		um.Where(models.AnimalSpecies.Columns.ID.EQ(psql.Arg(animalSpecieID))),
		um.Returning(models.AnimalSpecies.Columns),
	)

	animalSpecies, err := bob.One(ctx, pst.exec, q, scan.StructMapper[models.AnimalSpecie]())
	if err != nil {
		return models.AnimalSpecie{}, fmt.Errorf("update animal specie")
	}

	return animalSpecies, nil
}

func (pst *PgAnimalPersistor) DeleteAnimalSpecieByID(ctx context.Context, animalSpecieID int64) (int64, error) {
	q := models.AnimalSpecies.Delete(dm.Where(models.AnimalSpecies.Columns.ID.EQ(psql.Arg(animalSpecieID))))
	if _, err := bob.Exec(ctx, pst.exec, q); err != nil {
		return 0, fmt.Errorf("delete animal specie: %w", err)
	}

	return animalSpecieID, nil
}

func (pst *PgAnimalPersistor) DeleteAnimalSpecies(ctx context.Context, ids []int64) error {
	specieIDs := make([]any, len(ids))
	for i, id := range ids {
		specieIDs[i] = id
	}

	q := models.AnimalSpecies.Delete(dm.Where(models.AnimalSpecies.Columns.ID.In(psql.Arg(specieIDs...))))
	if _, err := bob.Exec(ctx, pst.exec, q); err != nil {
		return fmt.Errorf("delete animal species: %w", err)
	}

	return nil
}
