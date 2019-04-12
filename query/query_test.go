package query

import (
	"reflect"
	"strings"
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

func TestQuery_Map(t *testing.T) {
	type fields struct {
		warpscript string
		client     *b.Client
		Errs       []error
	}
	type args struct {
		m               mapper
		mapperParameter string
		pre             int
		post            int
		occurences      int
	}
	tests := []struct {
		name string
		args args
		want string
	}{{
		name: "no params",
		args: args{
			m:               MapperRate,
			mapperParameter: "",
			pre:             1,
			post:            0,
		},
		want: "[ SWAP  mapper.rate 1 0 0 ] MAP\n",
	}, {
		name: "with params",
		args: args{
			m:               MapperAdd,
			mapperParameter: "10.0",
			pre:             0,
			post:            0,
			occurences:      0,
		},
		want: "[ SWAP 10.0 mapper.add 0 0 0 ] MAP\n",
	}}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			q := NewQuery(nil)

			got := q.Map(tt.args.m, tt.args.mapperParameter, tt.args.pre, tt.args.post, tt.args.occurences)

			if !strings.HasSuffix(got.String(), tt.want) {
				t.Errorf("Query.Map() = %v, want %v", got, tt.want)
			}
		})
	}
}
