package lineSDK

import (
	"crypto/hmac"
	"crypto/sha256"
	"dongtzu/constant"
	"encoding/base64"
	"encoding/json"

	"github.com/line/line-bot-sdk-go/linebot"
	"gitlab.geax.io/demeter/gologger/logger"
)

func HandleRequest(req *Request) int {
	if !initialized {
		return constant.Module_Initialization_Notyet
	}

	c, ok := clientMap.Load(req.ChannelID)
	if !ok {
		logger.Warnf("[LineSDK] Can not found line bot by channel id: %v", req.ChannelID)
		return constant.LineSDK_ChannelID_NotFound
	}

	req.addChannelSecret(c.ChannelSecret)
	if !req.validateSignature() {
		return constant.LineSDK_Request_Invalid
	}

	reqChan <- req
	return constant.LineSDK_Success
}

type Request struct {
	ChannelID     string
	channelSecret string
	Signature     string
	Body          []byte
}

func NewRequest(channelID, signature string, body []byte) *Request {
	return &Request{
		ChannelID:     channelID,
		channelSecret: "",
		Signature:     signature,
		Body:          body,
	}
}

func (req *Request) addChannelSecret(channelSecret string) {
	req.channelSecret = channelSecret
}

func (req *Request) validateSignature() bool {
	decoded, err := base64.StdEncoding.DecodeString(req.Signature)
	if err != nil {
		return false
	}
	hash := hmac.New(sha256.New, []byte(req.channelSecret))

	_, err = hash.Write(req.Body)
	if err != nil {
		return false
	}

	return hmac.Equal(decoded, hash.Sum(nil))
}

func (req *Request) parseEvents() ([]*linebot.Event, int) {
	result := &struct {
		Events []*linebot.Event `json:"events"`
	}{}
	if err := json.Unmarshal(req.Body, result); err != nil {

		return nil, constant.LineSDK_Event_ParseFaild
	}

	return result.Events, constant.LineSDK_Success
}
