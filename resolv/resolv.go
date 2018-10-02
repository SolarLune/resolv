/*
Package resolv is a simple collision detection and resolution library. Its goal is to be lightweight, fast, simple, and easy-to-use
for game development. Its goal is to also to not become a physics engine or physics library itself, but to always leave the actual
physics implementation and "game feel" to the developer, while making it very easy to do so.

Usage of resolv essentially centers around two main concepts: Spaces and Shapes.

A Shape can be used to test for collisions against another Shape. That's really all they have to do, but that capability is powerful
when paired with the resolv.Resolve() function. You can then check to see if a Shape would have a collision if it attempted to move
in a specified direction. If so, the Resolve() function would return a Collision object, which tells you some information about the
Collision, like how far the checking Shape would have to move to come into contact with the other, and which Shape it comes into
contact with.

A Space is just a slice that holds Shapes for detection. It doesn't represent any real physical space, and so there aren't any
units of measurement to remember when using Spaces. Similar to Shapes, Spaces are simple, but also very powerful. Spaces allow
you to easily check for collision with, and resolve collision against multiple Shapes within that Space. A Space being just a
collection of Shapes means that you can manipulate and filter them as necessary.
*/
package resolv

import (
	"fmt"
	"math"
)

// Space represents a collection that holds Shapes for collision detection in the same common space. A Space is arbitrarily large -
// you can use one Space for a single level, room, or area in your game, or split it up if it makes more sense for your game design.
// Technically, a Space is just a slice of Shapes.
type Space []Shape

// NewSpace creates a new Space for shapes to exist in and be tested against in.
func NewSpace() Space {
	sp := Space{}
	sp = make(Space, 0)
	return sp
}

// AddShape adds the designated Shapes to the Space.
func (sp *Space) AddShape(shapes ...Shape) {
	*sp = append(*sp, shapes...)
}

// RemoveShape removes the designated Shapes from the Space.
func (sp *Space) RemoveShape(shapes ...Shape) {

	for _, shape := range shapes {

		for deleteIndex, s := range *sp {

			if s == shape {
				s := *sp
				s[deleteIndex] = nil
				s = append(s[:deleteIndex], s[deleteIndex+1:]...)
				*sp = s
				break
			}

		}

	}

}

// Clear "resets" the Space, cleaning out the Space of references to Shapes.
func (sp *Space) Clear() {
	*sp = make(Space, 0)
}

// IsColliding returns whether the provided Shape is colliding with something in this Space.
func (sp *Space) IsColliding(shape Shape) bool {

	for _, other := range *sp {

		if other != shape {

			if shape.IsColliding(other) {
				return true
			}

		}

	}

	return false

}

// GetCollidingShapes returns a Space comprised of Shapes that collide with the checking Shape.
func (sp *Space) GetCollidingShapes(shape Shape) Space {

	newSpace := Space{}

	for _, other := range *sp {
		if other != shape {
			if shape.IsColliding(other) {
				newSpace = append(newSpace, other)
			}
		}
	}

	return newSpace

}

// Resolve runs Resolve() using the checking Shape, checking against all other Shapes in the Space. The first Collision
// that returns true is the Collision that gets returned.
func (sp *Space) Resolve(checkingShape Shape, deltaX, deltaY int32) Collision {

	res := Collision{}

	for _, other := range *sp {

		if other != checkingShape && checkingShape.WouldBeColliding(other, int32(deltaX), int32(deltaY)) {
			res = Resolve(checkingShape, other, deltaX, deltaY)
			if res.Colliding() {
				break
			}
		}

	}

	return res

}

// Filter filters out a Space, returning a new Space comprised of Shapes that return true for the boolean function you provide.
// This can be used to focus on a set of object for collision testing or resolution, or lower the number of Shapes to test
// by filtering some out beforehand.
func (sp *Space) Filter(filterFunc func(Shape) bool) Space {
	subSpace := make(Space, 0)
	for _, shape := range *sp {
		if filterFunc(shape) {
			subSpace.AddShape(shape)
		}
	}
	return subSpace
}

// FilterByTags filters a Space out, creating a new Space that has just the Shapes that have all of the specified tags.
func (sp *Space) FilterByTags(tags ...string) Space {
	return sp.Filter(func(s Shape) bool {
		if s.HasTags(tags...) {
			return true
		}
		return false
	})
}

// Contains returns true if the Shape provided exists within the Space.
func (sp *Space) Contains(shape Shape) bool {
	for _, s := range *sp {
		if s == shape {
			return true
		}
	}
	return false
}

func (sp *Space) String() string {
	str := ""
	for _, s := range *sp {
		str += fmt.Sprintf("%v   ", s)
	}
	return str
}

// Shape is a basic interface that describes a Shape that can be passed to collision resolution functions and exist in the same
// Space.
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
}

// basicShape isn't to be used; it just has some basic functions and data, common to all structs that embed it, like and position
// and collide-ability.
type basicShape struct {
	X, Y        int32
	tags        []string
	Collideable bool
	Data        interface{}
}

// GetTags returns the tags on the Shape.
func (b *basicShape) GetTags() []string {
	return b.tags
}

// SetTags sets the tags on the Shape.
func (b *basicShape) SetTags(tags ...string) {
	b.tags = tags
}

// If the Shape has all of the tags provided.
func (b *basicShape) HasTags(tags ...string) bool {

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
func (b *basicShape) IsCollideable() bool {
	return b.Collideable
}

// SetCollideable sets the Shape's collide-ability.
func (b *basicShape) SetCollideable(on bool) {
	b.Collideable = on
}

// GetData returns the data on the Shape.
func (b *basicShape) GetData() interface{} {
	return b.Data
}

// SetData sets the data on the Shape.
func (b *basicShape) SetData(data interface{}) {
	b.Data = data
}

// GetXY returns the position of the Shape.
func (b *basicShape) GetXY() (int32, int32) {
	return b.X, b.Y
}

// SetXY sets the position of the Shape.
func (b *basicShape) SetXY(x, y int32) {
	b.X = x
	b.Y = y
}

// Collision describes the collision found when a Shape attempted to resolve a movement into another Shape, or in the same Space as
// other existing Shapes.
type Collision struct {
	ResolveX, ResolveY int32
	// ResolveX and ResolveY represent the displacement of the Shape to the point of collision. How far along the Shape got when
	// attempting to move along the direction given by deltaX and deltaY in the Resolve() function before touching another Shape.
	Teleporting bool
	// Teleporting is if moving according to ResolveX and ResolveY might be considered teleporting, which is moving greater than the
	// X or deltaY provided to the Resolve function * 1.5 (this is arbitrary, but can be useful).
	OtherShape Shape
	// OtherShape should be a pointer to the Shape that the colliding object collided with.
}

// Colliding returns whether the Collision actually was valid because of a collision against another Shape.
func (c Collision) Colliding() bool {
	return c.OtherShape != nil
}

// Resolve attempts to move the checking Shape with the specified X and Y values, returning a Collision object if it collides with
// the specified other Shape. The deltaX and deltaY arguments are the movement displacement in pixels. For most situations, you
// would want to resolve on the X and Y axes separately.
func Resolve(firstShape Shape, other Shape, deltaX, deltaY int32) Collision {

	out := Collision{}
	out.ResolveX = deltaX
	out.ResolveY = deltaY

	if !firstShape.IsCollideable() || !other.IsCollideable() || (deltaX == 0 && deltaY == 0) {
		return out
	}

	x := float32(deltaX)
	y := float32(deltaY)

	primeX := true
	slope := float32(0)

	if deltaY != 0 && deltaX != 0 {
		slope = float32(deltaY) / float32(deltaX)
	}

	if math.Abs(float64(deltaY)) > math.Abs(float64(deltaX)) {
		primeX = false
		if deltaY != 0 && deltaX != 0 {
			slope = float32(deltaX) / float32(deltaY)
		}
	}

	for true {

		if firstShape.WouldBeColliding(other, out.ResolveX, out.ResolveY) {

			if primeX {

				if deltaX > 0 {
					x--
				} else if deltaX < 0 {
					x++
				}

				if deltaY > 0 {
					y -= slope
				} else if deltaY < 0 {
					y += slope
				}

			} else {

				if deltaY > 0 {
					y--
				} else if deltaY < 0 {
					y++
				}

				if deltaX > 0 {
					x -= slope
				} else if deltaX < 0 {
					x += slope
				}

			}

			out.ResolveX = int32(x)
			out.ResolveY = int32(y)

			out.OtherShape = other

		} else {
			break
		}

	}

	if math.Abs(float64(deltaX-out.ResolveX)) > math.Abs(float64(deltaX)*1.5) || math.Abs(float64(deltaY-out.ResolveY)) > math.Abs(float64(deltaY)*1.5) {
		out.Teleporting = true
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

// IsColliding returns whether the Rectangle is colliding with the specified other Shape or not.
func (r *Rectangle) IsColliding(other Shape) bool {

	if !r.Collideable || !other.IsCollideable() {
		return false
	}

	b, ok := other.(*Rectangle)

	if ok {
		return r.X > b.X-r.W && r.Y > b.Y-r.H && r.X < b.X+b.W && r.Y < b.Y+b.H
	}

	c, ok := other.(*Circle)

	if ok {
		return c.IsColliding(r)
	}

	fmt.Println("WARNING! Object ", other, " isn't a valid shape for collision testing against a Rectangle ", r, "!")

	return false
}

// WouldBeColliding returns whether the Rectangle would be colliding with the other Shape if it were to move in the
// specified direction.
func (r *Rectangle) WouldBeColliding(other Shape, dx, dy int32) bool {
	r.X += dx
	r.Y += dy
	isColliding := r.IsColliding(other)
	r.X -= dx
	r.Y -= dy
	return isColliding
}

// IsZero returns whether the Rectangle has been initialized or not.
func (r *Rectangle) IsZero() bool {
	return r.X == 0 && r.Y == 0 && r.W == 0 && r.H == 0
}

// Center returns the center point of the Rectangle.
func (r *Rectangle) Center() (int32, int32) {

	x := r.X + r.W/2
	y := r.Y + r.H/2

	return x, y

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

// IsColliding returns true if the Circle is colliding with the specified other Shape.
func (c *Circle) IsColliding(other Shape) bool {

	if !c.Collideable || !other.IsCollideable() {
		return false
	}

	b, ok := other.(*Circle)

	if ok {

		return Distance(c.X, c.Y, b.X, b.Y) <= (c.Radius+b.Radius)*(c.Radius+b.Radius)

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

		return Distance(c.X, c.Y, closestX, closestY) <= c.Radius*c.Radius

	}

	fmt.Println("WARNING! Object ", other, " isn't a valid shape for collision testing against Circle ", c, "!")

	return false

}

// WouldBeColliding returns whether the Rectangle would be colliding with the specified other Shape if it were to move
// in the specified direction.
func (c *Circle) WouldBeColliding(other Shape, dx, dy int32) bool {
	c.X += dx
	c.Y += dy
	isColliding := c.IsColliding(other)
	c.X -= dx
	c.Y -= dy
	return isColliding
}

// GetBoundingRect returns a Rectangle which has a width and height of 2*Radius.
func (c *Circle) GetBoundingRect() *Rectangle {
	r := &Rectangle{}
	r.W = c.Radius * 2
	r.H = c.Radius * 2
	r.X = c.X - r.W/2
	r.Y = c.Y - r.H/2
	return r
}

// Distance returns the distance from one pair of X and Y values to another.
func Distance(x, y, x2, y2 int32) int32 {

	dx := x - x2
	dy := y - y2
	ds := (dx * dx) + (dy * dy)
	return int32(math.Abs(float64(ds)))

}
