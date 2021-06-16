package arangodb

import (
	"context"
	"dongtzu/constant"
	"dongtzu/pkg/model"

	"gitlab.geax.io/demeter/gologger/logger"
)

func CreateConsumer(ctx context.Context, doc *model.Consumer) int {
	col, err := db.Collection(ctx, "Consumers")
	if err != nil {
		logger.Errorf("[ArangoDB][CreateConsumer] get collection failed: %v", err)
		return constant.ArangoDB_Driver_Failed
	}

	_, err = col.CreateDocument(ctx, doc)
	if err != nil {
		logger.Errorf("[ArangoDB][CreateConsumer] create document failed: %v", err)
		return constant.ArangoDB_Driver_Failed
	}

	return constant.ArangoDB_Success
}

func UpdateConsumerByLineUserId(ctx context.Context, userId string, updates map[string]interface{}) int {
	col, err := db.Collection(ctx, "Consumers")
	if err != nil {
		logger.Errorf("[ArangoDB][UpdateConsumer] get collection failed: %v", err)
		return constant.ArangoDB_Driver_Failed
	}

	_, err = col.UpdateDocument(ctx, userId, updates)
	if err != nil {
		logger.Errorf("[ArangoDB][UpdateConsumer] create document failed: %v", err)
		return constant.ArangoDB_Driver_Failed
	}

	return constant.ArangoDB_Success
}
