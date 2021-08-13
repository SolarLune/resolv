package main

import (
	"fmt"
	"image/color"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
	"github.com/solarlune/resolv"
)

type WorldPrecision struct {
	Space              *resolv.Space
	Game               *Game
	Geometry           []*resolv.Object
	MovingRect         *resolv.Object
	ShowHelpText       bool
	PreciseCollisionOn bool
}

func NewWorldPrecision(g *Game) *WorldPrecision {
	w := &WorldPrecision{Game: g, ShowHelpText: true}
	w.Init()
	return w
}

func (world *WorldPrecision) Init() {

	gw := float64(world.Game.Width)
	gh := float64(world.Game.Height)
	cellSize := 32

	world.Space = resolv.NewSpace(int(gw), int(gh), cellSize, cellSize)

	world.Geometry = []*resolv.Object{
		resolv.NewObject(0, 0, 16, gh),
		resolv.NewObject(gw-16, 0, 16, gh),
		resolv.NewObject(0, 0, gw, 16),
		resolv.NewObject(0, gh-24, gw, 32),
		resolv.NewObject(320, 185, 19, 1600),
	}

	world.Space.Add(world.Geometry...)

	world.MovingRect = resolv.NewObject(320, 32, 16, 16)
	world.Space.Add(world.MovingRect)
	// world.MovingRect.PreciseCollision = true

	// world.Bouncers = []*Bouncer{}

	// world.SpawnObject()
}

func (world *WorldPrecision) Update() {

	if inpututil.IsKeyJustPressed(ebiten.KeyF2) {
		world.ShowHelpText = !world.ShowHelpText
	}

	if inpututil.IsKeyJustPressed(ebiten.KeySpace) {
		world.PreciseCollisionOn = !world.PreciseCollisionOn
		world.Init() // Restart the world
	}

	dy := 2.0

	if col := world.MovingRect.Check(0, dy); col.Valid() {

		if world.PreciseCollisionOn {

			if other := col.Objects()[0].ToRectangle().Intersection(world.MovingRect.ToRectangle()); other.Valid {
				dy = -other.Vector[1]
			}

		} else {
			dy = col.Contact.Vector[1]
		}

	}

	world.MovingRect.Y += dy

	world.MovingRect.Update()

}

func (world *WorldPrecision) Draw(screen *ebiten.Image) {

	for _, o := range world.Geometry {
		ebitenutil.DrawRect(screen, o.X, o.Y, o.W, o.H, color.RGBA{60, 60, 60, 255})
	}

	o := world.MovingRect
	ebitenutil.DrawRect(screen, o.X, o.Y, o.W, o.H, color.RGBA{255, 60, 20, 255})

	if world.Game.Debug {
		world.Game.DebugDraw(screen, world.Space)
	}

	if world.ShowHelpText {

		world.Game.DrawText(screen, 16, 16,
			"~ Precision Demo ~",
			"",
			fmt.Sprint("Precise Collision: ", world.PreciseCollisionOn),
			"",
			"For precise collisions, one can use resolv's built-in shape-checking functions ",
			"to see if a true collision happens between objects within the same Cell.",
			"Otherwise, objects collide based on overall cellular locations.",
			"",
			"Space: Turn on or off precise collisions",
			"on the moving red rectangle.",
			"F1: Toggle Debug View",
			"F2: Show / Hide help text",
			"R: Restart world",
			"E: Next world",
			"Q: Previous world",
			fmt.Sprintf("%d FPS (frames per second)", int(ebiten.CurrentFPS())),
			fmt.Sprintf("%d TPS (ticks per second)", int(ebiten.CurrentTPS())),
		)

	}

}
