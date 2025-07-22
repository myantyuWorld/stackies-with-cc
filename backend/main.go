package main

import (
	"log"
	"net/http"
	"os"
	"stackies-backend/infra/external"
	"stackies-backend/infra/persistence"
	"stackies-backend/presentation/handler"
	authMiddleware "stackies-backend/presentation/middleware"
	"stackies-backend/registry"
	"stackies-backend/usecase"

	"github.com/joho/godotenv"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

func main() {
	// 環境変数を読み込み
	if err := godotenv.Load(); err != nil {
		log.Println("Warning: .env file not found, using system environment variables")
	}

	e := echo.New()

	e.Use(middleware.Logger())
	e.Use(middleware.Recover())

	// CORS設定を追加
	e.Use(middleware.CORSWithConfig(middleware.CORSConfig{
		AllowOrigins: []string{"http://localhost:3000", "http://localhost:5173"},
		AllowMethods: []string{http.MethodGet, http.MethodPost, http.MethodPut, http.MethodDelete, http.MethodOptions},
		AllowHeaders: []string{echo.HeaderOrigin, echo.HeaderContentType, echo.HeaderAccept, echo.HeaderAuthorization},
	}))

	container := registry.NewContainer()
	userRepo := persistence.NewUserRepository()
	authRepo := persistence.NewAuthRepository()
	authUsecase := usecase.NewAuthUsecase(userRepo, authRepo, container.GetGoogleService(), container.GetJWTService())
	googleSvc := external.NewGoogleService(os.Getenv("GOOGLE_CLIENT_ID"), os.Getenv("GOOGLE_CLIENT_SECRET"))
	jwtSvc := external.NewJWTService(os.Getenv("JWT_SECRET"))
	authMiddleware := authMiddleware.NewAuthMiddleware(jwtSvc)

	container.SetGoogleService(googleSvc)
	container.SetJWTService(jwtSvc)

	authHandler := handler.NewAuthHandler(authUsecase, userRepo, googleSvc)

	e.GET("/health", healthCheck)
	e.GET("/auth/google/url", authHandler.GoogleAuthURL)
	e.POST("/auth/google/login", authHandler.GoogleLogin)
	e.POST("/auth/refresh", authHandler.RefreshToken)
	e.POST("/auth/logout", authHandler.Logout)
	e.GET("/auth/me", authHandler.GetMe, authMiddleware.Authenticate)

	// ポート設定（環境変数から取得、デフォルトは8080）
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	e.Logger.Fatal(e.Start(":" + port))
}

func healthCheck(c echo.Context) error {
	return c.JSON(http.StatusOK, map[string]string{
		"status":  "OK",
		"message": "Server is running",
	})
}
