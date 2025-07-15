package middleware

import (
	"errors"
	"net/http"
	"net/http/httptest"
	"stackies-backend/domain/service"
	"testing"

	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// MockJWTService はJWTServiceのモック
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

func TestAuthMiddleware_Authenticate(t *testing.T) {
	tests := []struct {
		testName       string
		authHeader     string
		setupMocks     func(*MockJWTService)
		expectedStatus int
		expectNext     bool
		expectedUserID string
	}{
		{
			testName:   "正常な認証",
			authHeader: "Bearer valid_token",
			setupMocks: func(jwtSvc *MockJWTService) {
				jwtSvc.On("ValidateToken", "valid_token").Return("user_123", nil)
			},
			expectedStatus: http.StatusOK,
			expectNext:     true,
			expectedUserID: "user_123",
		},
		{
			testName:   "Authorizationヘッダーなし",
			authHeader: "",
			setupMocks: func(jwtSvc *MockJWTService) {
				// モックの設定なし
			},
			expectedStatus: http.StatusUnauthorized,
			expectNext:     false,
			expectedUserID: "",
		},
		{
			testName:   "無効なAuthorizationヘッダー形式",
			authHeader: "invalid_header",
			setupMocks: func(jwtSvc *MockJWTService) {
				// モックの設定なし
			},
			expectedStatus: http.StatusUnauthorized,
			expectNext:     false,
			expectedUserID: "",
		},
		{
			testName:   "無効なトークン",
			authHeader: "Bearer invalid_token",
			setupMocks: func(jwtSvc *MockJWTService) {
				jwtSvc.On("ValidateToken", "invalid_token").Return("", errors.New("invalid token"))
			},
			expectedStatus: http.StatusUnauthorized,
			expectNext:     false,
			expectedUserID: "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.testName, func(t *testing.T) {
			jwtSvc := new(MockJWTService)
			tt.setupMocks(jwtSvc)

			middleware := NewAuthMiddleware(jwtSvc)

			// next関数が呼ばれたかチェックするフラグ
			nextCalled := false
			var capturedUserID interface{}

			next := func(c echo.Context) error {
				nextCalled = true
				capturedUserID = c.Get("user_id")
				return c.JSON(http.StatusOK, map[string]string{"message": "success"})
			}

			e := echo.New()
			req := httptest.NewRequest(http.MethodGet, "/protected", nil)
			if tt.authHeader != "" {
				req.Header.Set("Authorization", tt.authHeader)
			}
			rec := httptest.NewRecorder()
			c := e.NewContext(req, rec)

			handler := middleware.Authenticate(next)
			err := handler(c)

			if tt.expectNext {
				assert.NoError(t, err)
				assert.True(t, nextCalled)
				assert.Equal(t, tt.expectedUserID, capturedUserID)
				assert.Equal(t, tt.expectedStatus, rec.Code)
			} else {
				assert.Error(t, err)
				assert.False(t, nextCalled)
				
				// Echo HTTPErrorの場合はステータスコードをチェック
				if httpErr, ok := err.(*echo.HTTPError); ok {
					assert.Equal(t, tt.expectedStatus, httpErr.Code)
				}
			}

			jwtSvc.AssertExpectations(t)
		})
	}
}