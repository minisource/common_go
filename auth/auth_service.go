package auth

import (
	"encoding/json"
	"errors"
	"fmt"

	"github.com/casdoor/casdoor-go-sdk/casdoorsdk"
	"github.com/minisource/common_go/auth/models"
	"github.com/minisource/common_go/http/helper"
)

var authService *AuthService

type AuthService struct {
	CasdoorClient *casdoorsdk.Client
	APIClient     *helper.APIClient
	Application   string
	Organization  string
	ClientID      string
	ClientSecret  string
	Endpoint      string
	Certificate   string
}

func NewAuthService(endpoint, clientID, clientSecret, certificate, organization, application string) *AuthService {
	casdoor := casdoorsdk.NewClient(
		endpoint,     // Casdoor server URL
		clientID,     // Client ID
		clientSecret, // Client Secret
		certificate,  // JWT public key or certificate
		organization, // Casdoor organization
		application,  // Casdoor application
	)

	authService = &AuthService{
		CasdoorClient: casdoor,
		APIClient:     helper.NewAPIClient(endpoint).SetBasicAuth(clientID, clientSecret),
		Application:   application,
		Organization:  organization,
		ClientID:      clientID,
		ClientSecret:  clientSecret,
		Endpoint:      endpoint,
		Certificate:   certificate,
	}
	return authService
}

func GetAuthService() *AuthService {
	return authService
}

// RegisterUser registers a new user with forbidden status
func (s *AuthService) RegisterUser(countryCode, phoneNumber string) error {
	// Fetch the newly created user to get the ID
	existedUser, err := s.CasdoorClient.GetUserByPhone(countryCode + phoneNumber)
	if existedUser != nil {
		return nil
	}

	// Create user with IsForbidden = true (inactive)
	user := &casdoorsdk.User{
		CountryCode: countryCode,
		Phone:       countryCode + phoneNumber,
		Name:        "user_" + phoneNumber,
	}
	success, err := s.CasdoorClient.AddUser(user)
	if err != nil || !success {
		return errors.New("failed to register user in Casdoor")
	}

	return nil
}

// SendOTP generates and sends an OTP to the user's phone
func (s *AuthService) SendOTP(phone string) error {
	params := map[string]string{
		"dest":          phone,
		"method":        "login",
		"type":          "phone",
		"captchaType":   "none",
		"applicationId": "admin/" + s.Application,
	}
	response, err := s.CasdoorClient.DoPost("send-verification-code", params, nil, false, false)
	if err != nil {
		return err
	}
	fmt.Println("Verification code sent successfully:", response)
	return nil
}

// VerifyOTP verifies the OTP
func (s *AuthService) VerifyCode(phone, code string) (bool, error) {
    params := map[string]interface{}{
        "username": phone,
        "code":     code,
    }

    var response map[string]interface{}

    err := s.APIClient.PostJSON("/api/verify-code", params, &response)

    if err != nil {
        return false, fmt.Errorf("request failed: %v", err)
    }

    fmt.Println("Verification response:", response)

    if status, ok := response["status"].(string); ok && status == "ok" {
        return true, nil
    }

    return false, fmt.Errorf("verification failed: %v", response)
}

// GenerateJWT generates a JWT token for the user
func (s *AuthService) GenerateJWT(username string) (*models.AccessTokenResponse, error) {
	// Placeholder using Casdoor's token generation
    formData := map[string]string{
        "grant_type":    "token",
        "username":   username,
        // "client_id":     s.ClientID,
        // "client_secret": s.ClientSecret,
    }
	resp := s.APIClient.Post("/api/login/oauth/access_token", formData)

    if resp.Error != nil {
        return nil, resp.Error
    }

    if resp.StatusCode != 200 {
        return nil, fmt.Errorf("token request failed with status %d: %s", 
            resp.StatusCode, string(resp.Body))
    }

    var tokenResp models.AccessTokenResponse
    if err := json.Unmarshal(resp.Body, &tokenResp); err != nil {
        return nil, fmt.Errorf("failed to parse token response: %v", err)
    }

    return &tokenResp, nil
}

func (s *AuthService) ValidateToken(accessToken string) (*casdoorsdk.IntrospectTokenResult, error) {
    // Introspect the token
    tokenInfo, err := s.CasdoorClient.IntrospectToken(accessToken, "access_token")
    if err != nil {
        return nil, fmt.Errorf("token introspection failed: %v", err)
    }

    if !tokenInfo.Active {
        return nil, fmt.Errorf("token is not active")
    }

    return tokenInfo, nil
}

func (s *AuthService) GetUserInfoByUsername(username string) (*casdoorsdk.User, error) {
    // Get user by username
    user, err := s.CasdoorClient.GetUser(username)
    if err != nil {
        return nil, fmt.Errorf("failed to get user info: %v", err)
    }

    if user == nil {
        return nil, fmt.Errorf("user '%s' not found", username)
    }

    return user, nil
}