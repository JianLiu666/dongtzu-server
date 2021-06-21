package arangodb

import (
	"context"
	"dongtzu/pkg/model"

	"github.com/arangodb/go-driver"
	"gitlab.geax.io/demeter/gologger/logger"
)

func Migration() {
	ctx := context.TODO()
	ensureCollection(ctx, collectionAppointments)
	ensureCollection(ctx, collectionCourses)
	ensureCollection(ctx, collectionConsumers)
	ensureCollection(ctx, collectionFeedbacks)
	ensureCollection(ctx, collectionMonthReceipts)
	ensureCollection(ctx, collectionOrders)
	ensureCollection(ctx, collectionPayments)
	ensureCollection(ctx, collectionPaymentMethods)
	ensureCollection(ctx, collectionProviders)
	ensureCollection(ctx, collectionSchedules)
	ensureCollection(ctx, collectionServiceProducts)
	ensureCollection(ctx, collectionZoomAccounts)

	seedCollectionProviders(ctx)
	seedCollectionZoomAccounts(ctx)
}

// helper to check if a collection exists and create it if needed.
func ensureCollection(ctx context.Context, name string) {
	_, err := db.Collection(ctx, name)
	if driver.IsNotFound(err) {
		_, err = db.CreateCollection(ctx, name, nil)
		if err != nil {
			logger.Errorf("[ArangoDB][ensureCollection] failed to create collection '%v': %v", name, err)
		}
	} else if err != nil {
		logger.Errorf("[ArangoDB][ensureCollection] failed to open collection '%v': %v", name, err)
	}

	logger.Debugf("[ArangoDB][ensureCollection] %v ok.", name)
}

func seedCollectionProviders(ctx context.Context) {
	c, err := db.Collection(ctx, collectionProviders)
	if err != nil {
		logger.Errorf("[ArangoDB][seedCollectionProviders] failed to open collection: %v", err)
		return
	}

	dataSet := []*model.Provider{
		{
			LineAtID:            "DongTzu@v0.0.1",
			LineAtChannelID:     "1656097222",
			LineAtChannelSecret: "ed2091c7a39700df2cc70c15336ad3e1",
			LineAtAccessToken:   "gYmJDfEHSnNRd9U+ZVSghhBLd3LgKmpeptcvA2X8XOjCRK68EM/uBBSNYYcwidfudfZ3twgqO4TuHjFiNh7mKXaxm37HGAvw22Adc9s7JZkR5S9qPEFwDfWtCC4oK305RfxjlUR9Y2rHf8ioNMA2oAdB04t89/1O/w1cDnyilFU=",
		},
		{
			LineAtID:            "DongTzu@v0.0.2",
			LineAtChannelID:     "1656120045",
			LineAtChannelSecret: "a7d96a4ba33935ed2396a54b46034955",
			LineAtAccessToken:   "xVcvMO7NZHmMrJcY8ZYVV+nJKIrgReK47mffUV9+GymDDp7/urgHaSMTacrE/HrXEkoRJb2BtgqNkGtYZK+76mGFTc6XubI2Y+MDhR538xsRCG8ueMEmnYpTRNq/vTed6domTrVGx0G02vzBXAWXqAdB04t89/1O/w1cDnyilFU=",
		},
	}

	for _, data := range dataSet {
		_, err = c.CreateDocument(ctx, data)
		if err != nil {
			logger.Errorf("[ArangoDB][seedCollectionProviders] failed to create document: %v", err)
		}
	}

	logger.Debugf("[ArangoDB][seedCollectionProviders] done.")
}

func seedCollectionZoomAccounts(ctx context.Context) {
	c, err := db.Collection(ctx, collectionZoomAccounts)
	if err != nil {
		logger.Errorf("[ArangoDB][seedCollectionZoomAccounts] failed to open collection: %v", err)
		return
	}

	dataSet := []*model.ZoomAccount{
		{
			UserID:    "jianliou.6@gmail.com",
			APIKey:    "XIZoX2EVRTeZ6WGyoS36Og",
			APISecret: "hPwXgPnXkZk0TPDy5GjCkdWtglCPlGg9IRrV",
		},
		{
			UserID:    "jianliu0616@gmail.com",
			APIKey:    "ugiN0wadQFqHc5tlv1bYGg",
			APISecret: "Y0fXMLEDFa3KUbkrd7aPY19YN2OoyDGGwvWv",
		},
	}

	for _, data := range dataSet {
		_, err = c.CreateDocument(ctx, data)
		if err != nil {
			logger.Errorf("[ArangoDB][seedCollectionZoomAccounts] failed to create document: %v", err)
		}
	}

	logger.Debugf("[ArangoDB][seedCollectionZoomAccounts] done.")
}
