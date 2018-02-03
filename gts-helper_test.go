package warp10

import (
	"reflect"
	"testing"

	"github.com/go-test/deep"
)

func TestNewEmptyGTS(t *testing.T) {
	tests := []struct {
		name string
		want *GTS
	}{{
		name: "New empty GTS",
		want: &GTS{
			ClassName:    "",
			Labels:       map[string]string{},
			Attributes:   map[string]string{},
			LastActivity: 0,
			Values:       [][]interface{}{},
		},
	}}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := NewEmptyGTS()
			if pErr := deep.Equal(got, tt.want); len(pErr) > 0 {
				t.Errorf("NewEmptyGTS() = %+v, want %+v\n%v", got, tt.want, pErr)
			}
		})
	}
}

func TestNewGTS(t *testing.T) {
	className := "my.metric"
	g := NewGTS(className)

	if g.ClassName != className {
		t.Fatalf("NewGTS expect %s instead of %s as name", className, g.ClassName)
	}
}

func TestNewGTSWithLabels(t *testing.T) {
	type args struct {
		className string
		labels    map[string]string
	}
	tests := []struct {
		name string
		args args
		want *GTS
	}{{
		name: "New GTS with labels",
		args: args{
			className: "my.metric",
			labels: map[string]string{
				"a": "b",
				"c": "d",
			},
		},
		want: &GTS{
			ClassName: "my.metric",
			Labels: map[string]string{
				"a": "b",
				"c": "d",
			},
		},
	}}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := NewGTSWithLabels(tt.args.className, tt.args.labels); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewGTSWithLabels() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestParseGTSFromString(t *testing.T) {
	type args struct {
		sensisionLine string
	}
	tests := []struct {
		name    string
		args    args
		wantGts *GTS
		wantErr bool
	}{{
		name: "simple datapoint",
		args: args{
			sensisionLine: "1234// my.metric{} 10",
		},
		wantGts: &GTS{
			ClassName:    "my.metric",
			Labels:       map[string]string{},
			Attributes:   map[string]string{},
			LastActivity: 0,
			Values:       [][]interface{}{{int64(1234), float64(0), float64(0), float64(0), 10}},
		},
		wantErr: false,
	}, {
		name: "labelled datapoint",
		args: args{
			sensisionLine: "1234// my.metric{a=b,c=1} 10",
		},
		wantGts: &GTS{
			ClassName:    "my.metric",
			Labels:       map[string]string{"a": "b", "c": "1"},
			Attributes:   map[string]string{},
			LastActivity: 0,
			Values:       [][]interface{}{{int64(1234), float64(0), float64(0), float64(0), 10}},
		},
		wantErr: false,
	}, {
		name: "string datapoint",
		args: args{
			sensisionLine: "1234// my.metric{} 'my awesome metric'",
		},
		wantGts: &GTS{
			ClassName:    "my.metric",
			Labels:       map[string]string{},
			Attributes:   map[string]string{},
			LastActivity: 0,
			Values:       [][]interface{}{{int64(1234), float64(0), float64(0), float64(0), "my awesome metric"}},
		},
		wantErr: false,
	}, {
		name: "Geo datapoint",
		args: args{
			sensisionLine: "1234/12.3456:4.345678/1230 my.metric{} 10",
		},
		wantGts: &GTS{
			ClassName:    "my.metric",
			Labels:       map[string]string{},
			Attributes:   map[string]string{},
			LastActivity: 0,
			Values:       [][]interface{}{{int64(1234), float64(12.3456), float64(4.345678), float64(1230), 10}},
		},
		wantErr: false,
	}, {
		name: "No value datapoint",
		args: args{
			sensisionLine: "1234// my.metric{}",
		},
		wantGts: &GTS{},
		wantErr: true,
	}, {
		name: "No timestamp datapoint",
		args: args{
			sensisionLine: "// my.metric{} 10",
		},
		wantGts: &GTS{},
		wantErr: true,
	}}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotGts, err := ParseGTSFromString(tt.args.sensisionLine)

			if tt.wantErr {
				if err == nil {
					t.Fatalf("ParseGTSFromString() wantErr %v", tt.wantErr)
				}
				return
			} else if gotGts == nil {
				t.Fatalf("ParseGTSFromString() = return is undefined")
				return
			}

			if gotGts == nil {
				t.Fatalf("ParseGTSFromString() = return is nil")
				return
			}

			if differences := deep.Equal(*gotGts, *tt.wantGts); differences != nil {
				t.Fatalf("ParseGTSFromString() = %+v, want %+v\n%v", gotGts, tt.wantGts, differences)
			}
		})
	}
}
func TestParseGTSFromBytes(t *testing.T) {
	type args struct {
		in []byte
	}
	tests := []struct {
		name    string
		args    args
		wantGts *GTS
		wantErr bool
	}{{
		name: "Parse bytes sensision metric",
		args: args{
			in: []byte("1234// my.metric{} 50"),
		},
		wantGts: &GTS{
			ClassName:    "my.metric",
			Labels:       map[string]string{},
			Attributes:   map[string]string{},
			LastActivity: 0,
			Values:       [][]interface{}{{int64(1234), float64(0), float64(0), float64(0), 50}},
		},
		wantErr: false,
	}}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotGts, err := ParseGTSFromBytes(tt.args.in)
			if (err != nil) != tt.wantErr {
				t.Errorf("ParseGTSFromBytes() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if pErr := deep.Equal(gotGts, tt.wantGts); len(pErr) > 0 {
				t.Errorf("ParseGTSFromBytes() = %+v, want %+v %v", gotGts, tt.wantGts, pErr)
			}
		})
	}
}

func TestParseGTSArrayFromString(t *testing.T) {
	type args struct {
		in string
	}
	tests := []struct {
		name     string
		args     args
		wantGtss []*GTS
		wantErr  bool
	}{{
		name: "Parse bytes sensision metrics",
		args: args{
			in: "1234// my.metric{} 50\n5678// my.metric2{} 100",
		},
		wantGtss: []*GTS{&GTS{
			ClassName:    "my.metric",
			Labels:       map[string]string{},
			Attributes:   map[string]string{},
			LastActivity: 0,
			Values:       [][]interface{}{{int64(1234), float64(0), float64(0), float64(0), 50}},
		}, &GTS{
			ClassName:    "my.metric2",
			Labels:       map[string]string{},
			Attributes:   map[string]string{},
			LastActivity: 0,
			Values:       [][]interface{}{{int64(5678), float64(0), float64(0), float64(0), 100}},
		}},
		wantErr: false,
	}, {
		name: "Parse Error bytes sensision metrics",
		args: args{
			in: "1234// my.metric{} 50\n5678/00",
		},
		wantErr: true,
	}}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotGtss, err := ParseGTSArrayFromString(tt.args.in)
			if tt.wantErr {
				if err == nil {
					t.Errorf("ParseGTSArrayFromString() error = %v, wantErr %v", err, tt.wantErr)
				}
				return
			}
			if pErr := deep.Equal(gotGtss, tt.wantGtss); len(pErr) > 0 {
				t.Errorf("ParseGTSArrayFromString() = %+v, want %+v\n%v", gotGtss, tt.wantGtss, pErr)
			}
		})
	}
}

func Test_parseSensisionLine(t *testing.T) {
	type args struct {
		in string
	}
	tests := []struct {
		name     string
		args     args
		wantTs   int64
		wantLat  float64
		wantLong float64
		wantAlt  float64
		wantC    string
		wantL    map[string]string
		wantA    map[string]string
		wantV    interface{}
		wantErr  bool
	}{
	// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotTs, gotLat, gotLong, gotAlt, gotC, gotL, gotA, gotV, err := parseSensisionLine(tt.args.in)
			if (err != nil) != tt.wantErr {
				t.Errorf("parseSensisionLine() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if gotTs != tt.wantTs {
				t.Errorf("parseSensisionLine() gotTs = %v, want %v", gotTs, tt.wantTs)
			}
			if gotLat != tt.wantLat {
				t.Errorf("parseSensisionLine() gotLat = %v, want %v", gotLat, tt.wantLat)
			}
			if gotLong != tt.wantLong {
				t.Errorf("parseSensisionLine() gotLong = %v, want %v", gotLong, tt.wantLong)
			}
			if gotAlt != tt.wantAlt {
				t.Errorf("parseSensisionLine() gotAlt = %v, want %v", gotAlt, tt.wantAlt)
			}
			if gotC != tt.wantC {
				t.Errorf("parseSensisionLine() gotC = %v, want %v", gotC, tt.wantC)
			}
			if !reflect.DeepEqual(gotL, tt.wantL) {
				t.Errorf("parseSensisionLine() gotL = %v, want %v", gotL, tt.wantL)
			}
			if !reflect.DeepEqual(gotA, tt.wantA) {
				t.Errorf("parseSensisionLine() gotA = %v, want %v", gotA, tt.wantA)
			}
			if !reflect.DeepEqual(gotV, tt.wantV) {
				t.Errorf("parseSensisionLine() gotV = %v, want %v", gotV, tt.wantV)
			}
		})
	}
}

func TestGTS_Sensision(t *testing.T) {
	type fields struct {
		ClassName    string
		Labels       map[string]string
		Attributes   map[string]string
		LastActivity int64
		Values       [][]interface{}
	}
	tests := []struct {
		name   string
		fields fields
		wantS  string
	}{{
		name: "Format metric",
		fields: fields{
			ClassName:    "my.metric",
			Labels:       map[string]string{},
			Attributes:   map[string]string{},
			LastActivity: 0,
			Values:       [][]interface{}{{1234, 0, 0, 0, 100}},
		},
		wantS: "1234// my.metric{} 100\n",
	}}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gts := &GTS{
				ClassName:    tt.fields.ClassName,
				Labels:       tt.fields.Labels,
				Attributes:   tt.fields.Attributes,
				LastActivity: tt.fields.LastActivity,
				Values:       tt.fields.Values,
			}
			if gotS := gts.Sensision(); gotS != tt.wantS {
				t.Errorf("GTS.Sensision() = %v, want %v", gotS, tt.wantS)
			}
		})
	}
}
