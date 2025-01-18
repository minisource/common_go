package middleware

// import (
// 	"fmt"
// 	"net/http"
// 	"strings"

// 	"github.com/gofiber/fiber/v2"
// 	"github.com/golang-jwt/jwt"
// 	"github.com/minisource/common_go/constants"
// 	"github.com/minisource/common_go/http/services"
// )

// func Authentication(cfg *services.JWTConfig) fiber.Handler {
// 	// var tokenUsecase = usecase.NewTokenUsecase(cfg)

// 	return func(c *fiber.Ctx) error {
// 		// var err error
// 		// claimMap := map[string]interface{}{}
// 		// auth := c.GetHeader(constants.AuthorizationHeaderKey)
// 		// token := strings.Split(auth, " ")
// 		// if auth == "" || len(token) < 2 {
// 		// 	err = &service_errors.ServiceError{EndUserMessage: service_errors.TokenRequired}
// 		// } else {
// 		// 	claimMap, err = tokenUsecase.GetClaims(token[1])
// 		// 	if err != nil {
// 		// 		switch err.(*jwt.ValidationError).Errors {
// 		// 		case jwt.ValidationErrorExpired:
// 		// 			err = &service_errors.ServiceError{EndUserMessage: service_errors.TokenExpired}
// 		// 		default:
// 		// 			err = &service_errors.ServiceError{EndUserMessage: service_errors.TokenInvalid}
// 		// 		}
// 		// 	}
// 		// }
// 		// if err != nil {
// 		// 	c.AbortWithStatusJSON(http.StatusUnauthorized, helper.GenerateBaseResponseWithError(
// 		// 		nil, false, helper.AuthError, err,
// 		// 	))
// 		// 	return
// 		// }

// 		// c.Set(constant.UserIdKey, claimMap[constant.UserIdKey])
// 		// c.Set(constant.FirstNameKey, claimMap[constant.FirstNameKey])
// 		// c.Set(constant.LastNameKey, claimMap[constant.LastNameKey])
// 		// c.Set(constant.UsernameKey, claimMap[constant.UsernameKey])
// 		// c.Set(constant.EmailKey, claimMap[constant.EmailKey])
// 		// c.Set(constant.MobileNumberKey, claimMap[constant.MobileNumberKey])
// 		// c.Set(constant.RolesKey, claimMap[constant.RolesKey])
// 		// c.Set(constant.ExpireTimeKey, claimMap[constant.ExpireTimeKey])

// 		return c.Next()
// 	}
// }