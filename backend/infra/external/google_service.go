package external

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"stackies-backend/domain/model"
	"stackies-backend/domain/service"

	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
)

// GoogleServiceImpl はGoogleService interfaceの実装
type GoogleServiceImpl struct {
	config *oauth2.Config
}

// NewGoogleService は新しいGoogleServiceを作成する
func NewGoogleService(clientID, clientSecret string) service.GoogleService {
	config := &oauth2.Config{
		ClientID:     clientID,
		ClientSecret: clientSecret,
		RedirectURL:  "http://localhost:8080/auth/google/callback",
		Scopes: []string{
			"https://www.googleapis.com/auth/userinfo.email",
			"https://www.googleapis.com/auth/userinfo.profile",
		},
		Endpoint: google.Endpoint,
	}

	return &GoogleServiceImpl{
		config: config,
	}
}

// ExchangeCode は認証コードをアクセストークンに交換する
func (g *GoogleServiceImpl) ExchangeCode(ctx context.Context, code, redirectURI string) (*model.AuthToken, error) {
	token, err := g.config.Exchange(ctx, code)
	if err != nil {
		return nil, fmt.Errorf("failed to exchange code: %w", err)
	}

	return &model.AuthToken{
		AccessToken:  token.AccessToken,
		RefreshToken: token.RefreshToken,
		ExpiresIn:    int(token.Expiry.Unix()),
		TokenType:    token.TokenType,
	}, nil
}

// GetUserInfo はアクセストークンを使用してユーザー情報を取得する
func (g *GoogleServiceImpl) GetUserInfo(ctx context.Context, accessToken string) (*model.GoogleUserInfo, error) {
	client := &http.Client{}
	req, err := http.NewRequestWithContext(ctx, "GET", "https://www.googleapis.com/oauth2/v2/userinfo", nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("Authorization", "Bearer "+accessToken)

	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to get user info: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("google API returned status %d", resp.StatusCode)
	}

	var userInfo model.GoogleUserInfo
	if err := json.NewDecoder(resp.Body).Decode(&userInfo); err != nil {
		return nil, fmt.Errorf("failed to decode user info: %w", err)
	}

	return &userInfo, nil
}
