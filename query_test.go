package warp10

import (
	"reflect"
	"testing"
)

func TestClient_NewQuery(t *testing.T) {
	tests := []struct {
		name string
		c    *Client
		want *Query
	}{
	// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := &Client{}
			if got := c.NewQuery(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Client.NewQuery() = %v, want %v", got, tt.want)
			}
		})
	}
}
