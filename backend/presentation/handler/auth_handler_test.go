package handler

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"stackies-backend/domain/model"
	"stackies-backend/domain/repository"
	"stackies-backend/usecase"
	"testing"
	"time"

	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// MockAuthUsecase はAuthUsecaseのモック
type MockAuthUsecase struct {
	mock.Mock
}

func (m *MockAuthUsecase) GoogleLogin(ctx context.Context, input *usecase.GoogleLoginInput) (*usecase.GoogleLoginOutput, error) {
	args := m.Called(ctx, input)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*usecase.GoogleLoginOutput), args.Error(1)
}

func (m *MockAuthUsecase) RefreshToken(ctx context.Context, input *usecase.RefreshTokenInput) (*usecase.RefreshTokenOutput, error) {
	args := m.Called(ctx, input)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*usecase.RefreshTokenOutput), args.Error(1)
}

func (m *MockAuthUsecase) Logout(ctx context.Context, input *usecase.LogoutInput) error {
	args := m.Called(ctx, input)
	return args.Error(0)
}

// MockUserRepository はUserRepositoryのモック
type MockUserRepository struct {
	mock.Mock
}

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

func TestAuthHandler_GoogleLogin(t *testing.T) {
	tests := []struct {
		testName       string
		requestBody    interface{}
		setupMocks     func(*MockAuthUsecase, *MockUserRepository)
		expectedStatus int
		expectError    bool
	}{
		{
			testName: "正常なGoogleログイン",
			requestBody: GoogleLoginRequest{
				AuthorizationCode: "valid_code",
				RedirectURI:       "http://localhost:3000/callback",
			},
			setupMocks: func(authUC *MockAuthUsecase, userRepo *MockUserRepository) {
				user := &model.User{
					ID:    "user_123",
					Email: "test@example.com",
					Name:  "Test User",
				}
				output := &usecase.GoogleLoginOutput{
					User:         user,
					AccessToken:  "access_token",
					RefreshToken: "refresh_token",
					ExpiresIn:    3600,
				}
				authUC.On("GoogleLogin", mock.Anything, mock.AnythingOfType("*usecase.GoogleLoginInput")).Return(output, nil)
			},
			expectedStatus: http.StatusOK,
			expectError:    false,
		},
		{
			testName: "無効なリクエストボディ",
			requestBody: map[string]interface{}{
				"invalid": "data",
			},
			setupMocks: func(authUC *MockAuthUsecase, userRepo *MockUserRepository) {
				authUC.On("GoogleLogin", mock.Anything, mock.AnythingOfType("*usecase.GoogleLoginInput")).Return(nil, errors.New("auth error"))
			},
			expectedStatus: http.StatusInternalServerError,
			expectError:    true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.testName, func(t *testing.T) {
			authUC := new(MockAuthUsecase)
			userRepo := new(MockUserRepository)
			tt.setupMocks(authUC, userRepo)

			handler := NewAuthHandler(authUC, userRepo)

			e := echo.New()
			body, _ := json.Marshal(tt.requestBody)
			req := httptest.NewRequest(http.MethodPost, "/auth/google/login", bytes.NewReader(body))
			req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
			rec := httptest.NewRecorder()
			c := e.NewContext(req, rec)

			err := handler.GoogleLogin(c)

			if tt.expectError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.expectedStatus, rec.Code)

				if rec.Code == http.StatusOK {
					var response GoogleLoginResponse
					err := json.Unmarshal(rec.Body.Bytes(), &response)
					assert.NoError(t, err)
					assert.NotEmpty(t, response.AccessToken)
					assert.NotEmpty(t, response.RefreshToken)
				}
			}

			authUC.AssertExpectations(t)
			userRepo.AssertExpectations(t)
		})
	}
}

func TestAuthHandler_RefreshToken(t *testing.T) {
	tests := []struct {
		testName       string
		requestBody    RefreshTokenRequest
		setupMocks     func(*MockAuthUsecase, *MockUserRepository)
		expectedStatus int
		expectError    bool
	}{
		{
			testName: "正常なトークンリフレッシュ",
			requestBody: RefreshTokenRequest{
				RefreshToken: "valid_refresh_token",
			},
			setupMocks: func(authUC *MockAuthUsecase, userRepo *MockUserRepository) {
				output := &usecase.RefreshTokenOutput{
					AccessToken:  "new_access_token",
					RefreshToken: "new_refresh_token",
					ExpiresIn:    3600,
				}
				authUC.On("RefreshToken", mock.Anything, mock.AnythingOfType("*usecase.RefreshTokenInput")).Return(output, nil)
			},
			expectedStatus: http.StatusOK,
			expectError:    false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.testName, func(t *testing.T) {
			authUC := new(MockAuthUsecase)
			userRepo := new(MockUserRepository)
			tt.setupMocks(authUC, userRepo)

			handler := NewAuthHandler(authUC, userRepo)

			e := echo.New()
			body, _ := json.Marshal(tt.requestBody)
			req := httptest.NewRequest(http.MethodPost, "/auth/refresh", bytes.NewReader(body))
			req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
			rec := httptest.NewRecorder()
			c := e.NewContext(req, rec)

			err := handler.RefreshToken(c)

			if tt.expectError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.expectedStatus, rec.Code)
			}

			authUC.AssertExpectations(t)
			userRepo.AssertExpectations(t)
		})
	}
}

func TestAuthHandler_Logout(t *testing.T) {
	tests := []struct {
		testName       string
		userID         string
		setupMocks     func(*MockAuthUsecase, *MockUserRepository)
		expectedStatus int
		expectError    bool
	}{
		{
			testName: "正常なログアウト",
			userID:   "user_123",
			setupMocks: func(authUC *MockAuthUsecase, userRepo *MockUserRepository) {
				authUC.On("Logout", mock.Anything, mock.AnythingOfType("*usecase.LogoutInput")).Return(nil)
			},
			expectedStatus: http.StatusOK,
			expectError:    false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.testName, func(t *testing.T) {
			authUC := new(MockAuthUsecase)
			userRepo := new(MockUserRepository)
			tt.setupMocks(authUC, userRepo)

			handler := NewAuthHandler(authUC, userRepo)

			e := echo.New()
			req := httptest.NewRequest(http.MethodPost, "/auth/logout", nil)
			rec := httptest.NewRecorder()
			c := e.NewContext(req, rec)
			c.Set("user_id", tt.userID)

			err := handler.Logout(c)

			if tt.expectError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.expectedStatus, rec.Code)
			}

			authUC.AssertExpectations(t)
			userRepo.AssertExpectations(t)
		})
	}
}

func TestAuthHandler_GetMe(t *testing.T) {
	tests := []struct {
		testName       string
		userID         string
		setupMocks     func(*MockAuthUsecase, *MockUserRepository)
		expectedStatus int
		expectError    bool
	}{
		{
			testName: "正常なユーザー情報取得",
			userID:   "user_123",
			setupMocks: func(authUC *MockAuthUsecase, userRepo *MockUserRepository) {
				user := &model.User{
					ID:        "user_123",
					Email:     "test@example.com",
					Name:      "Test User",
					Picture:   "https://example.com/picture.jpg",
					CreatedAt: time.Now(),
					UpdatedAt: time.Now(),
				}
				userRepo.On("FindByID", mock.Anything, "user_123").Return(user, nil)
			},
			expectedStatus: http.StatusOK,
			expectError:    false,
		},
		{
			testName: "存在しないユーザー",
			userID:   "notfound",
			setupMocks: func(authUC *MockAuthUsecase, userRepo *MockUserRepository) {
				userRepo.On("FindByID", mock.Anything, "notfound").Return((*model.User)(nil), repository.ErrUserNotFound)
			},
			expectedStatus: http.StatusNotFound,
			expectError:    true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.testName, func(t *testing.T) {
			authUC := new(MockAuthUsecase)
			userRepo := new(MockUserRepository)
			tt.setupMocks(authUC, userRepo)

			handler := NewAuthHandler(authUC, userRepo)

			e := echo.New()
			req := httptest.NewRequest(http.MethodGet, "/auth/me", nil)
			rec := httptest.NewRecorder()
			c := e.NewContext(req, rec)
			c.Set("user_id", tt.userID)

			err := handler.GetMe(c)

			if tt.expectError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.expectedStatus, rec.Code)

				if rec.Code == http.StatusOK {
					var user model.User
					err := json.Unmarshal(rec.Body.Bytes(), &user)
					assert.NoError(t, err)
					assert.Equal(t, tt.userID, user.ID)
				}
			}

			authUC.AssertExpectations(t)
			userRepo.AssertExpectations(t)
		})
	}
}