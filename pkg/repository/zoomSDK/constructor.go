package zoomSDK

import (
	"dongtzu/config"

	"github.com/himalayan-institute/zoom-lib-golang"
	"gitlab.geax.io/demeter/gologger/logger"
)

var client *zoom.Client

func Init() {
	defer logger.Debugf("[ZoomSDK] Initialized.")

	client = zoom.NewClient(
		config.GetGlobalConfig().Zoom.ApiKey,
		config.GetGlobalConfig().Zoom.ApiSecret,
	)
}
