package routes

import (
	"backend-koda-shortlink/internal/handlers"

	"github.com/gin-gonic/gin"
)

func userRouter(r *gin.RouterGroup, userHandler *handlers.UserHandler) {
	r.GET("", userHandler.GetUserDetail)
}
