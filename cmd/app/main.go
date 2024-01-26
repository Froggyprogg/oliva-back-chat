package main

import (
	"log"
	"oliva-back-chat/config"
	"oliva-back-chat/internal/app"
)

func main() {
	if err := config.InitConfig(); err != nil {
		log.Fatalf("error initializing configs: %s", err.Error())
	}

	cfg := config.NewConfig()

	app.Run(cfg)
}
