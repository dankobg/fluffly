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

func (pst *PgOrganizationPersistor) ListOrganizationSocials(ctx context.Context, animalID int64, filters dbtype.ListOrganizationSocialsFilters) (dbtype.PagedResult[models.OrganizationSocial], error) {
	q := psql.Select(
		sm.Columns(models.OrganizationSocials.Columns),
		sm.From(models.OrganizationSocials.Name()),
		sm.GroupBy(models.OrganizationSocials.Columns.ID),
	)
	addOrderBy(&q, filters.Sort, models.OrganizationSocials.Columns.Names())
	addPagination(&q, filters.Page, filters.PageSize)

	if hasAnyLogicFilters(&filters.ListOrganizationSocialsParams) {
		if filters.ID != nil {
			ids := make([]any, len(*filters.ID))
			for i, id := range *filters.ID {
				ids[i] = id
			}

			q.Apply(sm.Where(models.OrganizationSocials.Columns.ID.In(psql.Arg(ids...))))
		}
	}

	type ListOrganizationSocialsRow struct {
		models.OrganizationSocial
		TotalCount int64
	}

	organizationSocials, err := bob.All(ctx, pst.exec, q, scan.StructMapper[ListOrganizationSocialsRow]())
	if err != nil {
		return dbtype.PagedResult[models.OrganizationSocial]{}, fmt.Errorf("query animal photos")
	}

	result := dbtype.PagedResult[models.OrganizationSocial]{
		Data: make([]models.OrganizationSocial, len(organizationSocials)),
	}
	for i, row := range organizationSocials {
		result.Data[i] = row.OrganizationSocial
	}

	if len(organizationSocials) > 0 {
		result.TotalCount = organizationSocials[0].TotalCount
	}

	return result, nil
}

func (pst *PgOrganizationPersistor) GetOrganizationSocial(ctx context.Context, animalID, photoID int64) (models.OrganizationSocial, error) {
	q := psql.Select(
		sm.Columns(models.OrganizationSocials.Columns),
		sm.From(models.OrganizationSocials.Name()),
		sm.Where(models.OrganizationSocials.Columns.ID.EQ(psql.Arg(photoID))),
	)

	organizationSocial, err := bob.One(ctx, pst.exec, q, scan.StructMapper[models.OrganizationSocial]())
	if err != nil {
		return models.OrganizationSocial{}, fmt.Errorf("query animal photo")
	}

	return organizationSocial, nil
}

func (pst *PgOrganizationPersistor) CreateOrganizationSocials(ctx context.Context, animalID int64, in []models.OrganizationSocialSetter) ([]models.OrganizationSocial, error) {
	setters := make([]*models.OrganizationSocialSetter, len(in))
	for i, x := range in {
		setters[i] = &x
	}

	q := models.OrganizationSocials.Insert(bob.ToMods(setters...), im.Returning(models.OrganizationSocials.Columns))

	organizationSocials, err := bob.All(ctx, pst.exec, q, scan.StructMapper[models.OrganizationSocial]())
	if err != nil {
		return nil, fmt.Errorf("insert animal photos")
	}

	return organizationSocials, nil
}

func (pst *PgOrganizationPersistor) DeleteOrganizationSocials(ctx context.Context, animalID int64, ids []int64) error {
	photoIDs := make([]any, len(ids))
	for i, id := range ids {
		photoIDs[i] = id
	}

	q := models.OrganizationSocials.Delete(
		dm.Where(models.OrganizationSocials.Columns.OrganizationID.EQ(psql.Arg(animalID)).And(
			models.OrganizationSocials.Columns.ID.In(psql.Arg(photoIDs...)),
		)),
	)
	if _, err := bob.Exec(ctx, pst.exec, q); err != nil {
		return fmt.Errorf("delete animal photos: %w", err)
	}

	return nil
}

func (pst *PgOrganizationPersistor) UpdateOrganizationSocial(ctx context.Context, animalID, photoID int64, in models.OrganizationSocialSetter) (models.OrganizationSocial, error) {
	q := models.OrganizationSocials.Update(
		in.UpdateMod(),
		um.Where(models.OrganizationSocials.Columns.OrganizationID.EQ(psql.Arg(animalID)).And(
			models.OrganizationSocials.Columns.ID.EQ(psql.Arg(photoID)),
		)),
		um.Returning(models.OrganizationSocials.Columns),
	)

	organizationSocial, err := bob.One(ctx, pst.exec, q, scan.StructMapper[models.OrganizationSocial]())
	if err != nil {
		return models.OrganizationSocial{}, fmt.Errorf("update animal photo")
	}

	return organizationSocial, nil
}

func (pst *PgOrganizationPersistor) DeleteOrganizationSocial(ctx context.Context, animalID, photoID int64) (int64, error) {
	q := models.OrganizationSocials.Delete(
		dm.Where(models.OrganizationSocials.Columns.OrganizationID.EQ(psql.Arg(animalID)).And(
			models.OrganizationSocials.Columns.ID.EQ(psql.Arg(photoID)),
		)),
	)
	if _, err := bob.Exec(ctx, pst.exec, q); err != nil {
		return 0, fmt.Errorf("delete animal photo: %w", err)
	}

	return photoID, nil
}
