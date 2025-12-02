package middlewares

import (
	"backend-koda-shortlink/internal/repository"
	"backend-koda-shortlink/internal/utils"
	"strings"

	"github.com/gin-gonic/gin"
)

type OptionalAuthMiddleware struct {
	sessionRepo *repository.SessionRepository
}

func NewOptionalAuthMiddleware(sessionRepo *repository.SessionRepository) *OptionalAuthMiddleware {
	return &OptionalAuthMiddleware{
		sessionRepo: sessionRepo,
	}
}

func (m *OptionalAuthMiddleware) OptionalAuth() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		authHeader := ctx.Request.Header.Get("Authorization")

		tokenString, found := strings.CutPrefix(authHeader, "Bearer ")
		if !found || tokenString == "" {
			ctx.Next()
			return
		}

		claims, err := utils.VerifyAccessToken(tokenString)
		if err != nil {
			ctx.Next()
			return
		}

		isActive, err := m.sessionRepo.CheckActive(ctx.Request.Context(), claims.SessionId)
		if err != nil || !isActive {
			ctx.Next()
			return
		}

		ctx.Set("userId", claims.Id)
		ctx.Set("sessionId", claims.SessionId)

		ctx.Next()
	}
}
