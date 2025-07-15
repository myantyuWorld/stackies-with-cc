package external

import (
	"errors"
	"stackies-backend/domain/service"
)

// JWTServiceImpl はJWTService interfaceの実装
type JWTServiceImpl struct {
	secretKey string
}

// NewJWTService は新しいJWTServiceを作成する
func NewJWTService(secretKey string) service.JWTService {
	return &JWTServiceImpl{
		secretKey: secretKey,
	}
}

// GenerateToken はJWTアクセストークンを生成する
func (j *JWTServiceImpl) GenerateToken(userID string) (string, error) {
	// TODO: 実際のJWTライブラリを使用してトークンを生成
	// ここでは仮の実装
	if userID == "" {
		return "", errors.New("userID cannot be empty")
	}
	return "mock_jwt_access_token_" + userID, nil
}

// ValidateToken はJWTトークンを検証してユーザーIDを返す
func (j *JWTServiceImpl) ValidateToken(token string) (string, error) {
	// TODO: 実際のJWTライブラリを使用してトークンを検証
	// ここでは仮の実装
	if token == "" {
		return "", errors.New("token cannot be empty")
	}
	if token == "invalid_token" {
		return "", errors.New("invalid token")
	}
	// mock実装では単純にuserIDを返す
	return "mock_user_id", nil
}

// GenerateRefreshToken はJWTリフレッシュトークンを生成する
func (j *JWTServiceImpl) GenerateRefreshToken(userID string) (string, error) {
	// TODO: 実際のJWTライブラリを使用してリフレッシュトークンを生成
	// ここでは仮の実装
	if userID == "" {
		return "", errors.New("userID cannot be empty")
	}
	return "mock_jwt_refresh_token_" + userID, nil
}
