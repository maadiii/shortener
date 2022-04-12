package repositories

import (
	"context"
	"database/sql"
	"os"
	"path/filepath"
	"runtime"
	"testing"

	"shortener/config"

	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/mehdiazizii/fastcontroller"
	"github.com/pkg/errors"
)

func TestMain(m *testing.M) {
	cfg := config.Config()
	cfg.DbSession.DBName = cfg.DbSession.TestDBName
	createTestDB(cfg.DbSession)
	migrateDb(cfg.DbSession)

	os.Exit(m.Run())
}

func createTestDB(cfg fastcontroller.SessionConfig) {
	db, err := sql.Open(cfg.Driver, cfg.AdminDsn())
	if err != nil {
		panic(err)
	}
	defer db.Close()

	if err := db.Ping(); err != nil {
		panic(err)
	}

	_, _ = db.Exec("DROP DATABASE " + cfg.TestDBName)

	if _, err := db.Exec("CREATE DATABASE " + cfg.TestDBName); err != nil {
		panic(err)
	}
}

func newTestSession(cfg fastcontroller.SessionConfig) (*DbSession, error) {
	db, err := pgxpool.Connect(context.Background(), cfg.DsnWithSchema())
	if err != nil {
		return nil, errors.Wrap(err, "")
	}

	if err := db.Ping(context.Background()); err != nil {
		return nil, errors.Wrap(err, "")
	}

	return &DbSession{Pool: db}, nil
}

func migrateDb(cfg fastcontroller.SessionConfig) {
	db, err := sql.Open(cfg.Driver, cfg.Dsn())
	if err != nil {
		panic(err)
	}
	defer db.Close()

	if err := db.Ping(); err != nil {
		panic(err)
	}

	_, f, _, _ := runtime.Caller(1)
	cfg.MigrationsPath = "file://" + filepath.Dir(filepath.Dir(filepath.Dir(f))) + "/migrations"
	m, err := GetMigrator(db, cfg)
	if err != nil {
		panic(err)
	}

	if err := m.Up(); err != nil {
		panic(err)
	}
}
