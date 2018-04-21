package instrumentation

import (
	"reflect"
	"testing"
)

func TestGauge(t *testing.T) {
	var c *Gauge

	t.Run("Create gauge", func(t *testing.T) {
		c = NewGauge("jobs.current", nil, "Number of running jobs")
		if c == nil {
			t.Error("Gauge is nil")
		}
	})

	t.Run("Initial gauge value is 0", func(t *testing.T) {
		r := c.Get()
		if len(r) == 0 {
			t.Error("Gauge has no GTS")
		}
		g := r[0]
		if len(g.Values) == 0 {
			t.Error("Gauge has no value")
		}
		dp := g.Values[0]
		if dp[0] != float64(0) {
			t.Errorf("Gauge has not 0 as default value: %v (%s)", dp[0], reflect.TypeOf(dp[0]))
		}
	})

	t.Run("Inc Gauge by 1", func(t *testing.T) {
		c.Inc()
		v := c.Get()[0].Values[0][0]
		if v != float64(1) {
			t.Errorf("Gauge is not incremented by 1: %v (%s)", v, reflect.TypeOf(v))
		}
	})

	t.Run("Add 10 to Gauge", func(t *testing.T) {
		c.Add(10)
		v := c.Get()[0].Values[0][0]
		if v != float64(11) {
			t.Errorf("Gauge is not incremented by 10: %v (%s)", v, reflect.TypeOf(v))
		}
	})

	t.Run("Dec Gauge by 1", func(t *testing.T) {
		c.Dec()
		v := c.Get()[0].Values[0][0]
		if v != float64(10) {
			t.Errorf("Gauge value is not 10: %v (%s)", v, reflect.TypeOf(v))
		}
	})

	t.Run("Sub Gauge by 10", func(t *testing.T) {
		c.Sub(5)
		v := c.Get()[0].Values[0][0]
		if v != float64(5) {
			t.Errorf("Gauge value is not 5: %v (%s)", v, reflect.TypeOf(v))
		}
	})

	t.Run("Set Gauge value", func(t *testing.T) {
		c.Set(50)
		v := c.Get()[0].Values[0][0]
		if v != float64(50) {
			t.Errorf("Gauge value is not 5: %v (%s)", v, reflect.TypeOf(v))
		}
	})

	t.Run("Reset Gauge", func(t *testing.T) {
		c.Reset()
		v := c.Get()[0].Values[0][0]
		if v != float64(0) {
			t.Errorf("Gauge has not been resetted: %v (%s)", v, reflect.TypeOf(v))
		}
	})
}
