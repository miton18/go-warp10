package query

import "fmt"

type (
	reducer string
)

var (

	// REDUCERS

	// ReducerArgMax outputs for each tick, the tick and the concatenation separated by ‘,’ of the values of the labels for which the value is the maximum of Geo Time SeriesTM which are in the same equivalence class
	ReducerArgMax = func(i float64) reducer {
		return reducer(fmt.Sprintf("%f reducer.argmax", i))
	}

	// ReducerArgMin outputs for each tick, the tick and the concatenation separated by ‘,’ with the values of the labels for which the value is the minimum of Geo Time SeriesTM which are in the same equivalence class
	ReducerArgMin = func(i float64) reducer {
		return reducer(fmt.Sprintf("%f reducer.argmin", i))
	}

	// ReducerCount computes for each tick the number of measures of Geo Time SeriesTM which are in the same equivalence class
	ReducerCount reducer = "reducer.count"

	// ReducerCountExcludeNulls computes at each tick the number of measures of Geo Time SeriesTM which are in the same equivalence class
	ReducerCountExcludeNulls reducer = "reducer.count.exclude-nulls"

	// ReducerCountIncludeNulls computes at each tick the number of measures of Geo Time SeriesTM which are in the same equivalence class
	ReducerCountIncludeNulls reducer = "reducer.count.include-nulls"

	// ReducerJoin outputs for each tick of Geo Time SeriesTM which are in the same equivalence class, the concatenation of the string representation of values separated by the join string
	ReducerJoin = func(s string) reducer {
		return reducer(fmt.Sprintf("%s reducer.join", s))
	}

	// ReducerJoinForbidNulls outputs for each tick of Geo Time SeriesTM which are in the same equivalence class, the concatenation of the string representation of values separated by the join string
	ReducerJoinForbidNulls = func(s string) reducer {
		return reducer(fmt.Sprintf("%s reducer.join.forbid-nulls", s))
	}

	// ReducerMax outputs for each tick the maximum value of Geo Time SeriesTM which are in the same equivalence class.It operates on any type
	ReducerMax reducer = "reducer.max"

	// ReducerMaxForbidNulls outputs for each tick the maximum value of Geo Time SeriesTM which are in the same equivalence class
	ReducerMaxForbidNulls reducer = "reducer.max.forbid-nulls"

	// ReducerMean outputs for each tick the mean of the values of Geo Time SeriesTM which are in the same equivalence class
	ReducerMean reducer = "reducer.mean"

	// ReducerMeanExcludeNulls outputs for each tick the mean of the values of Geo Time SeriesTM which are in the same equivalence class, excluding nulls from the computation
	ReducerMeanExcludeNulls reducer = "reducer.mean.exclude-nulls"

	// ReducerMeanCircular push a mapper onto the stack which can then be used to compute the circular mean of all the values found in each sliding window. The associated location is the centroid of all the encountered locations. The associated elevation is the mean of the encountered elevations
	ReducerMeanCircular reducer = "reducer.mean.circular"

	// ReducerMeanCircularExcludeNulls push a mapper onto the stack which can then be used to compute the circular mean of all the values found in each sliding window. The associated location is the centroid of all the encountered locations. The associated elevation is the mean of the encountered elevations
	ReducerMeanCircularExcludeNulls reducer = "reducer.mean.circular.exclude-nulls"

	// ReducerMedian outputs for each tick the median of the values of Geo Time SeriesTM which are in the same equivalence class
	ReducerMedian reducer = "reducer.median"

	// ReducerMin outputs for each tick the minimum value of Geo Time SeriesTM which are in the same equivalence class. It operates on any type
	ReducerMin reducer = "reducer.min"

	// ReducerMinForbidNulls outputs for each tick the minimum value of Geo Time SeriesTM which are in the same equivalence class
	ReducerMinForbidNulls reducer = "reducer.min.forbid-nulls"

	// ReducerAnd outputs the result of the logical operator AND for each tick value of Geo Time SeriesTM which are in the same equivalence class
	ReducerAnd reducer = "reducer.and"

	// ReducerAndExcludeNulls outputs the result of the logical operator AND for each tick value of Geo Time SeriesTM which are in the same equivalence class, excluding nulls from the computation
	ReducerAndExcludeNulls reducer = "reducer.and.exclude-nulls"

	// ReducerOr outputs the result of the logical operator OReducer for each tick value of Geo Time SeriesTM which are in the same equivalence class
	ReducerOr reducer = "reducer.or"

	// ReducerOrExcludeNulls outputs the result of the logical operator OReducer for each tick value of Geo Time SeriesTM which are in the same equivalence class, excluding nulls from the computation
	ReducerOrExcludeNulls reducer = "reducer.or.exclude-nulls"

	// ReducerSd outputs for each tick the standard deviation of the values of Geo Time SeriesTM which are in the same equivalence class
	ReducerSd = func(b bool) reducer {
		return reducer(fmt.Sprintf("%t reducer.sd", b))
	}

	// ReducerShannonentropy0 computes the Shannon entropy of the values it receives from the framework ReducerEDUCE at each tick
	ReducerShannonentropy0 reducer = "reducer.shannonentropy.0"

	// ReducerShannonentropy1 computes the Shannon entropy of the values it receives from the framework ReducerEDUCE at each tick
	ReducerShannonentropy1 reducer = "reducer.shannonentropy.1"

	// ReducerSum computes at each tick the sum of the values of Geo Time SeriesTM which are in the same equivalence class
	ReducerSum reducer = "reducer.sum"

	// ReducerSumForbidNulls computes at each tick the sum of the values of Geo Time SeriesTM which are in the same equivalence class
	ReducerSumForbidNulls reducer = "reducer.sum.forbid-nulls"

	// ReducerVar outputs for each tick the variance of the values of Geo Time SeriesTM which are in the same equivalence class
	ReducerVar = func(b bool) reducer {
		return reducer(fmt.Sprintf("%t reducer.var", b))
	}
)
