package postgres

import (
	"context"
	"errors"
	"fmt"

	"github.com/dankobg/fluffly/db/gen/models"
	"github.com/dankobg/fluffly/persistence"
	"github.com/dankobg/fluffly/persistence/dbtype"
	"github.com/jackc/pgerrcode"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/stephenafamo/bob"
	"github.com/stephenafamo/bob/dialect/psql"
	"github.com/stephenafamo/bob/dialect/psql/dm"
	"github.com/stephenafamo/bob/dialect/psql/im"
	"github.com/stephenafamo/bob/dialect/psql/sm"
	"github.com/stephenafamo/bob/dialect/psql/um"
	"github.com/stephenafamo/scan"
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

type (
	ErrCountryIntegrityViolation struct{ errIntegrityViolation }
	ErrCountryUniqueViolation    struct{ errUniqueViolation }
)

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

func (pst *PgCountryPersistor) ListCountries(ctx context.Context, filters dbtype.ListCountriesFilters) (dbtype.PagedResult[models.Country], error) {
	q := psql.Select(
		sm.Columns(models.Countries.Columns),
		sm.From(models.Countries.Name()),
		sm.GroupBy(models.Countries.Columns.ID),
	)
	addOrderBy(&q, filters.Sort, models.Countries.Columns.Names())
	addPagination(&q, filters.Page, filters.PageSize)

	if hasAnyLogicFilters(&filters.ListCountriesParams) {
		if filters.ID != nil {
			ids := make([]any, len(*filters.ID))
			for i, id := range *filters.ID {
				ids[i] = id
			}

			q.Apply(sm.Where(models.Countries.Columns.ID.In(psql.Arg(ids...))))
		}

		if filters.Name != nil {
			q.Apply(sm.Where(models.Countries.Columns.Name.ILike(psql.Arg("%" + *filters.Name + "%"))))
		}

		if filters.IsoAlpha2 != nil {
			q.Apply(sm.Where(models.Countries.Columns.IsoAlpha2.ILike(psql.Arg("%" + *filters.IsoAlpha2 + "%"))))
		}

		if filters.IsoAlpha3 != nil {
			q.Apply(sm.Where(models.Countries.Columns.IsoAlpha3.ILike(psql.Arg("%" + *filters.IsoAlpha3 + "%"))))
		}
	}

	type ListCountriesRow struct {
		models.Country
		TotalCount int64
	}

	countries, err := bob.All(ctx, pst.exec, q, scan.StructMapper[ListCountriesRow]())
	if err != nil {
		return dbtype.PagedResult[models.Country]{}, fmt.Errorf("query countries")
	}

	result := dbtype.PagedResult[models.Country]{
		Data: make([]models.Country, len(countries)),
	}
	for i, row := range countries {
		result.Data[i] = row.Country
	}

	if len(countries) > 0 {
		result.TotalCount = countries[0].TotalCount
	}

	return result, nil
}

func (pst *PgCountryPersistor) GetCountryByID(ctx context.Context, countryID int64) (models.Country, error) {
	q := psql.Select(
		sm.Columns(models.Countries.Columns),
		sm.From(models.Countries.Name()),
		sm.Where(models.Countries.Columns.ID.EQ(psql.Arg(countryID))),
	)

	country, err := bob.One(ctx, pst.exec, q, scan.StructMapper[models.Country]())
	if err != nil {
		return models.Country{}, fmt.Errorf("query country")
	}

	return country, nil
}

func (pst *PgCountryPersistor) DeleteCountryByID(ctx context.Context, countryID int64) (int64, error) {
	q := models.Countries.Delete(dm.Where(models.Countries.Columns.ID.EQ(psql.Arg(countryID))))
	if _, err := bob.Exec(ctx, pst.exec, q); err != nil {
		return 0, fmt.Errorf("delete country: %w", err)
	}

	return countryID, nil
}

func (pst *PgCountryPersistor) DeleteCountries(ctx context.Context, ids []int64) error {
	countryIDs := make([]any, len(ids))
	for i, id := range ids {
		countryIDs[i] = id
	}

	q := models.Countries.Delete(dm.Where(models.Countries.Columns.ID.In(psql.Arg(countryIDs...))))
	if _, err := bob.Exec(ctx, pst.exec, q); err != nil {
		return fmt.Errorf("delete countries: %w", err)
	}

	return nil
}

func (pst *PgCountryPersistor) CreateCountry(ctx context.Context, in models.CountrySetter) (models.Country, error) {
	q := models.Countries.Insert(&in, im.Returning(models.Countries.Columns))

	country, err := bob.One(ctx, pst.exec, q, scan.StructMapper[models.Country]())
	if err != nil {
		return models.Country{}, fmt.Errorf("insert country")
	}

	return country, nil
}

func (pst *PgCountryPersistor) UpdateCountry(ctx context.Context, countryID int64, in models.CountrySetter) (models.Country, error) {
	q := models.Countries.Update(
		in.UpdateMod(),
		um.Where(models.Countries.Columns.ID.EQ(psql.Arg(countryID))),
		um.Returning(models.Countries.Columns),
	)

	country, err := bob.One(ctx, pst.exec, q, scan.StructMapper[models.Country]())
	if err != nil {
		return models.Country{}, fmt.Errorf("update country")
	}

	return country, nil
}
