package lineSDK

import (
	"context"
	"dongtzu/constant"
	"dongtzu/pkg/model"
	"dongtzu/pkg/repository/arangodb"
	"time"

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
			setConsumerToUnfollowing(event.Source.UserID)

		case linebot.EventTypeMessage:
			handleMessage(provider, event)

		case linebot.EventTypePostback:
			handlePostback(provider, event)

		default:
		}
	}
}

func handleMessage(provider *model.Provider, event *linebot.Event) {
	switch message := event.Message.(type) {
	case *linebot.TextMessage:
		switch message.Text {
		case "測試購買":
			createConsumer(provider, event.Source.UserID)
			buyDefaultProductExample(provider, event.Source.UserID, event.ReplyToken)

		case "測試預約":
			createConsumer(provider, event.Source.UserID)

		case "測試流程":
			replyFlexMessageExample(provider.LineAtChannelID, event.ReplyToken)
		}
	}
}

func handlePostback(provider *model.Provider, event *linebot.Event) {
	data := event.Postback.Data
	switch data {
	case "GET_FEEDBACK_URL":
		replyFeedbackUrl(provider.LineAtChannelID, event.ReplyToken)

	default:
		logger.Warnf("[LineSDK] postback data has unknown syntax: %v", data)
	}
}

func createConsumer(provider *model.Provider, userID string) {
	doc := &model.Consumer{
		LineUserID:              userID,
		ProviderID:              provider.ID,
		ProviderLineAtChannelID: provider.LineAtChannelID,
		LineFollowingStatus:     constant.Consumer_LineStatus_Following,
		CreatedAt:               int(time.Now().Unix()),
	}

	_ = arangodb.CreateConsumer(context.TODO(), doc)
}

func setConsumerToUnfollowing(userID string) {
	updates := map[string]interface{}{
		"lineFollowingStatus": constant.Consumer_LineStatus_Unfollowing,
	}

	_ = arangodb.UpdateConsumerByLineUserId(context.TODO(), userID, updates)
}

func buyDefaultProductExample(provider *model.Provider, userID, replyToken string) {
	products, code := arangodb.GetServiceProducts(context.TODO(), provider.ID)
	if code != constant.ArangoDB_Success {
		return
	}
	if len(products) == 0 {
		return
	}

	consumer, code := arangodb.GetConsumerByLineUserID(context.TODO(), userID)
	if code != constant.ArangoDB_Success {
		return
	}

	order, code := arangodb.CreateOrder(context.TODO(), consumer.ID, provider.ID, products[0].ID, 1)
	if code != constant.ArangoDB_Success {
		return
	}

	paymentMethods, code := arangodb.GetPaymentMethods(context.TODO())
	if code != constant.ArangoDB_Success {
		return
	}
	if len(paymentMethods) == 0 {
		return
	}

	payment := &model.Payment{
		OrderID:         order.ID,
		ConsumerID:      consumer.ID,
		PaymentMethodID: paymentMethods[0].ID,
		PaidPrice:       int64(products[0].Price),
		PlatformFee:     0,
		PaymentFee:      0,
		AgentFee:        0,
		AdFee:           0,
		TaxFee:          0,
		NetAmount:       0,
		Status:          0,
		RawParams:       "",
		CreatedAt:       time.Now().Unix(),
		UpdatedAt:       time.Now().Unix(),
	}

	code = arangodb.CreatePayment(context.TODO(), payment)
	if code != constant.ArangoDB_Success {
		return
	}

	replyTextMessage(provider.LineAtChannelID, replyToken, "購買成功")
}
