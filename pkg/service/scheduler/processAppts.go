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
		appts, code := arangodb.GetApptsByScheduleIDAndStatus(context.TODO(), s.ID, constant.Appointment_Status_Unsend_MeetingUrl)
		if code != constant.ArangoDB_Success {
			continue
		}

		for _, appt := range appts {
			lineSDK.SendMeetingUrl(appt.ProviderID, appt.ConsumerLineID, s.MeetingUrl)
			appt.Status = constant.Appointment_Status_Unsend_FeedbackUrl
		}

		_ = arangodb.UpdateApptsStatus(context.TODO(), appts, constant.Appointment_Status_Unsend_FeedbackUrl)
	}
}

func processReadyDismissAppts() {
	timestamp := utils.GetTimestampRoundToNextHalf()

	schedules, code := arangodb.GetReadyDismissSchedules(context.TODO(), timestamp)
	if code != constant.ArangoDB_Success {
		return
	}

	for _, s := range schedules {
		appts, code := arangodb.GetApptsByScheduleIDAndStatus(context.TODO(), s.ID, constant.Appointment_Status_Unsend_FeedbackUrl)
		if code != constant.ArangoDB_Success {
			continue
		}

		for _, appt := range appts {
			lineSDK.SendMeetingUrl(appt.ProviderID, appt.ConsumerLineID, "TODO: 回饋連結")
			appt.Status = constant.Appointment_Status_Unverified
		}

		_ = arangodb.UpdateApptsStatus(context.TODO(), appts, constant.Appointment_Status_Unverified)
	}
}
