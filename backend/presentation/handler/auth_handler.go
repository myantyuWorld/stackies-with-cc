package handler

import (
	"net/http"
	"stackies-backend/domain/model"
	"stackies-backend/domain/repository"
	"stackies-backend/usecase"

	"github.com/labstack/echo/v4"
)

// AuthHandler は認証関連のHTTPハンドラーを表す
type AuthHandler struct {
	authUsecase usecase.AuthUsecase
	userRepo    repository.UserRepository
}

// NewAuthHandler はAuthHandlerの新しいインスタンスを作成する
func NewAuthHandler(authUsecase usecase.AuthUsecase, userRepo repository.UserRepository) *AuthHandler {
	return &AuthHandler{
		authUsecase: authUsecase,
		userRepo:    userRepo,
	}
}

type (
	// GoogleLoginRequest はGoogleログインのリクエスト構造体を表す
	GoogleLoginRequest struct {
		AuthorizationCode string `json:"authorization_code" validate:"required"`
		RedirectURI       string `json:"redirect_uri" validate:"required"`
	}

	// GoogleLoginResponse はGoogleログインのレスポンス構造体を表す
	GoogleLoginResponse struct {
		User         *model.User `json:"user"`
		AccessToken  string      `json:"access_token"`
		RefreshToken string      `json:"refresh_token"`
		ExpiresIn    int64       `json:"expires_in"`
	}

	// RefreshTokenRequest はトークンリフレッシュのリクエスト構造体を表す
	RefreshTokenRequest struct {
		RefreshToken string `json:"refresh_token" validate:"required"`
	}

	// RefreshTokenResponse はトークンリフレッシュのレスポンス構造体を表す
	RefreshTokenResponse struct {
		AccessToken  string `json:"access_token"`
		RefreshToken string `json:"refresh_token"`
		ExpiresIn    int64  `json:"expires_in"`
	}
)

// GoogleLogin はGoogleログインのハンドラーメソッドを表す
func (h *AuthHandler) GoogleLogin(c echo.Context) error {
	var req GoogleLoginRequest
	if err := c.Bind(&req); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "Invalid request")
	}

	input := &usecase.GoogleLoginInput{
		AuthorizationCode: req.AuthorizationCode,
		RedirectURI:       req.RedirectURI,
	}

	output, err := h.authUsecase.GoogleLogin(c.Request().Context(), input)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	response := &GoogleLoginResponse{
		User:         output.User,
		AccessToken:  output.AccessToken,
		RefreshToken: output.RefreshToken,
		ExpiresIn:    output.ExpiresIn,
	}

	return c.JSON(http.StatusOK, response)
}

// RefreshToken はトークンリフレッシュのハンドラーメソッドを表す
func (h *AuthHandler) RefreshToken(c echo.Context) error {
	var req RefreshTokenRequest
	if err := c.Bind(&req); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "Invalid request")
	}

	input := &usecase.RefreshTokenInput{
		RefreshToken: req.RefreshToken,
	}

	output, err := h.authUsecase.RefreshToken(c.Request().Context(), input)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	response := &RefreshTokenResponse{
		AccessToken:  output.AccessToken,
		RefreshToken: output.RefreshToken,
		ExpiresIn:    output.ExpiresIn,
	}

	return c.JSON(http.StatusOK, response)
}

// Logout はログアウトのハンドラーメソッドを表す
func (h *AuthHandler) Logout(c echo.Context) error {
	userID := c.Get("user_id").(string)

	input := &usecase.LogoutInput{
		UserID: userID,
	}

	if err := h.authUsecase.Logout(c.Request().Context(), input); err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	return c.JSON(http.StatusOK, map[string]string{"message": "Logged out successfully"})
}

// GetMe は現在認証されているユーザーの情報を取得するハンドラーメソッドを表す
func (h *AuthHandler) GetMe(c echo.Context) error {
	userID := c.Get("user_id").(string)

	user, err := h.userRepo.FindByID(c.Request().Context(), userID)
	if err != nil {
		if err == repository.ErrUserNotFound {
			return echo.NewHTTPError(http.StatusNotFound, "User not found")
		}
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	return c.JSON(http.StatusOK, user)
}
