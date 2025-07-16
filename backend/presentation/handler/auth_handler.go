package handler

import (
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"net/http"
	"stackies-backend/domain/model"
	"stackies-backend/domain/repository"
	"stackies-backend/domain/service"
	"stackies-backend/usecase"

	"github.com/davecgh/go-spew/spew"
	"github.com/labstack/echo/v4"
)

// AuthHandler は認証関連のHTTPハンドラーを表す
type AuthHandler struct {
	authUsecase usecase.AuthUsecase
	userRepo    repository.UserRepository
	googleSvc   service.GoogleService
	stateStore  map[string]bool // 本来はRedisなどを使用
}

// NewAuthHandler はAuthHandlerの新しいインスタンスを作成する
func NewAuthHandler(authUsecase usecase.AuthUsecase, userRepo repository.UserRepository, googleSvc service.GoogleService) *AuthHandler {
	return &AuthHandler{
		authUsecase: authUsecase,
		userRepo:    userRepo,
		googleSvc:   googleSvc,
		stateStore:  make(map[string]bool),
	}
}

type (
	// GoogleAuthURLResponse はGoogle認証URL生成のレスポンス構造体を表す
	GoogleAuthURLResponse struct {
		AuthURL string `json:"auth_url"`
		State   string `json:"state"`
	}

	// GoogleLoginRequest はGoogleログインのリクエスト構造体を表す
	GoogleLoginRequest struct {
		State string `json:"state" query:"state" validate:"required"`
		Code  string `json:"code" query:"code" validate:"required"`
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

// generateState はCSRF対策用のランダムなstateを生成する
func (h *AuthHandler) generateState() (string, error) {
	bytes := make([]byte, 16)
	if _, err := rand.Read(bytes); err != nil {
		return "", err
	}
	state := hex.EncodeToString(bytes)
	h.stateStore[state] = true
	return state, nil
}

// validateState はstateの有効性を検証する
func (h *AuthHandler) validateState(state string) bool {
	if _, exists := h.stateStore[state]; exists {
		delete(h.stateStore, state) // 一度使用したら削除
		return true
	}
	return false
}

// GoogleAuthURL はGoogle認証URLを生成するハンドラーメソッドを表す
func (h *AuthHandler) GoogleAuthURL(c echo.Context) error {
	state, err := h.generateState()
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "Failed to generate state")
	}

	authURL := h.googleSvc.GenerateAuthURL(state)

	response := &GoogleAuthURLResponse{
		AuthURL: authURL,
		State:   state,
	}

	return c.Redirect(http.StatusTemporaryRedirect, response.AuthURL)
	// return c.JSON(http.StatusOK, response)
}

// GoogleLogin はGoogleログインのハンドラーメソッドを表す
func (h *AuthHandler) GoogleLogin(c echo.Context) error {
	var req GoogleLoginRequest
	if err := c.Bind(&req); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "Invalid request")
	}

	spew.Dump(req)

	// CSRF対策：stateを検証
	// if !h.validateState(req.State) {
	// 	return echo.NewHTTPError(http.StatusBadRequest, "Invalid state parameter")
	// }

	input := &usecase.GoogleLoginInput{
		AuthorizationCode: req.Code,
	}

	output, err := h.authUsecase.GoogleLogin(c.Request().Context(), input)
	if err != nil {
		spew.Dump(err)
		fmt.Println(err)
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
