package cmd

import (
	"fmt"
	"github.com/saleh-ghazimoradi/FilmFetch/config"
	"github.com/saleh-ghazimoradi/FilmFetch/utils"
	"log"

	"github.com/spf13/cobra"
)

// httpCmd represents the http command
var httpCmd = &cobra.Command{
	Use:   "http",
	Short: "A brief description of your command",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("http called")

		cfg, err := config.NewConfig()
		if err != nil {
			log.Fatalf("Error loading config: %v", err)
		}

		posgresql := utils.NewPostgresql(
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
		postDB, err := posgresql.Connect()
		if err != nil {
			log.Fatalf("Error connecting to PostgreSQL: %v", err)
		}
		fmt.Println(postDB)
	},
}

func init() {
	rootCmd.AddCommand(httpCmd)
}
