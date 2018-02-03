package warp10

import (
	"errors"
	"fmt"
	"strconv"
	"strings"
	"time"
)

// NewEmptyGTS return an empty GTS
func NewEmptyGTS() *GTS {
	return &GTS{
		ClassName:    "",
		Labels:       map[string]string{},
		Attributes:   map[string]string{},
		LastActivity: 0,
		Values:       [][]interface{}{},
	}
}

// NewGTS return a nammed GTS
func NewGTS(className string) *GTS {
	return &GTS{
		ClassName:    className,
		Labels:       map[string]string{},
		Attributes:   map[string]string{},
		LastActivity: 0,
		Values:       [][]interface{}{},
	}
}

// NewGTSWithLabels return a nammed and labelized GTS
func NewGTSWithLabels(className string, labels map[string]string) *GTS {
	return &GTS{
		ClassName: className,
		Labels:    labels,
	}
}

// ParseGTSFromString parse an unique sensision format line into a new GTS
func ParseGTSFromString(sensisionLine string) (gts *GTS, err error) {
	ts, lat, long, alt, c, l, a, v, err := parseSensisionLine(strings.TrimSpace(sensisionLine))
	if err != nil {
		return
	}
	gts = &GTS{
		ClassName:  c,
		Labels:     l,
		Attributes: a,
		Values:     [][]interface{}{{ts, lat, long, alt, v}},
	}
	return
}

// ParseGTSFromBytes parse sensision format into a new GTS
func ParseGTSFromBytes(in []byte) (gts *GTS, err error) {
	return ParseGTSFromString(string(in))
}

// ParseGTSArrayFromString parse sensision format into a new GTS array
func ParseGTSArrayFromString(in string) (gtss []*GTS, err error) {
	gtss = []*GTS{}
	for _, sensisionLine := range strings.Split(in, "\n") {
		gts, err := ParseGTSFromString(sensisionLine)
		if err != nil {
			return gtss, err
		}
		gtss = append(gtss, gts)
	}
	return
}

// ParseGTSArrayFromBytes parse sensision format into a new GTS array
func ParseGTSArrayFromBytes(in []byte) (gtss []*GTS, err error) {
	return ParseGTSArrayFromString(string(in))
}

func parseSensisionLine(in string) (ts int64, lat float64, long float64, alt float64, c string, l map[string]string, a map[string]string, v interface{}, err error) {
	ts = int64(time.Now().Nanosecond()) / 1000
	lat = 0
	long = 0
	alt = 0
	c = ""
	l = map[string]string{}
	a = map[string]string{}
	v = 0
	err = nil

	// 469756475/5496:54965/5496757 test{l=v,ufgsdg=vfivuvf}{a=b} 57486847
	g := strings.Split(in, " ")
	if len(g) < 3 {
		err = errors.New("Cannot parse datapoint: " + in)
		return
	}

	ctx := strings.Split(g[0], "/")
	if len(ctx) != 3 {
		err = errors.New("Cannot parse datapoint: " + in)
		return
	}

	if ctx[0] == "" {
		err = errors.New("No timestamp provided")
		return
	}
	ts, err = strconv.ParseInt(ctx[0], 10, 64)
	if err != nil {
		err = errors.New("Cannot parse " + ctx[0] + " as valid timestamp")
		return
	}

	if ctx[2] != "" {
		alt, err = strconv.ParseFloat(ctx[2], 64)
		if err != nil {
			err = errors.New("Cannot parse " + ctx[0] + " as valid altitude")
			return
		}
	}

	if ctx[1] != "" {
		latlong := strings.Split(ctx[1], ":")
		if len(latlong) != 2 {
			err = errors.New("Cannot parse " + ctx[1] + " as valid latitude:longitude")
			return
		}

		lat, err = strconv.ParseFloat(latlong[0], 64)
		if err != nil {
			err = errors.New("Cannot parse " + ctx[0] + " as valid lattitude")
			return
		}

		long, err = strconv.ParseFloat(latlong[1], 64)
		if err != nil {
			err = errors.New("Cannot parse " + ctx[0] + " as valid lattitude")
			return
		}
	}

	// g[1] = test{l=v,ufgsdg=vfivuvf}{a=b}
	classLabelsAttributes := strings.SplitN(g[1], "{", 2)
	// test l=v,ufgsdg=vfivuvf}{a=b}
	// test l=v,ufgsdg=vfivuvf}
	// test }
	//fmt.Println(fmt.Sprintf("%+v", classLabelsAttributes))
	if len(classLabelsAttributes) != 2 {
		err = errors.New("Cannot parse " + g[1] + " as valid class + labels + attributes")
		return
	}

	c = classLabelsAttributes[0]
	classLabelsAttributes[1] = strings.TrimSuffix(classLabelsAttributes[1], "}")

	// contains labels
	if classLabelsAttributes[1] != "" {
		if strings.Contains(classLabelsAttributes[1], "}{") {
			//TODO: handle attributes
			classLabelsAttributes[1] = strings.Split(classLabelsAttributes[1], "}{")[0]
		}
		// attributes cleaned

		labelsValue := strings.Split(classLabelsAttributes[1], ",")
		for _, labelPair := range labelsValue {
			keyVal := strings.Split(labelPair, "=")
			if len(keyVal) != 2 {
				err = errors.New("Cannot parse " + labelPair + " as valid key and value label")
				return
			}
			l[keyVal[0]] = keyVal[1]
		}
	}

	vStr := strings.Join(g[2:], " ")
	vStr = strings.Trim(vStr, "'")

	if len(vStr) == 0 {
		err = errors.New("Cannot parse " + strings.Join(g[2:], " ") + " as valid value")
		return
	}

	if vInt, err := strconv.Atoi(vStr); err == nil {
		v = vInt
	} else if vInt, err := strconv.ParseInt(vStr, 10, 64); err == nil {
		v = vInt
	} else {
		v = vStr
	}

	return
}

// Sensision return the sensision format of the GTS
func (gts *GTS) Sensision() (s string) {
	s = ""
	static := " " + gts.ClassName + formatLabels(gts.Labels)
	if len(gts.Attributes) > 0 {
		static += formatLabels(gts.Attributes)
	}
	static += " "

	for _, dps := range gts.Values {
		dp := ""
		dp += fmt.Sprintf("%v", dps[0]) + "/"

		if dps[1] != 0 && dps[2] != 0 {
			lat, oka := dps[1].(float64)
			long, okb := dps[1].(float64)
			if !oka || !okb {
				continue
			}
			dp += fmt.Sprintf("%v:%v", lat, long)
		}
		dp += "/"

		if dps[3] != 0 {
			alt, ok := dps[3].(float64)
			if !ok {
				continue
			}
			dp += strconv.FormatFloat(alt, 'f', -1, 64)
		}

		s += dp + static + fmt.Sprintf("%v", dps[4]) + "\n"
	}

	return
}

func formatLabels(labels map[string]string) (s string) {
	s = "{"
	for k, v := range labels {
		s += k + "=" + v
	}
	s += "}"
	return
}
