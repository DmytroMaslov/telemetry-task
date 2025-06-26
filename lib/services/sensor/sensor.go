package sensor

import (
	"context"
	"fmt"
	"sync"
	logUtil "telemetry-task/lib/logger"
	"telemetry-task/lib/ratecalculator"
	"time"

	"google.golang.org/grpc"
)

const (
	maxWorkers = 5
)

var (
	logger = logUtil.LoggerWithPrefix("Sensor")
)

type Sensor struct {
	MetricSender   *MetricSender
	RateCalculator ratecalculator.RateCalculator
	stopCh         chan struct{}
	wg             *sync.WaitGroup
}

func NewSensor(conn *grpc.ClientConn, rate int, name string) (*Sensor, error) {
	logger.Info("New sensor created", "name", name)
	ms := NewMetricSender(conn, name)

	rc, err := ratecalculator.NewRateCalculator(rate, 1*time.Second)
	if err != nil {
		return nil, err
	}

	return &Sensor{
		MetricSender:   ms,
		RateCalculator: rc,
		stopCh:         make(chan struct{}),
		wg:             &sync.WaitGroup{},
	}, nil
}

func (s *Sensor) Start(ctx context.Context) error {
	if err := s.MetricSender.EstablishConnection(ctx); err != nil {
		return fmt.Errorf("establish connection, err:%w", err)
	}
	trigger := make(chan struct{})
	errorCh := make(chan error)

	startTime := time.Now()
	workerCounter := 1

	// run first worker
	s.wg.Add(1)
	go s.Send(s.wg, trigger, errorCh)

	go func() {
		defer func() {
			close(trigger)
			close(errorCh)
		}()

		messageCounter := uint64(0)

		for {
			timeFromStart := time.Since(startTime)
			waitTime := s.RateCalculator.WaitToNextMessage(timeFromStart, messageCounter)
			time.Sleep(waitTime)
			if workerCounter < maxWorkers {
				select {
				case trigger <- struct{}{}:
					messageCounter++
					continue
				case <-s.stopCh:
					return
				default:
					workerCounter++
					s.wg.Add(1)
					logger.Debug("start new worker", "total workers", workerCounter)
					go s.Send(s.wg, trigger, errorCh)
				}
			}
			select {
			case trigger <- struct{}{}:
				messageCounter++
			case <-s.stopCh:
				return
			}
		}

	}()
	go func() {
		for err := range errorCh {
			logger.Error("failed to send metric", "err", err)
		}
	}()
	return nil

}

func (s *Sensor) Send(wg *sync.WaitGroup, trigger <-chan struct{}, errorCh chan<- error) {
	defer wg.Done()
	for range trigger {
		err := s.MetricSender.Send()
		if err != nil {
			errorCh <- err
		}
	}
}

func (s *Sensor) Stop() error {
	logger.Info("stopped sensor")
	s.stopCh <- struct{}{}
	close(s.stopCh)
	s.wg.Wait()

	if err := s.MetricSender.metricStream.CloseSend(); err != nil {
		return fmt.Errorf("failed to close metric stream: %w", err)
	}
	return nil
}
