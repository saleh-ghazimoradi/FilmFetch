package cmd

import (
	"fmt"
	"github.com/saleh-ghazimoradi/FilmFetch/config"
	"github.com/saleh-ghazimoradi/FilmFetch/migrations"
	"github.com/saleh-ghazimoradi/FilmFetch/utils"
	"log/slog"
	"os"

	"github.com/spf13/cobra"
)

// migrateUpCmd represents the migrateUp command
var migrateUpCmd = &cobra.Command{
	Use:   "migrateUp",
	Short: "A brief description of your command",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("migrateUp called")

		logger := slog.New(slog.NewJSONHandler(os.Stdout, nil))

		cfg, err := config.NewConfig()
		if err != nil {
			logger.Error(err.Error())
			os.Exit(1)
		}

		postgresql := utils.NewPostgresql(
			utils.WithHost(cfg.Postgresql.Host),
			utils.WithPort(cfg.Postgresql.Port),
			utils.WithUser(cfg.Postgresql.User),
			utils.WithPassword(cfg.Postgresql.Password),
			utils.WithName(cfg.Postgresql.Name),
			utils.WithMaxOpenConn(cfg.Postgresql.MaxOpenConn),
			utils.WithMaxIdleConn(cfg.Postgresql.MaxIdleConn),
			utils.WithMaxIdleTime(cfg.Postgresql.MaxIdleTime),
			utils.WithSSLMode(cfg.Postgresql.SSLMode),
			utils.WithTimeout(cfg.Postgresql.Timeout),
		)

		db, err := postgresql.Connect()
		if err != nil {
			logger.Error(err.Error())
			os.Exit(1)
		}

		defer func() {
			if err := db.Close(); err != nil {
				logger.Error(err.Error())
				os.Exit(1)
			}
		}()

		migrator, err := migrations.NewMigrate(db, postgresql.Name)
		if err != nil {
			logger.Error(err.Error())
			os.Exit(1)
		}

		if err := migrator.UP(); err != nil {
			logger.Error(err.Error())
			os.Exit(1)
		}

		defer func() {
			if err := migrator.Close(); err != nil {
				logger.Error(err.Error())
				os.Exit(1)
			}
		}()
	},
}

func init() {
	rootCmd.AddCommand(migrateUpCmd)
}
