package query

import (
	"reflect"
	"testing"

	b "github.com/miton18/go-warp10/base"
)

func TestClient_NewQuery(t *testing.T) {
	tests := []struct {
		name string
		c    *b.Client
		want string
	}{{
		name: "Empty query",
		c:    b.NewClient(""),
		want: "// GENERATED WARPSCRIPT\n",
	}}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := &b.Client{}

			got := NewQuery(c).String()
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Client.NewQuery() = %v, want %v", got, tt.want)
			}
		})
	}
}
