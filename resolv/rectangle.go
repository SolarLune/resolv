package resolv

import "fmt"

// Rectangle represents a rectangle.
type Rectangle struct {
	BasicShape
	W, H int32
}

// NewRectangle creates a new Rectangle and returns a pointer to it.
func NewRectangle(x, y, w, h int32) *Rectangle {
	r := &Rectangle{W: w, H: h}
	r.X = x
	r.Y = y
	r.Collideable = true
	return r
}

// IsColliding returns whether the Rectangle is colliding with the specified other Shape or not, including the other Shape
// being wholly contained within the Rectangle.
func (r *Rectangle) IsColliding(other Shape) bool {

	if !r.Collideable || !other.IsCollideable() {
		return false
	}

	b, ok := other.(*Rectangle)

	if ok {
		return r.X > b.X-r.W && r.Y > b.Y-r.H && r.X < b.X+b.W && r.Y < b.Y+b.H
	}

	c, ok := other.(*Circle)

	if ok {
		return c.IsColliding(r)
	}

	l, ok := other.(*Line)

	if ok {
		return l.IsColliding(r)
	}

	sp, ok := other.(*Space)

	if ok {
		return sp.IsColliding(r)
	}

	fmt.Println("WARNING! Object ", other, " isn't a valid shape for collision testing against a Rectangle ", r, "!")

	return false
}

// WouldBeColliding returns whether the Rectangle would be colliding with the other Shape if it were to move in the
// specified direction.
func (r *Rectangle) WouldBeColliding(other Shape, dx, dy int32) bool {
	r.X += dx
	r.Y += dy
	isColliding := r.IsColliding(other)
	r.X -= dx
	r.Y -= dy
	return isColliding
}

// Center returns the center point of the Rectangle.
func (r *Rectangle) Center() (int32, int32) {

	x := r.X + r.W/2
	y := r.Y + r.H/2

	return x, y

}

// GetBoundingCircle returns a circle that wholly contains the Rectangle.
func (r *Rectangle) GetBoundingCircle() *Circle {

	x, y := r.Center()
	c := NewCircle(x, y, Distance(x, y, r.X+r.W, r.Y))
	return c

}
