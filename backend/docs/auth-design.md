# GoogleログインAPI設計書

## 概要
Google OAuth2.0を使用したログイン機能を実装する。

## アーキテクチャ

### 1. Domain Layer

#### 1.1 エンティティ・値オブジェクト
```go
// domain/model/user.go
type User struct {
    ID        string
    Email     string
    Name      string
    Picture   string
    CreatedAt time.Time
    UpdatedAt time.Time
}

// domain/model/auth.go
type (
    AuthToken struct {
        AccessToken  string
        RefreshToken string
        ExpiresIn    int64
        TokenType    string
    }

    GoogleUserInfo struct {
        ID            string
        Email         string
        VerifiedEmail bool
        Name          string
        GivenName     string
        FamilyName    string
        Picture       string
        Locale        string
    }
)
```

#### 1.2 リポジトリインターフェース
```go
// domain/repository/user_repository.go
type UserRepository interface {
    Create(ctx context.Context, user *User) error
    FindByEmail(ctx context.Context, email string) (*User, error)
    FindByID(ctx context.Context, id string) (*User, error)
    Update(ctx context.Context, user *User) error
}

// domain/repository/auth_repository.go
type AuthRepository interface {
    SaveToken(ctx context.Context, userID string, token *AuthToken) error
    GetToken(ctx context.Context, userID string) (*AuthToken, error)
    DeleteToken(ctx context.Context, userID string) error
}
```

### 2. Usecase Layer

#### 2.1 インターフェース
```go
// usecase/auth_usecase.go
type AuthUsecase interface {
    GoogleLogin(ctx context.Context, input *GoogleLoginInput) (*GoogleLoginOutput, error)
    RefreshToken(ctx context.Context, input *RefreshTokenInput) (*RefreshTokenOutput, error)
    Logout(ctx context.Context, input *LogoutInput) error
}

type (
    GoogleLoginInput struct {
        AuthorizationCode string
        RedirectURI      string
    }

    GoogleLoginOutput struct {
        User         *User
        AccessToken  string
        RefreshToken string
        ExpiresIn    int64
    }

    RefreshTokenInput struct {
        RefreshToken string
    }

    RefreshTokenOutput struct {
        AccessToken  string
        RefreshToken string
        ExpiresIn    int64
    }

    LogoutInput struct {
        UserID string
    }
)
```

#### 2.2 実装
```go
// usecase/auth_usecase_impl.go
type authUsecase struct {
    userRepo    domain.UserRepository
    authRepo    domain.AuthRepository
    googleSvc   infra.GoogleService
    jwtSvc      infra.JWTService
}

func NewAuthUsecase(
    userRepo domain.UserRepository,
    authRepo domain.AuthRepository,
    googleSvc infra.GoogleService,
    jwtSvc infra.JWTService,
) AuthUsecase {
    return &authUsecase{
        userRepo:  userRepo,
        authRepo:  authRepo,
        googleSvc: googleSvc,
        jwtSvc:    jwtSvc,
    }
}
```

### 3. Presentation Layer

#### 3.1 ハンドラー
```go
// presentation/handler/auth_handler.go
type AuthHandler struct {
    authUsecase usecase.AuthUsecase
}

func NewAuthHandler(authUsecase usecase.AuthUsecase) *AuthHandler {
    return &AuthHandler{
        authUsecase: authUsecase,
    }
}

// リクエスト・レスポンス構造体
type (
    GoogleLoginRequest struct {
        AuthorizationCode string `json:"authorization_code" validate:"required"`
        RedirectURI      string `json:"redirect_uri" validate:"required"`
    }

    GoogleLoginResponse struct {
        User         *User  `json:"user"`
        AccessToken  string `json:"access_token"`
        RefreshToken string `json:"refresh_token"`
        ExpiresIn    int64  `json:"expires_in"`
    }

    RefreshTokenRequest struct {
        RefreshToken string `json:"refresh_token" validate:"required"`
    }

    RefreshTokenResponse struct {
        AccessToken  string `json:"access_token"`
        RefreshToken string `json:"refresh_token"`
        ExpiresIn    int64  `json:"expires_in"`
    }
)

// ハンドラーメソッド
func (h *AuthHandler) GoogleLogin(c echo.Context) error {
    var req GoogleLoginRequest
    if err := c.Bind(&req); err != nil {
        return echo.NewHTTPError(http.StatusBadRequest, "Invalid request")
    }
    
    input := &usecase.GoogleLoginInput{
        AuthorizationCode: req.AuthorizationCode,
        RedirectURI:      req.RedirectURI,
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
```

#### 3.2 ミドルウェア
```go
// presentation/middleware/auth_middleware.go
type AuthMiddleware struct {
    jwtSvc infra.JWTService
}

func NewAuthMiddleware(jwtSvc infra.JWTService) *AuthMiddleware {
    return &AuthMiddleware{
        jwtSvc: jwtSvc,
    }
}

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
        
        claims, err := m.jwtSvc.ValidateToken(token)
        if err != nil {
            return echo.NewHTTPError(http.StatusUnauthorized, "Invalid token")
        }
        
        c.Set("user_id", claims.UserID)
        return next(c)
    }
}
```

### 4. Infrastructure Layer

#### 4.1 外部サービス抽象化
```go
// infra/external/google_service.go
type GoogleService interface {
    ExchangeCodeForToken(ctx context.Context, code, redirectURI string) (*AuthToken, error)
    GetUserInfo(ctx context.Context, accessToken string) (*GoogleUserInfo, error)
}

type googleService struct {
    clientID     string
    clientSecret string
    httpClient   *http.Client
}

func NewGoogleService(clientID, clientSecret string) GoogleService {
    return &googleService{
        clientID:     clientID,
        clientSecret: clientSecret,
        httpClient:   &http.Client{Timeout: 10 * time.Second},
    }
}
```

#### 4.2 JWTサービス
```go
// infra/external/jwt_service.go
type JWTService interface {
    GenerateToken(userID string) (string, error)
    ValidateToken(tokenString string) (*JWTClaims, error)
    RefreshToken(refreshToken string) (string, error)
}

type JWTClaims struct {
    UserID string `json:"user_id"`
    jwt.RegisteredClaims
}

type jwtService struct {
    secretKey []byte
}

func NewJWTService(secretKey string) JWTService {
    return &jwtService{
        secretKey: []byte(secretKey),
    }
}
```

#### 4.3 DTO
```go
// infra/dto/user_dto.go
type UserDTO struct {
    ID        string    `db:"id"`
    Email     string    `db:"email"`
    Name      string    `db:"name"`
    Picture   string    `db:"picture"`
    CreatedAt time.Time `db:"created_at"`
    UpdatedAt time.Time `db:"updated_at"`
}

func (dto *UserDTO) ToDomain() *domain.User {
    return &domain.User{
        ID:        dto.ID,
        Email:     dto.Email,
        Name:      dto.Name,
        Picture:   dto.Picture,
        CreatedAt: dto.CreatedAt,
        UpdatedAt: dto.UpdatedAt,
    }
}

func (dto *UserDTO) FromDomain(user *domain.User) {
    dto.ID = user.ID
    dto.Email = user.Email
    dto.Name = user.Name
    dto.Picture = user.Picture
    dto.CreatedAt = user.CreatedAt
    dto.UpdatedAt = user.UpdatedAt
}
```

### 5. Registry Layer (DI)
```go
// registry/registry.go
type Registry struct {
    userRepo    domain.UserRepository
    authRepo    domain.AuthRepository
    googleSvc   infra.GoogleService
    jwtSvc      infra.JWTService
    authUsecase usecase.AuthUsecase
    authHandler *presentation.AuthHandler
    authMiddleware *presentation.AuthMiddleware
}

func NewRegistry() *Registry {
    // 依存関係の初期化
    googleSvc := infra.NewGoogleService(os.Getenv("GOOGLE_CLIENT_ID"), os.Getenv("GOOGLE_CLIENT_SECRET"))
    jwtSvc := infra.NewJWTService(os.Getenv("JWT_SECRET_KEY"))
    
    userRepo := infra.NewUserRepository(db)
    authRepo := infra.NewAuthRepository(redis)
    
    authUsecase := usecase.NewAuthUsecase(userRepo, authRepo, googleSvc, jwtSvc)
    authHandler := presentation.NewAuthHandler(authUsecase)
    authMiddleware := presentation.NewAuthMiddleware(jwtSvc)
    
    return &Registry{
        userRepo:        userRepo,
        authRepo:        authRepo,
        googleSvc:       googleSvc,
        jwtSvc:          jwtSvc,
        authUsecase:     authUsecase,
        authHandler:     authHandler,
        authMiddleware:  authMiddleware,
    }
}
```

## API エンドポイント

### 1. Googleログイン
```
POST /api/auth/google/login
Content-Type: application/json

{
  "authorization_code": "4/0AfJohXn...",
  "redirect_uri": "http://localhost:3000/callback"
}

Response:
{
  "user": {
    "id": "user_123",
    "email": "user@example.com",
    "name": "John Doe",
    "picture": "https://...",
    "created_at": "2024-01-01T00:00:00Z",
    "updated_at": "2024-01-01T00:00:00Z"
  },
  "access_token": "eyJhbGciOiJIUzI1NiIs...",
  "refresh_token": "refresh_token_123",
  "expires_in": 3600
}
```

### 2. トークンリフレッシュ
```
POST /api/auth/refresh
Content-Type: application/json

{
  "refresh_token": "refresh_token_123"
}

Response:
{
  "access_token": "eyJhbGciOiJIUzI1NiIs...",
  "refresh_token": "new_refresh_token_456",
  "expires_in": 3600
}
```

### 3. ログアウト
```
POST /api/auth/logout
Authorization: Bearer eyJhbGciOiJIUzI1NiIs...

Response:
{
  "message": "Logged out successfully"
}
```

## 必要な環境変数
- `GOOGLE_CLIENT_ID`: Google OAuth2.0 クライアントID
- `GOOGLE_CLIENT_SECRET`: Google OAuth2.0 クライアントシークレット
- `JWT_SECRET_KEY`: JWT署名用のシークレットキー
- `DATABASE_URL`: データベース接続URL
- `REDIS_URL`: Redis接続URL

## 実装順序
1. Domain Layer (エンティティ、リポジトリインターフェース)
2. Infrastructure Layer (外部サービス、DTO)
3. Usecase Layer (ビジネスロジック)
4. Presentation Layer (ハンドラー、ミドルウェア)
5. Registry Layer (DI)
6. メインアプリケーションでの統合 