package lineSDK

import (
	"github.com/line/line-bot-sdk-go/linebot"
	"gitlab.geax.io/demeter/gologger/logger"
)

func PushMessage(userId, msg string) {
	_, err := bot.PushMessage(userId, linebot.NewTextMessage(msg)).Do()
	if err != nil {
		logger.Errorf("[LineSDK] push message failed: %v", err)
	}
}
