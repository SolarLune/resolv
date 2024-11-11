package resolv

// Cell is used to contain and organize Object information.
type Cell struct {
	X, Y   int
	Shapes []IShape // The Objects that a Cell contains.
}

// newCell creates a new cell at the specified X and Y position. Should not be used directly.
func newCell(x, y int) *Cell {
	return &Cell{
		X:      x,
		Y:      y,
		Shapes: []IShape{},
	}
}

// register registers an object with a Cell. Should not be used directly.
func (cell *Cell) register(obj IShape) {
	if !cell.Contains(obj) {
		cell.Shapes = append(cell.Shapes, obj)
	}
}

// unregister unregisters an object from a Cell. Should not be used directly.
func (cell *Cell) unregister(obj IShape) {

	for i, o := range cell.Shapes {

		if o == obj {
			cell.Shapes[i] = cell.Shapes[len(cell.Shapes)-1]
			cell.Shapes = cell.Shapes[:len(cell.Shapes)-1]
			break
		}

	}

}

// Contains returns whether a Cell contains the specified Object at its position.
func (cell *Cell) Contains(obj IShape) bool {
	for _, o := range cell.Shapes {
		if o == obj {
			return true
		}
	}
	return false
}

// ContainsTags returns whether a Cell contains an Object that has the specified tag at its position.
func (cell *Cell) HasTags(tags Tags) bool {
	for _, o := range cell.Shapes {
		if o.Tags().Has(tags) {
			return true
		}
	}
	return false
}

// IsOccupied returns whether a Cell contains any Objects at all.
func (cell *Cell) IsOccupied() bool {
	return len(cell.Shapes) > 0
}

// CellSelection is a selection of cells. It is primarily used to filter down Shapes.
type CellSelection struct {
	StartX, StartY, EndX, EndY int // The start and end position of the Cell in cellular locations.
	space                      *Space
	excludeSelf                IShape
}

// FilterShapes returns a ShapeFilter of the shapes within the cell selection.
func (c CellSelection) FilterShapes() ShapeFilter {

	if c.space == nil {
		return ShapeFilter{}
	}

	return ShapeFilter{
		operatingOn: c,
	}

}

// ForEach loops through each shape in the CellSelection.
func (c CellSelection) ForEach(iterationFunction func(shape IShape) bool) {
	// Internally, this function allows us to pass a CellSelection as the operatingOn property in a ShapeFilter.

	cellSelectionForEachIDSet = cellSelectionForEachIDSet[:0]

	for y := c.StartY; y <= c.EndY; y++ {

		for x := c.StartX; x <= c.EndX; x++ {

			cell := c.space.Cell(x, y)

			if cell != nil {

				for _, s := range cell.Shapes {

					if s == c.excludeSelf {
						continue
					}

					if cellSelectionForEachIDSet.idInSet(s.ID()) {
						continue
					}
					if !iterationFunction(s) {
						break
					}
					cellSelectionForEachIDSet = append(cellSelectionForEachIDSet, s.ID())

				}

			}

		}

	}

}
