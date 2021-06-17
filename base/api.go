package base

import (
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
	"time"
	"strconv"
)

// Exec execute a WarpScript on the backend, returning resultat as byte array
func (c *Client) Exec(warpScript string) ([]byte, error) {
	r := strings.NewReader(warpScript)

	req, err := http.NewRequest("POST", c.Host+c.ExecPath, r)
	if err != nil {
		return nil, err
	}

	res, err := c.HTTPClient.Do(req)
	if err != nil {
		return nil, err
	}

	if res.StatusCode != http.StatusOK {
		return nil, errors.New(res.Header.Get(HeaderErrorMessage))
	}

	bts, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}

	return bts, res.Body.Close()
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

	bts, err := ioutil.ReadAll(res.Body)
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

	bts, err := ioutil.ReadAll(res.Body)
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

	bts, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}

	if res.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("Headers: %+v\nBody: %+v", res.Header, string(bts))
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

	bts, err := ioutil.ReadAll(res.Body)
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

	bts, err := ioutil.ReadAll(res.Body)
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
