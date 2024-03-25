package main

import (
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"promotion/config"
	"promotion/internal/model"
	"promotion/internal/pkg/gormdb"
	"promotion/internal/usecase"
	"promotion/internal/usecase/repository"
	"promotion/internal/usecase/worker"

	"github.com/hibiken/asynq"
)

func main() {
	cfg, err := config.NewConfig()
	if err != nil {
		panic(err)
	}

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

	mux := asynq.NewServeMux()

	// Usecase
	campaignUsecase := usecase.NewCampaignService(repository.NewCampaignRepository(db))
	campaignUserUsecase := usecase.NewCampaignUserService(repository.NewCampaignUserRepository(db))

	// Handler
	userFirstLoginHandler := worker.NewUserTaskHandler(campaignUserUsecase, campaignUsecase)
	mux.Handle(worker.EventTypeUserFirstLogin, userFirstLoginHandler)

	// Serve
	redisAddr := fmt.Sprintf("%s:%s", cfg.Redis.Host, cfg.Redis.Port)
	srv := asynq.NewServer(
		asynq.RedisClientOpt{Addr: redisAddr, Password: cfg.Redis.Password},
		asynq.Config{
			Concurrency: 1,
			// Queues: map[string]int{
			// 	"critical": 6,
			// 	"default":  3,
			// 	"low":      1,
			// },
		})

	go func() {
		if err := srv.Run(mux); err != nil {
			log.Fatalf("could not run server: %v", err)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	log.Println("Shutdown Server ...")

	srv.Shutdown()
	log.Println("Server exiting")
}
