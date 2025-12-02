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

	authService := services.NewAuthService(userRepo, sessionRepo)
	shortLinkService := services.NewShortLinkService(shortLinkRepo)

	authHandler := handlers.NewAuthHandler(authService)
	shortLinkHandler := handlers.NewShortLinkHandler(shortLinkService)

	authMiddleware := middlewares.NewAuthMiddleware(*sessionRepo)

	authRouter(r.Group("/api/v1/auth"), authHandler)
	shortLinkRoutes(r.Group("/api/v1/links", authMiddleware.Auth()), shortLinkHandler)
	r.GET("/:shortCode", shortLinkHandler.Redirect)

}
