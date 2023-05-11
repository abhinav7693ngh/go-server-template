package controllers

import (
	"vas/errorCodes"
	"vas/types"
	"vas/utils"

	"github.com/gofiber/fiber/v2"
	"github.com/gookit/validate"
)

func ValiadateSamplePayload(c *fiber.Ctx) error {
	var samplePayload types.SamplePayload
	err := c.BodyParser(&samplePayload)
	if err != nil {
		return utils.ErrorResponse(
			c,
			utils.AppError{
				DebugMessage: err.Error(),
				ErrorCode:    errorCodes.BAD_REQUEST,
			},
			samplePayload,
		)
	}

	v := validate.Struct(samplePayload)
	if !v.Validate() {
		return utils.ErrorResponse(
			c,
			utils.AppError{
				DebugMessage: v.Errors.One(),
				ErrorCode:    errorCodes.BAD_REQUEST,
			},
			samplePayload,
		)
	}

	c.Locals("samplePayload", samplePayload)
	return c.Next()
}

func Sample(c *fiber.Ctx) error {
	samplePayload := c.Locals("samplePayload").(types.SamplePayload)

	return utils.SuccessResponse(c, fiber.Map{
		"text": samplePayload.Text,
	}, samplePayload)
}
