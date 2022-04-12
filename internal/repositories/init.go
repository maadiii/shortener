package repositories

import (
	"context"
	"database/sql"

	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/postgres"
	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/mehdiazizii/fastcontroller"
	"github.com/pkg/errors"
)

type DbSession struct {
	*pgxpool.Pool
	config fastcontroller.SessionConfig
}

func NewSession(cfg fastcontroller.SessionConfig) (*DbSession, error) {
	db, err := pgxpool.Connect(context.Background(), cfg.DsnWithSchema())
	if err != nil {
		return nil, errors.Wrap(err, "")
	}

	if err := db.Ping(context.Background()); err != nil {
		return nil, errors.Wrap(err, "")
	}

	return &DbSession{Pool: db, config: cfg}, nil
}

func GetMigrator(db *sql.DB, cfg fastcontroller.SessionConfig) (*migrate.Migrate, error) {
	drv, err := postgres.WithInstance(db, &postgres.Config{SchemaName: "public", MigrationsTable: cfg.MigrationsTable})
	if err != nil {
		return nil, err
	}

	migrator, err := migrate.NewWithDatabaseInstance(cfg.MigrationsPath, cfg.Driver, drv)
	if err != nil {
		return nil, err
	}

	return migrator, nil
}

func CreateDB(cfg fastcontroller.SessionConfig) error {
	db, err := sql.Open(cfg.Driver, cfg.AdminDsn())
	if err != nil {
		return err
	}
	if err := db.Ping(); err != nil {
		return err
	}
	defer db.Close()

	if _, err := db.Exec("CREATE DATABASE " + cfg.DBName); err != nil {
		return err
	}

	return nil
}
