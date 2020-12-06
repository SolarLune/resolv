package resolv

import (
	"math"
	"sort"
)

// Object represents an object that can be spread across one or more Cells in a Space.
type Object struct {
	Space            *Space      // Reference to the Space the Object exists within
	X, Y, W, H       float64     // Position and size of the Object in the Space
	TouchingCells    []*Cell     // An array of Cells the Object is touching
	tags             []string    // A list of tags the Object has
	Data             interface{} // A pointer to a user-definable object
	PreciseCollision bool        // Whether collisions are precise or if just using cellular positions are good enough.
}

// NewObject returns a new Object of the specified position and size, and which exists in the Space provided.
func NewObject(x, y, w, h float64, space *Space) *Object {
	o := &Object{
		X:     x,
		Y:     y,
		W:     w,
		H:     h,
		tags:  []string{},
		Space: space,
	}
	o.Update()
	return o
}

// Update updates the object's association to the Cells in the Space. This should be called whenever an Object is moved.
// This is automatically called once when creating the Object, so you don't have to call it for static objects.
func (obj *Object) Update() {

	if obj.Space != nil {

		obj.Remove()

		cx, cy, ex, ey := obj.ShapeToSpace(0, 0)

		for y := cy; y <= ey; y++ {

			for x := cx; x <= ex; x++ {

				c := obj.Space.Cell(x, y)

				if c != nil {
					c.register(obj)
					obj.TouchingCells = append(obj.TouchingCells, c)
				}

			}

		}

		obj.Space.Objects = append(obj.Space.Objects, obj)

	}

}

// Remove removes the Object from being associated with the Space. This should be done whenever an Object is removed from the
// game.
func (obj *Object) Remove() {

	for _, cell := range obj.TouchingCells {
		cell.unregister(obj)
	}

	obj.TouchingCells = []*Cell{}

	for i, o := range obj.Space.Objects {
		if o == obj {
			obj.Space.Objects[i] = obj.Space.Objects[len(obj.Space.Objects)-1]
			obj.Space.Objects = obj.Space.Objects[:len(obj.Space.Objects)-1]
			break
		}
	}

}

// AddTag adds a tag to the Object.
func (obj *Object) AddTag(tags ...string) {
	obj.tags = append(obj.tags, tags...)
}

// RemoveTag removes a tag from the Object.
func (obj *Object) RemoveTag(tags ...string) {

	for _, tag := range tags {

		for i, t := range obj.tags {

			if t == tag {
				obj.tags = append(obj.tags[:i], obj.tags[i+1:]...)
				break
			}

		}

	}

}

// HasTags indicates if an Object has all of the tags indicated.
func (obj *Object) HasTags(tags ...string) bool {

	for _, tag := range tags {

		hasTag := false

		for _, t := range obj.tags {

			if t == tag {
				hasTag = true
				break
			}

		}

		if !hasTag {
			return false
		}

	}

	return true

}

// Tags returns the tags an Object has.
func (obj *Object) Tags() []string {
	tags := []string{}
	for _, t := range obj.tags {
		tags = append(tags, t)
	}
	return tags
}

// ShapeToSpace returns the Space coordinates of the shape, given its world position and size.
func (obj *Object) ShapeToSpace(dx, dy float64) (int, int, int, int) {
	cx, cy := obj.Space.WorldToSpace(obj.X+dx, obj.Y+dy)
	ex, ey := obj.Space.WorldToSpace(obj.X+obj.W+dx-1, obj.Y+obj.H+dy-1)
	return cx, cy, ex, ey
}

// SharesCells returns whether the Object occupies a cell shared by the specified other Object.
func (obj *Object) SharesCells(other *Object) bool {
	for _, cell := range obj.TouchingCells {
		if cell.Contains(other) {
			return true
		}
	}
	return false
}

// SharesCellsTags returns if the Cells the Object occupies have an object with the specified tags.
func (obj *Object) SharesCellsTags(tags ...string) bool {
	for _, cell := range obj.TouchingCells {
		if cell.ContainsTags(tags...) {
			return true
		}
	}
	return false
}

func (obj *Object) Center() (float64, float64) {
	return obj.X + (obj.W / 2.0), obj.Y + (obj.H / 2.0)
}

func (obj *Object) SetCenter(x, y float64) {
	obj.X = x - (obj.W / 2)
	obj.Y = y - (obj.H / 2)
}

func (obj *Object) CellPosition() (int, int) {
	cx, cy := obj.Center()
	return obj.Space.WorldToSpace(cx, cy)
}

// SetRight sets the X position of the Object so the right edge is at the X position given.
func (obj *Object) SetRight(x float64) {
	obj.X = x - obj.W
}

// SetBottom sets the Y position of the Object so that the bottom edge is at the Y position given.
func (obj *Object) SetBottom(y float64) {
	obj.Y = y - obj.H
}

// Check checks movement for this object from its current position using the designated delta movement direction (dx, dy). tag indicates what tag objects in
// cells that are detected must have for the cell to be returned in the check. A tag of "" indicates that any cell that has an object other than the calling
// one is acceptable. Returns a CellCheck object containing this information. Note that the delta movement takes a minimum of 1 pixel movement in any direction
// (a dx and dy of 0 is fine as well). This is to prevent jittering when moving.
func (obj *Object) Check(dx, dy float64, tags ...string) Collision {

	cc := Collision{
		checkingObject: obj,
	}

	if dx < 0 {
		dx = math.Min(dx, -1)
	} else if dx > 0 {
		dx = math.Max(dx, 1)
	}

	if dy < 0 {
		dy = math.Min(dy, -1)
	} else if dy > 0 {
		dy = math.Max(dy, 1)
	}

	cx, cy, ex, ey := obj.ShapeToSpace(dx, dy)

	for y := cy; y <= ey; y++ {

		for x := cx; x <= ex; x++ {

			if c := obj.Space.Cell(x, y); c != nil {

				for _, o := range c.Objects {

					if o == obj {
						continue
					}

					// We only want cells that have objects other than the checking object.

					if len(tags) == 0 || c.ContainsTags(tags...) {

						// We only want to add cells for collisions if the calling object's collision is imprecise OR they are and we aren't intersecting
						if !obj.PreciseCollision || intersecting(obj.X+dx, obj.Y+dy, obj.W, obj.H, o.X, o.Y, o.W, o.H) {
							cc.Cells = append(cc.Cells, c)
							break
						}

					}

					break
				}

			}

		}

	}

	cells := cc.Cells[:]
	ox := cc.checkingObject.X + (cc.checkingObject.W / 2)
	oy := cc.checkingObject.Y + (cc.checkingObject.W / 2)

	sort.Slice(cells, func(i, j int) bool {

		ix, iy := cc.checkingObject.Space.SpaceToWorld(cells[i].X, cells[i].Y)
		jx, jy := cc.checkingObject.Space.SpaceToWorld(cells[j].X, cells[j].Y)

		return distance(ix, iy, ox, oy) < distance(jx, jy, ox, oy)

	})

	cc.calculateContactDelta(dx, dy)
	cc.calculateSlideDelta(dx, dy, tags...)

	return cc

}

func intersecting(x, y, w, h, x2, y2, w2, h2 float64) bool {
	return x > x2-w && y > y2-h && x < x2+w2 && y < y2+h2
}

// func (obj *Object) Intersecting(other *Object) bool {

// 	return obj.X > other.X-obj.W && obj.Y > other.Y-obj.H && obj.X < other.X+other.W && obj.Y < other.Y+other.H

// }
