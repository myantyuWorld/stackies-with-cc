package service

// JWTService はJWT認証サービスを抽象化する
type JWTService interface {
	GenerateToken(userID string) (string, error)
	ValidateToken(token string) (string, error)
	GenerateRefreshToken(userID string) (string, error)
}
