package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"sync"
	"syscall"

	"github.com/Ahu-Tools/example/config"
	"github.com/Ahu-Tools/example/edge"
)

func main() {
	config.CheckConfigs()
	if err := config.ConfigInfras(); err != nil {
		log.Fatalf("Infrastructures configuration failed. Error: %e", err)
	}

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	signalChan := make(chan os.Signal, 1)
	signal.Notify(signalChan, syscall.SIGINT, syscall.SIGTERM)

	var wg sync.WaitGroup

	edge.Start(ctx, &wg)

	<-signalChan
	log.Println("Shutdown signal received. Shutting down edges gracefully...")

	cancel()

	wg.Wait()

	log.Println("All edges have been shutted down. Exiting.")
}
