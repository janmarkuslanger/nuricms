package main

import (
	"os"

	"github.com/janmarkuslanger/nuricms"
)

func main() {
	config := &nuricms.ServerConfig{
		Port: os.Getenv("PORT"),
	}

	if config.Port == "" {
		config.Port = "8080"
	}

	nuricms.StartServer(config)
}
