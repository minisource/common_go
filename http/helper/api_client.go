package helper

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/minisource/common_go/http/services"
)

type APIClient struct {
	BaseURL    string
	JWTManager *services.JWTManager
}

func NewAPIClient(baseURL string, jwtManager *services.JWTManager) *APIClient {
	return &APIClient{
		BaseURL:    baseURL,
		JWTManager: jwtManager,
	}
}

func (client *APIClient) GetResourceWithAuthorization(method, resourcePath string, bodyData interface{}) ([]byte, error) {
	// دریافت توکن
	token, err := client.JWTManager.GetToken()
	if err != nil {
		return nil, err
	}

	var reqBody io.Reader

	// اگر متد نیاز به بادی دارد، داده‌ها را تبدیل به JSON می‌کنیم.
	if bodyData != nil {
		jsonBody, err := json.Marshal(bodyData)
		if err != nil {
			return nil, fmt.Errorf("error marshaling body data: %v", err)
		}
		reqBody = bytes.NewReader(jsonBody)
	}

	// ایجاد درخواست HTTP
	req, err := http.NewRequest(method, fmt.Sprintf("%s%s", client.BaseURL, resourcePath), reqBody)
	if err != nil {
		return nil, fmt.Errorf("error creating request: %v", err)
	}

	// اضافه کردن هدر برای احراز هویت
	req.Header.Set("Authorization", "Bearer "+token.AccessToken)

	// اگر متد `POST`, `PUT`, `PATCH` است، باید Content-Type را تنظیم کنیم
	if method == http.MethodPost || method == http.MethodPut || method == http.MethodPatch {
		req.Header.Set("Content-Type", "application/json")
	}

	// ارسال درخواست
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("error executing request: %v", err)
	}
	defer resp.Body.Close()

	// بررسی وضعیت پاسخ
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("received non-200 response status: %d", resp.StatusCode)
	}

	// خواندن بدنه پاسخ
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("error reading response body: %v", err)
	}

	return body, nil
}

func (client *APIClient) MakeRequest(method, resourcePath string, bodyData interface{}) ([]byte, error) {
	var reqBody io.Reader

	// اگر متد نیاز به بادی دارد، داده‌ها را تبدیل به JSON می‌کنیم.
	if bodyData != nil {
		jsonBody, err := json.Marshal(bodyData)
		if err != nil {
			return nil, fmt.Errorf("error marshaling body data: %v", err)
		}
		reqBody = bytes.NewReader(jsonBody)
	}

	// ایجاد درخواست HTTP
	req, err := http.NewRequest(method, fmt.Sprintf("%s%s", client.BaseURL, resourcePath), reqBody)
	if err != nil {
		return nil, fmt.Errorf("error creating request: %v", err)
	}

	// اگر متد `POST`, `PUT`, `PATCH` است، باید Content-Type را تنظیم کنیم
	if method == http.MethodPost || method == http.MethodPut || method == http.MethodPatch {
		req.Header.Set("Content-Type", "application/json")
	}

	// ارسال درخواست
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("error executing request: %v", err)
	}
	defer resp.Body.Close()

	// بررسی وضعیت پاسخ
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("received non-200 response status: %d", resp.StatusCode)
	}

	// خواندن بدنه پاسخ
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("error reading response body: %v", err)
	}

	return body, nil
}
