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
	"log"
	"net/http"
	"os"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-migrate/migrate/v4"
	"github.com/joho/godotenv"

	_ "backend-koda-shortlink/docs"

	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"

	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

func main() {
	godotenv.Load()

	runMigrations()

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

func runMigrations() {
	databaseURL := os.Getenv("DATABASE_URL")
	if databaseURL == "" {
		log.Fatal("DATABASE_URL is not set")
	}

	m, err := migrate.New(
		"file://migrations",
		databaseURL,
	)
	if err != nil {
		log.Fatalf("Failed to create migrate instance: %v", err)
	}
	defer m.Close()

	if err := m.Up(); err != nil {
		if err == migrate.ErrNoChange {
			log.Println("No new migrations to apply")
		} else {
			log.Fatalf("Failed to run migrations: %v", err)
		}
	} else {
		log.Println("Migrations applied successfully")
	}
}
