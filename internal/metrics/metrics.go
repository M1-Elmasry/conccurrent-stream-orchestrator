package metrics

import (
	"sort"
	"sync"
	"time"
)

// Metrics holds all system metrics
type Metrics struct {
	mu        sync.RWMutex
	processed int64
	dropped   int64
	latencies []time.Duration
	startTime time.Time
}

func New() *Metrics {
	return &Metrics{
		latencies: make([]time.Duration, 0),
		startTime: time.Now(),
	}
}

func (m *Metrics) RecordProcessed(latency time.Duration) {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.processed++
	m.latencies = append(m.latencies, latency) // unbounded growth, consider using a ring buffer for production
}

func (m *Metrics) RecordDropped() {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.dropped++
}

func (m *Metrics) Summary() map[string]interface{} {
	m.mu.RLock()
	defer m.mu.RUnlock()

	elapsed := time.Since(m.startTime)
	result := map[string]interface{}{
		"processed":  m.processed,
		"dropped":    m.dropped,
		"elapsed":    elapsed,
		"throughput": 0.0,
	}

	if elapsed.Seconds() > 0 {
		result["throughput"] = float64(m.processed) / elapsed.Seconds()
	}

	if len(m.latencies) > 0 {
		sorted := make([]time.Duration, len(m.latencies))
		copy(sorted, m.latencies)
		sort.Slice(sorted, func(i, j int) bool { return sorted[i] < sorted[j] })

		var total time.Duration
		for _, lat := range sorted {
			total += lat
		}

		result["avg_latency"] = total / time.Duration(len(sorted))
		result["p50"] = sorted[len(sorted)/2]
		result["p95"] = sorted[int(float64(len(sorted))*0.95)]
		result["p99"] = sorted[int(float64(len(sorted))*0.99)]
	}

	return result
}

// MeasureLatency measures the latency of a function execution and records it in metrics.
// DataChunk needs to be imported from the parent package to use this function.
func MeasureLatency[T any](createdAt time.Time, m *Metrics, calledFunc func() T) T {
	result := calledFunc()
	latency := time.Since(createdAt)
	m.RecordProcessed(latency)
	return result
}
