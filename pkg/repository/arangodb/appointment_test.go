package arangodb

import (
	"context"
	"testing"

	"github.com/arangodb/go-driver"
	"github.com/arangodb/go-driver/http"
	"gitlab.geax.io/demeter/gologger/logger"
)

func TestCreateAppointment(t *testing.T) {
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

	// appt := &model.Appointment{
	// 	ID:         "",
	// 	ProviderID: "1",
	// }
}
