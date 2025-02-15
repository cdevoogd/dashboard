package main

import (
	"flag"
	"log/slog"
	"os"

	"github.com/cdevoogd/dashboard/internal/dashboard"
)

var configFilePath string

func createLogger(level slog.Leveler) *slog.Logger {
	handler := slog.NewTextHandler(os.Stderr, &slog.HandlerOptions{
		AddSource: false,
		Level:     level,
	})
	return slog.New(handler)
}

func main() {
	logger := createLogger(slog.LevelInfo)

	flag.StringVar(&configFilePath, "config", "config.yaml", "Path to the config file")
	flag.Parse()

	logger.Info("Loading config", "path", configFilePath)
	config, err := dashboard.LoadConfig(configFilePath)
	if err != nil {
		logger.Error("Error loading config", "err", err)
		os.Exit(1)
	}

	logger.Info("Setting the current log level to the configured level", "level", config.LogLevel)
	logger = createLogger(config.LogLevel.Level())

	server, err := dashboard.NewServer(config, logger)
	if err != nil {
		logger.Error("Error constructing a new server", "err", err)
		os.Exit(1)
	}

	logger.Info("Starting to listen for HTTP requests", "port", config.Port)
	err = server.ListenAndServe()
	if err != nil {
		logger.Error("Error listening for requests", "err", err)
		os.Exit(1)
	}
}
