package main

import (
	"fmt"
	"image/color"
	"math/rand"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/solarlune/resolv"
)

type WorldBouncer struct {
	Game         *Game
	Space        *resolv.Space
	Geometry     []*resolv.Object
	Bouncers     []*Bouncer
	MaxBouncers  int
	ShowHelpText bool
}

type Bouncer struct {
	Object *resolv.Object
	Speed  resolv.Vector
}

func NewWorldBouncer(game *Game) *WorldBouncer {

	w := &WorldBouncer{
		Game:         game,
		ShowHelpText: true,
		MaxBouncers:  3000,
	}

	w.Init()

	return w
}

func (world *WorldBouncer) Init() {

	gw := float64(world.Game.Width)
	gh := float64(world.Game.Height)
	cellSize := 8

	world.Space = resolv.NewSpace(int(gw), int(gh), cellSize, cellSize)

	world.Geometry = []*resolv.Object{
		resolv.NewObject(0, 0, 16, gh),
		resolv.NewObject(gw-16, 0, 16, gh),
		resolv.NewObject(0, 0, gw, 16),
		resolv.NewObject(0, gh-24, gw, 32),
	}

	world.Space.Add(world.Geometry...)

	world.Bouncers = []*Bouncer{}

	world.SpawnObject()

}

func (world *WorldBouncer) SpawnObject() {

	bouncer := &Bouncer{
		Object: resolv.NewObject(0, 0, 2, 2),
		Speed: resolv.NewVector(
			(rand.Float64()*8)-4,
			(rand.Float64()*8)-4,
		),
	}

	world.Space.Add(bouncer.Object)

	// Choose an unoccupied cell to spawn a bouncing object in
	var c *resolv.Cell
	for c == nil {
		rx := rand.Intn(world.Space.Width())
		ry := rand.Intn(world.Space.Height())
		c = world.Space.Cell(rx, ry)
		if c.Occupied() {
			c = nil
		} else {
			bouncer.Object.Position.X, bouncer.Object.Position.Y = world.Space.SpaceToWorld(c.X, c.Y)
		}
	}

	world.Bouncers = append(world.Bouncers, bouncer)

}

func (world *WorldBouncer) Update() {

	if ebiten.IsKeyPressed(ebiten.KeyUp) {
		for i := 0; i < 5; i++ {
			if len(world.Bouncers)+1 < world.MaxBouncers {
				world.SpawnObject()
			}
		}
	} else if ebiten.IsKeyPressed(ebiten.KeyDown) {
		if len(world.Bouncers) > 0 {
			b := world.Bouncers[0]
			world.Space.Remove(b.Object)
			world.Bouncers = world.Bouncers[1:]
		}
	}

	for _, b := range world.Bouncers {

		b.Speed.Y += 0.1

		dx := b.Speed.X
		dy := b.Speed.Y

		if check := b.Object.Check(dx, 0); check != nil {
			// We move a bouncer into contact with the owning cell rather than the object because we don't need to be that specific and
			// moving into contact with another moving object that bounces away can get them both stuck; it's easier to bounce off of the
			// "containing" cells, which are static.
			contact := check.ContactWithCell(check.Cells[0])
			dx = contact.X
			b.Speed.X *= -1
		}

		b.Object.Position.X += dx

		if check := b.Object.Check(0, dy); check != nil {
			contact := check.ContactWithCell(check.Cells[0])
			dy = contact.Y
			b.Speed.Y *= -1
		}

		b.Object.Position.Y += dy

		b.Object.Update()

	}

}

func (world *WorldBouncer) Draw(screen *ebiten.Image) {

	for _, o := range world.Geometry {
		ebitenutil.DrawRect(screen, o.Position.X, o.Position.Y, o.Size.X, o.Size.Y, color.RGBA{60, 60, 60, 255})
	}

	for _, b := range world.Bouncers {
		o := b.Object
		ebitenutil.DrawRect(screen, o.Position.X, o.Position.Y, o.Size.X, o.Size.Y, color.RGBA{0, 80, 255, 255})
	}

	if world.Game.Debug {
		world.Game.DebugDraw(screen, world.Space)
	}

	if world.Game.ShowHelpText {

		world.Game.DrawText(screen, 16, 16,
			"~ Bouncer Demo ~",
			"This demo showcases how objects can bounce",
			"off of each other or walls at a good performance.",
			"This is accomplished by each bouncer checking the cell it's",
			"heading into, rather than checking each other bouncer",
			"in play.",
			"",
			"Up Arrow: Add bouncer",
			"Down Arrow: Remove bouncer",
			"",
			fmt.Sprintf("%d Bouncers in the world.", len(world.Bouncers)),
			fmt.Sprintf("%d FPS (frames per second)", int(ebiten.CurrentFPS())),
			fmt.Sprintf("%d TPS (ticks per second)", int(ebiten.CurrentTPS())),
			"",
			"F1: Toggle Debug View",
			"F2: Show / Hide help text",
			"F4: Toggle fullscreen",
			"R: Restart world",
			"E: Next world",
			"Q: Previous world",
		)

	}

}
