package arangodb

import (
	"context"
	"dongtzu/constant"
	"dongtzu/pkg/model"

	"github.com/arangodb/go-driver"
	"gitlab.geax.io/demeter/gologger/logger"
)

// 建立預約
//
// 1. 檢查 schedule 預約人數是否已經達到上限
//
// 2. 檢查 consumer 在相同時段中是否已經有其他的 appointment 存在
//
// 3. 建立 appointment 並更新 schedule 統計人數
//
// @param ctx
//
// @param appt
//
// @return int status code
func CreateAppointment(ctx context.Context, appt *model.Appointment) int {
	trxId, err := db.BeginTransaction(
		ctx,
		driver.TransactionCollections{
			Read:  []string{"Schedules"},
			Write: []string{"Schedules", "Appointments"},
		},
		nil,
	)
	if err != nil {
		logger.Errorf("[ArangoDB][CreateAppointment] begin transaction failed: %s", err)
		return constant.ArangoDB_Driver_Failed
	}
	trxCtx := driver.WithTransactionID(ctx, trxId)

	// 取得 consumer 在指定時間區段內的 appointments
	cursor1, err := db.Query(
		trxCtx,
		"FOR d IN Appointments FILTER d.consumerId == @consumerId AND d.courseStartAt >= @courseStartAt AND d.courseStartAt < @courseEndAt RETURN d",
		map[string]interface{}{
			"consumerId":    appt.ConsumerID,
			"courseStartAt": appt.CourseStartAt,
			"courseEndAt":   appt.CourseEndAt,
		},
	)
	defer closeCursor(cursor1)
	if err != nil {
		logger.Errorf("[ArangoDB][CreateAppointment] qeury failed: %s", err)
		if err = db.AbortTransaction(ctx, trxId, nil); err != nil {
			logger.Errorf("[ArangoDB][CreateAppointment] abort transaction failed: %s", err)
		}
		return constant.ArangoDB_Driver_Failed
	}

	// 檢查 consumer 在指定時間區段內是否已經有其他的 appointment 存在
	// 若有則表示預約衝突, 中斷 transaction
	tmp := &model.Appointment{}
	_, err = cursor1.ReadDocument(ctx, tmp)
	if !driver.IsNoMoreDocuments(err) {
		if err == nil {
			logger.Debugf("[ArangoDB][CreateAppointment] consumer %s has other appointments already.", appt.ConsumerID)
		} else {
			logger.Errorf("[ArangoDB][CreateAppointment] read document failed: %s", err)
		}
		if err = db.AbortTransaction(ctx, trxId, nil); err != nil {
			logger.Errorf("[ArangoDB][CreateAppointment] abort transaction failed: %s", err)
			return constant.ArangoDB_Driver_Failed
		}
		return constant.ArangoDB_Invalid_Operation
	}

	// 取得 schedule
	result := &model.Schedule{}
	scheduleCol, _ := db.Collection(trxCtx, "Schedules")
	_, err = scheduleCol.ReadDocument(trxCtx, appt.ScheduleID, result)
	if err != nil {
		logger.Errorf("[ArangoDB][CreateAppointment] read document failed: %s", err)
		if err = db.AbortTransaction(ctx, trxId, nil); err != nil {
			logger.Errorf("[ArangoDB][CreateAppointment] abort transaction failed: %s", err)
		}
		return constant.ArangoDB_Driver_Failed
	}

	// 確保 schedule 人數尚未達到規定上限
	if result.Count+1 > result.MaxConsumerLimit {
		logger.Debugf("[ArangoDB][CreateAppointment] schedule %s is fulled.", appt.ScheduleID)
		if err = db.AbortTransaction(ctx, trxId, nil); err != nil {
			logger.Errorf("[ArangoDB][CreateAppointment] abort transaction failed: %s", err)
			return constant.ArangoDB_Driver_Failed
		}
		return constant.ArangoDB_Invalid_Operation
	}

	// 更新 schedule 預約人數
	result.Count++
	_, err = scheduleCol.UpdateDocument(trxCtx, result.ID, map[string]interface{}{"count": result.Count})
	if err != nil {
		logger.Errorf("[ArangoDB][CreateAppointment] update document failed: %s", err)
		return constant.ArangoDB_Driver_Failed
	}

	// 寫入 appointment
	apptCol, _ := db.Collection(trxCtx, "Appointments")
	_, err = apptCol.CreateDocument(trxCtx, appt)
	if err != nil {
		logger.Errorf("[ArangoDB][CreateAppointment] create document failed: %s", err)
		return constant.ArangoDB_Driver_Failed
	}

	if err = db.CommitTransaction(ctx, trxId, nil); err != nil {
		logger.Errorf("[ArangoDB][CreateAppointment] commit transaction failed: %s", err)
		return constant.ArangoDB_Driver_Failed
	}

	return 0
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
