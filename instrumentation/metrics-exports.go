package instrumentation

import (
	"bytes"
	"fmt"
	"net/http"
	"sort"
	"sync"
	"time"

	b "github.com/miton18/go-warp10/base"
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
	Client   *b.Client
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
func NewActiveExporter(warp10Client b.Client, period *time.Duration) (*ActiveExporter, chan error) {
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
	pe.Metrics = append(pe.Metrics, m)
	sort.Sort(pe.Metrics)
}

// AddMetricToNextBatch add a new metric to the next Flush call or Period
func (ae *ActiveExporter) AddMetricToNextBatch(m Metric) {
	ae.Metrics = append(ae.Metrics, m)
	sort.Sort(ae.Metrics)
}

// Register add a new metric to the bundle
func (he *HandlerExporter) Register(m Metric) {
	he.Metrics = append(he.Metrics, m)
	sort.Sort(he.Metrics)
}

// Flush send all metrics
func (ae *ActiveExporter) Flush() error {
	ae.flushing.Lock()
	defer ae.flushing.Unlock()

	if err := ae.Client.Update(getGTSList(ae.Metrics)); err != nil {
		return fmt.Errorf("Failed to send metrics: %s", err.Error())
	}

	ae.Metrics = ae.Metrics[:0]
	return nil
}

func getGTSList(m Metrics) b.GTSList {
	gl := b.GTSList{}
	for _, mPtr := range m {
		if mPtr != nil {
			metric := mPtr
			for _, gts := range metric.Get() {
				gl = append(gl, gts)
			}
		}
	}
	return gl
}

func sensisionFromGTS(m Metrics) string {
	p := ""
	for _, mPtr := range m {
		if mPtr != nil {
			metric := mPtr
			for _, gts := range metric.Get() {
				p += "# " + metric.Help() + "\n"
				p += gts.Sensision() + "\n"
			}
		}
	}
	return p
}
