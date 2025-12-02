package middlewares

import (
	"backend-koda-shortlink/internal/config"
	"backend-koda-shortlink/pkg/response"
	"time"

	"github.com/gin-gonic/gin"
)

func RateLimiter(limit int, window time.Duration) gin.HandlerFunc {
	return func(c *gin.Context) {
		ip := c.ClientIP()
		path := c.FullPath()

		key := "ratelimit:" + ip + ":" + path

		count, _ := config.Rdb.Incr(c, key).Result()

		if count == 1 {
			config.Rdb.Expire(c, key, window)
		}

		if count > int64(limit) {
			c.JSON(429, response.ResponseError{
				Success: false,
				Error:   "Too many requests",
			})
			c.Abort()
			return
		}

		c.Next()
	}
}
