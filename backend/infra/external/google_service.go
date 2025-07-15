package external

import (
	"context"
	"stackies-backend/domain/model"
	"stackies-backend/domain/service"
)

// GoogleServiceImpl はGoogleService interfaceの実装
type GoogleServiceImpl struct {
	// 実際のGoogle OAuth2.0ライブラリとの統合ポイント
	clientID     string
	clientSecret string
}

// NewGoogleService は新しいGoogleServiceを作成する
func NewGoogleService(clientID, clientSecret string) service.GoogleService {
	return &GoogleServiceImpl{
		clientID:     clientID,
		clientSecret: clientSecret,
	}
}

// ExchangeCode は認証コードをアクセストークンに交換する
func (g *GoogleServiceImpl) ExchangeCode(ctx context.Context, code, redirectURI string) (*model.AuthToken, error) {
	// TODO: 実際のGoogle OAuth2.0 APIとの統合実装
	// ここでは仮の実装
	return &model.AuthToken{
		AccessToken:  "mock_access_token",
		RefreshToken: "mock_refresh_token",
		ExpiresIn:    3600,
		TokenType:    "Bearer",
	}, nil
}

// GetUserInfo はアクセストークンを使用してユーザー情報を取得する
func (g *GoogleServiceImpl) GetUserInfo(ctx context.Context, accessToken string) (*model.GoogleUserInfo, error) {
	// TODO: 実際のGoogle User Info APIとの統合実装
	// ここでは仮の実装
	return &model.GoogleUserInfo{
		ID:            "mock_google_user_id",
		Email:         "mock@example.com",
		VerifiedEmail: true,
		Name:          "Mock User",
		GivenName:     "Mock",
		FamilyName:    "User",
		Picture:       "https://example.com/picture.jpg",
		Locale:        "ja",
	}, nil
}
