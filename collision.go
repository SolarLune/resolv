package resolv

import (
	"math"
)

// Collision contains the results of an Object.Check() call. The Objects array indicate the Objects collided with.
type Collision struct {
	checkingObject *Object   // The checking object
	dx, dy         float64   // The delta the checking object was moving on that caused this collision
	Objects        []*Object // List of objects that would be collided with
}

func NewCollision() *Collision {
	return &Collision{
		Objects: []*Object{},
	}
}

// HasTags returns whether any objects within the Collision have all of the specified tags. Skips the object the check was performed with, of course.
func (cc *Collision) HasTags(tags ...string) bool {

	for _, o := range cc.Objects {

		if o == cc.checkingObject {
			continue
		}
		if o.HasTags(tags...) {
			return true
		}

	}

	return false
}

// ObjectsByTags returns a slice of Objects from the cells reported by a Collision object by searching for a specific set of tags. This slice does not contain the checking object.
func (cc *Collision) ObjectsByTags(tags ...string) []*Object {

	objects := []*Object{}

	for _, o := range cc.Objects {

		if o == cc.checkingObject {
			continue
		}
		if o.HasTags(tags...) {
			objects = append(objects, o)
		}

	}

	return objects

}

// // Objects returns a list of objects contained in the cells (which are sorted according to distance to the checking Object) returned in the Collision.
// // This slice does not contain the checking object.
// func (cc *Collision) Objects() []*Object {

// 	objectsSet := map[*Object]bool{} // We only want to count each Object once
// 	objects := []*Object{}

// 	for _, cell := range cc.Cells {

// 		for _, o := range cell.Objects {

// 			_, inSet := objectsSet[o]

// 			if o != cc.checkingObject && !inSet {
// 				objects = append(objects, o)
// 				objectsSet[o] = true
// 			}

// 		}

// 	}

// 	return objects
// }

// A Delta represents a movement from one position to another, either to come into contact with another Cell (to come into contact) or to move on a perpendicular axis to avoid contact (to slide).
type Delta struct {
	X, Y float64 // A 2-dimensional vector [X, Y] representing how much to move
}

// ContactWithObject returns the delta to move to come into contact with the specified Object.
func (cc *Collision) ContactWithObject(object *Object) *Delta {

	delta := &Delta{}

	currentDx, currentDy := object.X, object.Y

	if cc.dx < 0 {
		currentDx += float64(object.W)
		currentDx -= cc.checkingObject.X
	} else if cc.dx > 0 {
		currentDx -= (cc.checkingObject.X + cc.checkingObject.W)
	} else {
		currentDx = 0
	}

	if cc.dy < 0 {
		currentDy += float64(object.H)
		currentDy -= cc.checkingObject.Y
	} else if cc.dy > 0 {
		currentDy -= (cc.checkingObject.Y + cc.checkingObject.H)
	} else {
		currentDy = 0
	}

	if currentDx > 0 {
		delta.X = math.Max(delta.X, currentDx)
	} else {
		delta.X = math.Min(delta.X, currentDx)
	}

	if currentDy > 0 {
		delta.Y = math.Max(delta.Y, currentDy)
	} else {
		delta.Y = math.Min(delta.Y, currentDy)
	}

	return delta

}

// // SlideDelta returns how much distance the calling Object can slide to avoid a collision with the targetObject. If no targetObject is given,
// // this works with the cells in the Space. This only works on vertical and horizontal axes (x and y directly). avoidTags is a sequence of
// // tags (as strings) to indicate when sliding is valid (i.e. if )
// func (cc *Collision) SlideToAvoidCell(cell *Cell, avoidTags ...string) {

// 	diffX, diffY := 0.0, 0.0

// 	sp := cc.checkingObject.Space

// 	collidingCell := cc.Cells[0]
// 	ccX, ccY := sp.SpaceToWorld(collidingCell.X, collidingCell.Y)
// 	hX := float64(sp.CellWidth) / 2.0
// 	hY := float64(sp.CellHeight) / 2.0

// 	// Previous precise collisions

// 	// if targetObject != nil {
// 	// 	ccX, ccY = targetObject.X, targetObject.Y
// 	// 	hX = targetObject.W / 2
// 	// 	hY = targetObject.H / 2
// 	// }

// 	ccX += hX
// 	ccY += hY

// 	oX, oY := cc.checkingObject.Center()

// 	diffX = oX - ccX
// 	diffY = oY - ccY

// 	left := sp.Cell(collidingCell.X-1, collidingCell.Y)
// 	right := sp.Cell(collidingCell.X+1, collidingCell.Y)
// 	up := sp.Cell(collidingCell.X, collidingCell.Y-1)
// 	down := sp.Cell(collidingCell.X, collidingCell.Y+1)

// 	if cc.dy != 0 {
// 		if diffX > 0 && (right == nil || !right.ContainsTags(avoidTags...)) {
// 			// Slide right
// 			cc.Slide.Vector[0] = ccX + hX - cc.checkingObject.X
// 			cc.Slide.Valid = true

// 		} else if diffX < 0 && (left == nil || !left.ContainsTags(avoidTags...)) {
// 			// Slide left
// 			cc.Slide.Vector[0] = ccX - hX - (cc.checkingObject.X + cc.checkingObject.W)
// 			cc.Slide.Valid = true

// 		}
// 	}

// 	if cc.dx != 0 {
// 		if diffY > 0 && (down == nil || !down.ContainsTags(avoidTags...)) {
// 			// Slide down
// 			cc.Slide.Vector[1] = ccY + hY - cc.checkingObject.Y
// 			cc.Slide.Valid = true

// 		} else if diffY < 0 && (up == nil || !up.ContainsTags(avoidTags...)) {
// 			// Slide up
// 			cc.Slide.Vector[1] = ccY - hY - (cc.checkingObject.Y + cc.checkingObject.H)
// 			cc.Slide.Valid = true

// 		}
// 	}

// }
