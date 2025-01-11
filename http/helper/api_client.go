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

func (client *APIClient) GetResourceWithAuthorization(method, resourcePath string) ([]byte, error) {
    // دریافت توکن
    token, err := client.JWTManager.GetToken()
    if err != nil {
        return nil, err
    }

    // ارسال درخواست به API
    req, err := http.NewRequest(method, fmt.Sprintf("%s%s", client.BaseURL, resourcePath), nil)
    if err != nil {
        return nil, err
    }

    // اضافه کردن هدر برای احراز هویت
    req.Header.Set("Authorization", "Bearer "+token.AccessToken)

    resp, err := http.DefaultClient.Do(req)
    if err != nil {
        return nil, err
    }
    defer resp.Body.Close()

    if resp.StatusCode != http.StatusOK {
        return nil, fmt.Errorf("received non-200 response status: %d", resp.StatusCode)
    }

    body, err := io.ReadAll(resp.Body)
    if err != nil {
        return nil, err
    }

    return body, nil
}


func (client *APIClient) GetResource(method, resourcePath string) ([]byte, error) {
    // ارسال درخواست به API
    req, err := http.NewRequest(method, fmt.Sprintf("%s%s", client.BaseURL, resourcePath), nil)
    if err != nil {
        return nil, err
    }

    resp, err := http.DefaultClient.Do(req)
    if err != nil {
        return nil, err
    }
    defer resp.Body.Close()

    if resp.StatusCode != http.StatusOK {
        return nil, fmt.Errorf("received non-200 response status: %d", resp.StatusCode)
    }

    body, err := io.ReadAll(resp.Body)
    if err != nil {
        return nil, err
    }

    return body, nil
}