package instrumentation

import (
	"sync/atomic"

	b "github.com/miton18/go-warp10/base"
)

// Histogram is a metric that can only grow
type Histogram struct {
	name, help string
	buckets    map[float64]uint64
	context    b.Labels
}

// NewHistogram initialize a new counter
func NewHistogram(name string, context b.Labels, buckets []float64, help string) (mc *Counter) {
	if context == nil {
		context = b.Labels{}
	}
	bucketsMap := map[float64]uint64{}
	for _, bucket := range buckets {
		bucketsMap[bucket] = 0
	}
	return &Counter{
		count:   0,
		name:    name,
		context: context,
		help:    help,
	}
}

// Name return the metric name of the Histogram
func (mh *Histogram) Name() string {
	return mh.name
}

// Help return informations about this metric
func (mh *Histogram) Help() string {
	return "Histogram: " + mh.help
}

// Put a value into the histogram
func (mh *Histogram) Put(value float64) {

}

// Get return a plain GTS
func (mh *Histogram) Get() b.GTSList {
	return b.GTSList{&b.GTS{
		ClassName: mh.name,
		Labels:    mh.context,
		Values:    [][]interface{}{{atomic.LoadUint64(&mh.count)}},
	}}
}
