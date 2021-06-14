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

type RegisterProviderReq struct {
	LineUserID  string `json:"lineUserId"`
	RealName    string `json:"realName"`
	LineID      string `json:"lineId"`
	LineAtName  string `json:"lineAtName"`
	LineAtID    string `json:"LineAtID"`
	CountryCode string `json:"countryCode"`
	PhoneNum    string `json:"phoneNum"`
	GmailAddr   string `json:"gmailAddr"`
	InivteCode  string `json:"inivteCode"`
	MemeberTerm bool   `json:"memeberTerm"`
	PrivacyTerm bool   `json:"privacyTerm"`
	Status      int    `json:"status"`
}

type UpdateProviderInfoReq struct {
	RealName    string `json:"realName"`
	LineID      string `json:"lineId"`
	LineAtName  string `json:"lineAtName"`
	LineAtID    string `json:"LineAtID"`
	CountryCode string `json:"countryCode"`
	PhoneNum    string `json:"phoneNum"`
	GmailAddr   string `json:"gmailAddr"`
	Status      int    `json:"status"`
}

type RegisterOrUpdateProviderRes struct {
	StatusCode string `json:"status"`
}

func GetProviderInfo() fiber.Handler {
	return func(c *fiber.Ctx) error {
		lineUserID := c.Params("lineUserId")
		if lineUserID == "" {
			return c.Status(fasthttp.StatusNotFound).JSON(ErrRes{
				Code:       "404001",
				StatusCode: "404001",
				Message:    "No given Line user ID",
			})
		}

		ctx := context.Background()
		providerProfile, err := arangodb.GetProviderProfileByLineUserID(ctx, lineUserID)
		if err != nil {
			return c.Status(fasthttp.StatusInternalServerError).JSON(ErrRes{
				Code:       "500001",
				StatusCode: "500001",
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
		// 1. parsing post body
		// 2. register provider tx
		//		- get provider from lineUserId && status
		//      - deal with the provider from get result
		//      - create provider file
		return nil
	}
}

// Todo
// - 1. 之後再實作手機修改簡訊驗證邏輯
// - 2. 之後再實作gmail修改email + calendar同步邏輯
// 目前只實做完全信任前端覆蓋資料的邏輯
func UpdateProviderInfo() fiber.Handler {
	return func(c *fiber.Ctx) error {
		// 1. parsing put body
		// 2. update provider profile
		// 3. if params status is 2 && update success -> notify notion
		return nil
	}
}
