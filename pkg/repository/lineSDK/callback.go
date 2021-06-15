package lineSDK

import (
	"net/http"

	"github.com/line/line-bot-sdk-go/linebot"
	"gitlab.geax.io/demeter/gologger/logger"
)

func WebhookCallback(w http.ResponseWriter, r *http.Request) {
	events, err := bot.ParseRequest(r)
	if err != nil {
		logger.Errorf("[LineSDK] parse request failed: %v", err)
	}

	for _, event := range events {
		logger.Debugf("[LineSDK] event payload: %v", event)
		switch event.Type {
		case linebot.EventTypeFollow:
			logger.Debugf("[LineSDK] Follower: %v", event.Source.UserID)

		case linebot.EventTypeUnfollow:
			logger.Debugf("[LineSDK] Unfollower: %v", event.Source.UserID)

		case linebot.EventTypeMessage:
			switch message := event.Message.(type) {
			case *linebot.TextMessage:
				logger.Debugf("[LineSDK] Receive message from %v, content is %v", event.Source.UserID, message.Text)
			}

		default:
			logger.Debugf("[LineSDK] something happend: %v", event.Type)
		}
	}
}
