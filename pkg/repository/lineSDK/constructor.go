package lineSDK

import (
	"context"
	"dongtzu/constant"
	"dongtzu/pkg/repository/arangodb"
	"sync"

	"github.com/line/line-bot-sdk-go/linebot"
	"gitlab.geax.io/demeter/gologger/logger"
)

var once sync.Once
var initialized bool
var clientMap syncmap // map[string]*linebot.Client
var reqChan chan *Request

type client struct {
	ChannelID     string
	ChannelSecret string
	AccessToken   string
	Client        *linebot.Client
}

func newClient(channelID, channelSecret, accessToken string) *client {
	bot, err := linebot.New(channelSecret, accessToken)
	if err != nil {
		logger.Errorf("[LineSDK] Init line bot failed: %v", err)
		return nil
	}
	return &client{
		ChannelID:     channelID,
		ChannelSecret: channelSecret,
		AccessToken:   accessToken,
		Client:        bot,
	}
}

func Init() int {
	if initialized {
		return constant.Module_Initialization_Already
	}

	statusCode := constant.Module_Initialization_Success
	once.Do(func() {
		providers, code := arangodb.GetProviders(context.TODO())
		if code != constant.ArangoDB_Success {
			statusCode = constant.Module_Initialization_Failed
			return
		}

		for _, provider := range providers {
			bot := newClient(
				provider.LineAtChannelID,
				provider.LineAtChannelSecret,
				provider.LineAtAccessToken,
			)
			if bot != nil {
				clientMap.Store(provider.LineAtChannelID, bot)
			}
		}

		reqChan = make(chan *Request, 4096)
		go startEventHandler()

		initialized = true
		logger.Debugf("[LineSDK] Initialized.")
	})

	return statusCode
}
