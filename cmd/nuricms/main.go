package main

import (
	"os"

	"github.com/janmarkuslanger/nuricms"
	"github.com/janmarkuslanger/nuricms/pkg/config"
)

func main() {
	config := config.Config{
		Port: os.Getenv("PORT"),
	}

	if config.Port == "" {
		config.Port = "8080"
	}

	nuricms.Run(config)
}
