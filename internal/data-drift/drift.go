package data_drift

import (
	"math"
	"slices"
	"sort"
)

const (
	KolmogorovSmirnov        DetectionMetric = "ks"
	PopulationStabilityIndex DetectionMetric = "psi"
)

type DetectionMetric string

func Detect(referenceData, newData []float64, metric DetectionMetric, threshold float64, buckets int) Result {
	allPoints := append(referenceData, newData...)

	for _, v := range allPoints {
		if math.IsNaN(v) || math.IsInf(v, 0) {
			panic("Input data contains NaN or Inf") // Or return custom error
		}
	}

	switch metric {
	case KolmogorovSmirnov:
		if len(referenceData) == 0 && len(newData) == 0 {
			return Result{Detected: false, Value: 0}
		}

		if len(referenceData) == 0 || len(newData) == 0 {
			return Result{Detected: true, Value: 1} // Assume drift if one is empty
		}

		// Compute KS statistic
		sort.Float64s(allPoints)
		uniquePoints := make([]float64, 0, len(allPoints))
		for i := range allPoints {
			if i == 0 || allPoints[i] != allPoints[i-1] {
				uniquePoints = append(uniquePoints, allPoints[i])
			}
		}

		// Sort copies for ECDF
		refSorted := make([]float64, len(referenceData))
		copy(refSorted, referenceData)
		sort.Float64s(refSorted)
		newSorted := make([]float64, len(newData))
		copy(newSorted, newData)
		sort.Float64s(newSorted)

		// ECDF functions (proportion <= x)
		var maxDiff float64
		for _, x := range uniquePoints {
			refEcdf := upperBound(refSorted, x) / float64(len(refSorted))
			newEcdf := upperBound(newSorted, x) / float64(len(newSorted))
			diff := math.Abs(refEcdf - newEcdf)
			if diff > maxDiff {
				maxDiff = diff
			}
		}
		return Result{Detected: maxDiff > threshold, Value: maxDiff}
	case PopulationStabilityIndex:
		if len(referenceData) == 0 && len(newData) == 0 {
			return Result{Detected: false, Value: 0}
		}

		if len(referenceData) == 0 || len(newData) == 0 {
			return Result{Detected: true, Value: 1} // Assume drift if one is empty
		}

		// Custom PSI implementation
		minVal, maxVal := slices.Min(referenceData), slices.Max(referenceData)
		if minVal == maxVal {
			maxVal += 1e-10 // Avoid zero-width bins
		}

		breakpoints := make([]float64, buckets+1)
		for i := 0; i <= buckets; i++ {
			breakpoints[i] = minVal + (maxVal-minVal)*float64(i)/float64(buckets)
		}

		refPercents := histogramPercentages(referenceData, breakpoints)
		newPercents := histogramPercentages(newData, breakpoints)

		var metricValue float64
		for i := range refPercents {
			ePerc := refPercents[i]
			aPerc := newPercents[i]
			if aPerc == 0 {
				aPerc = 0.0001
			}
			if ePerc == 0 {
				ePerc = 0.0001
			}
			metricValue += (ePerc - aPerc) * math.Log(ePerc/aPerc)
		}
		return Result{Detected: metricValue > threshold, Value: metricValue}
	default:
		panic("Metric must be 'ks' or 'psi'")
	}
}
