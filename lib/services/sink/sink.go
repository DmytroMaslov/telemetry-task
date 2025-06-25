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
	mc          *MetricCollector
	rateLimited rateLimiter.RateLimited
	validator   validator.Validator
}

func NewSink(cfg *Config) (*Sink, error) {

	if cfg.RateLimit <= 0 {
		cfg.RateLimit = 1024 * 1024 // Default to 1 MB/s if not set
	}
	mc, err := NewMetricCollector(cfg)
	if err != nil {
		return nil, fmt.Errorf("crate new metric collector, err:%w", err)
	}

	return &Sink{
		mc:          mc,
		rateLimited: rateLimiter.NewRateLimiter(cfg.RateLimit),
		validator:   validator.NewMetricValidator(),
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
		s.mc.metricsCh <- res
	}
}

func (s *Sink) Start() error {
	return s.mc.Start()
}

func (s *Sink) Stop() error {
	return s.mc.Stop()
}
