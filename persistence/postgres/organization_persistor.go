package postgres

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/dankobg/fluffly/db/gen/test/public/model"
	t "github.com/dankobg/fluffly/db/gen/test/public/table"
	"github.com/dankobg/fluffly/persistence"
	p "github.com/go-jet/jet/v2/postgres"
	"github.com/oapi-codegen/nullable"
)

var _ persistence.OrganizationPersistor = (*PgOrganizationPersistor)(nil)

type PgOrganizationPersistor struct {
	*PgPersistor
}

func NewPgOrganizationPersistor(ps *PgPersistor) *PgOrganizationPersistor {
	return &PgOrganizationPersistor{
		PgPersistor: ps,
	}
}

func (po *PgOrganizationPersistor) ListOrganizations(ctx context.Context, filters persistence.OrganizationFilters) (persistence.PagedResult[persistence.OrganizationWithJoinData], error) {
	q := p.SELECT(
		t.Organization.AllColumns,
		t.OrganizationContact.AllColumns,
		t.Address.AllColumns,
		t.Country.AllColumns,
		t.OrganizationWorkHour.AllColumns,
		p.SELECT_JSON_ARR(
			t.OrganizationPhoto.AllColumns).
			FROM(t.OrganizationPhoto).
			WHERE(t.OrganizationPhoto.OrganizationID.EQ(t.Organization.ID)).
			AS("photos"),
		p.SELECT_JSON_ARR(
			t.OrganizationSocial.AllColumns).
			FROM(t.OrganizationSocial).
			WHERE(t.OrganizationSocial.OrganizationID.EQ(t.Organization.ID)).
			AS("socials"),
		getSelectTotalCount(filters.Pagination),
	).
		FROM(
			t.Organization.
				LEFT_JOIN(t.OrganizationContact, t.Organization.ID.EQ(t.OrganizationContact.OrganizationID)).
				LEFT_JOIN(t.Address, t.Address.ID.EQ(t.OrganizationContact.AddressID)).
				LEFT_JOIN(t.Country, t.Country.ID.EQ(t.Address.CountryID)).
				LEFT_JOIN(t.OrganizationWorkHour, t.Organization.ID.EQ(t.OrganizationWorkHour.OrganizationID)),
		).
		GROUP_BY(
			t.Organization.ID,
			t.OrganizationContact.ID,
			t.Address.ID,
			t.Country.ID,
			t.OrganizationWorkHour.ID,
		)
	q = getLimitOffset(q, filters.Pagination)

	var dest []struct {
		persistence.OrganizationWithJoinData
		TotalCount int64 `db:"total_count"`
	}
	if err := q.QueryContext(ctx, po.db, &dest); err != nil {
		return persistence.PagedResult[persistence.OrganizationWithJoinData]{}, err
	}
	result := persistence.PagedResult[persistence.OrganizationWithJoinData]{
		Data: make([]persistence.OrganizationWithJoinData, len(dest)),
	}
	for i, row := range dest {
		result.Data[i] = row.OrganizationWithJoinData
	}
	if len(dest) > 0 {
		result.TotalCount = dest[0].TotalCount
	}
	return result, nil
}

func (po *PgOrganizationPersistor) GetOrganizationByID(ctx context.Context, organizationID int64) (persistence.OrganizationWithJoinData, error) {
	q := p.SELECT(
		t.Organization.AllColumns,
		t.OrganizationContact.AllColumns,
		t.Address.AllColumns,
		t.Country.AllColumns,
		t.OrganizationWorkHour.AllColumns,
		p.SELECT_JSON_ARR(
			t.OrganizationPhoto.AllColumns).
			FROM(t.OrganizationPhoto).
			WHERE(t.OrganizationPhoto.OrganizationID.EQ(t.Organization.ID)).
			AS("photos"),
		p.SELECT_JSON_ARR(
			t.OrganizationSocial.AllColumns).
			FROM(t.OrganizationSocial).
			WHERE(t.OrganizationSocial.OrganizationID.EQ(t.Organization.ID)).
			AS("socials"),
	).
		FROM(
			t.Organization.
				LEFT_JOIN(t.OrganizationContact, t.Organization.ID.EQ(t.OrganizationContact.OrganizationID)).
				LEFT_JOIN(t.Address, t.Address.ID.EQ(t.OrganizationContact.AddressID)).
				LEFT_JOIN(t.Country, t.Country.ID.EQ(t.Address.CountryID)).
				LEFT_JOIN(t.OrganizationWorkHour, t.Organization.ID.EQ(t.OrganizationWorkHour.OrganizationID)),
		).
		WHERE(t.Organization.ID.EQ(p.Int64(organizationID))).
		GROUP_BY(
			t.Organization.ID,
			t.OrganizationContact.ID,
			t.Address.ID,
			t.Country.ID,
			t.OrganizationWorkHour.ID,
		)
	var dest persistence.OrganizationWithJoinData
	if err := q.QueryContext(ctx, po.db, &dest); err != nil {
		return persistence.OrganizationWithJoinData{}, err
	}
	return dest, nil
}

func (po *PgOrganizationPersistor) CreateOrganization(ctx context.Context, in persistence.OrganizationCreateSetter) (model.Organization, error) {
	var insertedOrganization model.Organization

	txErr := po.WithTx(ctx, func(tx *sql.Tx) error {
		orgCols, org := in.Organization.ToModel()
		q1 := t.Organization.INSERT(orgCols).
			MODEL(org).
			RETURNING(t.Organization.AllColumns)
		if err := q1.QueryContext(ctx, tx, &insertedOrganization); err != nil {
			return fmt.Errorf("failed to insert an organization: %w", err)
		}

		var insertedAddress model.Address
		addrCols, addr := in.Address.ToModel()
		q2 := t.Address.INSERT(addrCols).
			MODEL(addr).
			RETURNING(t.Address.AllColumns)
		if err := q2.QueryContext(ctx, tx, &insertedAddress); err != nil {
			return fmt.Errorf("failed to insert organization address: %w", err)
		}

		var insertedContact model.OrganizationContact
		in.Contact.OrganizationID = nullable.NewNullableWithValue(insertedOrganization.ID)
		in.Contact.AddressID = nullable.NewNullableWithValue(insertedAddress.ID)
		contactCols, contact := in.Contact.ToModel()
		q3 := t.OrganizationContact.INSERT(contactCols).
			MODEL(contact).
			RETURNING(t.OrganizationContact.AllColumns)
		if err := q3.QueryContext(ctx, tx, &insertedContact); err != nil {
			return fmt.Errorf("failed to insert organization contact: %w", err)
		}

		if in.WorkHour.IsSpecified() && !in.WorkHour.IsNull() {
			var insertedWorkHour model.OrganizationWorkHour
			inWorkHour := in.WorkHour.MustGet()
			inWorkHour.OrganizationID = nullable.NewNullableWithValue(insertedOrganization.ID)
			workHourCols, workHour := inWorkHour.ToModel()
			q4 := t.OrganizationWorkHour.INSERT(workHourCols).
				MODEL(workHour).
				RETURNING(t.OrganizationWorkHour.AllColumns)
			if err := q4.QueryContext(ctx, tx, &insertedWorkHour); err != nil {
				return fmt.Errorf("failed to insert organization work hour: %w", err)
			}
		}

		if in.Photos.IsSpecified() && !in.Photos.IsNull() {
			photosInput := in.Photos.MustGet()
			if len(photosInput) > 0 {
				var insertedPhotos []model.OrganizationPhoto
				photos := make([]model.OrganizationPhoto, len(photosInput))
				var photoCols p.ColumnList
				for i, photo := range photosInput {
					photo.OrganizationID = nullable.NewNullableWithValue(insertedOrganization.ID)
					cols, m := photo.ToModel()
					if i == 0 {
						photoCols = cols
					}
					photos[i] = m
				}
				q5 := t.OrganizationPhoto.INSERT(photoCols).
					MODELS(photos).
					RETURNING(t.OrganizationPhoto.AllColumns)
				if err := q5.QueryContext(ctx, tx, &insertedPhotos); err != nil {
					return fmt.Errorf("failed to insert organization photos: %w", err)
				}
			}
		}

		if in.Socials.IsSpecified() && !in.Socials.IsNull() {
			socialsInput := in.Socials.MustGet()
			if len(socialsInput) > 0 {
				var insertedSocials []model.OrganizationSocial
				socials := make([]model.OrganizationSocial, len(socialsInput))
				var socialsCols p.ColumnList
				for i, social := range socialsInput {
					social.OrganizationID = nullable.NewNullableWithValue(insertedOrganization.ID)
					cols, m := social.ToModel()
					if i == 0 {
						socialsCols = cols
					}
					socials[i] = m
				}
				q6 := t.OrganizationSocial.INSERT(socialsCols).
					MODELS(socials).
					RETURNING(t.OrganizationSocial.AllColumns)
				if err := q6.QueryContext(ctx, tx, &insertedSocials); err != nil {
					return fmt.Errorf("failed to insert organization social platforms: %w", err)
				}
			}
		}

		return nil
	})
	if txErr != nil {
		return model.Organization{}, txErr
	}
	return insertedOrganization, nil
}

func (po *PgOrganizationPersistor) UpdateOrganization(ctx context.Context, organizationID int64, in persistence.OrganizationSetter) (model.Organization, error) {
	cols, m := in.ToModel(true)

	q := t.Organization.UPDATE(cols).
		MODEL(m).
		WHERE(t.Organization.ID.EQ(p.Int64(organizationID))).
		RETURNING(t.Organization.AllColumns)

	var dest model.Organization
	if err := q.QueryContext(ctx, po.db, &dest); err != nil {
		return dest, fmt.Errorf("failed to update an organization: %w", err)
	}
	return dest, nil
}

func (po *PgOrganizationPersistor) DeleteOrganizationByID(ctx context.Context, organizationID int64) (int64, error) {
	q := t.Organization.DELETE().WHERE(t.Organization.ID.EQ(p.Int64(organizationID)))
	if _, err := q.ExecContext(ctx, po.db); err != nil {
		return 0, fmt.Errorf("failed to delete an organization: %w", err)
	}
	return organizationID, nil
}

// type Filters struct {
// 	WithContact  bool
// 	WithWorkHour bool
// 	WithPhotos   bool
// 	WithSocials  bool
// }

// func (po *PgOrganizationPersistor) DYNAMIC_LIST_ORGANIZATIONS(ctx context.Context, filters Filters) ([]persistence.OrganizationWithJoinData, error) {
// 	var (
// 		selectClause  p.ProjectionList
// 		fromClause    p.ReadableTable = t.Organization
// 		groupByClause                 = []p.GroupByClause{t.Organization.ID}
// 	)

// 	if filters.WithContact {
// 		selectClause = append(selectClause, t.OrganizationContact.AllColumns, t.Address.AllColumns, t.Country.AllColumns)
// 		fromClause = fromClause.
// 			LEFT_JOIN(t.OrganizationContact, t.Organization.ID.EQ(t.OrganizationContact.OrganizationID)).
// 			LEFT_JOIN(t.Address, t.Address.ID.EQ(t.OrganizationContact.AddressID)).
// 			LEFT_JOIN(t.Country, t.Country.ID.EQ(t.Address.CountryID))
// 		groupByClause = append(groupByClause, t.OrganizationContact.ID, t.Address.ID, t.Country.ID)
// 	}
// 	if filters.WithWorkHour {
// 		selectClause = append(selectClause, t.OrganizationWorkHour.AllColumns)
// 		fromClause = fromClause.LEFT_JOIN(t.OrganizationWorkHour, t.Organization.ID.EQ(t.OrganizationWorkHour.OrganizationID))
// 		groupByClause = append(groupByClause, t.OrganizationWorkHour.ID)
// 	}
// 	if filters.WithPhotos {
// 		selectClause = append(selectClause, t.OrganizationPhoto.AllColumns)
// 		fromClause = fromClause.LEFT_JOIN(t.OrganizationPhoto, t.Organization.ID.EQ(t.OrganizationPhoto.OrganizationID))
// 		groupByClause = append(groupByClause, t.OrganizationPhoto.ID)
// 	}
// 	if filters.WithSocials {
// 		selectClause = append(selectClause, t.OrganizationSocial.AllColumns)
// 		fromClause = fromClause.LEFT_JOIN(t.OrganizationSocial, t.Organization.ID.EQ(t.OrganizationSocial.OrganizationID))
// 		groupByClause = append(groupByClause, t.OrganizationSocial.OrganizationID, t.OrganizationSocial.Platform)
// 	}

// 	q := p.SELECT(t.Organization.AllColumns, selectClause...).
// 		FROM(fromClause).
// 		GROUP_BY(groupByClause...)

// 	var dest []persistence.OrganizationWithJoinData
// 	if err := q.QueryContext(ctx, po.db, &dest); err != nil {
// 		return nil, err
// 	}
// 	return dest, nil
// }
