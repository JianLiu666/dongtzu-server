package arangodb

import (
	"context"
	"dongtzu/constant"
	"dongtzu/pkg/model"
	"encoding/json"

	"github.com/arangodb/go-driver"
	"gitlab.geax.io/demeter/gologger/logger"
)

// 建立預約
// TODO: 還沒有檢查 consumer 的剩餘堂數
//
// @param ctx
//
// @param appt
//
// @return int status code
func CreateAppointment(ctx context.Context, appt *model.Appointment) int {
	jsonData, err := json.Marshal(appt)
	if err != nil {
		logger.Errorf("[ArangoDB][CreateAppointment] failed to marshal: %s", err)
	}

	aql := `
	function (Params) {
		const db = require('@arangodb').db;
		const savedObj = JSON.parse(Params[0]);
		const scheduleCol = db._collection("Schedules");
		const apptCol = db._collection("Appointments");

		// 檢查 schedule 預約人數是否已經達到上限
		scheduleDoc = scheduleCol.firstExample({_key: savedObj.scheduleId});
		if (!scheduleDoc || !scheduleDoc._key || scheduleDoc._key.length == 0) {
			return -1;
		}
		if (scheduleDoc.count+1 > scheduleDoc.maxConsumerLimit) {
			return -2;
		}

		// 檢查 consumer 在相同時段中是否已經有其他的 appointment 存在
		apptDocs = apptCol.closedRange("courseStartAt", savedObj.courseStartAt, savedObj.courseEndAt).toArray();
		var found = false;
		for (var i = 0; i < apptDocs.length; i++) {
			if (apptDocs[i].consumerId == savedObj.consumerId) {
				found = true;
				break;
			}
		}
		if (found) {
			return -3;
		}

		// 建立 appointment 並更新 schedule 統計人數
		scheduleCol.update(scheduleDoc._key, {count: scheduleDoc.count+1});
		apptCol.insert(savedObj);

		return 1;
	}`

	options := &driver.TransactionOptions{
		MaxTransactionSize: 100000,
		WriteCollections:   []string{collectionSchedules, collectionAppointments},
		ReadCollections:    []string{collectionSchedules, collectionAppointments},
		Params:             []interface{}{string(jsonData)},
		WaitForSync:        false,
	}

	code, err := db.Transaction(ctx, aql, options)
	if err != nil {
		logger.Errorf("[ArangoDB][CreateAppointment] transaction failed: %v", err)
		return constant.ArangoDB_Driver_Failed
	}
	if code.(float64) != 1 {
		logger.Errorf("[ArangoDB][CreateAppointment] invalid transaction operation: %v", code)
		return constant.ArangoDB_Invalid_Operation
	}

	return constant.ArangoDB_Success
}

// 取得與條件相符的預約
//
// @param scheduleId
//
// @param status
//
// @return []*model.Appointment
//
// @return int status code
func GetApptsByScheduleIDAndStatus(ctx context.Context, scheduleId string, status int) ([]*model.Appointment, int) {
	result := []*model.Appointment{}

	query := `
		FOR d IN Appointments
			FILTER d.scheduleId == @scheduleId 
				AND d.status == @status
			RETURN d`
	bindVars := map[string]interface{}{
		"scheduleId": scheduleId,
		"status":     status,
	}
	cursor, err := db.Query(ctx, query, bindVars)
	defer closeCursor(cursor)
	if err != nil {
		logger.Errorf("[ArangoDB][GetApptsByScheduleIDAndStatus] query failed: %v", err)
		return result, constant.ArangoDB_Driver_Failed
	}

	for {
		doc := &model.Appointment{}
		_, err := cursor.ReadDocument(ctx, doc)
		if driver.IsNoMoreDocuments(err) {
			break
		} else if err != nil {
			logger.Errorf("[ArangoDB][GetApptsByScheduleIDAndStatus] cursor failed: %v", err)
			return result, constant.ArangoDB_Driver_Failed
		}

		result = append(result, doc)
	}

	return result, constant.ArangoDB_Success
}

// 批次更新預約
//
// @param docs
//
// @param checkStatus 檢查準備更新的狀態是否相符
//
// @return int status code
func UpdateApptsStatus(ctx context.Context, docs []*model.Appointment, checkStatus int) int {
	col, err := db.Collection(ctx, "Appointments")
	if err != nil {
		logger.Errorf("[ArangoDB][UpdateApptsStatus] get collection failed: %v", err)
		return constant.ArangoDB_Driver_Failed
	}

	keys := []string{}
	updates := []map[string]interface{}{}

	for _, doc := range docs {
		if doc.ID == "" {
			logger.Warnf("[ArangoDB][UpdateApptsStatus] the doc without key: %v", doc)
			continue
		}
		if doc.Status != checkStatus {
			logger.Warnf("[ArangoDB][UpdateApptsStatus] the doc status is not equal to %v : %v", checkStatus, doc)
			continue
		}

		keys = append(keys, doc.ID)
		updates = append(updates, map[string]interface{}{
			"status": doc.Status,
		})
	}
	_, _, err = col.UpdateDocuments(ctx, keys, updates)
	if err != nil {
		logger.Errorf("[ArangoDB][UpdateApptsStatus] update documents failed: %v", err)
		return constant.ArangoDB_Driver_Failed
	}

	return constant.ArangoDB_Success
}
