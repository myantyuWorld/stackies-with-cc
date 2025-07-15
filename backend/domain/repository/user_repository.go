package repository

import (
	"context"
	"stackies-backend/domain/model"
)

// UserRepository はユーザー関連のデータアクセスを抽象化する
type UserRepository interface {
	Create(ctx context.Context, user *model.User) error
	FindByEmail(ctx context.Context, email string) (*model.User, error)
	FindByID(ctx context.Context, id string) (*model.User, error)
	Update(ctx context.Context, user *model.User) error
}
