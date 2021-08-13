package resolv

import (
	"math"

	"github.com/kvartborg/vector"
)

type Shape interface {
	Intersection(other Shape) *Delta
	SetPosition(x, y float64)
}

type Line struct {
	Start, End vector.Vector
}

func NewLine(x, y, x2, y2 float64) *Line {
	return &Line{
		Start: vector.Vector{x, y},
		End:   vector.Vector{x2, y2},
	}
}

func (line *Line) Project(axis vector.Vector) vector.Vector {
	return line.Vector().Scale(axis.Dot(line.Start.Sub(line.End)))
}

func (line *Line) Normal() vector.Vector {
	v := line.Vector()
	return vector.Vector{v[1], -v[0]}.Unit()
}

func (line *Line) Vector() vector.Vector {
	return line.End.Clone().Sub(line.Start).Unit()
}

type ConvexPolygon struct {
	Vertices []*Vertex
	X, Y     float64
}

// NewConvexPolygon creates a new convex polygon from the provided set of 2D points. Should be ordered clockwise.
func NewConvexPolygon(vertexPositions ...float64) *ConvexPolygon {

	if len(vertexPositions)/2 < 3 {
		return nil
	}

	cp := &ConvexPolygon{Vertices: []*Vertex{}}

	for v := 0; v < len(vertexPositions); v += 2 {
		cp.AddPoints(NewVertex(vertexPositions[v], vertexPositions[v+1]))
	}

	return cp
}

func (cp *ConvexPolygon) AddPoints(points ...*Vertex) {
	cp.Vertices = append(cp.Vertices, points...)
}

func (cp *ConvexPolygon) Lines() []*Line {

	lines := []*Line{}

	vertices := cp.Transformed()

	for i := 0; i < len(vertices); i++ {

		start, end := vertices[i], vertices[0]

		if i < len(vertices)-1 {
			end = vertices[i+1]
		}

		line := NewLine(start.X, start.Y, end.X, end.Y)

		lines = append(lines, line)

	}

	return lines

}

func (cp *ConvexPolygon) Transformed() []*Vertex {
	transformed := []*Vertex{}
	for _, point := range cp.Vertices {
		transformed = append(transformed, NewVertex(point.X+cp.X, point.Y+cp.Y))
	}
	return transformed
}

func (cp *ConvexPolygon) SetPosition(x, y float64) {
	cp.X = x
	cp.Y = y
}

func (cp *ConvexPolygon) Project(axis vector.Vector) Projection {
	axis = axis.Unit()
	vertices := cp.Transformed()
	min := axis.Dot(vertices[0].Vector())
	max := min
	for i := 1; i < len(vertices); i++ {
		p := axis.Dot(vertices[i].Vector())
		if p < min {
			min = p
		} else if p > max {
			max = p
		}
	}
	return Projection{min, max}
}

func (cp *ConvexPolygon) SATAxes() []vector.Vector {

	axes := []vector.Vector{}
	for _, line := range cp.Lines() {
		axes = append(axes, line.Normal())
	}
	return axes

}

// func (cp *ConvexPolygon) AddLine(line *Line) {
// 	cp.Lines = append(cp.Lines, line)
// }

func (cp *ConvexPolygon) Intersection(otherShape Shape) *Delta {

	delta := &Delta{}

	smallest := vector.Vector{math.MaxFloat64, 0}

	switch other := otherShape.(type) {

	case *ConvexPolygon:

		for _, axis := range cp.SATAxes() {
			if !cp.Project(axis).IsOverlapping(other.Project(axis)) {
				return nil
			}

			overlap := cp.Project(axis).Overlap(other.Project(axis))

			if smallest.Magnitude() > overlap {
				smallest = axis.Scale(overlap)
			}

		}

		for _, axis := range other.SATAxes() {

			if !cp.Project(axis).IsOverlapping(other.Project(axis)) {
				return nil
			}

			overlap := cp.Project(axis).Overlap(other.Project(axis))

			if smallest.Magnitude() > overlap {
				smallest = axis.Scale(overlap)
			}

		}

	}

	delta.X = smallest[0]
	delta.Y = smallest[1]

	return delta
}

func (cp *ConvexPolygon) ContainedBy(otherShape Shape) *Delta {

	delta := &Delta{}
	switch other := otherShape.(type) {

	case *ConvexPolygon:

		for _, axis := range cp.SATAxes() {
			if !cp.Project(axis).IsInside(other.Project(axis)) {
				return nil
			}
		}

		for _, axis := range other.SATAxes() {
			if !cp.Project(axis).IsInside(other.Project(axis)) {
				return nil
			}
		}

	}

	return delta
}

func (cp *ConvexPolygon) FlipH() {

	for _, v := range cp.Vertices {
		v.X = -v.X
	}
	// We have to reverse vertex order after flipping the vertices to ensure the winding order is consistent between Objects (so that the normals are consistently outside or inside, which is important
	// when doing Intersection tests). If we assume that the normal of a line, going from vertex A to vertex B, is one direction, then the normal would be inverted if the vertices were flipped in position,
	// but not in order. This would make Intersection tests drive objects into each other, instead of giving the delta to move away.
	cp.ReverseVertexOrder()

}

func (cp *ConvexPolygon) FlipV() {

	for _, v := range cp.Vertices {
		v.Y = -v.Y
	}
	cp.ReverseVertexOrder()

}

func (cp *ConvexPolygon) ReverseVertexOrder() {

	verts := []*Vertex{cp.Vertices[0]}

	for i := len(cp.Vertices) - 1; i >= 1; i-- {
		verts = append(verts, cp.Vertices[i])
	}

	cp.Vertices = verts

}

// NewRectangle just returns a rectangular ConvexPolygon, for simplicity, for now. The vertices are in clockwise order. In actuality, an AABBRectangle should be its own
// "thing" with its own optimized Intersection code check.
func NewRectangle(x, y, w, h float64) *ConvexPolygon {
	return NewConvexPolygon(
		x, y,
		x+w, y,
		x+w, y+h,
		x, y+h,
	)
}

// type Rectangle struct {
// 	X, Y, W, H float64
// }

// func NewRectangle(x, y, w, h float64) *Rectangle {
// 	return &Rectangle{x, y, w, h}
// }

// func (rect *Rectangle) Clone() *Rectangle {
// 	return NewRectangle(rect.X, rect.Y, rect.W, rect.H)
// }

// func (rect *Rectangle) Intersection(other Shape) Delta {

// 	delta := NewDelta()

// 	switch o := other.(type) {

// 	case *Rectangle:
// 		delta.Valid = rect.X > o.X-rect.W && rect.Y > o.Y-rect.H && rect.X < o.X+o.W && rect.Y < o.Y+o.H

// 		// case *Point:
// 		// 	delta.Valid = o.X >= rect.X && o.X <= rect.X+rect.W && o.Y >= rect.Y && o.Y <= rect.Y+rect.H

// 	}

// 	return delta

// }

type Vertex struct {
	X, Y float64
}

func NewVertex(x, y float64) *Vertex {
	return &Vertex{x, y}
}

func (vert *Vertex) Intersection(other Shape) *Delta {

	delta := &Delta{}

	switch o := other.(type) {
	// case *Rectangle:
	// 	return o.Intersection(point)
	case *Vertex:
		if vert.X != o.X || vert.Y != o.Y {
			return nil
		}
	}

	return delta

}

func (vert *Vertex) Vector() vector.Vector {
	return vector.Vector{vert.X, vert.Y}
}

func (vert *Vertex) SetPosition(x, y float64) {
	vert.X = x
	vert.Y = y
}

// type Circle struct {
// 	X, Y, Radius float64
// }

// func NewCircle(x, y, radius float64) *Circle {
// 	return &Circle{x, y, radius}
// }

// func (circle *Circle) Intersecting(other Shape) bool {

// }

type Projection struct {
	Min, Max float64
}

// Credit to https://www.sevenson.com.au/programming/sat/
func (projection Projection) IsOverlapping(other Projection) bool {
	return projection.Overlap(other) > 0
}

// Credit to https://dyn4j.org/2010/01/sat/#sat-nointer
func (projection Projection) Overlap(other Projection) float64 {
	return math.Min(projection.Max, other.Max) - math.Max(projection.Min, other.Min)
}

func (projection Projection) IsInside(other Projection) bool {
	return projection.Min >= other.Min && projection.Max <= other.Max
}
