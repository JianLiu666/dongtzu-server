package zoomapi

import (
	"dongtzu/config"

	"github.com/donvito/zoom-go/zoomAPI"
	"gitlab.geax.io/demeter/gologger/logger"
)

var zoom zoomAPI.Client

func Init() {
	defer logger.Debugf("[ZoomAPI] Initialized.")

	zoom = zoomAPI.NewClient(
		config.GetGlobalConfig().Zoom.BaseUrl,
		config.GetGlobalConfig().Zoom.JwtToken,
	)
}
