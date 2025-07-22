package usecase

import (
	"context"
	"errors"
	"stackies-backend/domain/model"
	"stackies-backend/domain/repository"
	"stackies-backend/domain/service"
	"time"
)

// AuthUsecase は認証関連のビジネスロジックを抽象化する
type AuthUsecase interface {
	GoogleLogin(ctx context.Context, input *GoogleLoginInput) (*GoogleLoginOutput, error)
	RefreshToken(ctx context.Context, input *RefreshTokenInput) (*RefreshTokenOutput, error)
	Logout(ctx context.Context, input *LogoutInput) error
}

type (
	// GoogleLoginInput はGoogleログインの入力パラメータを表す
	GoogleLoginInput struct {
		AuthorizationCode string
		RedirectURI       string
	}

	// GoogleLoginOutput はGoogleログインの出力パラメータを表す
	GoogleLoginOutput struct {
		User         *model.User
		AccessToken  string
		RefreshToken string
		ExpiresIn    int64
	}

	// RefreshTokenInput はトークンリフレッシュの入力パラメータを表す
	RefreshTokenInput struct {
		RefreshToken string
	}

	// RefreshTokenOutput はトークンリフレッシュの出力パラメータを表す
	RefreshTokenOutput struct {
		AccessToken  string
		RefreshToken string
		ExpiresIn    int64
	}

	// LogoutInput はログアウトの入力パラメータを表す
	LogoutInput struct {
		UserID string
	}

	// AuthUsecaseImpl はAuthUsecaseの実装
	AuthUsecaseImpl struct {
		userRepo  repository.UserRepository
		authRepo  repository.AuthRepository
		googleSvc service.GoogleService
		jwtSvc    service.JWTService
	}
)

// NewAuthUsecase は新しいAuthUsecaseを作成する
func NewAuthUsecase(
	userRepo repository.UserRepository,
	authRepo repository.AuthRepository,
	googleSvc service.GoogleService,
	jwtSvc service.JWTService,
) AuthUsecase {
	return &AuthUsecaseImpl{
		userRepo:  userRepo,
		authRepo:  authRepo,
		googleSvc: googleSvc,
		jwtSvc:    jwtSvc,
	}
}

// GoogleLogin はGoogle OAuth2.0を使用したログインを処理する
func (a *AuthUsecaseImpl) GoogleLogin(ctx context.Context, input *GoogleLoginInput) (*GoogleLoginOutput, error) {
	// 1. 認証コードをアクセストークンに交換
	googleToken, err := a.googleSvc.ExchangeCode(ctx, input.AuthorizationCode, input.RedirectURI)
	if err != nil {
		return nil, err
	}

	// 2. Googleからユーザー情報を取得
	googleUser, err := a.googleSvc.GetUserInfo(ctx, googleToken.AccessToken)
	if err != nil {
		return nil, err
	}

	// 3. メールアドレスが認証済みかチェック
	if !googleUser.IsVerified() {
		return nil, errors.New("email not verified")
	}

	// 4. 既存ユーザーかどうかチェック
	var user *model.User
	existingUser, err := a.userRepo.FindByEmail(ctx, googleUser.Email)
	if err != nil && err != repository.ErrUserNotFound {
		return nil, err
	}

	if existingUser != nil {
		// 既存ユーザーの場合、プロフィールを更新
		err = existingUser.UpdateProfile(googleUser.Name, googleUser.Picture)
		if err != nil {
			return nil, err
		}
		err = a.userRepo.Update(ctx, existingUser)
		if err != nil {
			return nil, err
		}
		user = existingUser
	} else {
		// 新規ユーザーの場合、作成
		user = googleUser.ToUser()
		err = a.userRepo.Save(ctx, user)
		if err != nil {
			return nil, err
		}
	}

	// 5. JWTトークンを生成
	accessToken, err := a.jwtSvc.GenerateToken(user.ID)
	if err != nil {
		return nil, err
	}

	refreshToken, err := a.jwtSvc.GenerateRefreshToken(user.ID)
	if err != nil {
		return nil, err
	}

	// 6. トークンを保存
	expiresIn := time.Now().Add(time.Hour).Unix()
	authToken, err := model.NewAuthToken(accessToken, refreshToken, expiresIn, "Bearer")
	if err != nil {
		return nil, err
	}

	err = a.authRepo.SaveToken(ctx, user.ID, authToken)
	if err != nil {
		return nil, err
	}

	return &GoogleLoginOutput{
		User:         user,
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
		ExpiresIn:    expiresIn,
	}, nil
}

// RefreshToken はリフレッシュトークンを使用してアクセストークンを更新する
func (a *AuthUsecaseImpl) RefreshToken(ctx context.Context, input *RefreshTokenInput) (*RefreshTokenOutput, error) {
	// 1. リフレッシュトークンを検証
	userID, err := a.jwtSvc.ValidateToken(input.RefreshToken)
	if err != nil {
		return nil, errors.New("invalid refresh token")
	}

	// 2. 新しいトークンを生成
	accessToken, err := a.jwtSvc.GenerateToken(userID)
	if err != nil {
		return nil, err
	}

	refreshToken, err := a.jwtSvc.GenerateRefreshToken(userID)
	if err != nil {
		return nil, err
	}

	// 3. トークンを保存
	expiresIn := time.Now().Add(time.Hour).Unix()
	authToken, err := model.NewAuthToken(accessToken, refreshToken, expiresIn, "Bearer")
	if err != nil {
		return nil, err
	}

	err = a.authRepo.SaveToken(ctx, userID, authToken)
	if err != nil {
		return nil, err
	}

	return &RefreshTokenOutput{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
		ExpiresIn:    expiresIn,
	}, nil
}

// Logout はユーザーのログアウトを処理する
func (a *AuthUsecaseImpl) Logout(ctx context.Context, input *LogoutInput) error {
	// トークンを削除
	return a.authRepo.DeleteToken(ctx, input.UserID)
}
