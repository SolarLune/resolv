package resolv

import (
	"math"
)

// Distance returns the distance from one pair of X and Y values to another.
func Distance(x, y, x2, y2 float64) float64 {

	dx := x - x2
	dy := y - y2
	ds := (dx * dx) + (dy * dy)
	return float64(math.Sqrt(math.Abs(float64(ds))))

}
