package instrumentation

import (
	"sync/atomic"

	b "github.com/miton18/go-warp10/base"
)

// Counter is a metric that can only grow
type Counter struct {
	name, help string
	count      uint64
	context    b.Labels
}

// NewCounter initialize a new counter
func NewCounter(name string, context b.Labels, help string) (mc *Counter) {
	if context == nil {
		context = b.Labels{}
	}
	return &Counter{
		count:   0,
		name:    name,
		context: context,
		help:    help,
	}
}

// Name return the metric name of the counter
func (mc *Counter) Name() string {
	return mc.name
}

// Help return informations about this metric
func (mc *Counter) Help() string {
	return "Counter: " + mc.help
}

// Get return a plain GTS
func (mc *Counter) Get() b.GTSList {
	return b.GTSList{&b.GTS{
		ClassName: mc.name,
		Labels:    mc.context,
		Values:    [][]interface{}{{atomic.LoadUint64(&mc.count)}},
	}}
}

// Reset set to 0 the counter
func (mc *Counter) Reset() {
	atomic.AddUint64(&mc.count, ^uint64(mc.count-1))
}

// Inc add 1 to the counter
func (mc *Counter) Inc() {
	atomic.AddUint64(&mc.count, 1)
}

// Add N to the counter
func (mc *Counter) Add(n uint64) {
	atomic.AddUint64(&mc.count, n)
}
