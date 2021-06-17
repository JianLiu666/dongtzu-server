package zoomSDK

import (
	"context"
	"dongtzu/constant"
	"dongtzu/pkg/repository/arangodb"
	"sync"

	"github.com/himalayan-institute/zoom-lib-golang"
	"gitlab.geax.io/demeter/gologger/logger"
)

var once sync.Once
var initialized bool
var clients []*handler
var clientIdx int
var mu sync.Mutex

type handler struct {
	UserId string
	Client *zoom.Client
}

func newHandler(userId, apiKey, apiSecret string) *handler {
	return &handler{
		UserId: userId,
		Client: zoom.NewClient(apiKey, apiSecret),
	}
}

func Init() int {
	if initialized {
		return constant.Initialization_Already
	}

	statusCode := constant.Initialization_Success
	once.Do(func() {
		zoomAccounts, code := arangodb.GetZoomAccounts(context.TODO())
		if code != constant.ArangoDB_Success {
			statusCode = constant.Initialization_Failed
			return
		}

		clients = []*handler{}
		for _, account := range zoomAccounts {
			clients = append(clients, newHandler(
				account.UserID,
				account.APIKey,
				account.APISecret,
			))
		}

		logger.Debugf("[ZoomSDK] Initialized.")
	})

	return statusCode
}

// round-robin
func getHandler() *handler {
	if len(clients) == 1 {
		return clients[0]
	}

	mu.Lock()
	defer mu.Unlock()

	clientIdx++
	if clientIdx >= len(clients) {
		clientIdx = 0
	}

	return clients[clientIdx]
}
