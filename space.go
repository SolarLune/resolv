package resolv

import "math"

// Space represents a collision space. Internally, each Space contains a 2D array of Cells, with each Cell being the same size. Cells contain information on which
// Shapes occupy those spaces and are used to speed up intersection testing across multiple Shapes that could be in dynamic locations.
type Space struct {
	cells                 [][]*Cell // The cells present in the Space
	shapes                ShapeCollection
	cellWidth, cellHeight int // Width and Height of each Cell in "world-space" / pixels / whatever
}

// NewSpace creates a new Space. spaceWidth and spaceHeight is the width and height of the Space (usually in pixels), which is then populated with cells of size
// cellWidth by cellHeight. Generally, you want cells to be the size of a "normal object".
// You want to move Objects at a maximum speed of one cell size per collision check to avoid missing any possible collisions.
func NewSpace(spaceWidth, spaceHeight, cellWidth, cellHeight int) *Space {

	sp := &Space{
		cellWidth:  cellWidth,
		cellHeight: cellHeight,
	}

	sp.Resize(int(math.Ceil(float64(spaceWidth)/float64(cellWidth))), int(math.Ceil(float64(spaceHeight)/float64(cellHeight))))

	// sp.Resize(int(math.Ceil(float64(spaceWidth)/float64(cellWidth))),
	// 	int(math.Ceil(float64(spaceHeight)/float64(cellHeight))))

	return sp

}

// Add adds the specified Objects to the Space, updating the Space's cells to refer to the Object.
func (s *Space) Add(shapes ...IShape) {

	for _, shape := range shapes {

		shape.setSpace(s)

		// We call Update() once to make sure the object gets its cells added.
		shape.update()

	}

	s.shapes = append(s.shapes, shapes...)

}

// Remove removes the specified Shapes from the Space.
// This should be done whenever a game object (and its Shape) is removed from the game.
func (s *Space) Remove(shapes ...IShape) {

	for _, shape := range shapes {

		shape.removeFromTouchingCells()

		for i, o := range s.shapes {
			if o == shape {
				s.shapes[i] = nil
				s.shapes = append(s.shapes[:i], s.shapes[i+1:]...)
				break
			}
		}

	}

}

// RemoveAll removes all Shapes from the Space (and from its internal Cells).
func (s *Space) RemoveAll() {

	for i := range s.shapes {
		s.shapes[i] = nil
	}

	s.shapes = s.shapes[:0]
	for y := range s.cells {
		for x := range s.cells[y] {
			for i := range s.cells[y][x].Shapes {
				s.cells[y][x].Shapes[i] = nil
			}
			s.cells[y][x].Shapes = s.cells[y][x].Shapes[:0]
		}
	}

}

// Shapes returns a new slice consisting of all of the shapes present in the Space.
func (s *Space) Shapes() []IShape {
	return append(make([]IShape, 0, len(s.shapes)), s.shapes...)
}

// ForEachShape iterates through each Object in the Space and runs the provided function on them, passing the Shape, its index in the
// Space's shapes slice, and the maximum number of shapes in the space.
// If the function returns false, the iteration ends. If it returns true, it continues.
func (s *Space) ForEachShape(forEach func(object IShape, index, maxCount int) bool) {

	for i, o := range s.shapes {
		if !forEach(o, i, len(s.shapes)) {
			break
		}
	}

}

// FilterShapes returns a ShapeFilter consisting of all shapes present in the Space.
func (s *Space) FilterShapes() ShapeFilter {
	return ShapeFilter{
		operatingOn: s.shapes,
	}
}

// Resize resizes the internal Cells array.
func (s *Space) Resize(width, height int) {

	s.cells = [][]*Cell{}

	for y := 0; y < height; y++ {

		s.cells = append(s.cells, []*Cell{})

		for x := 0; x < width; x++ {
			s.cells[y] = append(s.cells[y], newCell(x, y))
		}

	}

	for _, s := range s.shapes {
		s.update()
	}

}

// Cell returns the Cell at the given cellular / spatial (not world) X and Y position in the Space. If the X and Y position are
// out of bounds, Cell() will return nil. This does not flush shape vicinities beforehand.
func (s *Space) Cell(cx, cy int) *Cell {

	if cy >= 0 && cy < len(s.cells) && cx >= 0 && cx < len(s.cells[cy]) {
		return s.cells[cy][cx]
	}
	return nil

}

// Height returns the height of the Space grid in Cells (so a 320x240 Space with 16x16 cells would have a height of 15).
func (s *Space) Height() int {
	return len(s.cells)
}

// Width returns the width of the Space grid in Cells (so a 320x240 Space with 16x16 cells would have a width of 20).
func (s *Space) Width() int {
	if len(s.cells) > 0 {
		return len(s.cells[0])
	}
	return 0
}

// CellWidth returns the width of each cell in the Space.
func (s *Space) CellWidth() int {
	return s.cellWidth
}

// CellHeight returns the height of each cell in the Space.
func (s *Space) CellHeight() int {
	return s.cellHeight
}

// FilterCells selects a selection of cells.
func (s *Space) FilterCells(bounds Bounds) CellSelection {

	bounds.space = s

	fx, fy, fx2, fy2 := bounds.toCellSpace()

	return CellSelection{
		space:  s,
		StartX: fx,
		StartY: fy,
		EndX:   fx2,
		EndY:   fy2,
	}

}

// func (s *Space) FilterCellsInLine(start, end Vector) CellSelection {

// 	cells := CellSelection{}

// 	startX := int(math.Floor(start.X / float64(s.CellWidth)))
// 	startY := int(math.Floor(start.Y / float64(s.CellHeight)))

// 	endX := int(math.Floor(end.X / float64(s.CellWidth)))
// 	endY := int(math.Floor(end.Y / float64(s.CellHeight)))

// 	cell := s.Cell(startX, startY)
// 	endCell := s.Cell(endX, endY)

// 	if cell != nil && endCell != nil {

// 		dv := Vector{float64(endX - startX), float64(endY - startY)}.Unit()
// 		dv.X *= float64(s.CellWidth / 2)
// 		dv.Y *= float64(s.CellHeight / 2)

// 		pX := float64(startX * s.CellWidth)
// 		pY := float64(startY * s.CellHeight)

// 		p := Vector{pX + float64(s.CellWidth/2), pY + float64(s.CellHeight/2)}

// 		alternate := false

// 		for cell != nil {

// 			if cell == endCell {
// 				cells = append(cells, cell)
// 				break
// 			}

// 			cells = append(cells, cell)

// 			if alternate {
// 				p.Y += dv.Y
// 			} else {
// 				p.X += dv.X
// 			}

// 			cx := int(math.Floor(p.X / float64(s.CellWidth)))
// 			cy := int(math.Floor(p.Y / float64(s.CellHeight)))

// 			c := s.Cell(cx, cy)
// 			if c != cell {
// 				cell = c
// 			}
// 			alternate = !alternate

// 		}

// 	}

// 	return cells

// }
