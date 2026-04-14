package main

import (
	"context"
	"fmt"
	"os"

	"go.uber.org/fx"

	"lucasaladino.com/semanticlog/logger"
)

const file = "/tmp/go-semantic-log-test.log"

func main() {
	app := fx.New(
		fx.NopLogger,
		fx.Provide(newFileLogger),
		fx.Provide(newConsoleLogger),
		fx.Invoke(cleanup),
		fx.Invoke(start),
		fx.Invoke(check),
		fx.Invoke(cleanup),
	)

	ctx := context.Background()
	err := app.Start(ctx)
	if err != nil {
		os.Exit(1)
	}

	err = app.Stop(ctx)
	if err != nil {
		os.Exit(1)
	}
}

func cleanup(lc fx.Lifecycle, l *logger.ConsoleLogger) {
	lc.Append(fx.Hook{
		OnStart: func(ctx context.Context) error {
			if file == "" {
				return nil
			}

			if _, err := os.Stat(file); os.IsNotExist(err) {
				return nil
			}

			err := os.Remove(file)
			if err != nil {
				l.Error(fmt.Sprintf("Failed to remove log file: %v\n", err))
			}

			return nil
		},
	})
}

func newFileLogger(lc fx.Lifecycle) *logger.FileLogger {
	l := logger.NewFileLogger(file, "example", logger.LevelDebug)

	lc.Append(fx.Hook{
		OnStart: func(ctx context.Context) error {
			return l.Init()
		},
		OnStop: func(ctx context.Context) error {
			return l.Close()
		},
	})

	return l
}

func newConsoleLogger(lc fx.Lifecycle) *logger.ConsoleLogger {
	l := logger.NewConsoleLogger(logger.LevelDebug)

	lc.Append(fx.Hook{
		OnStart: func(ctx context.Context) error {
			return l.Init()
		},
		OnStop: func(ctx context.Context) error {
			return l.Close()
		},
	})

	return l
}

func start(lc fx.Lifecycle, l *logger.FileLogger) {
	lc.Append(fx.Hook{
		OnStart: func(ctx context.Context) error {
			l.Info("App starting...")

			l.Debug("This is a debug message")
			l.Info("This is an info message")
			l.Warn("This is a warning message")
			l.Error("This is an error message")

			l.SetLevel(logger.LevelWarn)

			l.Debug("This should NOT appear (level set to Warn)")
			l.Info("This should NOT appear (level set to Warn)")
			l.Warn("This IS a warning message")
			l.Error("This IS an error message")

			return nil
		},
	})
}

func check(lc fx.Lifecycle, l *logger.ConsoleLogger) {
	lc.Append(fx.Hook{
		OnStart: func(ctx context.Context) error {
			l.Debug("Reading file saved by FileLogger")
			content, err := os.ReadFile(file)
			if err != nil {
				l.Error("Failed to read file")
				return err
			}

			l.Info(fmt.Sprintf("File content length: %d", len(content)))

			l.Info("Check function executed successfully - logger is working!")
			l.Warn("This is a warning from check function")
			l.Error("This is an error from check function")

			l.Debug("Check the below list:")
			l.List([]string{
				"First line",
				"Second line",
				"Third line",
			})

			return nil
		},
	})
}
