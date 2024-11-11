package resolv

// Circle represents a circle (naturally), and is essentially a point with a radius.
type Circle struct {
	ShapeBase
	radius float64
}

// NewCircle returns a new Circle, with its center at the X and Y position given, and with the defined radius.
func NewCircle(x, y, radius float64) *Circle {
	circle := &Circle{
		ShapeBase: newShapeBase(x, y),
		radius:    radius,
	}
	circle.ShapeBase.owner = circle
	return circle
}

// Clone clones the Circle.
func (c *Circle) Clone() IShape {
	newCircle := NewCircle(c.position.X, c.position.Y, c.radius)
	newCircle.tags.Set(*c.tags)
	newCircle.ShapeBase = c.ShapeBase
	newCircle.id = globalShapeID
	globalShapeID++
	newCircle.ShapeBase.space = nil
	newCircle.ShapeBase.touchingCells = []*Cell{}
	newCircle.ShapeBase.owner = newCircle
	return newCircle
}

// Bounds returns the top-left and bottom-right corners of the Circle.
func (c *Circle) Bounds() Bounds {
	return Bounds{
		Min:   Vector{c.position.X - c.radius, c.position.Y - c.radius},
		Max:   Vector{c.position.X + c.radius, c.position.Y + c.radius},
		space: c.space,
	}
}

func (c *Circle) Project(axis Vector) Projection {
	axis = axis.Unit()
	projectedCenter := axis.Dot(c.position)
	projectedRadius := axis.Magnitude() * c.radius

	min := projectedCenter - projectedRadius
	max := projectedCenter + projectedRadius

	if min > max {
		min, max = max, min
	}

	return Projection{min, max}
}

// Radius returns the radius of the Circle.
func (c *Circle) Radius() float64 {
	return c.radius
}

// SetRadius sets the radius of the Circle, updating the scale multiplier to reflect this change.
func (c *Circle) SetRadius(radius float64) {
	c.radius = radius
	c.update()
}

// Intersection returns an IntersectionSet for the other Shape provided.
// If no intersection is detected, the IntersectionSet returned is empty.
func (c *Circle) Intersection(other IShape) IntersectionSet {

	switch otherShape := other.(type) {
	case *ConvexPolygon:
		return circleConvexTest(c, otherShape)

	case *Circle:
		return circleCircleTest(c, otherShape)
	}

	// This should never happen
	panic("Unimplemented intersection")

}
