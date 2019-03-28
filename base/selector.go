package base

// Selector is a Warp10 selector
type Selector string

const (
	// WildCardSelector select all available GTS
	WildCardSelector Selector = "~.*{}"
)

// NewSelector Build a new Selector
func NewSelector(className string, labels Labels) Selector {
	return Selector(className + formatLabels(labels))
}
