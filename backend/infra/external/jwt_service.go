package external

import (
	"errors"
	"stackies-backend/domain/service"
	"time"

	"github.com/golang-jwt/jwt"
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
	if userID == "" {
		return "", errors.New("userID cannot be empty")
	}

	claims := jwt.MapClaims{
		"user_id": userID,
		"type":    "access",
		"exp":     time.Now().Add(time.Hour * 24).Unix(),
		"iat":     time.Now().Unix(),
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString([]byte(j.secretKey))
	if err != nil {
		return "", err
	}

	return tokenString, nil
}

// ValidateToken はJWTトークンを検証してユーザーIDを返す
func (j *JWTServiceImpl) ValidateToken(token string) (string, error) {
	parsedToken, err := jwt.Parse(token, func(token *jwt.Token) (interface{}, error) {
		return []byte(j.secretKey), nil
	})
	if err != nil {
		return "", err
	}
	claims, ok := parsedToken.Claims.(jwt.MapClaims)
	if !ok {
		return "", errors.New("invalid token")
	}
	userID := claims["user_id"].(string)
	return userID, nil
}

// GenerateRefreshToken はJWTリフレッシュトークンを生成する
func (j *JWTServiceImpl) GenerateRefreshToken(userID string) (string, error) {
	if userID == "" {
		return "", errors.New("userID cannot be empty")
	}

	claims := jwt.MapClaims{
		"user_id": userID,
		"type":    "refresh",
		"exp":     time.Now().Add(time.Hour * 24 * 30).Unix(),
		"iat":     time.Now().Unix(),
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString([]byte(j.secretKey))
	if err != nil {
		return "", err
	}
	return tokenString, nil
}
