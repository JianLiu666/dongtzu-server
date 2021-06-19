package lineSDK

import (
	"github.com/line/line-bot-sdk-go/linebot"
	"gitlab.geax.io/demeter/gologger/logger"
)

func SendMeetingUrl(providerID, consumerLineID, meetingUrl string) {
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

	template, err := linebot.UnmarshalFlexMessageJSON([]byte(getMeetingFlexTemplate(meetingUrl)))
	if err != nil {
		logger.Errorf("[LineSDK] unmarshal JSON template failed: %v", err)
		return
	}

	_, err = c.Bot.PushMessage(consumerLineID, linebot.NewFlexMessage("template alt text", template)).Do()
	if err != nil {
		logger.Errorf("[LineSDK] push buttons template message failed: %v", err)
	}
}
