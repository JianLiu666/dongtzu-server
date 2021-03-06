package arangodb

import (
	"context"
	"dongtzu/constant"
	"dongtzu/pkg/model"
	"fmt"

	"github.com/arangodb/go-driver"
	"gitlab.geax.io/demeter/gologger/logger"
)

func GetZoomAccounts(ctx context.Context) ([]*model.ZoomAccount, int) {
	results := []*model.ZoomAccount{}

	query := fmt.Sprintf(`
		FOR d IN %s
			FILTER d.userId != "" AND
				d.apiKet != "" AND
				d.apiSecret != ""
			RETURN d`, collectionZoomAccounts)
	bindVars := map[string]interface{}{}
	cursor, err := db.Query(ctx, query, bindVars)
	defer closeCursor(cursor)
	if err != nil {
		logger.Errorf("[ArangoDB][GetZoomAccounts] failed to query: %v", err)
		return []*model.ZoomAccount{}, constant.ArangoDB_Driver_Failed
	}

	for {
		var doc *model.ZoomAccount
		_, err := cursor.ReadDocument(ctx, &doc)
		if driver.IsNoMoreDocuments(err) {
			break
		} else if err != nil {
			logger.Errorf("[ArangoDB][GetZoomAccounts] failed to read doc: %v", err)
			return []*model.ZoomAccount{}, constant.ArangoDB_Driver_Failed
		}

		results = append(results, doc)
	}

	return results, constant.ArangoDB_Success
}
