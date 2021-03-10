package resolv

import (
	"fmt"
	"math"
	"sort"
)

// Line represents a line, from one point to another.
type Line struct {
	BasicShape
	X2, Y2 float64
}

// NewLine returns a  Line instance.
func NewLine(x, y, x2, y2 float64) *Line {
	l := &Line{}
	l.tags = NewTags()
	l.X = x
	l.Y = y
	l.X2 = x2
	l.Y2 = y2
	return l
}

// BUG(SolarLune): Line.IsColliding() and Line.GetIntersectionPoints() doesn't work with Circles.
// BUG(SolarLune): Line.IsColliding() and Line.GetIntersectionPoints() fail if testing two lines that intersect along the exact same slope.

// IsColliding returns if the Line is colliding with the other Shape. Currently, Circle-Line collision is missing.
func (l *Line) IsColliding(other Shape) bool {

	intersection := l.Check(other, 0, 0)

	r, ok := other.(*Rectangle)
	if ok && !intersection.Colliding() {
		return (l.X >= r.X && l.Y >= r.Y && l.X < r.X+r.W && l.Y < r.Y+r.H) || (l.X2 >= r.X && l.Y2 >= r.Y && l.X2 < r.X+r.W && l.Y2 < r.Y+r.H)
	}

	return intersection.Colliding()

}

// IntersectionPoint represents a point of intersection from a Line with another Shape, in absolute coordinates.
type IntersectionPoint struct {
	X, Y  float64
	Shape Shape
}

// IntersectionPoints returns the intersection points of a Line with another Shape as an array of IntersectionPoints.
// The returned list of intersection points are always sorted in order of distance from the start of the casting Line to each intersection.
// Currently, Circle-Line collision is missing.
func (l *Line) IntersectionPoints(other Shape) []IntersectionPoint {

	intersections := []IntersectionPoint{}

	switch b := other.(type) {

	case *Line:

		det := (l.X2-l.X)*(b.Y2-b.Y) - (b.X2-b.X)*(l.Y2-l.Y)

		if det != 0 {

			// MAGIC MATH; the extra + 1 here makes it so that corner cases work.

			lambda := (float32(((l.Y-b.Y)*(b.X2-b.X))-((l.X-b.X)*(b.Y2-b.Y))) + 1) / float32(det)

			gamma := (float32(((l.Y-b.Y)*(l.X2-l.X))-((l.X-b.X)*(l.Y2-l.Y))) + 1) / float32(det)

			if (0 < lambda && lambda < 1) && (0 < gamma && gamma < 1) {
				dx, dy := l.Delta()
				intersections = append(intersections, IntersectionPoint{l.X + float64(lambda*float32(dx)), l.Y + float64(lambda*float32(dy)), other})
			}

		}
	case *Rectangle:

		for _, side := range b.ToLines() {
			intersections = append(intersections, l.IntersectionPoints(side)...)
		}

	case *Space:
		for _, shape := range *b {
			intersections = append(intersections, l.IntersectionPoints(shape)...)
		}
		// case *Circle:
		// 	// 	TO-DO: Add this later, because this is kinda hard and would necessitate some complex vector math that, for whatever
		// 	//  reason, is not even readily available in a Golang library as far as I can tell???
		// 	break
	}

	// fmt.Println("WARNING! Object ", other, " isn't a valid shape for collision testing against Line ", l, "!")

	sort.Slice(intersections, func(i, j int) bool {
		return Distance(l.X, l.Y, intersections[i].X, intersections[i].Y) < Distance(l.X, l.Y, intersections[j].X, intersections[j].Y)
	})

	return intersections

}

// Check returns a Movement object indicating how far the Line can move before colliding with another Shape (or the original delta values, if there was no collision).
// Note that the Collision determines "how far in" a Line is for a Line-Line intersection by comparing the intersection points against the center point of the Line.
func (l *Line) Check(other Shape, dx, dy float64) *MovementCheck {

	col := newMovementCheck(l, other)

	switch b := other.(type) {

	case *Line:

		if intersectionPoints := b.IntersectionPoints(l); len(intersectionPoints) > 0 {

			point := intersectionPoints[0]

			col.addPoint(point.X, point.Y)

			ix, iy := -0.5, -0.5

			if dx != 0 {
				col.Dx = ix
			}

			if dy != 0 {
				col.Dy = iy
			}

			fmt.Println(col)

		}

	case *Rectangle:

		highestDx, highestDy := 0.0, 0.0

		for _, side := range b.ToLines() {

			sideCollision := l.Check(side, dx, dy)
			col.Points = append(col.Points, sideCollision.Points...)

			if sideCollision.Colliding() {

				// We want to use the collision's returned delta values, but only if it's higher than any other line. Imagine a situation where a rectangle is
				// moving up a slope. The left corner of the rectangle could return a lower or higher delta sliding value than the right one - we want to use
				// whichever is higher, since that would indicate a "stronger" movement.

				if math.Abs(sideCollision.Dx) > math.Abs(highestDx) {
					col.Dx = sideCollision.Dx
					highestDx = col.Dx
				}

				if math.Abs(sideCollision.Dy) > math.Abs(highestDy) {
					col.Dy = sideCollision.Dy
					highestDy = col.Dy
				}

			}

		}

	case *Circle:

		// TODO: Implement

		break

	case *Space:

		for _, shape := range *b {

			shapeCollision := l.Check(shape, dx, dy)
			col.Points = append(col.Points, shapeCollision.Points...)

			if shapeCollision.Colliding() {
				col.Dx = shapeCollision.Dx
				col.Dy = shapeCollision.Dy
				break
			}

		}

		// default:

		// 	for _, point := range l.IntersectionPoints(other) {
		// 		col.addPoint(point.X-l.X, point.Y-l.Y)
		// 	}

	}

	return col

}

// SetPosition sets the position of the Line, also moving the end point of the line (so it wholly moves the line to the
// specified position).
func (l *Line) SetPosition(x, y float64) {
	dx := x - l.X
	dy := y - l.Y
	l.X = x
	l.Y = y
	l.X2 += dx
	l.Y2 += dy
}

// Move moves the Line by the values specified.
func (l *Line) Move(dx, dy float64) {
	l.X += dx
	l.Y += dy
	l.X2 += dx
	l.Y2 += dy
}

// Center returns the center X and Y values of the Line.
func (l *Line) Center() (float64, float64) {

	x := l.X + ((l.X2 - l.X) / 2)
	y := l.Y + ((l.Y2 - l.Y) / 2)
	return x, y

}

// Length returns the length of the Line.
func (l *Line) Length() float64 {
	return Distance(l.X, l.Y, l.X2, l.Y2)
}

// SetLength sets the length of the Line to the value provided.
func (l *Line) SetLength(length float64) {

	ln := l.Length()
	xd := float64(float32(l.X2-l.X) / float32(ln) * float32(length))
	yd := float64(float32(l.Y2-l.Y) / float32(ln) * float32(length))

	l.X2 = l.X + xd
	l.Y2 = l.Y + yd
}

// GetBoundingRectangle returns a rectangle centered on the center point of the Line that would fully contain the Line.
func (l *Line) GetBoundingRectangle() *Rectangle {

	w := float64(math.Abs(float64(l.X2 - l.X)))
	h := float64(math.Abs(float64(l.Y2 - l.Y)))

	x := l.X

	if l.X2 < l.X {
		x = l.X2
	}

	y := l.Y

	if l.Y2 < l.Y {
		y = l.Y2
	}

	return NewRectangle(x, y, w, h)

}

// GetBoundingCircle returns a circle centered on the Line's central point that would fully contain the Line.
// func (l *Line) GetBoundingCircle() *Circle {

// 	x, y := l.Center()

// 	radius := float64(math.Abs(float64(l.X2 - l.X)))
// 	r2 := float64(math.Abs(float64(l.Y2 - l.Y)))

// 	if r2 > radius {
// 		radius = r2
// 	}

// 	return NewCircle(x, y, radius/2)

// }

// Delta returns the delta (or difference) between the start and end point of a Line.
func (l *Line) Delta() (float64, float64) {
	dx := l.X2 - l.X
	dy := l.Y2 - l.Y
	return dx, dy
}

// Slope returns the slope of the line.
func (l *Line) Slope() float64 {
	deltaX, deltaY := l.Delta()
	return deltaX / deltaY
}
