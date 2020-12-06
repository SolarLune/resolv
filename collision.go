package resolv

import (
	"math"
)

// Collision contains the results of an Object.Check() call. The Cells array indicate the Cells returned in the check, and
// ContactX and ContactY indicate the delta between the calling object's leading edge and the checked Cells' closest edges.
type Collision struct {
	checkingObject *Object
	Cells          []*Cell // Cells returned in the collision check.
	ContactX       float64 // Distance to come into contact with closest edge of the closest Cell on the X-axis.
	ContactY       float64 // Distance to come into contact with closest edge of the closest Cell on the Y-axis.
	CanSlide       bool    // If the checking object can slide on the perpendicular axis to the movement delta to evade the collision (i.e. if moving upwards to collide with a wall, CanSlide indicates if it's possible to move left or right to avoid the wall).
	SlideX         float64 // Distance to move to evade the closest edge of the closest Cell on the X-axis.
	SlideY         float64 // Distance to move to evade the closest edge of the closest Cell on the Y-axis.
}

// HasTags returns whether any objects within the CellCheck have all of the specified tags. Skips the object the check was performed with, of course
func (cc *Collision) HasTags(tags ...string) bool {

	for _, c := range cc.Cells {

		for _, o := range c.Objects {

			if o == cc.checkingObject {
				continue
			}
			if o.HasTags(tags...) {
				return true
			}

		}

	}

	return false
}

// Valid returns whether a CellCheck contains any Cells.
func (cc *Collision) Valid() bool {
	return len(cc.Cells) > 0
}

// ObjectsByTags returns an object from a cell reported by a Collision object by searching for a set of tags.
func (cc *Collision) ObjectsByTags(tags ...string) []*Object {

	objects := []*Object{}

	for _, c := range cc.Cells {

		for _, o := range c.Objects {

			if o == cc.checkingObject {
				continue
			}
			if o.HasTags(tags...) {
				objects = append(objects, o)
			}

		}

	}

	return objects

}

// Objects returns a list of objects contained in the cells returned in the Collision.
func (cc *Collision) Objects() []*Object {
	objects := []*Object{}
	for _, cell := range cc.Cells {

		for _, o := range cell.Objects {

			if o != cc.checkingObject {
				objects = append(objects, o)
			}

		}

	}
	return objects
}

func (cc *Collision) calculateContactDelta(dx, dy float64) {

	if cc.Valid() {

		// By definition, because we're checking multiple cells but generally only one direction at a time
		// and the only objects that exist currently are axis-aligned rectangles, any cell is as good as any other.
		// If the shapes were customizeable, then things would change and we would need to grab the cell that's, perhaps,
		// closest to the center of the Object.

		target := cc.Cells[0]

		deltaX, deltaY := cc.checkingObject.Space.SpaceToWorld(target.X, target.Y)

		if cc.checkingObject.PreciseCollision {
			targetObj := cc.Objects()[0]

			deltaX, deltaY = targetObj.X, targetObj.Y

			if dx < 0 {
				deltaX += float64(targetObj.W)
				deltaX -= cc.checkingObject.X
			} else if dx > 0 {
				deltaX -= (cc.checkingObject.X + cc.checkingObject.W)
			} else {
				deltaX = 0
			}

			if dy < 0 {
				deltaY += float64(targetObj.H)
				deltaY -= cc.checkingObject.Y
			} else if dy > 0 {
				deltaY -= (cc.checkingObject.Y + cc.checkingObject.H)
			} else {
				deltaY = 0
			}

		} else {

			if dx < 0 {
				deltaX += float64(cc.checkingObject.Space.CellWidth)
				deltaX -= cc.checkingObject.X
			} else if dx > 0 {
				deltaX -= (cc.checkingObject.X + cc.checkingObject.W)
			} else {
				deltaX = 0
			}

			if dy < 0 {
				deltaY += float64(cc.checkingObject.Space.CellHeight)
				deltaY -= cc.checkingObject.Y
			} else if dy > 0 {
				deltaY -= (cc.checkingObject.Y + cc.checkingObject.H)
			} else {
				deltaY = 0
			}

		}

		cc.ContactX, cc.ContactY = deltaX, deltaY

	}

}

func (cc *Collision) calculateSlideDelta(dx, dy float64, tags ...string) {

	if cc.Valid() {

		sp := cc.checkingObject.Space

		collidingCell := cc.Cells[0]
		ccX, ccY := sp.SpaceToWorld(collidingCell.X, collidingCell.Y)
		hX := float64(sp.CellWidth) / 2.0
		hY := float64(sp.CellHeight) / 2.0

		if cc.checkingObject.PreciseCollision {
			obj := cc.Objects()[0]
			ccX, ccY = obj.X, obj.Y
			hX = obj.W / 2
			hY = obj.H / 2
		}

		ccX += hX
		ccY += hY

		oX, oY := cc.checkingObject.Center()

		diffX := oX - ccX
		diffY := oY - ccY

		left := sp.Cell(collidingCell.X-1, collidingCell.Y)
		right := sp.Cell(collidingCell.X+1, collidingCell.Y)
		up := sp.Cell(collidingCell.X, collidingCell.Y-1)
		down := sp.Cell(collidingCell.X, collidingCell.Y+1)

		cc.CanSlide = false

		if dy != 0 {
			if diffX > 0 && (right == nil || !right.ContainsTags(tags...)) {
				// Slide right
				diffX = ccX + hX - cc.checkingObject.X
				cc.CanSlide = true
			} else if diffX < 0 && (left == nil || !left.ContainsTags(tags...)) {
				// Slide left
				diffX = ccX - hX - (cc.checkingObject.X + cc.checkingObject.W)
				cc.CanSlide = true
			}
		}

		if dx != 0 {
			if diffY > 0 && (down == nil || !down.ContainsTags(tags...)) {
				// Slide down
				diffY = ccY + hY - cc.checkingObject.Y
				cc.CanSlide = true
			} else if diffY < 0 && (up == nil || !up.ContainsTags(tags...)) {
				// Slide up
				diffY = ccY - hY - (cc.checkingObject.Y + cc.checkingObject.H)
				cc.CanSlide = true
			}
		}

		cc.SlideX = diffX
		cc.SlideY = diffY

	}

}

func distance(x, y, x2, y2 float64) float64 {

	dx := x - x2
	dy := y - y2
	ds := (dx * dx) + (dy * dy)
	return math.Sqrt(math.Abs(ds))

}
