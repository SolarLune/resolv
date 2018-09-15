package resolv

import (
	"fmt"
	"math"
)

// Space represents a collection that holds Shapes for collision detection in the same common space. A Space is arbitrarily large -
// you can use one Space for a single level, room, or area in your game, or split it up if it makes more sense for your game design.
type Space []Shape

// NewSpace creates a new Space for shapes to exist in and be tested against in.
func NewSpace() Space {
	sp := Space{}
	sp = make(Space, 0)
	return sp
}

// AddShape adds the designated Shape to the Space collection.
func (sp *Space) AddShape(shape Shape) {
	*sp = append(*sp, shape)
}

// RemoveShape removes the designated Shape from the Space.
func (sp *Space) RemoveShape(shape Shape) {

	i := 0

	oldSpace := *sp

	for _, s := range oldSpace {
		if s == shape {
			break
		}
		i++
	}

	*sp = make(Space, 0)

	if len(oldSpace) > i+1 {
		*sp = append(*sp, oldSpace[:i]...)
		*sp = append(*sp, oldSpace[i+1:]...) // Don't forget to EXTRACT the elements in the list~
	} else {
		*sp = oldSpace[:i]
	}

}

// Clear "resets" the Space, cleaning out the Space of references to Shapes.
func (sp *Space) Clear() {
	*sp = make(Space, 0)
}

// IsColliding returns true if the designated Shape collides with another Shape in the Space with the specified tag.
func (sp *Space) IsColliding(shape Shape, tag string) bool {

	for _, other := range *sp {
		if other != shape && (tag == "" || other.GetTag() == tag) {
			if shape.IsColliding(other) {
				return true
			}
		}
	}
	return false

}

// Resolve attempts to move the checking shape through space, returning a Collision object with its findings if it collides with
// any other shapes in the Space. Speed is the movement displacement in pixels, and xAxis determines if the speed given should be
// tested on the X-axis, or Y-axis. This is done because for simple arcade-like movement, it tends to be a good idea to resolve
// collisions on the X and Y axes separately, one at a time, rather than all at once. The tag argument allows you to focus the
// function to look at only Shapes that have the specified tag.
// If the tag argument is a blank string, it will search all other Shapes.
func (sp *Space) Resolve(checkingShape Shape, speed float32, xAxis bool, tag string) Collision {

	for _, other := range *sp {
		if (tag == "" || other.GetTag() == tag) && other != checkingShape {
			res := checkingShape.Resolve(other, speed, xAxis)
			if res.Direction != CollisionNone {
				return res
			}
		}
	}

	return Collision{}

}

// Filter filters out a space, returning a "sub-space" of Shapes that return true for the boolean function you pass in that takes
// a Shape. Basically, you can use this to pick out specific Shapes from a single Space.
func (sp Space) Filter(filterFunc func(Shape) bool) Space {
	subSpace := make(Space, 0)
	for _, shape := range sp {
		if filterFunc(shape) {
			subSpace = append(subSpace, shape)
		}
	}
	return subSpace
}

func (sp *Space) String() string {
	str := ""
	for _, s := range *sp {
		str += fmt.Sprintf("%v   ", s)
	}
	return str
}

// CollisionDirection constants are to be used when checking the Collision object returned from Resolve functions.
type CollisionDirection int

func (cd CollisionDirection) String() string {

	switch cd {

	case CollisionDown:
		return "Collision Down"
	case CollisionUp:
		return "Collision Up"
	case CollisionRight:
		return "Collision Right"
	case CollisionLeft:
		return "Collision Left"
	case CollisionNone:
		return "No Collision"
	}

	return ""
}

// CollisionDirection constant definitions
const (
	CollisionNone CollisionDirection = iota
	CollisionDown
	CollisionUp
	CollisionRight
	CollisionLeft
)

// Collision describes the collision found when a Shape attempted to resolve a movement into another Shape, or in the same Space as
// other existing Shapes.
type Collision struct {
	Direction CollisionDirection
	// Direction is what direction a collision was encountered in, if any. One of the CollisionDirection constants.
	ResolveDistance int32
	// ResolveDistance is the distance to move to come into contact with the object.
	StartedFree bool
	// StartedFree is whether the Resolve() function was called with the calling object already in a non-colliding state.
}

// Colliding returns whether the Collision was valid; this is just checking to see if the Direction returned is not CollisionNone.
func (c Collision) Colliding() bool {
	return c.Direction != CollisionNone
}

// Shape is a basic interface that describes a Shape that can be passed to collision resolution functions and exist in the same Space.
type Shape interface {
	IsColliding(Shape) bool
	IsCollideable() bool
	SetCollideable(bool)
	Resolve(Shape, float32, bool) Collision
	GetTag() string
	SetTag(string)
	GetData() interface{}
	SetData(interface{})
	Move(int32, int32)
}

// basicShape isn't to be used; it just has some basic functions and data, common to all structs that embed it, like and position
// and collide-ability.
type basicShape struct {
	X, Y        int32
	Tag         string
	Collideable bool
	Data        interface{}
}

func (b basicShape) GetTag() string {
	return b.Tag
}

func (b *basicShape) SetTag(tag string) {
	b.Tag = tag
}

func (b basicShape) IsCollideable() bool {
	return b.Collideable
}

func (b *basicShape) SetCollideable(on bool) {
	b.Collideable = on
}

func (b basicShape) GetData() interface{} {
	return b.Data
}

func (b *basicShape) SetData(data interface{}) {
	b.Data = data
}

func (b *basicShape) Move(x, y int32) {
	b.X += x
	b.Y += y
}

// resolve is a generic function to resolve the attempt of one shape to move up against another one. Individual Shapes' Resolve()
// functions just point to this function for ease of use.
func resolve(firstShape Shape, other Shape, speed float32, xAxis bool) Collision {

	out := Collision{}

	if !other.IsCollideable() {
		return out
	}

	d := -1
	if speed > 0 {
		d = 1
	}

	for i := 0; i < int(math.Ceil(math.Abs(float64(speed))))+1; i++ {

		if xAxis {
			firstShape.Move(int32(d), 0)
		} else {
			firstShape.Move(0, int32(d))
		}

		out.ResolveDistance += int32(d)

		if firstShape.IsColliding(other) {

			if xAxis && d > 0 {
				out.Direction = CollisionRight
			} else if xAxis && d < 0 {
				out.Direction = CollisionLeft
			} else if !xAxis && d > 0 {
				out.Direction = CollisionDown
			} else {
				out.Direction = CollisionUp
			}

			if math.Abs(float64(out.ResolveDistance)) >= float64(d) {
				out.ResolveDistance -= int32(d)
			}

			break

		}

	}

	return out

}

// Rectangle represents a rectangle.
type Rectangle struct {
	basicShape
	W, H int32
}

// NewRectangle creates a new Rectangle and returns a pointer to it.
func NewRectangle(x, y, w, h int32) *Rectangle {
	r := &Rectangle{W: w, H: h}
	r.X = x
	r.Y = y
	r.Collideable = true
	return r
}

// IsColliding returns a boolean, indicating if the Rectangle is colliding with the specified other Shape or not.
func (r Rectangle) IsColliding(other Shape) bool {

	if !other.IsCollideable() {
		return false
	}

	b, ok := other.(*Rectangle)

	if ok {
		return r.X > b.X-r.W && r.Y > b.Y-r.H && r.X < b.X+b.W && r.Y < b.Y+b.H
	}

	c, ok := other.(*Circle)

	if ok {
		return c.IsColliding(&r)
	}

	fmt.Println("WARNING! Object ", other, " isn't a valid shape for collision testing against a Rectangle ", r, "!")

	return false
}

// IsZero returns whether the Rectangle has been initialized or not.
func (r Rectangle) IsZero() bool {
	return r.X == 0 && r.Y == 0 && r.W == 0 && r.H == 0
}

// Center returns the center point of the Rectangle.
func (r Rectangle) Center() (int32, int32) {

	x := r.X + r.W/2
	y := r.Y + r.H/2

	return x, y

}

// Resolve attempts to move the checking shape through space, returning a Collision object with its findings if it collides with
// the specified other Shape. Speed is the movement displacement in pixels, and xAxis determines if the speed given should be
// tested on the X-axis, or Y-axis. This is done because for simple arcade-like movement, it tends to be a good idea to resolve
// collisions on the X and Y axes separately, one at a time, rather than all at once.
func (r *Rectangle) Resolve(other Shape, speed float32, xAxis bool) Collision {

	// Because the resolve function is the same for all the shapes, essentially
	rect := *r
	return resolve(&rect, other, speed, xAxis)

}

// A Circle represents an ordinary circle, and has a radius, in addition to normal shape properties.
type Circle struct {
	basicShape
	Radius int32
}

// NewCircle returns a pointer to a new Circle object.
func NewCircle(x, y, radius int32) *Circle {
	c := &Circle{Radius: radius}
	c.X = x
	c.Y = y
	c.Collideable = true
	return c
}

func distance(x, y, x2, y2 int32) int32 {

	dx := x - x2
	dy := y - y2
	ds := (dx * dx) + (dy * dy)
	return int32(math.Abs(float64(ds)))

}

// IsColliding returns true if the Circle is colliding with the specified other Shape.
func (c Circle) IsColliding(other Shape) bool {

	if !other.IsCollideable() {
		return false
	}

	b, ok := other.(*Circle)

	if ok {

		return distance(c.X, c.Y, b.X, b.Y) <= (c.Radius+b.Radius)*(c.Radius+b.Radius)

	}

	r, ok := other.(*Rectangle)

	if ok {

		closestX := c.X
		closestY := c.Y

		if c.X < r.X {
			closestX = r.X
		} else if c.X > r.X+r.W {
			closestX = r.X + r.W
		}

		if c.Y < r.Y {
			closestY = r.Y
		} else if c.Y > r.Y+r.H {
			closestY = r.Y + r.H
		}

		return distance(c.X, c.Y, closestX, closestY) <= c.Radius*c.Radius

	}

	fmt.Println("WARNING! Object ", other, " isn't a valid shape for collision testing against Circle ", c, "!")

	return false

}

// Resolve attempts to move the checking shape through space, returning a Collision object with its findings if it collides with
// the specified other Shape. Speed is the movement displacement in pixels, and xAxis determines if the speed given should be
// tested on the X-axis, or Y-axis. This is done because for simple arcade-like movement, it tends to be a good idea to resolve
// collisions on the X and Y axes separately, one at a time, rather than all at once.
func (c *Circle) Resolve(other Shape, speed float32, xAxis bool) Collision {

	circle := *c
	return resolve(&circle, other, speed, xAxis)

}
