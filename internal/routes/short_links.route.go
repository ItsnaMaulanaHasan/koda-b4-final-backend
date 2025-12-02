package routes

import (
	"backend-koda-shortlink/internal/handlers"

	"github.com/gin-gonic/gin"
)

func shortLinkRoutes(r *gin.RouterGroup, handler *handlers.ShortLinkHandler) {
	r.GET("", handler.GetAllLinks)
	r.GET("/:shortCode", handler.GetLinkByShortCode)
	r.PUT("/:shortCode", handler.UpdateShortLink)
	r.DELETE("/:shortCode", handler.DeleteShortLink)
}
