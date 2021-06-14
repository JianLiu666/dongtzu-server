package server

import (
	"github.com/gofiber/fiber/v2"
)

func appointment() fiber.Handler {
	return func(c *fiber.Ctx) error {
		// TODO:
		// 1. 檢查封包內容是否完整
		// ---- tx start
		// 2. 檢查 Schedule 預約人數是否已經達到上限
		// 3. 檢查 Consumer 在相同時段中是否已經有其他的 appointment 存在
		// 4. 建立 appointment, 並更新 schedule 統計人數
		// ---- tx end
		return nil
	}
}
