package helper

import (
	"fmt"
	"io"
	"net/http"

	"github.com/minisource/common_go/http/services"
)

type APIClient struct {
    BaseURL     string
    JWTManager  *services.JWTManager
}

func NewAPIClient(baseURL string, jwtManager *services.JWTManager) *APIClient {
    return &APIClient{
        BaseURL:    baseURL,
        JWTManager: jwtManager,
    }
}

func (client *APIClient) GetResourceWithAuthorization(method, resourcePath string) (string, error) {
    // دریافت توکن
    token, err := client.JWTManager.GetToken()
    if err != nil {
        return "", err
    }

    // ارسال درخواست به API
    req, err := http.NewRequest(method, fmt.Sprintf("%s%s", client.BaseURL, resourcePath), nil)
    if err != nil {
        return "", err
    }

    // اضافه کردن هدر برای احراز هویت
    req.Header.Set("Authorization", "Bearer "+token.AccessToken)

    resp, err := http.DefaultClient.Do(req)
    if err != nil {
        return "", err
    }
    defer resp.Body.Close()

    body, err := io.ReadAll(resp.Body)
    if err != nil {
        return "", err
    }

    return string(body), nil
}


func (client *APIClient) GetResource(method, resourcePath string) (string, error) {
    // ارسال درخواست به API
    req, err := http.NewRequest(method, fmt.Sprintf("%s%s", client.BaseURL, resourcePath), nil)
    if err != nil {
        return "", err
    }

    resp, err := http.DefaultClient.Do(req)
    if err != nil {
        return "", err
    }
    defer resp.Body.Close()

    body, err := io.ReadAll(resp.Body)
    if err != nil {
        return "", err
    }

    return string(body), nil
}