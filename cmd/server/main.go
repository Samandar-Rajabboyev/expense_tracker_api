package main

import (
	"log"
	"os"

	"expense-tracker-api/internal/cache"
	"expense-tracker-api/internal/config"
	"expense-tracker-api/internal/database"
	"expense-tracker-api/internal/handler"
	"expense-tracker-api/internal/jwt"
	"expense-tracker-api/internal/middleware"
	"expense-tracker-api/internal/repository"
	"expense-tracker-api/internal/service"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

func main() {
	godotenv.Load()

	cfg, _ := config.Load()

	db, err := database.NewDB(cfg.DatabaseURL)
	if err != nil {
		log.Fatalf("failed to connect to database: %v", err)
	}

	redisCache := cache.NewCache(cfg.RedisURL)

	userRepo := repository.NewUserRepository(db)
	expenseRepo := repository.NewExpenseRepository(db)

	userService := service.NewUserService(userRepo)
	expenseService := service.NewExpenseService(expenseRepo, redisCache)
	jwtService := jwt.NewJWTService(cfg.JWTSecret)

	userHandler := handler.NewUserHandler(userService, jwtService)
	expenseHandler := handler.NewExpenseHandler(expenseService)

	authMiddleware := middleware.NewAuthMiddleware(jwtService)

	r := gin.Default()

	corsConfig := cors.DefaultConfig()
	corsConfig.AllowOrigins = []string{"http://localhost:3000", "http://localhost:5173"}
	corsConfig.AllowMethods = []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"}
	corsConfig.AllowHeaders = []string{"Origin", "Content-Type", "Authorization"}
	r.Use(cors.New(corsConfig))

	r.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{"status": "ok"})
	})

	r.POST("/users/register", userHandler.Register)
	r.POST("/users/login", userHandler.Login)
	r.POST("/auth/refresh", userHandler.Refresh)

	expenses := r.Group("/expenses")
	expenses.Use(authMiddleware.Authenticate())
	{
		expenses.POST("", expenseHandler.Create)
		expenses.GET("", expenseHandler.GetAll)
		expenses.GET("/:id", expenseHandler.GetByID)
		expenses.PUT("/:id", expenseHandler.Update)
		expenses.DELETE("/:id", expenseHandler.Delete)
	}

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	log.Printf("Server starting on :%s", port)
	if err := r.Run(":" + port); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
