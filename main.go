package main

import (
	"log"
	"os"

	"github.com/ivanterekh/go-skeleton/server"
	"github.com/xlab/closer"
)

func main() {
	env := os.Getenv("ENV")
	if env == "" {
		env = "unknown"
	}
	log.Printf("running app in %s environment", env) // TODO: change logging

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	log.Printf("starting server on port %v", port) // TODO: change logging
	server.Start(":" + port)
	closer.Hold()
}
