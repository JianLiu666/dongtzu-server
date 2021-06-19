package lineSDK

import (
	"github.com/line/line-bot-sdk-go/linebot"
	"gitlab.geax.io/demeter/gologger/logger"
)

func PushTextMessage(providerID, consumerLineID, msg string) {
	provider, ok := providerMapping.Load(providerID)
	if !ok {
		logger.Warnf("[LineSDK] can not found channel id by provider id: %v", providerID)
		return
	}

	c, ok := clientMap.Load(provider.LineAtChannelID)
	if !ok {
		logger.Warnf("[LineSDK] can not found line bot by channel id: %v", provider.LineAtChannelID)
		return
	}

	_, err := c.Bot.PushMessage(consumerLineID, linebot.NewTextMessage(msg)).Do()
	if err != nil {
		logger.Errorf("[LineSDK] push text message failed: %v", err)
	}
}
