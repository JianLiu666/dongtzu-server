package arangodb

import (
	"context"
	"dongtzu/constant"
	"dongtzu/pkg/model"
	"fmt"

	"github.com/arangodb/go-driver"
	"gitlab.geax.io/demeter/gologger/logger"
)

// TODO: 之後應該還要針對 expiredDuration 與 deleteAt 做判斷
func GetServiceProducts(ctx context.Context, providerID string) ([]*model.ServiceProduct, int) {
	results := []*model.ServiceProduct{}

	query := fmt.Sprintf(`
		FOR d IN %s
			FILTER d.providerId == @providerId
			RETURN d`, collectionServiceProducts)
	bindVars := map[string]interface{}{
		"providerId": providerID,
	}
	cursor, err := db.Query(ctx, query, bindVars)
	defer closeCursor(cursor)
	if err != nil {
		logger.Errorf("[ArangoDB][GetServiceProducts] failed to query: %v", err)
		return []*model.ServiceProduct{}, constant.ArangoDB_Driver_Failed
	}

	for {
		var doc *model.ServiceProduct
		_, err := cursor.ReadDocument(ctx, &doc)
		if driver.IsNoMoreDocuments(err) {
			break
		} else if err != nil {
			logger.Errorf("[ArangoDB][GetServiceProducts] failed to read doc: %v", err)
			return []*model.ServiceProduct{}, constant.ArangoDB_Driver_Failed
		}

		results = append(results, doc)
	}

	return results, constant.ArangoDB_Success
}
