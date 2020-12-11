package resolv

// Shape is a basic interface that describes a Shape that can be passed to collision testing and resolution functions and
// exist in the same Space.
type Shape interface {
	IsColliding(Shape) bool
	WouldBeColliding(Shape, float64, float64) bool
	GetTags() []string
	ClearTags()
	AddTags(...string)
	RemoveTags(...string)
	HasTags(...string) bool
	GetData() interface{}
	SetData(interface{})
	GetXY() (float64, float64)
	SetXY(float64, float64)
	Move(float64, float64)
}

// BasicShape isn't to be used directly; it just has some basic functions and data, common to all structs that embed it, like
// position and tags. It is embedded in other Shapes.
type BasicShape struct {
	X, Y float64
	tags []string
	Data interface{}
}

// GetTags returns a reference to the the string array representing the tags on the BasicShape.
func (b *BasicShape) GetTags() []string {
	return b.tags
}

// AddTags adds the specified tags to the BasicShape.
func (b *BasicShape) AddTags(tags ...string) {
	if b.tags == nil {
		b.tags = []string{}
	}
	b.tags = append(b.tags, tags...)
}

// RemoveTags removes the specified tags from the BasicShape.
func (b *BasicShape) RemoveTags(tags ...string) {

	for _, t := range tags {

		for i := len(b.tags) - 1; i >= 0; i-- {

			if t == b.tags[i] {
				b.tags = append(b.tags[:i], b.tags[i+1:]...)
			}

		}

	}

}

// ClearTags clears the tags active on the BasicShape.
func (b *BasicShape) ClearTags() {
	b.tags = []string{}
}

// HasTags returns true if the Shape has all of the tags provided.
func (b *BasicShape) HasTags(tags ...string) bool {

	hasTags := true

	for _, t1 := range tags {
		found := false
		for _, shapeTag := range b.tags {
			if t1 == shapeTag {
				found = true
				continue
			}
		}
		if !found {
			hasTags = false
			break
		}
	}

	return hasTags
}

// GetData returns the data on the Shape.
func (b *BasicShape) GetData() interface{} {
	return b.Data
}

// SetData sets the data on the Shape.
func (b *BasicShape) SetData(data interface{}) {
	b.Data = data
}

// GetXY returns the position of the Shape.
func (b *BasicShape) GetXY() (float64, float64) {
	return b.X, b.Y
}

// SetXY sets the position of the Shape.
func (b *BasicShape) SetXY(x, y float64) {
	b.X = x
	b.Y = y
}

// Move moves the Shape by the delta X and Y values provided.
func (b *BasicShape) Move(x, y float64) {
	b.X += x
	b.Y += y
}
