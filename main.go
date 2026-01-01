package main

import (
	"context"
	"os"
	"os/signal"
	"sync"
	"syscall"

	"github.com/Ahu-Tools/example/config"
	"github.com/Ahu-Tools/example/edge"
	"github.com/Ahu-Tools/example/log"
)

func main() {
	config.CheckConfigs()
	if err := config.ConfigInfras(); err != nil {
		log.Logger.Error("Infrastructures configuration failed.", "error", err)
		os.Exit(1)
	}

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	signalChan := make(chan os.Signal, 1)
	signal.Notify(signalChan, syscall.SIGINT, syscall.SIGTERM)

	var wg sync.WaitGroup

	edge.Start(ctx, &wg)

	<-signalChan
	log.Logger.Info("Shutdown signal received. Shutting down edges gracefully...")

	cancel()

	wg.Wait()

	log.Logger.Info("All edges have been shutted down. Exiting.")
}
