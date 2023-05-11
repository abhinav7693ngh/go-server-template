package middlewares

import (
	"vas/constants"
	"vas/logger"

	"github.com/gofiber/fiber/v2"
)

func LogRequestMiddleware() fiber.Handler {
	return func(c *fiber.Ctx) error {
		if FindStringInSlice(constants.NonMetricRoutes, c.Path()) {
			return c.Next()
		}
		logger.LogInfo(c, "REQUEST", nil)
		return c.Next()
	}
}
