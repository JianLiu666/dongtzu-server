package arangodb

import (
	"context"
	"dongtzu/config"
	"net"
	defaulthttp "net/http"
	"sync"
	"time"

	"github.com/arangodb/go-driver"
	"github.com/arangodb/go-driver/http"
	"gitlab.geax.io/demeter/gologger/logger"
)

const (
	collectionAppointments    = "Appointments"
	collectionCourses         = "Courses"
	collectionConsumers       = "Consumers"
	collectionFeedbacks       = "Feedbacks"
	collectionMonthReceipts   = "MonthReceipts"
	collectionOrders          = "Orders"
	collectionPayments        = "Payments"
	collectionPaymentMethods  = "PaymentMethods"
	collectionProviders       = "Providers"
	collectionSchedules       = "Schedules"
	collectionScheduleRules   = "ScheduleRules"
	collectionServiceProducts = "ServiceProducts"
	collectionZoomAccounts    = "ZoomAccounts"
)

var once sync.Once
var db driver.Database

func Init() {
	once.Do(func() {
		transport := &defaulthttp.Transport{
			DialContext: (&net.Dialer{
				KeepAlive: 60 * time.Second,
				DualStack: true}).DialContext,
			MaxIdleConns:          0,
			IdleConnTimeout:       30 * time.Second,
			ExpectContinueTimeout: 1 * time.Second,
		}

		conn, err := http.NewConnection(http.ConnectionConfig{
			Endpoints: []string{config.GetGlobalConfig().ArangoDB.Addr},
			ConnLimit: config.GetGlobalConfig().ArangoDB.ConnLimit,
			Transport: transport,
		})
		if err != nil {
			logger.Errorf("[ArangoDB] Init http connection failed: %v", err)
			return
		}

		c, err := driver.NewClient(driver.ClientConfig{
			Connection: conn,
			Authentication: driver.BasicAuthentication(
				config.GetGlobalConfig().ArangoDB.Username,
				config.GetGlobalConfig().ArangoDB.Password),
		})
		if err != nil {
			logger.Errorf("[ArangoDB] New arangodb client failed: %v", err)
			return
		}

		db, err = c.Database(context.TODO(), config.GetGlobalConfig().ArangoDB.DBName)
		if err != nil {
			logger.Errorf("[ArangoDB] New arangodb client failed: %v", err)
			return
		}

		logger.Debugf("[ArangoDB] Initialized.")
	})
}

func closeCursor(cursor driver.Cursor) {
	if err := cursor.Close(); err != nil {
		logger.Errorf("[ArangoDB] failed to close cursor: %v", err)
	}
}
