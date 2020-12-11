package resolv

import "fmt"

// A Circle represents an ordinary circle, and has a radius, in addition to normal shape properties.
type Circle struct {
	BasicShape
	Radius float64
}

// NewCircle returns a pointer to a new Circle object.
func NewCircle(x, y, radius float64) *Circle {
	c := &Circle{Radius: radius}
	c.X = x
	c.Y = y
	return c
}

// IsColliding returns true if the Circle is colliding with the specified other Shape, including the other Shape
// being wholly within the Circle.
func (c *Circle) IsColliding(other Shape) bool {

	switch b := other.(type) {

	case *Circle:
		return Distance(c.X, c.Y, b.X, b.Y) <= c.Radius+b.Radius
	case *Rectangle:
		closestX := c.X
		closestY := c.Y

		if c.X < b.X {
			closestX = b.X
		} else if c.X > b.X+b.W {
			closestX = b.X + b.W
		}

		if c.Y < b.Y {
			closestY = b.Y
		} else if c.Y > b.Y+b.H {
			closestY = b.Y + b.H
		}

		return Distance(c.X, c.Y, closestX, closestY) <= c.Radius
	case *Line:
		return b.IsColliding(c)
	case *Space:
		return b.IsColliding(c)

	}

	fmt.Println("WARNING! Object ", other, " isn't a valid shape for collision testing against Circle ", c, "!")

	return false

}

// WouldBeColliding returns whether the Circle would be colliding with the specified other Shape if it were to move
// in the specified direction.
func (c *Circle) WouldBeColliding(other Shape, dx, dy float64) bool {
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
