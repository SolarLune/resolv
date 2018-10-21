package resolv

import (
	"fmt"
	"math"
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

// BUG(SolarLune): Line.IsColliding() doesn't work with Circles.
// BUG(SolarLune): Line.IsColliding() fails if testing two lines who intersect along the exact same slope.

// IsColliding returns if the Line is colliding with the other Shape. Currently, Circle-Line collision is missing.
func (l *Line) IsColliding(other Shape) bool {

	if !l.Collideable || !other.IsCollideable() {
		return false
	}

	b, ok := other.(*Line)

	if ok {

		det := (l.X2-l.X)*(b.Y2-b.Y) - (b.X2-b.X)*(l.Y2-l.Y)

		if det == 0 {
			return false
		}

		// MAGIC MATH; the extra + 1 here makes it so that corner cases (literally aiming the line through the corners of the
		// hollow square in world5) works!

		lambda := (float32(((l.Y-b.Y)*(b.X2-b.X))-((l.X-b.X)*(b.Y2-b.Y))) + 1) / float32(det)

		gamma := (float32(((l.Y-b.Y)*(l.X2-l.X))-((l.X-b.X)*(l.Y2-l.Y))) + 1) / float32(det)

		return (0 < lambda && lambda < 1) && (0 < gamma && gamma < 1)

	}

	r, ok := other.(*Rectangle)

	if ok {

		side := NewLine(r.X, r.Y, r.X, r.Y+r.H)
		if l.IsColliding(side) {
			return true
		}

		side.Y = r.Y + r.H
		side.X2 = r.X + r.W
		side.Y2 = r.Y + r.H
		if l.IsColliding(side) {
			return true
		}

		side.X = r.X + r.W
		side.Y2 = r.Y
		if l.IsColliding(side) {
			return true
		}

		side.Y = r.Y
		side.X2 = r.X
		side.Y2 = r.Y
		if l.IsColliding(side) {
			return true
		}

		return (l.X >= r.X && l.Y >= r.Y && l.X < r.X+r.W && l.Y < r.Y+r.H) || (l.X2 >= r.X && l.Y2 >= r.Y && l.X2 < r.X+r.W && l.Y2 < r.Y+r.H)

	}

	_, ok = other.(*Circle)

	if ok {

		return false

		// 	TO-DO: Add this later, because this is kinda hard and would necessitate some complex vector math that, for whatever
		//  reason, is not even readily available in a Golang library as far as I can tell???

	}

	sp, ok := other.(*Space)

	if ok {
		return sp.IsColliding(l)
	}

	fmt.Println("WARNING! Object ", other, " isn't a valid shape for collision testing against Line ", l, "!")

	return false

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
