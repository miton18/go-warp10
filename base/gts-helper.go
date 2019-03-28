package base

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
		Labels:       Labels{},
		Attributes:   Attributes{},
		LastActivity: 0,
		Values:       [][]interface{}{},
	}
}

// NewGTS return a nammed GTS
func NewGTS(className string) *GTS {
	return &GTS{
		ClassName:    className,
		Labels:       Labels{},
		Attributes:   Attributes{},
		LastActivity: 0,
		Values:       [][]interface{}{},
	}
}

// NewGTSWithLabels return a nammed and labelized GTS
func NewGTSWithLabels(className string, labels Labels) *GTS {
	return &GTS{
		ClassName: className,
		Labels:    labels,
	}
}

// ParseGTSFromString parse an unique sensision format line into a new GTS
func ParseGTSFromString(sensisionLine string) (*GTS, error) {
	ts, lat, long, alt, c, l, a, v, err := parseSensisionLine(strings.TrimSpace(sensisionLine))
	if err != nil {
		return nil, err
	}
	gts := &GTS{
		ClassName:  c,
		Labels:     l,
		Attributes: a,
		Values:     [][]interface{}{{ts, lat, long, alt, v}},
	}
	return gts, nil
}

// ParseGTSFromBytes parse sensision format into a new GTS
func ParseGTSFromBytes(in []byte) (*GTS, error) {
	return ParseGTSFromString(string(in))
}

// ParseGTSArrayFromString parse sensision format into a new GTS array
func ParseGTSArrayFromString(in string) (GTSList, error) {
	gtss := GTSList{}
	for _, sensisionLine := range strings.Split(in, "\n") {
		gts, err := ParseGTSFromString(sensisionLine)
		if err != nil {
			return nil, err
		}
		gtss = append(gtss, gts)
	}
	return gtss, nil
}

// ParseGTSArrayFromBytes parse sensision format into a new GTS array
func ParseGTSArrayFromBytes(in []byte) (GTSList, error) {
	return ParseGTSArrayFromString(string(in))
}

func parseSensisionLine(in string) (int64, float64, float64, float64, string, Labels, Attributes, interface{}, error) {
	var err error
	var v interface{}
	v = 0
	ts := int64(time.Now().Nanosecond()) / 1000
	lat := float64(0)
	long := float64(0)
	alt := float64(0)
	c := ""
	l := Labels{}
	a := Attributes{}

	// 469756475/5496:54965/5496757 test{l=v,ufgsdg=vfivuvf}{a=b} 57486847
	g := strings.Split(in, " ")
	if len(g) < 3 {
		err = errors.New("Cannot parse datapoint: " + in)
		return ts, lat, long, alt, c, l, a, v, err
	}

	ctx := strings.Split(g[0], "/")
	if len(ctx) != 3 {
		err = errors.New("Cannot parse datapoint: " + in)
		return ts, lat, long, alt, c, l, a, v, err
	}

	if ctx[0] == "" {
		err = errors.New("No timestamp provided")
		return ts, lat, long, alt, c, l, a, v, err
	}
	ts, err = strconv.ParseInt(ctx[0], 10, 64)
	if err != nil {
		err = errors.New("Cannot parse " + ctx[0] + " as valid timestamp")
		return ts, lat, long, alt, c, l, a, v, err
	}

	if ctx[2] != "" {
		alt, err = strconv.ParseFloat(ctx[2], 64)
		if err != nil {
			err = errors.New("Cannot parse " + ctx[0] + " as valid altitude")
			return ts, lat, long, alt, c, l, a, v, err
		}
	}

	if ctx[1] != "" {
		latlong := strings.Split(ctx[1], ":")
		if len(latlong) != 2 {
			err = errors.New("Cannot parse " + ctx[1] + " as valid latitude:longitude")
			return ts, lat, long, alt, c, l, a, v, err
		}

		lat, err = strconv.ParseFloat(latlong[0], 64)
		if err != nil {
			err = errors.New("Cannot parse " + ctx[0] + " as valid latitude")
			return ts, lat, long, alt, c, l, a, v, err
		}

		long, err = strconv.ParseFloat(latlong[1], 64)
		if err != nil {
			err = errors.New("Cannot parse " + ctx[0] + " as valid latitude")
			return ts, lat, long, alt, c, l, a, v, err
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
		return ts, lat, long, alt, c, l, a, v, err
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
				return ts, lat, long, alt, c, l, a, v, err
			}
			l[keyVal[0]] = keyVal[1]
		}
	}

	vStr := strings.Join(g[2:], " ")
	vStr = strings.Trim(vStr, "'")

	if len(vStr) == 0 {
		err = errors.New("Cannot parse " + strings.Join(g[2:], " ") + " as valid value")
		return ts, lat, long, alt, c, l, a, v, err
	}

	var vInt int
	vInt, err = strconv.Atoi(vStr)
	if err == nil {
		v = vInt
	} else if vInt, err := strconv.ParseInt(vStr, 10, 64); err == nil {
		v = vInt
	} else {
		v = vStr
	}

	return ts, lat, long, alt, c, l, a, v, nil
}

// SensisionSelector return the GTS selector (class + labels (+ attributes) )
func (gts *GTS) SensisionSelector(withAttributes bool) (string) {
	s := gts.ClassName + formatLabels(gts.Labels)
	if withAttributes {
		s += formatAttributes(gts.Attributes)
	}
	return s
}

// SensisionSelectors return the GTSList selectors (class + labels (+ attributes) )
func (gtsList GTSList) SensisionSelectors(withAttributes bool) (string) {
	s := ""
	for _, gts := range gtsList {
		s += gts.SensisionSelector(withAttributes) + "\n"
	}
	return s
}

// Sensision return the sensision format of the GTS
func (gts *GTS) Sensision() (string) {
	s := ""
	static := gts.SensisionSelector(false)

	for _, dps := range gts.Values {
		l := len(dps)
		if l == 0 {
			continue
		}

		ts := ""
		lat := ""
		lng := ""
		alt := ""
		val := getVal(dps[l-1])

		if l >= 2 {
			ts = fmt.Sprintf("%v", dps[0])
		}
		if l >= 4 {
			lat = fmt.Sprintf("%v", dps[1])
			lng = fmt.Sprintf("%v", dps[2])
		}
		if l == 5 {
			alt = fmt.Sprintf("%v", dps[3])
		}

		if lat != "" && lng != "" {
			s += fmt.Sprintf("%s/%s:%s/%s %s %s", ts, lat, lng, alt, static, val)
		} else {
			s += fmt.Sprintf("%s//%s %s %s", ts, alt, static, val)
		}

		s += "\n"
	}
	return s
}

// Sensision return the sensision format of all GTS
func (gtsList GTSList) Sensision() (string) {
	s := ""
	for _, gts := range gtsList {
		s += gts.Sensision() + "\n"
	}
	return s
}

func getVal(i interface{}) string {
	switch i.(type) {
	case int, int8, int16, int32, int64, uint, uint8, uint16, uint32, uint64:
		return fmt.Sprintf("%d", i)
	case float32, float64:
		return fmt.Sprintf("%g", i)
	case string:
		return fmt.Sprintf("'%s'", i)
	case fmt.Stringer:
		return fmt.Sprintf("'%s'", i.(fmt.Stringer).String())
	}
	return fmt.Sprintf("'%v'", i)
}

func formatLabels(labels Labels) (string) {
	s := "{"
	if labels != nil && len(labels) > 0 {
		pairs := []string{}
		for k, v := range labels {
			pairs = append(pairs, k+"="+v)
		}
		s += strings.Join(pairs, ",")
	}
	return s + "}"
}

func formatAttributes(attrs Attributes) (string) {
	s := "{"
	if attrs != nil {
		pairs := []string{}
		for k, v := range attrs {
			pairs = append(pairs, k+"="+v)
		}
		s += strings.Join(pairs, ",")
	}
	s += "}"
	return s
}
