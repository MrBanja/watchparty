package main

import (
	"context"
	"github.com/gofiber/fiber/v2"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"os"
	"os/signal"
	app "stream/app"
	"sync"
	"time"
)

func init() {
	cfg := zapcore.EncoderConfig{
		MessageKey:     "message",
		LevelKey:       "level",
		TimeKey:        "timestamp",
		NameKey:        "name",
		CallerKey:      "caller",
		FunctionKey:    "function",
		StacktraceKey:  "stacktrace",
		EncodeLevel:    zapcore.CapitalLevelEncoder,
		EncodeTime:     zapcore.ISO8601TimeEncoder,
		EncodeDuration: zapcore.SecondsDurationEncoder,
		EncodeCaller:   zapcore.ShortCallerEncoder,
		EncodeName:     zapcore.FullNameEncoder,
	}
	core := zapcore.NewCore(
		zapcore.NewConsoleEncoder(cfg),
		zapcore.AddSync(os.Stdout),
		zap.DebugLevel,
	)
	logger := zap.New(
		core,
		zap.AddCaller(),
		zap.AddStacktrace(zapcore.ErrorLevel),
	)
	zap.ReplaceGlobals(logger)
}

func main() {
	server := app.New()
	serveGracefully(server)
}

func serveGracefully(app *fiber.App) {
	ctx, cancel := context.WithCancel(context.Background())
	errors := make(chan error, 2)
	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, os.Interrupt, os.Kill)

	var wg sync.WaitGroup
	wg.Add(1)
	go func() {
		defer wg.Done()
		defer cancel()
		if err := app.Listen("0.0.0.0:8000"); err != nil {
			errors <- err
		}
	}()

	wg.Add(1)
	go func() {
		defer wg.Done()
		defer cancel()

		select {
		case <-ctx.Done():
			zap.S().Error("Server error")
		case <-sigs:
			zap.S().Warn("Received signal, shutting down gracefully")
		}
		zap.S().Info("Shutting down gracefully")

		ctxShutdown, cancelShutdown := context.WithTimeout(context.Background(), time.Second*10)
		defer cancelShutdown()
		if err := app.Server().ShutdownWithContext(ctxShutdown); err != nil {
			errors <- err
		}
	}()
	wg.Wait()
	close(errors)

	for err := range errors {
		zap.S().Error(err)
	}
}
