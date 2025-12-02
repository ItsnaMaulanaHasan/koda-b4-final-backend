package middlewares

import (
	"log"
	"time"

	"github.com/gin-gonic/gin"
)

func RequestLogger() gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()

		path := c.Request.URL.Path
		method := c.Request.Method
		ip := c.ClientIP()

		c.Next()
		status := c.Writer.Status()
		latency := time.Since(start)

		log.Printf("[REQUEST] %s %s | %d | %s | %s",
			method, path, status, latency, ip)
	}
}
