package warp10

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"net/http"
	"sort"
	"strings"
	"sync"
	"time"
)

// PassiveExporter Expose your metrics on a web page in Sensision format
// Tools like Beamium cqn scrappe it
type PassiveExporter struct {
	Metrics Metrics
	listen  string
	path    string
}

// ActiveExporter Push metrics to your favourite Warp10 instance
// You can manually tell this exporter to flush metrics or let a ticker do it for you
type ActiveExporter struct {
	Metrics  Metrics
	Client   *Client
	flushing sync.Mutex
}

// HandlerExporter allow you to add a route on your application
type HandlerExporter struct {
	Metrics Metrics
	Handler http.HandlerFunc
}

// NewPassiveExporter Return an instanciated PassiveExporter
func NewPassiveExporter(listen, path string) (*PassiveExporter, error) {
	pe := PassiveExporter{
		Metrics: Metrics{},
		listen:  listen,
		path:    path,
	}

	server := http.NewServeMux()
	server.HandleFunc(path, func(res http.ResponseWriter, req *http.Request) {
		res.Write([]byte(sensisionFromGTS(pe.Metrics)))
	})
	server.HandleFunc("/", func(res http.ResponseWriter, req *http.Request) {
		res.Write([]byte("<a href=\"" + path + "\">Metrics</a>"))
	})

	go func() {
		if err := http.ListenAndServe(listen, server); err != nil {
			fmt.Println(err)
		}
	}()

	return &pe, nil
}

// NewActiveExporter Return an instanciated ActiveExporter
// Default behaviour is to send batch of metric at Period interval
// If period is nil, you have to manually Flush metrics.
// If the exporter failed to push metrics, it will keep them for the next batch and send you an error in the chan
func NewActiveExporter(warp10Client Client, period *time.Duration) (*ActiveExporter, chan error) {
	ae := &ActiveExporter{
		Metrics:  Metrics{},
		Client:   &warp10Client,
		flushing: sync.Mutex{},
	}
	e := make(chan error)

	if period != nil {
		go func() {
			t := time.NewTicker(*period)
			select {
			case <-t.C:
				if err := ae.Flush(); err != nil {
					e <- err
				}
			}
		}()
	}

	return ae, e
}

// NewHandlerExporter Return an instanciated HandlerExporter
func NewHandlerExporter() *HandlerExporter {
	he := HandlerExporter{
		Metrics: Metrics{},
	}

	he.Handler = func(req http.ResponseWriter, res *http.Request) {
		res.Write(bytes.NewBufferString(sensisionFromGTS(he.Metrics)))
	}
	return &he
}

// Register add a new metric to the bundle
func (pe *PassiveExporter) Register(m Metric) {
	pe.Metrics = append(pe.Metrics, &m)
	sort.Sort(pe.Metrics)
}

// AddMetricToNextBatch add a new metric to the next Flush call or Period
func (ae *ActiveExporter) AddMetricToNextBatch(m Metric) {
	ae.Metrics = append(ae.Metrics, &m)
	sort.Sort(ae.Metrics)
}

// Register add a new metric to the bundle
func (he *HandlerExporter) Register(m Metric) {
	he.Metrics = append(he.Metrics, &m)
	sort.Sort(he.Metrics)
}

// Flush send all metrics
func (ae *ActiveExporter) Flush() error {
	ae.flushing.Lock()
	defer ae.flushing.Unlock()

	r := strings.NewReader(sensisionFromGTS(ae.Metrics))

	req, err := http.NewRequest("POST", ae.Client.Host+ae.Client.UpdatePath, r)
	if err != nil {
		return err
	}

	req.Header.Add(ae.Client.Warp10Header, ae.Client.WriteToken)

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return err
	}
	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
		resBody, err := ioutil.ReadAll(res.Body)
		if err != nil {
			return fmt.Errorf("Failed to send metrics: %s", err.Error())
		}
		return fmt.Errorf("Failed to send metrics: %s", string(resBody))
	}

	ae.Metrics = ae.Metrics[:0]
	return nil
}

func sensisionFromGTS(m Metrics) string {
	p := ""
	for _, mPtr := range m {
		if mPtr != nil {
			metric := (*mPtr)
			for _, gts := range metric.Get() {
				p += "# " + metric.Help() + "\n"
				p += gts.Sensision() + "\n"
			}
		}
	}
	return p
}
