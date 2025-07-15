package usecase

import (
	"context"
	"stackies-backend/domain/model"
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
)
