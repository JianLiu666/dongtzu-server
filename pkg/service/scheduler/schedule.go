package scheduler

import (
	"context"
	"dongtzu/pkg/repository/arangodb"
)

func updateScheduleAndCreateZoomUrl() {
	schedules, _ := arangodb.GetWithoutMeetingUrlSchedules(context.TODO())
	for _, s := range schedules {
		// TODO: 串 zoom api
		s.MeetingUrl = "https://www.google.com/"
	}

	_ = arangodb.UpdateSchedulesMeetingUrl(context.TODO(), schedules)
}

func createZoomMeetingUrl() string {
	return ""
}
