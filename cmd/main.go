package main

import (
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/elct9620/poc-raft-api/internal/app"
	"github.com/elct9620/poc-raft-api/internal/config"
	"github.com/elct9620/poc-raft-api/internal/server"
)

func main() {
	cfg := config.New()
	state := app.NewState()
	r, err := app.NewRaft(cfg, state)
	if err != nil {
		log.Fatal(fmt.Errorf("failed to create raft: %w", err))
	}

	server := server.NewServer(r, state)
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
