package app

import (
	"github.com/gofiber/contrib/fiberzap"
	"github.com/gofiber/contrib/websocket"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/template/html/v2"
	"go.uber.org/zap"
	"stream/internal/api/middlware"
	"stream/internal/api/service"
	"time"
)

type Options struct {
	PublicAddr string `env:"PUBLIC_ADDR,required"`
	LocalAddr  string `env:"LOCAL_ADDR" envDefault:":8000"`
}

func New(o Options) *fiber.App {
	engine := html.New("./static", ".html")

	app := fiber.New(fiber.Config{
		ReadTimeout:  time.Second * 10,
		WriteTimeout: time.Second * 10,
		Views:        engine,
	})

	srv := service.New(o.PublicAddr)

	app.Use(cors.New())
	app.Use(fiberzap.New(fiberzap.Config{Logger: zap.L()}))
	app.Use("/ws", middlware.UpgradeToWebsocket)

	app.Get("/ws/:id", websocket.New(srv.Room))
	app.Get("/magnet", srv.GetMagnet)
	app.Get("/:id/peer_status/:partID", srv.GetStatusBadge)

	return app
}
