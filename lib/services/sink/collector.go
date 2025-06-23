package sink

import (
	"bytes"
	"fmt"
	"os"
	"sync"
	pb "telemetry-task/protocol"
	"time"
)

// TODO: refactor metric collector

type MetricCollector struct {
	metricsCh     chan *pb.Metric
	lock          *sync.Mutex
	buffer        bytes.Buffer
	doneCh        chan struct{}
	flushInterval int
	file          *os.File
}

func NewMetricCollector(cfg *Config) (*MetricCollector, error) {
	// TODO: add cfg validation
	file, err := os.OpenFile(cfg.FilePath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return nil, fmt.Errorf("opening file: %v", err)
	}

	return &MetricCollector{
		lock:          &sync.Mutex{},
		buffer:        *bytes.NewBuffer(make([]byte, 0, cfg.BufferSize)),
		flushInterval: cfg.FlushInterval,
		doneCh:        make(chan struct{}),
		metricsCh:     make(chan *pb.Metric),
		file:          file,
	}, nil
}

func (bw *MetricCollector) Start() error {
	logger.Info("Starting metric collector...")
	ticker := time.NewTicker(time.Duration(bw.flushInterval) * time.Millisecond)

	go func() {
		for {
			select {
			case metric := <-bw.metricsCh:
				err := bw.Write(metric)
				if err != nil {
					fmt.Println("Error writing metric:", err)
				}

			case <-ticker.C:
				fmt.Println("Ticker triggered, flushing buffer...")
				if err := bw.Flush(); err != nil {
					fmt.Printf("Error writing to file in ticker case: %v\n", err)
				}

			case <-bw.doneCh:
				fmt.Println("Received done signal, flushing buffer and stopping...")
				if err := bw.Flush(); err != nil {
					fmt.Printf("Error writing to file in done case: %v\n", err)
				}
				bw.file.Close()
				return
			}
		}
	}()
	return nil
}

func (bw *MetricCollector) Write(metric *pb.Metric) error {
	line := fmt.Sprintf("%d,%s,%d\n", metric.Value, metric.Name, metric.Timestamp)
	fmt.Println("Writing metric:", line)

	if len([]byte(line)) > bw.buffer.Available() {
		fmt.Println("Buffer is full, flushing...")
		bw.Flush() // Flush if the line exceeds available buffer space
	}
	bw.lock.Lock()
	defer bw.lock.Unlock()
	_, err := bw.buffer.Write([]byte(line))
	if err != nil {
		return fmt.Errorf("error writing to buffer: %v", err)
	}
	return nil
}

func (bw *MetricCollector) Flush() error {
	fmt.Println("Flushing buffer to file...")
	bw.lock.Lock()
	defer bw.lock.Unlock()

	if bw.buffer.Len() == 0 {
		return nil // Nothing to write
	}

	_, err := bw.file.WriteString(bw.buffer.String()) // Write the buffer to the file
	if err != nil {
		return fmt.Errorf("error writing to file: %v", err)
	}
	bw.buffer.Reset() // Clear the buffer after writing
	return nil
}

func (bw *MetricCollector) Stop() error {
	fmt.Println("Stopping MetricCollector...")
	bw.doneCh <- struct{}{}
	return nil
}
