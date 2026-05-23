package main

import (
	"log"
	"path/filepath"
	"runtime"

	"github.com/gin-gonic/gin"
	"github.com/zhaozeguang/timecapsule-api/internal/config"
	"github.com/zhaozeguang/timecapsule-api/internal/handler"
	"github.com/zhaozeguang/timecapsule-api/internal/middleware"
	"github.com/zhaozeguang/timecapsule-api/internal/repository"
	"github.com/zhaozeguang/timecapsule-api/internal/service"
	"github.com/zhaozeguang/timecapsule-api/pkg/database"
)

func main() {
	cfg := config.Load()

	if err := database.Connect(cfg.DatabaseURL); err != nil {
		log.Fatalf("Database connection failed: %v", err)
	}

	_, b, _, _ := runtime.Caller(0)
	migrationsDir := filepath.Join(filepath.Dir(b), "../../migrations")
	if err := database.RunMigrations(migrationsDir); err != nil {
		log.Fatalf("Migration failed: %v", err)
	}

	r := gin.Default()
	r.Use(middleware.CORS(cfg.FrontendURL))

	// Repositories
	userRepo := repository.NewUserRepo()
	spaceRepo := repository.NewSpaceRepo()
	memoryRepo := repository.NewMemoryRepo()

	// Services
	authService := service.NewAuthService(userRepo, cfg)
	spaceService := service.NewSpaceService(spaceRepo, memoryRepo)
	memoryService := service.NewMemoryService(memoryRepo)
	aiService := service.NewAIService(memoryRepo, cfg)

	// Handlers
	handler.NewHealthHandler().Register(r)
	handler.NewAuthHandler(authService, cfg.JWTSecret).Register(r)
	handler.NewSpaceHandler(spaceService, cfg.JWTSecret).Register(r)
	handler.NewMemoryHandler(memoryService, aiService, cfg.JWTSecret).Register(r)

	log.Printf("TimeCapsule API starting on :%s\n", cfg.Port)
	if err := r.Run(":" + cfg.Port); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
