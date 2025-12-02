package routes

import (
	"backend-koda-shortlink/internal/controllers"

	"github.com/gin-gonic/gin"
)

func authRouter(r *gin.RouterGroup) {
	r.POST("/register", controllers.Register)
	r.POST("/login", controllers.Login)
	r.POST("/refresh", controllers.RefreshToken)
	r.POST("/logout", controllers.Logout)
}
