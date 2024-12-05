package resolv

import (
	"math"
	"sort"
)

// IShape represents an interface that all Shapes fulfill.
type IShape interface {
	ID() uint32 // The unique ID of the Shape
	Clone() IShape
	Tags() *Tags

	Position() Vector
	SetPosition(x, y float64)
	SetPositionVec(vec Vector)

	Move(x, y float64)
	MoveVec(vec Vector)

	SetX(float64)
	SetY(float64)

	SelectTouchingCells(margin int) CellSelection

	update()

	SetData(data any)
	Data() any

	setSpace(space *Space)
	Space() *Space

	addToTouchingCells()
	removeFromTouchingCells()
	Bounds() Bounds

	IsLeftOf(other IShape) bool
	IsRightOf(other IShape) bool
	IsAbove(other IShape) bool
	IsBelow(other IShape) bool

	IntersectionTest(settings IntersectionTestSettings) bool
	IsIntersecting(other IShape) bool
	Intersection(other IShape) IntersectionSet

	VecTo(other IShape) Vector
	DistanceTo(other IShape) float64
	DistanceSquaredTo(other IShape) float64
}

// ShapeBase implements many of the common methods that Shapes need to implement to fulfill IShape
// (but not the Shape-specific ones, like rotating for ConvexPolygons or setting the radius for Circles).
type ShapeBase struct {
	position      Vector
	space         *Space
	touchingCells []*Cell
	tags          *Tags
	data          any    // Data represents some helper data present on the shape.
	owner         IShape // The owning shape; this allows ShapeBase to call overridden functions (i.e. owner.Bounds()).
	id            uint32
}

var globalShapeID = uint32(0)

func newShapeBase(x, y float64) ShapeBase {
	t := Tags(0)
	id := globalShapeID
	globalShapeID++
	return ShapeBase{
		position: NewVector(x, y),
		tags:     &t,
		id:       id,
	}
}

// ID returns the unique ID of the Shape.
func (s *ShapeBase) ID() uint32 {
	return s.id
}

// Data returns any auxiliary data set on the Circle shape.
func (s *ShapeBase) Data() any {
	return s.data
}

// SetData sets any auxiliary data on the Circle shape.
func (s *ShapeBase) SetData(data any) {
	s.data = data
}

// Tags returns the tags applied to the shape.
func (s *ShapeBase) Tags() *Tags {
	return s.tags
}

// Move translates the Circle by the designated X and Y values.
func (s *ShapeBase) Move(x, y float64) {
	s.position.X += x
	s.position.Y += y
	s.update()
}

// MoveVec translates the ShapeBase by the designated Vector.
func (s *ShapeBase) MoveVec(vec Vector) {
	s.Move(vec.X, vec.Y)
}

// Position() returns the X and Y position of the ShapeBase.
func (s *ShapeBase) Position() Vector {
	return s.position
}

// SetPosition sets the center position of the ShapeBase using the X and Y values given.
func (s *ShapeBase) SetPosition(x, y float64) {
	s.position.X = x
	s.position.Y = y
	s.update()
}

// SetPosition sets the center position of the ShapeBase using the Vector given.
func (c *ShapeBase) SetPositionVec(vec Vector) {
	c.SetPosition(vec.X, vec.Y)
}

// SetX sets the X position of the Shape.
func (c *ShapeBase) SetX(x float64) {
	pos := c.position
	pos.X = x
	c.SetPosition(pos.X, pos.Y)
}

// SetY sets the Y position of the Shape.
func (c *ShapeBase) SetY(y float64) {
	pos := c.position
	pos.Y = y
	c.SetPosition(pos.X, pos.Y)
}

func (s *ShapeBase) Space() *Space {
	return s.space
}

func (s *ShapeBase) setSpace(space *Space) {
	s.space = space
}

// IsLeftOf returns true if the Shape is to the left of the other shape.
func (s *ShapeBase) IsLeftOf(other IShape) bool {
	return s.owner.Bounds().Min.X < other.Bounds().Min.X
}

// IsRightOf returns true if the Shape is to the right of the other shape.
func (s *ShapeBase) IsRightOf(other IShape) bool {
	return s.owner.Bounds().Max.X > other.Bounds().Max.X
}

// IsAbove returns true if the Shape is above the other shape.
func (s *ShapeBase) IsAbove(other IShape) bool {
	return s.owner.Bounds().Min.Y < other.Bounds().Min.Y
}

// IsBelow returns true if the Shape is below the other shape.
func (s *ShapeBase) IsBelow(other IShape) bool {
	return s.owner.Bounds().Max.Y > other.Bounds().Max.Y
}

// VecTo returns a vector from the given shape to the other Shape.
func (s *ShapeBase) VecTo(other IShape) Vector {
	return s.position.Sub(other.Position())
}

// DistanceSquaredTo returns the distance from the given shape's center to the other Shape.
func (s *ShapeBase) DistanceTo(other IShape) float64 {
	return s.owner.Position().Distance(other.Position())
}

// DistanceSquaredTo returns the squared distance from the given shape's center to the other Shape.
func (s *ShapeBase) DistanceSquaredTo(other IShape) float64 {
	return s.owner.Position().DistanceSquared(other.Position())
}

func (s *ShapeBase) removeFromTouchingCells() {
	for _, cell := range s.touchingCells {
		cell.unregister(s.owner)
	}

	s.touchingCells = s.touchingCells[:0]
}

func (s *ShapeBase) addToTouchingCells() {

	if s.space != nil {

		cx, cy, ex, ey := s.owner.Bounds().toCellSpace()

		for y := cy; y <= ey; y++ {

			for x := cx; x <= ex; x++ {

				cell := s.space.Cell(x, y)

				if cell != nil {
					cell.register(s.owner)
					s.touchingCells = append(s.touchingCells, cell)
				}

			}

		}

	}

}

// SelectTouchingCells returns a CellSelection of the cells in the Space that the Shape is touching.
// margin sets the cellular margin - the higher the margin, the further away candidate Shapes can be to be considered for
// collision. A margin of 1 is a good default. To help visualize which cells contain Shapes, it would be good to implement some kind of debug
// drawing in your game, like can be seen in resolv's examples.
func (s *ShapeBase) SelectTouchingCells(margin int) CellSelection {

	cx, cy, ex, ey := s.owner.Bounds().toCellSpace()

	cx -= margin
	cy -= margin
	ex += margin
	ey += margin

	return CellSelection{
		StartX:      cx,
		StartY:      cy,
		EndX:        ex,
		EndY:        ey,
		space:       s.space,
		excludeSelf: s.owner,
	}
}

func (s *ShapeBase) update() {
	s.removeFromTouchingCells()
	s.addToTouchingCells()
}

// IsIntersecting returns if the shape is intersecting with the other given Shape.
func (s *ShapeBase) IsIntersecting(other IShape) bool {
	return !s.owner.Intersection(other).IsEmpty()
}

// IntersectionTestSettings is a struct that contains settings to control intersection tests.
type IntersectionTestSettings struct {
	TestAgainst ShapeIterator // The collection of shapes to test against
	// OnIntersect is a callback to be called for each intersection found between the calling Shape and any of the other shapes given in TestAgainst.
	// The callback should be called in order of distance to the testing object.
	// Moving the object can influence whether it intersects with future surrounding objects.
	// set is the intersection set that contains information about the intersection.
	// The boolean the callback returns indicates whether the LineTest function should continue testing or stop at the currently found intersection.
	OnIntersect func(set IntersectionSet) bool
}

type possibleIntersection struct {
	Shape    IShape
	Distance float64
}

var possibleIntersections []possibleIntersection

// IntersectionTest tests to see if the calling shape intersects with shapes specified in
// the given settings struct, checked in order of distance to the calling shape's center point.
// Internally, the function checks to see what Shapes are nearby, and tests against them in order
// of distance. If the testing Shape moves, then that will influence the result of testing future
// Shapes in the current game frame.
// If the test succeeds in finding at least one intersection, it returns true.
func (s *ShapeBase) IntersectionTest(settings IntersectionTestSettings) bool {

	possibleIntersections = possibleIntersections[:0]

	settings.TestAgainst.ForEach(func(other IShape) bool {

		if other == s.owner {
			return true
		}

		possibleIntersections = append(possibleIntersections, possibleIntersection{
			Shape:    other,
			Distance: other.DistanceSquaredTo(s.owner),
		})
		return true

	})

	sort.Slice(possibleIntersections, func(i, j int) bool {
		return possibleIntersections[i].Distance < possibleIntersections[j].Distance
	})

	collided := false

	for _, p := range possibleIntersections {

		result := s.owner.Intersection(p.Shape)

		if !result.IsEmpty() {
			collided = true
			if settings.OnIntersect != nil {
				if !settings.OnIntersect(result) {
					break
				}
			} else {
				break
			}
		}

	}

	return collided

}

func circleConvexTest(circle *Circle, convex *ConvexPolygon) IntersectionSet {

	intersectionSet := IntersectionSet{}

	if !convex.owner.Bounds().IsIntersecting(circle.Bounds()) {
		return intersectionSet
	}

	for _, line := range convex.Lines() {

		if res := line.IntersectionPointsCircle(circle); len(res) > 0 {

			for _, point := range res {
				intersectionSet.Intersections = append(intersectionSet.Intersections, Intersection{
					Point:  point,
					Normal: line.Normal(),
				})
			}
		}

	}

	if !intersectionSet.IsEmpty() {

		intersectionSet.OtherShape = convex

		// No point in sorting circle -> convex intersection tests because the circle's center is necessarily equidistant from any and all points of intersection

		if mtv, ok := convex.calculateMTV(circle); ok {
			intersectionSet.MTV = mtv.Invert()
		}

	}

	return intersectionSet

}

func convexCircleTest(convex *ConvexPolygon, circle *Circle) IntersectionSet {

	intersectionSet := IntersectionSet{}

	if !convex.owner.Bounds().IsIntersecting(circle.Bounds()) {
		return intersectionSet
	}

	for _, line := range convex.Lines() {

		res := line.IntersectionPointsCircle(circle)

		if len(res) > 0 {

			for _, point := range res {
				intersectionSet.Intersections = append(intersectionSet.Intersections, Intersection{
					Point:  point,
					Normal: point.Sub(circle.position).Unit(),
				})
			}

		}

	}

	if !intersectionSet.IsEmpty() {

		intersectionSet.OtherShape = circle

		sort.Slice(intersectionSet.Intersections, func(i, j int) bool {
			return intersectionSet.Intersections[i].Point.DistanceSquared(circle.position) < intersectionSet.Intersections[j].Point.DistanceSquared(circle.position)
		})

		if mtv, ok := convex.calculateMTV(circle); ok {
			intersectionSet.MTV = mtv
		}

	}

	return intersectionSet

}

func circleCircleTest(circleA, circleB *Circle) IntersectionSet {

	intersectionSet := IntersectionSet{}

	if !circleA.owner.Bounds().IsIntersecting(circleB.Bounds()) {
		return intersectionSet
	}

	d := math.Sqrt(math.Pow(circleB.position.X-circleA.position.X, 2) + math.Pow(circleB.position.Y-circleA.position.Y, 2))

	if d > circleA.radius+circleB.radius || d < math.Abs(circleA.radius-circleB.radius) || d == 0 {
		return intersectionSet
	}

	a := (math.Pow(circleA.radius, 2) - math.Pow(circleB.radius, 2) + math.Pow(d, 2)) / (2 * d)
	h := math.Sqrt(math.Pow(circleA.radius, 2) - math.Pow(a, 2))

	x2 := circleA.position.X + a*(circleB.position.X-circleA.position.X)/d
	y2 := circleA.position.Y + a*(circleB.position.Y-circleA.position.Y)/d

	intersectionSet.Intersections = []Intersection{
		{Point: Vector{x2 + h*(circleB.position.Y-circleA.position.Y)/d, y2 - h*(circleB.position.X-circleA.position.X)/d}},
		{Point: Vector{x2 - h*(circleB.position.Y-circleA.position.Y)/d, y2 + h*(circleB.position.X-circleA.position.X)/d}},
	}

	for i := range intersectionSet.Intersections {
		intersectionSet.Intersections[i].Normal = intersectionSet.Intersections[i].Point.Sub(circleA.position).Unit()
	}

	intersectionSet.MTV = Vector{circleA.position.X - circleB.position.X, circleA.position.Y - circleB.position.Y}
	dist := intersectionSet.MTV.Magnitude()
	intersectionSet.MTV = intersectionSet.MTV.Unit().Scale(circleA.radius + circleB.radius - dist)

	intersectionSet.OtherShape = circleB

	return intersectionSet

}

func convexConvexTest(convexA, convexB *ConvexPolygon) IntersectionSet {

	intersectionSet := IntersectionSet{}

	if !convexA.owner.Bounds().IsIntersecting(convexB.Bounds()) {
		return intersectionSet
	}

	for _, otherLine := range convexB.Lines() {

		for _, line := range convexA.Lines() {

			if point, ok := line.IntersectionPointsLine(otherLine); ok {
				intersectionSet.Intersections = append(intersectionSet.Intersections, Intersection{
					Point:  point,
					Normal: otherLine.Normal(),
				})
			}

		}

	}

	if !intersectionSet.IsEmpty() {

		intersectionSet.OtherShape = convexB

		center := convexA.Center()

		sort.Slice(intersectionSet.Intersections, func(i, j int) bool {
			return intersectionSet.Intersections[i].Point.DistanceSquared(center) < intersectionSet.Intersections[j].Point.DistanceSquared(center)
		})

		if mtv, ok := convexA.calculateMTV(convexB); ok {
			intersectionSet.MTV = mtv
		}

	}

	return intersectionSet

}
