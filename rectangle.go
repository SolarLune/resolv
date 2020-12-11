package resolv

// Rectangle represents a rectangle.
type Rectangle struct {
	BasicShape
	W, H float64
}

// NewRectangle creates a new Rectangle and returns a pointer to it.
func NewRectangle(x, y, w, h float64) *Rectangle {
	r := &Rectangle{W: w, H: h}
	r.X = x
	r.Y = y
	return r
}

// IsColliding returns whether the Rectangle is colliding with the specified other Shape or not, including the other Shape
// being wholly contained within the Rectangle.
func (r *Rectangle) IsColliding(other Shape) bool {

	switch b := other.(type) {
	case *Rectangle:
		return r.X > b.X-r.W && r.Y > b.Y-r.H && r.X < b.X+b.W && r.Y < b.Y+b.H
	default:
		return b.IsColliding(r)
	}

}

// WouldBeColliding returns whether the Rectangle would be colliding with the other Shape if it were to move in the
// specified direction.
func (r *Rectangle) WouldBeColliding(other Shape, dx, dy float64) bool {
	r.X += dx
	r.Y += dy
	isColliding := r.IsColliding(other)
	r.X -= dx
	r.Y -= dy
	return isColliding
}

// Center returns the center point of the Rectangle.
func (r *Rectangle) Center() (float64, float64) {

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
