package base

import (
	"fmt"
	"strings"
)

// OneOf craft a label value where a metrics must have one of the given labels
// Example:
//  l :=Labels{}
//  l.OneOf("datacenter", []string{"BHS3", "GRA1"})
//  // l = {datacenter~(BHS3|GRA1)}
// Or:
//  l :=Labels.OneOf("datacenter", []string{"BHS3", "GRA1"})
func (l Labels) OneOf(name string, values []string) Labels {
	if l == nil {
		l = Labels{}
	}
	if name == "" {
		return l
	}
	if len(values) == 0 {
		return l
	}

	l[name] = fmt.Sprintf("~(%s)", strings.Join(values, "|"))
	return l
}
