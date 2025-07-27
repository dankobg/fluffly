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
	"github.com/stephenafamo/bob"
)

// var _ persistence.Persistor = (*PgPersistor)(nil)

type PgPersistor struct {
	db bob.DB
}

// Connect creates the new connection
func Connect(dbSettings config.DatabaseConfig) (bob.DB, error) {
	dsn := url.URL{
		Scheme:   "postgresql",
		User:     url.UserPassword(dbSettings.User, dbSettings.Password),
		Host:     net.JoinHostPort(dbSettings.Host, strconv.Itoa(dbSettings.Port)),
		Path:     dbSettings.DB,
		RawQuery: url.Values{"sslmode": []string{dbSettings.SSLMode}}.Encode(),
	}

	var bobdb bob.DB
	var err error

	bobdb, err = bob.Open("pgx", dsn.String())
	for err != nil {
		log.Printf("failed to connect to database")

		if dbSettings.RetriesNum > 0 {
			dbSettings.RetriesNum--
			log.Printf("retrying the database connection. retries left: (%d)", dbSettings.RetriesNum)
			time.Sleep(dbSettings.RetriesDelay)
			bobdb, err = bob.Open("pgx", dsn.String())
			continue
		}

		return bob.DB{}, fmt.Errorf("failed to connect to the database: %w", err)
	}

	return bobdb, nil
}

// New creates the new persistence
func New(db bob.DB) *PgPersistor {
	return &PgPersistor{db: db}
}

func (ps *PgPersistor) User() persistence.UserPersistor {
	return NewPgUserPersistor(ps)
}

func WithTx(ctx context.Context, db *sql.DB, fn func(tx *sql.Tx) error) (err error) {
	tx, err := db.BeginTx(ctx, nil)
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
