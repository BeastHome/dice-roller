package dice

//
// ------------------------------------------------------------
// NUMERIC HELPERS
// ------------------------------------------------------------
//

// min returns the smaller of two ints.
func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

// max returns the larger of two ints.
func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}

//
// ------------------------------------------------------------
// SLICE HELPERS
// ------------------------------------------------------------
//

// sum returns the sum of all values in a slice.
func sum(v []int) int {
	total := 0
	for _, x := range v {
		total += x
	}
	return total
}

// countIf counts how many values satisfy a predicate.
func countIf(v []int, fn func(int) bool) int {
	n := 0
	for _, x := range v {
		if fn(x) {
			n++
		}
	}
	return n
}
