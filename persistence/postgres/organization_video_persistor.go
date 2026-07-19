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

func (pst *PgOrganizationPersistor) ListOrganizationVideos(ctx context.Context, organizationID int64, filters dbtype.ListOrganizationVideosFilters) (dbtype.PagedResult[models.OrganizationVideo], error) {
	q := psql.Select(
		sm.Columns(models.OrganizationVideos.Columns),
		sm.From(models.OrganizationVideos.Name()),
		sm.Where(models.OrganizationVideos.Columns.OrganizationID.EQ(psql.Arg(organizationID))),
		sm.GroupBy(models.OrganizationVideos.Columns.ID),
	)
	addOrderBy(&q, filters.Sort, models.OrganizationVideos.Columns.Names())
	addPagination(&q, filters.Page, filters.PageSize)

	if hasAnyLogicFilters(&filters.ListOrganizationVideosParams) {
		if filters.ID != nil {
			ids := make([]any, len(*filters.ID))
			for i, id := range *filters.ID {
				ids[i] = id
			}

			q.Apply(sm.Where(models.OrganizationVideos.Columns.ID.In(psql.Arg(ids...))))
		}
	}

	type ListOrganizationVideosRow struct {
		models.OrganizationVideo
		TotalCount int64
	}

	organizationVideos, err := bob.All(ctx, pst.exec, q, scan.StructMapper[ListOrganizationVideosRow]())
	if err != nil {
		return dbtype.PagedResult[models.OrganizationVideo]{}, fmt.Errorf("query organization videos")
	}

	result := dbtype.PagedResult[models.OrganizationVideo]{
		Data: make([]models.OrganizationVideo, len(organizationVideos)),
	}
	for i, row := range organizationVideos {
		result.Data[i] = row.OrganizationVideo
	}

	if len(organizationVideos) > 0 {
		result.TotalCount = organizationVideos[0].TotalCount
	}

	return result, nil
}

func (pst *PgOrganizationPersistor) GetOrganizationVideo(ctx context.Context, organizationID, videoID int64) (models.OrganizationVideo, error) {
	q := psql.Select(
		sm.Columns(models.OrganizationVideos.Columns),
		sm.From(models.OrganizationVideos.Name()),
		sm.Where(models.OrganizationVideos.Columns.OrganizationID.EQ(psql.Arg(organizationID))),
		sm.Where(models.OrganizationVideos.Columns.ID.EQ(psql.Arg(videoID))),
	)

	organizationVideo, err := bob.One(ctx, pst.exec, q, scan.StructMapper[models.OrganizationVideo]())
	if err != nil {
		return models.OrganizationVideo{}, fmt.Errorf("query organization video")
	}

	return organizationVideo, nil
}

func (pst *PgOrganizationPersistor) CreateOrganizationVideos(ctx context.Context, organizationID int64, in []models.OrganizationVideoSetter) ([]models.OrganizationVideo, error) {
	setters := make([]*models.OrganizationVideoSetter, len(in))
	for i, x := range in {
		x.OrganizationID = omitnull.From(organizationID)
		setters[i] = &x
	}

	q := models.OrganizationVideos.Insert(bob.ToMods(setters...), im.Returning(models.OrganizationVideos.Columns))

	organizationVideos, err := bob.All(ctx, pst.exec, q, scan.StructMapper[models.OrganizationVideo]())
	if err != nil {
		return nil, fmt.Errorf("insert animal videos")
	}

	return organizationVideos, nil
}

func (pst *PgOrganizationPersistor) DeleteOrganizationVideos(ctx context.Context, organizationID int64, ids []int64) error {
	videoIDs := make([]any, len(ids))
	for i, id := range ids {
		videoIDs[i] = id
	}

	q := models.OrganizationVideos.Delete(
		dm.Where(models.OrganizationVideos.Columns.OrganizationID.EQ(psql.Arg(organizationID)).And(
			models.OrganizationVideos.Columns.ID.In(psql.Arg(videoIDs...)),
		)),
	)
	if _, err := bob.Exec(ctx, pst.exec, q); err != nil {
		return fmt.Errorf("delete animal videos: %w", err)
	}

	return nil
}

func (pst *PgOrganizationPersistor) UpdateOrganizationVideo(ctx context.Context, organizationID, videoID int64, in models.OrganizationVideoSetter) (models.OrganizationVideo, error) {
	q := models.OrganizationVideos.Update(
		in.UpdateMod(),
		um.Where(models.OrganizationVideos.Columns.OrganizationID.EQ(psql.Arg(organizationID)).And(
			models.OrganizationVideos.Columns.ID.EQ(psql.Arg(videoID)),
		)),
		um.Returning(models.OrganizationVideos.Columns),
	)

	organizationVideo, err := bob.One(ctx, pst.exec, q, scan.StructMapper[models.OrganizationVideo]())
	if err != nil {
		return models.OrganizationVideo{}, fmt.Errorf("update animal video")
	}

	return organizationVideo, nil
}

func (pst *PgOrganizationPersistor) DeleteOrganizationVideo(ctx context.Context, organizationID, videoID int64) (int64, error) {
	q := models.OrganizationVideos.Delete(
		dm.Where(models.OrganizationVideos.Columns.OrganizationID.EQ(psql.Arg(organizationID)).And(
			models.OrganizationVideos.Columns.ID.EQ(psql.Arg(videoID)),
		)),
	)
	if _, err := bob.Exec(ctx, pst.exec, q); err != nil {
		return 0, fmt.Errorf("delete animal video: %w", err)
	}

	return videoID, nil
}
