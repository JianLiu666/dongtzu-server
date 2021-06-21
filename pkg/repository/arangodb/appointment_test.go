package arangodb

import (
	"context"
	"dongtzu/pkg/model"
	"testing"

	"gitlab.geax.io/demeter/gologger/logger"
)

func TestCreateAppointment(t *testing.T) {
	initConfig()

	appt := &model.Appointment{
		ID:            "",
		ProviderID:    "1",
		ScheduleID:    "103847",
		ConsumerID:    "1",
		FeedbackID:    "",
		CourseStartAt: 1623776400,
		CourseEndAt:   1623780000,
		Note:          "",
		Status:        0,
	}

	status := CreateAppointment(context.TODO(), appt)
	logger.Debugf("%v", status)
}
