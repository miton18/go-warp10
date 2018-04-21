package instrumentation

import b "github.com/miton18/go-warp10/base"

// Metric is a thing you can monitor
type Metric interface {
	Name() string
	Help() string
	Get() b.GTSList
	Reset()
}

// Metrics is a collection of Metric
type Metrics []*Metric

// Len Allow to sort metrics by name
func (m Metrics) Len() int {
	return len(m)
}

// Swap Allow to sort metrics by name
func (m Metrics) Swap(i, j int) {
	m[i], m[j] = m[j], m[i]
}

// Less Allow to sort metrics by name
func (m Metrics) Less(i, j int) bool {
	if m[i] == nil {
		return false
	}
	if m[j] == nil {
		return true
	}
	return (*m[i]).Name() > (*m[j]).Name()
}
