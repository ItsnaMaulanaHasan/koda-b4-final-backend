package middlewares

import (
	"backend-koda-shortlink/pkg/response"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
)

func RateLimiter(limit int, window time.Duration) gin.HandlerFunc {
	var mu sync.Mutex
	tokens := limit
	lastRefill := time.Now()

	return func(c *gin.Context) {
		mu.Lock()
		now := time.Now()

		if now.Sub(lastRefill) >= window {
			tokens = limit
			lastRefill = now
		}

		if tokens > 0 {
			tokens--
			mu.Unlock()
			c.Next()
		} else {
			mu.Unlock()
			c.JSON(429, response.ResponseError{
				Success: false,
				Error:   "Too many requests",
			})
			c.Abort()
		}
	}
}
