package utils

import (
	"encoding/json"
	"fmt"
	"os"
	"sync"
	"time"
)

type Metric struct {
	Timestamp string `json:"timestamp"`
	LevelName string `json:"level_name"`
	Value     int    `json:"value"`
}

// Entry represents each stream. If the stream fails, an error will be present.
type Entry struct {
	Error  error
	Metric Metric
}

// Stream helps transmit each streams withing a channel.
type Stream struct {
	stream chan Entry
}

// NewJSONStream returns a new `Stream` type.
func NewJSONStream() Stream {
	return Stream{
		stream: make(chan Entry),
	}
}

// Watch watches JSON streams. Each stream entry will either have an error or a
// User object. Client code does not need to explicitly exit after catching an
// error as the `Start` method will close the channel automatically.
func (s Stream) Watch() <-chan Entry {
	return s.stream
}

// Start starts streaming JSON file line by line. If an error occurs, the channel
// will be closed.
func (s Stream) Start(path string) {
	// Stop streaming channel as soon as nothing left to read in the file.
	defer close(s.stream)

	// Open file to read.
	file, err := os.Open(path)
	if err != nil {
		s.stream <- Entry{Error: fmt.Errorf("open file: %w", err)}
		return
	}
	defer file.Close()

	decoder := json.NewDecoder(file)

	// Read opening delimiter. `[` or `{`
	if _, err := decoder.Token(); err != nil {
		s.stream <- Entry{Error: fmt.Errorf("decode opening delimiter: %w", err)}
		return
	}

	// Read file content as long as there is something.
	i := 1
	for decoder.More() {
		var metric Metric
		if err := decoder.Decode(&metric); err != nil {
			s.stream <- Entry{Error: fmt.Errorf("decode line %d: %w", i, err)}
			return
		}
		s.stream <- Entry{Metric: metric}

		i++
	}

	// Read closing delimiter. `]` or `}`
	if _, err := decoder.Token(); err != nil {
		s.stream <- Entry{Error: fmt.Errorf("decode closing delimiter: %w", err)}
		return
	}
}

func JsonFile(path string, startTime, endTime time.Time, cMap chan map[string]int) {
	fmt.Println(path)
	wg := sync.WaitGroup{}
	wg.Add(2)

	resMap := make(map[string]int)
	cMetrix := make(chan Metric, 10)
	stream := NewJSONStream()
	go func() {
		defer wg.Done()
		for data := range stream.Watch() {
			if data.Error != nil {
				fmt.Println(data.Error)
			}
			checkTime, err := time.Parse(time.RFC3339, data.Metric.Timestamp)
			if err == nil && InTimeSpan(startTime, endTime, checkTime) {
				cMetrix <- data.Metric
			}
		}
		close(cMetrix)
	}()
	go func() {
		defer wg.Done()
		for c := range cMetrix {
			resMap[c.LevelName] += c.Value
		}
	}()
	stream.Start(path)
	wg.Wait()
	cMap <- resMap
}
