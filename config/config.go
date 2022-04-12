package config

import (
	"os"
	"path/filepath"
	"runtime"

	"github.com/joho/godotenv"
	"github.com/mehdiazizii/fastcontroller"
)

func Config() fastcontroller.Config {
	cfg := fastcontroller.Config{
		HTTPPort: 8000,
	}

	mod := os.Getenv("SHORTENER_ENV")
	if mod != "RELEASE" {
		_, f, _, _ := runtime.Caller(0)
		if err := godotenv.Load(filepath.Dir(f) + "/envs"); err != nil {
			panic(err)
		}
		// logrus.Warning("Application run in DEVELOPER mode, set SHORTENER_ENV to RELEASE for production mod")
		cfg.DevMode = true
	}

	cfg.DbSession = fastcontroller.SessionConfig{
		Driver:          "postgres",
		Host:            os.Getenv("POSTGRES_HOST"),
		Port:            os.Getenv("POSTGRES_PORT"),
		User:            os.Getenv("POSTGRES_USER"),
		Password:        os.Getenv("POSTGRES_PASSWORD"),
		DBName:          os.Getenv("POSTGRES_DB"),
		Schema:          os.Getenv("POSTGRES_SCHEMA"),
		TestDBName:      "shortener_test",
		AdminDBName:     "postgres",
		SslMode:         "disable",
		TimeZone:        "Asia/Tehran",
		MigrationsPath:  "file://migrations",
		MigrationsTable: "shortener_migrations",
	}

	return cfg
}
