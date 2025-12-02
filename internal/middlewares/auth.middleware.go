package middlewares

import (
	"backend-koda-shortlink/internal/repository"
	"backend-koda-shortlink/internal/utils"
	"backend-koda-shortlink/pkg/response"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

type AuthMiddleware struct {
	sessionRepo *repository.SessionRepository
}

func NewAuthMiddleware(sessionRepo *repository.SessionRepository) *AuthMiddleware {
	return &AuthMiddleware{
		sessionRepo: sessionRepo,
	}
}

func (m *AuthMiddleware) Auth() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		authHeader := ctx.Request.Header.Get("Authorization")
		tokenString, found := strings.CutPrefix(authHeader, "Bearer ")
		if !found {
			ctx.JSON(http.StatusUnauthorized, response.ResponseError{
				Success: false,
				Error:   "Authorization header required or invalid format",
			})
			ctx.Abort()
			return
		}

		claims, err := utils.VerifyAccessToken(tokenString)
		if err != nil {
			message := "Invalid or expired token"
			switch err {
			case jwt.ErrTokenExpired:
				message = "Token expired. Please refresh your token"
			case jwt.ErrSignatureInvalid:
				message = "Invalid token signature"
			}

			ctx.JSON(http.StatusUnauthorized, response.ResponseError{
				Success: false,
				Error:   message,
			})
			ctx.Abort()
			return
		}

		isActive, err := m.sessionRepo.CheckActive(ctx.Request.Context(), claims.SessionId)
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, response.ResponseError{
				Success: false,
				Error:   "Failed to verify session",
			})
			ctx.Abort()
			return
		}

		if !isActive {
			ctx.JSON(http.StatusUnauthorized, response.ResponseError{
				Success: false,
				Error:   "Session has been terminated. Please login again",
			})
			ctx.Abort()
			return
		}

		ctx.Set("userId", claims.Id)
		ctx.Set("sessionId", claims.SessionId)

		ctx.Next()
	}
}
