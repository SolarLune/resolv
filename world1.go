package main

import (
	"fmt"
	"math/rand"
	"strconv"

	"github.com/SolarLune/resolv/resolv"
	"github.com/veandco/go-sdl2/sdl"
)

type World1 struct{}

type Bouncer struct {
	Rect   *resolv.Rectangle
	SpeedX float32
	SpeedY float32
}

var squares []*Bouncer

func MakeNewBouncer() *Bouncer {

	var cell int32 = 16

	bouncer := &Bouncer{Rect: resolv.NewRectangle(cell*2+rand.Int31n(screenWidth-cell*4), cell*2+rand.Int31n(screenHeight-cell*4), cell, cell),
		SpeedX: (0.5 - rand.Float32()) * 8,
		SpeedY: (0.5 - rand.Float32()) * 8}

	// Attempt to not spawn a Bouncer in an occupied location
	for i := 0; i < 100; i++ {

		if space.IsColliding(bouncer.Rect) {

			bouncer.Rect.X = cell*2 + rand.Int31n(screenWidth-cell*4)
			bouncer.Rect.Y = cell*2 + rand.Int31n(screenHeight-cell*4)

		}

	}

	bouncer.Rect.SetTags("bouncer", "solid")

	squares = append(squares, bouncer)

	space.AddShape(bouncer.Rect)

	return bouncer

}

func (w World1) Create() {

	squares = make([]*Bouncer, 0)

	var cell int32 = 16
	var screenCellWidth = screenWidth / cell
	var screenCellHeight = screenHeight / cell

	// Just so nobody gets confused, yes, this isn't "true" fidelity because while I'm using floats for the speed variables,
	// I'm putting them into ints in the rectangle rather than having extra X and Y variables (just to make it easier to follow).

	space = resolv.NewSpace()
	space.AddShape(resolv.NewRectangle(0, 0, screenWidth, cell))
	space.AddShape(resolv.NewRectangle(0, cell, cell, screenHeight-cell))
	space.AddShape(resolv.NewRectangle(screenWidth-cell, cell, cell, screenHeight-cell))
	space.AddShape(resolv.NewRectangle(cell, screenHeight-cell, screenWidth-(cell*2), cell))

	for i := 0; i < 5; i++ {
		x := rand.Int31n(screenCellWidth - 2)
		y := rand.Int31n(screenCellHeight - 2)
		space.AddShape(resolv.NewRectangle(cell+(x*cell), cell+(y*cell), cell*(1+rand.Int31n(8)), cell))
	}

	for _, shape := range space {
		shape.SetTags("solid")
	}

	MakeNewBouncer()

}

func (w World1) Update() {

	solids := space.FilterByTags("solid")

	for _, bouncer := range squares {

		bouncer.SpeedY += 0.25

		// The additional teleporting check means that it won't resolve in a way that would cause it to move inordinately far (i.e.
		// teleporting). See the docs in resolv.go to see exactly what Teleporting is defined as.
		if res := solids.Resolve(bouncer.Rect, bouncer.SpeedX, 0); res.Colliding() && !res.Teleporting {
			bouncer.Rect.X += res.ResolveX
			bouncer.SpeedX *= -1
		} else {
			bouncer.Rect.X += int32(bouncer.SpeedX)
		}

		if res := solids.Resolve(bouncer.Rect, 0, bouncer.SpeedY); res.Colliding() && !res.Teleporting {
			bouncer.Rect.Y += res.ResolveY
			bouncer.SpeedY *= -1
		} else {
			bouncer.Rect.Y += int32(bouncer.SpeedY)
		}

	}

	if keyboard.KeyDown(sdl.K_UP) {
		MakeNewBouncer()
		fmt.Println(len(squares), " bouncers in the world now.")
	}

	if keyboard.KeyDown(sdl.K_DOWN) {

		bouncers := space.FilterByTags("bouncer")

		if len(bouncers) > 0 {

			space.RemoveShape(bouncers[0])

			for i, b := range squares {

				if b.Rect == bouncers[0] {
					squares[i] = nil
					squares = append(squares[:i], squares[i+1:]...)
				}

			}

			fmt.Println(len(squares), " bouncers in the world now.")

		}

	}

}

func (w World1) Draw() {

	for _, shape := range space {

		rect, ok := shape.(*resolv.Rectangle)

		renderer.SetDrawColor(255, 255, 255, 255)

		if rect.HasTags("bouncer") {
			renderer.SetDrawColor(60, 180, 255, 255)
		}

		if ok {
			renderer.DrawRect(&sdl.Rect{X: rect.X, Y: rect.Y, W: rect.W, H: rect.H})
		}

	}

	DrawText("Press Up to spawn bouncers", 32, 16)
	DrawText("Press Down to remove bouncers", 32, 32)
	DrawText("Press 'R' to restart with a new", 32, 48)
	DrawText("layout", 32, 64)
	DrawText(strconv.Itoa(len(squares))+" bouncers in the world", 32, 80)

}
