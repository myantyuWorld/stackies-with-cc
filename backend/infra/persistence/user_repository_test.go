package persistence

import (
	"context"
	"stackies-backend/domain/model"
	"stackies-backend/domain/repository"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestUserRepositoryImpl_Save(t *testing.T) {
	tests := []struct {
		testName    string
		user        *model.User
		expectError bool
	}{
		{
			testName: "正常なユーザー保存",
			user: &model.User{
				ID:        "user_123",
				Email:     "test@example.com",
				Name:      "Test User",
				Picture:   "https://example.com/picture.jpg",
				CreatedAt: time.Now(),
				UpdatedAt: time.Now(),
			},
			expectError: false,
		},
		{
			testName:    "nilユーザーでエラー",
			user:        nil,
			expectError: true,
		},
		{
			testName: "空のIDでエラー",
			user: &model.User{
				ID:    "",
				Email: "test@example.com",
				Name:  "Test User",
			},
			expectError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.testName, func(t *testing.T) {
			repo := NewUserRepository()
			err := repo.Save(context.Background(), tt.user)

			if tt.expectError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestUserRepositoryImpl_FindByEmail(t *testing.T) {
	repo := NewUserRepository()
	user := &model.User{
		ID:        "user_123",
		Email:     "test@example.com",
		Name:      "Test User",
		Picture:   "https://example.com/picture.jpg",
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}
	repo.Save(context.Background(), user)

	tests := []struct {
		testName    string
		email       string
		expectError bool
		expectUser  bool
	}{
		{
			testName:    "正常なメール検索",
			email:       "test@example.com",
			expectError: false,
			expectUser:  true,
		},
		{
			testName:    "存在しないメール",
			email:       "notfound@example.com",
			expectError: true,
			expectUser:  false,
		},
		{
			testName:    "空のメールでエラー",
			email:       "",
			expectError: true,
			expectUser:  false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.testName, func(t *testing.T) {
			result, err := repo.FindByEmail(context.Background(), tt.email)

			if tt.expectError {
				assert.Error(t, err)
				if err == repository.ErrUserNotFound {
					assert.Equal(t, repository.ErrUserNotFound, err)
				}
			} else {
				assert.NoError(t, err)
			}

			if tt.expectUser {
				assert.NotNil(t, result)
				assert.Equal(t, user.ID, result.ID)
				assert.Equal(t, user.Email, result.Email)
			} else {
				assert.Nil(t, result)
			}
		})
	}
}

func TestUserRepositoryImpl_FindByID(t *testing.T) {
	repo := NewUserRepository()
	user := &model.User{
		ID:        "user_123",
		Email:     "test@example.com",
		Name:      "Test User",
		Picture:   "https://example.com/picture.jpg",
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}
	repo.Save(context.Background(), user)

	tests := []struct {
		testName    string
		id          string
		expectError bool
		expectUser  bool
	}{
		{
			testName:    "正常なID検索",
			id:          "user_123",
			expectError: false,
			expectUser:  true,
		},
		{
			testName:    "存在しないID",
			id:          "notfound",
			expectError: true,
			expectUser:  false,
		},
		{
			testName:    "空のIDでエラー",
			id:          "",
			expectError: true,
			expectUser:  false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.testName, func(t *testing.T) {
			result, err := repo.FindByID(context.Background(), tt.id)

			if tt.expectError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}

			if tt.expectUser {
				assert.NotNil(t, result)
				assert.Equal(t, user.ID, result.ID)
			} else {
				assert.Nil(t, result)
			}
		})
	}
}

func TestUserRepositoryImpl_Update(t *testing.T) {
	repo := NewUserRepository()
	user := &model.User{
		ID:        "user_123",
		Email:     "test@example.com",
		Name:      "Test User",
		Picture:   "https://example.com/picture.jpg",
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}
	repo.Save(context.Background(), user)

	tests := []struct {
		testName    string
		user        *model.User
		expectError bool
	}{
		{
			testName: "正常なユーザー更新",
			user: &model.User{
				ID:        "user_123",
				Email:     "test@example.com",
				Name:      "Updated User",
				Picture:   "https://example.com/new-picture.jpg",
				CreatedAt: user.CreatedAt,
				UpdatedAt: time.Now(),
			},
			expectError: false,
		},
		{
			testName:    "nilユーザーでエラー",
			user:        nil,
			expectError: true,
		},
		{
			testName: "存在しないユーザーでエラー",
			user: &model.User{
				ID:    "notfound",
				Email: "test@example.com",
				Name:  "Test User",
			},
			expectError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.testName, func(t *testing.T) {
			err := repo.Update(context.Background(), tt.user)

			if tt.expectError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				// 更新されたかチェック
				updated, _ := repo.FindByID(context.Background(), tt.user.ID)
				assert.Equal(t, tt.user.Name, updated.Name)
			}
		})
	}
}