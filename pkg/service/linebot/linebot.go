package linebot

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
	defer logger.Debugf("[LineBot] Initialized.")
	loginLineBot()

	server = fiber.New()
	server.Post("/", adaptor.HTTPHandlerFunc(lineBotCallbacks))
}

func Start() {
	go func() {
		if err := server.Listen(config.GetGlobalConfig().LineBot.Port); err != nil {
			logger.Errorf("[LineBot] enable fiber server failed: %v", err)
			return
		}
		logger.Debugf("[LineBot] Enabled.")
	}()
}

func loginLineBot() {
	b, err := linebot.New(
		config.GetGlobalConfig().LineBot.ChannelSecret,
		config.GetGlobalConfig().LineBot.ChannelAccessToken)
	if err != nil {
		logger.Errorf("[LineBot] Init line bot failed: %v", err)
		return
	}

	bot = b
}
