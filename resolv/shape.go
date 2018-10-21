package resolv

// Shape is a basic interface that describes a Shape that can be passed to collision testing and resolution functions and
// exist in the same Space.
type Shape interface {
	IsColliding(Shape) bool
	WouldBeColliding(Shape, int32, int32) bool
	IsCollideable() bool
	SetCollideable(bool)
	GetTags() []string
	SetTags(...string)
	HasTags(...string) bool
	GetData() interface{}
	SetData(interface{})
	GetXY() (int32, int32)
	SetXY(int32, int32)
	Move(int32, int32)
}

// BasicShape isn't to be used directly; it just has some basic functions and data, common to all structs that embed it, like
// position and collide-ability. It is embedded in other Shapes.
type BasicShape struct {
	X, Y        int32
	tags        []string
	Collideable bool
	Data        interface{}
}

// GetTags returns the tags on the Shape.
func (b *BasicShape) GetTags() []string {
	return b.tags
}

// SetTags sets the tags on the Shape.
func (b *BasicShape) SetTags(tags ...string) {
	b.tags = tags
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

// IsCollideable returns whether the Shape is currently collide-able or not.
func (b *BasicShape) IsCollideable() bool {
	return b.Collideable
}

// SetCollideable sets the Shape's collide-ability.
func (b *BasicShape) SetCollideable(on bool) {
	b.Collideable = on
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
func (b *BasicShape) GetXY() (int32, int32) {
	return b.X, b.Y
}

// SetXY sets the position of the Shape.
func (b *BasicShape) SetXY(x, y int32) {
	b.X = x
	b.Y = y
}

// Move moves the Shape by the delta X and Y values provided.
func (b *BasicShape) Move(x, y int32) {
	b.X += x
	b.Y += y
}
