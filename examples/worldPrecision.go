package main

import (
	"fmt"
	"image/color"

	"github.com/SolarLune/resolv"
	"github.com/hajimehoshi/ebiten"
	"github.com/hajimehoshi/ebiten/ebitenutil"
	"github.com/hajimehoshi/ebiten/inpututil"
)

type WorldPrecision struct {
	Space        *resolv.Space
	Game         *Game
	Geometry     []*resolv.Object
	MovingRect   *resolv.Object
	ShowHelpText bool
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

	world.Space = resolv.NewSpace(int(gw)/cellSize, int(gh)/cellSize, cellSize, cellSize)

	world.Geometry = []*resolv.Object{
		resolv.NewObject(0, 0, 16, gh, world.Space),
		resolv.NewObject(gw-16, 0, 16, gh, world.Space),
		resolv.NewObject(0, 0, gw, 16, world.Space),
		resolv.NewObject(0, gh-24, gw, 32, world.Space),
		resolv.NewObject(320, 185, 19, 1600, world.Space),
	}

	world.MovingRect = resolv.NewObject(320, 32, 16, 16, world.Space)
	world.MovingRect.PreciseCollision = true

	// world.Bouncers = []*Bouncer{}

	// world.SpawnObject()
}

func (world *WorldPrecision) Update(screen *ebiten.Image) {

	if inpututil.IsKeyJustPressed(ebiten.KeyF2) {
		world.ShowHelpText = !world.ShowHelpText
	}

	if inpututil.IsKeyJustPressed(ebiten.KeySpace) {
		precise := !world.MovingRect.PreciseCollision
		world.Init()
		world.MovingRect.PreciseCollision = precise
	}

	dy := 2.0

	if col := world.MovingRect.Check(0, dy); col.Valid() {
		world.MovingRect.Y += col.ContactY
	} else {
		world.MovingRect.Y += dy
	}

	world.MovingRect.Update()

}

func (world *WorldPrecision) Draw(screen *ebiten.Image) {

	screen.Fill(color.RGBA{20, 20, 40, 255})

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
			"Space: Turn on or off precise collisions",
			"on the moving red rectangle.",
			"",
			fmt.Sprint("Precise Collision: ", world.MovingRect.PreciseCollision),
			"",
			"When colliding precisely, objects collide based on",
			"overlapping cellular locations AND objects' individual rectangles.",
			"Otherwise, objects collide based on",
			"overall cellular locations.",
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
