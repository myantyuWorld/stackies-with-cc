package persistence

import (
	"context"
	"errors"
	"stackies-backend/domain/model"
	"stackies-backend/domain/repository"
	"sync"
)

// UserRepositoryImpl はUserRepository interfaceの実装
// TODO: 実際のデータベース統合時にこのin-memory実装を置き換える
type UserRepositoryImpl struct {
	users map[string]*model.User
	mutex sync.RWMutex
}

// NewUserRepository は新しいUserRepositoryを作成する
func NewUserRepository() repository.UserRepository {
	return &UserRepositoryImpl{
		users: make(map[string]*model.User),
	}
}

// Save はユーザーを保存する
func (r *UserRepositoryImpl) Save(ctx context.Context, user *model.User) error {
	if user == nil {
		return errors.New("user cannot be nil")
	}
	if user.ID == "" {
		return errors.New("user ID cannot be empty")
	}

	r.mutex.Lock()
	defer r.mutex.Unlock()

	r.users[user.ID] = user
	return nil
}

// FindByEmail はメールアドレスでユーザーを検索する
func (r *UserRepositoryImpl) FindByEmail(ctx context.Context, email string) (*model.User, error) {
	if email == "" {
		return nil, errors.New("email cannot be empty")
	}

	r.mutex.RLock()
	defer r.mutex.RUnlock()

	for _, user := range r.users {
		if user.Email == email {
			return user, nil
		}
	}

	return nil, repository.ErrUserNotFound
}

// FindByID はIDでユーザーを検索する
func (r *UserRepositoryImpl) FindByID(ctx context.Context, id string) (*model.User, error) {
	if id == "" {
		return nil, errors.New("id cannot be empty")
	}

	r.mutex.RLock()
	defer r.mutex.RUnlock()

	user, exists := r.users[id]
	if !exists {
		return nil, repository.ErrUserNotFound
	}

	return user, nil
}

// Update はユーザー情報を更新する
func (r *UserRepositoryImpl) Update(ctx context.Context, user *model.User) error {
	if user == nil {
		return errors.New("user cannot be nil")
	}
	if user.ID == "" {
		return errors.New("user ID cannot be empty")
	}

	r.mutex.Lock()
	defer r.mutex.Unlock()

	if _, exists := r.users[user.ID]; !exists {
		return repository.ErrUserNotFound
	}

	r.users[user.ID] = user
	return nil
}
