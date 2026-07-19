package postgres

import (
	"context"
	"errors"
	"fmt"

	api "github.com/dankobg/fluffly/api/gen"
	"github.com/dankobg/fluffly/db/gen/models"
	"github.com/dankobg/fluffly/persistence"
	"github.com/dankobg/fluffly/persistence/dbtype"
	"github.com/stephenafamo/bob"
	"github.com/stephenafamo/bob/dialect/psql"
	"github.com/stephenafamo/bob/dialect/psql/sm"
	"github.com/stephenafamo/bob/orm"
	"github.com/stephenafamo/scan"
)

var _ persistence.AdoptionPersistor = (*PgAdoptionPersistor)(nil)

type PgAdoptionPersistor struct {
	*PgPersistor
}

func NewPgAdoptionPersistor(ps *PgPersistor) *PgAdoptionPersistor {
	return &PgAdoptionPersistor{
		PgPersistor: ps,
	}
}

var ErrAdoptionNotFound = errors.New("adoption not found")

func (pst *PgAdoptionPersistor) ListAdoptions(ctx context.Context, filters dbtype.ListAdoptionsFilters) (dbtype.PagedResult[dbtype.AdoptionWithJoinData], error) {
	q := psql.Select(
		sm.Columns(models.Adoptions.Columns),
		sm.From(models.Adoptions.Name()),
		sm.GroupBy(models.Adoptions.Columns.ID),
	)
	addOrderBy(&q, filters.Sort, models.Adoptions.Columns.Names())
	addPagination(&q, filters.Page, filters.PageSize)

	if hasAnyLogicFilters(&filters.ListAdoptionsParams) {
		if filters.Status != nil {
			statuses := make([]any, len(*filters.Status))
			for i, age := range *filters.Status {
				statuses[i] = age
			}

			q.Apply(sm.Where(models.Adoptions.Columns.Status.In(psql.Arg(statuses...))))
		}
		if filters.IsPermanent != nil {
			q.Apply(sm.Where(models.Adoptions.Columns.IsPermanent.EQ(psql.Arg(*filters.IsPermanent))))
		}
		if filters.AnimalID != nil {
			animalIDs := make([]any, len(*filters.AnimalID))
			for i, id := range *filters.AnimalID {
				animalIDs[i] = id
			}

			q.Apply(sm.Where(models.Adoptions.Columns.AnimalID.In(psql.Arg(animalIDs...))))
		}
		if filters.OrganizationID != nil {
			organizationIDs := make([]any, len(*filters.OrganizationID))
			for i, id := range *filters.OrganizationID {
				organizationIDs[i] = id
			}

			q.Apply(sm.Where(models.Adoptions.Columns.OrganizationID.In(psql.Arg(organizationIDs...))))
		}
		if filters.UserID != nil {
			userIDs := make([]any, len(*filters.UserID))
			for i, id := range *filters.UserID {
				userIDs[i] = id
			}

			q.Apply(sm.Where(models.Adoptions.Columns.UserID.In(psql.Arg(userIDs...))))
		}
		if filters.Embed != nil {
			for _, embed := range *filters.Embed {
				switch embed {
				case api.ListAdoptionsParamsEmbedOrganization:
					q.Apply(sm.Columns(
						models.Organizations.Columns.WithPrefix("organization."),
						models.OrganizationContacts.Columns.WithPrefix("organization_contact."),
						models.OrganizationWorkHours.Columns.WithPrefix("organization_work_hour."),
						models.Addresses.Columns.Except("coords").WithPrefix("organization_contact_address."),
						psql.Raw(`ST_AsBinary("address"."coords") AS "organization_contact_address.coords"`),
						models.Countries.Columns.WithPrefix("organization_contact_address_country."),
					))
					q.Apply(
						sm.LeftJoin(models.Organizations.Name()).
							On(models.Organizations.Columns.ID.EQ(models.Adoptions.Columns.OrganizationID)),
						sm.LeftJoin(models.OrganizationContacts.Name()).
							On(models.Organizations.Columns.ID.EQ(models.OrganizationContacts.Columns.OrganizationID)),
						sm.LeftJoin(models.Addresses.Name()).
							On(models.Addresses.Columns.ID.EQ(models.OrganizationContacts.Columns.AddressID)),
						sm.LeftJoin(models.Countries.Name()).
							On(models.Countries.Columns.ID.EQ(models.Addresses.Columns.CountryID)),
						sm.LeftJoin(models.OrganizationWorkHours.Name()).
							On(models.Organizations.Columns.ID.EQ(models.OrganizationWorkHours.Columns.OrganizationID)),
					)
					q.Apply(sm.GroupBy(models.Organizations.Columns.ID))
					q.Apply(sm.GroupBy(models.OrganizationContacts.Columns.ID))
					q.Apply(sm.GroupBy(models.Addresses.Columns.ID))
					q.Apply(sm.GroupBy(models.Countries.Columns.ID))
					q.Apply(sm.GroupBy(models.OrganizationWorkHours.Columns.ID))

				case api.ListAdoptionsParamsEmbedUser:
					// fetched via kratos

				}
			}
		}

		if filters.ID != nil {
			ids := make([]any, len(*filters.ID))
			for i, id := range *filters.ID {
				ids[i] = id
			}

			q.Apply(sm.Where(models.Adoptions.Columns.ID.In(psql.Arg(ids...))))
		}

		if filters.Status != nil {
			statuses := make([]any, len(*filters.Status))
			for i, id := range *filters.Status {
				statuses[i] = id
			}

			q.Apply(sm.Where(models.Adoptions.Columns.Status.In(psql.Arg(statuses...))))
		}
	}

	type ListAdoptionsRow struct {
		dbtype.AdoptionWithJoinData
		TotalCount int64
	}

	adoptions, err := bob.All(ctx, pst.exec, q, scan.StructMapper[ListAdoptionsRow](scan.WithTypeConverter(orm.NullTypeConverter{})))
	if err != nil {
		return dbtype.PagedResult[dbtype.AdoptionWithJoinData]{}, fmt.Errorf("query adoptions")
	}

	result := dbtype.PagedResult[dbtype.AdoptionWithJoinData]{
		Data: make([]dbtype.AdoptionWithJoinData, len(adoptions)),
	}
	for i, row := range adoptions {
		result.Data[i] = row.AdoptionWithJoinData
	}

	if len(adoptions) > 0 {
		result.TotalCount = adoptions[0].TotalCount
	}

	return result, nil
}

func (pst *PgAdoptionPersistor) GetAdoptionByID(ctx context.Context, adoptionID int64, filters dbtype.GetAdoptionByIDFilters) (dbtype.AdoptionWithJoinData, error) {
	q := psql.Select(
		sm.Columns(models.Adoptions.Columns),
		sm.From(models.Adoptions.Name()),
		sm.Where(models.Adoptions.Columns.ID.EQ(psql.Arg(adoptionID))),
		sm.GroupBy(models.Adoptions.Columns.ID),
	)

	if hasAnyLogicFilters(&filters.GetAdoptionParams) {
		if filters.Embed != nil {
			for _, embed := range *filters.Embed {
				switch embed {
				case api.GetAdoptionParamsEmbedOrganization:
					q.Apply(sm.Columns(
						models.Organizations.Columns.WithPrefix("organization."),
						models.OrganizationContacts.Columns.WithPrefix("organization_contact."),
						models.OrganizationWorkHours.Columns.WithPrefix("organization_work_hour."),
						models.Addresses.Columns.Except("coords").WithPrefix("organization_contact_address."),
						psql.Raw(`ST_AsBinary("address"."coords") AS "organization_contact_address.coords"`),
						models.Countries.Columns.WithPrefix("organization_contact_address_country."),
					))
					q.Apply(
						sm.LeftJoin(models.Organizations.Name()).
							On(models.Organizations.Columns.ID.EQ(models.Adoptions.Columns.OrganizationID)),
						sm.LeftJoin(models.OrganizationContacts.Name()).
							On(models.Organizations.Columns.ID.EQ(models.OrganizationContacts.Columns.OrganizationID)),
						sm.LeftJoin(models.Addresses.Name()).
							On(models.Addresses.Columns.ID.EQ(models.OrganizationContacts.Columns.AddressID)),
						sm.LeftJoin(models.Countries.Name()).
							On(models.Countries.Columns.ID.EQ(models.Addresses.Columns.CountryID)),
						sm.LeftJoin(models.OrganizationWorkHours.Name()).
							On(models.Organizations.Columns.ID.EQ(models.OrganizationWorkHours.Columns.OrganizationID)),
					)
					q.Apply(sm.GroupBy(models.Organizations.Columns.ID))
					q.Apply(sm.GroupBy(models.OrganizationContacts.Columns.ID))
					q.Apply(sm.GroupBy(models.Addresses.Columns.ID))
					q.Apply(sm.GroupBy(models.Countries.Columns.ID))
					q.Apply(sm.GroupBy(models.OrganizationWorkHours.Columns.ID))

				case api.GetAdoptionParamsEmbedUser:
					// fetched via kratos
				}
			}
		}
	}

	adoptionWithJoinData, err := bob.One(ctx, pst.exec, q, scan.StructMapper[dbtype.AdoptionWithJoinData](scan.WithTypeConverter(orm.NullTypeConverter{})))
	if err != nil {
		return dbtype.AdoptionWithJoinData{}, fmt.Errorf("query adoption")
	}

	return adoptionWithJoinData, nil
}
