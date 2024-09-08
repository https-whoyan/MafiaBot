package main

import (
	"context"
	"github.com/https-whoyan/MafiaBot/internal/app"
	"github.com/https-whoyan/MafiaBot/internal/app/config"
	"log"
)

func main() {
	ctx := context.Background()
	cfg := config.LoadConfig()
	apl, err := app.InitApp(ctx, cfg)
	if err != nil {
		log.Fatal(err)
	}
	err = apl.Run(ctx)
	if err != nil {
		log.Fatal(err)
	}
}
