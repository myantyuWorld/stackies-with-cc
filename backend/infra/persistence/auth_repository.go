package persistence

import (
	"context"
	"errors"
	"stackies-backend/domain/model"
	"stackies-backend/domain/repository"
	"sync"
)

// AuthRepositoryImpl はAuthRepository interfaceの実装
// TODO: 実際のRedis統合時にこのin-memory実装を置き換える
type AuthRepositoryImpl struct {
	tokens map[string]*model.AuthToken
	mutex  sync.RWMutex
}

// NewAuthRepository は新しいAuthRepositoryを作成する
func NewAuthRepository() repository.AuthRepository {
	return &AuthRepositoryImpl{
		tokens: make(map[string]*model.AuthToken),
	}
}

// SaveToken はトークンを保存する
func (r *AuthRepositoryImpl) SaveToken(ctx context.Context, userID string, token *model.AuthToken) error {
	if userID == "" {
		return errors.New("userID cannot be empty")
	}
	if token == nil {
		return errors.New("token cannot be nil")
	}

	r.mutex.Lock()
	defer r.mutex.Unlock()

	r.tokens[userID] = token
	return nil
}

// GetToken はユーザーIDでトークンを取得する
func (r *AuthRepositoryImpl) GetToken(ctx context.Context, userID string) (*model.AuthToken, error) {
	if userID == "" {
		return nil, errors.New("userID cannot be empty")
	}

	r.mutex.RLock()
	defer r.mutex.RUnlock()

	token, exists := r.tokens[userID]
	if !exists {
		return nil, errors.New("token not found")
	}

	return token, nil
}

// DeleteToken はトークンを削除する
func (r *AuthRepositoryImpl) DeleteToken(ctx context.Context, userID string) error {
	if userID == "" {
		return errors.New("userID cannot be empty")
	}

	r.mutex.Lock()
	defer r.mutex.Unlock()

	delete(r.tokens, userID)
	return nil
}

// ValidateToken はトークンを検証してユーザーIDを返す
func (r *AuthRepositoryImpl) ValidateToken(ctx context.Context, token string) (string, error) {
	if token == "" {
		return "", errors.New("token cannot be empty")
	}

	r.mutex.RLock()
	defer r.mutex.RUnlock()

	// in-memory実装では、すべてのトークンを検索してマッチするものを探す
	for userID, storedToken := range r.tokens {
		if storedToken.AccessToken == token || storedToken.RefreshToken == token {
			// トークンの有効期限チェック
			if storedToken.IsExpired() {
				return "", errors.New("token expired")
			}
			return userID, nil
		}
	}

	return "", errors.New("invalid token")
}