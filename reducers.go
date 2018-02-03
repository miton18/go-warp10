package warp10

import "fmt"

var (

	// REDUCERS

	// RArgMax outputs for each tick, the tick and the concatenation separated by ‘,’ of the values of the labels for which the value is the maximum of Geo Time SeriesTM which are in the same equivalence class
	RArgMax = func(i float64) reducer {
		return reducer(fmt.Sprintf("%f reducer.argmax", i))
	}
	// RArgMin outputs for each tick, the tick and the concatenation separated by ‘,’ with the values of the labels for which the value is the minimum of Geo Time SeriesTM which are in the same equivalence class
	RArgMin = func(i float64) reducer {
		return reducer(fmt.Sprintf("%f reducer.argmin", i))
	}
	// RCount computes for each tick the number of measures of Geo Time SeriesTM which are in the same equivalence class
	RCount reducer = "reducer.count"
	// RCountExcludeNulls computes at each tick the number of measures of Geo Time SeriesTM which are in the same equivalence class
	RCountExcludeNulls reducer = "reducer.count.exclude-nulls"
	// RCountIncludeNulls computes at each tick the number of measures of Geo Time SeriesTM which are in the same equivalence class
	RCountIncludeNulls reducer = "reducer.count.include-nulls"
	// RJoin outputs for each tick of Geo Time SeriesTM which are in the same equivalence class, the concatenation of the string representation of values separated by the join string
	RJoin = func(s string) reducer {
		return reducer(fmt.Sprintf("%s reducer.join", s))
	}
	// RJoinForbidNulls outputs for each tick of Geo Time SeriesTM which are in the same equivalence class, the concatenation of the string representation of values separated by the join string
	RJoinForbidNulls = func(s string) reducer {
		return reducer(fmt.Sprintf("%s reducer.join.forbid-nulls", s))
	}
	// RMax outputs for each tick the maximum value of Geo Time SeriesTM which are in the same equivalence class.It operates on any type
	RMax reducer = "reducer.max"
	// RMaxForbidNulls outputs for each tick the maximum value of Geo Time SeriesTM which are in the same equivalence class
	RMaxForbidNulls reducer = "reducer.max.forbid-nulls"
	// RMean outputs for each tick the mean of the values of Geo Time SeriesTM which are in the same equivalence class
	RMean reducer = "reducer.mean"
	// RMeanExcludeNulls outputs for each tick the mean of the values of Geo Time SeriesTM which are in the same equivalence class, excluding nulls from the computation
	RMeanExcludeNulls reducer = "reducer.mean.exclude-nulls"
	// RMeanCircular push a mapper onto the stack which can then be used to compute the circular mean of all the values found in each sliding window. The associated location is the centroid of all the encountered locations. The associated elevation is the mean of the encountered elevations
	RMeanCircular reducer = "reducer.mean.circular"
	// RMeanCircularExcludeNulls push a mapper onto the stack which can then be used to compute the circular mean of all the values found in each sliding window. The associated location is the centroid of all the encountered locations. The associated elevation is the mean of the encountered elevations
	RMeanCircularExcludeNulls reducer = "reducer.mean.circular.exclude-nulls"
	// RMedian outputs for each tick the median of the values of Geo Time SeriesTM which are in the same equivalence class
	RMedian reducer = "reducer.median"
	// RMin outputs for each tick the minimum value of Geo Time SeriesTM which are in the same equivalence class. It operates on any type
	RMin reducer = "reducer.min"
	// RMinForbidNulls outputs for each tick the minimum value of Geo Time SeriesTM which are in the same equivalence class
	RMinForbidNulls reducer = "reducer.min.forbid-nulls"
	// RAnd outputs the result of the logical operator AND for each tick value of Geo Time SeriesTM which are in the same equivalence class
	RAnd reducer = "reducer.and"
	// RAndExcludeNulls outputs the result of the logical operator AND for each tick value of Geo Time SeriesTM which are in the same equivalence class, excluding nulls from the computation
	RAndExcludeNulls reducer = "reducer.and.exclude-nulls"
	// ROr outputs the result of the logical operator OR for each tick value of Geo Time SeriesTM which are in the same equivalence class
	ROr reducer = "reducer.or"
	// ROrExcludeNulls outputs the result of the logical operator OR for each tick value of Geo Time SeriesTM which are in the same equivalence class, excluding nulls from the computation
	ROrExcludeNulls reducer = "reducer.or.exclude-nulls"
	// RSd outputs for each tick the standard deviation of the values of Geo Time SeriesTM which are in the same equivalence class
	RSd = func(b bool) reducer {
		return reducer(fmt.Sprintf("%t reducer.sd", b))
	}
	// RShannonentropy0 computes the Shannon entropy of the values it receives from the framework REDUCE at each tick
	RShannonentropy0 reducer = "reducer.shannonentropy.0"
	// RShannonentropy1 computes the Shannon entropy of the values it receives from the framework REDUCE at each tick
	RShannonentropy1 reducer = "reducer.shannonentropy.1"
	// RSum computes at each tick the sum of the values of Geo Time SeriesTM which are in the same equivalence class
	RSum reducer = "reducer.sum"
	// RSumForbidNulls computes at each tick the sum of the values of Geo Time SeriesTM which are in the same equivalence class
	RSumForbidNulls reducer = "reducer.sum.forbid-nulls"
	// RVar outputs for each tick the variance of the values of Geo Time SeriesTM which are in the same equivalence class
	RVar = func(b bool) reducer {
		return reducer(fmt.Sprintf("%t reducer.var", b))
	}
)
