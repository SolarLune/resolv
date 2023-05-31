package resolv

import (
	"math"

	"github.com/quartercastle/vector"
)

// ToRadians is a helper function to easily convert degrees to radians.
func ToRadians(degrees float64) float64 {
	return math.Pi * degrees / 180
}

// ToDegrees is a helper function to easily convert radians to degrees for human readability.
func ToDegrees(radians float64) float64 {
	return radians / math.Pi * 180
}

func dot(a, b vector.Vector) float64 {
	result := a[0]*b[0] + a[1]*b[1]
	return result
}
