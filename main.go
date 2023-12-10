package main

import (
	"flag"

	"github.com/cdevoogd/dashboard/internal/dashboard"
	"github.com/charmbracelet/log"
)

var configFilePath string

func main() {
	logger := log.Default()

	flag.StringVar(&configFilePath, "config", "config.yaml", "Path to the config file")
	flag.Parse()

	logger.Info("Loading config", "path", configFilePath)
	config, err := dashboard.LoadConfig(configFilePath)
	if err != nil {
		logger.Fatal("Error loading config", "err", err)
	}

	level, err := log.ParseLevel(config.LogLevel)
	if err != nil {
		logger.Fatal("Error parsing log level", "configured_level", config.LogLevel, "err", err)
	}

	logger.Info("Setting the current log level to the configured level", "level", level)
	log.SetLevel(level)

	server, err := dashboard.NewServer(config, logger.WithPrefix("http-server"))
	if err != nil {
		logger.Fatal("Error constructing a new server", "err", err)
	}

	logger.Info("Starting to listen for HTTP requests", "port", config.Port)
	err = server.ListenAndServe()
	if err != nil {
		logger.Fatal("Error listening for requests", "err", err)
	}
}
