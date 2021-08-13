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

	Triangle     *resolv.ConvexPolygon
	Intersection resolv.Delta
	Solid        *resolv.ConvexPolygon
}

func NewWorldShapeTest(game *Game) *WorldShapeTest {
	return &WorldShapeTest{
		Game: game,
		Triangle: resolv.NewConvexPolygon(
			-10, -10,
			10, 10,
			-10, 10,
		),
		Solid: resolv.NewConvexPolygon(
			100, 100,
			300, 80,
			350, 150,
			300, 300,
			200, 350,
			80, 150,
		),
		// Space: resolv.NewSpace(game.Width, game.Height, 16, 16)
	}
}

func (world *WorldShapeTest) Init() {

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

	if ebiten.IsKeyPressed(ebiten.KeySpace) {
		if world.Intersection.Valid {
			world.Triangle.Move(world.Intersection.Vector[0], world.Intersection.Vector[1])
		}
	}

	world.Intersection = world.Triangle.Intersection(world.Solid)

	world.Triangle.Move(dx, dy)

}

func (world *WorldShapeTest) Draw(screen *ebiten.Image) {

	triangleColor := color.RGBA{0, 255, 0, 255}
	if world.Intersection.Valid {
		triangleColor = color.RGBA{255, 0, 0, 255}
	}

	world.DrawPolygon(screen, world.Triangle, triangleColor)
	world.DrawPolygon(screen, world.Solid, color.RGBA{255, 255, 255, 255})

	world.Game.DrawText(screen, 16, 16,
		"~World Shape Test~",
	)
	// if world.Game.Debug {
	// }

}

func (world *WorldShapeTest) DrawPolygon(screen *ebiten.Image, shape *resolv.ConvexPolygon, color color.Color) {

	for i := 0; i < len(shape.Vertices); i++ {
		vert := shape.Vertices[i]
		next := shape.Vertices[0]

		if i < len(shape.Vertices)-1 {
			next = shape.Vertices[i+1]
		}
		ebitenutil.DrawLine(screen, vert.X, vert.Y, next.X, next.Y, color)

	}

}
