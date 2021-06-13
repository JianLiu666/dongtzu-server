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

	case constant.ApptStatus_Unstarted_Unconfirmed:
		// 1. 傳送 MeetingUrl 到指定頻道
		// 2. 押到 ApptStatus_Starting
		appt.Status = constant.ApptStatus_Starting
		arangodb.UpdateAppointment(context.TODO(), appt.ID, appt)

	case constant.ApptStatus_Unstarted_Confirmed:
		// 1. 傳送 MeetingUrl 到指定頻道
		// 2. 押到 ApptStatus_Starting
		appt.Status = constant.ApptStatus_Starting
		arangodb.UpdateAppointment(context.TODO(), appt.ID, appt)

	case constant.ApptStatus_Starting:
		// 異常情形
		logger.Errorf("[Scheduler] wrong tatus error: %v", appt.Status)

	case constant.ApptStatus_End_Unconfirmed:
		// 異常情形
		logger.Errorf("[Scheduler] wrong tatus error: %v", appt.Status)

	case constant.ApptStatus_End_Confirmed:
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

	case constant.ApptStatus_Unstarted_Unconfirmed:
		// 異常情形
		logger.Errorf("[Scheduler] wrong tatus error: %v", appt.Status)

	case constant.ApptStatus_Unstarted_Confirmed:
		// 異常情形
		logger.Errorf("[Scheduler] wrong tatus error: %v", appt.Status)

	case constant.ApptStatus_Starting:
		// 1. 傳送 FeedbackUrl 到指定頻道
		// 2. 押到 ApptStatus_End_WithoutFeedback
		appt.Status = constant.ApptStatus_End_Unconfirmed
		arangodb.UpdateAppointment(context.TODO(), appt.ID, appt)

	case constant.ApptStatus_End_Unconfirmed:
		// 異常情形
		logger.Errorf("[Scheduler] wrong tatus error: %v", appt.Status)

	case constant.ApptStatus_End_Confirmed:
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
