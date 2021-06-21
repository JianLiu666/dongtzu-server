package arangodb

import (
	"context"
	"dongtzu/constant"
	"dongtzu/pkg/model"
	"encoding/json"

	"github.com/arangodb/go-driver"
	"gitlab.geax.io/demeter/gologger/logger"
)

func CreatePayment(ctx context.Context, doc *model.Payment) int {
	jsonData, err := json.Marshal(doc)
	if err != nil {
		logger.Errorf("[ArangoDB][CreatePayment] failed to marshal: %v", err)
		return constant.ArangoDB_Driver_Failed
	}

	aql := `
	function (Params) {
		const db = require('@arangodb').db;
		const savedObj = JSON.parse(Params[0]);
		const orderCol = db._collection("Orders");
		const consumerCol = db._collection("Consumers");
		const methodCol = db._collection("PaymentMethods");
		const paymentCol = db._collection("Payments");

		orderDoc = orderCol.firstExample({_key: savedObj.orderId});
		if (!orderDoc || !orderDoc._key || orderDoc._key.length == 0) {
			return -1;
		}
		if (orderDoc.status != 0) {
			return -2;
		}

		consumerDoc = consumerCol.firstExample({_key: savedObj.consumerId});
		if (!consumerDoc || !consumerDoc._key || consumerDoc._key.length == 0) {
			return -3;
		}

		methodDoc = methodCol.firstExample({_key: savedObj.paymentMethodId});
		if (!methodDoc || !methodDoc._key || methodDoc._key.length == 0) {
			return -4;
		}

		var orderUpdates = {
			paymentMethodId: savedObj.paymentMethodId,
			status: 2,
			updatedAt: savedObj.createdAt,
			paidAt: savedObj.createdAt,
		}
		orderCol.update(orderDoc._key, orderUpdates);

		paymentCol.insert(savedObj);

		return 1;
	}`

	options := &driver.TransactionOptions{
		MaxTransactionSize: 100000,
		WriteCollections:   []string{collectionOrders, collectionPayments},
		ReadCollections:    []string{collectionOrders, collectionConsumers, collectionPaymentMethods},
		Params:             []interface{}{string(jsonData)},
		WaitForSync:        false,
	}

	code, err := db.Transaction(ctx, aql, options)
	if err != nil {
		logger.Errorf("[ArangoDB][CreatePayment] transaction failed: %v", err)
		return constant.ArangoDB_Driver_Failed
	}
	if code.(float64) != 1 {
		logger.Errorf("[ArangoDB][CreatePayment] invalid transaction operation: %v", code)
		return constant.ArangoDB_Invalid_Operation
	}

	return constant.ArangoDB_Success
}
