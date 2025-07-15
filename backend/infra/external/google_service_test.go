package external

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
)

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
			testName:    "空のコードでもエラーにならない（mock実装）",
			code:        "",
			redirectURI: "http://localhost:3000/callback",
			expectError: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.testName, func(t *testing.T) {
			service := NewGoogleService("test_client_id", "test_client_secret")
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
			testName:    "空のトークンでもエラーにならない（mock実装）",
			accessToken: "",
			expectError: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.testName, func(t *testing.T) {
			service := NewGoogleService("test_client_id", "test_client_secret")
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