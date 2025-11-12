package data

import (
	"errors"
	"fmt"
	"io/fs"

	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	"github.com/golang-migrate/migrate/v4/source/iofs"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
	"github.com/okunix/prservice/internal/app/config"
)

var (
	postgresDB *sqlx.DB
)

func InitPostgres(postgresConfig config.PostgresConfig, migrations fs.FS) error {
	addr := fmt.Sprintf(
		"postgres://%s:%s@%s/%s?sslmode=%s",
		postgresConfig.User,
		postgresConfig.Password,
		postgresConfig.Addr,
		postgresConfig.DB,
		postgresConfig.SSLMode,
	)
	db, err := sqlx.Open("postgres", addr)
	if err != nil {
		return err
	}
	if err := db.Ping(); err != nil {
		return err
	}

	postgresDB = db

	return migratePostgres(addr, migrations)
}

func migratePostgres(addr string, migrations fs.FS) error {
	sourceDriver, err := iofs.New(migrations, "migrations")
	if err != nil {
		return err
	}
	dbDriver, err := database.Open(addr)
	if err != nil {
		return err
	}
	migrator, err := migrate.NewWithInstance("iofs", sourceDriver, "postgres", dbDriver)
	if err != nil {
		return err
	}
	defer migrator.Close()
	if err := migrator.Up(); err != nil && !errors.Is(err, migrate.ErrNoChange) {
		return err
	}
	return nil
}

func PostgreSQL() *sqlx.DB {
	return postgresDB
}
