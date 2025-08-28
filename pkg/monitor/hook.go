package monitor

import (
	"context"
)

// Hook is a generic processor that can be used in any monitor
type Hook[T any] func(ctx context.Context, entity T) error
