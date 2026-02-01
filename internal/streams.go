package internal

import (
	"context"
	cryptoRand "crypto/rand"
	mathRand "math/rand"
	"time"

	"github.com/m1-elmasry/conccurrent-stream-orchestrator/internal/config"
	"github.com/m1-elmasry/conccurrent-stream-orchestrator/internal/metrics"
)

type DataChunk struct {
	ID        int
	StreamID  int
	Data      []byte
	CreatedAt time.Time
}

func generateRandomBytes(numberOfBytes int, rng *mathRand.Rand) []byte {
	randomSize := rng.Intn(numberOfBytes + 1)
	b := make([]byte, randomSize)
	cryptoRand.Read(b)
	return b
}

func Stream(ctx context.Context, id int, dataChan chan<- DataChunk, m *metrics.Metrics, cfg *config.Configuration) {
	chunkID := 0
	rng := mathRand.New(mathRand.NewSource(time.Now().UnixNano() + int64(id)))
	for {
		chunk := DataChunk{
			ID:        chunkID,
			StreamID:  id,
			Data:      generateRandomBytes(cfg.ChunkSize, rng),
			CreatedAt: time.Now(),
		}

		select {
		case <-ctx.Done():
			return
		case dataChan <- chunk:
			// success
		default:
			m.RecordDropped()
		}
		chunkID++
		randomInterval := rng.Intn(cfg.StreamInterval + 1)
		time.Sleep(time.Duration(randomInterval) * time.Millisecond)
	}
}
