package arangodb

import (
	"context"
	"dongtzu/constant"
	"dongtzu/pkg/model"
	"fmt"

	"github.com/arangodb/go-driver"
	"gitlab.geax.io/demeter/gologger/logger"
)

const (
	CollectionAppointments = "Appointments"
)

func GetAndConfirmApptsByStartTimestamp(ctx context.Context, startTimestamp, endTimestamp int64) []*model.Appointment {
	result := []*model.Appointment{}

	query := fmt.Sprintf("FOR d IN %s FILTER d.startTimestamp >= @startTimestamp AND d.startTimestamp < @endTimestamp AND d.status == @status RETURN d", CollectionAppointments)
	bindVars := map[string]interface{}{
		"startTimestamp": startTimestamp,
		"endTimestamp":   endTimestamp,
		"status":         constant.ApptStatus_Unstarted_Unconfirmed,
	}
	cursor, err := db.Query(ctx, query, bindVars)
	defer func() {
		if err := cursor.Close(); err != nil {
			logger.Errorf("[ArangoDB] cursor close falied: %v", err)
		}
	}()

	if err != nil {
		logger.Errorf("[ArangoDB] GetAndConfirmApptsByStartTimestamp query falied: %v", err)
		return []*model.Appointment{}
	}

	for {
		var data model.Appointment
		_, err := cursor.ReadDocument(ctx, &data)
		if driver.IsNoMoreDocuments(err) {
			break
		} else if err != nil {
			logger.Errorf("[ArangoDB] GetAndConfirmApptsByStartTimestamp cursor falied: %v", err)
			return []*model.Appointment{}
		}

		result = append(result, &data)
	}

	// TODO: 這裡應該要再將這些 unconfirmed 的 status 押成 confirmed 後在寫回 db

	return result
}

func GetAndConfirmApptsByEndTimestamp(ctx context.Context, startTimestamp, endTimestamp int64) []*model.Appointment {
	result := []*model.Appointment{}

	query := fmt.Sprintf("FOR d IN %s FILTER d.endTimestamp >= @startTimestamp AND d.endTimestamp < @endTimestamp AND d.status == @status RETURN d", CollectionAppointments)
	bindVars := map[string]interface{}{
		"startTimestamp": startTimestamp,
		"endTimestamp":   endTimestamp,
		"status":         constant.ApptStatus_End_Unconfirmed,
	}
	cursor, err := db.Query(ctx, query, bindVars)
	defer func() {
		if err := cursor.Close(); err != nil {
			logger.Errorf("[ArangoDB] cursor close falied: %v", err)
		}
	}()

	if err != nil {
		logger.Errorf("[ArangoDB] GetAndConfirmApptsByEndTimestamp query falied: %v", err)
		return []*model.Appointment{}
	}

	for {
		var data model.Appointment
		_, err := cursor.ReadDocument(ctx, &data)
		if driver.IsNoMoreDocuments(err) {
			break
		} else if err != nil {
			logger.Errorf("[ArangoDB] GetAndConfirmApptsByEndTimestamp cursor falied: %v", err)
			return []*model.Appointment{}
		}

		result = append(result, &data)
	}

	// TODO: 這裡應該要再將這些 unconfirmed 的 status 押成 confirmed 後在寫回 db

	return result
}

func UpdateAppointment(ctx context.Context, key string, appt *model.Appointment) {
	col, err := db.Collection(ctx, CollectionAppointments)
	if err != nil {
		if err != nil {
			logger.Errorf("[ArangoDB] UpdateAppointment falied: %v", err)
			return
		}
	}

	_, err = col.UpdateDocument(ctx, key, appt)
	if err != nil {
		logger.Errorf("[ArangoDB] UpdateAppointment falied: %v", err)
		return
	}
}

// 建立預約
// 1. 檢查 schedule 預約人數是否已經達到上限
// 2. 檢查 consumer 在相同時段中是否已經有其他的 appointment 存在
// 3. 建立 appointment 並更新 schedule 統計人數
func CreateAppointment(appt *model.Appointment) int {
	ctx := context.TODO()

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
		return 1
	}
	trxCtx := driver.WithTransactionID(ctx, trxId)

	// 取得 consumer 在指定時間區段內的 appointments
	cursor1, err := db.Query(
		trxCtx,
		"FOR d IN Appointments FILTER d.consumerId == @consumerId AND d.startTimestamp >= @startTimestamp AND d.startTimestamp < @endTimestamp RETURN d",
		map[string]interface{}{
			"consumerId":     appt.ConsumerID,
			"startTimestamp": appt.StartTimestamp,
			"endTimestamp":   appt.EndTimestamp,
		},
	)
	defer closeCursor(cursor1)
	if err != nil {
		logger.Errorf("[ArangoDB][CreateAppointment] qeury failed: %s", err)
		if err = db.AbortTransaction(ctx, trxId, nil); err != nil {
			logger.Errorf("[ArangoDB][CreateAppointment] abort transaction failed: %s", err)
		}
		return 1
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
		}
		return 1
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
		return 1
	}

	// 確保 schedule 人數尚未達到規定上限
	if result.Count+1 > result.MaxConsumerLimit {
		logger.Debugf("[ArangoDB][CreateAppointment] schedule %s is fulled.", appt.ScheduleID)
		if err = db.AbortTransaction(ctx, trxId, nil); err != nil {
			logger.Errorf("[ArangoDB][CreateAppointment] abort transaction failed: %s", err)
		}
		return 1
	}

	// 更新 schedule 預約人數
	result.Count++
	_, err = scheduleCol.UpdateDocument(trxCtx, result.ID, map[string]interface{}{"count": result.Count})
	if err != nil {
		logger.Errorf("[ArangoDB][CreateAppointment] update document failed: %s", err)
		return 1
	}

	// 寫入 appointment
	apptCol, _ := db.Collection(trxCtx, "Appointments")
	_, err = apptCol.CreateDocument(trxCtx, appt)
	if err != nil {
		logger.Errorf("[ArangoDB][CreateAppointment] create document failed: %s", err)
		return 1
	}

	if err = db.CommitTransaction(ctx, trxId, nil); err != nil {
		logger.Errorf("[ArangoDB][CreateAppointment] commit transaction failed: %s", err)
		return 1
	}

	return 0
}

func closeCursor(cursor driver.Cursor) {
	if err := cursor.Close(); err != nil {
		logger.Errorf("[ArangoDB] close cursor failed: %v", err)
	}
}
