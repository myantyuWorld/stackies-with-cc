package model

import (
	"errors"
	"strings"
	"time"
)

type (
	// AuthToken は認証トークンを表す
	AuthToken struct {
		AccessToken  string `json:"access_token"`
		RefreshToken string `json:"refresh_token"`
		ExpiresIn    int64  `json:"expires_in"`
		TokenType    string `json:"token_type"`
	}

	// GoogleUserInfo はGoogleから取得するユーザー情報を表す
	GoogleUserInfo struct {
		ID            string `json:"id"`
		Email         string `json:"email"`
		VerifiedEmail bool   `json:"verified_email"`
		Name          string `json:"name"`
		GivenName     string `json:"given_name"`
		FamilyName    string `json:"family_name"`
		Picture       string `json:"picture"`
		Locale        string `json:"locale"`
	}
)

// NewAuthToken は新しい認証トークンを作成する
func NewAuthToken(accessToken, refreshToken string, expiresIn int64, tokenType string) (*AuthToken, error) {
	if strings.TrimSpace(accessToken) == "" {
		return nil, errors.New("access token cannot be empty")
	}
	if strings.TrimSpace(refreshToken) == "" {
		return nil, errors.New("refresh token cannot be empty")
	}
	if expiresIn <= 0 {
		return nil, errors.New("expires in must be positive")
	}

	return &AuthToken{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
		ExpiresIn:    expiresIn,
		TokenType:    tokenType,
	}, nil
}

// IsExpired はトークンが期限切れかどうかを確認する
func (a *AuthToken) IsExpired() bool {
	return time.Now().Unix() >= a.ExpiresIn
}

// ToUser はGoogleユーザー情報からユーザーエンティティを作成する
func (g *GoogleUserInfo) ToUser() *User {
	return &User{
		ID:        g.ID,
		Email:     g.Email,
		Name:      g.Name,
		Picture:   g.Picture,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}
}

// IsVerified はメールアドレスが認証済みかどうかを確認する
func (g *GoogleUserInfo) IsVerified() bool {
	return g.VerifiedEmail
}
