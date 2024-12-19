package main

import (
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/solarlune/resolv"
)

type WorldCircle struct {
	space  *resolv.Space
	Circle *Circle
}

func NewWorldCircle() *WorldCircle {
	w := &WorldCircle{}
	w.Init()
	return w
}

func (w *WorldCircle) Init() {

	if w.space != nil {
		w.space.RemoveAll()
	}

	// Create the space. It is 640x360 large (the size of the screen), and divided into 16x16 cells.
	// The cell division makes it more efficient to check for shapes.
	w.space = resolv.NewSpace(640, 360, 16, 16)

	solids := resolv.ShapeCollection{
		resolv.NewRectangleFromTopLeft(0, 0, 640, 16),
		resolv.NewRectangleFromTopLeft(0, 360-16, 640, 16),
		resolv.NewRectangleFromTopLeft(0, 16, 16, 360-16),
		resolv.NewRectangleFromTopLeft(640-16, 16, 16, 360-16),
		resolv.NewRectangleFromTopLeft(64, 128, 16, 200),
		resolv.NewRectangleFromTopLeft(120, 300, 200, 8),

		resolv.NewLine(256-32, 180-32, 256+32, 180+32),
		resolv.NewLine(256-32, 180+32, 256+32, 180-32),

		resolv.NewCircle(128, 128, 64),
	}

	solids.SetTags(TagSolidWall | TagPlatform)

	w.space.Add(solids...)

	w.Circle = NewCircle(w)

}

func (w *WorldCircle) Update() {
	w.Circle.Update()
}

func (w *WorldCircle) Draw(screen *ebiten.Image) {
	CommonDraw(screen, w)
	if GlobalGame.ShowHelpText {
		GlobalGame.DrawText(screen, 0, 128,
			"Circle Movement Test",
			"Arrow keys to move",
		)
	}
}

// To allow the world's physical state to be drawn using the debug draw function.
func (w *WorldCircle) Space() *resolv.Space {
	return w.space
}

type Circle struct {
	Object   *resolv.Circle
	Movement resolv.Vector
}

func NewCircle(world *WorldCircle) *Circle {

	circle := &Circle{
		Object: resolv.NewCircle(320, 64, 8),
	}
	circle.Object.Tags().Set(TagPlayer)
	circle.Object.SetData(circle)

	world.space.Add(circle.Object)
	return circle

}

func (c *Circle) Update() {

	movement := resolv.NewVectorZero()
	maxSpd := 4.0
	friction := 0.5
	accel := 0.5 + friction

	if ebiten.IsKeyPressed(ebiten.KeyLeft) {
		movement.X -= 1
	}
	if ebiten.IsKeyPressed(ebiten.KeyRight) {
		movement.X += 1
	}

	if ebiten.IsKeyPressed(ebiten.KeyUp) {
		movement.Y -= 1
	}
	if ebiten.IsKeyPressed(ebiten.KeyDown) {
		movement.Y += 1
	}

	c.Movement = c.Movement.Add(movement.Scale(accel)).SubMagnitude(friction).ClampMagnitude(maxSpd)

	c.Object.MoveVec(c.Movement)

	c.Object.IntersectionTest(resolv.IntersectionTestSettings{
		TestAgainst: c.Object.SelectTouchingCells(1).FilterShapes(),
		OnIntersect: func(set resolv.IntersectionSet) bool {
			c.Object.MoveVec(set.MTV)
			return true
		},
	})

}
