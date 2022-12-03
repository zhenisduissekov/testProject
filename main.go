package main

import (
	"github.com/ansrivas/fiberprometheus/v2"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/swagger"
	"github.com/gofiber/websocket/v2"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	_ "github.com/zhenisduissekov/testProject/docs"
	"github.com/zhenisduissekov/testProject/pkg/config"
	"github.com/zhenisduissekov/testProject/pkg/connection"
	"github.com/zhenisduissekov/testProject/pkg/cryptocompare"
	"github.com/zhenisduissekov/testProject/pkg/handler"
	lg "github.com/zhenisduissekov/testProject/pkg/logger"
	"github.com/zhenisduissekov/testProject/pkg/repository"
	"github.com/zhenisduissekov/testProject/pkg/socket"
	"time"
)

var (
	listenAddress = ":3000"
	reqLogFormat  = "[${time}] ${status} - ${latency} ${method} ${path} ${ip} ${url} in ${bytesReceived} bytes/ out ${bytesSent} bytes\n"
	shutDownDelay = 2 * time.Second
)

// @title TEST service API
// @contact.name API Support
// @contact.email zduisekov@gmail.com
// @host localhost:3000
// @BasePath /
// @schemes http
func main() {
	cnf := config.New()
	zerolog.SetGlobalLevel(lg.New(cnf.LogLevel))
	conn := connection.New(cnf.DB)
	repo := repository.New(conn)
	crypto := cryptocompare.New(cnf.CryptoCompare, repo)
	socket := socket.New(repo)
	h := handler.New(crypto, socket)

	//go func() { //todo: uncomment this
	//	scheduler.Run(cnf.Scheduler, crypto)
	//}()

	if err := server(h, cnf.DB.Service).Listen(listenAddress); err != nil {
		log.Fatal().Err(err).Msg("Server down")
	}
}

func server(h *handler.Handler, serviceName string) *fiber.App {
	app := fiber.New()
	prometheus := fiberprometheus.New(serviceName)
	prometheus.RegisterAt(app, "/metrics")
	app.Use(prometheus.Middleware)
	app.Use(cors.New(cors.Config{
		AllowOrigins: "*",
		AllowMethods: "GET,POST",
		AllowHeaders: "*",
	}))

	app.Get("/health", h.HealthCheck)

	app.Use(logger.New(logger.Config{
		Format:       reqLogFormat,
		TimeFormat:   "15:04:05",
		TimeZone:     "Local",
		TimeInterval: 0,
		Output:       nil,
	}))

	app.Use("/ws", func(c *fiber.Ctx) error {
		if websocket.IsWebSocketUpgrade(c) {
			c.Locals("allowed", true)
			return c.Next()
		}
		return fiber.ErrUpgradeRequired
	})

	app.Get("/swagger/*", swagger.New(swagger.Config{
		URL:          "/swagger/doc.json",
		DeepLinking:  false,
		DocExpansion: "none",
		OAuth: &swagger.OAuthConfig{
			AppName:  "OAuth Provider",
			ClientId: "21bb4edc-05a7-4afc-86f1-2e151e4ba6e2",
		},
	}))

	api := app.Group("/service")
	{
		api.Get("/price", h.GetPrice)
		api.Get("/ws/price", websocket.New(h.Publish))
	}
	return app
}
