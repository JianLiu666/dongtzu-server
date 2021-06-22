package scheduler

import (
	"context"
	"dongtzu/constant"
	"dongtzu/pkg/repository/arangodb"
	"dongtzu/pkg/repository/zoomSDK"
	"time"
)

func createMeetingUrl() {
	schedules, code := arangodb.GetUncreatedMeetingUrlSchedules(context.TODO())
	if code != constant.ArangoDB_Success {
		return
	}

	for _, s := range schedules {
		scheduleTime := time.Unix(s.CourseStartAt, 0)
		minuteInteger := (s.CourseEndAt - s.CourseStartAt) / 60
		meetingUrl, code := zoomSDK.GetMeetingUrl(scheduleTime, int(minuteInteger))

		if code != constant.ZoomSDK_Success {
			continue
		}
		s.MeetingUrl = meetingUrl
	}

	_ = arangodb.UpdateSchedulesMeetingUrl(context.TODO(), schedules)
}
