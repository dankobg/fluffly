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

func (pst *PgAnimalPersistor) ListAnimalVideos(ctx context.Context, animalID int64, filters dbtype.ListAnimalVideosFilters) (dbtype.PagedResult[models.AnimalVideo], error) {
	q := psql.Select(
		sm.Columns(models.AnimalVideos.Columns),
		sm.From(models.AnimalVideos.Name()),
		sm.Where(models.AnimalPhotos.Columns.ID.EQ(psql.Arg(animalID))),
		sm.GroupBy(models.AnimalVideos.Columns.ID),
	)
	addOrderBy(&q, filters.Sort, models.AnimalVideos.Columns.Names())
	addPagination(&q, filters.Page, filters.PageSize)

	if hasAnyLogicFilters(&filters.ListAnimalVideosParams) {
		if filters.ID != nil {
			ids := make([]any, len(*filters.ID))
			for i, id := range *filters.ID {
				ids[i] = id
			}

			q.Apply(sm.Where(models.AnimalVideos.Columns.ID.In(psql.Arg(ids...))))
		}
	}

	type ListAnimalVideosRow struct {
		models.AnimalVideo
		TotalCount int64
	}

	animalVideos, err := bob.All(ctx, pst.exec, q, scan.StructMapper[ListAnimalVideosRow]())
	if err != nil {
		return dbtype.PagedResult[models.AnimalVideo]{}, fmt.Errorf("query animal videos")
	}

	result := dbtype.PagedResult[models.AnimalVideo]{
		Data: make([]models.AnimalVideo, len(animalVideos)),
	}
	for i, row := range animalVideos {
		result.Data[i] = row.AnimalVideo
	}

	if len(animalVideos) > 0 {
		result.TotalCount = animalVideos[0].TotalCount
	}

	return result, nil
}

func (pst *PgAnimalPersistor) GetAnimalVideo(ctx context.Context, animalID, videoID int64) (models.AnimalVideo, error) {
	q := psql.Select(
		sm.Columns(models.AnimalVideos.Columns),
		sm.From(models.AnimalVideos.Name()),
		sm.Where(models.AnimalVideos.Columns.AnimalID.EQ(psql.Arg(animalID)).
			And(models.AnimalVideos.Columns.ID.EQ(psql.Arg(videoID))),
		),
	)

	animalVideo, err := bob.One(ctx, pst.exec, q, scan.StructMapper[models.AnimalVideo]())
	if err != nil {
		return models.AnimalVideo{}, fmt.Errorf("query animal video")
	}

	return animalVideo, nil
}

func (pst *PgAnimalPersistor) CreateAnimalVideos(ctx context.Context, animalID int64, in []models.AnimalVideoSetter) ([]models.AnimalVideo, error) {
	setters := make([]*models.AnimalVideoSetter, len(in))
	for i, x := range in {
		x.AnimalID = omitnull.From(animalID)
		setters[i] = &x
	}

	q := models.AnimalVideos.Insert(bob.ToMods(setters...), im.Returning(models.AnimalVideos.Columns))

	animalVideos, err := bob.All(ctx, pst.exec, q, scan.StructMapper[models.AnimalVideo]())
	if err != nil {
		return nil, fmt.Errorf("insert animal videos")
	}

	return animalVideos, nil
}

func (pst *PgAnimalPersistor) DeleteAnimalVideos(ctx context.Context, animalID int64, ids []int64) error {
	videoIDs := make([]any, len(ids))
	for i, id := range ids {
		videoIDs[i] = id
	}

	q := models.AnimalVideos.Delete(
		dm.Where(models.AnimalVideos.Columns.AnimalID.EQ(psql.Arg(animalID)).And(
			models.AnimalVideos.Columns.ID.In(psql.Arg(videoIDs...)),
		)),
	)
	if _, err := bob.Exec(ctx, pst.exec, q); err != nil {
		return fmt.Errorf("delete animal videos: %w", err)
	}

	return nil
}

func (pst *PgAnimalPersistor) UpdateAnimalVideo(ctx context.Context, animalID, videoID int64, in models.AnimalVideoSetter) (models.AnimalVideo, error) {
	q := models.AnimalVideos.Update(
		in.UpdateMod(),
		um.Where(models.AnimalVideos.Columns.AnimalID.EQ(psql.Arg(animalID)).And(
			models.AnimalVideos.Columns.ID.EQ(psql.Arg(videoID)),
		)),
		um.Returning(models.AnimalVideos.Columns),
	)

	animalVideo, err := bob.One(ctx, pst.exec, q, scan.StructMapper[models.AnimalVideo]())
	if err != nil {
		return models.AnimalVideo{}, fmt.Errorf("update animal video")
	}

	return animalVideo, nil
}

func (pst *PgAnimalPersistor) DeleteAnimalVideo(ctx context.Context, animalID, videoID int64) (int64, error) {
	q := models.AnimalVideos.Delete(
		dm.Where(models.AnimalVideos.Columns.AnimalID.EQ(psql.Arg(animalID)).And(
			models.AnimalVideos.Columns.ID.EQ(psql.Arg(videoID)),
		)),
	)
	if _, err := bob.Exec(ctx, pst.exec, q); err != nil {
		return 0, fmt.Errorf("delete animal video: %w", err)
	}

	return videoID, nil
}
