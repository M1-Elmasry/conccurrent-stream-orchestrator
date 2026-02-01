package internal

import (
	"context"
	cryptoRand "crypto/rand"
	// "fmt"
	mathRand "math/rand"
	"time"

	"m1-elmasry/conccurrent-stream-orchestrator/internal/config"
)

type DataChunk struct {
	ID        int
	StreamID  int
	Data      []byte
	CreatedAt time.Time
}

func generateRandomBytes(numberOfBytes int) []byte {
	randomSize := mathRand.Intn(numberOfBytes + 1)
	b := make([]byte, randomSize)
	cryptoRand.Read(b)
	return b
}

func Stream(ctx context.Context, id int, dataChan chan<- DataChunk) {
	chunkID := 0
	// to avoid race conditions in random number generation
	rng := mathRand.New(mathRand.NewSource(time.Now().UnixNano() + int64(id)))
	for {
		// Simulate data generation
		chunk := DataChunk{
			ID:        chunkID,
			StreamID:  id,
			Data:      generateRandomBytes(config.CHUNK_SIZE),
			CreatedAt: time.Now(),
		}
		// fmt.Printf("Stream %d generated chunk %d\n", id, chunkID)

		select {
		case <-ctx.Done():
			return
		case dataChan <- chunk:
			// fmt.Printf("chunk %d from stream %d sent Successfully\n", chunkID, id)
		default:
			// fmt.Printf("chunk %d from stream %d dropped due to channel being full\n", chunkID, id)
		}
		chunkID++
		// Simulate variable interval between data chunks
		randomInterval := rng.Intn(config.STREAM_INTERVAL + 1)
		time.Sleep(time.Duration(randomInterval) * time.Millisecond)
	}
}
