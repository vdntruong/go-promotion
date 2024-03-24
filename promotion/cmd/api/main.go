package main

import (
	"fmt"
	"log"
	"promotion/config"
	"promotion/internal/server"
)

func main() {
	cfg, err := config.NewConfig()
	if err != nil {
		panic(err)
	}

	log.Println(fmt.Sprintf("config: %+v", cfg))
	server.Run(cfg)
}
