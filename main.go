package main

// @title           API Koda Shortlink Documentation
// @version         1.0
// @description     Dokumentasi REST API menggunakan Gin dan Swagger

// @BasePath  /api/v1

// @securityDefinitions.apikey BearerAuth
// @in header
// @name Authorization

import (
	"backend-koda-shortlink/internal/config"
	"backend-koda-shortlink/internal/database"
	"backend-koda-shortlink/internal/middlewares"
	"backend-koda-shortlink/internal/routes"
	"backend-koda-shortlink/pkg/response"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"

	_ "backend-koda-shortlink/docs"

	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

func main() {
	godotenv.Load()
	database.InitDatabase()
	config.InitRedis()

	r := gin.Default()
	r.Use(gin.Recovery())
	r.Use(middlewares.CorsMiddleware())
	r.Use(middlewares.RateLimiter(60, time.Minute))
	r.Use(middlewares.RequestLogger())

	r.GET("/", func(ctx *gin.Context) {
		ctx.JSON(http.StatusOK, response.ResponseSuccess{
			Success: true,
			Message: "Backend is running well!",
		})
	})

	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	routes.SetUpRoutes(r)

	r.Run(":8080")
}
