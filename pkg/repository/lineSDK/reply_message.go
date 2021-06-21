package lineSDK

import (
	"github.com/line/line-bot-sdk-go/linebot"
	"gitlab.geax.io/demeter/gologger/logger"
)

func replyTextMessage(channelID, replyToken string, msg string) {
	c, ok := clientMap.Load(channelID)
	if !ok {
		logger.Warnf("[LineSDK] can not found line bot by channel id: %v", channelID)
		return
	}

	_, err := c.Bot.ReplyMessage(replyToken, linebot.NewTextMessage(msg)).Do()
	if err != nil {
		logger.Errorf("[LineSDK] failed to reply text message: %v", err)
	}
}

func replyFlexMessageExample(channelID, replyToken string) {
	c, ok := clientMap.Load(channelID)
	if !ok {
		logger.Warnf("[LineSDK] can not found line bot by channel id: %v", channelID)
		return
	}

	// TODO: mock data
	template, err := linebot.UnmarshalFlexMessageJSON([]byte(getMeetingFlexTemplate("https://www.google.com/")))
	if err != nil {
		logger.Errorf("[LineSDK] failed to unmarshal JSON template: %v", err)
		return
	}

	_, err = c.Bot.ReplyMessage(replyToken, linebot.NewFlexMessage("template alt text", template)).Do()
	if err != nil {
		logger.Errorf("[LineSDK] failed to reply buttons template message: %v", err)
	}
}

func replyFeedbackUrl(channelID, replyToken string) {
	c, ok := clientMap.Load(channelID)
	if !ok {
		logger.Warnf("[LineSDK] can not found line bot by channel id: %v", channelID)
		return
	}

	// TODO: mock data
	template, err := linebot.UnmarshalFlexMessageJSON([]byte(getFeedbackFlexTemplate("https://www.google.com/")))
	if err != nil {
		logger.Errorf("[LineSDK] failed to unmarshal JSON template: %v", err)
		return
	}

	_, err = c.Bot.ReplyMessage(replyToken, linebot.NewFlexMessage("template alt text", template)).Do()
	if err != nil {
		logger.Errorf("[LineSDK] failed to reply buttons template message: %v", err)
	}
}
