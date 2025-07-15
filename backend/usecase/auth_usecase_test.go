package usecase

import (
	"context"
	"errors"
	"stackies-backend/domain/model"
	"stackies-backend/domain/repository"
	"stackies-backend/domain/service"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// MockUserRepository はUserRepositoryのモック
type MockUserRepository struct {
	mock.Mock
}

var _ repository.UserRepository = (*MockUserRepository)(nil)

func (m *MockUserRepository) Save(ctx context.Context, user *model.User) error {
	args := m.Called(ctx, user)
	return args.Error(0)
}

func (m *MockUserRepository) FindByID(ctx context.Context, id string) (*model.User, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*model.User), args.Error(1)
}

func (m *MockUserRepository) FindByEmail(ctx context.Context, email string) (*model.User, error) {
	args := m.Called(ctx, email)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*model.User), args.Error(1)
}

func (m *MockUserRepository) Update(ctx context.Context, user *model.User) error {
	args := m.Called(ctx, user)
	return args.Error(0)
}

// MockAuthRepository はAuthRepositoryのモック
type MockAuthRepository struct {
	mock.Mock
}

var _ repository.AuthRepository = (*MockAuthRepository)(nil)

func (m *MockAuthRepository) SaveToken(ctx context.Context, userID string, token *model.AuthToken) error {
	args := m.Called(ctx, userID, token)
	return args.Error(0)
}

func (m *MockAuthRepository) GetToken(ctx context.Context, userID string) (*model.AuthToken, error) {
	args := m.Called(ctx, userID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*model.AuthToken), args.Error(1)
}

func (m *MockAuthRepository) DeleteToken(ctx context.Context, userID string) error {
	args := m.Called(ctx, userID)
	return args.Error(0)
}

func (m *MockAuthRepository) ValidateToken(ctx context.Context, token string) (string, error) {
	args := m.Called(ctx, token)
	return args.String(0), args.Error(1)
}

// MockGoogleService はGoogleサービスのモック
type MockGoogleService struct {
	mock.Mock
}

var _ service.GoogleService = (*MockGoogleService)(nil)

func (m *MockGoogleService) ExchangeCode(ctx context.Context, code, redirectURI string) (*model.AuthToken, error) {
	args := m.Called(ctx, code, redirectURI)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*model.AuthToken), args.Error(1)
}

func (m *MockGoogleService) GetUserInfo(ctx context.Context, accessToken string) (*model.GoogleUserInfo, error) {
	args := m.Called(ctx, accessToken)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*model.GoogleUserInfo), args.Error(1)
}

// MockJWTService はJWTサービスのモック
type MockJWTService struct {
	mock.Mock
}

var _ service.JWTService = (*MockJWTService)(nil)

func (m *MockJWTService) GenerateToken(userID string) (string, error) {
	args := m.Called(userID)
	return args.String(0), args.Error(1)
}

func (m *MockJWTService) ValidateToken(token string) (string, error) {
	args := m.Called(token)
	return args.String(0), args.Error(1)
}

func (m *MockJWTService) GenerateRefreshToken(userID string) (string, error) {
	args := m.Called(userID)
	return args.String(0), args.Error(1)
}

func TestAuthUsecaseImpl_GoogleLogin(t *testing.T) {
	tests := []struct {
		testName    string
		input       *GoogleLoginInput
		setupMocks  func(*MockUserRepository, *MockAuthRepository, *MockGoogleService, *MockJWTService)
		expectError bool
		expectUser  bool
	}{
		{
			testName: "正常なGoogleログイン - 新規ユーザー",
			input: &GoogleLoginInput{
				AuthorizationCode: "test_code",
				RedirectURI:       "http://localhost:3000/callback",
			},
			setupMocks: func(userRepo *MockUserRepository, authRepo *MockAuthRepository, googleSvc *MockGoogleService, jwtSvc *MockJWTService) {
				googleToken := &model.AuthToken{
					AccessToken:  "google_access_token",
					RefreshToken: "google_refresh_token",
					ExpiresIn:    time.Now().Add(time.Hour).Unix(),
					TokenType:    "Bearer",
				}
				googleUser := &model.GoogleUserInfo{
					ID:            "google_123",
					Email:         "test@example.com",
					VerifiedEmail: true,
					Name:          "Test User",
					Picture:       "https://example.com/picture.jpg",
				}

				googleSvc.On("ExchangeCode", mock.Anything, "test_code", "http://localhost:3000/callback").Return(googleToken, nil)
				googleSvc.On("GetUserInfo", mock.Anything, "google_access_token").Return(googleUser, nil)
				userRepo.On("FindByEmail", mock.Anything, "test@example.com").Return((*model.User)(nil), repository.ErrUserNotFound)
				userRepo.On("Save", mock.Anything, mock.AnythingOfType("*model.User")).Return(nil)
				jwtSvc.On("GenerateToken", "google_123").Return("jwt_access_token", nil)
				jwtSvc.On("GenerateRefreshToken", "google_123").Return("jwt_refresh_token", nil)
				authRepo.On("SaveToken", mock.Anything, "google_123", mock.AnythingOfType("*model.AuthToken")).Return(nil)
			},
			expectError: false,
			expectUser:  true,
		},
		{
			testName: "正常なGoogleログイン - 既存ユーザー",
			input: &GoogleLoginInput{
				AuthorizationCode: "test_code",
				RedirectURI:       "http://localhost:3000/callback",
			},
			setupMocks: func(userRepo *MockUserRepository, authRepo *MockAuthRepository, googleSvc *MockGoogleService, jwtSvc *MockJWTService) {
				googleToken := &model.AuthToken{
					AccessToken:  "google_access_token",
					RefreshToken: "google_refresh_token",
					ExpiresIn:    time.Now().Add(time.Hour).Unix(),
					TokenType:    "Bearer",
				}
				googleUser := &model.GoogleUserInfo{
					ID:            "google_123",
					Email:         "test@example.com",
					VerifiedEmail: true,
					Name:          "Test User Updated",
					Picture:       "https://example.com/new-picture.jpg",
				}
				existingUser := &model.User{
					ID:        "google_123",
					Email:     "test@example.com",
					Name:      "Test User",
					Picture:   "https://example.com/picture.jpg",
					CreatedAt: time.Now().Add(-time.Hour),
					UpdatedAt: time.Now().Add(-time.Hour),
				}

				googleSvc.On("ExchangeCode", mock.Anything, "test_code", "http://localhost:3000/callback").Return(googleToken, nil)
				googleSvc.On("GetUserInfo", mock.Anything, "google_access_token").Return(googleUser, nil)
				userRepo.On("FindByEmail", mock.Anything, "test@example.com").Return(existingUser, nil)
				userRepo.On("Update", mock.Anything, mock.AnythingOfType("*model.User")).Return(nil)
				jwtSvc.On("GenerateToken", "google_123").Return("jwt_access_token", nil)
				jwtSvc.On("GenerateRefreshToken", "google_123").Return("jwt_refresh_token", nil)
				authRepo.On("SaveToken", mock.Anything, "google_123", mock.AnythingOfType("*model.AuthToken")).Return(nil)
			},
			expectError: false,
			expectUser:  true,
		},
		{
			testName: "Google認証コード交換エラー",
			input: &GoogleLoginInput{
				AuthorizationCode: "invalid_code",
				RedirectURI:       "http://localhost:3000/callback",
			},
			setupMocks: func(userRepo *MockUserRepository, authRepo *MockAuthRepository, googleSvc *MockGoogleService, jwtSvc *MockJWTService) {
				googleSvc.On("ExchangeCode", mock.Anything, "invalid_code", "http://localhost:3000/callback").Return((*model.AuthToken)(nil), errors.New("invalid code"))
			},
			expectError: true,
			expectUser:  false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.testName, func(t *testing.T) {
			userRepo := new(MockUserRepository)
			authRepo := new(MockAuthRepository)
			googleSvc := new(MockGoogleService)
			jwtSvc := new(MockJWTService)

			tt.setupMocks(userRepo, authRepo, googleSvc, jwtSvc)

			usecase := NewAuthUsecase(userRepo, authRepo, googleSvc, jwtSvc)
			result, err := usecase.GoogleLogin(context.Background(), tt.input)

			if tt.expectError {
				assert.Error(t, err)
				assert.Nil(t, result)
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, result)
				if tt.expectUser {
					assert.NotNil(t, result.User)
					assert.NotEmpty(t, result.AccessToken)
					assert.NotEmpty(t, result.RefreshToken)
					assert.Greater(t, result.ExpiresIn, int64(0))
				}
			}

			userRepo.AssertExpectations(t)
			authRepo.AssertExpectations(t)
			googleSvc.AssertExpectations(t)
			jwtSvc.AssertExpectations(t)
		})
	}
}

func TestAuthUsecaseImpl_RefreshToken(t *testing.T) {
	tests := []struct {
		testName    string
		input       *RefreshTokenInput
		setupMocks  func(*MockUserRepository, *MockAuthRepository, *MockGoogleService, *MockJWTService)
		expectError bool
	}{
		{
			testName: "正常なトークンリフレッシュ",
			input: &RefreshTokenInput{
				RefreshToken: "valid_refresh_token",
			},
			setupMocks: func(userRepo *MockUserRepository, authRepo *MockAuthRepository, googleSvc *MockGoogleService, jwtSvc *MockJWTService) {
				jwtSvc.On("ValidateToken", "valid_refresh_token").Return("user_123", nil)
				jwtSvc.On("GenerateToken", "user_123").Return("new_access_token", nil)
				jwtSvc.On("GenerateRefreshToken", "user_123").Return("new_refresh_token", nil)
				authRepo.On("SaveToken", mock.Anything, "user_123", mock.AnythingOfType("*model.AuthToken")).Return(nil)
			},
			expectError: false,
		},
		{
			testName: "無効なリフレッシュトークン",
			input: &RefreshTokenInput{
				RefreshToken: "invalid_refresh_token",
			},
			setupMocks: func(userRepo *MockUserRepository, authRepo *MockAuthRepository, googleSvc *MockGoogleService, jwtSvc *MockJWTService) {
				jwtSvc.On("ValidateToken", "invalid_refresh_token").Return("", errors.New("invalid token"))
			},
			expectError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.testName, func(t *testing.T) {
			userRepo := new(MockUserRepository)
			authRepo := new(MockAuthRepository)
			googleSvc := new(MockGoogleService)
			jwtSvc := new(MockJWTService)

			tt.setupMocks(userRepo, authRepo, googleSvc, jwtSvc)

			usecase := NewAuthUsecase(userRepo, authRepo, googleSvc, jwtSvc)
			result, err := usecase.RefreshToken(context.Background(), tt.input)

			if tt.expectError {
				assert.Error(t, err)
				assert.Nil(t, result)
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, result)
				assert.NotEmpty(t, result.AccessToken)
				assert.NotEmpty(t, result.RefreshToken)
			}

			userRepo.AssertExpectations(t)
			authRepo.AssertExpectations(t)
			googleSvc.AssertExpectations(t)
			jwtSvc.AssertExpectations(t)
		})
	}
}

func TestAuthUsecaseImpl_Logout(t *testing.T) {
	tests := []struct {
		testName    string
		input       *LogoutInput
		setupMocks  func(*MockUserRepository, *MockAuthRepository, *MockGoogleService, *MockJWTService)
		expectError bool
	}{
		{
			testName: "正常なログアウト",
			input: &LogoutInput{
				UserID: "user_123",
			},
			setupMocks: func(userRepo *MockUserRepository, authRepo *MockAuthRepository, googleSvc *MockGoogleService, jwtSvc *MockJWTService) {
				authRepo.On("DeleteToken", mock.Anything, "user_123").Return(nil)
			},
			expectError: false,
		},
		{
			testName: "トークン削除エラー",
			input: &LogoutInput{
				UserID: "user_123",
			},
			setupMocks: func(userRepo *MockUserRepository, authRepo *MockAuthRepository, googleSvc *MockGoogleService, jwtSvc *MockJWTService) {
				authRepo.On("DeleteToken", mock.Anything, "user_123").Return(errors.New("delete error"))
			},
			expectError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.testName, func(t *testing.T) {
			userRepo := new(MockUserRepository)
			authRepo := new(MockAuthRepository)
			googleSvc := new(MockGoogleService)
			jwtSvc := new(MockJWTService)

			tt.setupMocks(userRepo, authRepo, googleSvc, jwtSvc)

			usecase := NewAuthUsecase(userRepo, authRepo, googleSvc, jwtSvc)
			err := usecase.Logout(context.Background(), tt.input)

			if tt.expectError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}

			userRepo.AssertExpectations(t)
			authRepo.AssertExpectations(t)
			googleSvc.AssertExpectations(t)
			jwtSvc.AssertExpectations(t)
		})
	}
}