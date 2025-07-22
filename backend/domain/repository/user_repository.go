package repository

import (
	"context"
	"errors"
	"stackies-backend/domain/model"
)

var (
	ErrUserNotFound = errors.New("user not found")
)

// UserRepository はユーザー関連のデータアクセスを抽象化する
type UserRepository interface {
	Save(ctx context.Context, user *model.User) error
	FindByEmail(ctx context.Context, email string) (*model.User, error)
	FindByID(ctx context.Context, id string) (*model.User, error)
	Update(ctx context.Context, user *model.User) error
}
