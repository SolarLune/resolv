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
	Contacts     []*resolv.ContactSet
	Polygon      *resolv.ConvexPolygon
	PlayerCircle *resolv.Circle
	CircleTwo    *resolv.Circle
}

func NewWorldShapeTest(game *Game) *WorldShapeTest {
	world := &WorldShapeTest{Game: game}
	world.Init()
	return world
}

func (world *WorldShapeTest) Init() {

	world.Polygon = resolv.NewConvexPolygon(
		100, 100,
		250, 80,
		300, 150,
		250, 250,
		150, 300,
		80, 150,
	)
	world.PlayerCircle = resolv.NewCircle(500, 200, 32)
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

	world.Contacts = world.Contacts[:0]

	world.PlayerCircle.IntersectionForEach(
		0, 0,
		func(c *resolv.ContactSet) bool {

			if ebiten.IsKeyPressed(ebiten.KeySpace) {
				world.PlayerCircle.MoveVec(c.MTV)
			}

			world.Contacts = append(world.Contacts, c)
			return true
		},
		world.CircleTwo,
		world.Polygon,
	)

	world.PlayerCircle.Move(dx, dy)
}

func (world *WorldShapeTest) Draw(screen *ebiten.Image) {

	controllingColor := color.RGBA{0, 255, 80, 255}
	if len(world.Contacts) > 0 {
		controllingColor = color.RGBA{160, 0, 0, 255}
	}

	DrawPolygon(screen, world.Polygon, color.White)

	DrawCircle(screen, world.PlayerCircle, controllingColor)
	DrawCircle(screen, world.CircleTwo, color.White)

	for _, c := range world.Contacts {

		for _, p := range c.Points {
			DrawBigDot(screen, p, color.RGBA{255, 255, 0, 255})
		}

		ebitenutil.DrawLine(screen, c.Center.X, c.Center.Y, c.Center.X+c.MTV.X, c.Center.Y+c.MTV.Y, color.RGBA{255, 128, 0, 255})

		DrawBigDot(screen, c.Center, color.RGBA{255, 128, 255, 255})

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
