package external

import (
	"context"
	"stackies-backend/domain/model"
)

// GoogleService はGoogle OAuth2.0関連の外部サービスを抽象化する
type GoogleService interface {
	ExchangeCodeForToken(ctx context.Context, code, redirectURI string) (*model.AuthToken, error)
	GetUserInfo(ctx context.Context, accessToken string) (*model.GoogleUserInfo, error)
}
