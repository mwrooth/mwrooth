package main

import (
	"log"

	"github.com/joho/godotenv"
	"github.com/mwrooth/gateway/pkg/api/general"
	"github.com/mwrooth/gateway/pkg/service"
)

func main() {
	if err := godotenv.Load(); err != nil {
		log.Fatalf("ENV: %+v", err)
	}

	s := service.New()
	general.Start(s)
}
