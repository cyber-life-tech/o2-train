package drift

import (
	"math"
	"math/rand"
	"testing"

	"github.com/brianvoe/gofakeit/v7"
	"github.com/stretchr/testify/assert"
)

func TestDetectDataDrift(t *testing.T) {
	r := rand.New(rand.NewSource(42))

	type testCase struct {
		name           string
		input          Input
		expectedResult Output
		expectedError  error
	}

	happy := []testCase{
		{
			name: "Happy - KS - No Drift",
			input: Input{
				RefVector: generateNormalData(t, 1000, 0, 1, r),
				NewVector: generateNormalData(t, 1000, 0, 1, r),
				Metric:    KolmogorovSmirnov,
				Threshold: 0.1,
			},
			expectedResult: Output{Detected: false, MetricValue: 0.036},
		},
		{
			name: "Happy - KS - Drift",
			input: Input{
				RefVector: generateNormalData(t, 1000, 0, 1, r),
				NewVector: generateNormalData(t, 1000, 0.5, 1, r),
				Metric:    KolmogorovSmirnov,
				Threshold: 0.1,
			},
			expectedResult: Output{Detected: true, MetricValue: 0.228},
		},
		{
			name: "Happy - PSI - No Drift",
			input: Input{
				RefVector: generateNormalData(t, 1000, 0, 1, r),
				NewVector: generateNormalData(t, 1000, 0, 1, r),
				Metric:    PopulationStabilityIndex,
				Threshold: 0.25,
				PSI:       DefaultInputPSI(),
			},
			expectedResult: Output{Detected: false, MetricValue: 0.012},
		},
		{
			name: "Happy - PSI - Drift",
			input: Input{
				RefVector: generateNormalData(t, 1000, 0, 1, r),
				NewVector: generateNormalData(t, 1000, 0.5, 1, r),
				Metric:    PopulationStabilityIndex,
				Threshold: 0.311,
				PSI:       DefaultInputPSI(),
			},
			expectedResult: Output{Detected: true, MetricValue: 0.311},
		},
	}
	empty := []testCase{
		{
			name: "Empty - Both - No Drift",
			input: Input{
				RefVector: []float64{},
				NewVector: []float64{},
			},
			expectedResult: Output{Detected: false, MetricValue: 0},
		},
		{
			name: "Empty - Reference - Drift",
			input: Input{
				RefVector: []float64{},
				NewVector: []float64{1, 2, 3},
			},
			expectedResult: Output{Detected: true, MetricValue: 1},
		},
		{
			name: "Empty - New - Drift",
			input: Input{
				RefVector: []float64{1, 2, 3},
				NewVector: []float64{},
			},
			expectedResult: Output{Detected: true, MetricValue: 1},
		},
	}
	edge := []testCase{
		{
			name: "All Values Equal/MinMax Epsilon Adjustment - PSI - No Drift",
			input: Input{
				RefVector: []float64{5, 5, 5, 5},
				NewVector: []float64{5, 5, 5, 5},
				Metric:    PopulationStabilityIndex,
				Threshold: 0.25,
				PSI:       DefaultInputPSI(),
			},
			expectedResult: Output{Detected: false, MetricValue: 0},
		},
		{
			name: "NewVector is out of RefVector's range - PSI - Drift",
			input: Input{
				RefVector: []float64{1, 2, 3},
				NewVector: []float64{4, 5, 6},
				Metric:    PopulationStabilityIndex,
				Threshold: 0.25,
				PSI:       DefaultInputPSI(),
			},
			expectedResult: Output{Detected: true, MetricValue: 6.138},
		},
		{
			name: "Duplicates in Vectors - KS - No Drift",
			input: Input{
				RefVector: []float64{1, 1, 2, 2, 3, 3},
				NewVector: []float64{1, 1, 2, 2, 3, 3},
				Metric:    KolmogorovSmirnov,
				Threshold: 0.1,
			},
			expectedResult: Output{Detected: false, MetricValue: 0},
		},
		{
			name: "Short Vectors - KS - No Drift",
			input: Input{
				RefVector: []float64{3},
				NewVector: []float64{3},
				Metric:    KolmogorovSmirnov,
				Threshold: 0.1,
			},
			expectedResult: Output{Detected: false, MetricValue: 0},
		},
		{
			name: "Negative Values - PSI - No Drift",
			input: Input{
				RefVector: []float64{-5, -4, -3},
				NewVector: []float64{-5, -4, -3},
				Metric:    PopulationStabilityIndex,
				Threshold: 0.25,
				PSI:       DefaultInputPSI(),
			},
			expectedResult: Output{Detected: false, MetricValue: 0},
		},
		{
			name: "Large Values - KS - Drift",
			input: Input{
				RefVector: []float64{1e10, 2e10},
				NewVector: []float64{3e10, 4e10},
				Metric:    KolmogorovSmirnov,
				Threshold: 0.1,
			},
			expectedResult: Output{Detected: true, MetricValue: 1},
		},
	}
	faulty := []testCase{
		{
			name: "Faulty - Reference - ErrInvalidValue",
			input: Input{
				RefVector: []float64{math.NaN(), math.Inf(-1), math.Inf(1)},
				NewVector: []float64{1, 2, 3},
			},
			expectedError: ErrInvalidValue,
		},
		{
			name: "Faulty - New - ErrInvalidValue",
			input: Input{
				RefVector: []float64{1, 2, 3},
				NewVector: []float64{math.NaN(), math.Inf(-1), math.Inf(1)},
			},
			expectedError: ErrInvalidValue,
		},
		{
			name: "Faulty - ErrInvalidMetric",
			input: Input{
				RefVector: []float64{1, 2, 3},
				NewVector: []float64{1, 2, 3},
				Metric:    Metric(gofakeit.HackerNoun()),
			},
			expectedError: ErrInvalidMetric,
		},
		{
			name: "Faulty - PSI - ErrEmptyExtraConfig",
			input: Input{
				RefVector: []float64{1, 2, 3},
				NewVector: []float64{1, 2, 3},
				Metric:    PopulationStabilityIndex,
				Threshold: 0.25,
				PSI:       nil,
			},
			expectedError: ErrEmptyExtraConfig,
		},
	}

	for _, tests := range [][]testCase{happy, empty, edge, faulty} {
		for _, tt := range tests {
			t.Run(tt.name, func(t *testing.T) {
				t.Parallel()

				if tt.expectedError != nil {
					_, err := Detect(tt.input)
					assert.ErrorIs(t, err, tt.expectedError)

					return
				}

				result, err := Detect(tt.input)

				assert.NoError(t, err)
				assert.Equal(t, tt.expectedResult.Detected, result.Detected)
				assert.InDelta(t, tt.expectedResult.MetricValue, result.MetricValue, 1e-3)
			})
		}
	}
}

func generateNormalData(t *testing.T, n int, mean, std float64, r *rand.Rand) []float64 {
	t.Helper()

	data := make([]float64, n)
	for i := range data {
		data[i] = mean + std*r.NormFloat64()
	}

	return data
}
