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
var reqChan chan *Request
var clientMap syncmap                // map[string]*linebot.Client, key=lineAtChannelID
var providerMapping provider_mapping // map[string]*model.Provider, key=Provider.ID

type client struct {
	ProviderID    string
	ChannelID     string
	ChannelSecret string
	AccessToken   string
	Bot           *linebot.Client
}

func newClient(providerID, channelID, channelSecret, accessToken string) *client {
	bot, err := linebot.New(channelSecret, accessToken)
	if err != nil {
		logger.Errorf("[LineSDK] Init line bot failed: %v", err)
		return nil
	}
	return &client{
		ProviderID:    providerID,
		ChannelID:     channelID,
		ChannelSecret: channelSecret,
		AccessToken:   accessToken,
		Bot:           bot,
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
				provider.ID,
				provider.LineAtChannelID,
				provider.LineAtChannelSecret,
				provider.LineAtAccessToken,
			)
			if bot != nil {
				clientMap.Store(provider.LineAtChannelID, bot)
				providerMapping.Store(provider.ID, provider)
			}
		}

		reqChan = make(chan *Request, 4096)
		go startEventHandler()

		initialized = true
		logger.Debugf("[LineSDK] Initialized.")
	})

	return statusCode
}
