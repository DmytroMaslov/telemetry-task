package validator

import (
	"errors"
	pb "telemetry-task/protocol/telemetry"
)

type Validator interface {
	IsValid(*pb.Metric) error
}

type MetricValidator struct{}

func NewMetricValidator() *MetricValidator {
	return &MetricValidator{}
}

func (v *MetricValidator) IsValid(metric *pb.Metric) error {
	if metric == nil {
		return errors.New("metric cannot be nil")
	}
	if metric.Name == "" {
		return errors.New("metric name cannot be empty")
	}
	if metric.Value < 0 {
		return errors.New("metric value cannot be negative")
	}
	if metric.Timestamp == 0 {
		return errors.New("metric timestamp cannot be zero")
	}
	return nil
}
