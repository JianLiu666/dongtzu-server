package lineapi

import (
	"dongtzu/config"

	"github.com/line/line-bot-sdk-go/linebot"
	"gitlab.geax.io/demeter/gologger/logger"
)

var bot *linebot.Client

func Init() {
	b, err := linebot.New(
		config.GetGlobalConfig().LineBot.ChannelSecret,
		config.GetGlobalConfig().LineBot.ChannelAccessToken)
	if err != nil {
		logger.Errorf("[LineAPI] Init line bot failed: %v", err)
		return
	}

	bot = b
	logger.Debugf("[LineAPI] Initialized.")
}
