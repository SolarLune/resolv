package resolv

import "sort"

// CollisionPoint describes the point of intersection found when a Shape attempted to resolve a movement
// into another Shape in an isolated check, or when within the same Space as other existing Shapes.
// Dx and Dy represent the displacement of the Shape to the point of intersection - "how far in" the Shape
// got when attempting to move along the direction given by deltaX and deltaY in the Resolve() function.
// In other words, the delta X (dx) and delta y (dy).
// ShapeA is a pointer to the Shape that initiated the resolution check.
// ShapeB is a pointer to the Shape that the colliding object collided with, if the Collision was successful.
type CollisionPoint struct {
	X, Y float64
}

// MovementCheck represents a collision found when attempting to check movement of an object in a specified direction. The
type MovementCheck struct {
	Points []CollisionPoint
	ShapeA Shape
	ShapeB Shape
	Dx, Dy float64
}

func newMovementCheck(shapeA, shapeB Shape) *MovementCheck {
	return &MovementCheck{
		Points: []CollisionPoint{},
		ShapeA: shapeA,
		ShapeB: shapeB,
	}
}

func (movementCheck *MovementCheck) addPoint(x, y float64) {
	movementCheck.Points = append(movementCheck.Points, CollisionPoint{x, y})
}

func (movementCheck *MovementCheck) Colliding() bool {
	return len(movementCheck.Points) > 0
}

func (movementCheck *MovementCheck) sortByDistance() {
	sort.Slice(movementCheck.Points, func(i, j int) bool {
		return Distance(0, 0, movementCheck.Points[i].X, movementCheck.Points[i].Y) < Distance(0, 0, movementCheck.Points[j].X, movementCheck.Points[j].Y)
	})
}
