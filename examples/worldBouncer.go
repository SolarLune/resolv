package main

import (
	"fmt"
	"math/rand/v2"
	"time"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/solarlune/resolv"
)

var (
	TagBouncer = resolv.NewTag("bouncer")
)

type WorldBouncer struct {
	space    *resolv.Space
	Solids   resolv.ShapeCollection
	Bouncers []*Bouncer

	BouncerUpdateTime time.Duration
	UpdateTimeTick    int
}

func NewWorldBouncer() *WorldBouncer {
	w := &WorldBouncer{}
	w.Init()
	return w
}

func (w *WorldBouncer) Init() {

	if w.space != nil {
		w.space.RemoveAll()
	}

	for i := range w.Bouncers {
		w.Bouncers[i] = nil
	}

	w.Bouncers = w.Bouncers[:0]

	// Create the space.
	w.space = resolv.NewSpace(640, 360, 16, 16)

	// Create a selection of shapes that comprise the walls.
	w.Solids = resolv.ShapeCollection{
		resolv.NewRectangleTopLeft(0, 0, 640, 16),
		resolv.NewRectangleTopLeft(0, 360-16, 640, 16),
		resolv.NewRectangleTopLeft(0, 16, 16, 360-16),
		resolv.NewRectangleTopLeft(640-16, 16, 16, 360-16),
		resolv.NewRectangleTopLeft(64, 128, 16, 200),
		resolv.NewRectangleTopLeft(120, 300, 200, 8),
	}

	// Set their tags (not strictly necessary here because the bouncers bounce off of everything and anything)..
	w.Solids.SetTags(TagSolidWall)

	// Add them to the space.
	w.space.Add(w.Solids...)

}

func (w *WorldBouncer) Update() {

	t := time.Now()

	for _, b := range w.Bouncers {
		b.Update()
	}

	if ebiten.IsMouseButtonPressed(ebiten.MouseButtonLeft) {
		x, y := ebiten.CursorPosition()
		b := NewBouncer(float64(x), float64(y), w)
		w.Bouncers = append(w.Bouncers, b)
	}

	if len(w.Bouncers) > 0 {
		if ebiten.IsMouseButtonPressed(ebiten.MouseButtonRight) {
			x, y := ebiten.CursorPosition()
			for i := len(w.Bouncers) - 1; i >= 0; i-- {
				if w.Bouncers[i].Object.Position().Distance(resolv.NewVector(float64(x), float64(y))) < 64 {
					w.space.Remove(w.Bouncers[i].Object)
					w.Bouncers[i] = nil
					w.Bouncers = append(w.Bouncers[:i], w.Bouncers[i+1:]...)
				}
			}
		}
	}

	w.UpdateTimeTick++
	if w.UpdateTimeTick >= 10 {
		w.UpdateTimeTick = 0
		w.BouncerUpdateTime = time.Since(t)
	}

}

func (w *WorldBouncer) Draw(screen *ebiten.Image) {
	CommonDraw(screen, w)
	if GlobalGame.ShowHelpText {
		GlobalGame.DrawText(screen, 0, 128,
			"Bouncer Test",
			"Left click to add spheres",
			"Right click to remove spheres",
			fmt.Sprintf("%d Bouncers in the world", len(w.Bouncers)),
			fmt.Sprintf("Update Time: %s", w.BouncerUpdateTime.String()),
		)
	}
}

// To allow the world's physical state to be drawn using the debug draw function.
func (w *WorldBouncer) Space() *resolv.Space {
	return w.space
}

type Bouncer struct {
	Object   *resolv.Circle
	Movement resolv.Vector

	ColorChange float64
}

func NewBouncer(x, y float64, world *WorldBouncer) *Bouncer {

	bouncer := &Bouncer{
		Object:   resolv.NewCircle(x, y, 8),
		Movement: resolv.NewVector(rand.Float64()*2-1, 0),
	}

	bouncer.Object.Tags().Set(TagBouncer)
	bouncer.Object.SetData(bouncer)

	world.space.Add(bouncer.Object)
	return bouncer

}

func (b *Bouncer) Update() {

	gravity := 0.25
	b.Movement.Y += gravity

	// Clamp movement to the maximum speed of half the size of a ball (so at max speed, it can't go beyond halfway through a surface)
	b.Movement = b.Movement.ClampMagnitude(8)

	b.Object.MoveVec(b.Movement)

	b.ColorChange *= 0.98

	totalMTV := resolv.NewVector(0, 0)

	b.Object.IntersectionTest(resolv.IntersectionTestSettings{

		TestAgainst: b.Object.SelectTouchingCells(1).FilterShapes(),

		OnIntersect: func(set resolv.IntersectionSet) bool {
			b.Movement = b.Movement.Reflect(set.Intersections[0].Normal).Scale(0.9)
			b.ColorChange = b.Movement.Magnitude()
			// Collect all MTV values to apply together (mainly because bouncers might intersect each other and by moving,
			// phase into other Bouncers). By collecting all MTV values together and moving at once, it minimizes the chances
			// of Bouncers moving into each other and forcing lower ones through walls.
			totalMTV = totalMTV.Add(set.MTV)
			return true // Keep looping through intersection sets
		},
	})

	b.Object.MoveVec(totalMTV)

	if b.ColorChange > 1 {
		b.ColorChange = 1
	}

}
