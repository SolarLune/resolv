package resolv

import (
	"errors"
	"log"
	"math"
	"sort"
)

// ConvexPolygon represents a series of points, connected by lines, constructing a convex shape.
// The polygon has a position, a scale, a rotation, and may or may not be closed.
type ConvexPolygon struct {
	ShapeBase

	scale    Vector
	rotation float64  // How many radians the ConvexPolygon is rotated around in the viewing vector (Z).
	Points   []Vector // Points represents the points constructing the ConvexPolygon.
	Closed   bool     // Closed is whether the ConvexPolygon is closed or not; only takes effect if there are more than 2 points.
	bounds   Bounds
}

// NewConvexPolygon creates a new convex polygon at the position given, from the provided set of X and Y positions of 2D points (or vertices).
// You don't need to pass any points at this stage, but if you do, you should pass whole pairs. The points should generally be ordered clockwise,
// from X and Y of the first, to X and Y of the last.
// For example: NewConvexPolygon(30, 20, 0, 0, 10, 0, 10, 10, 0, 10) would create a 10x10 convex
// polygon square, with the vertices at {0,0}, {10,0}, {10, 10}, and {0, 10}, with the polygon itself occupying a position of 30, 20.
// You can also pass the points using vectors with ConvexPolygon.AddPointsVec().
func NewConvexPolygon(x, y float64, points []float64) *ConvexPolygon {

	cp := &ConvexPolygon{
		ShapeBase: newShapeBase(x, y),
		scale:     NewVector(1, 1),
		Points:    []Vector{},
		Closed:    true,
	}

	cp.owner = cp

	if len(points) > 0 {
		err := cp.AddPoints(points...)
		if err != nil {
			log.Println(err)
		}
	}

	return cp
}

func NewConvexPolygonVec(position Vector, points []Vector) *ConvexPolygon {

	cp := &ConvexPolygon{
		ShapeBase: newShapeBase(position.X, position.Y),
		scale:     NewVector(1, 1),
		Points:    []Vector{},
		Closed:    true,
	}

	cp.owner = cp

	if len(points) > 0 {
		cp.AddPointsVec(points...)
	}

	return cp

}

// Clone returns a clone of the ConvexPolygon as an IShape.
func (cp *ConvexPolygon) Clone() IShape {

	points := append(make([]Vector, 0, len(cp.Points)), cp.Points...)

	newPoly := NewConvexPolygonVec(cp.position, points)
	newPoly.tags.Set(*cp.tags)

	newPoly.ShapeBase = cp.ShapeBase
	newPoly.id = globalShapeID
	globalShapeID++
	newPoly.ShapeBase.space = nil
	newPoly.ShapeBase.touchingCells = []*Cell{}
	newPoly.ShapeBase.owner = newPoly

	newPoly.rotation = cp.rotation
	newPoly.scale = cp.scale
	newPoly.Closed = cp.Closed

	return newPoly
}

// AddPoints allows you to add points to the ConvexPolygon with a slice or selection of float64s, with each pair indicating an X or Y value for
// a point / vertex (i.e. AddPoints(0, 1, 2, 3) would add two points - one at {0, 1}, and another at {2, 3}).
func (cp *ConvexPolygon) AddPoints(vertexPositions ...float64) error {
	if len(vertexPositions) < 4 {
		return errors.New("addpoints called with not enough passed vertex positions")
	}
	if len(vertexPositions)%2 == 1 {
		return errors.New("addpoints called with a non-even amount of vertex positions")
	}
	for v := 0; v < len(vertexPositions); v += 2 {
		cp.Points = append(cp.Points, Vector{vertexPositions[v], vertexPositions[v+1]})
	}

	// Call updateBounds first so that the bounds are updated to determine cellular location
	cp.updateBounds()
	cp.update()
	return nil
}

// AddPointsVec allows you to add points to the ConvexPolygon with a slice of Vectors, each indicating a point / vertex.
func (cp *ConvexPolygon) AddPointsVec(points ...Vector) {
	cp.Points = append(cp.Points, points...)
	cp.updateBounds()
	cp.update()
}

// Lines returns a slice of transformed internalLines composing the ConvexPolygon.
func (cp *ConvexPolygon) Lines() []collidingLine {

	lines := []collidingLine{}

	vertices := cp.Transformed()

	for i := 0; i < len(vertices); i++ {

		start, end := vertices[i], vertices[0]

		if i < len(vertices)-1 {
			end = vertices[i+1]
		} else if !cp.Closed || len(cp.Points) <= 2 {
			break
		}

		line := newCollidingLine(start.X, start.Y, end.X, end.Y)

		lines = append(lines, line)

	}

	return lines

}

// Transformed returns the ConvexPolygon's points / vertices, transformed according to the ConvexPolygon's position.
func (cp *ConvexPolygon) Transformed() []Vector {
	transformed := []Vector{}
	for _, point := range cp.Points {
		p := Vector{point.X * cp.scale.X, point.Y * cp.scale.Y}
		if cp.rotation != 0 {
			p = p.Rotate(-cp.rotation)
		}
		transformed = append(transformed, Vector{p.X + cp.position.X, p.Y + cp.position.Y})
	}
	return transformed
}

// Bounds returns two Vectors, comprising the top-left and bottom-right positions of the bounds of the
// ConvexPolygon, post-transformation.
func (cp *ConvexPolygon) Bounds() Bounds {
	cp.bounds.space = cp.space
	return cp.bounds.MoveVec(cp.position)
}

func (cp *ConvexPolygon) updateBounds() {

	transformed := cp.Transformed()

	topLeft := Vector{transformed[0].X, transformed[0].Y}
	bottomRight := topLeft

	for i := 0; i < len(transformed); i++ {

		point := transformed[i]

		if point.X < topLeft.X {
			topLeft.X = point.X
		} else if point.X > bottomRight.X {
			bottomRight.X = point.X
		}

		if point.Y < topLeft.Y {
			topLeft.Y = point.Y
		} else if point.Y > bottomRight.Y {
			bottomRight.Y = point.Y
		}

	}

	cp.bounds = Bounds{
		Min:   topLeft,
		Max:   bottomRight,
		space: cp.space,
	}

	// Untransform those points so that we don't have to update it whenever it moves
	cp.bounds = cp.bounds.Move(-cp.position.X, -cp.position.Y)

}

// Center returns the transformed Center of the ConvexPolygon.
func (cp *ConvexPolygon) Center() Vector {

	// pos := Vector{0, 0}

	// for _, v := range cp.Transformed() {
	// 	pos = pos.Add(v)
	// }

	// pos.X /= float64(len(cp.Transformed()))
	// pos.Y /= float64(len(cp.Transformed()))

	// return pos

	return cp.Bounds().Center()

}

// Project projects (i.e. flattens) the ConvexPolygon onto the provided axis.
func (cp *ConvexPolygon) Project(axis Vector) Projection {
	axis = axis.Unit()
	vertices := cp.Transformed()
	min := axis.Dot(vertices[0])
	max := min
	for i := 1; i < len(vertices); i++ {
		p := axis.Dot(vertices[i])
		if p < min {
			min = p
		} else if p > max {
			max = p
		}
	}
	return Projection{min, max}
}

// SATAxes returns the axes of the ConvexPolygon for SAT intersection testing.
func (cp *ConvexPolygon) SATAxes() []Vector {

	axes := []Vector{}
	for _, line := range cp.Lines() {
		axes = append(axes, line.Normal())
	}
	return axes

}

// Rotation returns the rotation (in radians) of the ConvexPolygon.
func (p *ConvexPolygon) Rotation() float64 {
	return p.rotation
}

// SetRotation sets the rotation for the ConvexPolygon; note that the rotation goes counter-clockwise from 0 to pi, and then from -pi at 180 down, back to 0.
// This rotation scheme follows the way math.Atan2() works.
func (p *ConvexPolygon) SetRotation(radians float64) {
	p.rotation = radians
	if p.rotation > math.Pi {
		p.rotation -= math.Pi * 2
	} else if p.rotation < -math.Pi {
		p.rotation += math.Pi * 2
	}
	p.updateBounds()
	p.update()
}

// Rotate is a helper function to rotate a ConvexPolygon by the radians given.
func (p *ConvexPolygon) Rotate(radians float64) {
	p.SetRotation(p.Rotation() + radians)
}

// Scale returns the scale multipliers of the ConvexPolygon.
func (p *ConvexPolygon) Scale() Vector {
	return p.scale
}

// SetScale sets the scale multipliers of the ConvexPolygon.
func (p *ConvexPolygon) SetScale(x, y float64) {
	p.scale.X = x
	p.scale.Y = y
	p.updateBounds()
	p.update()
}

// SetScaleVec sets the scale multipliers of the ConvexPolygon using the provided Vector.
func (p *ConvexPolygon) SetScaleVec(vec Vector) {
	p.SetScale(vec.X, vec.Y)
}

// Intersection returns an IntersectionSet for the other Shape provided.
// If no intersection is detected, the IntersectionSet returned is empty.
func (p *ConvexPolygon) Intersection(other IShape) IntersectionSet {

	switch otherShape := other.(type) {
	case *ConvexPolygon:
		return convexConvexTest(p, otherShape)
	case *Circle:
		return convexCircleTest(p, otherShape)
	}

	// This should never happen
	panic("Unimplemented intersection")

}

// ShapeLineTestSettings is a struct of settings to be used when performing shape line tests
// (the equivalent of 3D hitscan ray tests for 2D, but emitted from each vertex of the Shape).
type ShapeLineTestSettings struct {
	StartOffset Vector        // An offset to use for casting rays from each vertex of the Shape.
	Vector      Vector        // The direction and distance vector to use for casting the lines.
	TestAgainst ShapeIterator // The shapes to test against.
	// OnIntersect is the callback to be called for each intersection between a line from the given Shape, ranging from its origin off towards the given Vector against each shape given in TestAgainst.
	// set is the intersection set that contains information about the intersection, index is the index of the current intersection out of the max number of intersections,
	// and count is the total number of intersections detected from the intersection test.
	// The boolean the callback returns indicates whether the line test should continue iterating through results or stop at the currently found intersection.
	OnIntersect      func(set IntersectionSet, index, count int) bool
	IncludeAllPoints bool  // Whether to cast lines from all points in the Shape (true), or just points from the leading edges (false, and the default). Only takes effect for ConvexPolygons.
	Lines            []int // Which line indices to cast from. If unset (which is the default), then all vertices from all lines will be used.
}

var lineTestResults []IntersectionSet
var lineTestVertices = newSet[Vector]()

// ShapeLineTest conducts a line test from each vertex of the ConvexPolygon using the settings passed.
// By default, lines are cast from each vertex of each leading edge in the ConvexPolygon.
func (cp *ConvexPolygon) ShapeLineTest(settings ShapeLineTestSettings) bool {

	lineTestResults = lineTestResults[:0]

	lineTestVertices.Clear()

	// We only have to test vertices from the leading lines, not all of them
	if !settings.IncludeAllPoints {

		for i, l := range cp.Lines() {

			found := true
			if len(settings.Lines) > 0 {
				found = false
				for lineIndex := range settings.Lines {
					if i == lineIndex {
						found = true
						break
					}
				}
			}

			if found {

				// If a line's normal points away from the checking vector, it isn't a leading edge
				if l.Normal().Dot(settings.Vector) < 0.01 {
					continue
				}

				// Kick the vertices in along the lines a bit to ensure they don't get snagged up on borders
				v := l.Vector().Scale(0.5).Invert()
				lineTestVertices.Add(l.Start.Sub(v))
				lineTestVertices.Add(l.End.Sub(v.Invert()))

			}

		}

	} else {
		lineTestVertices.Add(cp.Transformed()...)
	}

	vu := settings.Vector.Unit()

	for p := range lineTestVertices {

		start := p.Sub(vu.Add(settings.StartOffset))

		LineTest(LineTestSettings{
			Start:        start,
			End:          p.Add(settings.Vector),
			TestAgainst:  settings.TestAgainst,
			callingShape: cp,
			OnIntersect: func(set IntersectionSet, index, max int) bool {

				// Consolidate hits together across multiple objects
				for i := range lineTestResults {
					if lineTestResults[i].OtherShape == set.OtherShape {
						lineTestResults[i].Intersections = append(lineTestResults[i].Intersections, set.Intersections...)
						if set.MTV.MagnitudeSquared() < lineTestResults[i].MTV.MagnitudeSquared() {
							lineTestResults[i].MTV = set.MTV
						}
						return true
					}
				}

				lineTestResults = append(lineTestResults, set)

				return true

			},
		})

	}

	// Sort the results by smallest MTV because we can't really easily get the starting points of the ray test results
	sort.Slice(lineTestResults, func(i, j int) bool {
		return lineTestResults[i].MTV.MagnitudeSquared() < lineTestResults[j].MTV.MagnitudeSquared()
	})

	if settings.OnIntersect != nil {

		for i := range lineTestResults {

			if !settings.OnIntersect(lineTestResults[i], i, len(lineTestResults)) {
				break
			}

		}

	}

	return len(lineTestResults) > 0
}

// calculateMTV returns the MTV, if possible, and a bool indicating whether it was possible or not.
func (cp *ConvexPolygon) calculateMTV(otherShape IShape) (Vector, bool) {

	delta := Vector{0, 0}

	smallest := Vector{math.MaxFloat64, 0}

	switch other := otherShape.(type) {

	case *ConvexPolygon:

		for _, axis := range cp.SATAxes() {
			pa := cp.Project(axis)
			pb := other.Project(axis)

			overlap := pa.Overlap(pb)

			if overlap <= 0 {
				return Vector{}, false
			}

			if smallest.Magnitude() > overlap {
				smallest = axis.Scale(overlap)
			}

		}

		for _, axis := range other.SATAxes() {

			pa := cp.Project(axis)
			pb := other.Project(axis)

			overlap := pa.Overlap(pb)

			if overlap <= 0 {
				return Vector{}, false
			}

			if smallest.Magnitude() > overlap {
				smallest = axis.Scale(overlap)
			}

		}

		// If the direction from target to source points opposite to the separation, invert the separation vector.
		if cp.Center().Sub(other.Center()).Dot(smallest) < 0 {
			smallest = smallest.Invert()
		}

	case *Circle:

		verts := append([]Vector{}, cp.Transformed()...)
		center := other.position
		sort.Slice(verts, func(i, j int) bool { return verts[i].Sub(center).Magnitude() < verts[j].Sub(center).Magnitude() })

		axis := Vector{center.X - verts[0].X, center.Y - verts[0].Y}

		pa := cp.Project(axis)
		pb := other.Project(axis)
		overlap := pa.Overlap(pb)
		if overlap <= 0 {
			return Vector{}, false
		}
		smallest = axis.Unit().Scale(overlap)

		for _, axis := range cp.SATAxes() {
			pa := cp.Project(axis)
			pb := other.Project(axis)

			overlap := pa.Overlap(pb)

			if overlap <= 0 {
				return Vector{}, false
			}

			if smallest.Magnitude() > overlap {
				smallest = axis.Scale(overlap)
			}

		}

		// If the direction from target to source points opposite to the separation, invert the separation vector
		if cp.Center().Sub(other.position).Dot(smallest) < 0 {
			smallest = smallest.Invert()
		}

	}

	delta.X = smallest.X
	delta.Y = smallest.Y

	pointingDirection := otherShape.Position().Sub(cp.Position())
	if pointingDirection.Dot(delta) > 0 {
		delta = delta.Invert()
	}

	return delta, true
}

// IsContainedBy returns if the ConvexPolygon is wholly contained by the other shape provided.
// Note that only testing against ConvexPolygons is implemented currently.
func (cp *ConvexPolygon) IsContainedBy(otherShape IShape) bool {

	switch other := otherShape.(type) {

	case *ConvexPolygon:

		for _, axis := range cp.SATAxes() {
			if !cp.Project(axis).IsInside(other.Project(axis)) {
				return false
			}
		}

		for _, axis := range other.SATAxes() {
			if !cp.Project(axis).IsInside(other.Project(axis)) {
				return false
			}
		}

		// TODO: Implement this for Circles

	}

	return true
}

// FlipH flips the ConvexPolygon's vertices horizontally, across the polygon's width, according to their initial offset when adding the points.
func (cp *ConvexPolygon) FlipH() {

	for _, v := range cp.Points {
		v.X = -v.X
	}
	// We have to reverse vertex order after flipping the vertices to ensure the winding order is consistent between
	// Objects (so that the normals are consistently outside or inside, which is important when doing Intersection tests).
	// If we assume that the normal of a line, going from vertex A to vertex B, is one direction, then the normal would be
	// inverted if the vertices were flipped in position, but not in order. This would make Intersection tests drive objects
	// into each other, instead of giving the delta to move away.
	cp.ReverseVertexOrder()
	cp.updateBounds()
	cp.update()

}

// FlipV flips the ConvexPolygon's vertices vertically according to their initial offset when adding the points.
func (cp *ConvexPolygon) FlipV() {

	for _, v := range cp.Points {
		v.Y = -v.Y
	}
	cp.ReverseVertexOrder()
	cp.updateBounds()
	cp.update()

}

// RecenterPoints recenters the vertices in the polygon, such that they are all equidistant from the center.
// For example, say you had a polygon with the following three points: {0, 0}, {10, 0}, {0, 16}.
// After calling cp.RecenterPoints(), the polygon's points would be at {-5, -8}, {5, -8}, {-5, 8}.
func (cp *ConvexPolygon) RecenterPoints() {

	if len(cp.Points) <= 1 {
		return
	}

	offset := Vector{0, 0}
	for _, p := range cp.Points {
		offset = offset.Add(p)
	}

	offset = offset.Scale(1.0 / float64(len(cp.Points))).Invert()

	for i := range cp.Points {
		cp.Points[i] = cp.Points[i].Add(offset)
	}

	cp.position = cp.position.Sub(offset)

	cp.updateBounds()

	cp.update()

}

// ReverseVertexOrder reverses the vertex ordering of the ConvexPolygon.
func (cp *ConvexPolygon) ReverseVertexOrder() {

	verts := []Vector{cp.Points[0]}

	for i := len(cp.Points) - 1; i >= 1; i-- {
		verts = append(verts, cp.Points[i])
	}

	cp.Points = verts

}

// NewRectangle returns a rectangular ConvexPolygon at the position given with the vertices ordered in clockwise order.
// The Rectangle's origin will be the center of its shape (as is recommended for collision testing).
// The {x, y} is the center position of the Rectangle).
func NewRectangle(x, y, w, h float64) *ConvexPolygon {
	// TODO: In actuality, an AABBRectangle should be its own "thing" with its own optimized Intersection code check.

	hw := w / 2
	hh := h / 2

	return NewConvexPolygon(
		x, y,

		[]float64{
			-hw, -hh,
			hw, -hh,
			hw, hh,
			-hw, hh,
		},
	)
}

// NewRectangleTopLeft returns a rectangular ConvexPolygon at the position given with the vertices ordered in clockwise order.
// The Rectangle's origin will be the center of its shape (as is recommended for collision testing).
// Note that the rectangle will be positioned such that x, y is the top-left corner, though the center-point is still
// in the center of the ConvexPolygon shape.
func NewRectangleTopLeft(x, y, w, h float64) *ConvexPolygon {

	r := NewRectangle(x, y, w, h)
	r.Move(w/2, h/2)
	return r
}

// NewRectangleFromCorners returns a rectangluar ConvexPolygon properly centered with its corners at the given { x1, y1 } and { x2, y2 } coordinates.
// The Rectangle's origin will be the center of its shape (as is recommended for collision testing).
func NewRectangleFromCorners(x1, y1, x2, y2 float64) *ConvexPolygon {

	if x2 < x2 {
		x1, x2 = x2, x1
	}
	if y2 < y2 {
		y1, y2 = y2, y1
	}

	halfWidth := (x2 - x1) / 2
	halfHeight := (y2 - y1) / 2

	return NewConvexPolygon(
		x1+halfWidth, y1+halfHeight,

		[]float64{
			-halfWidth, -halfHeight,
			halfWidth, -halfHeight,
			halfWidth, halfHeight,
			-halfWidth, halfHeight,
		},
	)
}

func NewLine(x1, y1, x2, y2 float64) *ConvexPolygon {

	cx := x1 + ((x2 - x1) / 2)
	cy := y1 + ((y2 - y1) / 2)

	return NewConvexPolygon(
		cx, cy,

		[]float64{
			cx - x1, cy - y1,
			cx - x2, cy - y2,
		},
	)
}

/////

// A collidingLine is a helper shape used to determine if two ConvexPolygon lines intersect; you can't create a collidingLine to use as a Shape.
// Instead, you can create a ConvexPolygon, specify two points, and set its Closed value to false (or use NewLine(), as this does it for you).
type collidingLine struct {
	Start, End Vector
}

func newCollidingLine(x, y, x2, y2 float64) collidingLine {
	return collidingLine{
		Start: Vector{x, y},
		End:   Vector{x2, y2},
	}
}

func (line collidingLine) Project(axis Vector) Vector {
	return line.Vector().Scale(axis.Dot(line.Start.Sub(line.End)))
}

func (line collidingLine) Normal() Vector {
	v := line.Vector()
	return Vector{v.Y, -v.X}.Unit()
}

func (line collidingLine) Vector() Vector {
	return line.End.Sub(line.Start).Unit()
}

// IntersectionPointsLine returns the intersection point of a Line with another Line as a Vector, and if the intersection was found.
func (line collidingLine) IntersectionPointsLine(other collidingLine) (Vector, bool) {

	det := (line.End.X-line.Start.X)*(other.End.Y-other.Start.Y) - (other.End.X-other.Start.X)*(line.End.Y-line.Start.Y)

	if det != 0 {

		// MAGIC MATH; the extra + 1 here makes it so that corner cases (literally, lines going through corners) works.

		// lambda := (float32(((line.Y-b.Y)*(b.X2-b.X))-((line.X-b.X)*(b.Y2-b.Y))) + 1) / float32(det)
		lambda := (((line.Start.Y - other.Start.Y) * (other.End.X - other.Start.X)) - ((line.Start.X - other.Start.X) * (other.End.Y - other.Start.Y)) + 1) / det

		// gamma := (float32(((line.Y-b.Y)*(line.X2-line.X))-((line.X-b.X)*(line.Y2-line.Y))) + 1) / float32(det)
		gamma := (((line.Start.Y - other.Start.Y) * (line.End.X - line.Start.X)) - ((line.Start.X - other.Start.X) * (line.End.Y - line.Start.Y)) + 1) / det

		if (0 < lambda && lambda < 1) && (0 < gamma && gamma < 1) {

			// Delta
			dx := line.End.X - line.Start.X
			dy := line.End.Y - line.Start.Y

			// dx, dy := line.GetDelta()

			return Vector{line.Start.X + (lambda * dx), line.Start.Y + (lambda * dy)}, true
		}

	}

	return Vector{}, false

}

// IntersectionPointsCircle returns a slice of Vectors, each indicating the intersection point. If no intersection is found, it will return an empty slice.
func (line collidingLine) IntersectionPointsCircle(circle *Circle) []Vector {

	points := []Vector{}

	cp := circle.position
	lStart := line.Start.Sub(cp)
	lEnd := line.End.Sub(cp)
	diff := lEnd.Sub(lStart)

	a := diff.X*diff.X + diff.Y*diff.Y
	b := 2 * ((diff.X * lStart.X) + (diff.Y * lStart.Y))
	c := (lStart.X * lStart.X) + (lStart.Y * lStart.Y) - (circle.radius * circle.radius)

	det := b*b - (4 * a * c)

	if det < 0 {
		// Do nothing, no intersections
	} else if det == 0 {

		t := -b / (2 * a)

		if t >= 0 && t <= 1 {
			points = append(points, Vector{line.Start.X + t*diff.X, line.Start.Y + t*diff.Y})
		}

	} else {

		t := (-b + math.Sqrt(det)) / (2 * a)

		// We have to ensure t is between 0 and 1; otherwise, the collision points are on the circle as though the lines were infinite in length.
		if t >= 0 && t <= 1 {
			points = append(points, Vector{line.Start.X + t*diff.X, line.Start.Y + t*diff.Y})
		}
		t = (-b - math.Sqrt(det)) / (2 * a)
		if t >= 0 && t <= 1 {
			points = append(points, Vector{line.Start.X + t*diff.X, line.Start.Y + t*diff.Y})
		}

	}

	return points

}
