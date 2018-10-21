package resolv

import "fmt"

// A Circle represents an ordinary circle, and has a radius, in addition to normal shape properties.
type Circle struct {
	BasicShape
	Radius int32
}

// NewCircle returns a pointer to a new Circle object.
func NewCircle(x, y, radius int32) *Circle {
	c := &Circle{Radius: radius}
	c.X = x
	c.Y = y
	c.Collideable = true
	return c
}

// IsColliding returns true if the Circle is colliding with the specified other Shape, including the other Shape
// being wholly within the Circle.
func (c *Circle) IsColliding(other Shape) bool {

	if !c.Collideable || !other.IsCollideable() {
		return false
	}

	b, ok := other.(*Circle)

	if ok {

		return Distance(c.X, c.Y, b.X, b.Y) <= c.Radius+b.Radius

	}

	r, ok := other.(*Rectangle)

	if ok {

		closestX := c.X
		closestY := c.Y

		if c.X < r.X {
			closestX = r.X
		} else if c.X > r.X+r.W {
			closestX = r.X + r.W
		}

		if c.Y < r.Y {
			closestY = r.Y
		} else if c.Y > r.Y+r.H {
			closestY = r.Y + r.H
		}

		return Distance(c.X, c.Y, closestX, closestY) <= c.Radius

	}

	l, ok := other.(*Line)

	if ok {
		return l.IsColliding(c)
	}

	sp, ok := other.(*Space)

	if ok {
		return sp.IsColliding(r)
	}

	fmt.Println("WARNING! Object ", other, " isn't a valid shape for collision testing against Circle ", c, "!")

	return false

}

// WouldBeColliding returns whether the Rectangle would be colliding with the specified other Shape if it were to move
// in the specified direction.
func (c *Circle) WouldBeColliding(other Shape, dx, dy int32) bool {
	c.X += dx
	c.Y += dy
	isColliding := c.IsColliding(other)
	c.X -= dx
	c.Y -= dy
	return isColliding
}

// GetBoundingRect returns a Rectangle which has a width and height of 2*Radius.
func (c *Circle) GetBoundingRect() *Rectangle {
	r := &Rectangle{}
	r.W = c.Radius * 2
	r.H = c.Radius * 2
	r.X = c.X - r.W/2
	r.Y = c.Y - r.H/2
	return r
}
