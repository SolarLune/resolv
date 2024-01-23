package resolv

import (
	"math"
)

// Space represents a collision space. Internally, each Space contains a 2D array of Cells, with each Cell being the same size. Cells contain information on which
// Objects occupy those spaces.
type Space struct {
	Cells                 [][]*Cell
	CellWidth, CellHeight int // Width and Height of each Cell in "world-space" / pixels / whatever
}

// NewSpace creates a new Space. spaceWidth and spaceHeight is the width and height of the Space (usually in pixels), which is then populated with cells of size
// cellWidth by cellHeight. Generally, you want cells to be the size of the smallest collide-able objects in your game, and you want to move Objects at a maximum
// speed of one cell size per collision check to avoid missing any possible collisions.
func NewSpace(spaceWidth, spaceHeight, cellWidth, cellHeight int) *Space {

	sp := &Space{
		CellWidth:  cellWidth,
		CellHeight: cellHeight,
	}

	sp.Resize(spaceWidth/cellWidth, spaceHeight/cellHeight)

	// sp.Resize(int(math.Ceil(float64(spaceWidth)/float64(cellWidth))),
	// 	int(math.Ceil(float64(spaceHeight)/float64(cellHeight))))

	return sp

}

// Add adds the specified Objects to the Space, updating the Space's cells to refer to the Object.
func (sp *Space) Add(objects ...*Object) {

	if sp == nil {
		panic("ERROR: space is nil")
	}

	for _, obj := range objects {

		obj.Space = sp

		// We call Update() once to make sure the object gets its cells added.
		obj.Update()

	}

}

// Remove removes the specified Objects from being associated with the Space. This should be done whenever an Object is removed from the
// game.
func (sp *Space) Remove(objects ...*Object) {

	if sp == nil {
		panic("ERROR: space is nil")
	}

	for _, obj := range objects {

		for _, cell := range obj.TouchingCells {
			cell.unregister(obj)
		}

		obj.TouchingCells = []*Cell{}

		obj.Space = nil

	}

}

// Objects loops through all Cells in the Space (from top to bottom, and from left to right) to return all Objects
// that exist in the Space. Of course, each Object is counted only once.
func (sp *Space) Objects() []*Object {

	objectsAdded := map[*Object]bool{}
	objects := []*Object{}

	for cy := range sp.Cells {

		for cx := range sp.Cells[cy] {

			for _, o := range sp.Cells[cy][cx].Objects {

				if _, added := objectsAdded[o]; !added {
					objects = append(objects, o)
					objectsAdded[o] = true
				}

			}

		}

	}

	return objects

}

// Resize resizes the internal Cells array.
func (sp *Space) Resize(width, height int) {

	sp.Cells = [][]*Cell{}

	for y := 0; y < height; y++ {

		sp.Cells = append(sp.Cells, []*Cell{})

		for x := 0; x < width; x++ {
			sp.Cells[y] = append(sp.Cells[y], newCell(x, y))
		}

	}

}

// Cell returns the Cell at the given cellular / spatial (not world) X and Y position in the Space. If the X and Y position are
// out of bounds, Cell() will return nil.
func (sp *Space) Cell(x, y int) *Cell {

	if y >= 0 && y < len(sp.Cells) && x >= 0 && x < len(sp.Cells[y]) {
		return sp.Cells[y][x]
	}
	return nil

}

// CheckCells checks a set of cells (from x,y to x + w, y + h in cellular coordinates) and returns
// a slice of the objects found within those Cells.
// The objects must have any of the tags provided (if any are provided).
func (sp *Space) CheckCells(x, y, w, h int, tags ...string) []*Object {

	res := []*Object{}

	for ix := x; ix < x+w; ix++ {

		for iy := y; iy < y+h; iy++ {

			cell := sp.Cell(ix, iy)

			if cell != nil {

				if len(tags) > 0 {

					if cell.ContainsTags(tags...) {
						for _, obj := range cell.Objects {
							if obj.HasTags(tags...) {
								res = append(res, obj)
							}
						}
					}

				} else if cell.Occupied() {
					res = append(res, cell.Objects...)
				}

			}

		}

	}

	return res

}

// CheckWorld checks the cells of the Grid with the given world coordinates.
// Internally, this is just syntactic sugar for calling Space.WorldToSpace() on the
// position and size given.
func (sp *Space) CheckWorld(x, y, w, h float64, tags ...string) []*Object {

	sx, sy := sp.WorldToSpace(x, y)
	cw, ch := sp.WorldToSpace(w, h)

	return sp.CheckCells(sx, sy, cw, ch, tags...)

}

// CheckWorldVec checks the cells of the Grid with the given world coordinates.
// This function takes vectors for the position and size of the checked area.
// Internally, this is just syntactic sugar for calling Space.WorldToSpace() on the
// position and size given.
func (sp *Space) CheckWorldVec(pos, size Vector, tags ...string) []*Object {

	sx, sy := sp.WorldToSpace(pos.X, pos.Y)
	cw, ch := sp.WorldToSpace(size.X, size.Y)

	return sp.CheckCells(sx, sy, cw, ch, tags...)

}

// UnregisterAllObjects unregisters all Objects registered to Cells in the Space.
func (sp *Space) UnregisterAllObjects() {

	for y := 0; y < len(sp.Cells); y++ {

		for x := 0; x < len(sp.Cells[y]); x++ {
			cell := sp.Cells[y][x]
			sp.Remove(cell.Objects...)
		}

	}

}

// WorldToSpace converts from a world position (x, y) to a position in the Space (a grid-based position).
func (sp *Space) WorldToSpace(x, y float64) (int, int) {
	fx := int(math.Floor(x / float64(sp.CellWidth)))
	fy := int(math.Floor(y / float64(sp.CellHeight)))
	return fx, fy
}

// WorldToSpaceVec converts from a world position Vector to a position in the Space (a grid-based position).
func (sp *Space) WorldToSpaceVec(position Vector) (int, int) {
	return sp.WorldToSpace(position.X, position.Y)
}

// SpaceToWorld converts from a position in the Space (on a grid) to a world-based position, given the size of the Space when first created.
func (sp *Space) SpaceToWorld(x, y int) (float64, float64) {
	fx := float64(x * sp.CellWidth)
	fy := float64(y * sp.CellHeight)
	return fx, fy
}

func (sp *Space) SpaceToWorldVec(x, y int) Vector {
	outX, outY := sp.SpaceToWorld(x, y)
	return Vector{outX, outY}
}

// Height returns the height of the Space grid in Cells (so a 320x240 Space with 16x16 cells would have a height of 15).
func (sp *Space) Height() int {
	return len(sp.Cells)
}

// Width returns the width of the Space grid in Cells (so a 320x240 Space with 16x16 cells would have a width of 20).
func (sp *Space) Width() int {
	if len(sp.Cells) > 0 {
		return len(sp.Cells[0])
	}
	return 0
}

func (sp *Space) CellsInLine(startX, startY, endX, endY int) []*Cell {

	cells := []*Cell{}
	cell := sp.Cell(startX, startY)
	endCell := sp.Cell(endX, endY)

	if cell != nil && endCell != nil {

		dv := Vector{float64(endX - startX), float64(endY - startY)}.Unit()
		dv.X *= float64(sp.CellWidth / 2)
		dv.Y *= float64(sp.CellHeight / 2)

		pX, pY := sp.SpaceToWorld(startX, startY)
		p := Vector{pX + float64(sp.CellWidth/2), pY + float64(sp.CellHeight/2)}

		alternate := false

		for cell != nil {

			if cell == endCell {
				cells = append(cells, cell)
				break
			}

			cells = append(cells, cell)

			if alternate {
				p.Y += dv.Y
			} else {
				p.X += dv.X
			}

			cx, cy := sp.WorldToSpace(p.X, p.Y)
			c := sp.Cell(cx, cy)
			if c != cell {
				cell = c
			}
			alternate = !alternate

		}

	}

	return cells

}
