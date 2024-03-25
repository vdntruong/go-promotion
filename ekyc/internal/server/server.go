package server

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"ekyc/config"
	"ekyc/internal/model"
	"ekyc/internal/pkg/ejwt"
	"ekyc/internal/pkg/gormdb"
	healthHttp "ekyc/internal/server/health/delivery/http"
	v1 "ekyc/internal/server/v1"
	"ekyc/internal/usecase"
	"ekyc/internal/usecase/distributor"
	"ekyc/internal/usecase/repository"

	"github.com/gin-contrib/cors"
	"github.com/gin-contrib/requestid"
	limits "github.com/gin-contrib/size"
	"github.com/gin-gonic/gin"
	"github.com/hibiken/asynq"
)

func Run(cfg *config.Config) {
	// Infrastructure
	db, err := gormdb.DB(cfg)
	if err != nil {
		panic(err)
	}

	db.AutoMigrate(&model.User{})

	sqlDB, err := db.DB()
	if err != nil {
		panic(err)
	}
	defer sqlDB.Close()

	sqlDB.SetMaxIdleConns(10)
	sqlDB.SetMaxOpenConns(100)
	sqlDB.SetConnMaxLifetime(time.Hour)

	redisAddr := fmt.Sprintf("%s:%s", cfg.Redis.Host, cfg.Redis.Port)
	asynqClient := asynq.NewClient(asynq.RedisClientOpt{Addr: redisAddr, Password: cfg.Redis.Password})
	defer asynqClient.Close()

	handler := gin.New()
	handler.Use(cors.Default())
	handler.Use(requestid.New())
	handler.Use(limits.RequestSizeLimiter(500))
	handler.Use(gin.Logger())
	handler.Use(gin.Recovery())
	healthHandlers := healthHttp.NewHealthHandlers(cfg)
	healthHttp.MapHealthRoutes(handler, healthHandlers)

	// Usecase
	userUsecase := usecase.NewUserService(
		repository.NewUserRepository(db),
		usecase.NewAuth(ejwt.NewGenerator(cfg.JWT.Issuer, cfg.JWT.Secret)),
		distributor.NewUserDistributor(asynqClient),
	)
	// API - Version 1
	v1.NewRouter(handler, userUsecase)

	srv := http.Server{
		Handler:      handler,
		Addr:         fmt.Sprintf(":%s", cfg.HTTP.Port),
		ReadTimeout:  time.Second * cfg.HTTP.ReadTimeout,
		WriteTimeout: time.Second * cfg.HTTP.WriteTimeout,
	}

	go func() {
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("failed listen: %s\n", err)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	log.Println("Shutdown Server ...")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		if err == http.ErrServerClosed {
			log.Println("Server shutdown completed")
		} else {
			log.Fatal("Server Shutdown:", err)
		}
	}

	log.Println("Server exiting")
}
