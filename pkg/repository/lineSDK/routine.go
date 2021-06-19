package lineSDK

import (
	"context"
	"dongtzu/constant"
	"dongtzu/pkg/model"
	"dongtzu/pkg/repository/arangodb"

	"github.com/line/line-bot-sdk-go/linebot"
	"gitlab.geax.io/demeter/gologger/logger"
)

func startEventHandler() {
	logger.Debugf("[LineSDK] Start event handler.")

	for req := range reqChan {
		events, code := req.parseEvents()
		if code != constant.LineSDK_Success {
			logger.Warnf("[LineSDK] Prase line events failed")
			continue
		}

		c, ok := clientMap.Load(req.ChannelID)
		if !ok {
			logger.Warnf("[LineSDK] can not found line bot by channel id: %v", req.ChannelID)
			continue
		}

		provider, ok := providerMapping.Load(c.ProviderID)
		if !ok {
			logger.Warnf("[LineSDK] can not found provider by provider id: %v", c.ProviderID)
			continue
		}

		handleEvents(provider, events)
	}
}

func handleEvents(provider *model.Provider, events []*linebot.Event) {
	for _, event := range events {
		logger.Debugf("Receive new event:\n%v", event)
		switch event.Type {
		case linebot.EventTypeFollow:
			createConsumer(provider, event.Source.UserID)

		case linebot.EventTypeUnfollow:
			updateConsumer(event.Source.UserID)

		case linebot.EventTypeMessage:

		default:
		}
	}
}

func createConsumer(provider *model.Provider, userID string) {
	doc := &model.Consumer{
		LineUserID:              userID,
		ProviderID:              provider.ID,
		ProviderLineAtChannelID: provider.LineAtChannelID,
		LineFollowingStatus:     constant.Consumer_LineStatus_Following,
	}

	_ = arangodb.CreateConsumer(context.TODO(), doc)
}

func updateConsumer(userId string) {
	updates := map[string]interface{}{
		"lineFollowingStatus": constant.Consumer_LineStatus_Unfollowing,
	}

	_ = arangodb.UpdateConsumerByLineUserId(context.TODO(), userId, updates)
}
