package linebot

import (
	"net/http"

	"github.com/line/line-bot-sdk-go/linebot"
	"gitlab.geax.io/demeter/gologger/logger"
)

func lineBotCallbacks(w http.ResponseWriter, r *http.Request) {
	events, err := bot.ParseRequest(r)
	if err != nil {
		logger.Errorf("[LineBot] parse request failed: %v", err)
	}

	for _, event := range events {
		logger.Debugf("[LineBot] event payload: %v", event)
		switch event.Type {
		case linebot.EventTypeFollow:
			logger.Debugf("[LineBot] Follower: %v", event.Source.UserID)

		case linebot.EventTypeUnfollow:
			logger.Debugf("[LineBot] Unfollower: %v", event.Source.UserID)

		case linebot.EventTypeMessage:
			switch message := event.Message.(type) {
			case *linebot.TextMessage:
				logger.Debugf("[LineBot] Receive message from %v, content is %v", event.Source.UserID, message.Text)
			}

		default:
			logger.Debugf("[LineBot] something happend: %v", event.Type)
		}
	}
}
