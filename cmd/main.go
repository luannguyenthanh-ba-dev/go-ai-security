package main

//go:generate go run github.com/swaggo/swag/cmd/swag@latest init -g cmd/main.go -o ./docs

// @title           Go AI Security API
// @version         1.0
// @description     Backend API for AI Security project.
// @BasePath        /api/v1
// @schemes         http

import (
	"time"

	"github.com/gin-gonic/gin"
	config "github.com/luannguyenthanh-ba-dev/go-ai-security/config"

	// User
	userHttp "github.com/luannguyenthanh-ba-dev/go-ai-security/internal/users/delivery/http"
	userRepository "github.com/luannguyenthanh-ba-dev/go-ai-security/internal/users/repository"
	userUseCase "github.com/luannguyenthanh-ba-dev/go-ai-security/internal/users/usecase"

	// Auth
	authHttp "github.com/luannguyenthanh-ba-dev/go-ai-security/internal/auth/delivery/http"
	authUseCase "github.com/luannguyenthanh-ba-dev/go-ai-security/internal/auth/usecase"
	appLogger "github.com/luannguyenthanh-ba-dev/go-ai-security/pkg/logger"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
	"go.uber.org/zap"
)

func main() {
	// Initialize colorized console logger
	cfg, err := config.LoadConfig()

	logger, _ := appLogger.New(cfg.Env.AppEnv == "production")
	defer logger.Sync()
	zap.ReplaceGlobals(logger)

	if err != nil {
		zap.L().Fatal("failed to load config", zap.Error(err))
	}

	zap.L().Info("configuration loaded",
		zap.String("appName", cfg.Env.AppName),
		zap.String("env", cfg.Env.AppEnv),
		zap.String("port", cfg.Env.Port),
	)

	// MongoDB setup
	mongoDatabase, err := config.NewMongoDatabase(cfg.Env.MongoURI, cfg.Env.MongoDatabase)
	if err != nil {
		zap.L().Fatal("failed to create MongoDB database", zap.Error(err))
	}
	defer mongoDatabase.Close()

	userCollection := mongoDatabase.Database.Collection("users")

	// Gin setup
	if cfg.Env.AppEnv == "production" {
		// Set the Gin mode to release mode for production environment
		gin.SetMode(gin.ReleaseMode)
	} else {
		// Set the Gin mode to debug mode for development environment
		gin.SetMode(gin.DebugMode)
	}

	// Create a new Gin instance
	r := gin.New()
	// Use the Gin logger and recovery middleware
	r.Use(gin.Logger(), gin.Recovery())

	// Basic health endpoint
	r.GET("/health", healthHandler)

	// User routes
	api := r.Group("/api/v1")
	// User routes
	mongoUserRepository := userRepository.NewMongoUserRepository(userCollection)
	userService := userUseCase.NewUserService(mongoUserRepository, cfg.Env.PasswordHashSaltRounds)
	userHttp.RegisterUserRoutes(api, userService)

	// Auth routes
	jwtService := authUseCase.NewJWTService(cfg.Env.JWTSecret, time.Duration(cfg.Env.JWTExpiresIn) * time.Second)
	authService := authUseCase.NewAuthService(userService, jwtService)
	authHttp.RegisterAuthRoutes(api, authService)

	// Swagger UI Route (use local generated spec)
	r.Static("/docs", "./docs") // or: r.StaticFile("/docs/swagger.json", "./docs/swagger.json")
	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler, ginSwagger.URL("/docs/swagger.json")))

	// Set the address to the port specified in the environment variables
	addr := ":8080"
	if cfg.Env.Port != "" {
		addr = ":" + cfg.Env.Port
	}

	// Log the server starting
	zap.L().Info("starting HTTP server", zap.String("addr", addr))
	// Run the server
	if err := r.Run(addr); err != nil {
		zap.L().Fatal("server exited with error", zap.Error(err))
	}
}

// healthHandler godoc
// @Summary      Health check
// @Description  Returns service health status
// @Tags         health
// @Produce      json
// @Success      200  {object}  map[string]string
// @Router       /health [get]
func healthHandler(c *gin.Context) {
	c.JSON(200, gin.H{"status": "ok"})
}
