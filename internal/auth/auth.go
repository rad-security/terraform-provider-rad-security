package auth

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"
)

type Authenticator struct {
	ApiURL     string
	httpClient *http.Client
}

func New(apiURL string) *Authenticator {
	return &Authenticator{
		ApiURL: apiURL,
		httpClient: &http.Client{
			Timeout: 15 * time.Second,
		},
	}
}

type AuthRequest struct {
	AccessKeyID string `json:"access_key_id"`
	SecretKey   string `json:"secret_key"`
}

type AuthResponse struct {
	Token string `json:"token"`
}

func (a *Authenticator) Authenticate(ctx context.Context, accessKeyID, secretKey string) (string, error) {
	authURL := fmt.Sprintf("%s/authentication/authenticate", a.ApiURL)

	reqBytes, err := json.Marshal(&AuthRequest{
		AccessKeyID: accessKeyID,
		SecretKey:   secretKey,
	})
	if err != nil {
		return "", fmt.Errorf("marshaling auth request: %w", err)
	}

	req, err := http.NewRequest(http.MethodPost, authURL, bytes.NewReader(reqBytes))
	if err != nil {
		return "", err
	}
	req = req.WithContext(ctx)

	resp, err := a.httpClient.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	bodyBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	if resp.StatusCode >= 300 {
		return "", fmt.Errorf("cannot auth: %s", string(bodyBytes))
	}

	var authResp AuthResponse
	if err = json.Unmarshal(bodyBytes, &authResp); err != nil {
		return "", fmt.Errorf("unmarshaling auth response: %w", err)
	}

	return authResp.Token, nil
}
