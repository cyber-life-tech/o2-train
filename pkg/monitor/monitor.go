package monitor

import "context"

// Monitor ticks at a specified interval and executes the necessary logical steps
type Monitor interface {
	Tick(ctx context.Context) error
}
