package resolv

import (
	"math"
	"sort"

	"github.com/quartercastle/vector"
)

type IShape interface {
	// Intersection tests to see if a Shape intersects with the other given Shape. dx and dy are delta movement variables indicating
	// movement to be applied before the intersection check (thereby allowing you to see if a Shape would collide with another if it
	// were in a different relative location). If an Intersection is found, a ContactSet will be returned, giving information regarding
	// the intersection.
	Intersection(dx, dy float64, other IShape) *ContactSet
	// Bounds returns the top-left and bottom-right points of the Shape.
	Bounds() (vector.Vector, vector.Vector)
	// Position returns the X and Y position of the Shape.
	Position() (float64, float64)
	// SetPosition allows you to place a Shape at another location.
	SetPosition(x, y float64)

	Rotation() float64
	SetRotation(radians float64)

	// Rotate rotates the IShape by the radians provided.
	// Note that the rotation goes counter-clockwise from 0 at right to pi/2 in the upwards direction,
	// pi / -pi at left, -pi/2 in the downwards direction, and finally back to 0.
	// This can be visualized as follows:
	//
	//   U
	// L   R
	//   D
	//
	// R: 0
	// U: pi/2
	// L: pi / -pi
	// D: -pi/2
	//
	// This rotation scheme follows the way math.Atan2() works.
	// Note that Rotate(), of course, doesn't do anything for circles for obvious reasons.
	Rotate(radians float64)

	Scale() (float64, float64) // Returns the scale of the IShape (the radius for Circles).
	SetScale(w, h float64)     // Sets the overall scale of the IShape; 1.0 is 100% scale, 2.0 is 200%, and so on. The greater of these values is used for the radius for Circles.

	// Move moves the IShape by the x and y values provided.
	Move(x, y float64)
	// MoveVec moves the IShape by the movement values given in the vector provided.
	MoveVec(vec vector.Vector)

	// Clone duplicates the IShape.
	Clone() IShape
}

// A collidingLine is a helper shape used to determine if two ConvexPolygon lines intersect; you can't create a collidingLine to use as a Shape.
// Instead, you can create a ConvexPolygon, specify two points, and set its Closed value to false (or use NewLine(), as this does it for you).
type collidingLine struct {
	Start, End vector.Vector
}

func new_line(x, y, x2, y2 float64) *collidingLine {
	return &collidingLine{
		Start: vector.Vector{x, y},
		End:   vector.Vector{x2, y2},
	}
}

func (line *collidingLine) Project(axis vector.Vector) vector.Vector {
	return line.Vector().Scale(axis.Dot(line.Start.Sub(line.End)))
}

func (line *collidingLine) Normal() vector.Vector {
	v := line.Vector()
	return vector.Vector{v[1], -v[0]}.Unit()
}

func (line *collidingLine) Vector() vector.Vector {
	return line.End.Clone().Sub(line.Start).Unit()
}

// IntersectionPointsLine returns the intersection point of a Line with another Line as a vector.Vector. If no intersection is found, it will return nil.
func (line *collidingLine) IntersectionPointsLine(other *collidingLine) vector.Vector {

	det := (line.End[0]-line.Start[0])*(other.End[1]-other.Start[1]) - (other.End[0]-other.Start[0])*(line.End[1]-line.Start[1])

	if det != 0 {

		// MAGIC MATH; the extra + 1 here makes it so that corner cases (literally, lines going through corners) works.

		// lambda := (float32(((line.Y-b.Y)*(b.X2-b.X))-((line.X-b.X)*(b.Y2-b.Y))) + 1) / float32(det)
		lambda := (((line.Start[1] - other.Start[1]) * (other.End[0] - other.Start[0])) - ((line.Start[0] - other.Start[0]) * (other.End[1] - other.Start[1])) + 1) / det

		// gamma := (float32(((line.Y-b.Y)*(line.X2-line.X))-((line.X-b.X)*(line.Y2-line.Y))) + 1) / float32(det)
		gamma := (((line.Start[1] - other.Start[1]) * (line.End[0] - line.Start[0])) - ((line.Start[0] - other.Start[0]) * (line.End[1] - line.Start[1])) + 1) / det

		if (0 < lambda && lambda < 1) && (0 < gamma && gamma < 1) {

			// Delta
			dx := line.End[0] - line.Start[0]
			dy := line.End[1] - line.Start[1]

			// dx, dy := line.GetDelta()

			return vector.Vector{line.Start[0] + (lambda * dx), line.Start[1] + (lambda * dy)}
		}

	}

	return nil

}

// IntersectionPointsCircle returns a slice of vector.Vectors, each indicating the intersection point. If no intersection is found, it will return an empty slice.
func (line *collidingLine) IntersectionPointsCircle(circle *Circle) []vector.Vector {

	points := []vector.Vector{}

	cp := vector.Vector{circle.X, circle.Y}
	lStart := line.Start.Sub(cp)
	lEnd := line.End.Sub(cp)
	diff := lEnd.Sub(lStart)

	a := diff[0]*diff[0] + diff[1]*diff[1]
	b := 2 * ((diff[0] * lStart[0]) + (diff[1] * lStart[1]))
	c := (lStart[0] * lStart[0]) + (lStart[1] * lStart[1]) - (circle.radius * circle.radius)

	det := b*b - (4 * a * c)

	if det < 0 {
		// Do nothing, no intersections
	} else if det == 0 {

		t := -b / (2 * a)

		if t >= 0 && t <= 1 {
			points = append(points, vector.Vector{line.Start[0] + t*diff[0], line.Start[1] + t*diff[1]})
		}

	} else {

		t := (-b + math.Sqrt(det)) / (2 * a)

		// We have to ensure t is between 0 and 1; otherwise, the collision points are on the circle as though the lines were infinite in length.
		if t >= 0 && t <= 1 {
			points = append(points, vector.Vector{line.Start[0] + t*diff[0], line.Start[1] + t*diff[1]})
		}
		t = (-b - math.Sqrt(det)) / (2 * a)
		if t >= 0 && t <= 1 {
			points = append(points, vector.Vector{line.Start[0] + t*diff[0], line.Start[1] + t*diff[1]})
		}

	}

	return points

}

// ConvexPolygon represents a series of points, connected by lines, constructing a convex shape.
// The polygon has a position, a scale, a rotation, and may or may not be closed.
type ConvexPolygon struct {
	Points         []vector.Vector // Points represents the points constructing the ConvexPolygon.
	X, Y           float64         // X and Y are the position of the ConvexPolygon.
	ScaleW, ScaleH float64         // The width and height for scaling
	rotation       float64         // How many radians the ConvexPolygon is rotated around in the viewing vector (Z).
	Closed         bool            // Closed is whether the ConvexPolygon is closed or not; only takes effect if there are more than 2 points.
}

// NewConvexPolygon creates a new convex polygon at the position given, from the provided set of X and Y positions of 2D points (or vertices).
// You don't need to pass any points at this stage, but if you do, you should pass whole pairs. The points should generally be ordered clockwise,
// from X and Y of the first, to X and Y of the last.
// For example: NewConvexPolygon(30, 20, 0, 0, 10, 0, 10, 10, 0, 10) would create a 10x10 convex
// polygon square, with the vertices at {0,0}, {10,0}, {10, 10}, and {0, 10}, with the polygon itself occupying a position of 30, 20.
// You can also pass the points using vectors with ConvexPolygon.AddPointsVec().
func NewConvexPolygon(x, y float64, points ...float64) *ConvexPolygon {

	// if len(points)/2 < 2 {
	// 	return nil
	// }

	cp := &ConvexPolygon{
		X:      x,
		Y:      y,
		ScaleW: 1,
		ScaleH: 1,
		Points: []vector.Vector{},
		Closed: true,
	}

	if len(points) > 0 {
		cp.AddPoints(points...)
	}

	return cp
}

// Clone returns a clone of the ConvexPolygon as an IShape.
func (cp *ConvexPolygon) Clone() IShape {

	points := []vector.Vector{}

	for _, point := range cp.Points {
		points = append(points, point.Clone())
	}

	newPoly := NewConvexPolygon(cp.X, cp.Y)
	newPoly.rotation = cp.rotation
	newPoly.ScaleW = cp.ScaleW
	newPoly.ScaleH = cp.ScaleH
	newPoly.AddPointsVec(points...)
	newPoly.Closed = cp.Closed
	return newPoly
}

// AddPointsVec allows you to add points to the ConvexPolygon with a slice of vector.Vectors, each indicating a point / vertex.
func (cp *ConvexPolygon) AddPointsVec(points ...vector.Vector) {
	cp.Points = append(cp.Points, points...)
}

// AddPoints allows you to add points to the ConvexPolygon with a slice or selection of float64s, with each pair indicating an X or Y value for
// a point / vertex (i.e. AddPoints(0, 1, 2, 3) would add two points - one at {0, 1}, and another at {2, 3}).
func (cp *ConvexPolygon) AddPoints(vertexPositions ...float64) {
	if len(vertexPositions) == 0 {
		panic("Error: AddPoints called with 0 passed vertex positions.")
	}
	if len(vertexPositions)%2 == 1 {
		panic("Error: AddPoints called with a non-even amount of vertex positions.")
	}
	for v := 0; v < len(vertexPositions); v += 2 {
		cp.Points = append(cp.Points, vector.Vector{vertexPositions[v], vertexPositions[v+1]})
	}
}

// Lines returns a slice of transformed internalLines composing the ConvexPolygon.
func (cp *ConvexPolygon) Lines() []*collidingLine {

	lines := []*collidingLine{}

	vertices := cp.Transformed()

	for i := 0; i < len(vertices); i++ {

		start, end := vertices[i], vertices[0]

		if i < len(vertices)-1 {
			end = vertices[i+1]
		} else if !cp.Closed || len(cp.Points) <= 2 {
			break
		}

		line := new_line(start[0], start[1], end[0], end[1])

		lines = append(lines, line)

	}

	return lines

}

// Transformed returns the ConvexPolygon's points / vertices, transformed according to the ConvexPolygon's position.
func (cp *ConvexPolygon) Transformed() []vector.Vector {
	transformed := []vector.Vector{}
	for _, point := range cp.Points {
		p := vector.Vector{point[0] * cp.ScaleW, point[1] * cp.ScaleH}
		if cp.rotation != 0 {
			vector.In(p).Rotate(-cp.rotation)
		}
		transformed = append(transformed, vector.Vector{p[0] + cp.X, p[1] + cp.Y})
	}
	return transformed
}

// Bounds returns two Vectors, comprising the top-left and bottom-right positions of the bounds of the
// ConvexPolygon, post-transformation.
func (cp *ConvexPolygon) Bounds() (vector.Vector, vector.Vector) {

	transformed := cp.Transformed()

	topLeft := vector.Vector{transformed[0][0], transformed[0][1]}
	bottomRight := topLeft.Clone()

	for i := 0; i < len(transformed); i++ {

		point := transformed[i]

		if point[0] < topLeft[0] {
			topLeft[0] = point[0]
		} else if point[0] > bottomRight[0] {
			bottomRight[0] = point[0]
		}

		if point[1] < topLeft[1] {
			topLeft[1] = point[1]
		} else if point[1] > bottomRight[1] {
			bottomRight[1] = point[1]
		}

	}
	return topLeft, bottomRight
}

// Position returns the position of the ConvexPolygon.
func (cp *ConvexPolygon) Position() (float64, float64) {
	return cp.X, cp.Y
}

// SetPosition sets the position of the ConvexPolygon. The offset of the vertices compared to the X and Y position is relative to however
// you initially defined the polygon and added the vertices.
func (cp *ConvexPolygon) SetPosition(x, y float64) {
	cp.X = x
	cp.Y = y
}

// SetPositionVec allows you to set the position of the ConvexPolygon using a vector.Vector. The offset of the vertices compared to the X and Y
// position is relative to however you initially defined the polygon and added the vertices.
func (cp *ConvexPolygon) SetPositionVec(vec vector.Vector) {
	cp.X = vec.X()
	cp.Y = vec.Y()
}

// Move translates the ConvexPolygon by the designated X and Y values.
func (cp *ConvexPolygon) Move(x, y float64) {
	cp.X += x
	cp.Y += y
}

// MoveVec translates the ConvexPolygon by the designated vector.Vector.
func (cp *ConvexPolygon) MoveVec(vec vector.Vector) {
	cp.X += vec.X()
	cp.Y += vec.Y()
}

// Center returns the transformed Center of the ConvexPolygon.
func (cp *ConvexPolygon) Center() vector.Vector {

	pos := vector.Vector{0, 0}

	for _, v := range cp.Transformed() {
		pos.Add(v)
	}

	pos[0] /= float64(len(cp.Transformed()))
	pos[1] /= float64(len(cp.Transformed()))

	return pos

}

// Project projects (i.e. flattens) the ConvexPolygon onto the provided axis.
func (cp *ConvexPolygon) Project(axis vector.Vector) Projection {
	axis = axis.Unit()
	vertices := cp.Transformed()
	min := dot(axis, vertices[0]) // We use a manual dot function here instead of Vector.Dot() because some idiot (me) smashed the dot product to a range of -1 to 1
	max := min
	for i := 1; i < len(vertices); i++ {
		p := dot(axis, vertices[i])
		if p < min {
			min = p
		} else if p > max {
			max = p
		}
	}
	return Projection{min, max}
}

// SATAxes returns the axes of the ConvexPolygon for SAT intersection testing.
func (cp *ConvexPolygon) SATAxes() []vector.Vector {

	axes := []vector.Vector{}
	for _, line := range cp.Lines() {
		axes = append(axes, line.Normal())
	}
	return axes

}

// PointInside returns if a Point (a vector.Vector) is inside the ConvexPolygon.
func (polygon *ConvexPolygon) PointInside(point vector.Vector) bool {

	pointLine := new_line(point[0], point[1], point[0]+999999999999, point[1])

	contactCount := 0

	for _, line := range polygon.Lines() {

		if line.IntersectionPointsLine(pointLine) != nil {
			contactCount++
		}

	}

	return contactCount%2 == 1
}

// Rotation returns the rotation (in radians) of the ConvexPolygon.
func (polygon *ConvexPolygon) Rotation() float64 {
	return polygon.rotation
}

// SetRotation sets the rotation for the ConvexPolygon; note that the rotation goes counter-clockwise from 0 to pi, and then from -pi at 180 down, back to 0.
// This can be visualized as follows:
//
//	 	        (Pi / 2)
//						|
//						|
//
// (Pi / -Pi) ------------- (0)
//
//		|
//		|
//	(-Pi / 2)
//
// This rotation scheme follows the way math.Atan2() works.
func (polygon *ConvexPolygon) SetRotation(radians float64) {
	polygon.rotation = radians
	if polygon.rotation > math.Pi {
		polygon.rotation -= math.Pi * 2
	} else if polygon.rotation < -math.Pi {
		polygon.rotation += math.Pi * 2
	}
}

// Rotate is a helper function to rotate a ConvexPolygon by the radians given.
func (polygon *ConvexPolygon) Rotate(radians float64) {
	polygon.SetRotation(polygon.Rotation() + radians)
}

// Scale returns the scale multipliers of the ConvexPolygon.
func (polygon *ConvexPolygon) Scale() (float64, float64) {
	return polygon.ScaleW, polygon.ScaleH
}

// SetScale sets the scale multipliers of the ConvexPolygon.
func (polygon *ConvexPolygon) SetScale(w, h float64) {
	polygon.ScaleW = w
	polygon.ScaleH = h
}

type ContactSet struct {
	Points []vector.Vector // Slice of Points indicating contact between the two Shapes.
	MTV    vector.Vector   // Minimum Translation Vector; this is the vector to move a Shape on to move it outside of its contacting Shape.
	Center vector.Vector   // Center of the Contact set; this is the average of all Points contained within the Contact Set.
}

func NewContactSet() *ContactSet {
	return &ContactSet{
		Points: []vector.Vector{},
		MTV:    vector.Vector{0, 0},
		Center: vector.Vector{0, 0},
	}
}

// LeftmostPoint returns the left-most point out of the ContactSet's Points slice. If the Points slice is empty somehow, this returns nil.
func (cs *ContactSet) LeftmostPoint() vector.Vector {

	var left vector.Vector

	for _, point := range cs.Points {

		if left == nil || point[0] < left[0] {
			left = point
		}

	}

	return left

}

// RightmostPoint returns the right-most point out of the ContactSet's Points slice. If the Points slice is empty somehow, this returns nil.
func (cs *ContactSet) RightmostPoint() vector.Vector {

	var right vector.Vector

	for _, point := range cs.Points {

		if right == nil || point[0] > right[0] {
			right = point
		}

	}

	return right

}

// TopmostPoint returns the top-most point out of the ContactSet's Points slice. If the Points slice is empty somehow, this returns nil.
func (cs *ContactSet) TopmostPoint() vector.Vector {

	var top vector.Vector

	for _, point := range cs.Points {

		if top == nil || point[1] < top[1] {
			top = point
		}

	}

	return top

}

// BottommostPoint returns the bottom-most point out of the ContactSet's Points slice. If the Points slice is empty somehow, this returns nil.
func (cs *ContactSet) BottommostPoint() vector.Vector {

	var bottom vector.Vector

	for _, point := range cs.Points {

		if bottom == nil || point[1] > bottom[1] {
			bottom = point
		}

	}

	return bottom

}

// Intersection tests to see if a ConvexPolygon intersects with the other given Shape. dx and dy are delta movement variables indicating
// movement to be applied before the intersection check (thereby allowing you to see if a Shape would collide with another if it
// were in a different relative location). If an Intersection is found, a ContactSet will be returned, giving information regarding
// the intersection.
func (cp *ConvexPolygon) Intersection(dx, dy float64, other IShape) *ContactSet {

	contactSet := NewContactSet()

	ogX := cp.X
	ogY := cp.Y
	cp.X += dx
	cp.Y += dy

	if circle, isCircle := other.(*Circle); isCircle {

		for _, line := range cp.Lines() {
			contactSet.Points = append(contactSet.Points, line.IntersectionPointsCircle(circle)...)
		}

	} else if poly, isPoly := other.(*ConvexPolygon); isPoly {

		for _, line := range cp.Lines() {

			for _, otherLine := range poly.Lines() {

				if point := line.IntersectionPointsLine(otherLine); point != nil {
					contactSet.Points = append(contactSet.Points, point)
				}

			}

		}

	}

	if len(contactSet.Points) > 0 {

		for _, point := range contactSet.Points {
			contactSet.Center = contactSet.Center.Add(point)
		}

		contactSet.Center[0] /= float64(len(contactSet.Points))
		contactSet.Center[1] /= float64(len(contactSet.Points))

		if mtv := cp.calculateMTV(contactSet, other); mtv != nil {
			contactSet.MTV = mtv
		}

	} else {
		contactSet = nil
	}

	// If dx or dy aren't 0, then the MTV will be greater to compensate; this adjusts the vector back.
	if contactSet != nil && (dx != 0 || dy != 0) {
		deltaMagnitude := vector.Vector{dx, dy}.Magnitude()
		ogMagnitude := contactSet.MTV.Magnitude()
		contactSet.MTV = contactSet.MTV.Unit().Scale(ogMagnitude - deltaMagnitude)
	}

	cp.X = ogX
	cp.Y = ogY

	return contactSet

}

// calculateMTV returns the MTV, if possible, and a bool indicating whether it was possible or not.
func (cp *ConvexPolygon) calculateMTV(contactSet *ContactSet, otherShape IShape) vector.Vector {

	delta := vector.Vector{0, 0}

	smallest := vector.Vector{math.MaxFloat64, 0}

	switch other := otherShape.(type) {

	case *ConvexPolygon:

		for _, axis := range cp.SATAxes() {
			pa := cp.Project(axis)
			pb := other.Project(axis)

			overlap := pa.Overlap(pb)

			if overlap <= 0 {
				return nil
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
				return nil
			}

			if smallest.Magnitude() > overlap {
				smallest = axis.Scale(overlap)
			}

		}

	case *Circle:

		verts := append([]vector.Vector{}, cp.Transformed()...)
		// The center point of a contact could also be closer than the verts, particularly if we're testing from a Circle to another Shape.
		verts = append(verts, contactSet.Center)
		center := vector.Vector{other.X, other.Y}
		sort.Slice(verts, func(i, j int) bool { return verts[i].Sub(center).Magnitude() < verts[j].Sub(center).Magnitude() })

		smallest = vector.Vector{center[0] - verts[0][0], center[1] - verts[0][1]}
		smallest = smallest.Unit().Scale(smallest.Magnitude() - other.radius)

	}

	delta[0] = smallest[0]
	delta[1] = smallest[1]

	return delta
}

// ContainedBy returns if the ConvexPolygon is wholly contained by the other shape provided.
func (cp *ConvexPolygon) ContainedBy(otherShape IShape) bool {

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

	}

	return true
}

// FlipH flips the ConvexPolygon's vertices horizontally according to their initial offset when adding the points.
func (cp *ConvexPolygon) FlipH() {

	for _, v := range cp.Points {
		v[0] = -v[0]
	}
	// We have to reverse vertex order after flipping the vertices to ensure the winding order is consistent between Objects (so that the normals are consistently outside or inside, which is important
	// when doing Intersection tests). If we assume that the normal of a line, going from vertex A to vertex B, is one direction, then the normal would be inverted if the vertices were flipped in position,
	// but not in order. This would make Intersection tests drive objects into each other, instead of giving the delta to move away.
	cp.ReverseVertexOrder()

}

// FlipV flips the ConvexPolygon's vertices vertically according to their initial offset when adding the points.
func (cp *ConvexPolygon) FlipV() {

	for _, v := range cp.Points {
		v[1] = -v[1]
	}
	cp.ReverseVertexOrder()

}

// RecenterPoints recenters the vertices in the polygon, such that they are all equidistant from the center.
// For example, say you had a polygon with the following three points: {0, 0}, {10, 0}, {0, 16}.
// After calling cp.RecenterPoints(), the polygon's points would be at {-5, -8}, {5, -8}, {-5, 8}.
func (cp *ConvexPolygon) RecenterPoints() {

	if len(cp.Points) <= 1 {
		return
	}

	offset := vector.Vector{0, 0}
	for _, p := range cp.Points {
		vector.In(offset).Add(p)
	}

	vector.In(offset).Scale(1.0 / float64(len(cp.Points))).Invert()

	for _, p := range cp.Points {
		vector.In(p).Add(offset)
	}

}

// ReverseVertexOrder reverses the vertex ordering of the ConvexPolygon.
func (cp *ConvexPolygon) ReverseVertexOrder() {

	verts := []vector.Vector{cp.Points[0]}

	for i := len(cp.Points) - 1; i >= 1; i-- {
		verts = append(verts, cp.Points[i])
	}

	cp.Points = verts

}

// NewRectangle returns a rectangular ConvexPolygon at the {x, y} position given with the vertices ordered in clockwise order,
// positioned at {0, 0}, {w, 0}, {w, h}, {0, h}.
// TODO: In actuality, an AABBRectangle should be its own "thing" with its own optimized Intersection code check.
func NewRectangle(x, y, w, h float64) *ConvexPolygon {
	return NewConvexPolygon(
		x, y,

		0, 0,
		w, 0,
		w, h,
		0, h,
	)
}

// NewLine is a helper function that returns a ConvexPolygon composed of a single line. The Polygon has a position of x1, y1, and has a width and height
// equivalent to x2-x1 and y2-y1 (so the end of the line is at x2, y2).
func NewLine(x1, y1, x2, y2 float64) *ConvexPolygon {
	newLine := NewConvexPolygon(x1, y1,
		0, 0,
		x2-x1, y2-y1,
	)
	newLine.Closed = false // This actually isn't necessary for a one-sided polygon
	return newLine
}

type Circle struct {
	X, Y, radius   float64
	originalRadius float64
	scale          float64
}

// NewCircle returns a new Circle, with its center at the X and Y position given, and with the defined radius.
func NewCircle(x, y, radius float64) *Circle {
	circle := &Circle{
		X:              x,
		Y:              y,
		radius:         radius,
		originalRadius: radius,
		scale:          1,
	}
	return circle
}

func (circle *Circle) Clone() IShape {
	newCircle := NewCircle(circle.X, circle.Y, circle.radius)
	newCircle.originalRadius = circle.originalRadius
	newCircle.scale = circle.scale
	return newCircle
}

// Bounds returns the top-left and bottom-right corners of the Circle.
func (circle *Circle) Bounds() (vector.Vector, vector.Vector) {
	return vector.Vector{circle.X - circle.radius, circle.Y - circle.radius}, vector.Vector{circle.X + circle.radius, circle.Y + circle.radius}
}

// Intersection tests to see if a Circle intersects with the other given Shape. dx and dy are delta movement variables indicating
// movement to be applied before the intersection check (thereby allowing you to see if a Shape would collide with another if it
// were in a different relative location). If an Intersection is found, a ContactSet will be returned, giving information regarding
// the intersection.
func (circle *Circle) Intersection(dx, dy float64, other IShape) *ContactSet {

	var contactSet *ContactSet

	ox := circle.X
	oy := circle.Y

	circle.X += dx
	circle.Y += dy

	// here

	switch shape := other.(type) {
	case *ConvexPolygon:
		// Maybe this would work?
		contactSet = shape.Intersection(-dx, -dy, circle)
		if contactSet != nil {
			contactSet.MTV = contactSet.MTV.Scale(-1)
		}
	case *Circle:

		contactSet = NewContactSet()

		contactSet.Points = circle.IntersectionPointsCircle(shape)

		if len(contactSet.Points) == 0 {
			return nil
		}

		contactSet.MTV = vector.Vector{circle.X - shape.X, circle.Y - shape.Y}
		dist := contactSet.MTV.Magnitude()
		contactSet.MTV = contactSet.MTV.Unit().Scale(circle.radius + shape.radius - dist)

		for _, point := range contactSet.Points {
			contactSet.Center = contactSet.Center.Add(point)
		}

		contactSet.Center[0] /= float64(len(contactSet.Points))
		contactSet.Center[1] /= float64(len(contactSet.Points))

		// if contactSet != nil {
		// 	contactSet.MTV[0] -= dx
		// 	contactSet.MTV[1] -= dy
		// }

		// contactSet.MTV = vector.Vector{circle.X - shape.X, circle.Y - shape.Y}
	}

	circle.X = ox
	circle.Y = oy

	return contactSet
}

// Move translates the Circle by the designated X and Y values.
func (circle *Circle) Move(x, y float64) {
	circle.X += x
	circle.Y += y
}

// MoveVec translates the Circle by the designated vector.Vector.
func (circle *Circle) MoveVec(vec vector.Vector) {
	circle.X += vec.X()
	circle.Y += vec.Y()
}

// SetPosition sets the center position of the Circle using the X and Y values given.
func (circle *Circle) SetPosition(x, y float64) {
	circle.X = x
	circle.Y = y
}

// SetPosition sets the center position of the Circle using the vector.Vector given.
func (circle *Circle) SetPositionVec(vec vector.Vector) {
	circle.X = vec.X()
	circle.Y = vec.Y()
}

// Position() returns the X and Y position of the Circle.
func (circle *Circle) Position() (float64, float64) {
	return circle.X, circle.Y
}

// PointInside returns if the given vector.Vector is inside of the circle.
func (circle *Circle) PointInside(point vector.Vector) bool {
	return point.Sub(vector.Vector{circle.X, circle.Y}).Magnitude() <= circle.radius
}

// IntersectionPointsCircle returns the intersection points of the two circles provided.
func (circle *Circle) IntersectionPointsCircle(other *Circle) []vector.Vector {

	d := math.Sqrt(math.Pow(other.X-circle.X, 2) + math.Pow(other.Y-circle.Y, 2))

	if d > circle.radius+other.radius || d < math.Abs(circle.radius-other.radius) || d == 0 && circle.radius == other.radius {
		return nil
	}

	a := (math.Pow(circle.radius, 2) - math.Pow(other.radius, 2) + math.Pow(d, 2)) / (2 * d)
	h := math.Sqrt(math.Pow(circle.radius, 2) - math.Pow(a, 2))

	x2 := circle.X + a*(other.X-circle.X)/d
	y2 := circle.Y + a*(other.Y-circle.Y)/d

	return []vector.Vector{
		{x2 + h*(other.Y-circle.Y)/d, y2 - h*(other.X-circle.X)/d},
		{x2 - h*(other.Y-circle.Y)/d, y2 + h*(other.X-circle.X)/d},
	}

}

// Circles can't rotate, of course. This function is just a stub to make them acceptable as IShapes.
func (circle *Circle) Rotate(radians float64) {}

// Circles can't rotate, of course. This function is just a stub to make them acceptable as IShapes.
func (circle *Circle) SetRotation(rotation float64) {}

// Circles can't rotate, of course. This function is just a stub to make them acceptable as IShapes.
func (circle *Circle) Rotation() float64 {
	return 0
}

// Scale returns the scale multiplier of the Circle, twice; this is to have it adhere to the
func (circle *Circle) Scale() (float64, float64) {
	return circle.scale, circle.scale
}

// SetScale sets the scale multiplier of the Circle (this is W / H to have it adhere to IShape as a
// contract; in truth, the Circle will be set to 0.5 * the maximum out of the width and height
// height values given).
func (circle *Circle) SetScale(w, h float64) {
	circle.scale = math.Max(w, h)
	circle.radius = circle.originalRadius * circle.scale
}

// Radius returns the radius of the Circle.
func (circle *Circle) Radius() float64 {
	return circle.radius
}

// SetRadius sets the radius of the Circle, updating the scale multiplier to reflect this change.
func (circle *Circle) SetRadius(radius float64) {
	circle.radius = radius
	circle.scale = circle.radius / circle.originalRadius
}

// // MultiShape is a Shape comprised of other sub-shapes.
// type MultiShape struct {
// 	Shapes []Shape
// }

// // NewMultiShape returns a new MultiShape.
// func NewMultiShape() *MultiShape {
// 	return &MultiShape{
// 		Shapes: []Shape{},
// 	}
// }

// // Add adds Shapes to the MultiShape.
// func (ms *MultiShape) Add(shapes ...Shape) {
// 	ms.Shapes = append(ms.Shapes, shapes...)
// }

// // Remove removes Shapes from the MultiShape.
// func (ms *MultiShape) Remove(shapes ...Shape) {

// 	for _, toRemove := range shapes {

// 		for i, s := range ms.Shapes {

// 			if toRemove == s {
// 				ms.Shapes[i] = nil
// 				ms.Shapes = append(ms.Shapes[:i], ms.Shapes[i+1:]...)
// 				break
// 			}

// 		}

// 	}

// }

// func (ms *MultiShape) Clone() Shape {
// 	newMS := NewMultiShape()
// 	for _, shape := range ms.Shapes {
// 		newMS.Add(shape.Clone())
// 	}
// 	return newMS
// }

// // SetPosition sets the position of the MultiShape by first setting the position of the first
// // sub-Shape it "owns", and then offsetting every other shape by the same direction.
// func (ms *MultiShape) SetPosition(x, y float64) {

// 	if len(ms.Shapes) > 0 {

// 		center := ms.Shapes[0]

// 		cx, cy := center.Position()
// 		deltaX := x - cx
// 		deltaY := y - cy

// 		center.SetPosition(x, y)

// 		for i := 1; i < len(ms.Shapes); i++ {
// 			posX, posY := ms.Shapes[i].Position()
// 			posX += deltaX
// 			posY += deltaY
// 			ms.Shapes[i].SetPosition(posX, posY)
// 		}

// 	}

// }

// // Move moves all Shapes contained in the MultiShape by the delta X and Y values given.
// func (ms *MultiShape) Move(dx, dy float64) {
// 	for _, shape := range ms.Shapes {
// 		x, y := shape.Position()
// 		shape.SetPosition(x+dx, y+dy)
// 	}
// }

// // Position returns the position of the first Shape added to the MultiShape. If the MultiShape doesn't contain any other Shapes,
// // 0, 0 is returned.
// func (ms *MultiShape) Position() (float64, float64) {
// 	if len(ms.Shapes) > 0 {
// 		return ms.Shapes[0].Position()
// 	}
// 	return 0, 0
// }

// // Intersection tests to see if a MultiShape intersects with the other given Shape. dx and dy are delta movement variables indicating
// // movement to be applied before the intersection check (thereby allowing you to see if a Shape would collide with another if it
// // were in a different relative location). If an Intersection is found, a ContactSet will be returned, giving information regarding
// // the intersection. Intersection returns the first intersection between any of its sub-Shapes, and the other Shape provided.
// func (ms *MultiShape) Intersection(dx, dy float64, other Shape) *ContactSet {

// 	for _, shape := range ms.Shapes {
// 		if intersection := shape.Intersection(dx, dy, other); intersection != nil {
// 			return intersection
// 		}
// 	}

// 	return nil

// }

// // Bounds returns a slice of points comprising the top-left-most and bottom-right-most
// // positions of all shapes contained within the MultiShape.
// func (ms *MultiShape) Bounds() (vector.Vector, vector.Vector) {

// 	if len(ms.Shapes) == 0 {
// 		return vector.Vector{}, vector.Vector{}
// 	}

// 	topLeft, bottomRight := ms.Shapes[0].Bounds()

// 	for i := 1; i < len(ms.Shapes); i++ {

// 		tl, br := ms.Shapes[i].Bounds()

// 		if tl[0] < topLeft[0] {
// 			topLeft[0] = tl[0]
// 		} else if br[0] > bottomRight[0] {
// 			bottomRight[0] = br[0]
// 		}

// 		if tl[1] < topLeft[1] {
// 			topLeft[1] = tl[1]
// 		} else if br[1] > bottomRight[1] {
// 			bottomRight[1] = br[1]
// 		}
// 	}

// 	return topLeft, bottomRight

// }

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
type Projection struct {
	Min, Max float64
}

// Overlapping returns whether a Projection is overlapping with the other, provided Projection. Credit to https://www.sevenson.com.au/programming/sat/
func (projection Projection) Overlapping(other Projection) bool {
	return projection.Overlap(other) > 0
}

// Overlap returns the amount that a Projection is overlapping with the other, provided Projection. Credit to https://dyn4j.org/2010/01/sat/#sat-nointer
func (projection Projection) Overlap(other Projection) float64 {
	return math.Min(projection.Max, other.Max) - math.Max(projection.Min, other.Min)
}

// IsInside returns whether the Projection is wholly inside of the other, provided Projection.
func (projection Projection) IsInside(other Projection) bool {
	return projection.Min >= other.Min && projection.Max <= other.Max
}
