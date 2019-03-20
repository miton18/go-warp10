package instrumentation

import (
	b "github.com/miton18/go-warp10/base"
)

type State struct {
	name, help string
	text       string
	context    b.Labels
}

func NewState(name string, ctx b.Labels, help string) *State {
	if ctx == nil {
		ctx = b.Labels{}
	}

	return &State{
		name:    name,
		help:    help,
		context: ctx,
		text:    "",
	}
}

// Name return the metric name of the annotation
func (s *State) Name() string {
	return s.name
}

// Help return informations about this metric
func (s *State) Help() string {
	return "state: " + s.help
}

// Get return a plain GTS
func (s *State) Get() b.GTSList {
	return b.GTSList{&b.GTS{
		ClassName: s.name,
		Labels:    s.context,
		Values:    [][]interface{}{{s.text}},
	}}
}

// Set a new state value
func (s *State) Set(newState string) {
	s.text = newState
}

// Reset set to "" the state
func (s *State) Reset() {
	s.text = ""
}
