package model

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestAuthToken_NewAuthToken(t *testing.T) {
	tests := []struct {
		testName     string
		accessToken  string
		refreshToken string
		expiresIn    int64
		tokenType    string
		wantErr      bool
	}{
		{
			testName:     "正常なトークン作成",
			accessToken:  "access_token_123",
			refreshToken: "refresh_token_456",
			expiresIn:    3600,
			tokenType:    "Bearer",
			wantErr:      false,
		},
		{
			testName:     "空のアクセストークンでエラー",
			accessToken:  "",
			refreshToken: "refresh_token_456",
			expiresIn:    3600,
			tokenType:    "Bearer",
			wantErr:      true,
		},
		{
			testName:     "空のリフレッシュトークンでエラー",
			accessToken:  "access_token_123",
			refreshToken: "",
			expiresIn:    3600,
			tokenType:    "Bearer",
			wantErr:      true,
		},
		{
			testName:     "無効な有効期限でエラー",
			accessToken:  "access_token_123",
			refreshToken: "refresh_token_456",
			expiresIn:    0,
			tokenType:    "Bearer",
			wantErr:      true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.testName, func(t *testing.T) {
			got, err := NewAuthToken(tt.accessToken, tt.refreshToken, tt.expiresIn, tt.tokenType)
			if tt.wantErr {
				assert.Error(t, err)
				assert.Nil(t, got)
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, got)
				assert.Equal(t, tt.accessToken, got.AccessToken)
				assert.Equal(t, tt.refreshToken, got.RefreshToken)
				assert.Equal(t, tt.expiresIn, got.ExpiresIn)
				assert.Equal(t, tt.tokenType, got.TokenType)
			}
		})
	}
}

func TestAuthToken_IsExpired(t *testing.T) {
	tests := []struct {
		testName  string
		expiresIn int64
		want      bool
	}{
		{
			testName:  "まだ有効",
			expiresIn: time.Now().Add(time.Hour).Unix(),
			want:      false,
		},
		{
			testName:  "期限切れ",
			expiresIn: time.Now().Add(-time.Hour).Unix(),
			want:      true,
		},
		{
			testName:  "ちょうど期限",
			expiresIn: time.Now().Unix(),
			want:      true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.testName, func(t *testing.T) {
			token := &AuthToken{
				AccessToken:  "test_token",
				RefreshToken: "test_refresh",
				ExpiresIn:    tt.expiresIn,
				TokenType:    "Bearer",
			}
			got := token.IsExpired()
			assert.Equal(t, tt.want, got)
		})
	}
}

func TestGoogleUserInfo_ToUser(t *testing.T) {
	googleUser := &GoogleUserInfo{
		ID:            "google_123",
		Email:         "test@example.com",
		VerifiedEmail: true,
		Name:          "Test User",
		GivenName:     "Test",
		FamilyName:    "User",
		Picture:       "https://example.com/picture.jpg",
		Locale:        "ja",
	}

	user := googleUser.ToUser()

	assert.Equal(t, googleUser.ID, user.ID)
	assert.Equal(t, googleUser.Email, user.Email)
	assert.Equal(t, googleUser.Name, user.Name)
	assert.Equal(t, googleUser.Picture, user.Picture)
	assert.WithinDuration(t, time.Now(), user.CreatedAt, time.Second)
	assert.WithinDuration(t, time.Now(), user.UpdatedAt, time.Second)
}

func TestGoogleUserInfo_IsVerified(t *testing.T) {
	tests := []struct {
		testName string
		verified bool
		want     bool
	}{
		{
			testName: "認証済み",
			verified: true,
			want:     true,
		},
		{
			testName: "未認証",
			verified: false,
			want:     false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.testName, func(t *testing.T) {
			googleUser := &GoogleUserInfo{
				VerifiedEmail: tt.verified,
			}
			got := googleUser.IsVerified()
			assert.Equal(t, tt.want, got)
		})
	}
}
