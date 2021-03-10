package resolv

import (
	"math"
)

// Rectangle represents a rectangle.
type Rectangle struct {
	BasicShape
	W, H float64
}

// NewRectangle creates a new Rectangle and returns a pointer to it.
func NewRectangle(x, y, w, h float64) *Rectangle {
	r := &Rectangle{W: w, H: h}
	r.tags = NewTags()
	r.X = x
	r.Y = y
	return r
}

// IsColliding returns whether the Rectangle is colliding with the specified other Shape or not, including the other Shape
// being wholly contained within the Rectangle.
func (r *Rectangle) IsColliding(other Shape) bool {
	return r.Check(other, 0, 0).Colliding()
}

// Check returns a MovementCheck indicating how far "in" the Rectangle would be into the other Shape if the movement was made.
func (r *Rectangle) Check(other Shape, dx, dy float64) *MovementCheck {

	movement := newMovementCheck(r, other)
	movement.Dx, movement.Dy = dx, dy

	switch other := other.(type) {

	case *Rectangle:

		x, y := math.Max(r.X, other.X), math.Max(r.Y, other.Y)
		x2, y2 := math.Min(r.X+r.W, other.X+other.W), math.Min(r.Y+r.H, other.Y+other.H)
		w, h := x2-x, y2-y

		if w < 0 {
			w = 0
		}

		if h < 0 {
			h = 0
		}

		if dx < 0 {
			movement.Dx = w
		} else if dx > 0 {
			movement.Dx = -w
		}

		if dy < 0 {
			movement.Dy = h
		} else if dy > 0 {
			movement.Dy = -h
		}

		if w > 0 && h > 0 {

			movement.addPoint(x, y)
			movement.addPoint(x+w, y)
			movement.addPoint(x+w, y+h)
			movement.addPoint(x, y+h)

		}

	case *Circle:

		

	case *Line:

		highestDx, highestDy := 0.0, 0.0

		for _, side := range r.ToLines() {

			sideCollision := other.Check(side, dx, dy)
			movement.Points = append(movement.Points, sideCollision.Points...)

			if sideCollision.Colliding() {

				// We want to use the collision's returned delta values, but only if it's higher than any other line. Imagine a situation where a rectangle is
				// moving up a slope. The left corner of the rectangle could return a lower or higher delta sliding value than the right one - we want to use
				// whichever is higher, since that would indicate a "stronger" movement.

				if math.Abs(sideCollision.Dx) > math.Abs(highestDx) {
					movement.Dx = sideCollision.Dx
					highestDx = movement.Dx
				}

				if math.Abs(sideCollision.Dy) > math.Abs(highestDy) {
					movement.Dy = sideCollision.Dy
					highestDy = movement.Dy
				}

			}

		}

	case *Space:

		for _, shape := range *other {
			if shape == r {
				continue
			}
			if test := r.Check(shape, dx, dy); test.Colliding() {
				movement = test
				break
			}
		}

	}

	return movement

}

func (r *Rectangle) Valid() bool {
	return r.W > 0 && r.H > 0
}

func (r *Rectangle) ToLines() []*Line {
	return []*Line{
		NewLine(r.X, r.Y, r.X+r.W, r.Y),
		NewLine(r.X+r.W, r.Y, r.X+r.W, r.Y+r.H),
		NewLine(r.X+r.W, r.Y+r.H, r.X, r.Y+r.H),
		NewLine(r.X, r.Y+r.H, r.X, r.Y),
	}
}

// Center returns the center point of the Rectangle.
func (r *Rectangle) Center() (float64, float64) {

	x := r.X + r.W/2
	y := r.Y + r.H/2

	return x, y

}

// GetBoundingCircle returns a circle that wholly contains the Rectangle.
// func (r *Rectangle) GetBoundingCircle() *Circle {

// 	x, y := r.Center()
// 	c := NewCircle(x, y, Distance(x, y, r.X+r.W, r.Y))
// 	return c

// }
