package external

// JWTService はJWT関連の処理を抽象化する
type JWTService interface {
	GenerateToken(userID string) (string, error)
	ValidateToken(tokenString string) (*JWTClaims, error)
	RefreshToken(refreshToken string) (string, error)
}

// JWTClaims はJWTのクレーム情報を表す
type JWTClaims struct {
	UserID string `json:"user_id"`
}
