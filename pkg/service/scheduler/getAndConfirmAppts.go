package scheduler

import (
	"context"
	"dongtzu/pkg/repository/arangodb"
	"dongtzu/pkg/utils"
	"time"
)

func getAndConfirmAppts() {
	startTimestamp, endTimestamp := utils.GetTimeRange()

	apptsList := arangodb.GetAndConfirmApptsByStartTimestamp(context.TODO(), startTimestamp, endTimestamp)
	for _, appt := range apptsList {
		t := time.Unix(appt.StartTimestamp, 0).UTC()
		startTimeWheel[t.Minute()] = append(startTimeWheel[t.Minute()], appt)
	}

	apptsList = arangodb.GetAndConfirmApptsByStartTimestamp(context.TODO(), startTimestamp, endTimestamp)
	for _, appt := range apptsList {
		t := time.Unix(appt.EndTimestamp, 0).UTC()
		endTimeWheel[t.Minute()] = append(endTimeWheel[t.Minute()], appt)
	}
}
