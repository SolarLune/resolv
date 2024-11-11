package resolv

// ShapeIterator is an interface that defines a method to iterate through Shapes.
// Any object that has such a function (e.g. a ShapeFilter or a ShapeCollection (which is essentially just a slice of Shapes)) fulfills the ShapeIterator interface.
type ShapeIterator interface {
	// ForEach is a function that can iterate through a collection of Shapes, controlled by a function's return value.
	// If the function returns true, the iteration continues to the end. If it returns false, the iteration ends.
	ForEach(iterationFunction func(shape IShape) bool)
}

// ShapeFilter is a selection of Shapes, primarily used to filter them out to select only some (i.e. Shapes with specific tags or placement).
// Usually one would use a ShapeFilter to select Shapes that are near a moving Shape (e.g. the player character).
type ShapeFilter struct {
	Filters     []func(s IShape) bool
	operatingOn ShapeIterator
}

// ForEach is a function that can run a customizeable function on each Shape contained within the filter.
// If the shape passes the filters, the forEachFunc will run with the shape as an argument.
// If the function returns true, the iteration will continue; if it doesn't, the iteration will end.
func (s ShapeFilter) ForEach(forEachFunc func(shape IShape) bool) {

	s.operatingOn.ForEach(func(shape IShape) bool {

		pass := true

		for _, f := range s.Filters {
			if !f(shape) {
				pass = false
				return true
			}
		}

		if pass && !forEachFunc(shape) {
			return true
		}

		return true
	})

}

// ByTags adds a filter to the ShapeFilter that filters out Shapes by tags (so only Shapes that have the specified Tag(s) pass the filter).
// The function returns the ShapeFiler for easy method chaining.
func (s ShapeFilter) ByTags(tags Tags) ShapeFilter {
	s.Filters = append(s.Filters, func(s IShape) bool {
		return s.Tags().Has(tags)
	})
	return s
}

// ByDistance adds a filter to the ShapeFilter that filters out Shapes distance to a given point.
// The shapes have to be at least min and at most max distance from the given point Vector.
// The function returns the ShapeFiler for easy method chaining.
func (s ShapeFilter) ByDistance(point Vector, min, max float64) ShapeFilter {
	s.Filters = append(s.Filters, func(s IShape) bool {
		d := s.Position().Distance(point)
		return d > min && d < max
	})
	return s
}

// ByFunc adds a filter to the ShapeFilter that filters out Shapes using a function if it returns true, the Shape passes the ShapeFilter.
// The function returns the ShapeFiler for easy method chaining.
func (s ShapeFilter) ByFunc(filterFunc func(s IShape) bool) ShapeFilter {
	s.Filters = append(s.Filters, filterFunc)
	return s
}

// Not adds a filter to the ShapeFilter that specifcally does not allow specified Shapes in.
// The function returns the ShapeFiler for easy method chaining.
func (s ShapeFilter) Not(shapes ...IShape) ShapeFilter {
	s.Filters = append(s.Filters, func(s IShape) bool {
		for _, shape := range shapes {
			if shape == s {
				return false
			}
		}
		return true
	})
	return s
}

// Shapes returns all shapes that pass the filters as a ShapeCollection.
func (s ShapeFilter) Shapes() ShapeCollection {

	collection := ShapeCollection{}

	s.operatingOn.ForEach(func(shape IShape) bool {
		for _, filter := range s.Filters {
			if filter(shape) {
				collection = append(collection, shape)
			}
		}
		return true
	})

	return collection
}

// ShapeCollection is a slice of Shapes.
type ShapeCollection []IShape

// ForEach allows you to iterate through each shape in the ShapeCollection; if the function returns false, the iteration ends.
func (s ShapeCollection) ForEach(forEachFunc func(shape IShape) bool) {
	for _, shape := range s {
		if !forEachFunc(shape) {
			break
		}
	}
}

// SetTags sets the tag(s) on all Shapes present in the Shapecollection.
func (s ShapeCollection) SetTags(tags Tags) {
	for _, shape := range s {
		shape.Tags().Set(tags)
	}
}

// UnsetTags unsets the tag(s) on all Shapes present in the Shapecollection.
func (s ShapeCollection) UnsetTags(tags Tags) {
	for _, shape := range s {
		shape.Tags().Unset(tags)
	}
}
