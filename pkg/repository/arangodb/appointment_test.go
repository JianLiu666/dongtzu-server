package arangodb

import (
	"context"
	"dongtzu/pkg/model"
	"testing"

	"github.com/arangodb/go-driver"
	"github.com/arangodb/go-driver/http"
	"gitlab.geax.io/demeter/gologger/logger"
)

func TestCreateAppointment(t *testing.T) {
	logger.Init("debug")

	conn, err := http.NewConnection(http.ConnectionConfig{
		Endpoints: []string{"http://127.0.0.1:8529/"},
	})
	if err != nil {
		logger.Errorf("[ArangoDB] Init http connection failed: %v", err)
		return
	}

	c, err := driver.NewClient(driver.ClientConfig{
		Connection: conn,
	})
	if err != nil {
		logger.Errorf("[ArangoDB] New arangodb client failed: %v", err)
		return
	}

	db, err = c.Database(context.TODO(), "_system")
	if err != nil {
		logger.Errorf("[ArangoDB] New arangodb client failed: %v", err)
		return
	}

	appt := &model.Appointment{
		ID:             "",
		ProviderID:     "1",
		ScheduleID:     "103847",
		ConsumerID:     "1",
		FeedbackID:     "",
		StartTimestamp: 1623776400,
		EndTimestamp:   1623780000,
		Note:           "",
		Status:         0,
	}

	status := CreateAppointment(context.TODO(), appt)
	logger.Debugf("%v", status)
}
