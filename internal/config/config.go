package config

import (
	"fmt"
	"os"
	"strconv"
)

// Configuration holds all application settings loaded from environment variables
type Configuration struct {
	WorkersCount   int
	StreamsCount   int
	StreamInterval int
	ChunkSize      int
	BufferSize     int
}

// Config is the global configuration instance
var Config *Configuration

func LoadConfig() *Configuration {
	return &Configuration{
		WorkersCount:   getEnvInt("WORKERS_COUNT", 50),
		StreamsCount:   getEnvInt("STREAMS_COUNT", 100),
		StreamInterval: getEnvInt("STREAM_INTERVAL", 1000), // milliseconds
		ChunkSize:      getEnvInt("CHUNK_SIZE", 1024),      // bytes
		BufferSize:     getEnvInt("BUFFER_SIZE", 500),
	}
}

func getEnvInt(key string, defaultValue int) int {
	valueStr := os.Getenv(key)
	if valueStr == "" {
		return defaultValue
	}

	value, err := strconv.Atoi(valueStr)
	if err != nil {
		fmt.Printf("Warning: Invalid value for %s: %v, using default: %d\n", key, err, defaultValue)
		return defaultValue
	}

	return value
}

// GetConfig returns the current configuration, loading it if necessary
func GetConfig() *Configuration {
	if Config == nil {
		Config = LoadConfig()
	}
	return Config
}

func (c *Configuration) PrintConfig() {
	fmt.Println("\n=== Configuration ===")
	fmt.Printf("Workers Count:    %d\n", c.WorkersCount)
	fmt.Printf("Streams Count:    %d\n", c.StreamsCount)
	fmt.Printf("Stream Interval:  %d ms\n", c.StreamInterval)
	fmt.Printf("Chunk Size:       %d bytes\n", c.ChunkSize)
	fmt.Printf("Buffer Size:      %d\n", c.BufferSize)
	fmt.Println("====================")
}
