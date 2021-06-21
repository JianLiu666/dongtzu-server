package arangodb

import (
	"context"
	"dongtzu/constant"
	"dongtzu/pkg/model"
	"encoding/json"
	"time"

	"github.com/arangodb/go-driver"
	"gitlab.geax.io/demeter/gologger/logger"
)

func CreateOrder(ctx context.Context, consumerId, providerId, serviceProductId string, amount int) (*model.Order, int) {
	doc := &model.Order{
		ConsumerID:       consumerId,
		ProviderID:       providerId,
		ServiceProductID: serviceProductId,
		PaymentMethodID:  "",
		Amount:           1,
		Status:           0,
		CreatedAt:        time.Now().Unix(),
		UpdatedAt:        time.Now().Unix(),
		PaidAt:           0,
	}
	jsonData, err := json.Marshal(doc)
	if err != nil {
		logger.Errorf("[ArangoDB][CreateOrder] failed to marshal: %v", err)
		return nil, constant.ArangoDB_Driver_Failed
	}

	aql := `
	function (Params) {
		const db = require('@arangodb').db;
		const savedObj = JSON.parse(Params[0]);
		const consumerCol = db._collection("Consumers");
		const providerCol = db._collection("Providers");
		const productCol = db._collection("ServiceProducts");
		const orderCol = db._collection("Orders");

		consumerDoc = consumerCol.firstExample({_key: savedObj.consumerId});
		if (!consumerDoc || !consumerDoc._key || consumerDoc._key.length == 0) {
			throw("consumer not found.");
		}

		providerDoc = providerCol.firstExample({_key: savedObj.providerId});
		if (!providerDoc || !providerDoc._key || providerDoc._key.length == 0) {
			throw("provider not found.");
		}

		productrDoc = productCol.firstExample({_key: savedObj.serviceProductId});
		if (!productrDoc || !productrDoc._key || productrDoc._key.length == 0) {
			throw("service product not found.");
		}

		var result = orderCol.insert(savedObj, {returnNew: true});
		return JSON.stringify(result);
	}`

	options := &driver.TransactionOptions{
		MaxTransactionSize: 100000,
		WriteCollections:   []string{collectionOrders},
		ReadCollections:    []string{collectionConsumers, collectionProviders, collectionServiceProducts},
		Params:             []interface{}{string(jsonData)},
		WaitForSync:        false,
	}

	resp, err := db.Transaction(ctx, aql, options)
	if err != nil {
		logger.Errorf("[ArangoDB][CreateOrder] transaction failed: %v", err)
		return nil, constant.ArangoDB_Driver_Failed
	}

	result := &model.Order{}
	if err = json.Unmarshal([]byte(resp.(string)), result); err != nil {
		logger.Errorf("[ArangoDB][CreateOrder] failed to unmarshal response data: %v", err)
		return nil, constant.ArangoDB_Driver_Failed
	}

	return result, constant.ArangoDB_Success
}
