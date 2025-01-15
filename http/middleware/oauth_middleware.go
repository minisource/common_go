package middleware

import (
	"strings"

	"github.com/gofiber/fiber/v2"
	"github.com/minisource/apiclients/auth"
	"github.com/minisource/apiclients/auth/models"
	"github.com/minisource/common_go/constants"
	helper "github.com/minisource/common_go/http/helper"
	"github.com/minisource/common_go/service_errors"
)

func OAuthValidationMiddleware(authApiClient *helper.APIClient) fiber.Handler {
	auth.NewAuthService(authApiClient)
	service := auth.GetAuthService()

	return func(c *fiber.Ctx) error {
		var err error
		tokenInfo := &models.ValidateTokenRes{}

		// Get the Authorization header
		authHeader := c.Get(constants.AuthorizationHeaderKey)
		tokenParts := strings.Split(authHeader, " ")

		// Check if the Authorization header is missing or invalid
		if authHeader == "" {
			err = &service_errors.ServiceError{EndUserMessage: service_errors.TokenRequired}
		} else {
			// Validate the token
			tokenInfo, err = service.ValidateToken(models.ValidateTokenReq{Token: tokenParts[1]})
		}

		// Handle validation errors
		if err != nil {
			return c.Status(fiber.StatusUnauthorized).JSON(
				helper.GenerateBaseResponseWithError(nil, false, helper.AuthError, err),
			)
		}

		// Store the token information in the context
		c.Locals("token_info", tokenInfo)

		// Proceed to the next middleware/handler
		return c.Next()
	}
}