package postgres

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"net"
	"net/url"
	"strconv"
	"time"

	"github.com/dankobg/fluffly/config"
	"github.com/dankobg/fluffly/persistence"
	_ "github.com/jackc/pgx/v5/stdlib"
)

// var _ persistence.Persistor = (*PgPersistor)(nil)

type PgPersistor struct {
	db *sql.DB
}

// Connect creates the new connection
func Connect(dbSettings config.DatabaseConfig) (*sql.DB, error) {
	dsn := url.URL{
		Scheme:   "postgresql",
		User:     url.UserPassword(dbSettings.User, dbSettings.Password),
		Host:     net.JoinHostPort(dbSettings.Host, strconv.Itoa(dbSettings.Port)),
		Path:     dbSettings.DB,
		RawQuery: url.Values{"sslmode": []string{dbSettings.SSLMode}}.Encode(),
	}

	var db *sql.DB
	var err error

	db, err = sql.Open("pgx", dsn.String())
	for err != nil {
		log.Printf("failed to connect to database")

		if dbSettings.RetriesNum > 0 {
			dbSettings.RetriesNum--
			log.Printf("retrying the database connection. retries left: (%d)", dbSettings.RetriesNum)
			time.Sleep(dbSettings.RetriesDelay)
			db, err = sql.Open("pgx", dsn.String())
			continue
		}

		return nil, fmt.Errorf("failed to connect to the database: %w", err)
	}

	return db, nil
}

// New creates the new persistence
func New(db *sql.DB) *PgPersistor {
	return &PgPersistor{db: db}
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

func (ps *PgPersistor) WithTx(ctx context.Context, fn func(tx *sql.Tx) error) (err error) {
	tx, err := ps.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}

	defer func() {
		if r := recover(); r != nil {
			_ = tx.Rollback()
			panic(r)
		} else if err != nil {
			_ = tx.Rollback()
		} else {
			err = tx.Commit()
		}
	}()

	return fn(tx)
}
