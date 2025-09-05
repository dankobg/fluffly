package postgres

import (
	"context"
	"errors"
	"fmt"

	"github.com/dankobg/fluffly/db/gen/test/public/model"
	t "github.com/dankobg/fluffly/db/gen/test/public/table"
	"github.com/dankobg/fluffly/persistence"
	"github.com/dankobg/fluffly/persistence/dbtype"
	p "github.com/go-jet/jet/v2/postgres"
	"github.com/go-jet/jet/v2/qrm"
	"github.com/jackc/pgerrcode"
	"github.com/jackc/pgx/v5/pgconn"
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

type ErrCountryIntegrityViolation struct{ errIntegrityViolation }
type ErrCountryUniqueViolation struct{ errUniqueViolation }

var (
	ErrCountryNotFound         = errors.New("country not found")
	errCountryIntegrity        = ErrCountryIntegrityViolation{}
	errCountryUniqueName       = ErrCountryUniqueViolation{errUniqueViolation: errUniqueViolation{Name: "name"}}
	errCountryUniqueIsoAlpha2  = ErrCountryUniqueViolation{errUniqueViolation: errUniqueViolation{Name: "iso_alpha2"}}
	errCountryUniqueIsoAlpha3  = ErrCountryUniqueViolation{errUniqueViolation: errUniqueViolation{Name: "iso_alpha3"}}
	errCountryUniqueIsoNumeric = ErrCountryUniqueViolation{errUniqueViolation: errUniqueViolation{Name: "iso_numeric"}}
)

func convertCountryPgError(pgErr *pgconn.PgError) error {
	if pgerrcode.IsIntegrityConstraintViolation(pgErr.Code) {
		switch pgErr.ConstraintName {
		case "uq_country_name":
			return errCountryUniqueName
		case "uq_country_iso_alpha2":
			return errCountryUniqueIsoAlpha2
		case "uq_country_iso_alpha3":
			return errCountryUniqueIsoAlpha3
		case "uq_country_iso_numeric":
			return errCountryUniqueIsoNumeric
		}
		return errCountryIntegrity
	}
	return pgErr
}

func (pc *PgCountryPersistor) ListCountries(ctx context.Context, filters dbtype.CountryFilters) (dbtype.PagedResult[model.Country], error) {
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
		return dbtype.PagedResult[model.Country]{}, fmt.Errorf("could not query countries")
	}

	result := dbtype.PagedResult[model.Country]{
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
		if errors.Is(err, qrm.ErrNoRows) {
			return dest, ErrCountryNotFound
		}
		return dest, fmt.Errorf("could not query a country")
	}
	return dest, nil
}

func (pc *PgCountryPersistor) DeleteCountryByID(ctx context.Context, countryID int64) (int64, error) {
	q := t.Country.DELETE().WHERE(t.Country.ID.EQ(p.Int64(countryID)))
	if _, err := q.ExecContext(ctx, pc.db); err != nil {
		return 0, fmt.Errorf("could not delete a country: %w", err)
	}
	return countryID, nil
}

func (pc *PgCountryPersistor) CreateCountry(ctx context.Context, in dbtype.CountrySetter) (model.Country, error) {
	cols, m := in.ToModel()
	q := t.Country.INSERT(cols).
		MODEL(m).
		RETURNING(t.Country.AllColumns)

	var dest model.Country
	if err := q.QueryContext(ctx, pc.db, &dest); err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) {
			return dest, convertCountryPgError(pgErr)
		}
		return dest, fmt.Errorf("could not insert a country: %w", err)
	}
	return dest, nil
}

func (pc *PgCountryPersistor) UpdateCountry(ctx context.Context, countryID int64, in dbtype.CountrySetter) (model.Country, error) {
	cols, m := in.ToModel(true)
	q := t.Country.UPDATE(cols).
		MODEL(m).
		WHERE(t.Country.ID.EQ(p.Int64(countryID))).
		RETURNING(t.Country.AllColumns)

	var dest model.Country
	if err := q.QueryContext(ctx, pc.db, &dest); err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) {
			return dest, convertCountryPgError(pgErr)
		}
		return dest, fmt.Errorf("could not update a country: %w", err)
	}
	return dest, nil
}
