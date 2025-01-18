package middleware

import (
	"strings"

	"github.com/gofiber/fiber/v2"
	"github.com/minisource/apiclients/auth/auth"
	"github.com/minisource/apiclients/auth/auth/models"
	"github.com/minisource/common_go/constants"
	helper "github.com/minisource/common_go/http/helper"
	"github.com/minisource/common_go/service_errors"
)

func Authentication(authApiClient *helper.APIClient) fiber.Handler {
	auth.NewAuthService(authApiClient)
	service := auth.GetAuthService()

	return func(c *fiber.Ctx) error {
		var err error
		var claimMap *models.ValidateAuthTokenRes
		auth := c.Get(constants.AuthorizationHeaderKey)
		token := strings.Split(auth, " ")
		if auth == "" || len(token) < 2 {
			err = &service_errors.ServiceError{EndUserMessage: service_errors.TokenRequired}
		} else {
			claimMap, err = service.ValidateAccessToken(models.ValidateAccessTokenRequest{AccessToken: token[1]})
			if err != nil {
				return err
			}
		}
		if err != nil {
			return c.Status(fiber.StatusUnauthorized).JSON(
				helper.GenerateBaseResponseWithError(nil, false, helper.AuthError, err),
			)
		}

		if userId, ok := claimMap.Claims[constants.UserIdKey].(string); ok {
			c.Set(constants.UserIdKey, userId)
		}
		if firstName, ok := claimMap.Claims[constants.FirstNameKey].(string); ok {
			c.Set(constants.FirstNameKey, firstName)
		}
		if lastName, ok := claimMap.Claims[constants.LastNameKey].(string); ok {
			c.Set(constants.LastNameKey, lastName)
		}
		if username, ok := claimMap.Claims[constants.UsernameKey].(string); ok {
			c.Set(constants.UsernameKey, username)
		}
		if email, ok := claimMap.Claims[constants.EmailKey].(string); ok {
			c.Set(constants.EmailKey, email)
		}
		if phoneNumber, ok := claimMap.Claims[constants.PhoneNumberKey].(string); ok {
			c.Set(constants.PhoneNumberKey, phoneNumber)
		}
		if expireTime, ok := claimMap.Claims[constants.ExpireTimeKey].(string); ok {
			c.Set(constants.ExpireTimeKey, expireTime)
		}

		return c.Next()
	}
}