package services

import (
    "encoding/json"
    "errors"
    "net/http"
    "net/url"
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

    data := url.Values{}
    data.Set("grant_type", "client_credentials")
    data.Set("client_id", j.clientID)
    data.Set("client_secret", j.clientSecret)

    resp, err := http.PostForm(j.tokenURL, data)
    if err != nil {
        return nil, err
    }
    defer resp.Body.Close()

    if resp.StatusCode != http.StatusOK {
        return nil, errors.New("failed to retrieve token")
    }

    var token Token
    json.NewDecoder(resp.Body).Decode(&token)
    token.tokenExpiry = time.Now().Add(time.Duration(token.ExpiresIn) * time.Second)
    j.token = &token

    return &token, nil
}

func (j *JWTManager) GetAuthToken() (string, error) {
    token, err := j.GetToken()
    if err != nil {
        return "", err
    }
    return token.AccessToken, nil
}
