package base

import (
	"context"
	"errors"
	"fmt"
	"io"
	"net/http"
	"strconv"
	"strings"
	"time"
)

// Exec execute a WarpScript on the backend, returning resultat as byte array
func (c *Client) Exec(ctx context.Context, warpScript string) (*ExecResult, error) {
	execRes := &ExecResult{}
	r := strings.NewReader(warpScript)

	req, err := http.NewRequestWithContext(ctx, "POST", c.Host+c.ExecPath, r)
	if err != nil {
		return execRes, err
	}

	res, err := c.HTTPClient.Do(req)
	execRes.response = res
	if err != nil {
		return execRes, err
	}

	if res.StatusCode != http.StatusOK {

		if errMsg := execRes.ErrorMessage(); errMsg != "" {
			return execRes, errors.New(errMsg)
		}

		if body, err := io.ReadAll(res.Body); err == nil {
			return execRes, errors.New(string(body))
		}

		return execRes, fmt.Errorf("received status code '%s'", res.Status)
	}

	bts, err := io.ReadAll(res.Body)
	if err != nil {
		return execRes, err
	}

	execRes.raw = bts

	return execRes, res.Body.Close()
}

type ExecResult struct {
	response *http.Response
	raw      []byte
}

func (e *ExecResult) Raw() []byte { return e.raw }

func (e *ExecResult) getHeader(name, defaultValue string) string {
	if e.response == nil {
		return defaultValue
	}
	if e.response.Header == nil {
		return defaultValue
	}

	return e.response.Header.Get(name)
}

func (e *ExecResult) Elapsed() int64 {
	h := e.getHeader(HeaderElapsed, "0")
	i, _ := strconv.ParseInt(h, 10, 64)
	return i
}

func (e *ExecResult) ErrorLine() int64 {
	h := e.getHeader(HeaderErrorLine, "0")
	i, _ := strconv.ParseInt(h, 10, 64)
	return i
}

// You should not use this method as the WarpScript error should be returned as 2nd method return parameter
func (e *ExecResult) ErrorMessage() string {
	return e.getHeader(HeaderErrorMessage, "")
}

func (e *ExecResult) Operations() int64 {
	h := e.getHeader(HeaderOperations, "0")
	i, _ := strconv.ParseInt(h, 10, 64)
	return i
}

func (e *ExecResult) Fetched() int64 {
	h := e.getHeader(HeaderFetched, "0")
	i, _ := strconv.ParseInt(h, 10, 64)
	return i
}

func (e *ExecResult) TimeUnit() int64 {
	h := e.getHeader(HeaderTimeUNit, "0")
	i, _ := strconv.ParseInt(h, 10, 64)
	return i
}

// Find execute a WarpScript on the backend, returning resultat as byte array
// Enhance parse response
func (c *Client) Find(selector Selector) ([]byte, error) {
	req, err := http.NewRequest("GET", c.Host+c.FindPath, nil)
	if err != nil {
		return nil, err
	}

	err = needReadAccess(req, c)
	if err != nil {
		return nil, err
	}

	q := req.URL.Query()
	q.Add("selector", string(selector))
	q.Add("format", "fulltext")
	req.URL.RawQuery = q.Encode()

	res, err := c.HTTPClient.Do(req)
	if err != nil {
		return nil, err
	}

	bts, err := io.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}

	if res.StatusCode != http.StatusOK {
		return nil, errors.New(string(bts))
	}

	return bts, res.Body.Close()
}

// Update execute a WarpScript on the backend, returning resultat as byte array
func (c *Client) Update(gts GTSList) error {
	r := strings.NewReader(gts.Sensision())

	req, err := http.NewRequest("POST", c.Host+c.UpdatePath, r)
	if err != nil {
		return err
	}

	err = needWriteAccess(req, c)
	if err != nil {
		return err
	}

	res, err := c.HTTPClient.Do(req)
	if err != nil {
		return err
	}

	bts, err := io.ReadAll(res.Body)
	if err != nil {
		return err
	}

	if res.StatusCode != http.StatusOK {
		return errors.New(string(bts))
	}

	return res.Body.Close()
}

// Meta execute a WarpScript on the backend, returning resultat as byte array
// Enhance parse response
func (c *Client) Meta(gtsList GTSList) ([]byte, error) {
	r := strings.NewReader(gtsList.SensisionSelectors(true))

	req, err := http.NewRequest("POST", c.Host+c.MetaPath, r)
	if err != nil {
		return nil, err
	}

	err = needWriteAccess(req, c)
	if err != nil {
		return nil, err
	}

	res, err := c.HTTPClient.Do(req)
	if err != nil {
		return nil, err
	}

	bts, err := io.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}

	if res.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("headers: %+v\nBody: %+v", res.Header, string(bts))
	}

	return nil, res.Body.Close()
}

// Fetch execute a WarpScript on the backend, returning resultat as byte array
// Enhance parse response
func (c *Client) Fetch(selector Selector, start time.Time, stop time.Time) ([]byte, error) {
	req, err := http.NewRequest("GET", c.Host+c.FetchPath, nil)
	if err != nil {
		return nil, err
	}

	err = needReadAccess(req, c)
	if err != nil {
		return nil, err
	}

	q := req.URL.Query()
	q.Add("selector", string(selector))
	q.Add("format", "fulltext")
	q.Add("start", start.UTC().Format(time.RFC3339Nano))
	q.Add("stop", stop.UTC().Format(time.RFC3339Nano))
	req.URL.RawQuery = q.Encode()

	res, err := c.HTTPClient.Do(req)
	if err != nil {
		return nil, err
	}

	bts, err := io.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}

	if res.StatusCode != http.StatusOK {
		return nil, errors.New(string(bts))
	}

	err = res.Body.Close()
	if err != nil {
		return nil, err
	}

	return bts, nil
}

// Delete execute a WarpScript on the backend, returning resultat as byte array
// if start and end are nil, assume user want to delete all datapoints
// Enhance parse response
func (c *Client) Delete(selector Selector, start time.Time, stop time.Time) ([]byte, error) {
	req, err := http.NewRequest("GET", c.Host+c.DeletePath, nil)
	if err != nil {
		return nil, err
	}

	err = needWriteAccess(req, c)
	if err != nil {
		return nil, err
	}

	q := req.URL.Query()
	q.Add("selector", string(selector))
	if start.IsZero() && stop.IsZero() {
		q.Add("deleteall", "")
	} else {
		q.Add("start", strconv.FormatInt(start.UnixNano()/1000, 10))
		q.Add("end", strconv.FormatInt(stop.UnixNano()/1000, 10))
	}
	req.URL.RawQuery = q.Encode()

	res, err := c.HTTPClient.Do(req)
	if err != nil {
		return nil, err
	}

	bts, err := io.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}

	if res.StatusCode != http.StatusOK {
		return nil, errors.New(string(bts))
	}

	return nil, res.Body.Close()
}

func needReadAccess(req *http.Request, c *Client) error {
	if c.ReadToken == "" {
		return NoReadTokenError
	}
	req.Header.Add(c.Warp10Header, c.ReadToken)
	return nil
}

func needWriteAccess(req *http.Request, c *Client) error {
	if c.WriteToken == "" {
		return NoWriteTokenError
	}
	req.Header.Add(c.Warp10Header, c.WriteToken)
	return nil
}
