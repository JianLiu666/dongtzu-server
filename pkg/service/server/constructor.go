package server

import (
	"dongtzu/config"
	"dongtzu/pkg/repository/lineSDK"

	"github.com/gofiber/adaptor/v2"
	"github.com/gofiber/fiber/v2"
	"gitlab.geax.io/demeter/gologger/logger"
)

var server *fiber.App

func Init() {
	defer logger.Debugf("[Server] Initialized.")

	server = fiber.New()
	server.Post("/line/webhook", adaptor.HTTPHandlerFunc(lineSDK.WebhookCallback))

	dt := server.Group("/dt")
	dt.Post("/appointment", appointment())

	// Provider Registration
	dt.Get("/providers/:lineUserId", GetProviderInfo())
	dt.Post("/providers/register", RegisterProvider())
	dt.Put("/providers/:lineUserId", UpdateProviderInfo())
}

func Start() {
	go func() {
		if err := server.Listen(config.GetGlobalConfig().Fiber.Port); err != nil {
			logger.Errorf("[Server] enable fiber server failed: %v", err)
			return
		}
		logger.Debugf("[Server] Enabled.")
	}()
}