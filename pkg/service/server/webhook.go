package server

import (
	"context"
	"dongtzu/constant"
	"dongtzu/pkg/model"
	"dongtzu/pkg/repository/arangodb"
	"dongtzu/pkg/repository/lineSDK"
	"net/http"

	"github.com/line/line-bot-sdk-go/linebot"
)

func lineWebhook(w http.ResponseWriter, r *http.Request) {
	events := lineSDK.ParseRequest(w, r)

	for _, event := range events {
		switch event.Type {
		case linebot.EventTypeFollow:
			createConsumer(event.Source.UserID)

		case linebot.EventTypeUnfollow:
			updateConsumer(event.Source.UserID)

		case linebot.EventTypeMessage:

		default:
		}
	}
}

func createConsumer(userId string) {
	doc := &model.Consumer{
		LineUserID:          userId,
		LineFollowingStatus: constant.Consumer_LineStatus_Following,
	}

	_ = arangodb.CreateConsumer(context.TODO(), doc)
}

func updateConsumer(userId string) {
	// updates := map[string]interface{}{
	// 	"lineFollowingStatus": constant.Consumer_LineStatus_Unfollowing,
	// }

	// _ = arangodb.UpdateConsumerByLineUserId(context.TODO(), userId, updates)
}
