package core

import (
	"sync"
	"time"
)

// MetricsCollector collects and aggregates metrics
type MetricsCollector struct {
	metrics map[string]*Metric
	mu      sync.RWMutex
}

// Metric represents a single metric
type Metric struct {
	Name      string
	Count     int64
	Sum       float64
	Min       float64
	Max       float64
	Avg       float64
	LastValue float64
	LastTime  time.Time
	Values    []float64 // For percentile calculation
}

// NewMetricsCollector creates a new metrics collector
func NewMetricsCollector() *MetricsCollector {
	return &MetricsCollector{
		metrics: make(map[string]*Metric),
	}
}

// Record records a metric value
func (mc *MetricsCollector) Record(name string, value float64) {
	mc.mu.Lock()
	defer mc.mu.Unlock()

	metric, exists := mc.metrics[name]
	if !exists {
		metric = &Metric{
			Name:   name,
			Min:    value,
			Max:    value,
			Values: make([]float64, 0),
		}
		mc.metrics[name] = metric
	}

	metric.Count++
	metric.Sum += value
	metric.LastValue = value
	metric.LastTime = time.Now()

	if value < metric.Min {
		metric.Min = value
	}
	if value > metric.Max {
		metric.Max = value
	}

	metric.Avg = metric.Sum / float64(metric.Count)
	metric.Values = append(metric.Values, value)

	// Keep only last 1000 values for percentile calculation
	if len(metric.Values) > 1000 {
		metric.Values = metric.Values[len(metric.Values)-1000:]
	}
}

// GetMetric gets a metric by name
func (mc *MetricsCollector) GetMetric(name string) (*Metric, bool) {
	mc.mu.RLock()
	defer mc.mu.RUnlock()
	metric, exists := mc.metrics[name]
	return metric, exists
}

// GetAllMetrics gets all metrics
func (mc *MetricsCollector) GetAllMetrics() map[string]*Metric {
	mc.mu.RLock()
	defer mc.mu.RUnlock()

	result := make(map[string]*Metric)
	for name, metric := range mc.metrics {
		result[name] = metric
	}
	return result
}

// Reset resets all metrics
func (mc *MetricsCollector) Reset() {
	mc.mu.Lock()
	defer mc.mu.Unlock()
	mc.metrics = make(map[string]*Metric)
}

// Percentile calculates percentile for a metric
func (m *Metric) Percentile(p float64) float64 {
	if len(m.Values) == 0 {
		return 0
	}

	// Simple percentile calculation
	// Would use proper algorithm in production
	index := int(float64(len(m.Values)) * p / 100.0)
	if index >= len(m.Values) {
		index = len(m.Values) - 1
	}
	return m.Values[index]
}

// Counter represents a counter metric
type Counter struct {
	name  string
	value int64
	mu    sync.RWMutex
}

// NewCounter creates a new counter
func NewCounter(name string) *Counter {
	return &Counter{name: name}
}

// Inc increments the counter
func (c *Counter) Inc() {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.value++
}

// Add adds a value to the counter
func (c *Counter) Add(delta int64) {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.value += delta
}

// Value returns the counter value
func (c *Counter) Value() int64 {
	c.mu.RLock()
	defer c.mu.RUnlock()
	return c.value
}

// Reset resets the counter
func (c *Counter) Reset() {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.value = 0
}

// Gauge represents a gauge metric
type Gauge struct {
	name  string
	value float64
	mu    sync.RWMutex
}

// NewGauge creates a new gauge
func NewGauge(name string) *Gauge {
	return &Gauge{name: name}
}

// Set sets the gauge value
func (g *Gauge) Set(value float64) {
	g.mu.Lock()
	defer g.mu.Unlock()
	g.value = value
}

// Value returns the gauge value
func (g *Gauge) Value() float64 {
	g.mu.RLock()
	defer g.mu.RUnlock()
	return g.value
}

// Histogram represents a histogram metric
type Histogram struct {
	name   string
	buckets []float64
	counts []int64
	mu     sync.RWMutex
}

// NewHistogram creates a new histogram
func NewHistogram(name string, buckets []float64) *Histogram {
	return &Histogram{
		name:    name,
		buckets: buckets,
		counts:  make([]int64, len(buckets)+1), // +1 for overflow bucket
	}
}

// Observe records an observation
func (h *Histogram) Observe(value float64) {
	h.mu.Lock()
	defer h.mu.Unlock()

	bucketIndex := len(h.buckets)
	for i, bucket := range h.buckets {
		if value <= bucket {
			bucketIndex = i
			break
		}
	}

	h.counts[bucketIndex]++
}

// GetCounts returns bucket counts
func (h *Histogram) GetCounts() []int64 {
	h.mu.RLock()
	defer h.mu.RUnlock()
	result := make([]int64, len(h.counts))
	copy(result, h.counts)
	return result
}

// Timer represents a timer metric
type Timer struct {
	name      string
	durations []time.Duration
	mu        sync.RWMutex
}

// NewTimer creates a new timer
func NewTimer(name string) *Timer {
	return &Timer{
		name:      name,
		durations: make([]time.Duration, 0),
	}
}

// Record records a duration
func (t *Timer) Record(duration time.Duration) {
	t.mu.Lock()
	defer t.mu.Unlock()
	t.durations = append(t.durations, duration)
	if len(t.durations) > 1000 {
		t.durations = t.durations[len(t.durations)-1000:]
	}
}

// Time records time for a function
func (t *Timer) Time(fn func()) {
	start := time.Now()
	fn()
	t.Record(time.Since(start))
}

// Average returns average duration
func (t *Timer) Average() time.Duration {
	t.mu.RLock()
	defer t.mu.RUnlock()
	if len(t.durations) == 0 {
		return 0
	}
	var sum time.Duration
	for _, d := range t.durations {
		sum += d
	}
	return sum / time.Duration(len(t.durations))
}

// Min returns minimum duration
func (t *Timer) Min() time.Duration {
	t.mu.RLock()
	defer t.mu.RUnlock()
	if len(t.durations) == 0 {
		return 0
	}
	min := t.durations[0]
	for _, d := range t.durations[1:] {
		if d < min {
			min = d
		}
	}
	return min
}

// Max returns maximum duration
func (t *Timer) Max() time.Duration {
	t.mu.RLock()
	defer t.mu.RUnlock()
	if len(t.durations) == 0 {
		return 0
	}
	max := t.durations[0]
	for _, d := range t.durations[1:] {
		if d > max {
			max = d
		}
	}
	return max
}

// RepositoryMetrics tracks repository operation metrics
type RepositoryMetrics struct {
	operationCounters map[string]*Counter
	operationTimers   map[string]*Timer
	errorCounters     map[string]*Counter
	mu                sync.RWMutex
}

// NewRepositoryMetrics creates new repository metrics
func NewRepositoryMetrics() *RepositoryMetrics {
	return &RepositoryMetrics{
		operationCounters: make(map[string]*Counter),
		operationTimers:   make(map[string]*Timer),
		errorCounters:     make(map[string]*Counter),
	}
}

// RecordOperation records an operation
func (rm *RepositoryMetrics) RecordOperation(operation string, duration time.Duration, err error) {
	rm.mu.Lock()
	defer rm.mu.Unlock()

	// Counter
	counter, exists := rm.operationCounters[operation]
	if !exists {
		counter = NewCounter(operation)
		rm.operationCounters[operation] = counter
	}
	counter.Inc()

	// Timer
	timer, exists := rm.operationTimers[operation]
	if !exists {
		timer = NewTimer(operation)
		rm.operationTimers[operation] = timer
	}
	timer.Record(duration)

	// Error counter
	if err != nil {
		errorCounter, exists := rm.errorCounters[operation]
		if !exists {
			errorCounter = NewCounter(operation + "_errors")
			rm.errorCounters[operation] = errorCounter
		}
		errorCounter.Inc()
	}
}

// GetOperationStats gets statistics for an operation
func (rm *RepositoryMetrics) GetOperationStats(operation string) map[string]interface{} {
	rm.mu.RLock()
	defer rm.mu.RUnlock()

	stats := make(map[string]interface{})

	if counter, exists := rm.operationCounters[operation]; exists {
		stats["count"] = counter.Value()
	}

	if timer, exists := rm.operationTimers[operation]; exists {
		stats["avg_duration"] = timer.Average()
		stats["min_duration"] = timer.Min()
		stats["max_duration"] = timer.Max()
	}

	if errorCounter, exists := rm.errorCounters[operation]; exists {
		stats["error_count"] = errorCounter.Value()
	}

	return stats
}

// GetAllStats gets all operation statistics
func (rm *RepositoryMetrics) GetAllStats() map[string]map[string]interface{} {
	rm.mu.RLock()
	defer rm.mu.RUnlock()

	allStats := make(map[string]map[string]interface{})
	for operation := range rm.operationCounters {
		allStats[operation] = rm.GetOperationStats(operation)
	}
	return allStats
}

