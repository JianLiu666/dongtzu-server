package arangodb

import (
	"context"
	"dongtzu/constant"
	"dongtzu/pkg/model"
	"time"

	"github.com/arangodb/go-driver"
	"gitlab.geax.io/demeter/gologger/logger"
)

// NOTE: 如果未來量過大時, 應該要分批撈取
func GetWithoutMeetingUrlSchedules(ctx context.Context) ([]*model.Schedule, int) {
	result := []*model.Schedule{}

	query := `
		FOR d IN Schedules
			FILTER d.startTimestamp <= @startTimestamp 
				AND d.count >= d.minConsumerLimit 
				AND d.meetingUrl == ""
			RETURN d`
	bindVars := map[string]interface{}{
		"startTimestamp": time.Now().Add(30 * time.Minute).UTC().Unix(),
	}
	cursor, err := db.Query(ctx, query, bindVars)
	defer closeCursor(cursor)
	if err != nil {
		logger.Errorf("[ArangoDB][GetSchedules] query failed: %v", err)
		return result, constant.ArangoDB_Driver_Failed
	}

	for {
		doc := &model.Schedule{}
		_, err := cursor.ReadDocument(ctx, doc)
		if driver.IsNoMoreDocuments(err) {
			break
		} else if err != nil {
			logger.Errorf("[ArangoDB][GetSchedules] cursor failed: %v", err)
			return result, constant.ArangoDB_Driver_Failed
		}

		result = append(result, doc)
	}

	return result, constant.ArangoDB_Success
}

func UpdateSchedulesMeetingUrl(ctx context.Context, docs []*model.Schedule) int {
	col, err := db.Collection(ctx, "Schedules")
	if err != nil {
		logger.Errorf("[ArangoDB][UpdateSchedulesMeetingUrl] get collection failed: %v", err)
		return constant.ArangoDB_Driver_Failed
	}

	keys := []string{}
	updates := []map[string]interface{}{}

	for _, doc := range docs {
		if doc.ID == "" {
			logger.Warnf("[ArangoDB][UpdateSchedulesMeetingUrl] the doc without key: %v", doc)
			continue
		}
		if doc.MeetingUrl == "" {
			logger.Warnf("[ArangoDB][UpdateSchedulesMeetingUrl] the doc without meetingUrl: %v", doc)
			continue
		}

		keys = append(keys, doc.ID)
		updates = append(updates, map[string]interface{}{
			"meetingUrl": doc.MeetingUrl,
		})
	}
	_, _, err = col.UpdateDocuments(ctx, keys, updates)
	if err != nil {
		logger.Errorf("[ArangoDB][UpdateSchedulesMeetingUrl] update documents failed: %v", err)
		return constant.ArangoDB_Driver_Failed
	}

	return constant.ArangoDB_Success
}
