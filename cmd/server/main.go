package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

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

	r.GET("/ready", func(c *gin.Context) {
		// Check DB
		sqlDB, err := db.DB()
		if err != nil {
			c.JSON(503, gin.H{"status": "unhealthy", "error": "db unavailable"})
			return
		}
		if err := sqlDB.Ping(); err != nil {
			c.JSON(503, gin.H{"status": "unhealthy", "error": "db ping failed"})
			return
		}
		// Check Redis
		ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
		defer cancel()
		if err := redisCache.Ping(ctx); err != nil {
			c.JSON(503, gin.H{"status": "unhealthy", "error": "redis unavailable"})
			return
		}
		c.JSON(200, gin.H{"status": "ready"})
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

	srv := &http.Server{
		Addr:    ":" + port,
		Handler: r,
	}

	go func() {
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Failed to start server: %v", err)
		}
	}()

	log.Printf("Server starting on :%s", port)

	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()
	<-ctx.Done()

	log.Println("Shutting down gracefully...")

	shutdownCtx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := srv.Shutdown(shutdownCtx); err != nil {
		log.Fatalf("Server forced to shutdown: %v", err)
	}

	log.Println("Server exited properly")
}
