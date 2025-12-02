package middlewares

import (
	"backend-koda-shortlink/internal/models"
	"backend-koda-shortlink/internal/utils"
	"backend-koda-shortlink/pkg/response"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

func Auth() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		authHeader := ctx.Request.Header.Get("Authorization")
		tokenString, found := strings.CutPrefix(authHeader, "Bearer ")
		if !found {
			ctx.JSON(http.StatusUnauthorized, response.ResponseError{
				Success: false,
				Message: "Authorization header required or invalid format",
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
				Message: message,
			})
			ctx.Abort()
			return
		}

		isActive, err := models.CheckSessionActive(claims.SessionId)
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, response.ResponseError{
				Success: false,
				Message: "Failed to verify session",
			})
			ctx.Abort()
			return
		}

		if !isActive {
			ctx.JSON(http.StatusUnauthorized, response.ResponseError{
				Success: false,
				Message: "Session has been terminated. Please login again",
			})
			ctx.Abort()
			return
		}

		ctx.Set("userId", claims.Id)
		ctx.Set("sessionId", claims.SessionId)

		ctx.Next()
	}
}
