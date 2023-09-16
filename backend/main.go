package main

import (
	"fmt"
	"log/slog"
	"os"

	"github.com/cdevoogd/dashboard/backend/internal/api"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

func setupLogger() *zap.Logger {
	config := zap.Config{
		Level:             zap.NewAtomicLevelAt(zap.InfoLevel),
		DisableCaller:     true,
		DisableStacktrace: true,
		Encoding:          "console",
		EncoderConfig: zapcore.EncoderConfig{
			TimeKey:          "T",
			LevelKey:         "L",
			NameKey:          "N",
			CallerKey:        zapcore.OmitKey,
			FunctionKey:      zapcore.OmitKey,
			MessageKey:       "M",
			StacktraceKey:    zapcore.OmitKey,
			LineEnding:       zapcore.DefaultLineEnding,
			EncodeLevel:      zapcore.CapitalLevelEncoder,
			EncodeTime:       zapcore.RFC3339TimeEncoder,
			EncodeDuration:   zapcore.StringDurationEncoder,
			EncodeCaller:     zapcore.ShortCallerEncoder,
			ConsoleSeparator: "  ",
		},
		OutputPaths:      []string{"stderr"},
		ErrorOutputPaths: []string{"stderr"},
	}
	return zap.Must(config.Build())
}

func run() error {
	logger := setupLogger()
	sugaredLogger := logger.Sugar()

	config, err := api.LoadServerConfig()
	if err != nil {
		return fmt.Errorf("error loading server config: %w", err)
	}

	server, err := api.NewServer(config, sugaredLogger)
	if err != nil {
		return fmt.Errorf("error creating server: %w", err)
	}

	return server.Serve()
}

func main() {
	slog.SetDefault(slog.New(slog.NewTextHandler(os.Stderr, nil)))

	err := run()
	if err != nil {
		slog.Error(err.Error())
		os.Exit(1)
	}
}
