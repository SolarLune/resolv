package main

import (
	"image/color"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/solarlune/resolv"
)

type WorldShapeTest struct {
	Game *Game
	// Space *resolv.Space
	Contact   *resolv.ContactSet
	Solid     *resolv.ConvexPolygon
	CircleOne *resolv.Circle
	CircleTwo *resolv.Circle
}

func NewWorldShapeTest(game *Game) *WorldShapeTest {
	world := &WorldShapeTest{Game: game}
	world.Init()
	return world
}

func (world *WorldShapeTest) Init() {

	world.Solid = resolv.NewConvexPolygon(
		100, 100,
		250, 80,
		300, 150,
		250, 250,
		150, 300,
		80, 150,
	)
	world.CircleOne = resolv.NewCircle(500, 200, 32)
	world.CircleTwo = resolv.NewCircle(400, 250, 32)

}

func (world *WorldShapeTest) Update() {

	dx := 0.0
	dy := 0.0

	if ebiten.IsKeyPressed(ebiten.KeyRight) {
		dx += 1
	}

	if ebiten.IsKeyPressed(ebiten.KeyLeft) {
		dx -= 1
	}

	if ebiten.IsKeyPressed(ebiten.KeyUp) {
		dy -= 1
	}

	if ebiten.IsKeyPressed(ebiten.KeyDown) {
		dy += 1
	}

	if ebiten.IsKeyPressed(ebiten.KeySpace) && world.Contact != nil {
		world.CircleOne.MoveVec(world.Contact.MTV)
	}

	world.Contact = world.CircleOne.Intersection(0, 0, world.Solid)
	if world.Contact == nil {
		world.Contact = world.CircleOne.Intersection(0, 0, world.CircleTwo)
	}

	world.CircleOne.Move(dx, dy)
}

func (world *WorldShapeTest) Draw(screen *ebiten.Image) {

	controllingColor := color.RGBA{0, 255, 80, 255}
	if world.Contact != nil {
		controllingColor = color.RGBA{160, 0, 0, 255}
	}

	DrawPolygon(screen, world.Solid, color.White)

	DrawCircle(screen, world.CircleOne, controllingColor)
	DrawCircle(screen, world.CircleTwo, color.White)

	if world.Contact != nil {

		for _, p := range world.Contact.Points {
			DrawBigDot(screen, p.X(), p.Y(), color.RGBA{255, 255, 0, 255})
		}

		ebitenutil.DrawLine(screen, world.Contact.Center.X(), world.Contact.Center.Y(), world.Contact.Center.X()+world.Contact.MTV.X(), world.Contact.Center.Y()+world.Contact.MTV.Y(), color.RGBA{255, 128, 0, 255})

		DrawBigDot(screen, world.Contact.Center.X(), world.Contact.Center.Y(), color.RGBA{255, 128, 255, 255})

	}

	if world.Game.ShowHelpText {

		world.Game.DrawText(screen, 16, 16,
			"~World Shape Test~",
			"Move green Circle: Arrow keys",
			"Move along MTV (Minimum Translation Vector) to avoid collision: Space key",
			"",
			"The circle turns red when intersecting with another Shape.",
			"Yellow dots indicate contact points.",
			"The pink dot is the center of the contact points.",
			"The orange line indicates the MTV. This is how far the Shape",
			"must move in whatever direction to avoid intersection.",
			"This gives best results when not very far into another Shape.",
			"",
			"F2: Show / Hide help text",
			"F4: Toggle fullscreen",
			"R: Restart world",
			"E: Next world",
			"Q: Previous world",
		)

	}

}
