package services

import (
	"time"

	"github.com/golang-jwt/jwt"
	"github.com/minisource/common_go/constants"
	"github.com/minisource/common_go/service_errors"
)

type JWTConfig struct {
	AccessTokenExpireDuration  time.Duration `env:"JWT_ACCESSTOKEN_EXPIRE_DURATION"`
	RefreshTokenExpireDuration time.Duration `env:"JWT_REFRESHTOKEN_EXPIRE_DURATION"`
	Secret                     string        `env:"JWT_ACCESSTOKEN_SECRET"`
	RefreshSecret              string        `env:"JWT_REFRESHTOKEN_SECRET"`
}

type TokenService struct {
	cfg *JWTConfig
}

type TokenDetail struct {
	AccessToken            string `json:"accessToken"`
	RefreshToken           string `json:"refreshToken"`
	AccessTokenExpireTime  int64  `json:"accessTokenExpireTime"`
	RefreshTokenExpireTime int64  `json:"refreshTokenExpireTime"`
}

type TokenDto struct {
	UserId      int
	FirstName   string
	LastName    string
	Username    string
	PhoneNumber string
	Email       string
}

func NewTokenService(cfg *JWTConfig) *TokenService {
	return &TokenService{
		cfg: cfg,
	}
}

func (s *TokenService) GenerateToken(token TokenDto) (*TokenDetail, error) {
	td := &TokenDetail{}
	td.AccessTokenExpireTime = time.Now().Add(s.cfg.AccessTokenExpireDuration * time.Minute).Unix()
	td.RefreshTokenExpireTime = time.Now().Add(s.cfg.RefreshTokenExpireDuration * time.Minute).Unix()

	atc := jwt.MapClaims{}

	atc[constants.UserIdKey] = token.UserId
	atc[constants.FirstNameKey] = token.FirstName
	atc[constants.LastNameKey] = token.LastName
	atc[constants.UsernameKey] = token.Username
	atc[constants.EmailKey] = token.Email
	atc[constants.PhoneNumberKey] = token.PhoneNumber
	atc[constants.ExpireTimeKey] = td.AccessTokenExpireTime

	at := jwt.NewWithClaims(jwt.SigningMethodHS256, atc)

	var err error
	td.AccessToken, err = at.SignedString([]byte(s.cfg.Secret))

	if err != nil {
		return nil, err
	}

	rtc := jwt.MapClaims{}

	rtc[constants.UserIdKey] = token.UserId
	rtc[constants.ExpireTimeKey] = td.RefreshTokenExpireTime

	rt := jwt.NewWithClaims(jwt.SigningMethodHS256, rtc)

	td.RefreshToken, err = rt.SignedString([]byte(s.cfg.RefreshSecret))

	if err != nil {
		return nil, err
	}

	return td, nil
}

func (s *TokenService) VerifyToken(token string) (*jwt.Token, error) {
	at, err := jwt.Parse(token, func(token *jwt.Token) (interface{}, error) {
		_, ok := token.Method.(*jwt.SigningMethodHMAC)
		if !ok {
			return nil, &service_errors.ServiceError{EndUserMessage: service_errors.UnExpectedError}
		}
		return []byte(s.cfg.Secret), nil
	})
	if err != nil {
		return nil, err
	}
	return at, nil
}

func (s *TokenService) GetClaims(token string) (claimMap map[string]interface{}, err error) {
	claimMap = map[string]interface{}{}

	verifyToken, err := s.VerifyToken(token)
	if err != nil {
		return nil, err
	}
	claims, ok := verifyToken.Claims.(jwt.MapClaims)
	if ok && verifyToken.Valid {
		for k, v := range claims {
			claimMap[k] = v
		}
		return claimMap, nil
	}
	return nil, &service_errors.ServiceError{EndUserMessage: service_errors.ClaimsNotFound}
}
