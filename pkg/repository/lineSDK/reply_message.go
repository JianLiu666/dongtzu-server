package lineSDK

import (
	"github.com/line/line-bot-sdk-go/linebot"
	"gitlab.geax.io/demeter/gologger/logger"
)

func replyButtonTemplateMessageExample(channelID, replyToken string) {
	c, ok := clientMap.Load(channelID)
	if !ok {
		logger.Warnf("[LineSDK] can not found line bot by channel id: %v", channelID)
		return
	}

	template := linebot.NewButtonsTemplate(
		"https://www.charlestonchronicle.net/wp-content/uploads/2020/01/kobe-bryant-memorial-v1-nba-1.jpg",
		"ButtonsTemplate Test",
		"description.",
		linebot.NewURIAction("開始上課", "https://www.google.com/"),
		linebot.NewPostbackAction("結束簽退", "GET_FEEDBACK_URL", "", ""),
	)

	_, err := c.Bot.ReplyMessage(replyToken, linebot.NewTemplateMessage("template alt text", template)).Do()
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

	template := linebot.NewImageCarouselTemplate(
		linebot.NewImageCarouselColumn(
			"https://memegenerator.net/img/instances/60468871.jpg",
			linebot.NewURIAction("", "https://www.google.com/"),
		),
	)

	_, err := c.Bot.ReplyMessage(replyToken, linebot.NewTemplateMessage("template alt text", template)).Do()
	if err != nil {
		logger.Errorf("[LineSDK] reply buttons template message failed: %v", err)
	}
}
