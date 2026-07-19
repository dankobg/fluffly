package postgres

import (
	"context"
	"fmt"
	"net"
	"net/url"
	"strconv"

	"github.com/aarondl/opt/omitnull"
	"github.com/dankobg/fluffly/config"
	"github.com/dankobg/fluffly/persistence"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	bobpgx "github.com/stephenafamo/bob/drivers/pgx"
	pgxgeom "github.com/twpayne/pgx-geom"
)

// var _ persistence.Persistor = (*PgPersistor)(nil)

type PgPersistor struct {
	pool *pgxpool.Pool
	exec bobpgx.Pool
}

func NewPool(ctx context.Context, dbSettings config.DatabaseConfig) (*pgxpool.Pool, error) {
	dsn := url.URL{
		Scheme:   "postgresql",
		User:     url.UserPassword(dbSettings.User, dbSettings.Password),
		Host:     net.JoinHostPort(dbSettings.Host, strconv.Itoa(dbSettings.Port)),
		Path:     dbSettings.DB,
		RawQuery: url.Values{"sslmode": []string{dbSettings.SSLMode}}.Encode(),
	}

	poolCfg, err := pgxpool.ParseConfig(dsn.String())
	if err != nil {
		return nil, fmt.Errorf("pgxpool.ParseConfig: %w", err)
	}

	pool, err := pgxpool.NewWithConfig(ctx, poolCfg)
	if err != nil {
		return nil, fmt.Errorf("pgxpool.New: %w", err)
	}

	pool.Config().AfterConnect = func(ctx context.Context, c *pgx.Conn) error {
		if err := pgxgeom.Register(ctx, c); err != nil {
			return fmt.Errorf("pgxgeom.Register: %w", err)
		}

		return nil
	}

	return pool, nil
}

// New creates the new persistence
func New(pool *pgxpool.Pool) *PgPersistor {
	return &PgPersistor{pool: pool, exec: bobpgx.NewPool(pool)}
}

func (ps *PgPersistor) User() persistence.UserPersistor {
	return NewPgUserPersistor(ps)
}

func (ps *PgPersistor) Country() persistence.CountryPersistor {
	return NewPgCountryPersistor(ps)
}

func (ps *PgPersistor) Organization() persistence.OrganizationPersistor {
	return NewPgOrganizationPersistor(ps)
}

func (ps *PgPersistor) Animal() persistence.AnimalPersistor {
	return NewPgAnimalPersistor(ps)
}

func (ps *PgPersistor) Geocoding() persistence.GeocodingResultPersistor {
	return NewPgGeocodingResultPersistor(ps)
}

func (ps *PgPersistor) Analytics() persistence.AnalyticsPersistor {
	return NewPgAnalyticsPersistor(ps)
}

func (ps *PgPersistor) Adoption() persistence.AdoptionPersistor {
	return NewPgAdoptionPersistor(ps)
}

func (ps *PgPersistor) WithTx(
	ctx context.Context,
	fn func(tx bobpgx.Tx) error,
) (err error) {
	tx, err := ps.pool.BeginTx(ctx, pgx.TxOptions{})
	if err != nil {
		return err
	}

	bobTx := bobpgx.NewTx(tx, func() {})

	defer func() {
		if r := recover(); r != nil {
			_ = tx.Rollback(ctx)

			panic(r)
		} else if err != nil {
			_ = tx.Rollback(ctx)
		} else {
			err = tx.Commit(ctx)
		}
	}()

	return fn(bobTx)
}

func valOrNil[T any](v omitnull.Val[T]) *T {
	var out *T
	if !v.IsUnset() && !v.IsNull() {
		out = v.MustPtr()
	}

	return out
}
