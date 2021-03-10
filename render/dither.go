package render

import (
	"math"
)

// There are several cases where we need to allocate width across an uneven
// number of buckets. Naive approaches tend to clump the error unnaturally. This
// is a general purpose function for solving this problem in a way that is both
// deterministic and aesthetically pleasing.
//
// Given an amount of space to divide up and a number of buckets N, produces a
// generator that can be called N times to get a nicely divided allocation.

func dither(amount, buckets int) func() int {
	// The amount we'd like to allocate if we had floats
	idealPerBucket := float64(amount) / float64(buckets)
	// How far off of our target last time
	lastError := 0.0

	return func() int {
		// Diffuse the accumulated error into this allocation
		correctedAllocation := idealPerBucket - lastError
		// Snap to integer
		actualAllocation := int(math.Round(correctedAllocation))
		// Propagate the error forward
		lastError = float64(actualAllocation) - correctedAllocation
		return actualAllocation
	}
}
