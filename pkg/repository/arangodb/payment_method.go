package arangodb

import (
	"context"
	"dongtzu/constant"
	"dongtzu/pkg/model"
	"fmt"

	"github.com/arangodb/go-driver"
	"gitlab.geax.io/demeter/gologger/logger"
)

func GetPaymentMethods(ctx context.Context) ([]*model.PaymentMethod, int) {
	results := []*model.PaymentMethod{}

	query := fmt.Sprintf(`
		FOR d IN %s
			RETURN d`, collectionPaymentMethods)
	cursor, err := db.Query(ctx, query, nil)
	defer closeCursor(cursor)
	if err != nil {
		logger.Errorf("[ArangoDB][GetPaymentMethods] failed to query: %v", err)
		return nil, constant.ArangoDB_Driver_Failed
	}

	for {
		var doc *model.PaymentMethod
		_, err := cursor.ReadDocument(ctx, &doc)
		if driver.IsNoMoreDocuments(err) {
			break
		} else if err != nil {
			logger.Errorf("[ArangoDB][GetPaymentMethods] failed to read doc: %v", err)
			return nil, constant.ArangoDB_Driver_Failed
		}

		results = append(results, doc)
	}

	return results, constant.ArangoDB_Success
}
