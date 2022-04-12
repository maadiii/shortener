package main

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"os"
	"os/signal"
	"strings"

	"shortener/config"
	"shortener/internal/repositories"

	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/source/file"

	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

var rootCli = &cobra.Command{
	Use:   "shortener",
	Short: "URL shortener web application.",
	PersistentPreRun: func(cli *cobra.Command, args []string) {
		logrus.SetFormatter(&logrus.TextFormatter{
			ForceColors:   true,
			FullTimestamp: true,
		})
	},
}

var serveCli = &cobra.Command{
	Use:   "serve",
	Short: "Serve the application",
	RunE: func(cli *cobra.Command, args []string) error {
		ctx, cancel := context.WithCancel(context.Background())
		go func() {
			defer cancel()
			ch := make(chan os.Signal, 1)
			signal.Notify(ch, os.Interrupt)
			<-ch
			fmt.Println()
			logrus.Info("Signal caught. Shutting down...")
		}()

		serve(ctx)

		return nil
	},
}

var dbCli = &cobra.Command{
	Use:   "db",
	Short: "Database management.",
}

var migrateCli = &cobra.Command{
	Use:   "migrate",
	Short: "Migrate database.",
	RunE: func(cli *cobra.Command, args []string) error {
		version, err := cli.Flags().GetUint("version")
		if err != nil {
			return err
		}

		cfg := config.Config()
		db, err := sql.Open(cfg.DbSession.Driver, cfg.DbSession.Dsn())
		if err != nil {
			return err
		}
		defer db.Close()
		if err := db.Ping(); err != nil {
			if strings.Contains(err.Error(), `database "shortener" does not exist`) {
				if err := repositories.CreateDB(cfg.DbSession); err != nil {
					return err
				}
			} else {
				return err
			}
		}

		m, err := repositories.GetMigrator(db, cfg.DbSession)
		if err != nil {
			return err
		}

		if err := m.Migrate(version); err != nil && !errors.Is(err, migrate.ErrNoChange) {
			return err
		}

		return nil
	},
}
