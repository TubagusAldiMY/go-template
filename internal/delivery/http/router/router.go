package router

import (
	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"

	"github.com/TubagusAldiMY/go-template/internal/delivery/http/middleware"
	userHttp "github.com/TubagusAldiMY/go-template/internal/domain/user/delivery/http"
	"github.com/TubagusAldiMY/go-template/internal/infrastructure/config"
	"github.com/TubagusAldiMY/go-template/internal/shared/constants"
	"github.com/TubagusAldiMY/go-template/pkg/jwt"
	"github.com/TubagusAldiMY/go-template/pkg/response"
)

type RouterConfig struct {
	Config      *config.Config
	JWTManager  *jwt.Manager
	UserHandler *userHttp.UserHandler
}

func SetupRouter(cfg *RouterConfig) *gin.Engine {
	// Set gin mode
	if !cfg.Config.App.Debug {
		gin.SetMode(gin.ReleaseMode)
	}

	router := gin.New()

	// Global middleware
	router.Use(middleware.Recovery())
	router.Use(middleware.RequestLogger())
	router.Use(middleware.CORS(cfg.Config.CORS))
	router.Use(middleware.RateLimit(cfg.Config.RateLimit))

	// Health check
	router.GET("/health", func(c *gin.Context) {
		response.OK(c, "Service is healthy", gin.H{
			"service": cfg.Config.App.Name,
			"version": "1.0.0",
		})
	})

	// Swagger documentation
	if cfg.Config.App.Debug {
		router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
	}

	// API v1 routes
	v1 := router.Group("/api/v1")
	{
		// Auth routes (public)
		auth := v1.Group("/auth")
		{
			auth.POST("/register", cfg.UserHandler.Register)
			auth.POST("/login", cfg.UserHandler.Login)
			auth.POST("/refresh", cfg.UserHandler.RefreshToken)
		}

		// User routes (protected)
		users := v1.Group("/users")
		users.Use(middleware.AuthMiddleware(cfg.JWTManager))
		{
			users.GET("/profile", cfg.UserHandler.GetProfile)
			users.PUT("/profile", cfg.UserHandler.UpdateProfile)
			users.POST("/change-password", cfg.UserHandler.ChangePassword)

			// Admin only routes
			users.GET("", middleware.RequireRole(constants.RoleAdmin), cfg.UserHandler.ListUsers)
			users.DELETE("/:id", middleware.RequireRole(constants.RoleAdmin), cfg.UserHandler.DeleteUser)
		}
	}

	return router
}
