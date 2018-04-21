package instrumentation

import (
	"sync/atomic"

	b "github.com/miton18/go-warp10/base"
)

// Gauge is a metric that can only grow
type Gauge struct {
	name, help string
	count      atomic.Value
	context    b.Labels
}

// NewGauge initialize a new counter
func NewGauge(name string, context b.Labels, help string) (mc *Gauge) {
	if context == nil {
		context = b.Labels{}
	}
	g := &Gauge{
		count:   atomic.Value{},
		name:    name,
		context: context,
		help:    help,
	}
	g.count.Store(float64(0))
	return g
}

// Name return the metric name of the counter
func (mc *Gauge) Name() string {
	return mc.name
}

// Help return informations about this metric
func (mc *Gauge) Help() string {
	return "Gauge: " + mc.help
}

// Get return a plain GTS
func (mc *Gauge) Get() b.GTSList {
	return b.GTSList{&b.GTS{
		ClassName: mc.name,
		Labels:    mc.context,
		Values:    [][]interface{}{{mc.count.Load().(float64)}},
	}}
}

// Reset set to 0 the counter
func (mc *Gauge) Reset() {
	mc.count.Store(float64(0))
}

// Inc add 1 to the counter
func (mc *Gauge) Inc() {
	mc.count.Store(mc.count.Load().(float64) + float64(1))
}

// Dec N to the counter
func (mc *Gauge) Dec() {
	mc.count.Store(mc.count.Load().(float64) - float64(1))
}

// Add N to the counter
func (mc *Gauge) Add(n float64) {
	mc.count.Store(mc.count.Load().(float64) + n)
}

// Sub N to the counter
func (mc *Gauge) Sub(n float64) {
	mc.count.Store(mc.count.Load().(float64) - n)
}

// Set N to the counter
func (mc *Gauge) Set(n float64) {
	mc.count.Store(n)
}
