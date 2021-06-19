package server

import (
	"dongtzu/config"
	"fmt"
	"runtime"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/recover"
	"gitlab.geax.io/demeter/gologger/logger"
)

var server *fiber.App

func Init() {
	defer logger.Debugf("[Server] Initialized.")

	server = fiber.New()
	setMiddelWare(server)
	server.Post("/webhook/:channelId", lineWebhook())

	dt := server.Group("/dt")
	dt.Post("/appointment", appointment())

	// Provider Registration
	dt.Get("/providers/:lineUserId", getProviderInfo())
	dt.Post("/providers/register", registerProvider())
	dt.Put("/providers/:lineUserId", updateProviderInfo())

	// Provider Dashboard
	dt.Get("/providers/:lineUserId/eventSchedule", getProviderEventSchedule())
	dt.Get("/providers/:lineUserId/incomeSummary", getProviderIncomeSummary())
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

/**
 * Private Method
 */
func setMiddelWare(fiberInstance *fiber.App) {
	fiberInstance.Use(recover.New(
		recover.Config{
			EnableStackTrace: true,
			StackTraceHandler: func(e interface{}) {
				buf := make([]byte, 1024)
				buf = buf[:runtime.Stack(buf, false)]
				logger.Errorf(fmt.Sprintf("catch panic error: %v\n%s\n", e, buf))
			},
		},
	))
	fiberInstance.Use(cors.New())
}
