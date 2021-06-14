package server

import (
	"dongtzu/config"
	"net/http"

	"github.com/line/line-bot-sdk-go/linebot"
	"gitlab.geax.io/demeter/gologger/logger"
)

func loginLineBot() {
	b, err := linebot.New(
		config.GetGlobalConfig().LineBot.ChannelSecret,
		config.GetGlobalConfig().LineBot.ChannelAccessToken)
	if err != nil {
		logger.Errorf("[Server] Init line bot failed: %v", err)
		return
	}

	bot = b
}

func lineBotCallbacks(w http.ResponseWriter, r *http.Request) {
	events, err := bot.ParseRequest(r)
	if err != nil {
		logger.Errorf("[Server] parse request failed: %v", err)
	}

	for _, event := range events {
		logger.Debugf("[Server] event payload: %v", event)
		switch event.Type {
		case linebot.EventTypeFollow:
			logger.Debugf("[Server] Follower: %v", event.Source.UserID)

		case linebot.EventTypeUnfollow:
			logger.Debugf("[Server] Unfollower: %v", event.Source.UserID)

		case linebot.EventTypeMessage:
			switch message := event.Message.(type) {
			case *linebot.TextMessage:
				logger.Debugf("[Server] Receive message from %v, content is %v", event.Source.UserID, message.Text)
			}

		default:
			logger.Debugf("[Server] something happend: %v", event.Type)
		}
	}
}
