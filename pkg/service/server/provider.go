package server

import (
	"context"
	"dongtzu/constant"
	"dongtzu/pkg/model"
	"dongtzu/pkg/repository/arangodb"
	"dongtzu/pkg/repository/githubSDK"

	"github.com/gofiber/fiber/v2"
	"github.com/valyala/fasthttp"
	"gitlab.geax.io/demeter/gologger/logger"
)

func GetProviderInfo() fiber.Handler {
	return func(c *fiber.Ctx) error {
		lineUserID := c.Params("lineUserId")
		if lineUserID == "" {
			return c.Status(fasthttp.StatusNotFound).JSON(model.ErrRes{
				Code:       "404001",
				StatusCode: "404001",
				Message:    "No given Line user ID",
			})
		}

		ctx := context.Background()
		providerProfile, err := arangodb.GetProviderProfileByLineUserID(ctx, lineUserID)
		if err != nil {
			return c.Status(fasthttp.StatusInternalServerError).JSON(model.ErrRes{
				Code:       "500001",
				StatusCode: "500001",
				Message:    "SERVER_ERROR",
			})
		}

		return c.Status(fasthttp.StatusOK).JSON(model.GetProviderInfoRes{
			StatusCode: "200",
			Data:       providerProfile,
		})
	}
}

func RegisterProvider() fiber.Handler {
	return func(c *fiber.Ctx) error {
		// 1. parsing post body
		var body model.RegisterProviderReq
		if err := c.BodyParser(&body); err != nil {
			return c.Status(fasthttp.StatusNotFound).JSON(model.ErrRes{
				Code:       "404001",
				StatusCode: "404001",
				Message:    "No validate registeration input",
			})
		}

		// 2. register provider tx
		ctx := context.Background()
		err := arangodb.CreateProviderProfile(ctx, body)
		if err != nil {
			return c.Status(fasthttp.StatusInternalServerError).JSON(model.ErrRes{
				Code:       "500001",
				StatusCode: "500001",
				Message:    "SERVER_ERROR",
			})
		}

		return c.Status(fasthttp.StatusOK).JSON(model.RegisterOrUpdateProviderRes{
			StatusCode: "200001",
		})
	}
}

// Todo
// - 1. 之後再實作手機修改簡訊驗證邏輯
// - 2. 之後再實作gmail修改email + calendar同步邏輯
// 目前只實做完全信任前端覆蓋資料的邏輯
func UpdateProviderInfo() fiber.Handler {
	return func(c *fiber.Ctx) error {
		lineUserID := c.Params("lineUserId")
		if lineUserID == "" {
			return c.Status(fasthttp.StatusNotFound).JSON(model.ErrRes{
				Code:       "404001",
				StatusCode: "404001",
				Message:    "No given Line user ID",
			})
		}

		var body model.UpdateProviderInfoReq
		if err := c.BodyParser(&body); err != nil {
			return c.Status(fasthttp.StatusNotFound).JSON(model.ErrRes{
				Code:       "404001",
				StatusCode: "404001",
				Message:    "No validate edit provider input",
			})
		}

		// 2. update provider profile
		ctx := context.Background()
		err := arangodb.UpdateProviderByLineUserID(ctx, lineUserID, body)
		if err != nil {
			return c.Status(fasthttp.StatusInternalServerError).JSON(model.ErrRes{
				Code:       "500001",
				StatusCode: "500001",
				Message:    "SERVER_ERROR",
			})
		}

		// 3. if params status is 2 && update success -> github 串接
		if body.Status == constant.Provider_Status_Auditing {
			profile, _ := arangodb.GetProviderProfileByLineUserID(ctx, lineUserID)
			err = githubSDK.CreateIssueForProvider(*profile)
			if err != nil {
				logger.Errorf("[githubSDK] create issue failure : %v", err)
			}
		}

		return c.Status(fasthttp.StatusOK).JSON(model.RegisterOrUpdateProviderRes{
			StatusCode: "200001",
		})
	}
}
