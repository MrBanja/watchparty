package app

import (
	"github.com/gofiber/contrib/fiberzap"
	"github.com/gofiber/contrib/websocket"
	"github.com/gofiber/fiber/v2"
	"go.uber.org/zap"
	"path/filepath"
	"stream/internal/api/middlware"
	"stream/internal/api/service"
	"time"
)

func New() *fiber.App {
	app := fiber.New(fiber.Config{
		ReadTimeout:  time.Second * 10,
		WriteTimeout: time.Second * 10,
	})

	srv := service.New()

	app.Use(fiberzap.New(fiberzap.Config{Logger: zap.L()}))
	app.Use("/ws", middlware.UpgradeToWebsocket)

	app.Get("/", srv.Index)
	app.Get("/ws/:id", websocket.New(srv.Room))
	app.Get("/magnet", srv.GetMagnet)

	contentPath, err := filepath.Abs("./static/content")
	if err != nil {
		zap.S().Panicf("Error getting absolute path for content: %v", err)
	}
	app.Static("/static/content", contentPath, fiber.Static{
		Browse:   true,
		Download: true,
	})
	return app
}
