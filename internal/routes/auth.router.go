package routes

import (
	"backend-koda-shortlink/internal/handlers"

	"github.com/gin-gonic/gin"
)

func authRouter(r *gin.RouterGroup, authHandler *handlers.AuthHandler) {
	r.POST("/register", authHandler.Register)
	r.POST("/login", authHandler.Login)
	r.POST("/refresh", authHandler.RefreshToken)
	r.POST("/logout", authHandler.Logout)
}
