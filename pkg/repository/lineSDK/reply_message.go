package lineSDK

import (
	"github.com/line/line-bot-sdk-go/linebot"
	"gitlab.geax.io/demeter/gologger/logger"
)

func replyFlexMessageExample(channelID, replyToken string) {
	c, ok := clientMap.Load(channelID)
	if !ok {
		logger.Warnf("[LineSDK] can not found line bot by channel id: %v", channelID)
		return
	}

	// TODO: mock data
	template, err := linebot.UnmarshalFlexMessageJSON([]byte(getMeetingFlexTemplate("https://www.google.com/")))
	if err != nil {
		logger.Errorf("[LineSDK] unmarshal JSON template failed: %v", err)
		return
	}

	_, err = c.Bot.ReplyMessage(replyToken, linebot.NewFlexMessage("template alt text", template)).Do()
	if err != nil {
		logger.Errorf("[LineSDK] reply buttons template message failed: %v", err)
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
		logger.Errorf("[LineSDK] unmarshal JSON template failed: %v", err)
		return
	}

	_, err = c.Bot.ReplyMessage(replyToken, linebot.NewFlexMessage("template alt text", template)).Do()
	if err != nil {
		logger.Errorf("[LineSDK] reply buttons template message failed: %v", err)
	}
}
