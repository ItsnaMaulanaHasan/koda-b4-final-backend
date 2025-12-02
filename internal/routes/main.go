package routes

import (
	"backend-koda-shortlink/internal/database"
	"backend-koda-shortlink/internal/handlers"
	"backend-koda-shortlink/internal/middlewares"
	"backend-koda-shortlink/internal/repository"
	"backend-koda-shortlink/internal/services"

	"github.com/gin-gonic/gin"
)

func SetUpRoutes(r *gin.Engine) {
	userRepo := repository.NewUserRepository(database.DB)
	sessionRepo := repository.NewSessionRepository(database.DB)
	shortLinkRepo := repository.NewShortLinkRepository(database.DB)
	clickRepo := repository.NewClickRepository(database.DB)
	dashboardRepo := repository.NewDashboardRepository(database.DB)

	authService := services.NewAuthService(userRepo, sessionRepo)
	shortLinkService := services.NewShortLinkService(shortLinkRepo, clickRepo)
	dashboardService := services.NewDashboardService(dashboardRepo)

	authHandler := handlers.NewAuthHandler(authService)
	shortLinkHandler := handlers.NewShortLinkHandler(shortLinkService)
	dashboardHandler := handlers.NewDashboardHandler(dashboardService)

	authMiddleware := middlewares.NewAuthMiddleware(*sessionRepo)
	optionalAuth := middlewares.NewOptionalAuthMiddleware(*sessionRepo)

	authRouter(r.Group("/api/v1/auth"), authHandler)
	shortLinkRoutes(r.Group("/api/v1/links", authMiddleware.Auth()), shortLinkHandler)

	r.POST("/api/v1/links", optionalAuth.OptionalAuth(), shortLinkHandler.CreateShortLink)

	r.GET("/:shortCode", shortLinkHandler.Redirect)

	r.GET("/api/v1/dashboard/stats", authMiddleware.Auth(), dashboardHandler.Stats)
}
