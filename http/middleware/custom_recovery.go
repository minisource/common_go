package middleware

import (
	"github.com/gofiber/fiber/v2"
	helper "github.com/minisource/common_go/http/helper"
)

func ErrorHandler(c *fiber.Ctx, err any) error {
	if err, ok := err.(error); ok {
		httpResponse := helper.GenerateBaseResponseWithError(nil, false, helper.CustomRecovery, err)
		return c.Status(fiber.StatusInternalServerError).JSON(httpResponse)
	}
	httpResponse := helper.GenerateBaseResponseWithAnyError(nil, false, helper.CustomRecovery, err)
	return c.Status(fiber.StatusInternalServerError).JSON(httpResponse)
}