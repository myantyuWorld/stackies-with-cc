package repository

import (
	"context"
	"stackies-backend/domain/model"
)

// AuthRepository は認証関連のデータアクセスを抽象化する
type AuthRepository interface {
	SaveToken(ctx context.Context, userID string, token *model.AuthToken) error
	GetToken(ctx context.Context, userID string) (*model.AuthToken, error)
	DeleteToken(ctx context.Context, userID string) error
	ValidateToken(ctx context.Context, token string) (string, error)
}
