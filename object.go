package resolv

import (
	"math"
	"sort"

	"github.com/quartercastle/vector"
)

// Object represents an object that can be spread across one or more Cells in a Space. An Object is essentially an AABB (Axis-Aligned Bounding Box) Rectangle.
type Object struct {
	Shape         IShape      // A shape for more specific collision-checking.
	Space         *Space      // Reference to the Space the Object exists within
	X, Y, W, H    float64     // Position and size of the Object in the Space
	TouchingCells []*Cell     // An array of Cells the Object is touching
	Data          interface{} // A pointer to a user-definable object
	ignoreList    []*Object   // Set of Objects to ignore when checking for collisions
	tags          []string    // A list of tags the Object has
}

// NewObject returns a new Object of the specified position and size.
func NewObject(x, y, w, h float64, tags ...string) *Object {
	o := &Object{
		X:    x,
		Y:    y,
		W:    w,
		H:    h,
		tags: []string{},
	}

	if len(tags) > 0 {
		o.AddTags(tags...)
	}

	return o
}

// Clone clones the Object with its properties into another Object. It also clones the Object's Shape (if it has one).
func (obj *Object) Clone() *Object {
	newObj := NewObject(obj.X, obj.Y, obj.W, obj.H, obj.Tags()...)
	newObj.Data = obj.Data
	if obj.Shape != nil {
		newObj.SetShape(obj.Shape.Clone())
	}
	for _, k := range obj.ignoreList {
		newObj.AddToIgnoreList(k)
	}
	return newObj
}

// Update updates the object's association to the Cells in the Space. This should be called whenever an Object is moved.
// This is automatically called once when creating the Object, so you don't have to call it for static objects.
func (obj *Object) Update() {

	if obj.Space != nil {

		// Object.Space.Remove() sets the removed object's Space to nil, indicating it's been removed. Because we're updating
		// the Object (which is essentially removing it from its previous Cells / position and re-adding it to the new Cells /
		// position), we store the original Space to re-set it.

		space := obj.Space

		obj.Space.Remove(obj)

		obj.Space = space

		cx, cy, ex, ey := obj.BoundsToSpace(0, 0)

		for y := cy; y <= ey; y++ {

			for x := cx; x <= ex; x++ {

				c := obj.Space.Cell(x, y)

				if c != nil {
					c.register(obj)
					obj.TouchingCells = append(obj.TouchingCells, c)
				}

			}

		}

	}

	if obj.Shape != nil {
		obj.Shape.SetPosition(obj.X, obj.Y)
	}

}

// AddTags adds tags to the Object.
func (obj *Object) AddTags(tags ...string) {
	obj.tags = append(obj.tags, tags...)
}

// RemoveTags removes tags from the Object.
func (obj *Object) RemoveTags(tags ...string) {

	for _, tag := range tags {

		for i, t := range obj.tags {

			if t == tag {
				obj.tags = append(obj.tags[:i], obj.tags[i+1:]...)
				break
			}

		}

	}

}

// HasTags indicates if an Object has any of the tags indicated.
func (obj *Object) HasTags(tags ...string) bool {

	for _, tag := range tags {

		for _, t := range obj.tags {

			if t == tag {
				return true
			}

		}

	}

	return false

}

// Tags returns the tags an Object has.
func (obj *Object) Tags() []string {
	return append([]string{}, obj.tags...)
}

// SetShape sets the Shape on the Object, in case you need to use precise per-Shape intersection detection. SetShape calls Object.Update() as well, so that it's able to
// update the Shape's position to match its Object as necessary. (If you don't use this, the Shape's position might not match the Object's, depending on if you set the Shape
// after you added the Object to a Space and if you don't call Object.Update() yourself afterwards.)
func (obj *Object) SetShape(shape IShape) {
	if obj.Shape != shape {
		obj.Shape = shape
		obj.Update()
	}
}

// BoundsToSpace returns the Space coordinates of the shape (x, y, w, and h), given its world position and size, and a supposed movement of dx and dy.
func (obj *Object) BoundsToSpace(dx, dy float64) (int, int, int, int) {
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

// Center returns the center position of the Object.
func (obj *Object) Center() (float64, float64) {
	return obj.X + (obj.W / 2.0), obj.Y + (obj.H / 2.0)
}

// SetCenter sets the Object such that its center is at the X and Y position given.
func (obj *Object) SetCenter(x, y float64) {
	obj.X = x - (obj.W / 2)
	obj.Y = y - (obj.H / 2)
}

// CellPosition returns the cellular position of the Object's center in the Space.
func (obj *Object) CellPosition() (int, int) {
	return obj.Space.WorldToSpace(obj.Center())
}

// SetRight sets the X position of the Object so the right edge is at the X position given.
func (obj *Object) SetRight(x float64) {
	obj.X = x - obj.W
}

// SetBottom sets the Y position of the Object so that the bottom edge is at the Y position given.
func (obj *Object) SetBottom(y float64) {
	obj.Y = y - obj.H
}

// Bottom returns the bottom Y coordinate of the Object (i.e. object.Y + object.H).
func (obj *Object) Bottom() float64 {
	return obj.Y + obj.H
}

// Right returns the right X coordinate of the Object (i.e. object.X + object.W).
func (obj *Object) Right() float64 {
	return obj.X + obj.W
}

func (obj *Object) SetBounds(topLeft, bottomRight vector.Vector) {
	obj.X = topLeft[0]
	obj.Y = topLeft[1]
	obj.W = bottomRight[0] - obj.X
	obj.H = bottomRight[1] - obj.Y
}

// Check checks the space around the object using the designated delta movement (dx and dy). This is done by querying the containing Space's Cells
// so that it can see if moving it would coincide with a cell that houses another Object (filtered using the given selection of tag strings). If so,
// Check returns a Collision. If no objects are found or the Object does not exist within a Space, this function returns nil.
func (obj *Object) Check(dx, dy float64, tags ...string) *Collision {

	if obj.Space == nil {
		return nil
	}

	cc := NewCollision()
	cc.checkingObject = obj

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

	cc.dx = dx
	cc.dy = dy

	cx, cy, ex, ey := obj.BoundsToSpace(dx, dy)

	var objectsAdded []*Object
	var cellsAdded []*Cell

	isObjectAdded := func(o *Object) bool {
		for _, o2 := range objectsAdded {
			if o2 == o {
				return true
			}
		}
		return false
	}

	isCellAdded := func(c *Cell) bool {
		for _, c2 := range cellsAdded {
			if c2 == c {
				return true
			}
		}
		return false
	}

	for y := cy; y <= ey; y++ {

		for x := cx; x <= ex; x++ {

			if c := obj.Space.Cell(x, y); c != nil {

				for _, o := range c.Objects {

					// We only want cells that have objects other than the checking object, or that aren't on the ignore list.
					if ignored, _ := obj.InIgnoreList(o); o == obj || ignored {
						continue
					}

					if added := isObjectAdded(o); (len(tags) == 0 || o.HasTags(tags...)) && !added {

						cc.Objects = append(cc.Objects, o)
						objectsAdded = append(objectsAdded, o)
						if added := isCellAdded(c); !added {
							cc.Cells = append(cc.Cells, c)
							cellsAdded = append(cellsAdded, c)
						}
						continue

					}

				}

			}

		}

	}

	if len(cc.Objects) == 0 {
		return nil
	}

	// ox := cc.checkingObject.X + (cc.checkingObject.W / 2)
	// oy := cc.checkingObject.Y + (cc.checkingObject.H / 2)

	ox, oy := cc.checkingObject.Center()
	oc := vector.Vector{ox, oy}

	sort.Slice(cc.Objects, func(i, j int) bool {

		ix, iy := cc.Objects[i].Center()
		jx, jy := cc.Objects[j].Center()
		return vector.Vector{ix, iy}.Sub(oc).Magnitude() < vector.Vector{jx, jy}.Sub(oc).Magnitude()

	})

	cw := cc.checkingObject.Space.CellWidth
	ch := cc.checkingObject.Space.CellHeight

	sort.Slice(cc.Cells, func(i, j int) bool {

		return vector.Vector{float64(cc.Cells[i].X*cw + (cw / 2)), float64(cc.Cells[i].Y*ch + (ch / 2))}.Sub(oc).Magnitude() <
			vector.Vector{float64(cc.Cells[j].X*cw + (cw / 2)), float64(cc.Cells[j].Y*ch + (ch / 2))}.Sub(oc).Magnitude()

	})

	return cc

}

// Overlaps returns if an Object overlaps another Object.
func (obj *Object) Overlaps(other *Object) bool {
	return other.X <= obj.X+obj.W && other.X+other.W >= obj.X && other.Y <= obj.Y+obj.H && other.Y+other.H >= obj.Y
}

// AddToIgnoreList adds the specified Object to the Object's internal collision ignoral list. Cells that contain the specified Object will not be counted when calling Check().
func (obj *Object) AddToIgnoreList(ignoreObj *Object) {
	if exists, _ := obj.InIgnoreList(ignoreObj); !exists {
		obj.ignoreList = append(obj.ignoreList, ignoreObj)
	}
}

// RemoveFromIgnoreList removes the specified Object from the Object's internal collision ignoral list. Objects removed from this list will once again be counted for Check().
func (obj *Object) RemoveFromIgnoreList(ignoreObj *Object) {
	if exists, i := obj.InIgnoreList(ignoreObj); exists {
		obj.ignoreList = append(obj.ignoreList[:i], obj.ignoreList[i+1:]...)
	}
}

// InIgnoreList returns if and where, by index, the specified Object is in the Object's internal collision ignoral list.
func (obj *Object) InIgnoreList(ignoreObj *Object) (bool, int) {
	for i, o := range obj.ignoreList {
		if o == ignoreObj {
			return true, i
		}
	}
	return false, -1
}
