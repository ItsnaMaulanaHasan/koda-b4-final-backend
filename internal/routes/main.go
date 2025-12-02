package routes

import (
	"backend-koda-shortlink/internal/handlers"
	"backend-koda-shortlink/internal/repository"
	"backend-koda-shortlink/internal/services"

	"github.com/gin-gonic/gin"
)

func SetUpRoutes(r *gin.Engine) {
	userRepo := repository.NewUserRepository()
	sessionRepo := repository.NewSessionRepository()

	authService := services.NewAuthService(userRepo, sessionRepo)

	authHandler := handlers.NewAuthHandler(authService)

	// authMiddleware := middlewares.NewAuthMiddleware(sessionRepo)

	authRouter(r.Group("/api/v1/auth"), authHandler)
}
