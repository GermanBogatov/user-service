package application

import (
	"github.com/gojuno/goose"
	_ "github.com/jackc/pgx/v4/stdlib"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
	"github.com/pkg/errors"
)

const (
	migrationsPath = "./internal/repository/postgres/migration"
)

func migrateUP(url string) error {
	conn, err := sqlx.Connect("pgx", url)
	if err != nil {
		return errors.Wrap(err, "failed to connect to postgres")
	}

	err = goose.Up(conn.DB, migrationsPath)
	if err != nil {
		return errors.Wrap(err, "failed goose up")
	}

	err = conn.Close()
	if err != nil {
		return errors.Wrap(err, "failed close migration client")
	}

	return nil
}
