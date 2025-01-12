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
    AccessToken string `json:"access_token"`
    ExpiresIn   int    `json:"expires_in"`
    tokenExpiry time.Time
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

    if j.token != nil && time.Now().Before(j.token.tokenExpiry) {
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
        body, _ := io.ReadAll(resp.Body) // for debugging purposes
        return nil, fmt.Errorf("failed to retrieve token, status: %d, body: %s", resp.StatusCode, string(body))
    }

    var token Token
    if err := json.NewDecoder(resp.Body).Decode(&token); err != nil {
        return nil, err
    }

    token.tokenExpiry = time.Now().Add(time.Duration(token.ExpiresIn) * time.Second)
    j.token = &token

    return &token, nil
}
