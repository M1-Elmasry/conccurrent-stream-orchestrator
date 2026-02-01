package internal

import (
	"context"
	"fmt"
	"math/rand"
	"time"
)

func ProcessChunk(chunk DataChunk) int {
	// simulate 20-50ms API Call. simulate "an external processing engine"
	time.Sleep(time.Duration(20+rand.Intn(30)) * time.Millisecond)
	return len(chunk.Data)
}

func Worker(ctx context.Context, id int, dataChan <-chan DataChunk) {
	for {
		select {
		case <-ctx.Done():
			return
		case chunk := <-dataChan:

			// fmt.Printf("Worker %d received chunk %d from stream %d for processing\n", id, chunk.ID, chunk.StreamID)

			// Measure latency and process the chunk

			result := MeasureLatency(chunk, func() int {
				return ProcessChunk(chunk)
			})

			// Log the processed chunk
			fmt.Printf("Worker %d processed chunk %d from stream %d, Length of chunk's data %d\n", id, chunk.ID, chunk.StreamID, result)
		}
	}
}
