package server

import (
	"dongtzu/config"

	"github.com/gofiber/adaptor/v2"
	"github.com/gofiber/fiber/v2"
	"github.com/line/line-bot-sdk-go/linebot"
	"gitlab.geax.io/demeter/gologger/logger"
)

var server *fiber.App
var bot *linebot.Client

func Init() {
	defer logger.Debugf("[Server] Initialized.")
	loginLineBot()

	server = fiber.New()
	server.Post("/line/webhook", adaptor.HTTPHandlerFunc(lineBotCallbacks))

	dt := server.Group("/dt")
	dt.Post("/appointment", appointment())
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
