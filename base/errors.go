package base

import "errors"

var (
	NoReadTokenError  = errors.New("this Warp10 call need a READ token access on the data")
	NoWriteTokenError = errors.New("this Warp10 call need a WRITE token access on the data")
)
