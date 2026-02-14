package main

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	_ "github.com/TubagusAldiMY/go-template/docs" // Import swagger docs
	"github.com/TubagusAldiMY/go-template/internal/delivery/http/router"
	userHttp "github.com/TubagusAldiMY/go-template/internal/domain/user/delivery/http"
	userRepo "github.com/TubagusAldiMY/go-template/internal/domain/user/repository"
	userUsecase "github.com/TubagusAldiMY/go-template/internal/domain/user/usecase"
	"github.com/TubagusAldiMY/go-template/internal/infrastructure/cache"
	"github.com/TubagusAldiMY/go-template/internal/infrastructure/config"
	"github.com/TubagusAldiMY/go-template/internal/infrastructure/database"
	"github.com/TubagusAldiMY/go-template/internal/infrastructure/messaging"
	"github.com/TubagusAldiMY/go-template/pkg/crypto"
	"github.com/TubagusAldiMY/go-template/pkg/jwt"
	"github.com/TubagusAldiMY/go-template/pkg/logger"
	"github.com/TubagusAldiMY/go-template/pkg/validator"
	"go.uber.org/zap"
)

// @title Golang DDD Template API
// @version 1.0
// @description Production-ready REST API built with Clean Architecture and DDD
// @termsOfService http://swagger.io/terms/

// @contact.name API Support
// @contact.email support@example.com

// @license.name MIT
// @license.url https://opensource.org/licenses/MIT

// @host localhost:8080
// @BasePath /api/v1

// @securityDefinitions.apikey Bearer
// @in header
// @name Authorization
// @description Type "Bearer" followed by a space and JWT token.

func main() {
	// Load configuration
	cfg, err := config.Load()
	if err != nil {
		fmt.Printf("Failed to load config: %v\n", err)
		os.Exit(1)
	}

	// Initialize logger
	if err := logger.Init(logger.Config{
		Level:  cfg.Log.Level,
		Format: cfg.Log.Format,
		Output: cfg.Log.Output,
	}); err != nil {
		fmt.Printf("Failed to initialize logger: %v\n", err)
		os.Exit(1)
	}
	defer logger.Sync()

	logger.Info("starting application",
		zap.String("app", cfg.App.Name),
		zap.String("env", cfg.App.Env),
		zap.Int("port", cfg.App.Port),
	)

	// Initialize validator
	if err := validator.Init(); err != nil {
		logger.Fatal("failed to initialize validator", zap.Error(err))
	}

	// Initialize database
	db, err := database.NewPostgreSQL(cfg.Database)
	if err != nil {
		logger.Fatal("failed to connect to database", zap.Error(err))
	}
	defer db.Close()

	// Initialize Redis
	redisClient, err := cache.NewRedis(cfg.Redis)
	if err != nil {
		logger.Fatal("failed to connect to redis", zap.Error(err))
	}
	defer redisClient.Close()

	// Initialize RabbitMQ
	rabbitmq, err := messaging.NewRabbitMQ(cfg.RabbitMQ)
	if err != nil {
		logger.Warn("failed to connect to rabbitmq", zap.Error(err))
		// RabbitMQ is optional, continue without it
	} else {
		defer rabbitmq.Close()
	}

	// Initialize utilities
	passwordHasher := crypto.NewPasswordHasher(cfg.Security.BcryptCost)
	jwtManager := jwt.NewManager(
		cfg.JWT.Secret,
		cfg.JWT.AccessTokenExpiry,
		cfg.JWT.RefreshTokenExpiry,
	)

	// Initialize repositories
	userRepository := userRepo.NewPostgresUserRepository(db.GetPool())

	// Initialize use cases
	userUsecaseImpl := userUsecase.NewUserUsecase(
		userRepository,
		passwordHasher,
		jwtManager,
		redisClient,
	)

	// Initialize handlers
	userHandler := userHttp.NewUserHandler(userUsecaseImpl)

	// Setup router
	routerCfg := &router.RouterConfig{
		Config:      cfg,
		JWTManager:  jwtManager,
		UserHandler: userHandler,
	}
	r := router.SetupRouter(routerCfg)

	// Create HTTP server
	srv := &http.Server{
		Addr:         fmt.Sprintf(":%d", cfg.App.Port),
		Handler:      r,
		ReadTimeout:  cfg.Server.ReadTimeout,
		WriteTimeout: cfg.Server.WriteTimeout,
		IdleTimeout:  cfg.Server.IdleTimeout,
	}

	// Start server in goroutine
	go func() {
		logger.Info("server started",
			zap.String("address", srv.Addr),
		)
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			logger.Fatal("failed to start server", zap.Error(err))
		}
	}()

	// Wait for interrupt signal
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	logger.Info("shutting down server...")

	// Graceful shutdown
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		logger.Error("server forced to shutdown", zap.Error(err))
	}

	logger.Info("server exited")
}
