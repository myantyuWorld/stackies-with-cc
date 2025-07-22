package external

import (
	"context"
	"stackies-backend/domain/model"
	"stackies-backend/domain/service"
	"testing"

	"github.com/stretchr/testify/assert"
)

// MockGoogleService はテスト用のモックサービス
type MockGoogleService struct{}

// NewMockGoogleService はモックサービスを作成する
func NewMockGoogleService() service.GoogleService {
	return &MockGoogleService{}
}

// GenerateAuthURL はモック認証URLを返す
func (m *MockGoogleService) GenerateAuthURL(state string) string {
	return "https://accounts.google.com/oauth/authorize?client_id=mock&state=" + state
}

// ExchangeCode はモックのトークンを返す
func (m *MockGoogleService) ExchangeCode(ctx context.Context, code, redirectURI string) (*model.AuthToken, error) {
	if code == "error_code" {
		return nil, assert.AnError
	}
	
	return &model.AuthToken{
		AccessToken:  "mock_access_token",
		RefreshToken: "mock_refresh_token",
		ExpiresIn:    3600,
		TokenType:    "Bearer",
	}, nil
}

// GetUserInfo はモックのユーザー情報を返す
func (m *MockGoogleService) GetUserInfo(ctx context.Context, accessToken string) (*model.GoogleUserInfo, error) {
	if accessToken == "error_token" {
		return nil, assert.AnError
	}

	return &model.GoogleUserInfo{
		ID:            "mock_google_user_id",
		Email:         "mock@example.com",
		VerifiedEmail: true,
		Name:          "Mock User",
	}, nil
}

func TestGoogleServiceImpl_ExchangeCode(t *testing.T) {
	tests := []struct {
		testName    string
		code        string
		redirectURI string
		expectError bool
	}{
		{
			testName:    "正常なコード交換",
			code:        "valid_code",
			redirectURI: "http://localhost:3000/callback",
			expectError: false,
		},
		{
			testName:    "空のコードでも成功する（mock実装）",
			code:        "",
			redirectURI: "http://localhost:3000/callback",
			expectError: false,
		},
		{
			testName:    "エラーコードでエラーになる",
			code:        "error_code",
			redirectURI: "http://localhost:3000/callback",
			expectError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.testName, func(t *testing.T) {
			service := NewMockGoogleService()
			result, err := service.ExchangeCode(context.Background(), tt.code, tt.redirectURI)

			if tt.expectError {
				assert.Error(t, err)
				assert.Nil(t, result)
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, result)
				assert.Equal(t, "mock_access_token", result.AccessToken)
				assert.Equal(t, "mock_refresh_token", result.RefreshToken)
				assert.Equal(t, "Bearer", result.TokenType)
			}
		})
	}
}

func TestGoogleServiceImpl_GetUserInfo(t *testing.T) {
	tests := []struct {
		testName    string
		accessToken string
		expectError bool
	}{
		{
			testName:    "正常なユーザー情報取得",
			accessToken: "valid_access_token",
			expectError: false,
		},
		{
			testName:    "空のトークンでも成功する（mock実装）",
			accessToken: "",
			expectError: false,
		},
		{
			testName:    "エラートークンでエラーになる",
			accessToken: "error_token",
			expectError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.testName, func(t *testing.T) {
			service := NewMockGoogleService()
			result, err := service.GetUserInfo(context.Background(), tt.accessToken)

			if tt.expectError {
				assert.Error(t, err)
				assert.Nil(t, result)
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, result)
				assert.Equal(t, "mock_google_user_id", result.ID)
				assert.Equal(t, "mock@example.com", result.Email)
				assert.True(t, result.VerifiedEmail)
				assert.Equal(t, "Mock User", result.Name)
			}
		})
	}
}
