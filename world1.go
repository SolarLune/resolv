package main

import (
	"fmt"
	"math/rand"
	"strconv"

	"github.com/SolarLune/resolv/resolv"
	"github.com/veandco/go-sdl2/sdl"
)

type World1 struct {
	DrawInfo bool
}

type Bouncer struct {
	Rect        *resolv.Rectangle
	SpeedX      float32
	SpeedY      float32
	BounceFrame float32
}

var squares []*Bouncer

func MakeNewBouncer() *Bouncer {

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

func (w *World1) Create() {

	w.DrawInfo = true

	squares = make([]*Bouncer, 0)

	var screenCellWidth = screenWidth / cell
	var screenCellHeight = screenHeight / cell

	// Just so nobody gets confused, yes, this isn't "true" fidelity because while I'm using floats for the speed variables,
	// I'm putting them into ints in the rectangle rather than having extra X and Y variables (just to make it easier to follow).

	space = resolv.NewSpace()
	space.AddShape(resolv.NewRectangle(0, 0, screenWidth, cell))
	space.AddShape(resolv.NewRectangle(0, cell, cell, screenHeight-cell))
	space.AddShape(resolv.NewRectangle(screenWidth-cell, cell, cell, screenHeight-cell))
	space.AddShape(resolv.NewRectangle(cell, screenHeight-cell, screenWidth-(cell*2), cell))

	for i := 0; i < 10; i++ {
		x := rand.Int31n(screenCellWidth - 2)
		y := rand.Int31n(screenCellHeight - 2)
		space.AddShape(resolv.NewRectangle(cell+(x*cell), cell+(y*cell), cell*(1+rand.Int31n(16)), cell*(1+rand.Int31n(16))))
	}

	for _, shape := range space {
		shape.SetTags("solid")
	}

	MakeNewBouncer()

}

func (w *World1) Update() {

	solids := space.FilterByTags("solid")

	for _, bouncer := range squares {

		bouncer.SpeedY += 0.25
		bouncer.BounceFrame *= .9

		if bouncer.SpeedY > float32(cell) {
			bouncer.SpeedY = float32(cell)
		} else if bouncer.SpeedY < -float32(cell) {
			bouncer.SpeedY = -float32(cell)
		}

		if bouncer.SpeedX > float32(cell) {
			bouncer.SpeedX = float32(cell)
		} else if bouncer.SpeedX < -float32(cell) {
			bouncer.SpeedX = -float32(cell)
		}

		// The additional teleporting check means that it won't resolve in a way that would cause it to move inordinately far (i.e.
		// teleporting). See the docs in resolv.go to see exactly what Teleporting is defined as.
		if res := solids.Resolve(bouncer.Rect, int32(bouncer.SpeedX), 0); res.Colliding() && !res.Teleporting {
			bouncer.Rect.X += res.ResolveX
			bouncer.SpeedX *= -1
			bouncer.BounceFrame = 1
		} else {
			bouncer.Rect.X += int32(bouncer.SpeedX)
		}

		if res := solids.Resolve(bouncer.Rect, 0, int32(bouncer.SpeedY)); res.Colliding() && !res.Teleporting {
			bouncer.Rect.Y += res.ResolveY
			bouncer.SpeedY *= -1
			bouncer.BounceFrame = 1
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

	if keyboard.KeyPressed(sdl.K_F1) {
		w.DrawInfo = !w.DrawInfo
	}

}

func (w World1) Draw() {

	for _, shape := range space {

		// Living on the edge~~~
		// We know that this Space just has Rectangles, so we'll just assume they are

		rect := shape.(*resolv.Rectangle)

		if !rect.HasTags("bouncer") {

			renderer.SetDrawColor(255, 255, 255, 255)

			renderer.DrawRect(&sdl.Rect{X: rect.X, Y: rect.Y, W: rect.W, H: rect.H})

		}

	}

	for _, b := range squares {

		g := uint8(60) + uint8((255-60)*b.BounceFrame)

		renderer.SetDrawColor(60, g, 255, 255)

		renderer.DrawRect(&sdl.Rect{X: b.Rect.X, Y: b.Rect.Y, W: b.Rect.W, H: b.Rect.H})

	}

	if w.DrawInfo {
		DrawText("Press Up to spawn bouncers", 32, 16)
		DrawText("Press Down to remove bouncers", 32, 32)
		DrawText("Press 'R' to restart with a new", 32, 48)
		DrawText("layout", 32, 64)
		DrawText(strconv.Itoa(len(squares))+" bouncers in the world", 32, 80)
		DrawText("Press F1 to turn on or off this text", 32, 96)
	}
}
