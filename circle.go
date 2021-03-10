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

// IsColliding returns if the Circle is colliding with the other Shape.
func (c *Circle) IsColliding(other Shape) bool {
	return c.Check(other, 0, 0).Colliding()
}

// Check returns a MovementCheck object for a proposed movement from the
func (c *Circle) Check(other Shape, dx, dy float64) *MovementCheck {

	col := newMovementCheck(c, other)

	switch b := other.(type) {

	case *Circle:

		d := Distance(c.X, c.Y, b.X, b.Y)

		if d <= c.Radius+b.Radius {
			col.addPoint(c.X, c.Y)
			// col.Dx, col.Dy
		}

	// return Distance(c.X, c.Y, b.X, b.Y) <= c.Radius+b.Radius
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

		if Distance(c.X, c.Y, closestX, closestY) <= c.Radius {
			col.addPoint(closestX, closestY)
			// col.Dx, col.Dy
		}

		// case *Line:
		// 	return b.IsColliding(c)
		// case *Space:
		// 	return b.IsColliding(c)

	}

	fmt.Println("WARNING! Object ", other, " isn't a valid shape for collision testing against Circle ", c, "!")

	return col

}

// BoundingRectangle returns a Rectangle which has a width and height of 2*Radius (and so, contains the Circle).
func (c *Circle) BoundingRectangle() *Rectangle {
	r := &Rectangle{}
	r.W = c.Radius * 2
	r.H = c.Radius * 2
	r.X = c.X - r.W/2
	r.Y = c.Y - r.H/2
	return r
}
