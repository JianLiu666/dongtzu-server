package arangodb

import (
	"context"
	"dongtzu/pkg/model"
	"fmt"

	"gitlab.geax.io/demeter/gologger/logger"
)

const (
	CollectionProviders = "Providers"
)

func GetProviderProfileByLineUserID(ctx context.Context, lineUserID string) (*model.Provider, error) {
	var result model.Provider

	query := fmt.Sprintf(`
		FOR p IN %s 
			FILTER p.lineUserId == @lineUserId
		RETURN p`, CollectionProviders)
	bindVars := map[string]interface{}{
		"lineUserId": lineUserID,
	}
	cursor, err := db.Query(ctx, query, bindVars)
	defer func() {
		if err := cursor.Close(); err != nil {
			logger.Errorf("[ArangoDB] cursor close falied: %v", err)
		}
	}()

	if err != nil {
		logger.Errorf("[ArangoDB] GetProviderProfileByLineUserID query falied: %v", err)
		return nil, err
	}

	_, err = cursor.ReadDocument(ctx, &result)
	if err != nil {
		logger.Errorf("[ArangoDB] GetProviderProfileByLineUserID cursor falied: %v", err)
		return nil, err
	}
	return &result, nil
}
