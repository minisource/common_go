package services

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"sync"
	"time"
)

type Token struct {
    AccessToken string    `json:"access_token"`
    TokenType   string    `json:"token_type"`
    Expiry      time.Time `json:"expiry"`
}

type TokenResponse struct {
    Result          Token `json:"result"`
    Success         bool  `json:"success"`
    ResultCode      int   `json:"resultCode"`
    ValidationErrors any   `json:"validationErrors"`
    Error           any   `json:"error"`
}

type JWTManager struct {
    clientID     string
    clientSecret string
    tokenURL     string
    token        *Token
    mu           sync.Mutex
}

func NewJWTManager(clientID, clientSecret, tokenURL string) *JWTManager {
    return &JWTManager{
        clientID:     clientID,
        clientSecret: clientSecret,
        tokenURL:     tokenURL,
    }
}

func (j *JWTManager) GetToken() (*Token, error) {
    j.mu.Lock()
    defer j.mu.Unlock()

    if j.token != nil && time.Now().Before(j.token.Expiry) {
        return j.token, nil
    }

    payload := map[string]string{
        "client_id":     j.clientID,
        "client_secret": j.clientSecret,
    }

    payloadBytes, err := json.Marshal(payload)
    if err != nil {
        return nil, err
    }

    req, err := http.NewRequest("POST", j.tokenURL, bytes.NewBuffer(payloadBytes))
    if err != nil {
        return nil, err
    }

    req.Header.Set("Content-Type", "application/json")

    client := &http.Client{}
    resp, err := client.Do(req)
    if err != nil {
        return nil, err
    }
    defer resp.Body.Close()

    if resp.StatusCode != http.StatusOK {
        body, _ := io.ReadAll(resp.Body)
        return nil, fmt.Errorf("failed to retrieve token, status: %d, body: %s", resp.StatusCode, string(body))
    }

    var tokenResponse TokenResponse
    if err := json.NewDecoder(resp.Body).Decode(&tokenResponse); err != nil {
        return nil, err
    }

    if !tokenResponse.Success {
        return nil, fmt.Errorf("token request unsuccessful, resultCode: %d", tokenResponse.ResultCode)
    }

    j.token = &tokenResponse.Result
    return &tokenResponse.Result, nil
}

func (j *JWTManager) GetAuthToken() (string, error) {
    token, err := j.GetToken()
    if err != nil {
        return "", err
    }
    return token.AccessToken, nil
}