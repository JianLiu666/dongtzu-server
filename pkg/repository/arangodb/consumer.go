package arangodb

import (
	"context"
	"dongtzu/constant"
	"dongtzu/pkg/model"
	"encoding/json"

	"github.com/arangodb/go-driver"
	"gitlab.geax.io/demeter/gologger/logger"
)

const (
	CollectionConsumers = "Consumers"
)

func CreateConsumer(ctx context.Context, doc *model.Consumer) int {
	jsonData, err := json.Marshal(doc)
	if err != nil {
		return constant.ArangoDB_Driver_Failed
	}

	aql := `
	function (Params) {
		const db = require('@arangodb').db;
		const consumerCol = db._collection("Consumers");
		const savedObj = JSON.parse(Params[0]);

		consumerDoc = consumerCol.firstExample({
			lineUserId: savedObj.lineUserId, 
			providerLineAtChannelId: savedObj.providerLineAtChannelId
		});

		if (consumerDoc && consumerDoc._key && consumerDoc._key.length > 0) {
			consumerCol.update(consumerDoc._key, savedObj);

		} else {
			consumerCol.insert(savedObj);
		}
	
		return 1;
	}`

	options := &driver.TransactionOptions{
		MaxTransactionSize: 100000,
		WriteCollections:   []string{CollectionConsumers},
		ReadCollections:    []string{CollectionConsumers},
		Params:             []interface{}{string(jsonData)},
		WaitForSync:        false,
	}

	_, err = db.Transaction(ctx, aql, options)
	if err != nil {
		logger.Errorf("[ArangoDB][CreateConsumer] transaction failed: %v", err)
		return constant.ArangoDB_Driver_Failed
	}

	return constant.ArangoDB_Success
}

func UpdateConsumerByLineUserId(ctx context.Context, userId string, updates map[string]interface{}) int {
	aql := `
	function (Params) {
		const db = require('@arangodb').db;
		const consumerCol = db._collection("Consumers");
		
		consumerDoc = consumerCol.firstExample({lineUserId: Params[0]});
		
		if (consumerDoc && consumerDoc._key && consumerDoc._key.length > 0) {
			consumerCol.update(consumerDoc._key, Params[1]);
		}
		
		return 1;
	}`

	options := &driver.TransactionOptions{
		MaxTransactionSize: 100000,
		WriteCollections:   []string{CollectionConsumers},
		ReadCollections:    []string{CollectionConsumers},
		Params:             []interface{}{userId, updates},
		WaitForSync:        false,
	}

	_, err := db.Transaction(ctx, aql, options)
	if err != nil {
		logger.Errorf("[ArangoDB][CreateConsumer] transaction failed: %v", err)
		return constant.ArangoDB_Driver_Failed
	}

	return constant.ArangoDB_Success
}
