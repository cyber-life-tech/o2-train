package drift

import "errors"

var (
	ErrInvalidValue  = errors.New("a value cannot be NaN or Inf")
	ErrInvalidMetric = errors.New("invalid metric")
)
