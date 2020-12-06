package main

import (
	"fmt"
	"image/color"
	"math/rand"

	"github.com/kvartborg/vector"

	"github.com/SolarLune/resolv"
	"github.com/hajimehoshi/ebiten"
	"github.com/hajimehoshi/ebiten/ebitenutil"
	"github.com/hajimehoshi/ebiten/inpututil"
)

type WorldBouncer struct {
	Game         *Game
	Space        *resolv.Space
	Geometry     []*resolv.Object
	Bouncers     []*Bouncer
	ShowHelpText bool
}

type Bouncer struct {
	Object *resolv.Object
	Speed  vector.Vector
}

func NewWorldBouncer(game *Game) *WorldBouncer {

	w := &WorldBouncer{
		Game:         game,
		ShowHelpText: true,
	}

	w.Init()

	return w
}

func (world *WorldBouncer) Init() {

	gw := float64(world.Game.Width)
	gh := float64(world.Game.Height)
	cellSize := 8

	world.Space = resolv.NewSpace(int(gw)/cellSize, int(gh)/cellSize, cellSize, cellSize)

	world.Geometry = []*resolv.Object{
		resolv.NewObject(0, 0, 16, gh, world.Space),
		resolv.NewObject(gw-16, 0, 16, gh, world.Space),
		resolv.NewObject(0, 0, gw, 16, world.Space),
		resolv.NewObject(0, gh-24, gw, 32, world.Space),
	}

	world.Bouncers = []*Bouncer{}

	world.SpawnObject()

}

func (world *WorldBouncer) SpawnObject() {

	bouncer := &Bouncer{
		Object: resolv.NewObject(0, 0, 2, 2, world.Space),
		Speed:  vector.Vector{(rand.Float64() * 8) - 4, (rand.Float64() * 8) - 4},
	}

	var c *resolv.Cell
	for c == nil {
		rx := rand.Intn(world.Space.Width())
		ry := rand.Intn(world.Space.Height())
		c = world.Space.Cell(rx, ry)
		if c.Occupied() {
			c = nil
		} else {
			bouncer.Object.X, bouncer.Object.Y = world.Space.SpaceToWorld(c.X, c.Y)
		}
	}

	world.Bouncers = append(world.Bouncers, bouncer)

}

func (world *WorldBouncer) Update(screen *ebiten.Image) {

	if inpututil.IsKeyJustPressed(ebiten.KeyF2) {
		world.ShowHelpText = !world.ShowHelpText
	}

	if ebiten.IsKeyPressed(ebiten.KeyUp) {
		world.SpawnObject()
	} else if ebiten.IsKeyPressed(ebiten.KeyDown) {
		if len(world.Bouncers) > 0 {
			b := world.Bouncers[0]
			b.Object.Remove()
			world.Bouncers = world.Bouncers[1:]
		}
	}

	for _, b := range world.Bouncers {

		b.Speed[1] += 0.1

		if check := b.Object.Check(b.Speed[0], 0); check.Valid() {
			b.Speed[0] *= -1
			b.Object.X += check.ContactX
		} else {
			b.Object.X += b.Speed[0]
		}

		if check := b.Object.Check(0, b.Speed[1]); check.Valid() {
			b.Speed[1] *= -1
			b.Object.Y += check.ContactY
		} else {
			b.Object.Y += b.Speed[1]
		}

		b.Object.Update()

	}

}

func (world *WorldBouncer) Draw(screen *ebiten.Image) {

	screen.Fill(color.RGBA{20, 20, 40, 255})

	for _, o := range world.Geometry {
		ebitenutil.DrawRect(screen, o.X, o.Y, o.W, o.H, color.RGBA{60, 60, 60, 255})
	}

	for _, b := range world.Bouncers {
		o := b.Object
		ebitenutil.DrawRect(screen, o.X, o.Y, o.W, o.H, color.RGBA{0, 80, 255, 255})
	}

	if world.Game.Debug {
		world.Game.DebugDraw(screen, world.Space)
	}

	if world.ShowHelpText {

		world.Game.DrawText(screen, 16, 16,
			"~ Bouncer Demo ~",
			"Up Arrow: Add bouncer",
			"Down Arrow: Remove bouncer",
			"",
			"F1: Toggle Debug View",
			"F2: Show / Hide help text",
			"R: Restart world",
			"E: Next world",
			"Q: Previous world",
			fmt.Sprintf("%d Bouncers in the world.", len(world.Bouncers)),
			fmt.Sprintf("%d FPS (frames per second)", int(ebiten.CurrentFPS())),
			fmt.Sprintf("%d TPS (ticks per second)", int(ebiten.CurrentTPS())),
		)

	}

}
