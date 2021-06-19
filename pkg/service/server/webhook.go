package server

import (
	"dongtzu/constant"
	"dongtzu/pkg/repository/lineSDK"

	"github.com/gofiber/fiber/v2"
	"gitlab.geax.io/demeter/gologger/logger"
)

func lineWebhook() fiber.Handler {
	return func(c *fiber.Ctx) error {
		channelID := c.Params("channelId")

		code := lineSDK.HandleRequest(lineSDK.NewRequest(
			channelID,
			c.Get("X-Line-Signature"),
			c.Request().Body(),
		))
		if code != constant.LineSDK_Success {
			logger.Warnf("[Server] handle line request failed: %v", code)
		}

		return c.JSON(fiber.Map{"code": 200})
	}
}
