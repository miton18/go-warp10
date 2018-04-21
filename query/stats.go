package query

import b "github.com/miton18/go-warp10/base"

type GTSStatsResult struct {
	GTSCount                 uint32            `json:"gts.estimate"`
	GTSCountByClass          map[string]uint32 `json:"per.class.estimate"`
	ClassCount               uint32            `json:"classes.estimate"`
	LabelsNameCount          uint32            `json:"labelnames.estimate"`
	LabelsValuesCount        uint32            `json:"labelvalues.estimate"`
	LabelsValuesCountByLabel map[string]uint32 `json:"per.label.value.estimate"`
	Approximation            float64           `json:"error.rate"`
}

func GTSStats(client *b.Client, class b.Selector, labels b.Labels) {


	
}
