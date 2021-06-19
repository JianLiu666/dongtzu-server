package lineSDK

import (
	"dongtzu/constant"

	"github.com/line/line-bot-sdk-go/linebot"
	"gitlab.geax.io/demeter/gologger/logger"
)

func startEventHandler() {
	logger.Debugf("[LineSDK] Start event handler.")

	for req := range reqChan {
		events, code := req.parseEvents()
		if code != constant.LineSDK_Success {
			continue
		}

		handleEvents(events)
	}
}

func handleEvents(events []*linebot.Event) {
	for _, event := range events {
		logger.Debugf("Receive new event:\n%v", event)
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
	// doc := &model.Consumer{
	// 	LineUserID:          userId,
	// 	LineFollowingStatus: constant.Consumer_LineStatus_Following,
	// }

	// _ = arangodb.CreateConsumer(context.TODO(), doc)
}

func updateConsumer(userId string) {
	// updates := map[string]interface{}{
	// 	"lineFollowingStatus": constant.Consumer_LineStatus_Unfollowing,
	// }

	// _ = arangodb.UpdateConsumerByLineUserId(context.TODO(), userId, updates)
}
