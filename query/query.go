package query

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strings"
	"time"

	b "github.com/miton18/go-warp10/base"
	"github.com/moul/http2curl"
)

// Query is a WarpScript
type Query struct {
	warpscript string
	client     *b.Client
	Errs       []error
}

// NewQuery should be called to build a new Warpscript
func NewQuery(c *b.Client) *Query {
	return &Query{
		warpscript: "// GENERATED WARPSCRIPT\n",
		client:     c,
	}
}

// Fetch grep some datapoints on the backend
// if token is empty, the client one is used
func (q *Query) Fetch(token, class string, labels b.Labels, start time.Time, interval time.Duration) *Query {
	if token == "" {
		token = q.client.ReadToken
	}
	ls, err := json.Marshal(labels)
	if err != nil {
		q.Errs = append(q.Errs, errors.New("Fetch(): "+err.Error()))
	}
	q.warpscript += fmt.Sprintf("[ '%s' '%s' '%s' JSON-> %d %d ] FETCH\n", token, class, ls, start.UnixNano()/1000, interval.Nanoseconds()/1000)
	return q
}

// Bucketize perform a temporel aggregation
// Expect a GTS array as previous element
func (q *Query) Bucketize(buc bucketizer, start time.Time, bucket time.Duration, buckets int) *Query {
	q.warpscript += fmt.Sprintf("[ SWAP %s %d %d %d ] BUCKETIZE\n", buc, start.UnixNano()/1000, bucket.Nanoseconds()/1000, buckets)
	return q
}

// Reduce perfom a series aggregation
// Expect a GTS array as previous element
func (q *Query) Reduce(red reducer, equivalentLabels []string) *Query {
	els := ""
	if len(equivalentLabels) > 0 {
		for _, el := range equivalentLabels {
			els += fmt.Sprintf("'%s' ", el)
		}
		els = strings.TrimSuffix(els, " ")
	}
	q.warpscript += fmt.Sprintf("[ SWAP [ %s ] %s ] REDUCE\n", els, red)
	return q
}

// Map apply a transformation on datapoints from each GTS
// Expect a GTS array as previous element
func (q *Query) Map(token, selector string, start, stop time.Time) *Query {
	return q
}

// Dedup removes GTS datapoints duplicates
// Expect a GTS array as previous element
func (q *Query) Dedup() *Query {
	q.warpscript += "DEDUP\n"
	return q
}

// NonEmpty retains only GTS that have at least one value
// Expect a GTS array as previous element
func (q *Query) NonEmpty() *Query {
	q.warpscript += "NONEMPTY\n"
	return q
}

// isIncrementalCounter compensates for possible counter resets by adding the last value before the rest to all values after the reset
// Expect a GTS array as previous element
func (q *Query) isIncrementalCounter() *Query {
	q.warpscript += "false RESETS\n"
	return q
}

// isDecrementalCounter compensates for possible counter resets by sub the last value before the rest to all values after the reset
// Expect a GTS array as previous element
func (q *Query) isDecrementalCounter() *Query {
	q.warpscript += "true RESETS\n"
	return q
}

// Sort sort GTS values by timestamp
// Expect a GTS array as previous element
func (q *Query) Sort() *Query {
	q.warpscript += "SORT\n"
	return q
}

// Exec send the WarpScript and parse the response
func (q *Query) Exec() (gtsList b.GTSList, err error) {

	if len(q.Errs) > 0 {
		errs := []string{}
		for _, err := range q.Errs {
			errs = append(errs, err.Error())
		}
		return nil, fmt.Errorf("Can't execute query with errors: %s", strings.Join(errs, "\n"))
	}

	body, err := q.client.Exec(q.warpscript)
	if err != nil {
		return
	}

	var stack []b.GTSList
	if err = json.Unmarshal(body, &stack); err != nil {
		return
	}

	gtsList = stack[0]
	return
}

// Debug output the computed WarpScript
func (q *Query) Debug() (string, error) {

	r := bytes.NewReader([]byte(q.warpscript))
	req, err := http.NewRequest("POST", q.client.Host+q.client.ExecPath, r)
	if err != nil {
		return "", err
	}

	curl, err := http2curl.GetCurlCommand(req)
	if err != nil {
		return "", err
	}

	return curl.String(), nil
}

// Debug output the computed WarpScript
func (q *Query) String() string {
	return q.warpscript
}
