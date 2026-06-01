package dice

import (
	"math"
	"sort"
)

// Stats summarizes the per-roll totals of a multi-roll batch. All
// "numeric" fields are computed once from the same input; callers
// (formatters, history viewers, embedded consumers) should construct
// this via ComputeStats or StatsFromResults rather than recomputing
// individual fields ad-hoc.
//
// Variance/StdDev use the POPULATION formula (divide by N), not the
// sample formula (N-1), since the dice batch IS the full population
// being summarized.
type Stats struct {
	Count   int
	Sum     int
	Min     int
	Max     int
	Average float64
	Median  float64
	StdDev  float64
}

// ComputeStats reduces a slice of integer totals to a Stats summary.
// Returns the zero value for an empty input.
func ComputeStats(totals []int) Stats {
	if len(totals) == 0 {
		return Stats{}
	}

	minVal := totals[0]
	maxVal := totals[0]
	sum := 0
	for _, v := range totals {
		sum += v
		if v < minVal {
			minVal = v
		}
		if v > maxVal {
			maxVal = v
		}
	}

	avg := float64(sum) / float64(len(totals))

	return Stats{
		Count:   len(totals),
		Sum:     sum,
		Min:     minVal,
		Max:     maxVal,
		Average: avg,
		Median:  median(totals),
		StdDev:  stdDev(totals, avg),
	}
}

// StatsFromResults is a convenience for the common case of summarizing
// MultiRollResult.Rolls. Equivalent to ComputeStats called with each
// result's Total.
func StatsFromResults(rolls []Result) Stats {
	totals := make([]int, len(rolls))
	for i, r := range rolls {
		totals[i] = r.Total
	}
	return ComputeStats(totals)
}

func median(values []int) float64 {
	sorted := append([]int(nil), values...)
	sort.Ints(sorted)
	mid := len(sorted) / 2
	if len(sorted)%2 == 1 {
		return float64(sorted[mid])
	}
	return float64(sorted[mid-1]+sorted[mid]) / 2
}

func stdDev(values []int, mean float64) float64 {
	if len(values) == 0 {
		return 0
	}
	var variance float64
	for _, v := range values {
		delta := float64(v) - mean
		variance += delta * delta
	}
	variance /= float64(len(values))
	return math.Sqrt(variance)
}
