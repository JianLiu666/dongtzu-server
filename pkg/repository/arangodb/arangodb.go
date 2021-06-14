package arangodb

import (
	"context"
	"dongtzu/config"
	"sync"

	"github.com/arangodb/go-driver"
	"github.com/arangodb/go-driver/http"
	"gitlab.geax.io/demeter/gologger/logger"
)

var once sync.Once
var db driver.Database

func Init() {
	once.Do(func() {
		conn, err := http.NewConnection(http.ConnectionConfig{
			Endpoints: []string{config.GetGlobalConfig().ArangoDB.Addr},
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
		logger.Errorf("[ArangoDB] close cursor failed: %v", err)
	}
}
