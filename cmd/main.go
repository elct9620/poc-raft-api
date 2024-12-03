package main

import (
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/elct9620/poc-raft-api/internal/app"
	"github.com/elct9620/poc-raft-api/internal/server"
)

func main() {
	r, err := app.NewRaft(
		os.Getenv("HOSTNAME"),
		"/data",
		os.Getenv("RAFT_ADDRESS"),
	)
	if err != nil {
		log.Fatal(fmt.Errorf("failed to create raft: %w", err))
	}

	server := server.NewServer(r)

	signalChan := make(chan os.Signal, 1)
	signal.Notify(signalChan, os.Interrupt, syscall.SIGTERM)

	go func() {
		<-signalChan
		if err := server.Stop(); err != nil {
			log.Fatal(fmt.Errorf("failed to stop server: %w", err))
		}
	}()

	if err := server.Start(); err != nil {
		log.Fatal(fmt.Errorf("failed to start server: %w", err))
	}
}
