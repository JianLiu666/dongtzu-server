package scheduler

import (
	"context"
	"dongtzu/constant"
	"dongtzu/pkg/repository/arangodb"
	"dongtzu/pkg/repository/lineSDK"
	"dongtzu/pkg/utils"
)

func processReadyStartAppts() {
	timestamp := utils.GetTimestampRoundToNextHalf()

	schedules, code := arangodb.GetReadyStartSchedules(context.TODO(), timestamp)
	if code != constant.ArangoDB_Success {
		return
	}

	for _, s := range schedules {
		appts, code := arangodb.GetApptsByScheduleIDAndStatus(context.TODO(), s.ID, constant.ApptStatus_Unsend_MeetingUrl)
		if code != constant.ArangoDB_Success {
			continue
		}

		for _, appt := range appts {
			lineSDK.PushMessage(appt.ConsumerLineID, s.MeetingUrl)
			appt.Status = constant.ApptStatus_Unsend_FeedbackUrl
		}

		_ = arangodb.UpdateApptsStatus(context.TODO(), appts, constant.ApptStatus_Unsend_FeedbackUrl)
	}
}

func processReadyDismissAppts() {
	timestamp := utils.GetTimestampRoundToNextHalf()

	schedules, code := arangodb.GetReadyDismissSchedules(context.TODO(), timestamp)
	if code != constant.ArangoDB_Success {
		return
	}

	for _, s := range schedules {
		appts, code := arangodb.GetApptsByScheduleIDAndStatus(context.TODO(), s.ID, constant.ApptStatus_Unsend_FeedbackUrl)
		if code != constant.ArangoDB_Success {
			continue
		}

		for _, appt := range appts {
			lineSDK.PushMessage(appt.ConsumerLineID, "TODO: 回饋連結")
			appt.Status = constant.ApptStatus_Unverified
		}

		_ = arangodb.UpdateApptsStatus(context.TODO(), appts, constant.ApptStatus_Unverified)
	}
}
