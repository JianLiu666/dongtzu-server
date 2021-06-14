package server

import (
	"github.com/gofiber/fiber/v2"
)

func appointment() fiber.Handler {
	return func(c *fiber.Ctx) error {

		// TODO: 實作預約流程
		// arangodb.UpdateAppointment()

		return nil
	}
}
