package base

import (
	"time"
)

// Add a new point to a GTS
func (dps *Datapoints) Add(ts time.Time, value interface{}) {
	if dps == nil {
		dps = &Datapoints{}
	}

	*dps = append(*dps, []interface{}{
		ts.UnixNano() / 1000,
		nil,
		nil,
		nil,
		value,
	})
}

// AddWithGeo a new point to a GTS with geolocation
func (dps *Datapoints) AddWithGeo(ts time.Time, lattitude, longitude, altitude float64, value interface{}) {
	if dps == nil {
		dps = &Datapoints{}
	}

	*dps = append(*dps, []interface{}{
		ts.UnixNano() / 1000,
		lattitude,
		longitude,
		altitude,
		value,
	})
}

// Has look for datapoint at this time
func (dps *Datapoints) Has(ts time.Time) bool {
	if dps == nil {
		return false
	}

	for _, dp := range *dps {
		if dp[0].(int64) == ts.UnixNano()/1000 {
			return true
		}
	}
	return false
}

// Remove datapoint at this time (only the first)
func (dps *Datapoints) Remove(ts time.Time) {
	if dps == nil {
		return
	}

	for i, dp := range *dps {
		if dp[0].(int64) == ts.UnixNano()/1000 {
			*dps = append((*dps)[:i-1], (*dps)[i+1:]...)
			return
		}
	}
}
