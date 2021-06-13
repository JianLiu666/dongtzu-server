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
