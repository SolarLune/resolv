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

// Colliding takes a checking Shape, and checks to see if any of the Shapes in the Space are colliding with it. If so, it adds it
// to a new Space, and returns it.
func (sp *Space) Colliding(shape Shape) Space {

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

// Resolve attempts to move the checking shape through space, returning a Collision object with its findings if it collides with
// any other shapes in the Space. xSpeed and ySpeed are the movement displacement in pixels for the frame. You should generally
// check for collision resolutions on the X and Y axes separately.
func (sp *Space) Resolve(checkingShape Shape, xSpeed, ySpeed float32) Collision {

	for _, other := range *sp {
		if other != checkingShape {
			res := checkingShape.Resolve(other, xSpeed, ySpeed)
			if res.Colliding() {
				return res
			}
		}
	}

	return Collision{}

}

// Filter filters out a Space, returning a new Space comprised of Shapes that return true for the boolean function you provide.
// This can be used to focus on a set of object for collision testing or resolution, or lower the number of Shapes to test
// by filtering some out beforehand.
func (sp Space) Filter(filterFunc func(Shape) bool) Space {
	subSpace := make(Space, 0)
	for _, shape := range sp {
		if filterFunc(shape) {
			subSpace.AddShape(shape)
		}
	}
	return subSpace
}

// FilterByTags filters a Space out, creating a new Space that has just the Shapes that have all of the specified tags.
func (sp Space) FilterByTags(tags ...string) Space {
	return sp.Filter(func(s Shape) bool {
		if s.HasTags(tags...) {
			return true
		}
		return false
	})
}

// Contains returns true if the Shape provided exists within the Space.
func (sp Space) Contains(shape Shape) bool {
	for _, s := range sp {
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
	IsCollideable() bool
	SetCollideable(bool)
	Resolve(Shape, float32, float32) Collision
	GetTags() []string
	SetTags(...string)
	HasTags(...string) bool
	GetData() interface{}
	SetData(interface{})
	SetXY(int32, int32)
	GetXY() (int32, int32)
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
func (b basicShape) GetTags() []string {
	return b.tags
}

// SetTags sets the tags on the Shape.
func (b *basicShape) SetTags(tags ...string) {
	b.tags = tags
}

// If the Shape has all of the tags provided.
func (b basicShape) HasTags(tags ...string) bool {

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
func (b basicShape) IsCollideable() bool {
	return b.Collideable
}

// SetCollideable sets the Shape's collide-ability.
func (b *basicShape) SetCollideable(on bool) {
	b.Collideable = on
}

// GetData returns the data on the Shape.
func (b basicShape) GetData() interface{} {
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
	// attempting to move along the direction given by xSpeed and ySpeed in the Resolve() function before touching another Shape.
	Teleporting bool
	// Teleporting is if moving according to ResolveX and ResolveY might be considered teleporting, which is moving greater than the
	// X or Yspeed provided to the Resolve function * 1.5 (this is arbitrary, but can be useful).
	OtherShape Shape
	// OtherShape should be a pointer to the Shape that the colliding object collided with.
}

// Colliding returns whether the Collision actually was valid because of a collision against another Shape.
func (c Collision) Colliding() bool {
	return c.OtherShape != nil
}

// resolve is a generic function to resolve the attempt of one shape to move up against another one. Individual Shapes' Resolve()
// functions just point to this function for ease of use.
func resolve(firstShape Shape, other Shape, xSpeed, ySpeed float32) Collision {

	out := Collision{}
	out.ResolveX, out.ResolveY = firstShape.GetXY()
	out.ResolveX += int32(xSpeed)
	out.ResolveY += int32(ySpeed)

	if !other.IsCollideable() || (xSpeed == 0 && ySpeed == 0) {
		return out
	}

	xv, yv := firstShape.GetXY()

	firstShape.SetXY(xv+int32(xSpeed), yv+int32(ySpeed))

	x := float32(xv) + xSpeed
	y := float32(yv) + ySpeed

	primeX := true
	var slope float32

	if ySpeed != 0 && xSpeed != 0 {
		slope = ySpeed / xSpeed
	}

	if math.Abs(float64(ySpeed)) > math.Abs(float64(xSpeed)) {
		primeX = false
		if ySpeed != 0 && xSpeed != 0 {
			slope = xSpeed / ySpeed
		}
	}

	colliding := true

	for colliding {

		if firstShape.IsColliding(other) {

			if primeX {

				if xSpeed > 0 {
					x--
				} else if xSpeed < 0 {
					x++
				}

				if ySpeed > 0 {
					y -= slope
				} else if ySpeed < 0 {
					y += slope
				}

			} else {

				if ySpeed > 0 {
					y--
				} else if ySpeed < 0 {
					y++
				}

				if xSpeed > 0 {
					x -= slope
				} else if xSpeed < 0 {
					x += slope
				}

			}

			out.ResolveX = int32(x)
			out.ResolveY = int32(y)

			firstShape.SetXY(out.ResolveX, out.ResolveY)

			out.OtherShape = other

		} else {
			colliding = false
		}

	}

	out.ResolveX -= xv
	out.ResolveY -= yv

	if math.Abs(float64(xSpeed-float32(out.ResolveX))) > math.Abs(float64(xSpeed*1.5)) || math.Abs(float64(ySpeed-float32(out.ResolveY))) > math.Abs(float64(ySpeed*1.5)) {
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
// the specified other Shape. The xSpeed and ySpeed arguments are the movement displacement in pixels (note that collision resolution
// operates on, naturally, whole pixels still). For most situations, you would want to resolve on the X and Y axes separately.
func (r *Rectangle) Resolve(other Shape, xSpeed, ySpeed float32) Collision {

	// Because the resolve function is the same for all the shapes, essentially
	rect := *r
	return resolve(&rect, other, xSpeed, ySpeed)

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
// the specified other Shape. The xSpeed and ySpeed arguments are the movement displacement in pixels (note that collision resolution
// operates on, naturally, whole pixels still). For most situations, you would want to resolve on the X and Y axes separately.
func (c *Circle) Resolve(other Shape, xSpeed, ySpeed float32) Collision {

	circle := *c
	return resolve(&circle, other, xSpeed, ySpeed)

}
