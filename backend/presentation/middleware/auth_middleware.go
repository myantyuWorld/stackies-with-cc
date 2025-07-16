package middleware

import (
	"net/http"
	"stackies-backend/domain/service"
	"strings"

	"github.com/labstack/echo/v4"
)

// AuthMiddleware は認証ミドルウェアを表す
type AuthMiddleware struct {
	jwtSvc service.JWTService
}

// NewAuthMiddleware はAuthMiddlewareの新しいインスタンスを作成する
func NewAuthMiddleware(jwtSvc service.JWTService) *AuthMiddleware {
	return &AuthMiddleware{
		jwtSvc: jwtSvc,
	}
}

// Authenticate は認証ミドルウェアを表す
func (m *AuthMiddleware) Authenticate(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		authHeader := c.Request().Header.Get("Authorization")
		if authHeader == "" {
			return echo.NewHTTPError(http.StatusUnauthorized, "Authorization header required")
		}

		token := strings.TrimPrefix(authHeader, "Bearer ")
		if token == authHeader {
			return echo.NewHTTPError(http.StatusUnauthorized, "Invalid authorization header format")
		}

		userID, err := m.jwtSvc.ValidateToken(token)
		if err != nil {
			return echo.NewHTTPError(http.StatusUnauthorized, "Invalid token")
		}

		c.Set("user_id", userID)
		return next(c)
	}
}
