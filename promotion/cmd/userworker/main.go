package main

import (
	"fmt"
	"log"
	"os"
	"os/signal"
	"time"

	"promotion/config"
	"promotion/internal/pkg/gormdb"
	"promotion/internal/usecase"
	"promotion/internal/usecase/repository"
	"promotion/internal/usecase/voucher"
	"promotion/internal/usecase/worker"

	"github.com/hibiken/asynq"
	"golang.org/x/sys/unix"
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

	// Usecase
	campaignUsecase := usecase.NewCampaignService(repository.NewCampaignRepository(db))
	campaignUserUsecase := usecase.NewCampaignUserService(repository.NewCampaignUserRepository(db))
	voucherClient := voucher.NewVoucherClient(cfg)
	// if voucherAvailable, err := util.Retry(voucherClient.Ping, 5, 5*time.Second); err != nil {
	// 	panic(err)
	// } else if !voucherAvailable {
	// 	panic(errors.New("voucher client unavailable"))
	// }

	userFirstLoginHandler := worker.NewUserTaskHandler(campaignUserUsecase, campaignUsecase, voucherClient)
	mux := asynq.NewServeMux()
	mux.Handle(worker.EventTypeUserFirstLogin, userFirstLoginHandler)

	redisAddr := fmt.Sprintf("%s:%s", cfg.Redis.Host, cfg.Redis.Port)
	srv := asynq.NewServer(
		asynq.RedisClientOpt{Addr: redisAddr, Password: cfg.Redis.Password},
		asynq.Config{Concurrency: 1},
	)

	if err := srv.Start(mux); err != nil {
		log.Fatalf("could not start server: %v", err)
	}

	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, unix.SIGTERM, unix.SIGINT, unix.SIGTSTP)
	for {
		s := <-sigs
		if s == unix.SIGTSTP {
			log.Println("Shutdown Server ...")
			srv.Stop() // stop processing new tasks
			continue
		}
		break
	}

	srv.Shutdown()
	log.Println("Server exiting")
}
