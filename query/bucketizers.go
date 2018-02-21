package query

import "fmt"

type (
	bucketizer string
)

var (
	// BUCKETIZERS

	// BucketizerSum computes the sum of all the values found in the interval to bucketize
	BucketizerSum bucketizer = "bucketizer.sum"

	// BucketizerMax returns the max of all the values found on the interval to bucketize
	BucketizerMax bucketizer = "bucketizer.max"

	// BucketizerMin returns the min of all the values found in the interval to bucketize
	BucketizerMin bucketizer = "bucketizer.min"

	// BucketizerMean returns the mean of all the values found in the interval to bucketize
	BucketizerMean bucketizer = "bucketizer.mean"

	// BucketizerMeanCircular push a mapper onto the stack which can then be used to compute the circular mean of all the values found in each sliding window
	BucketizerMeanCircular = func(i float64) bucketizer {
		return bucketizer(fmt.Sprintf("%f bucketizer.mean.circular", i))
	}

	// BucketizerMeanCircularExcludeNulls push a mapper onto the stack which can then be used to compute the circular mean of all the values found in each sliding window. The associated location is the centroid of all the encountered locations. The associated elevation is the mean of the encountered elevations
	BucketizerMeanCircularExcludeNulls = func(i float64) bucketizer {
		return bucketizer(fmt.Sprintf("%f bucketizer.mean.circular.exclude-nulls", i))
	}

	// BucketizerFirst returns the first value of the interval to bucketize with its associated location and elevation
	BucketizerFirst bucketizer = "bucketizer.first"

	// BucketizerLast returns the last value of the interval to bucketize with its associated location and elevation
	BucketizerLast bucketizer = "bucketizer.last"

	// BucketizerJoin outputs for the interval to bucketize of the Geo Time SeriesTM the concatenation of the string representation of values separated by the join string
	BucketizerJoin = func(s string) bucketizer {
		return bucketizer(fmt.Sprintf("'%s' bucketizer.join", s))
	}

	// BucketizerMedian returns the median of all the values found in the interval to bucketize
	BucketizerMedian bucketizer = "bucketizer.median"

	// BucketizerCount computes the number of non-null values found in the interval to bucketize
	BucketizerCount bucketizer = "bucketizer.count"

	// BucketizerAnd applies the logical operator AND on all the values found in the interval to bucketize
	BucketizerAnd bucketizer = "bucketizer.and"

	// BucketizerOr applies the logical operator OR on all the values found in the interval to bucketize
	BucketizerOr bucketizer = "bucketizer.or"
)
