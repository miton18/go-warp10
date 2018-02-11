package warp10

import (
	"sync"
)

// MetricCounter is a metric that can only grow
type MetricCounter struct {
	name, help string
	count      uint64
	context    Labels
	lock       *sync.Mutex
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
		lock:    &sync.Mutex{},
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
		Values:    [][]interface{}{{mc.count}},
	}}
}

// Reset set to 0 the counter
func (mc *MetricCounter) Reset() {
	mc.lock.Lock()
	mc.count = 0
	mc.lock.Unlock()
}

// Inc add 1 to the counter
func (mc *MetricCounter) Inc() {
	mc.lock.Lock()
	mc.count++
	mc.lock.Unlock()
}

// Add N to the counter
func (mc *MetricCounter) Add(n uint64) {
	mc.lock.Lock()
	mc.count += n
	mc.lock.Unlock()
}
