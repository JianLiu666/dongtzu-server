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

func PushButtonTemplateMessage(providerID, consumerLineID string) {
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

	template := linebot.NewButtonsTemplate(
		"https://www.charlestonchronicle.net/wp-content/uploads/2020/01/kobe-bryant-memorial-v1-nba-1.jpg",
		"ButtonsTemplate Test",
		"description.",
		linebot.NewURIAction("開始上課", "https://www.google.com/"),
		linebot.NewPostbackAction("結束簽退", "data", "text", ""),
	)

	_, err := c.Bot.PushMessage(consumerLineID, linebot.NewTemplateMessage("template alt text", template)).Do()
	if err != nil {
		logger.Errorf("[LineSDK] push buttons template message failed: %v", err)
	}
}
