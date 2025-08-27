package drift

import (
	"fmt"
	"math"
	"slices"
)

func detectPSI(input Input) (Output, error) {
	if input.PSI == nil {
		return Output{}, fmt.Errorf("no PSI configuration in the given input")
	}

	minVal, maxVal := slices.Min(input.RefVector), slices.Max(input.RefVector)
	if minVal == maxVal {
		maxVal += 1e-10 // Avoid zero-width buckets
	}

	breakpoints := make([]float64, input.PSI.Buckets+1)
	for i := 0; i <= input.PSI.Buckets; i++ {
		breakpoints[i] = minVal + (maxVal-minVal)*float64(i)/float64(input.PSI.Buckets)
	}

	var (
		refPercents = histogramPercentages(input.RefVector, breakpoints)
		newPercents = histogramPercentages(input.NewVector, breakpoints)
		metricValue = 0.0
	)

	for i := range refPercents {
		var (
			refPerc = max(refPercents[i], input.PSI.MinPercentage)
			newPerc = max(newPercents[i], input.PSI.MinPercentage)
		)

		metricValue += (refPerc - newPerc) * math.Log(refPerc/newPerc)
	}

	return Output{Detected: metricValue > input.Threshold, Value: metricValue}, nil
}

func histogramPercentages(data, breakpoints []float64) []float64 {
	var (
		nIn, nOut = len(data), len(breakpoints) - 1
		counts    = make([]int, nOut)
		percents  = make([]float64, nOut)
	)

	if nIn == 0 {
		return percents
	}

	for _, v := range data {
		placed := false

		for i := 0; i < nOut; i++ {
			if (i == 0 && v <= breakpoints[i+1]) || (v > breakpoints[i] && v <= breakpoints[i+1]) {
				counts[i]++
				placed = true

				break
			}
		}

		if !placed {
			counts[nOut-1]++ // Catch v > last breakpoint
		}
	}

	for i, c := range counts {
		percents[i] = float64(c) / float64(nIn)
	}

	return percents
}
