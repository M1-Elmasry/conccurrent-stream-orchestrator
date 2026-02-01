package main

import (
	"context"
	"fmt"
	"m1-elmasry/conccurrent-stream-orchestrator/internal"
	"m1-elmasry/conccurrent-stream-orchestrator/internal/config"
	"os"
	"os/signal"
	"sync"
	"syscall"
)

func main() {
	fmt.Println("Starting the application...")

	ctx, cancel := context.WithCancel(context.Background())

	fmt.Println("Application is running...")

	// communication channel for streams data
	dataChan := make(chan internal.DataChunk, config.BUFFER_SIZE)

	wg := sync.WaitGroup{}

	// start the streams and workers
	for workerID := range config.WORKERS_COUNT {
		wg.Add(1)
		go func(id int) {
			defer wg.Done()
			internal.Worker(ctx, id, dataChan)
		}(workerID)
	}

	fmt.Printf("%d workers started\n", config.WORKERS_COUNT)

	for streamID := range config.STREAMS_COUNT {
		wg.Add(1)
		go func(id int) {
			defer wg.Done()
			internal.Stream(ctx, id, dataChan)
		}(streamID)
	}

	fmt.Printf("%d streams started\n", config.STREAMS_COUNT)

	// signal handling and graceful shutdown logic would go here
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)

	<-sigChan
	fmt.Println("Shutdown signal received, gracefully shutting down...")

	cancel()
	wg.Wait()

	fmt.Println("Application has exited.")
}
