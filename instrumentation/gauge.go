package instrumentation

import (
	"sync/atomic"

	b "github.com/miton18/go-warp10/base"
)

// Gauge is a metric that can only grow
type Gauge struct {
	name, help string
	count      uint64
	context    b.Labels
}

// NewGauge initialize a new counter
func NewGauge(name string, context b.Labels, help string) (mc *Gauge) {
	if context == nil {
		context = b.Labels{}
	}
	return &Gauge{
		count:   0,
		name:    name,
		context: context,
		help:    help,
	}
}

// Name return the metric name of the counter
func (mc *Gauge) Name() string {
	return mc.name
}

// Help return informations about this metric
func (mc *Gauge) Help() string {
	return "counter: " + mc.help
}

// Get return a plain GTS
func (mc *Gauge) Get() b.GTSList {
	return b.GTSList{&b.GTS{
		ClassName: mc.name,
		Labels:    mc.context,
		Values:    [][]any{{atomic.LoadUint64(&mc.count)}},
	}}
}

// Reset set to 0 the counter
func (mc *Gauge) Reset() {
	atomic.AddUint64(&mc.count, ^uint64(mc.count-1))
}

// Inc add 1 to the counter
func (mc *Gauge) Inc() {
	atomic.AddUint64(&mc.count, 1)
}

// Dec N to the counter
func (mc *Gauge) Dec() {
	atomic.AddUint64(&mc.count, ^uint64(0))
}

// Add N to the counter
func (mc *Gauge) Add(n uint64) {
	atomic.AddUint64(&mc.count, n)
}

// Sub N to the counter
func (mc *Gauge) Sub(n uint64) {
	atomic.AddUint64(&mc.count, ^uint64(n-1))
}

// Set N to the counter
func (mc *Gauge) Set(n uint64) {
	atomic.StoreUint64(&mc.count, n)
}
