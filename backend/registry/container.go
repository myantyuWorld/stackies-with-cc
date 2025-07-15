package registry

import (
	"os"
	"stackies-backend/domain/service"
	"stackies-backend/infra/external"
)

// Container はDependency Injection用のコンテナ
type Container struct {
	googleService service.GoogleService
	jwtService    service.JWTService
}

// NewContainer は新しいコンテナを作成する
func NewContainer() *Container {
	return &Container{}
}

// GetGoogleService はGoogleServiceの実装を返す
func (c *Container) GetGoogleService() service.GoogleService {
	if c.googleService == nil {
		clientID := os.Getenv("GOOGLE_CLIENT_ID")
		clientSecret := os.Getenv("GOOGLE_CLIENT_SECRET")
		c.googleService = external.NewGoogleService(clientID, clientSecret)
	}
	return c.googleService
}

// GetJWTService はJWTServiceの実装を返す
func (c *Container) GetJWTService() service.JWTService {
	if c.jwtService == nil {
		secretKey := os.Getenv("JWT_SECRET")
		c.jwtService = external.NewJWTService(secretKey)
	}
	return c.jwtService
}

// SetGoogleService はテスト用にGoogleServiceをセットする
func (c *Container) SetGoogleService(gs service.GoogleService) {
	c.googleService = gs
}

// SetJWTService はテスト用にJWTServiceをセットする
func (c *Container) SetJWTService(js service.JWTService) {
	c.jwtService = js
}