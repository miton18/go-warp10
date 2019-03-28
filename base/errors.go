package base

import "errors"

var (
	NoReadTokenError  = errors.New("This Warp10 call need a READ token access on the data")
	NoWriteTokenError = errors.New("This Warp10 call need a WRITE token access on the data")
)