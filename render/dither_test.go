package render

import (
	"math"
	"testing"

	"github.com/thomaso-mirodin/intmath/intgr"
)

// Brute force test to make sure that over a large set of inputs, we obey the
// constraint that the total we get out is the same as what we put in.
func TestDitherBruteForce(t *testing.T) {
	for amount := 1; amount < 250; amount++ {
		for buckets := 1; buckets <= amount; buckets++ {
			generator := dither(amount, buckets)
			sum := 0
			firstHalfSum := 0
			secondHalfSum := 0
			max := 0
			min := math.MaxInt16 // just needs to be big
			for i := 0; i < buckets; i++ {
				val := generator()
				sum += val
				if i < buckets/2 {
					firstHalfSum += val
				} else {
					secondHalfSum += val
				}

				max = intgr.Max(max, val)
				min = intgr.Min(min, val)
			}

			if sum != amount {
				t.Errorf("dithert(%d, %d) summed to %d", amount, buckets, sum)
			}

			if max-min > 1 {
				t.Errorf(
					"dither(%d, %d) should have difference of no more than one. Min: %d, max: %d",
					amount, buckets, min, max)
			}

			if amount%buckets == 0 && min != max {
				t.Errorf(
					"dither(%d, %d) should have equal sized buckets. Min: %d, max: %d",
					amount, buckets, min, max)
			}

			// This is a rough assertion that we're dividing space up evenly
			if buckets%2 == 0 && amount%2 == 0 && firstHalfSum != secondHalfSum {
				t.Errorf(
					"dither(%d, %d) should have equal first and second half sums (even amount and buckets)\n"+
						"but got %d for first half and %d for second half",
					amount, buckets, firstHalfSum, secondHalfSum)
			} else if buckets%2 == 0 && intgr.Abs(firstHalfSum-secondHalfSum) > 1 {
				t.Errorf("dither(%d, %d) first and second half should differ by no more than 1 (even buckets)\n"+
					"but got %d for first half and %d for second half",
					amount, buckets, firstHalfSum, secondHalfSum)
			}
		}
	}
}
