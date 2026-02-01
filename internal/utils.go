package internal

import (
	"fmt"
	"time"
)

func MeasureLatency[T any](chunk DataChunk, calledFunc func() T) T {
	result := calledFunc()
	latency := time.Since(chunk.CreatedAt)
	if latency > 100*time.Millisecond {
		fmt.Printf("⚠️  HIGH LATENCY: %v (chunk %d from stream %d)\n", latency, chunk.ID, chunk.StreamID)
	} else {
		fmt.Printf("✅ Latency: %v (chunk %d from stream %d)\n", latency, chunk.ID, chunk.StreamID)
	}
	return result
}
