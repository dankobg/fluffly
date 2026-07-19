package postgres

import (
	"context"
	"fmt"

	"github.com/aarondl/opt/omitnull"
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

func (pst *PgAnimalPersistor) ListAnimalPhotos(ctx context.Context, animalID int64, filters dbtype.ListAnimalPhotosFilters) (dbtype.PagedResult[models.AnimalPhoto], error) {
	q := psql.Select(
		sm.Columns(models.AnimalPhotos.Columns),
		sm.From(models.AnimalPhotos.Name()),
		sm.Where(models.AnimalPhotos.Columns.AnimalID.EQ(psql.Arg(animalID))),
		sm.GroupBy(models.AnimalPhotos.Columns.ID),
	)
	addOrderBy(&q, filters.Sort, models.AnimalPhotos.Columns.Names())
	addPagination(&q, filters.Page, filters.PageSize)

	if hasAnyLogicFilters(&filters.ListAnimalPhotosParams) {
		if filters.ID != nil {
			ids := make([]any, len(*filters.ID))
			for i, id := range *filters.ID {
				ids[i] = id
			}

			q.Apply(sm.Where(models.AnimalPhotos.Columns.ID.In(psql.Arg(ids...))))
		}
	}

	type ListAnimalPhotosRow struct {
		models.AnimalPhoto
		TotalCount int64
	}

	animalPhotos, err := bob.All(ctx, pst.exec, q, scan.StructMapper[ListAnimalPhotosRow]())
	if err != nil {
		return dbtype.PagedResult[models.AnimalPhoto]{}, fmt.Errorf("query animal photos")
	}

	result := dbtype.PagedResult[models.AnimalPhoto]{
		Data: make([]models.AnimalPhoto, len(animalPhotos)),
	}
	for i, row := range animalPhotos {
		result.Data[i] = row.AnimalPhoto
	}

	if len(animalPhotos) > 0 {
		result.TotalCount = animalPhotos[0].TotalCount
	}

	return result, nil
}

func (pst *PgAnimalPersistor) GetAnimalPhoto(ctx context.Context, animalID, photoID int64) (models.AnimalPhoto, error) {
	q := psql.Select(
		sm.Columns(models.AnimalPhotos.Columns),
		sm.From(models.AnimalPhotos.Name()),
		sm.Where(models.AnimalPhotos.Columns.AnimalID.EQ(psql.Arg(animalID)).
			And(models.AnimalPhotos.Columns.ID.EQ(psql.Arg(photoID))),
		),
	)

	animalPhoto, err := bob.One(ctx, pst.exec, q, scan.StructMapper[models.AnimalPhoto]())
	if err != nil {
		return models.AnimalPhoto{}, fmt.Errorf("query animal photo")
	}

	return animalPhoto, nil
}

func (pst *PgAnimalPersistor) CreateAnimalPhotos(ctx context.Context, animalID int64, in []models.AnimalPhotoSetter) ([]models.AnimalPhoto, error) {
	setters := make([]*models.AnimalPhotoSetter, len(in))
	for i, x := range in {
		x.AnimalID = omitnull.From(animalID)
		setters[i] = &x
	}

	q := models.AnimalPhotos.Insert(bob.ToMods(setters...), im.Returning(models.AnimalPhotos.Columns))

	animalPhotos, err := bob.All(ctx, pst.exec, q, scan.StructMapper[models.AnimalPhoto]())
	if err != nil {
		return nil, fmt.Errorf("insert animal photos")
	}

	return animalPhotos, nil
}

func (pst *PgAnimalPersistor) DeleteAnimalPhotos(ctx context.Context, animalID int64, ids []int64) error {
	photoIDs := make([]any, len(ids))
	for i, id := range ids {
		photoIDs[i] = id
	}

	q := models.AnimalPhotos.Delete(
		dm.Where(models.AnimalPhotos.Columns.AnimalID.EQ(psql.Arg(animalID)).And(
			models.AnimalPhotos.Columns.ID.In(psql.Arg(photoIDs...)),
		)),
	)
	if _, err := bob.Exec(ctx, pst.exec, q); err != nil {
		return fmt.Errorf("delete animal photos: %w", err)
	}

	return nil
}

func (pst *PgAnimalPersistor) UpdateAnimalPhoto(ctx context.Context, animalID, photoID int64, in models.AnimalPhotoSetter) (models.AnimalPhoto, error) {
	q := models.AnimalPhotos.Update(
		in.UpdateMod(),
		um.Where(models.AnimalPhotos.Columns.AnimalID.EQ(psql.Arg(animalID)).And(
			models.AnimalPhotos.Columns.ID.EQ(psql.Arg(photoID)),
		)),
		um.Returning(models.AnimalPhotos.Columns),
	)

	animalPhoto, err := bob.One(ctx, pst.exec, q, scan.StructMapper[models.AnimalPhoto]())
	if err != nil {
		return models.AnimalPhoto{}, fmt.Errorf("update animal photo")
	}

	return animalPhoto, nil
}

func (pst *PgAnimalPersistor) DeleteAnimalPhoto(ctx context.Context, animalID, photoID int64) (int64, error) {
	q := models.AnimalPhotos.Delete(
		dm.Where(models.AnimalPhotos.Columns.AnimalID.EQ(psql.Arg(animalID)).And(
			models.AnimalPhotos.Columns.ID.EQ(psql.Arg(photoID)),
		)),
	)
	if _, err := bob.Exec(ctx, pst.exec, q); err != nil {
		return 0, fmt.Errorf("delete animal photo: %w", err)
	}

	return photoID, nil
}
