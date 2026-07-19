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

func (pst *PgOrganizationPersistor) ListOrganizationPhotos(ctx context.Context, organizationID int64, filters dbtype.ListOrganizationPhotosFilters) (dbtype.PagedResult[models.OrganizationPhoto], error) {
	q := psql.Select(
		sm.Columns(models.OrganizationPhotos.Columns),
		sm.From(models.OrganizationPhotos.Name()),
		sm.Where(models.OrganizationPhotos.Columns.OrganizationID.EQ(psql.Arg(organizationID))),
		sm.GroupBy(models.OrganizationPhotos.Columns.ID),
	)
	addOrderBy(&q, filters.Sort, models.OrganizationPhotos.Columns.Names())
	addPagination(&q, filters.Page, filters.PageSize)

	if hasAnyLogicFilters(&filters.ListOrganizationPhotosParams) {
		if filters.ID != nil {
			ids := make([]any, len(*filters.ID))
			for i, id := range *filters.ID {
				ids[i] = id
			}

			q.Apply(sm.Where(models.OrganizationPhotos.Columns.ID.In(psql.Arg(ids...))))
		}
	}

	type ListOrganizationPhotosRow struct {
		models.OrganizationPhoto
		TotalCount int64
	}

	organizationPhotos, err := bob.All(ctx, pst.exec, q, scan.StructMapper[ListOrganizationPhotosRow]())
	if err != nil {
		return dbtype.PagedResult[models.OrganizationPhoto]{}, fmt.Errorf("query organization photos")
	}

	result := dbtype.PagedResult[models.OrganizationPhoto]{
		Data: make([]models.OrganizationPhoto, len(organizationPhotos)),
	}
	for i, row := range organizationPhotos {
		result.Data[i] = row.OrganizationPhoto
	}

	if len(organizationPhotos) > 0 {
		result.TotalCount = organizationPhotos[0].TotalCount
	}

	return result, nil
}

func (pst *PgOrganizationPersistor) GetOrganizationPhoto(ctx context.Context, organizationID, photoID int64) (models.OrganizationPhoto, error) {
	q := psql.Select(
		sm.Columns(models.OrganizationPhotos.Columns),
		sm.From(models.OrganizationPhotos.Name()),
		sm.Where(models.OrganizationPhotos.Columns.OrganizationID.EQ(psql.Arg(organizationID)).
			And(models.OrganizationPhotos.Columns.ID.EQ(psql.Arg(photoID))),
		),
	)

	organizationPhoto, err := bob.One(ctx, pst.exec, q, scan.StructMapper[models.OrganizationPhoto]())
	if err != nil {
		return models.OrganizationPhoto{}, fmt.Errorf("query organization photo")
	}

	return organizationPhoto, nil
}

func (pst *PgOrganizationPersistor) CreateOrganizationPhotos(ctx context.Context, organizationID int64, in []models.OrganizationPhotoSetter) ([]models.OrganizationPhoto, error) {
	setters := make([]*models.OrganizationPhotoSetter, len(in))
	for i, x := range in {
		x.OrganizationID = omitnull.From(organizationID)
		setters[i] = &x
	}

	q := models.OrganizationPhotos.Insert(bob.ToMods(setters...), im.Returning(models.OrganizationPhotos.Columns))

	organizationPhotos, err := bob.All(ctx, pst.exec, q, scan.StructMapper[models.OrganizationPhoto]())
	if err != nil {
		return nil, fmt.Errorf("insert organization photos")
	}

	return organizationPhotos, nil
}

func (pst *PgOrganizationPersistor) DeleteOrganizationPhotos(ctx context.Context, organizationID int64, ids []int64) error {
	photoIDs := make([]any, len(ids))
	for i, id := range ids {
		photoIDs[i] = id
	}

	q := models.OrganizationPhotos.Delete(
		dm.Where(models.OrganizationPhotos.Columns.OrganizationID.EQ(psql.Arg(organizationID)).And(
			models.OrganizationPhotos.Columns.ID.In(psql.Arg(photoIDs...)),
		)),
	)
	if _, err := bob.Exec(ctx, pst.exec, q); err != nil {
		return fmt.Errorf("delete organization photos: %w", err)
	}

	return nil
}

func (pst *PgOrganizationPersistor) UpdateOrganizationPhoto(ctx context.Context, organizationID, photoID int64, in models.OrganizationPhotoSetter) (models.OrganizationPhoto, error) {
	q := models.OrganizationPhotos.Update(
		in.UpdateMod(),
		um.Where(models.OrganizationPhotos.Columns.OrganizationID.EQ(psql.Arg(organizationID)).And(
			models.OrganizationPhotos.Columns.ID.EQ(psql.Arg(photoID)),
		)),
		um.Returning(models.OrganizationPhotos.Columns),
	)

	organizationPhoto, err := bob.One(ctx, pst.exec, q, scan.StructMapper[models.OrganizationPhoto]())
	if err != nil {
		return models.OrganizationPhoto{}, fmt.Errorf("update organization photo")
	}

	return organizationPhoto, nil
}

func (pst *PgOrganizationPersistor) DeleteOrganizationPhoto(ctx context.Context, organizationID, photoID int64) (int64, error) {
	q := models.OrganizationPhotos.Delete(
		dm.Where(models.OrganizationPhotos.Columns.OrganizationID.EQ(psql.Arg(organizationID)).And(
			models.OrganizationPhotos.Columns.ID.EQ(psql.Arg(photoID)),
		)),
	)
	if _, err := bob.Exec(ctx, pst.exec, q); err != nil {
		return 0, fmt.Errorf("delete organization photo: %w", err)
	}

	return photoID, nil
}
