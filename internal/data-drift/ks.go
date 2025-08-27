package drift

import (
	"math"
	"sort"

	mapset "github.com/deckarep/golang-set/v2"
)

func detectKS(input Input) Output {
	var (
		combinedVector       = append(input.RefVector, input.NewVector...)
		uniqueVector         = mapset.NewThreadUnsafeSet(combinedVector...).ToSlice()
		refSorted, newSorted = copySort(input.RefVector), copySort(input.NewVector)
		maxDiff              = 0.0
	)

	for _, x := range uniqueVector {
		var (
			refECDF = upperBound(refSorted, x) / float64(len(refSorted))
			newECDF = upperBound(newSorted, x) / float64(len(newSorted))
			diff    = math.Abs(refECDF - newECDF)
		)

		if diff > maxDiff {
			maxDiff = diff
		}
	}

	return Output{Detected: maxDiff > input.Threshold, MetricValue: maxDiff}
}

func upperBound(sorted []float64, x float64) float64 {
	return float64(sort.Search(len(sorted), func(i int) bool { return sorted[i] > x }))
}

func copySort(in []float64) []float64 {
	out := make([]float64, len(in))

	copy(out, in)
	sort.Float64s(out)

	return out
}
