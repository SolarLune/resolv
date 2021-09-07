package resolv

// Cell is used to contain and organize Object information.
type Cell struct {
	X, Y    int       // The X and Y position of the cell in the Space - note that this is in Grid position, not World position.
	Objects []*Object // The Objects that a Cell contains.
}

// newCell creates a new cell at the specified X and Y position. Should not be used directly.
func newCell(x, y int) *Cell {
	return &Cell{
		X:       x,
		Y:       y,
		Objects: []*Object{},
	}
}

// register registers an object with a Cell. Should not be used directly.
func (cell *Cell) register(obj *Object) {
	if !cell.Contains(obj) {
		cell.Objects = append(cell.Objects, obj)
	}
}

// unregister unregisters an object from a Cell. Should not be used directly.
func (cell *Cell) unregister(obj *Object) {

	for i, o := range cell.Objects {

		if o == obj {
			cell.Objects[i] = cell.Objects[len(cell.Objects)-1]
			cell.Objects = cell.Objects[:len(cell.Objects)-1]
			break
		}

	}

}

// Contains returns whether a Cell contains the specified Object at its position.
func (cell *Cell) Contains(obj *Object) bool {
	for _, o := range cell.Objects {
		if o == obj {
			return true
		}
	}
	return false
}

// ContainsTags returns whether a Cell contains an Object that has the specified tag at its position.
func (cell *Cell) ContainsTags(tags ...string) bool {
	for _, o := range cell.Objects {
		if o.HasTags(tags...) {
			return true
		}
	}
	return false
}

// Occupied returns whether a Cell contains any Objects at all.
func (cell *Cell) Occupied() bool {
	return len(cell.Objects) > 0
}
