package sensor

import (
	"context"
	"fmt"
	"math/rand"
	pb "telemetry-task/protocol"
	"time"

	"google.golang.org/grpc"
)

type MetricSender struct {
	name         string
	client       pb.TelemetryServiceClient
	metricStream pb.TelemetryService_SendMetricsClient
}

func NewMetricSender(conn *grpc.ClientConn, name string) *MetricSender {
	return &MetricSender{
		name:   name,
		client: pb.NewTelemetryServiceClient(conn),
	}
}

func (ms *MetricSender) EstablishConnection(ctx context.Context) error {
	stream, err := ms.client.SendMetrics(ctx)
	if err != nil {
		return fmt.Errorf("failed to establish stream: %w", err)
	}
	ms.metricStream = stream
	logger.Debug("connection established")
	return nil
}

func (ms *MetricSender) Send() error {
	metric := &pb.Metric{
		Name:      ms.name,
		Value:     int64(rand.Int()),
		Timestamp: uint64(time.Now().Unix()),
	}

	if err := ms.metricStream.Send(metric); err != nil {
		return fmt.Errorf("send metric: %w", err)
	}
	logger.Debug("send metric", "metric", metric)
	return nil
}
