package warp10

import (
	"sync/atomic"
)

// MetricCounter is a metric that can only grow
type MetricCounter struct {
	name, help string
	count      uint64
	context    Labels
}

// NewMetricCounter initialize a new counter
func NewMetricCounter(name string, context Labels, help string) (mc *MetricCounter) {
	if context == nil {
		context = Labels{}
	}
	return &MetricCounter{
		count:   0,
		name:    name,
		context: context,
		help:    help,
	}
}

// Name return the metric name of the counter
func (mc *MetricCounter) Name() string {
	return mc.name
}

// Help return informations about this metric
func (mc *MetricCounter) Help() string {
	return "counter: " + mc.help
}

// Get return a plain GTS
func (mc *MetricCounter) Get() GTSList {
	return GTSList{&GTS{
		ClassName: mc.name,
		Labels:    mc.context,
		Values:    [][]interface{}{{atomic.LoadUint64(&mc.count)}},
	}}
}

// Reset set to 0 the counter
func (mc *MetricCounter) Reset() {
	atomic.AddUint64(&mc.count, ^uint64(mc.count))
}

// Inc add 1 to the counter
func (mc *MetricCounter) Inc() {
	atomic.AddUint64(&mc.count, 1)
}

// Add N to the counter
func (mc *MetricCounter) Add(n uint64) {
	atomic.AddUint64(&mc.count, n)
}
