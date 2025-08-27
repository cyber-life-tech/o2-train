package drift

// Input of a drift detection operation
type Input struct {
	RefVector []float64 `json:"ref_vector"`
	NewVector []float64 `json:"new_vector"`
	Metric    Metric    `json:"metric"`
	Threshold float64   `json:"threshold"`
	PSI       *InputPSI `json:"psi"`
}

// InputPSI is an additional configuration for a drift detection operation,
// used to calculate the PopulationStabilityIndex
type InputPSI struct {
	Buckets       int     `json:"buckets"`
	MinPercentage float64 `json:"min_percentage"`
}

// DefaultInputPSI is a default constructor used for proof-of-concept or testing purposes
func DefaultInputPSI() *InputPSI {
	const (
		defaultBuckets = 10
		minPercentage  = 1e-4
	)

	return &InputPSI{
		Buckets:       defaultBuckets,
		MinPercentage: minPercentage,
	}
}

// Output of a drift detection operation, including a detailed value for the magnitude of data drift
type Output struct {
	Detected bool    `json:"detected"`
	Value    float64 `json:"value"`
}
