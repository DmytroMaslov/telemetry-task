package sink

import (
	"fmt"
	"io"
	logUtil "telemetry-task/lib/logger"
	pb "telemetry-task/protocol/telemetry"

	"telemetry-task/lib/validator"

	rateLimiter "telemetry-task/lib/ratelimiter"

	"google.golang.org/protobuf/proto"
)

var (
	logger = logUtil.LoggerWithPrefix("SINK")
)

type Sink struct {
	pb.UnimplementedTelemetryServiceServer
	mc              *MetricCollector
	metricEncryptor *EncryptManager
	rateLimited     rateLimiter.RateLimited
	validator       validator.Validator
}

func NewSink(cfg *Config) (*Sink, error) {

	if cfg.RateLimit <= 0 {
		cfg.RateLimit = 1024 * 1024 // Default to 1 MB/s if not set
	}
	em := NewEncryptManager()
	mc, err := NewMetricCollector(cfg, em.outputCh)
	if err != nil {
		return nil, fmt.Errorf("crate new metric collector, err:%w", err)
	}

	return &Sink{
		mc:              mc,
		rateLimited:     rateLimiter.NewRateLimiter(cfg.RateLimit),
		validator:       validator.NewMetricValidator(),
		metricEncryptor: em,
	}, nil
}

func (s *Sink) SendMetrics(stream pb.TelemetryService_SendMetricsServer) error {
	for {
		res, err := stream.Recv()
		if err == io.EOF {
			// Stream has ended, return a result
			return stream.SendAndClose(&pb.Result{
				Success: true,
				Message: "Metrics received successfully",
			})

		}
		if err != nil {
			return fmt.Errorf("error receiving metrics: %v", err)
		}
		// rate check
		resB, err := proto.Marshal(res)
		if err != nil {
			logger.Error("failed to marshall metric", "err", err)
		}
		if !s.rateLimited.IsAllowed(resB) {
			logger.Debug("Dropping metric due to rate limit")
			continue
		}
		// validate
		if err := s.validator.IsValid(res); err != nil {
			logger.Debug("Dropping invalid metric", "metric", res)
			continue
		}

		s.metricEncryptor.inputCh <- fmt.Sprintf("%d,%s,%d", res.Value, res.Name, res.Timestamp)
	}
}

func (s *Sink) Start() error {
	if err := s.metricEncryptor.Run(); err != nil {
		return fmt.Errorf("run encryptor manager: %w", err)
	}
	if err := s.mc.Start(); err != nil {
		return fmt.Errorf("run metric collector: %w", err)
	}
	return nil
}

func (s *Sink) Stop() error {
	s.metricEncryptor.Stop()
	return s.mc.Stop()
}
