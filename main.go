package main

import (
	"context"
	"github.com/caarlos0/env/v9"
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
	o := app.Options{}
	if err := env.Parse(&o); err != nil {
		zap.S().Fatal(err)
	}
	server := app.New(o)
	serveGracefully(server, o.LocalAddr)
}

func serveGracefully(app *fiber.App, addr string) {
	ctx, cancel := context.WithCancel(context.Background())
	errors := make(chan error, 2)
	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, os.Interrupt, os.Kill)

	var wg sync.WaitGroup
	wg.Add(1)
	go func() {
		defer wg.Done()
		defer cancel()
		if err := app.Listen(addr); err != nil {
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
