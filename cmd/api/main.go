package main

import (
	"context"
	"flag"
	"log"
	"strconv"
	"time"

	"ai-listen/internal/config"
	"ai-listen/internal/handler"
	"ai-listen/internal/middleware"
	"ai-listen/internal/model"
	"ai-listen/internal/repository"
	"ai-listen/internal/router"
	"ai-listen/internal/service"
	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

func main() {
	configPath := flag.String("config", "config.yaml", "config file path")
	flag.Parse()

	cfg, err := config.Load(*configPath)
	if err != nil {
		log.Fatalf("load config failed: %v", err)
	}

	if cfg.App.Env == "prod" {
		gin.SetMode(gin.ReleaseMode)
	}

	db, err := gorm.Open(mysql.Open(cfg.Database.DSN()), &gorm.Config{})
	if err != nil {
		log.Fatalf("connect mysql failed: %v", err)
	}

	if err := db.AutoMigrate(&model.User{}, &model.ProviderApplication{}, &model.Provider{}); err != nil {
		log.Fatalf("migrate database failed: %v", err)
	}

	rdb := redis.NewClient(&redis.Options{
		Addr:     cfg.Redis.Addr,
		Password: cfg.Redis.Password,
		DB:       cfg.Redis.DB,
	})
	if err := rdb.Ping(context.Background()).Err(); err != nil {
		log.Fatalf("connect redis failed: %v", err)
	}

	userRepo := repository.NewUserRepository(db)
	providerApplicationRepo := repository.NewProviderApplicationRepository(db)
	providerRepo := repository.NewProviderRepository(db)

	authService := service.NewAuthService(
		userRepo,
		rdb,
		time.Duration(cfg.Auth.CodeTTLMinutes)*time.Minute,
		time.Duration(cfg.Auth.TokenTTLDays)*24*time.Hour,
		cfg.App.Env,
	)
	userService := service.NewUserService(userRepo)
	providerCenterService := service.NewProviderCenterService(userRepo, providerApplicationRepo, providerRepo)

	authHandler := handler.NewAuthHandler(authService)
	userHandler := handler.NewUserHandler(userService)
	providerCenterHandler := handler.NewProviderCenterHandler(providerCenterService)

	engine := gin.New()
	engine.Use(gin.Logger(), middleware.ErrorHandler(), middleware.Recovery())

	router.RegisterAPIRoutes(
		engine,
		authHandler,
		userHandler,
		providerCenterHandler,
		middleware.AuthRequired(rdb),
	)

	addr := ":" + strconv.Itoa(cfg.Server.APIPort)
	log.Printf("api server started at %s", addr)
	if err := engine.Run(addr); err != nil {
		log.Fatalf("run api failed: %v", err)
	}
}
