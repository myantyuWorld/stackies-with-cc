package external

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestJWTServiceImpl_GenerateToken(t *testing.T) {
	tests := []struct {
		testName    string
		userID      string
		expectError bool
	}{
		{
			testName:    "正常なトークン生成",
			userID:      "user_123",
			expectError: false,
		},
		{
			testName:    "空のユーザーIDでエラー",
			userID:      "",
			expectError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.testName, func(t *testing.T) {
			service := NewJWTService("test_secret_key")
			result, err := service.GenerateToken(tt.userID)

			if tt.expectError {
				assert.Error(t, err)
				assert.Empty(t, result)
			} else {
				assert.NoError(t, err)
				assert.NotEmpty(t, result)
				assert.Contains(t, result, "mock_jwt_access_token_")
				assert.Contains(t, result, tt.userID)
			}
		})
	}
}

func TestJWTServiceImpl_ValidateToken(t *testing.T) {
	tests := []struct {
		testName    string
		token       string
		expectError bool
		expectedUID string
	}{
		{
			testName:    "正常なトークン検証",
			token:       "valid_token",
			expectError: false,
			expectedUID: "mock_user_id",
		},
		{
			testName:    "無効なトークンでエラー",
			token:       "invalid_token",
			expectError: true,
			expectedUID: "",
		},
		{
			testName:    "空のトークンでエラー",
			token:       "",
			expectError: true,
			expectedUID: "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.testName, func(t *testing.T) {
			service := NewJWTService("test_secret_key")
			result, err := service.ValidateToken(tt.token)

			if tt.expectError {
				assert.Error(t, err)
				assert.Empty(t, result)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.expectedUID, result)
			}
		})
	}
}

func TestJWTServiceImpl_GenerateRefreshToken(t *testing.T) {
	tests := []struct {
		testName    string
		userID      string
		expectError bool
	}{
		{
			testName:    "正常なリフレッシュトークン生成",
			userID:      "user_123",
			expectError: false,
		},
		{
			testName:    "空のユーザーIDでエラー",
			userID:      "",
			expectError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.testName, func(t *testing.T) {
			service := NewJWTService("test_secret_key")
			result, err := service.GenerateRefreshToken(tt.userID)

			if tt.expectError {
				assert.Error(t, err)
				assert.Empty(t, result)
			} else {
				assert.NoError(t, err)
				assert.NotEmpty(t, result)
				assert.Contains(t, result, "mock_jwt_refresh_token_")
				assert.Contains(t, result, tt.userID)
			}
		})
	}
}