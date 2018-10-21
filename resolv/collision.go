package resolv

// Collision describes the collision found when a Shape attempted to resolve a movement into another Shape in an
// isolated check, or when within the same Space as other existing Shapes.
// ResolveX and ResolveY represent the displacement of the Shape to the point of collision. How far along the Shape
// got when attempting to move along the direction given by deltaX and deltaY in the Resolve() function before
// touching another Shape.
// Teleporting is if moving according to ResolveX and ResolveY might be considered teleporting, which is moving
// greater than the deltaX or deltaY provided to the Resolve function * 1.5 (this is arbitrary, but can be useful
// when attempting to see if a movement would be ).
// ShapeA is a pointer to the Shape that initiated the resolution check.
// ShapeB is a pointer to the Shape that the colliding object collided with, if the Collision was successful.
type Collision struct {
	ResolveX, ResolveY int32
	Teleporting        bool
	ShapeA             Shape
	ShapeB             Shape
}

// Colliding returns whether the Collision actually was valid because of a collision against another Shape.
func (c *Collision) Colliding() bool {
	return c.ShapeB != nil
}
