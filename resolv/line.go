package resolv

import (
	"math"
	"sort"
)

// Line represents a line, from one point to another.
type Line struct {
	BasicShape
	X2, Y2 int32
}

// NewLine returns a new Line instance.
func NewLine(x, y, x2, y2 int32) *Line {
	l := &Line{}
	l.X = x
	l.Y = y
	l.X2 = x2
	l.Y2 = y2
	l.Collideable = true
	return l
}

// BUG(SolarLune): Line.IsColliding() and Line.IntersectionPoints() doesn't work with Circles.
// BUG(SolarLune): Line.IsColliding() and Line.IntersectionPoints() fail if testing two lines that intersect along the exact same slope.

// IsColliding returns if the Line is colliding with the other Shape. Currently, Circle-Line collision is missing.
func (l *Line) IsColliding(other Shape) bool {

	intersectionPoints := l.IntersectionPoints(other)

	colliding := len(intersectionPoints) > 0

	r, ok := other.(*Rectangle)
	if ok && !colliding {
		return (l.X >= r.X && l.Y >= r.Y && l.X < r.X+r.W && l.Y < r.Y+r.H) || (l.X2 >= r.X && l.Y2 >= r.Y && l.X2 < r.X+r.W && l.Y2 < r.Y+r.H)
	}

	return colliding

}

// IntersectionPoint represents a point of intersection from a Line with another Shape.
type IntersectionPoint struct {
	X, Y  int32
	Shape Shape
}

// IntersectionPoints returns the intersection points of a Line with another Shape as an array of arrays, composed of X and Y position int32s.
// The returned list of intersection points are always sorted in order of distance from the start of the casting Line to the end.
// Currently, Circle-Line collision is missing.
func (l *Line) IntersectionPoints(other Shape) []IntersectionPoint {

	intersections := []IntersectionPoint{}

	if !l.Collideable || !other.IsCollideable() {
		return intersections
	}

	b, ok := other.(*Line)

	if ok {

		det := (l.X2-l.X)*(b.Y2-b.Y) - (b.X2-b.X)*(l.Y2-l.Y)

		if det != 0 {

			// MAGIC MATH; the extra + 1 here makes it so that corner cases (literally aiming the line through the corners of the
			// hollow square in world5) works!

			lambda := (float32(((l.Y-b.Y)*(b.X2-b.X))-((l.X-b.X)*(b.Y2-b.Y))) + 1) / float32(det)

			gamma := (float32(((l.Y-b.Y)*(l.X2-l.X))-((l.X-b.X)*(l.Y2-l.Y))) + 1) / float32(det)

			if (0 < lambda && lambda < 1) && (0 < gamma && gamma < 1) {
				dx, dy := l.GetDelta()
				intersections = append(intersections, IntersectionPoint{l.X + int32(lambda*float32(dx)), l.Y + int32(lambda*float32(dy)), other})
			}

		}

	}

	r, ok := other.(*Rectangle)

	if ok {

		side := NewLine(r.X, r.Y, r.X, r.Y+r.H)
		intersections = append(intersections, l.IntersectionPoints(side)...)

		side.Y = r.Y + r.H
		side.X2 = r.X + r.W
		side.Y2 = r.Y + r.H
		intersections = append(intersections, l.IntersectionPoints(side)...)

		side.X = r.X + r.W
		side.Y2 = r.Y
		intersections = append(intersections, l.IntersectionPoints(side)...)

		side.Y = r.Y
		side.X2 = r.X
		side.Y2 = r.Y
		intersections = append(intersections, l.IntersectionPoints(side)...)

	}

	_, ok = other.(*Circle)

	if ok {

		// return false

		// 	TO-DO: Add this later, because this is kinda hard and would necessitate some complex vector math that, for whatever
		//  reason, is not even readily available in a Golang library as far as I can tell???

	}

	sp, ok := other.(*Space)

	if ok {
		for _, shape := range *sp {
			intersections = append(intersections, l.IntersectionPoints(shape)...)
		}
	}

	// fmt.Println("WARNING! Object ", other, " isn't a valid shape for collision testing against Line ", l, "!")

	sort.Slice(intersections, func(i, j int) bool {
		return Distance(l.X, l.Y, intersections[i].X, intersections[i].Y) < Distance(l.X, l.Y, intersections[j].X, intersections[j].Y)
	})

	return intersections

}

// WouldBeColliding returns if the Line would be colliding if it were moved by the designated delta X and Y values.
func (l *Line) WouldBeColliding(other Shape, dx, dy int32) bool {
	l.X += dx
	l.Y += dy
	l.X2 += dx
	l.Y2 += dy
	isColliding := l.IsColliding(other)
	l.X -= dx
	l.Y -= dy
	l.X2 -= dx
	l.Y2 -= dy
	return isColliding
}

// SetXY sets the position of the Line, also moving the end point of the line (so it wholly moves the line to the
// specified position).
func (l *Line) SetXY(x, y int32) {
	dx := x - l.X
	dy := y - l.Y
	l.X = x
	l.Y = y
	l.X2 += dx
	l.Y2 += dy
}

// Move moves the Line by the values specified.
func (l *Line) Move(x, y int32) {
	l.X += x
	l.Y += y
	l.X2 += x
	l.Y2 += y
}

// Center returns the center X and Y values of the Line.
func (l *Line) Center() (int32, int32) {

	x := l.X + ((l.X2 - l.X) / 2)
	y := l.Y + ((l.Y2 - l.Y) / 2)
	return x, y

}

// GetLength returns the length of the Line.
func (l *Line) GetLength() int32 {
	return Distance(l.X, l.Y, l.X2, l.Y2)
}

// SetLength sets the length of the Line to the value provided.
func (l *Line) SetLength(length int32) {

	ln := l.GetLength()
	xd := int32(float32(l.X2-l.X) / float32(ln) * float32(length))
	yd := int32(float32(l.Y2-l.Y) / float32(ln) * float32(length))

	l.X2 = l.X + xd
	l.Y2 = l.Y + yd
}

// GetBoundingRectangle returns a rectangle centered on the center point of the Line that would fully contain the Line.
func (l *Line) GetBoundingRectangle() *Rectangle {

	w := int32(math.Abs(float64(l.X2 - l.X)))
	h := int32(math.Abs(float64(l.Y2 - l.Y)))

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
func (l *Line) GetBoundingCircle() *Circle {

	x, y := l.Center()

	radius := int32(math.Abs(float64(l.X2 - l.X)))
	r2 := int32(math.Abs(float64(l.Y2 - l.Y)))

	if r2 > radius {
		radius = r2
	}

	return NewCircle(x, y, radius/2)

}

// GetDelta returns the delta (or difference) between the start and end point of a Line.
func (l *Line) GetDelta() (int32, int32) {
	dx := l.X2 - l.X
	dy := l.Y2 - l.Y
	return dx, dy
}
