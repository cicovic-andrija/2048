package core

import "math/rand"

func min(a int, b int) int {
	if a < b {
		return a
	}
	return b
}

func max(a int, b int) int {
	if a > b {
		return a
	}
	return b
}

// assumes p is in [0.0, 1.0]
func ptrue(p float64, rng *rand.Rand) bool {
	if rng.Float64() < p {
		return true
	}
	return false
}
