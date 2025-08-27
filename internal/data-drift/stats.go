package data_drift

import (
	"sort"
)

// The slice must be sorted in ascending order.
func upperBound(sorted []float64, x float64) float64 {
	return float64(sort.Search(len(sorted), func(i int) bool { return sorted[i] > x }))
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
