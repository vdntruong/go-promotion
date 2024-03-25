package main

import (
	"ekyc/config"
	"ekyc/internal/server"
)

func main() {
	cfg, err := config.NewConfig()
	if err != nil {
		panic(err)
	}

	server.Run(cfg)
}
