package core

import (
	"fmt"
	"testing"
	"time"
)

func TestMetricsCollector(t *testing.T) {
	collector := NewMetricsCollector()

	// Record metrics
	collector.Record("query_duration", 100.5)
	collector.Record("query_duration", 150.3)
	collector.Record("query_duration", 80.2)

	// Get metric
	metric, exists := collector.GetMetric("query_duration")
	if !exists {
		t.Error("Metric should exist")
		return
	}

	if metric.Count != 3 {
		t.Errorf("Expected count 3, got %d", metric.Count)
	}
	if metric.Avg == 0 {
		t.Error("Average should be calculated")
	}
}

func TestCounter(t *testing.T) {
	counter := NewCounter("test_counter")

	if counter.Value() != 0 {
		t.Error("Counter should start at 0")
	}

	counter.Inc()
	if counter.Value() != 1 {
		t.Errorf("Expected 1, got %d", counter.Value())
	}

	counter.Add(5)
	if counter.Value() != 6 {
		t.Errorf("Expected 6, got %d", counter.Value())
	}

	counter.Reset()
	if counter.Value() != 0 {
		t.Error("Counter should be reset to 0")
	}
}

func TestGauge(t *testing.T) {
	gauge := NewGauge("test_gauge")

	if gauge.Value() != 0 {
		t.Error("Gauge should start at 0")
	}

	gauge.Set(42.5)
	if gauge.Value() != 42.5 {
		t.Errorf("Expected 42.5, got %f", gauge.Value())
	}
}

func TestHistogram(t *testing.T) {
	buckets := []float64{10, 50, 100, 500}
	histogram := NewHistogram("test_histogram", buckets)

	histogram.Observe(5)
	histogram.Observe(25)
	histogram.Observe(75)
	histogram.Observe(200)
	histogram.Observe(1000)

	counts := histogram.GetCounts()
	if len(counts) != len(buckets)+1 {
		t.Errorf("Expected %d buckets, got %d", len(buckets)+1, len(counts))
	}
}

func TestTimer(t *testing.T) {
	timer := NewTimer("test_timer")

	timer.Record(100 * time.Millisecond)
	timer.Record(200 * time.Millisecond)
	timer.Record(150 * time.Millisecond)

	avg := timer.Average()
	if avg == 0 {
		t.Error("Average should be calculated")
	}

	min := timer.Min()
	if min != 100*time.Millisecond {
		t.Errorf("Expected 100ms, got %v", min)
	}

	max := timer.Max()
	if max != 200*time.Millisecond {
		t.Errorf("Expected 200ms, got %v", max)
	}
}

func TestRepositoryMetrics(t *testing.T) {
	metrics := NewRepositoryMetrics()

	metrics.RecordOperation("FindByID", 50*time.Millisecond, nil)
	metrics.RecordOperation("FindByID", 60*time.Millisecond, nil)
	metrics.RecordOperation("FindByID", 70*time.Millisecond, fmt.Errorf("error"))

	stats := metrics.GetOperationStats("FindByID")
	if stats["count"] != int64(3) {
		t.Errorf("Expected count 3, got %v", stats["count"])
	}
	if stats["error_count"] != int64(1) {
		t.Errorf("Expected error count 1, got %v", stats["error_count"])
	}
}

