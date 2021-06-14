package server

import (
	"context"
	"dongtzu/pkg/model"
	"dongtzu/pkg/repository/arangodb"

	"github.com/gofiber/fiber/v2"
	"github.com/valyala/fasthttp"
)

type ErrRes struct {
	Code       string `json:"code"`    // 系統自定義錯誤代碼
	StatusCode string `json:"status"`  // http 狀態碼
	Message    string `json:"message"` // 錯誤訊息
}

type GetProviderInfoRes struct {
	StatusCode string          `json:"statusCode"`
	Data       *model.Provider `json:"data"`
}

func GetProviderInfo() fiber.Handler {
	return func(c *fiber.Ctx) error {
		lineUserID := c.Params("lineUserId")
		if lineUserID == "" {
			return c.Status(fasthttp.StatusNotFound).JSON(ErrRes{
				Code:       "404001",
				StatusCode: "404",
				Message:    "No given Line user ID",
			})
		}

		ctx := context.Background()
		providerProfile, err := arangodb.GetProviderProfileByLineUserID(ctx, lineUserID)
		if err != nil {
			return c.Status(fasthttp.StatusInternalServerError).JSON(ErrRes{
				Code:       "500001",
				StatusCode: "500",
				Message:    "SERVER_ERROR",
			})
		}

		return c.Status(fasthttp.StatusOK).JSON(GetProviderInfoRes{
			StatusCode: "200",
			Data:       providerProfile,
		})
	}
}

func RegisterProvider() fiber.Handler {
	return func(c *fiber.Ctx) error {
		return nil
	}
}

func UpdateProviderInfo() fiber.Handler {
	return func(c *fiber.Ctx) error {
		return nil
	}
}
