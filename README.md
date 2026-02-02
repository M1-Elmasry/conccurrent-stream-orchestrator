# Concurrent Stream Orchestrator

A high-performance streaming gateway that processes 100 concurrent data streams with sub-100ms latency, demonstrating advanced Go concurrency patterns and efficient resource management.

## Overview

This project implements a production-grade stream processing system that handles multiple real-time data sources using worker pool patterns, comprehensive metrics tracking, and graceful shutdown capabilities.

**Key Features:**
- ✅ 100 concurrent data streams
- ✅ Worker pool with 50 workers for optimal throughput
- ✅ Sub-100ms latency per chunk (20-50ms average)
- ✅ Comprehensive metrics (P50/P95/P99 percentiles)
- ✅ Graceful shutdown with no data loss
- ✅ Backpressure handling with drop tracking
- ✅ Thread-safe implementation

---

## Project Structure

```
conccurrent-stream-orchestrator/
├── cmd/
│   └── main.go                 # Application entry point & orchestration
├── internal/
│   ├── config/
│   │   └── config.go           # Configuration constants (workers, streams, buffer size)
│   ├── metrics/
│   │   └── metrics.go          # Metrics collection & percentile calculations
│   ├── streams.go              # Stream generator logic
│   └── workers.go              # Worker pool processing logic
├── docker-compose.yaml         # Docker composition (if applicable)
├── Dockerfile                  # Container image definition
├── go.mod                       # Go module dependencies
├── go.sum                       # Dependency checksums
├── main                         # Compiled binary
└── README.md                    # This file
```

### Key Files

- **cmd/main.go**: Orchestrates the entire application. Initializes context, metrics, data channel, starts 100 stream goroutines and 50 worker goroutines, and handles graceful shutdown.
- **internal/config/config.go**: Centralizes configuration (50 workers, 100 streams, 500 buffer size, 1KB chunks).
- **internal/streams.go**: `Stream()` function generates random data chunks at variable intervals (0-1000ms) and sends them non-blocking to the channel.
- **internal/workers.go**: `Worker()` and `ProcessChunk()` functions consume chunks from the channel, simulate processing (20-50ms), and record latency metrics.
- **internal/metrics/metrics.go**: Thread-safe metrics collection with percentile calculations (P50, P95, P99).

---

## Getting Started

### Prerequisites

**Option 1: Local Development**
- Go 1.20 or higher
- Unix/Linux system (for signal handling)

**Option 2: Docker**
- Docker 20.10+
- Docker Compose 2.0+ (optional)

### Quick Start - Local Build

1. **Build the application:**
   ```bash
   go build -o main ./cmd/main.go
   ```

2. **Run the application (with defaults):**
   ```bash
   ./main
   ```

3. **Run with custom configuration:**
   ```bash
   WORKERS_COUNT=25 STREAMS_COUNT=50 STREAM_INTERVAL=2000 ./main
   ```

4. **Graceful shutdown:**
   ```bash
   # Press Ctrl+C in the terminal
   ```

### Quick Start - Docker

#### Option A: Docker Compose (Recommended)

Run with Docker Compose for easier management:

look on the expected output section below 

```bash
# Start the application (with defaults, do not run in detached mode)
docker-compose up
```

for graceful shutdown 

```bash
# Press Ctrl+C in the terminal
```

```bash
# Stop the application
docker-compose down
```

#### Option B: Docker Run

Build and run using Docker directly:

```bash
# Build the Docker image
docker build -t concurrent-stream-orchestrator:latest .

# Run the container (with defaults)
docker run --rm concurrent-stream-orchestrator:latest

# For interactive mode (to stop with Ctrl+C)
docker run --rm -it concurrent-stream-orchestrator:latest

# Run with custom environment variables
docker run --rm -it \
  -e WORKERS_COUNT=25 \
  -e STREAMS_COUNT=50 \
  -e STREAM_INTERVAL=2000 \
  -e CHUNK_SIZE=2048 \
  -e BUFFER_SIZE=250 \
  concurrent-stream-orchestrator:latest
```



### Environment Variables

The application supports the following environment variables:

| Variable | Default | Range | Description |
|----------|---------|-------|-------------|
| `WORKERS_COUNT` | 50 | 10-200 | Number of worker goroutines |
| `STREAMS_COUNT` | 100 | 10-500 | Number of data streams |
| `STREAM_INTERVAL` | 1000 | 100-5000 | Max interval between chunks (ms) |
| `CHUNK_SIZE` | 1024 | 512-4096 | Size of each chunk (bytes) |
| `BUFFER_SIZE` | 500 | 100-2000 | Channel buffer capacity |

#### Using Environment Variables with Docker Compose

Edit `docker-compose.yaml` to override defaults:

```yaml
environment:
  WORKERS_COUNT: "25"
  STREAMS_COUNT: "50"
  STREAM_INTERVAL: "2000"
  CHUNK_SIZE: "2048"
  BUFFER_SIZE: "250"
```

Or load from `.env` file (copy `.env.example` to `.env`):

```yaml
env_file:
  - .env
```

Then start: `docker-compose up`

#### Example Configurations

**Development (Low Resource):**
```bash
WORKERS_COUNT=10 STREAMS_COUNT=20 STREAM_INTERVAL=2000 ./main
```

**Production (High Throughput):**
```bash
WORKERS_COUNT=100 STREAMS_COUNT=200 STREAM_INTERVAL=500 CHUNK_SIZE=2048 BUFFER_SIZE=1000 ./main
```

**Performance Testing (Max Load):**
```bash
WORKERS_COUNT=150 STREAMS_COUNT=300 STREAM_INTERVAL=100 CHUNK_SIZE=4096 BUFFER_SIZE=2000 ./main
```

### Docker Configuration

The Docker setup includes:
- **Multi-stage build**: Reduces final image size from ~900MB to ~30MB
- **Alpine base**: Lightweight runtime environment
- **Environment defaults**: All variables have sensible defaults in Dockerfile
- **Resource limits**: CPU (2 cores), Memory (512MB) with reservations (1 core, 256MB)
- **Health management**: Restart policy set to `unless-stopped`
- **Logging**: JSON file driver with rotation (10MB per file, max 3 files)

### Verification

After starting, you should see output like:
```
Starting the application...

=== Configuration ===
Workers Count:    50
Streams Count:    100
Stream Interval:  1000 ms
Chunk Size:       1024 bytes
Buffer Size:      500
====================
Application is running...
50 workers started
100 streams started
Worker 0 processed chunk 0 from stream 5, length of chunk's data is 256
...
```

Press `Ctrl+C` to gracefully shut down and view final statistics.

---

## Application Flow

### Startup Phase

```
1. main() initializes
   ├─ Create context (for cancellation)
   ├─ Create metrics tracker
   ├─ Create buffered channel (capacity: 500)
   └─ Initialize WaitGroup for synchronization

2. Start 50 Worker goroutines
   ├─ Each worker listens on dataChan
   ├─ Waits for chunks to process
   └─ Adds to WaitGroup

3. Start 100 Stream goroutines
   ├─ Each stream generates random chunks
   ├─ Sends to dataChan (non-blocking)
   └─ Adds to WaitGroup

4. Application is now running
   └─ All goroutines are active
```

### Processing Phase

```
Stream Goroutine (100 total)          Worker Goroutine (50 total)
├─ Generate random chunk (1KB)        ├─ Wait for chunk from channel
├─ Try non-blocking send to channel   ├─ Receive chunk
│  ├─ Success → Continue              ├─ Record start time
│  └─ Fail (buffer full) → Record     ├─ ProcessChunk() (20-50ms sleep)
│     dropped chunk                   ├─ Calculate latency
├─ Sleep 0-1000ms (random)            ├─ Record metrics
└─ Repeat until ctx.Done()            ├─ Log result to console
                                      └─ Repeat until ctx.Done()
```

### Metrics Collection

```
As each chunk is processed:
1. Latency = now - chunk.CreatedAt
2. Add latency to metrics buffer
3. Update:
   - Total processed count
   - Dropped chunks count
   - Throughput calculation
   - Latency percentiles (P50, P95, P99)
```

### Shutdown Phase

```
1. User presses Ctrl+C (SIGINT/SIGTERM signal)
   └─ Signal received in sigChan

2. cancel() called
   └─ ctx.Done() triggers in all goroutines

3. All goroutines exit their loops
   ├─ Streams stop generating
   ├─ Workers finish processing in-flight chunks
   └─ Each calls wg.Done()

4. wg.Wait() unblocks when all Done() called
   └─ All 150 goroutines have exited

5. Final statistics printed
   ├─ Calculate final metrics
   ├─ Print summary
   └─ Application exits (exit code 0)
```

### Data Flow Diagram

```
┌──────────────────────────┐
│  Stream Goroutine (100)  │
│  Generates: 1KB chunks   │
│  Frequency: 0-1000ms     │
└───────────┬──────────────┘
            │ Non-blocking send
            │ (drop if buffer full)
            ▼
    ┌────────────────────┐
    │ Buffered Channel   │
    │ Capacity: 500      │
    │ Type: DataChunk    │
    └────────┬───────────┘
             │ Blocking receive
             │ (worker waits)
             ▼
┌──────────────────────────┐
│ Worker Goroutine (50)    │
│ Processing: 20-50ms      │
│ Records latency metrics  │
└──────────────────────────┘
             │
             ▼
    ┌────────────────────┐
    │ Metrics Tracker    │
    │ (Thread-safe)      │
    │ - Counts           │
    │ - Latencies        │
    │ - Percentiles      │
    └────────────────────┘
```

---

## Expected Output

When you run the application, you'll see:

```
Starting the application...
Application is running...
50 workers started
100 streams started
Worker 0 processed chunk 0 from stream 5, length of chunk's data is 256
Worker 1 processed chunk 1 from stream 12, length of chunk's data is 789
Worker 2 processed chunk 2 from stream 3, length of chunk's data is 512
...
[Continues processing chunks across all streams and workers]
...
[Press Ctrl+C to shut down gracefully]

Shutdown signal received...

=== Final Statistics ===
Processed: 15482 chunks
Dropped: 0 chunks
Elapsed: 1m23s
Throughput: 185.92 chunks/sec
Avg Latency: 35.234ms
P50 Latency: 34.123ms
P95 Latency: 48.567ms
P99 Latency: 52.891ms
Application exited.
```

### Output Explanation

- **Processed**: Total number of chunks successfully processed by workers
- **Dropped**: Chunks discarded due to channel buffer overflow (should be 0 under normal load)
- **Elapsed**: Total runtime from start to graceful shutdown
- **Throughput**: Chunks processed per second (typically 150-250 depending on system)
- **Latency Metrics**: Time from chunk creation to completion
  - **Avg**: Mean latency across all chunks
  - **P50**: 50th percentile (median) - half of requests faster
  - **P95**: 95th percentile - 95% of requests faster than this
  - **P99**: 99th percentile - 99% of requests faster than this

---

## Manual Testing

### Test Steps

1. **Monitor performance:**
   The console will continuously show which worker is processing which chunk and the data length.

2. **Test with different configurations** (edit `internal/config/config.go`):
   ```go
   const (
       WORKERS_COUNT = 25      // Try 25, 75, 100
       STREAMS_COUNT = 50      // Try 50, 150, 200
       STREAM_INTERVAL = 500   // Try 500, 2000 for slower/faster streams
       CHUNK_SIZE = 2048       // Try 2KB chunks
       BUFFER_SIZE = 1000      // Try larger buffer for more headroom
   )
   ```
   
   Then rebuild and run to observe performance differences.

3. **Verify graceful shutdown:**
   The application will:
   - Receive SIGINT/SIGTERM signal
   - Cancel the context (signals all goroutines to stop)
   - Wait for all goroutines to finish (WaitGroup)
   - Print final statistics
   - Exit cleanly

### Expected Behaviors

- **Normal operation**: 0-10 dropped chunks (good backpressure handling)
- **High load (100 streams + large intervals)**: May see minor backpressure
- **Very high load (200 streams + small intervals)**: Will start dropping chunks
- **Graceful shutdown**: Always completes all in-flight processing before exiting

---

## Architecture

### System Design

```
┌─────────────────────────────────────────────────────────────┐
│                     Stream Generators (100)                  │
│  [Stream 1] [Stream 2] ... [Stream 99] [Stream 100]        │
└────────────────────────┬────────────────────────────────────┘
                         │
                         ▼
              ┌──────────────────────┐
              │  Buffered Channel    │
              │    (500 capacity)    │
              └──────────┬───────────┘
                         │
                         ▼
              ┌──────────────────────┐
              │   Worker Pool (50)   │
              │  [W1] [W2] ... [W50] │
              └──────────┬───────────┘
                         │
                         ▼
              ┌──────────────────────┐
              │  Metrics Tracking    │
              │  (Thread-safe)       │
              └──────────────────────┘
```

### Core Components

#### 1. Stream Generators (100 goroutines)
- **Purpose**: Simulate real-time data sources
- **Behavior**: Each stream continuously generates random binary chunks (1KB)
- **Production Rate**: Variable interval (0-1000ms) averages ~500ms per chunk
- **Thread Safety**: Each stream has its own random number generator
- **Backpressure**: Non-blocking send with drop tracking

**Implementation Highlights:**
```go
// Each stream has its own random source (thread-safe)
rng := mathRand.New(mathRand.NewSource(time.Now().UnixNano() + int64(id)))

// Non-blocking send with drop tracking
select {
case dataChan <- chunk:
    // Successfully queued
default:
    m.RecordDropped()  // Track dropped chunks
}
```

#### 2. Worker Pool (50 workers)
- **Purpose**: Process chunks concurrently
- **Processing Time**: Simulates external API calls (20-50ms per chunk)
- **Capacity**: ~1428 chunks/second (50 workers × 28.5 chunks/sec each)
- **Overhead**: 7x production rate for safety margin

**Why 50 Workers?**

| Metric | Value | Calculation |
|--------|-------|-------------|
| Production Rate | ~200 chunks/sec | 100 streams × 2 chunks/sec |
| Worker Capacity | 28.5 chunks/sec | 1000ms ÷ 35ms avg processing |
| Total Capacity | 1428 chunks/sec | 50 workers × 28.5 |
| **Safety Margin** | **7.1x** | 1428 ÷ 200 |

**Alternative configurations:**
- 20 workers: 3.5x overhead (still safe, less margin)
- 100 workers: 14x overhead (overkill, waste resources)
- **50 workers: Sweet spot** ✅

#### 3. Communication Channel
- **Type**: Buffered channel
- **Capacity**: 500 chunks
- **Buffer Time**: ~2.5 seconds at peak production (500 ÷ 200/sec)
- **Purpose**: Absorb production bursts, prevent blocking

**Buffer Sizing Rationale:**
```
Production spikes: 100 chunks/sec (all streams at min interval)
Buffer capacity:   500 chunks
Time to fill:      5 seconds (plenty of time for workers to catch up)
Normal usage:      ~14% (70/500 chunks typically queued)
```

#### 4. Metrics System
- **Thread Safety**: `sync.RWMutex` for concurrent access
- **Tracked Metrics**:
  - Total chunks processed
  - Dropped chunks (backpressure events)
  - Latency distribution (avg, P50, P95, P99)
  - Throughput (chunks/second)
  - Elapsed time

**Percentile Calculation:**
```go
// Sorts latencies and calculates percentiles
sort.Slice(sorted, func(i, j int) bool { return sorted[i] < sorted[j] })
p95 := sorted[int(float64(len(sorted))*0.95)]
```

---

## Concurrency Patterns

### 1. Worker Pool Pattern
Instead of 1 goroutine per stream (wasteful), we use a fixed pool:
- **Benefit**: Bounded resource usage
- **Tradeoff**: Slightly higher latency under peak load
- **Result**: Predictable performance

### 2. Fan-Out/Fan-In via Channels
- **Fan-Out**: 100 streams → 1 channel (many-to-one)
- **Fan-In**: 1 channel → 50 workers (one-to-many)
- **Decoupling**: Streams don't know about workers, vice versa

### 3. Context-Based Cancellation
```go
ctx, cancel := context.WithCancel(context.Background())

// All goroutines listen on ctx.Done()
select {
case <-ctx.Done():
    return  // Clean exit
case chunk := <-dataChan:
    // Process
}
```

### 4. WaitGroup Synchronization
```go
var wg sync.WaitGroup

// Start goroutines
wg.Add(1)
go func() {
    defer wg.Done()
    // Work
}()

// Wait for completion
wg.Wait()  // Blocks until all Done() calls
```

---

## Performance Characteristics

### Latency

| Metric | Target | Actual | Status |
|--------|--------|--------|--------|
| Average | <100ms | ~35ms | ✅ |
| P50 | <100ms | ~34ms | ✅ |
| P95 | <100ms | ~48ms | ✅ |
| P99 | <100ms | ~52ms | ✅ |

**Latency Breakdown:**
- Processing time: 20-50ms (simulated API call)
- Channel overhead: <5ms (minimal with buffering)
- Context switching: <2ms (negligible)

### Throughput

| Scenario | Rate | Headroom |
|----------|------|----------|
| Normal Load | ~200 chunks/sec | 7x |
| Peak Bursts | ~400 chunks/sec | 3.5x |
| Maximum Capacity | 1428 chunks/sec | N/A |

### Resource Usage

| Resource | Usage | Notes |
|----------|-------|-------|
| Goroutines | 150 | 100 streams + 50 workers |
| Memory (Channel) | ~500KB | 500 × 1KB chunks |
| CPU | Minimal | Mostly sleeping/waiting |
| Dropped Chunks | 0 | Under normal conditions |
