package resolv

import (
	"math"
)

// Resolve attempts to move the checking Shape with the specified X and Y values, returning a Collision object
// if it collides with the specified other Shape. The deltaX and deltaY arguments are the movement displacement
// in pixels. For platformers in particular, you would probably want to resolve on the X and Y axes separately.
func Resolve(firstShape Shape, other Shape, deltaX, deltaY int32) Collision {

	out := Collision{}
	out.ResolveX = deltaX
	out.ResolveY = deltaY
	out.ShapeA = firstShape

	if deltaX == 0 && deltaY == 0 {
		return out
	}

	x := float32(deltaX)
	y := float32(deltaY)

	primeX := true
	slope := float32(0)

	if math.Abs(float64(deltaY)) > math.Abs(float64(deltaX)) {
		primeX = false
		if deltaY != 0 && deltaX != 0 {
			slope = float32(deltaX) / float32(deltaY)
		}
	} else if deltaY != 0 && deltaX != 0 {
		slope = float32(deltaY) / float32(deltaX)
	}

	for true {

		if firstShape.WouldBeColliding(other, out.ResolveX, out.ResolveY) {

			if primeX {

				if deltaX > 0 {
					x--
				} else if deltaX < 0 {
					x++
				}

				if deltaY > 0 {
					y -= slope
				} else if deltaY < 0 {
					y += slope
				}

			} else {

				if deltaY > 0 {
					y--
				} else if deltaY < 0 {
					y++
				}

				if deltaX > 0 {
					x -= slope
				} else if deltaX < 0 {
					x += slope
				}

			}

			out.ResolveX = int32(x)
			out.ResolveY = int32(y)
			out.ShapeB = other

		} else {
			break
		}

	}

	if math.Abs(float64(deltaX-out.ResolveX)) > math.Abs(float64(deltaX)*1.5) || math.Abs(float64(deltaY-out.ResolveY)) > math.Abs(float64(deltaY)*1.5) {
		out.Teleporting = true
	}

	return out

}

// Distance returns the distance from one pair of X and Y values to another.
func Distance(x, y, x2, y2 int32) int32 {

	dx := x - x2
	dy := y - y2
	ds := (dx * dx) + (dy * dy)
	return int32(math.Sqrt(math.Abs(float64(ds))))

}
