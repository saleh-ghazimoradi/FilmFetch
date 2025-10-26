package cmd

import (
	"fmt"
	"github.com/saleh-ghazimoradi/FilmFetch/config"
	"github.com/saleh-ghazimoradi/FilmFetch/internal/gateway/handlers"
	"github.com/saleh-ghazimoradi/FilmFetch/internal/gateway/routes"
	"github.com/saleh-ghazimoradi/FilmFetch/internal/helper"
	"github.com/saleh-ghazimoradi/FilmFetch/internal/middleware"
	"github.com/saleh-ghazimoradi/FilmFetch/internal/repository"
	"github.com/saleh-ghazimoradi/FilmFetch/internal/server"
	"github.com/saleh-ghazimoradi/FilmFetch/internal/service"
	"github.com/saleh-ghazimoradi/FilmFetch/utils"
	"github.com/saleh-ghazimoradi/FilmFetch/utils/email"
	"github.com/wneessen/go-mail"
	"log/slog"
	"os"
	"time"

	"github.com/spf13/cobra"
)

// httpCmd represents the http command
var httpCmd = &cobra.Command{
	Use:   "http",
	Short: "A brief description of your command",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("http called")

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

		clientMail, err := mail.NewClient(
			cfg.Mail.Host,
			mail.WithSMTPAuth(mail.SMTPAuthLogin),
			mail.WithPort(cfg.Mail.Port),
			mail.WithUsername(cfg.Mail.Username),
			mail.WithPassword(cfg.Mail.Password),
			mail.WithTimeout(5*time.Second),
		)
		if err != nil {
			logger.Error(err.Error())
			os.Exit(1)
		}

		mailer := email.NewMailSender(
			email.WithClient(clientMail),
			email.WithHost(cfg.Mail.Host),
			email.WithPort(cfg.Mail.Port),
			email.WithUsername(cfg.Mail.Username),
			email.WithPassword(cfg.Mail.Password),
			email.WithSender(cfg.Mail.Sender),
		)

		customError := helper.NewCustomErr(logger)
		middleWare := middleware.NewMiddleware(cfg, customError)
		healthHandler := handlers.NewHealthHandler(cfg, logger, customError)
		healthRoutes := routes.NewHealthRoute(healthHandler)
		movieRepository := repository.NewMovieRepository(db, db)
		movieService := service.NewMovieService(movieRepository)
		movieHandler := handlers.NewMovieHandler(logger, customError, movieService)
		movieRoutes := routes.NewMovieRoutes(movieHandler)
		userRepository := repository.NewUserRepository(db, db)
		userService := service.NewUserService(userRepository)
		userHandler := handlers.NewUserHandler(customError, userService, mailer, logger)
		userRoutes := routes.NewUserRoutes(userHandler)

		registerRoutes := routes.NewRegister(
			routes.WithCustomError(customError),
			routes.WithMiddleware(middleWare),
			routes.WithHealthRoutes(healthRoutes),
			routes.WithMovieRoutes(movieRoutes),
			routes.WithUserRoutes(userRoutes),
		)

		httpServer := server.NewServer(
			server.WithHost(cfg.Server.Host),
			server.WithPort(cfg.Server.Port),
			server.WithHandler(registerRoutes.RegisterRoutes()),
			server.WithIdleTimeout(cfg.Server.IdleTimeout),
			server.WithReadTimeout(cfg.Server.ReadTimeout),
			server.WithWriteTimeout(cfg.Server.WriteTimeout),
			server.WithErrorLog(slog.NewLogLogger(logger.Handler(), slog.LevelError)),
		)

		logger.Info("starting server", "addr", cfg.Server.Host+":"+cfg.Server.Port, "env", cfg.Application.Environment)

		if err := httpServer.Connect(); err != nil {
			logger.Error(err.Error())
			os.Exit(1)
		}
	},
}

func init() {
	rootCmd.AddCommand(httpCmd)
}
