package main

import (
	"fmt"
	"log"

	"voucher/config"
	"voucher/internal/server"
)

func main() {
	cfg, err := config.NewConfig()
	if err != nil {
		panic(err)
	}

	log.Println(fmt.Sprintf("config: %+v", cfg))
	server.Run(cfg)
}
