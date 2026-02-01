package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"sync"
	"syscall"

	"github.com/m1-elmasry/conccurrent-stream-orchestrator/internal"
	"github.com/m1-elmasry/conccurrent-stream-orchestrator/internal/config"
	"github.com/m1-elmasry/conccurrent-stream-orchestrator/internal/metrics"
)

func main() {
	fmt.Println("Starting the application...")

	// Load configuration from environment variables
	cfg := config.LoadConfig()
	cfg.PrintConfig()

	ctx, cancel := context.WithCancel(context.Background())

	fmt.Println("Application is running...")

	m := metrics.New()
	dataChan := make(chan internal.DataChunk, cfg.BufferSize)
	wg := sync.WaitGroup{}

	// start workers
	for workerID := range cfg.WorkersCount {
		wg.Add(1)
		go func(id int) {
			defer wg.Done()
			internal.Worker(ctx, id, dataChan, m)
		}(workerID)
	}

	fmt.Printf("%d workers started\n", cfg.WorkersCount)

	// start streams
	for streamID := range cfg.StreamsCount {
		wg.Add(1)
		go func(id int) {
			defer wg.Done()
			internal.Stream(ctx, id, dataChan, m, cfg)
		}(streamID)
	}

	fmt.Printf("%d streams started\n", cfg.StreamsCount)

	// wait for shutdown signal
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)

	<-sigChan
	fmt.Println("\nShutdown signal received...")

	cancel()
	wg.Wait()

	// print final stats
	fmt.Println("\n=== Final Statistics ===")
	stats := m.Summary()
	fmt.Printf("Processed: %d chunks\n", stats["processed"])
	fmt.Printf("Dropped: %d chunks\n", stats["dropped"])
	fmt.Printf("Elapsed: %v\n", stats["elapsed"])
	fmt.Printf("Throughput: %.2f chunks/sec\n", stats["throughput"])
	fmt.Printf("Avg Latency: %v\n", stats["avg_latency"])
	fmt.Printf("P50 Latency: %v\n", stats["p50"])
	fmt.Printf("P95 Latency: %v\n", stats["p95"])
	fmt.Printf("P99 Latency: %v\n", stats["p99"])
	fmt.Println("Application exited.")
}
