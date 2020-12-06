package resolv

import (
	"math"

	"github.com/kvartborg/vector"
)

// Space represents a collision space. Internally, each Space contains a 2D array of Cells, with each Cell being the same size. Cells contain information on which
// Objects occupy those spaces.
type Space struct {
	Cells                 [][]*Cell
	CellWidth, CellHeight int
	Objects               []*Object
}

// NewSpace creates a new Space, dividing up the provided rectangle of spaceWidthxspaceHeight into an even 2D array grid of Cells, each of the specified width and height.
// Generally, you want cells to be the size of the smallest collide-able objects in your game, and you want to move Objects at a maximum speed of one cell size
// per collision check.
func NewSpace(spaceWidth, spaceHeight, cellWidth, cellHeight int) *Space {

	sp := &Space{
		CellWidth:  cellWidth,
		CellHeight: cellHeight,
		Objects:    []*Object{},
	}

	sp.Resize(spaceWidth, spaceHeight)

	// sp.Resize(int(math.Ceil(float64(spaceWidth)/float64(cellWidth))),
	// 	int(math.Ceil(float64(spaceHeight)/float64(cellHeight))))

	return sp

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

// Cell returns the Cell at the given X and Y position in the Space. If the X and Y position are
// out of bounds, Cell() will return nil.
func (sp *Space) Cell(x, y int) *Cell {

	if y >= 0 && y < len(sp.Cells) && x >= 0 && x < len(sp.Cells[y]) {
		return sp.Cells[y][x]
	}
	return nil

}

// UnregisterAllObjects unregisters all Objects registered to Cells in the Space.
func (sp *Space) UnregisterAllObjects() {

	for y := 0; y < len(sp.Cells); y++ {

		for x := 0; x < len(sp.Cells); x++ {
			cell := sp.Cells[y][x]
			for _, obj := range cell.Objects {
				obj.Remove()
			}
		}

	}

}

// WorldToSpace converts from a world position (x, y) to a position in the Space (a grid-based position).
func (sp *Space) WorldToSpace(x, y float64) (int, int) {
	fx := int(math.Floor(x / float64(sp.CellWidth)))
	fy := int(math.Floor(y / float64(sp.CellHeight)))
	return fx, fy
}

// SpaceToWorld converts from a position in the Space (on a grid) to a world-based position, given the size of the Space when first created.
func (sp *Space) SpaceToWorld(x, y int) (float64, float64) {
	fx := float64(x * sp.CellWidth)
	fy := float64(y * sp.CellHeight)
	return fx, fy
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

		dv := vector.Vector{float64(endX - startX), float64(endY - startY)}
		dv.Unit()
		dv[0] *= float64(sp.CellWidth / 2)
		dv[1] *= float64(sp.CellHeight / 2)

		pX, pY := sp.SpaceToWorld(startX, startY)
		p := vector.Vector{pX + float64(sp.CellWidth/2), pY + float64(sp.CellHeight/2)}

		alternate := false

		for cell != nil {

			if cell == endCell {
				cells = append(cells, cell)
				break
			}

			cells = append(cells, cell)

			if alternate {
				p[1] += dv[1]
			} else {
				p[0] += dv[0]
			}

			cx, cy := sp.WorldToSpace(p[0], p[1])
			c := sp.Cell(cx, cy)
			if c != cell {
				cell = c
			}
			alternate = !alternate

		}

	}

	return cells

}
