package warp10

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
	"time"

	"github.com/moul/http2curl"
)

// Bucketizer is used by BUVKETIZE framework
type (
	bucketizer string
	reducer    string
)

// Query is a WarpScript
type Query struct {
	warpscript string
	client     *Client
	Errs       []error
}

// NewQuery should be called to build a new Warpscript
func (c *Client) NewQuery() *Query {
	return &Query{
		warpscript: "// GENERATED WARPSCRIPT\n",
		client:     c,
	}
}

// Fetch grep some datapoints on the backend
// if token is empty, the client one is used
func (qi *Query) Fetch(token, class string, labels map[string]string, start time.Time, interval time.Duration) *Query {
	if token == "" {
		token = qi.client.ReadToken
	}
	ls, err := json.Marshal(labels)
	if err != nil {
		qi.Errs = append(qi.Errs, errors.New("Fetch(): "+err.Error()))
	}
	qi.warpscript += fmt.Sprintf("[ '%s' '%s' '%s' JSON-> %d %d ] FETCH\n", token, class, ls, start.UnixNano()/1000, interval.Nanoseconds()/1000)
	return qi
}

// Bucketize perform a temporel aggregation
// Expect a GTS array as previous element
func (qi *Query) Bucketize(buc bucketizer, start time.Time, bucket time.Duration, buckets int) *Query {
	qi.warpscript += fmt.Sprintf("[ SWAP %s %d %d %d ] BUCKETIZE\n", buc, start.UnixNano()/1000, bucket.Nanoseconds()/1000, buckets)
	return qi
}

// Reduce perfom a series aggregation
// Expect a GTS array as previous element
func (qi *Query) Reduce(red reducer, equivalentLabels []string) (qo *Query) {
	els := ""
	if len(equivalentLabels) > 0 {
		for _, el := range equivalentLabels {
			els += fmt.Sprintf("'%s' ", el)
		}
		els = strings.TrimSuffix(els, " ")
	}
	qi.warpscript += fmt.Sprintf("[ SWAP [ %s ] %s ] REDUCE\n", els, red)
	return qi
}

// Map apply a transformation on datapoints from each GTS
// Expect a GTS array as previous element
func (qi *Query) Map(token, selector string, start, stop time.Time) *Query {
	return qi
}

// Dedup removes GTS datapoints duplicates
// Expect a GTS array as previous element
func (qi *Query) Dedup() *Query {
	qi.warpscript += "DEDUP\n"
	return qi
}

// NonEmpty retains only GTS that have at least one value
// Expect a GTS array as previous element
func (qi *Query) NonEmpty() *Query {
	qi.warpscript += "NONEMPTY\n"
	return qi
}

// isIncrementalCounter compensates for possible counter resets by adding the last value before the rest to all values after the reset
// Expect a GTS array as previous element
func (qi *Query) isIncrementalCounter() *Query {
	qi.warpscript += "false RESETS\n"
	return qi
}

// isDecrementalCounter compensates for possible counter resets by sub the last value before the rest to all values after the reset
// Expect a GTS array as previous element
func (qi *Query) isDecrementalCounter() *Query {
	qi.warpscript += "true RESETS\n"
	return qi
}

// Sort sort GTS values by timestamp
// Expect a GTS array as previous element
func (qi *Query) Sort() *Query {
	qi.warpscript += "true RESETS\n"
	return qi
}

// Exec send the WarpScript and parse the response
func (qi *Query) Exec() (gtss []GTS, err error) {

	if len(qi.Errs) > 0 {
		errs := []string{}
		for _, err := range qi.Errs {
			errs = append(errs, err.Error())
		}
		return nil, fmt.Errorf("Can't execute query with errors: %s", strings.Join(errs, "\n"))
	}

	r := bytes.NewReader([]byte(qi.warpscript))
	req, err := http.NewRequest("POST", qi.client.Host+qi.client.ExecPath, r)
	if err != nil {
		return
	}

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return
	}
	defer res.Body.Close()

	bts, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return
	}

	var stack [][]GTS
	if err = json.Unmarshal(bts, &stack); err != nil {
		err = errors.New(err.Error() + "\n" + string(bts))
		return
	}

	gtss = stack[0]
	return
}

// Debug output the computed WarpScript
func (qi *Query) Debug() (string, error) {

	r := bytes.NewReader([]byte(qi.warpscript))
	req, err := http.NewRequest("POST", qi.client.Host+qi.client.ExecPath, r)
	if err != nil {
		return "", err
	}

	curl, err := http2curl.GetCurlCommand(req)
	if err != nil {
		return "", err
	}

	return curl.String(), nil
}
