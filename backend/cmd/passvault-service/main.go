package main

import (
	passvaultService "github.com/Novando/pintartek/internal/passvault-service"
	"github.com/Novando/pintartek/pkg/env"
	"github.com/Novando/pintartek/pkg/logger"
	"github.com/Novando/pintartek/pkg/postgresql/pgx/v5"
	"github.com/gofiber/fiber/v2"
	"github.com/spf13/viper"
	"os"
	"os/signal"
	"syscall"
)

func main() {
	// Channel to receive OS signals
	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, os.Interrupt, syscall.SIGTERM)

	// Logger configuration
	log := logger.InitZerolog(logger.Config{
		ConsoleLoggingEnabled: true,
		FileLoggingEnabled:    true,
		CallerSkip:            3,
		Directory:             "./log",
		Filename:              "logfile",
	})

	if os.Getenv("TZ") == "" {
		log.Info("TZ set to default (Asia/Jakarta)")
		os.Setenv("TZ", "Asia/Jakarta")
	} else {
		log.Infof("TZ = %v", os.Getenv("TZ"))
	}

	// Environment configuration
	if err := env.InitViper("./config/config.local.json", log); err != nil || os.Getenv("CONSUL_PATH") != "" {
		// If not on local configuration, use Consul and Sentry
		consul := env.InitConsul(
			"52.230.98.3",
			8800,
			"http",
			log,
		)

		consul.RetrieveConfiguration(
			os.Getenv("CONSUL_PATH"),
			"json",
		)

		log.Infof("Using Consul at %s", os.Getenv("CONSUL_PATH"))

	} else {
		log.Info("Using local")
	}

	// Init DB
	pgxpool, query, err := pgx.InitPGXv5(
		viper.GetString("postgres.username"),
		viper.GetString("postgres.password"),
		viper.GetString("postgres.host"),
		viper.GetInt("postgres.port"),
		viper.GetString("postgres.database"),
		viper.GetString("postgres.schema"),
		5,
	)
	if err != nil {
		log.Panic(err.Error())
	}
	defer pgxpool.Close()

	// Fiber configuration
	app := fiber.New()

	// Define a health check endpoint
	app.Get("/health", func(c *fiber.Ctx) error {
		return c.SendStatus(fiber.StatusOK)
	})

	// Module initialization
	v1 := app.Group("/v1")
	passvaultService.InitPassvaultService(v1, query, pgxpool, log)

	// Start Fiber
	go func() {
		err := app.Listen(":" + viper.GetString("application.port"))
		if err != nil {
			log.Fatalf("%s: %s", "Error starting server", err)
		}

		log.Infof("Server started on port %s", viper.GetString("application.port"))
	}()
	log.Infof("Service started")

	// Wait for termination signal
	<-sigCh
	log.Info("Received termination signal. Initiating shutdown...")
}
