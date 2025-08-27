package drift

import (
	"fmt"
	"math"
)

// Detect data drift between two data vectors, typically training(ref) and test/production(new) vectors.
// Choose a Metric that suits your specific dataset.
func Detect(input Input) (Output, error) {
	if len(input.RefVector) == 0 && len(input.NewVector) == 0 {
		return Output{}, nil
	}

	if len(input.RefVector) == 0 || len(input.NewVector) == 0 {
		return Output{Detected: true, MetricValue: 1}, nil // Assume drift if one is empty
	}

	if err := checkInvalidValues(input.RefVector); err != nil {
		return Output{}, fmt.Errorf("invalid reference vector: %w", err)
	}

	if err := checkInvalidValues(input.NewVector); err != nil {
		return Output{}, fmt.Errorf("invalid new vector: %w", err)
	}

	switch input.Metric {
	case KolmogorovSmirnov:
		return detectKS(input), nil
	case PopulationStabilityIndex:
		return detectPSI(input)
	default:
		return Output{}, fmt.Errorf(
			"'%s' is an %w; it must be one of: %s, %s",
			input.Metric, ErrInvalidMetric, KolmogorovSmirnov, PopulationStabilityIndex,
		)
	}
}

func checkInvalidValues(in []float64) error {
	for i, v := range in {
		if math.IsNaN(v) || math.IsInf(v, 0) {
			return fmt.Errorf("%w, found at %d", ErrInvalidValue, i) // Or return custom error
		}
	}

	return nil
}
