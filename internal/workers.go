package internal

import (
	"context"
	"fmt"
	"github.com/m1-elmasry/conccurrent-stream-orchestrator/internal/metrics"
	"math/rand"
	"time"
)

func ProcessChunk(chunk DataChunk) int {
	time.Sleep(time.Duration(20+rand.Intn(30)) * time.Millisecond)
	return len(chunk.Data)
}

func Worker(ctx context.Context, id int, dataChan <-chan DataChunk, m *metrics.Metrics) {
	for {
		select {
		case <-ctx.Done():
			return
		case chunk := <-dataChan:
			result := metrics.MeasureLatency(chunk.CreatedAt, m, func() int {
				return ProcessChunk(chunk)
			})
			fmt.Printf("Worker %d processed chunk %d from stream %d, length of chunk's data is %d\n", id, chunk.ID, chunk.StreamID, result)
		}
	}
}
