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

	"promotion/config"
	"promotion/internal/model"
	"promotion/internal/pkg/gormdb"
	healthHttp "promotion/internal/server/health/delivery/http"
	v1 "promotion/internal/server/v1"
	"promotion/internal/usecase"
	"promotion/internal/usecase/repository"

	"github.com/gin-contrib/cors"
	"github.com/gin-contrib/requestid"
	limits "github.com/gin-contrib/size"
	"github.com/gin-gonic/gin"
)

func Run(cfg *config.Config) {
	// Infrastructure
	db, err := gormdb.DB(cfg)
	if err != nil {
		panic(err)
	}

	sqlDB, err := db.DB()
	if err != nil {
		panic(err)
	}
	defer sqlDB.Close()

	sqlDB.SetMaxIdleConns(10)
	sqlDB.SetMaxOpenConns(100)
	sqlDB.SetConnMaxLifetime(time.Hour)

	db.AutoMigrate(&model.Campaign{})
	db.AutoMigrate(&model.CampaignUser{})

	handler := gin.New()
	handler.Use(cors.Default())
	handler.Use(requestid.New())
	handler.Use(limits.RequestSizeLimiter(500))
	handler.Use(gin.Logger())
	handler.Use(gin.Recovery())
	healthHandlers := healthHttp.NewHealthHandlers(cfg)
	healthHttp.MapHealthRoutes(handler, healthHandlers)

	// Usecase
	campaignUsecase := usecase.NewCampaignService(repository.NewCampaignRepository(db))
	campaignUserUsecase := usecase.NewCampaignUserService(repository.NewCampaignUserRepository(db))
	// API - Version 1
	v1.NewRouter(handler, campaignUsecase, campaignUserUsecase)

	srv := http.Server{
		Handler:      handler,
		Addr:         fmt.Sprintf(":%s", cfg.HTTP.Port),
		ReadTimeout:  time.Second * cfg.HTTP.ReadTimeout,
		WriteTimeout: time.Second * cfg.HTTP.WriteTimeout,
	}

	go func() {
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("listen: %s\n", err)
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
			log.Fatalf("Server Shutdown: %v", err)
		}
	}

	log.Println("Server exiting")
}
