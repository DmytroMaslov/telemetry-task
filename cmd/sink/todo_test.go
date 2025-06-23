package main

/*
import (
	"fmt"
	"os"
	"sync"
	pb "telemetry-task/protocol"
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_MetricCollector(t *testing.T) {
	tmpDir := t.TempDir()
	var testWg sync.WaitGroup

	mc, err := NewMetricCollector(tmpDir+"/test_metrics.txt", 1024, 10)
	assert.NoError(t, err)
	count := 100

	metrics := generateTestMetrics(count)

	go func() {
		defer testWg.Done()

		testWg.Add(1)
		for _, metric := range metrics {
			mc.metricsCh <- metric
		}
	}()

	go mc.Start()
	testWg.Wait()
	mc.Stop()

	file, err := os.ReadFile(tmpDir + "/test_metrics.txt")
	assert.NoError(t, err)
	fmt.Printf("File content:\n%s\n", string(file))
	t.Fail()
}

func generateTestMetrics(count int) []*pb.Metric {
	metrics := make([]*pb.Metric, 0, count)
	for i := 0; i < count; i++ {
		metrics = append(metrics, &pb.Metric{
			Name:      "test_metric",
			Value:     int64(i),
			Timestamp: uint64(i * 1000), // Simulating timestamps
		})
	}
	return metrics
}
*/
