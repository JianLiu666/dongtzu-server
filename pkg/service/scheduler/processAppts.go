package scheduler

import (
	"context"
	"dongtzu/constant"
	"dongtzu/pkg/model"
	"dongtzu/pkg/repository/arangodb"

	"gitlab.geax.io/demeter/gologger/logger"
)

func processAppts() {
	defer updateTimeWheelIdx()

	for _, appt := range endTimeWheel[timeWheelIdx] {
		sendFeedbackUrl(appt)
	}
	endTimeWheel[timeWheelIdx] = []*model.Appointment{}

	for _, appt := range startTimeWheel[timeWheelIdx] {
		sendMeetingUrl(appt)
	}
	startTimeWheel[timeWheelIdx] = []*model.Appointment{}
}

func sendMeetingUrl(appt *model.Appointment) {
	switch appt.Status {
	case constant.ApptStatus_Exception:
		// 異常情形
		logger.Errorf("[Scheduler] wrong tatus error: %v", appt.Status)

	case constant.ApptStatus_Cancelled:
		// 異常情形
		logger.Errorf("[Scheduler] wrong tatus error: %v", appt.Status)

	case constant.ApptStatus_Unsend_MeetingUrl:
		// 1. 傳送 MeetingUrl 到指定頻道
		// 2. 押到 ApptStatus_Unsend_FeedbackUrl
		appt.Status = constant.ApptStatus_Unsend_FeedbackUrl
		arangodb.UpdateAppointment(context.TODO(), appt.ID, appt)

	case constant.ApptStatus_Unsend_FeedbackUrl:
		// 異常情形
		logger.Errorf("[Scheduler] wrong tatus error: %v", appt.Status)

	case constant.ApptStatus_Unverified:
		// 異常情形
		logger.Errorf("[Scheduler] wrong tatus error: %v", appt.Status)

	case constant.ApptStatus_Varified:
		// 異常情形
		logger.Errorf("[Scheduler] wrong tatus error: %v", appt.Status)

	default:
		// 異常情形
		logger.Errorf("[Scheduler] unknown status error: %v", appt.Status)
	}
}

func sendFeedbackUrl(appt *model.Appointment) {
	switch appt.Status {
	case constant.ApptStatus_Exception:
		// 異常情形
		logger.Errorf("[Scheduler] wrong tatus error: %v", appt.Status)

	case constant.ApptStatus_Cancelled:
		// 異常情形
		logger.Errorf("[Scheduler] wrong tatus error: %v", appt.Status)

	case constant.ApptStatus_Unsend_MeetingUrl:
		// 異常情形
		logger.Errorf("[Scheduler] wrong tatus error: %v", appt.Status)

	case constant.ApptStatus_Unsend_FeedbackUrl:
		// 1. 傳送 FeedbackUrl 到指定頻道
		// 2. 押到 ApptStatus_Unverified

	case constant.ApptStatus_Unverified:
		// 異常情形
		logger.Errorf("[Scheduler] wrong tatus error: %v", appt.Status)

	case constant.ApptStatus_Varified:
		// 異常情形
		logger.Errorf("[Scheduler] wrong tatus error: %v", appt.Status)

	default:
		// 異常情形
		logger.Errorf("[Scheduler] unknown status error: %v", appt.Status)
	}
}

func updateTimeWheelIdx() {
	timeWheelIdx++
	if timeWheelIdx == 60 {
		timeWheelIdx = 0
	}
}
