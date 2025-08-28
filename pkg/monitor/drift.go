package monitor

import (
	"context"
	"encoding/json"
	"fmt"
	"io"

	"golang.org/x/sync/errgroup"

	drift "github.com/cyber-life-tech/o2-train/internal/data-drift"
)

// Enforce interface implementation at compile time
var _ Monitor = (*DataDrift)(nil)

// DataDrift is a monitor designed to detect data drift from any source,
// execute the corresponding logic, and provide a detailed detection response
type DataDrift struct {
	reader io.Reader
	writer io.Writer
	hooks  []Hook[drift.Output]
}

// NewDataDrift is a constructor for the DataDrift
func NewDataDrift(r io.Reader, w io.Writer, hooks ...Hook[drift.Output]) DataDrift {
	return DataDrift{
		reader: r,
		writer: w,
		hooks:  hooks,
	}
}

// Tick implementation for the Monitor
func (dd DataDrift) Tick(ctx context.Context) error {
	if dd.reader == nil {
		return fmt.Errorf("%w: reader is required", ErrTick)
	}

	data, err := io.ReadAll(dd.reader)
	if err != nil {
		return fmt.Errorf("failed to read the data: %w", err)
	}

	var (
		in  drift.Input
		out drift.Output
	)

	if err = json.Unmarshal(data, &in); err != nil {
		return fmt.Errorf("failed to json.Unmarshal data into drift.Input: %w", err)
	}

	out, err = drift.Detect(in)
	if err != nil {
		return fmt.Errorf("failed to detect data drift: %w", err)
	}

	if out.Detected {
		return dd.onDriftDetected(ctx, out)
	}

	return nil
}

func (dd DataDrift) onDriftDetected(ctx context.Context, out drift.Output) error {
	if len(dd.hooks) > 0 {
		eg, egCTX := errgroup.WithContext(ctx)

		for _, hook := range dd.hooks {
			eg.Go(func() error {
				if err := hook(egCTX, out); err != nil {
					return fmt.Errorf("data drift hook failed: %w", err)
				}

				return nil
			})
		}

		if err := eg.Wait(); err != nil {
			return fmt.Errorf("failed to execute all hooks: %w", err)
		}
	}

	if dd.writer != nil {
		response, err := json.Marshal(out)
		if err != nil {
			return fmt.Errorf("failed to json.Marshal the drift.Output: %w", err)
		}

		if _, err = dd.writer.Write(response); err != nil {
			return fmt.Errorf("failed to write the data: %w", err)
		}
	}

	return nil
}
