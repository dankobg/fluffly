package postgres

import (
	"context"
	"fmt"

	"github.com/dankobg/fluffly/db/gen/test/public/model"
	t "github.com/dankobg/fluffly/db/gen/test/public/table"
	"github.com/dankobg/fluffly/persistence"
	p "github.com/go-jet/jet/v2/postgres"
)

var _ persistence.CountryPersistor = (*PgCountryPersistor)(nil)

type PgCountryPersistor struct {
	*PgPersistor
}

func NewPgCountryPersistor(ps *PgPersistor) *PgCountryPersistor {
	return &PgCountryPersistor{
		PgPersistor: ps,
	}
}

func (pc *PgCountryPersistor) ListCountries(ctx context.Context, filters persistence.CountryFilters) (persistence.PagedResult[model.Country], error) {
	q := p.SELECT(
		t.Country.AllColumns,
		getSelectTotalCount(filters.Pagination),
	).
		FROM(t.Country)
	q = getLimitOffset(q, filters.Pagination)

	var dest []struct {
		model.Country
		TotalCount int64 `db:"total_count"`
	}
	if err := q.QueryContext(ctx, pc.db, &dest); err != nil {
		return persistence.PagedResult[model.Country]{}, err
	}

	result := persistence.PagedResult[model.Country]{
		Data: make([]model.Country, len(dest)),
	}
	for i, row := range dest {
		result.Data[i] = row.Country
	}
	if len(dest) > 0 {
		result.TotalCount = dest[0].TotalCount
	}

	return result, nil
}

func (pc *PgCountryPersistor) GetCountryByID(ctx context.Context, countryID int64) (model.Country, error) {
	q := p.SELECT(t.Country.AllColumns).
		FROM(t.Country).
		WHERE(t.Country.ID.EQ(p.Int64(countryID)))
	var dest model.Country
	if err := q.QueryContext(ctx, pc.db, &dest); err != nil {
		return model.Country{}, err
	}
	return dest, nil
}

func (pc *PgCountryPersistor) DeleteCountryByID(ctx context.Context, countryID int64) (int64, error) {
	q := t.Country.DELETE().WHERE(t.Country.ID.EQ(p.Int64(countryID)))
	if _, err := q.ExecContext(ctx, pc.db); err != nil {
		return 0, fmt.Errorf("failed to delete an country: %w", err)
	}
	return countryID, nil
}

func (pc *PgCountryPersistor) CreateCountry(ctx context.Context, in persistence.CountrySetter) (model.Country, error) {
	cols, m := in.ToModel()
	q := t.Country.INSERT(cols).
		MODEL(m).
		RETURNING(t.Country.AllColumns)

	var dest model.Country
	if err := q.QueryContext(ctx, pc.db, &dest); err != nil {
		return dest, fmt.Errorf("failed to create an country: %w", err)
	}
	return dest, nil
}

func (pc *PgCountryPersistor) UpdateCountry(ctx context.Context, countryID int64, in persistence.CountrySetter) (model.Country, error) {
	cols, m := in.ToModel(true)
	q := t.Country.UPDATE(cols).
		MODEL(m).
		WHERE(t.Country.ID.EQ(p.Int64(countryID))).
		RETURNING(t.Country.AllColumns)

	var dest model.Country
	if err := q.QueryContext(ctx, pc.db, &dest); err != nil {
		return dest, fmt.Errorf("failed to update an country: %w", err)
	}
	return dest, nil
}
