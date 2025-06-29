package sink

import (
	"bytes"
	"fmt"
	"os"
	"sync"
	"time"
)

type MetricCollector struct {
	metricsCh     chan string
	lock          *sync.Mutex
	buffer        bytes.Buffer
	doneCh        chan struct{}
	flushInterval int
	file          *os.File
}

func NewMetricCollector(cfg *Config, in chan string) (*MetricCollector, error) {
	if err := validateConfig(cfg); err != nil {
		return nil, fmt.Errorf("invalid configuration for metric collector: %w", err)
	}
	file, err := os.OpenFile(cfg.FilePath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return nil, fmt.Errorf("opening file: %v", err)
	}

	return &MetricCollector{
		lock:          &sync.Mutex{},
		buffer:        *bytes.NewBuffer(make([]byte, 0, cfg.BufferSize)),
		flushInterval: cfg.FlushInterval,
		doneCh:        make(chan struct{}),
		metricsCh:     in,
		file:          file,
	}, nil
}

func validateConfig(cfg *Config) error {
	if cfg.FilePath == "" {
		return fmt.Errorf("file path cannot be empty")
	}
	if cfg.BufferSize <= 0 {
		return fmt.Errorf("buffer size must be greater than 0")
	}
	if cfg.FlushInterval <= 0 {
		return fmt.Errorf("flush interval must be greater than 0")
	}
	return nil
}

func (bw *MetricCollector) Start() error {
	logger.Info("Starting metric collector...")
	ticker := time.NewTicker(time.Duration(bw.flushInterval) * time.Millisecond)

	errCh := make(chan error)
	go func() {
		defer func() {
			close(errCh)
		}()

		for {
			select {
			case metric, ok := <-bw.metricsCh:
				if ok {
					err := bw.Write(fmt.Sprintf("%s\n", metric))
					if err != nil {
						errCh <- err
					}
				}

			case <-ticker.C:
				if err := bw.FlushWithLock(); err != nil {
					errCh <- err
				}

			case <-bw.doneCh:
				logger.Debug("Received done signal, flushing buffer and stopping...")
				ticker.Stop()
				if err := bw.FlushWithLock(); err != nil {
					errCh <- err
				}
				bw.file.Close()
				return
			}
		}
	}()

	go func() {
		for err := range errCh {
			logger.Error("failed to receive metric", "err", err)
		}
	}()
	return nil
}

func (bw *MetricCollector) Write(line string) error {
	logger.Debug("write metric", "metric", line)

	bw.lock.Lock()
	defer bw.lock.Unlock()

	if len([]byte(line)) > bw.buffer.Available() {
		logger.Debug("Buffer is full, flushing...")
		bw.Flush() // Flush if the line exceeds available buffer space
	}
	_, err := bw.buffer.Write([]byte(line))
	if err != nil {
		return fmt.Errorf("writing to buffer: %v", err)
	}
	return nil
}

func (bw *MetricCollector) FlushWithLock() error {
	bw.lock.Lock()
	defer bw.lock.Unlock()

	return bw.Flush()
}

func (bw *MetricCollector) Flush() error {

	if bw.buffer.Len() == 0 {
		return nil // Nothing to write
	}

	_, err := bw.file.WriteString(bw.buffer.String()) // Write the buffer to the file
	if err != nil {
		return fmt.Errorf("writing string to file: %v", err)
	}
	bw.buffer.Reset() // Clear the buffer after writing
	return nil
}

func (bw *MetricCollector) Stop() error {
	fmt.Println("Stopping MetricCollector...")
	bw.doneCh <- struct{}{}
	close(bw.doneCh)
	return nil
}
