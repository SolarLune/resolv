package resolv

// Shape is a basic interface that describes a Shape that can be passed to collision testing and resolution functions and
// exist in the same Space.
type Shape interface {
	IsColliding(other Shape) bool
	Check(other Shape, dx, dy float64) *MovementCheck
	Tags() *Tags
	Data() interface{}
	SetData(interface{})
	Position() (float64, float64)
	SetPosition(float64, float64)
	Move(float64, float64)
}

// Tags contains string "tags" that identify and classify Shapes.
type Tags struct {
	names []string
}

// NewTags creates a new Tags object.
func NewTags() *Tags {
	tags := &Tags{}
	tags.Clear()
	return tags
}

// Add adds the specified tag to the Tags object.
func (tags *Tags) Add(tag string) {
	tags.names = append(tags.names, tag)
}

// Remove removes the specified tag from the Tags object.
func (tags *Tags) Remove(tag string) {

	for i := len(tags.names) - 1; i >= 0; i-- {

		if tag == tags.names[i] {
			tags.names = append(tags.names[:i], tags.names[i+1:]...)
		}

	}

}

// Has returns true if the Tags object has the specified Tag.
func (tags *Tags) Has(tag string) bool {

	for _, shapeTag := range tags.names {
		if tag == shapeTag {
			return true
		}
	}

	return false
}

func (tags *Tags) Clear() {
	tags.names = []string{}
}

// BasicShape isn't to be used directly; it just has some basic functions and data, common to all structs that embed it, like
// position and tags. It is embedded in other Shapes.
type BasicShape struct {
	X, Y float64
	tags *Tags
	data interface{}
}

// Tags returns the Tags object on the Shape, used to add or change tags on the Shape.
func (b *BasicShape) Tags() *Tags {
	return b.tags
}

// Data returns the custom user data on the Shape.
func (b *BasicShape) Data() interface{} {
	return b.data
}

// SetData sets the custom user data on the Shape.
func (b *BasicShape) SetData(data interface{}) {
	b.data = data
}

// Position returns the position of the Shape.
func (b *BasicShape) Position() (float64, float64) {
	return b.X, b.Y
}

// SetPosition sets the position of the Shape.
func (b *BasicShape) SetPosition(x, y float64) {
	b.X = x
	b.Y = y
}

// Move moves the Shape by the delta X and Y values provided.
func (b *BasicShape) Move(x, y float64) {
	b.X += x
	b.Y += y
}
