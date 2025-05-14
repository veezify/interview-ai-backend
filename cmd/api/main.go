package main

import (
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"github.com/veezify/interview-ai-backend/internal/api/middleware"
	"github.com/veezify/interview-ai-backend/internal/api/router"
	"github.com/veezify/interview-ai-backend/internal/database"
	"github.com/veezify/interview-ai-backend/internal/repository"
	"github.com/veezify/interview-ai-backend/internal/service"
	"github.com/veezify/interview-ai-backend/pkg/logger"
	"go.uber.org/zap"
)

// @title Interview AI API
// @version 1.0
// @description API Server for Interview AI Application
// @termsOfService http://swagger.io/terms/

// @contact.name API Support
// @contact.url http://www.swagger.io/support
// @contact.email support@swagger.io

// @license.name Apache 2.0
// @license.url http://www.apache.org/licenses/LICENSE-2.0.html

// @host localhost:8080
// @BasePath /api
// @schemes http https
func main() {
	if err := godotenv.Load(); err != nil {
		log.Printf("Warning: .env file not found or cannot be loaded. Using environment variables.")
	}

	env := os.Getenv("ENV")
	if env == "" {
		env = "development"
	}

	logger.Init(env)
	defer logger.Sync()

	logger.Info("Starting application",
		zap.String("environment", env),
	)

	if env == "production" {
		gin.SetMode(gin.ReleaseMode)
	} else {
		gin.SetMode(gin.DebugMode)
	}

	db, err := database.SetupDB()
	if err != nil {
		logger.Fatal("Failed to setup database", zap.Error(err))
	}

	userRepo := repository.NewUserRepository(db)
	interviewRepo := repository.NewInterviewRepository(db)

	authService := service.NewAuthService(userRepo)
	interviewService := service.NewInterviewService(interviewRepo, "uploads/cv")
	userService := service.NewUserService(userRepo)

	r := router.SetupRouter(authService, interviewService, userService)

	r.Use(middleware.Logger())

	port := os.Getenv("SERVER_PORT")
	if port == "" {
		port = "8080"
	}

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		logger.Info("Server starting", zap.String("port", port))
		if err := r.Run(":" + port); err != nil {
			logger.Fatal("Failed to start server", zap.Error(err))
		}
	}()

	<-quit
	logger.Info("Shutting down server...")
	logger.Info("Server shutdown complete")
}
