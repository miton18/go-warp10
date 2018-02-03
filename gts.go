package warp10

// GTS is the Warp10 representation of a GeoTimeSerie (GTS)
type GTS struct {
	// Name of the time serie
	ClassName string `json:"c"`
	// Key/value of the GTS labels (changing one key or his value create a new GTS)
	Labels map[string]string `json:"l"`
	// Key/value of the GTS attributes (can be setted/updated/removed without creating a new GTS)
	Attributes map[string]string `json:"a"`
	// Timestamp of the last datapoint received on this GTS (% last activity window)
	LastActivity int64 `json:"la"`
	// Array of datapoints of this GTS
	Values [][]interface{} `json:"v"`
}

// GTSList is an array of GTS
type GTSList []*GTS
