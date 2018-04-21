package instrumentation

import (
	"reflect"
	"testing"

	b "github.com/miton18/go-warp10/base"
)

func TestCounter(t *testing.T) {
	var c *Counter

	t.Run("Create counter", func(t *testing.T) {
		c = NewCounter("http.requests", b.Labels{"path": "/login"}, "Count of HTTP call on /login resource")
		if c == nil {
			t.Error("Counter is nil")
		}
	})

	t.Run("Initial counter value is 0", func(t *testing.T) {
		r := c.Get()
		if len(r) == 0 {
			t.Error("Counter has no GTS")
		}
		g := r[0]
		if len(g.Values) == 0 {
			t.Error("Counter has no value")
		}
		dp := g.Values[0]
		if dp[0] != uint64(0) {
			t.Errorf("Counter has not 0 as default value: %v (%s)", dp[0], reflect.TypeOf(dp[0]))
		}
	})

	t.Run("Inc Counter by 1", func(t *testing.T) {
		c.Inc()
		v := c.Get()[0].Values[0][0]
		if v != uint64(1) {
			t.Errorf("Counter is not incremented by 1: %v (%s)", v, reflect.TypeOf(v))
		}
	})

	t.Run("Add 10 to Counter", func(t *testing.T) {
		c.Add(10)
		v := c.Get()[0].Values[0][0]
		if v != uint64(11) {
			t.Errorf("Counter is not incremented by 10: %v (%s)", v, reflect.TypeOf(v))
		}
	})

	t.Run("Reset Counter", func(t *testing.T) {
		c.Reset()
		v := c.Get()[0].Values[0][0]
		if v != uint64(0) {
			t.Errorf("Counter has not been resetted: %v (%s)", v, reflect.TypeOf(v))
		}
	})
}
