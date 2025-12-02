package routes

import (
	"github.com/gin-gonic/gin"
)

func SetUpRoutes(r *gin.Engine) {
	authRouter(r.Group("/api/v1/auth"))
}
