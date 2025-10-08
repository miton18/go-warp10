package base

import (
	"time"
)

// Add a new point to a GTS
func (dps *Datapoints) Add(ts time.Time, value any) *Datapoints {
	if dps == nil {
		dps = &Datapoints{}
	}

	*dps = append(*dps, []any{ts.UnixNano() / 1000, value})

	return dps
}

// AddWithGeo a new point to a GTS with geolocation
func (dps *Datapoints) AddWithGeo(ts time.Time, lattitude, longitude, altitude float64, value any) *Datapoints {
	if dps == nil {
		dps = &Datapoints{}
	}

	*dps = append(*dps, []any{ts.UnixNano() / 1000, lattitude, longitude, altitude, value})

	return dps
}

// Has look for datapoint at this time
func (dps *Datapoints) Has(ts time.Time) bool {
	if dps == nil {
		dps = &Datapoints{}
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
func (dps *Datapoints) Remove(ts time.Time) *Datapoints {
	if dps == nil {
		dps = &Datapoints{}
		return dps
	}

	for i, dp := range *dps {
		if dp[0].(int64) == ts.UnixNano()/1000 {
			*dps = append((*dps)[:i-1], (*dps)[i+1:]...)
			return dps
		}
	}

	return dps
}
