package service

import (
	"context"
	"stackies-backend/domain/model"
)

// GoogleService はGoogle OAuth2.0サービスを抽象化する
type GoogleService interface {
	ExchangeCode(ctx context.Context, code, redirectURI string) (*model.AuthToken, error)
	GetUserInfo(ctx context.Context, accessToken string) (*model.GoogleUserInfo, error)
}