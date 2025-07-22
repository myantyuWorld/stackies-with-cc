package persistence

import (
	"context"
	"stackies-backend/domain/model"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestAuthRepositoryImpl_SaveToken(t *testing.T) {
	tests := []struct {
		testName    string
		userID      string
		token       *model.AuthToken
		expectError bool
	}{
		{
			testName: "正常なトークン保存",
			userID:   "user_123",
			token: &model.AuthToken{
				AccessToken:  "access_token",
				RefreshToken: "refresh_token",
				ExpiresIn:    time.Now().Add(time.Hour).Unix(),
				TokenType:    "Bearer",
			},
			expectError: false,
		},
		{
			testName: "空のユーザーIDでエラー",
			userID:   "",
			token: &model.AuthToken{
				AccessToken:  "access_token",
				RefreshToken: "refresh_token",
				ExpiresIn:    time.Now().Add(time.Hour).Unix(),
				TokenType:    "Bearer",
			},
			expectError: true,
		},
		{
			testName:    "nilトークンでエラー",
			userID:      "user_123",
			token:       nil,
			expectError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.testName, func(t *testing.T) {
			repo := NewAuthRepository()
			err := repo.SaveToken(context.Background(), tt.userID, tt.token)

			if tt.expectError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestAuthRepositoryImpl_GetToken(t *testing.T) {
	repo := NewAuthRepository()
	token := &model.AuthToken{
		AccessToken:  "access_token",
		RefreshToken: "refresh_token",
		ExpiresIn:    time.Now().Add(time.Hour).Unix(),
		TokenType:    "Bearer",
	}
	repo.SaveToken(context.Background(), "user_123", token)

	tests := []struct {
		testName    string
		userID      string
		expectError bool
		expectToken bool
	}{
		{
			testName:    "正常なトークン取得",
			userID:      "user_123",
			expectError: false,
			expectToken: true,
		},
		{
			testName:    "存在しないユーザーID",
			userID:      "notfound",
			expectError: true,
			expectToken: false,
		},
		{
			testName:    "空のユーザーIDでエラー",
			userID:      "",
			expectError: true,
			expectToken: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.testName, func(t *testing.T) {
			result, err := repo.GetToken(context.Background(), tt.userID)

			if tt.expectError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}

			if tt.expectToken {
				assert.NotNil(t, result)
				assert.Equal(t, token.AccessToken, result.AccessToken)
			} else {
				assert.Nil(t, result)
			}
		})
	}
}

func TestAuthRepositoryImpl_DeleteToken(t *testing.T) {
	repo := NewAuthRepository()
	token := &model.AuthToken{
		AccessToken:  "access_token",
		RefreshToken: "refresh_token",
		ExpiresIn:    time.Now().Add(time.Hour).Unix(),
		TokenType:    "Bearer",
	}
	repo.SaveToken(context.Background(), "user_123", token)

	tests := []struct {
		testName    string
		userID      string
		expectError bool
	}{
		{
			testName:    "正常なトークン削除",
			userID:      "user_123",
			expectError: false,
		},
		{
			testName:    "存在しないユーザーIDでも成功",
			userID:      "notfound",
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
			err := repo.DeleteToken(context.Background(), tt.userID)

			if tt.expectError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestAuthRepositoryImpl_ValidateToken(t *testing.T) {
	repo := NewAuthRepository()
	token := &model.AuthToken{
		AccessToken:  "valid_access_token",
		RefreshToken: "valid_refresh_token",
		ExpiresIn:    time.Now().Add(time.Hour).Unix(),
		TokenType:    "Bearer",
	}
	expiredToken := &model.AuthToken{
		AccessToken:  "expired_access_token",
		RefreshToken: "expired_refresh_token",
		ExpiresIn:    time.Now().Add(-time.Hour).Unix(),
		TokenType:    "Bearer",
	}
	repo.SaveToken(context.Background(), "user_123", token)
	repo.SaveToken(context.Background(), "user_456", expiredToken)

	tests := []struct {
		testName     string
		token        string
		expectError  bool
		expectedUser string
	}{
		{
			testName:     "正常なアクセストークン検証",
			token:        "valid_access_token",
			expectError:  false,
			expectedUser: "user_123",
		},
		{
			testName:     "正常なリフレッシュトークン検証",
			token:        "valid_refresh_token",
			expectError:  false,
			expectedUser: "user_123",
		},
		{
			testName:     "期限切れトークンでエラー",
			token:        "expired_access_token",
			expectError:  true,
			expectedUser: "",
		},
		{
			testName:     "存在しないトークンでエラー",
			token:        "invalid_token",
			expectError:  true,
			expectedUser: "",
		},
		{
			testName:     "空のトークンでエラー",
			token:        "",
			expectError:  true,
			expectedUser: "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.testName, func(t *testing.T) {
			result, err := repo.ValidateToken(context.Background(), tt.token)

			if tt.expectError {
				assert.Error(t, err)
				assert.Empty(t, result)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.expectedUser, result)
			}
		})
	}
}
