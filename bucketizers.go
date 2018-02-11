package warp10

import "fmt"

type (
	bucketizer string
)

var (
	// BUCKETIZERS

	// BSum computes the sum of all the values found in the interval to bucketize
	BSum bucketizer = "bucketizer.sum"
	// BMax returns the max of all the values found on the interval to bucketize
	BMax bucketizer = "bucketizer.max"
	// BMin returns the min of all the values found in the interval to bucketize
	BMin bucketizer = "bucketizer.min"
	// BMean returns the mean of all the values found in the interval to bucketize
	BMean bucketizer = "bucketizer.mean"
	// BMeanCircular push a mapper onto the stack which can then be used to compute the circular mean of all the values found in each sliding window
	BMeanCircular = func(i float64) bucketizer {
		return bucketizer(fmt.Sprintf("%f bucketizer.mean.circular", i))
	}
	// BMeanCircularExcludeNulls push a mapper onto the stack which can then be used to compute the circular mean of all the values found in each sliding window. The associated location is the centroid of all the encountered locations. The associated elevation is the mean of the encountered elevations
	BMeanCircularExcludeNulls = func(i float64) bucketizer {
		return bucketizer(fmt.Sprintf("%f bucketizer.mean.circular.exclude-nulls", i))
	}
	// BFirst returns the first value of the interval to bucketize with its associated location and elevation
	BFirst bucketizer = "bucketizer.first"
	// BLast returns the last value of the interval to bucketize with its associated location and elevation
	BLast bucketizer = "bucketizer.last"
	// BJoin outputs for the interval to bucketize of the Geo Time SeriesTM the concatenation of the string representation of values separated by the join string
	BJoin = func(s string) bucketizer {
		return bucketizer(fmt.Sprintf("'%s' bucketizer.join", s))
	}
	// BMedian returns the median of all the values found in the interval to bucketize
	BMedian bucketizer = "bucketizer.median"
	// BCount computes the number of non-null values found in the interval to bucketize
	BCount bucketizer = "bucketizer.count"
	// BAnd applies the logical operator AND on all the values found in the interval to bucketize
	BAnd bucketizer = "bucketizer.and"
	// BOr applies the logical operator OR on all the values found in the interval to bucketize
	BOr bucketizer = "bucketizer.or"
)
