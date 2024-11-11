package resolv

import (
	"math"
	"sort"
)

// ToDegrees is a helper function to easily convert radians to degrees for human readability.
func ToDegrees(radians float64) float64 {
	return radians / math.Pi * 180
}

// ToRadians is a helper function to easily convert degrees to radians (which is what the rotation-oriented functions in Tetra3D use).
func ToRadians(degrees float64) float64 {
	return math.Pi * degrees / 180
}

func min(a, b float64) float64 {
	if a < b {
		return a
	}
	return b
}

func max(a, b float64) float64 {
	if a > b {
		return a
	}
	return b
}

func clamp(value, min, max float64) float64 {
	if value < min {
		return min
	} else if value > max {
		return max
	}
	return value
}

// func pow(value float64, power int) float64 {
// 	x := value
// 	for i := 0; i < power; i++ {
// 		x += x
// 	}
// 	return x
// }

func round(value float64) float64 {

	iv := float64(int(value))

	if value > iv+0.5 {
		return iv + 1
	} else if value < iv-0.5 {
		return iv - 1
	}

	return iv

}

// Projection represents the projection of a shape (usually a ConvexPolygon) onto an axis for intersection testing.
// Normally, you wouldn't need to get this information, but it could be useful in some circumstances, I'm sure.
type Projection struct {
	Min, Max float64
}

// IsOverlapping returns whether a Projection is overlapping with the other, provided Projection. Credit to https://www.sevenson.com.au/programming/sat/
func (projection Projection) IsOverlapping(other Projection) bool {
	return projection.Overlap(other) > 0
}

// Overlap returns the amount that a Projection is overlapping with the other, provided Projection. Credit to https://dyn4j.org/2010/01/sat/#sat-nointer
func (projection Projection) Overlap(other Projection) float64 {
	return math.Min(projection.Max-other.Min, other.Max-projection.Min)
}

// IsInside returns whether the Projection is wholly inside of the other, provided Projection.
func (projection Projection) IsInside(other Projection) bool {
	return projection.Min >= other.Min && projection.Max <= other.Max
}

// Bounds represents the minimum and maximum bounds of a Shape.
type Bounds struct {
	Min, Max Vector
	space    *Space
}

func (b Bounds) toCellSpace() (int, int, int, int) {

	minX := int(math.Floor(b.Min.X / float64(b.space.cellWidth)))
	minY := int(math.Floor(b.Min.Y / float64(b.space.cellHeight)))
	maxX := int(math.Floor(b.Max.X / float64(b.space.cellWidth)))
	maxY := int(math.Floor(b.Max.Y / float64(b.space.cellHeight)))

	return minX, minY, maxX, maxY
}

// Center returns the center position of the Bounds.
func (b Bounds) Center() Vector {
	return b.Min.Add(b.Max.Sub(b.Min).Scale(0.5))
}

// Width returns the width of the Bounds.
func (b Bounds) Width() float64 {
	return b.Max.X - b.Min.X
}

// Height returns the height of the bounds.
func (b Bounds) Height() float64 {
	return b.Max.Y - b.Min.Y
}

// Move moves the Bounds, such that the center point is offset by {x, y}.
func (b Bounds) Move(x, y float64) Bounds {
	b.Min.X += x
	b.Min.Y += y
	b.Max.X += x
	b.Max.Y += y
	return b
}

// MoveVec moves the Bounds by the vector provided, such that the center point is offset by {x, y}.
func (b *Bounds) MoveVec(vec Vector) Bounds {
	return b.Move(vec.X, vec.Y)
}

// IsIntersecting returns if the Bounds is intersecting with the given other Bounds.
func (b Bounds) IsIntersecting(other Bounds) bool {
	bounds := b.Intersection(other)
	return !bounds.IsEmpty()
}

// Intersection returns the intersection between the two Bounds objects.
func (b Bounds) Intersection(other Bounds) Bounds {

	overlap := Bounds{}

	if other.Max.X < b.Min.X || other.Min.X > b.Max.X || other.Max.Y < b.Min.Y || other.Min.Y > b.Max.Y {
		return overlap
	}

	overlap.Min.X = math.Min(other.Max.X, b.Max.X)
	overlap.Max.X = math.Max(other.Min.X, b.Min.X)

	overlap.Min.Y = math.Min(other.Max.Y, b.Max.Y)
	overlap.Max.Y = math.Max(other.Min.Y, b.Min.Y)

	return overlap

}

// IsEmpty returns true if the Bounds's minimum and maximum corners are 0.
func (b Bounds) IsEmpty() bool {
	return b.Max.X-b.Min.X == 0 && b.Max.Y-b.Min.X == 0
}

/////

// Set represents a Set of elements.
type Set[E comparable] map[E]struct{}

// newSet creates a new set.
func newSet[E comparable]() Set[E] {
	return Set[E]{}
}

// Clone clones the Set.
func (s Set[E]) Clone() Set[E] {
	newSet := newSet[E]()
	newSet.Combine(s)
	return newSet
}

// Set sets the Set to have the same values as in the given other Set.
func (s Set[E]) Set(other Set[E]) {
	s.Clear()
	s.Combine(other)
}

// Add adds the given elements to a set.
func (s Set[E]) Add(elements ...E) {
	for _, element := range elements {
		s[element] = struct{}{}
	}
}

// Combine combines the given other elements to the set.
func (s Set[E]) Combine(otherSet Set[E]) {
	for element := range otherSet {
		s.Add(element)
	}
}

// Contains returns if the set contains the given element.
func (s Set[E]) Contains(element E) bool {
	_, ok := s[element]
	return ok
}

// Remove removes the given element from the set.
func (s Set[E]) Remove(elements ...E) {
	for _, element := range elements {
		delete(s, element)
	}
}

// Clear clears the set.
func (s Set[E]) Clear() {
	for v := range s {
		delete(s, v)
	}
}

// ForEach runs the provided function for each element in the set.
func (s Set[E]) ForEach(f func(element E) bool) {
	for element := range s {
		if !f(element) {
			break
		}
	}
}

/////

// shapeIDSet is an easy way to determine if a shape has been iterated over before (used for filtering through shapes from CellSelections).
type shapeIDSet []uint32

func (s shapeIDSet) idInSet(id uint32) bool {
	for _, v := range s {
		if v == id {
			return true
		}
	}
	return false
}

var cellSelectionForEachIDSet = shapeIDSet{}

/////

// LineTestSettings is a struct of settings to be used when performing line tests (the equivalent of 3D hitscan ray tests for 2D)
type LineTestSettings struct {
	Start       Vector        // The start of the line to test shapes against
	End         Vector        // The end of the line to test chapes against
	TestAgainst ShapeIterator // The collection of shapes to test against
	// The callback to be called for each intersection between the given line, ranging from start to end, and each shape given in TestAgainst.
	// set is the intersection set that contains information about the intersection, index is the index of the current index
	// and count is the total number of intersections detected from the intersection test.
	// The boolean the callback returns indicates whether the LineTest function should continue testing or stop at the currently found intersection.
	OnIntersect  func(set IntersectionSet, index, max int) bool
	callingShape IShape
}

var intersectionSets []IntersectionSet

// LineTest instantly tests a selection of shapes against a ray / line.
// Note that there is no MTV for these results.
func LineTest(settings LineTestSettings) bool {

	castMargin := 0.01 // Basically, the line cast starts are a smidge back so that moving to contact doesn't make future line casts fail
	vu := settings.End.Sub(settings.Start).Unit()
	start := settings.Start.Sub(vu.Scale(castMargin))

	line := newCollidingLine(start.X, start.Y, settings.End.X, settings.End.Y)

	intersectionSets = intersectionSets[:0]

	i := 0

	settings.TestAgainst.ForEach(func(other IShape) bool {

		if other == settings.callingShape {
			return true
		}

		i++

		contactSet := newIntersectionSet()

		switch shape := other.(type) {

		case *Circle:

			res := line.IntersectionPointsCircle(shape)

			if len(res) > 0 {
				for _, contactPoint := range res {
					contactSet.Intersections = append(contactSet.Intersections, Intersection{
						Point:  contactPoint,
						Normal: contactPoint.Sub(shape.position).Unit(),
					})
				}
			}

		case *ConvexPolygon:

			for _, otherLine := range shape.Lines() {

				if point, ok := line.IntersectionPointsLine(otherLine); ok {
					contactSet.Intersections = append(contactSet.Intersections, Intersection{
						Point:  point,
						Normal: otherLine.Normal(),
					})
				}

			}

		}

		if len(contactSet.Intersections) > 0 {

			contactSet.OtherShape = other

			for _, contact := range contactSet.Intersections {
				contactSet.Center = contactSet.Center.Add(contact.Point)
			}

			contactSet.Center.X /= float64(len(contactSet.Intersections))
			contactSet.Center.Y /= float64(len(contactSet.Intersections))

			// Sort the points by distance to line start
			sort.Slice(contactSet.Intersections, func(i, j int) bool {
				return contactSet.Intersections[i].Point.DistanceSquared(settings.Start) < contactSet.Intersections[j].Point.DistanceSquared(settings.Start)
			})

			contactSet.MTV = contactSet.Intersections[0].Point.Sub(settings.Start).Sub(vu.Scale(castMargin))

			intersectionSets = append(intersectionSets, contactSet)

		}

		return true

	})

	// Sort intersection sets by distance from closest hit to line start
	sort.Slice(intersectionSets, func(i, j int) bool {
		return intersectionSets[i].Intersections[0].Point.DistanceSquared(line.Start) < intersectionSets[j].Intersections[0].Point.DistanceSquared(line.Start)
	})

	// Loop through all intersections and iterate through them
	for i, c := range intersectionSets {
		if !settings.OnIntersect(c, i, len(intersectionSets)) {
			break
		}
	}

	return len(intersectionSets) > 0

}
