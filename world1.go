package main

import (
	"fmt"
	"math/rand"
	"strconv"

	"github.com/SolarLune/resolv/resolv"
	"github.com/veandco/go-sdl2/sdl"
)

type WorldInterface interface {
	Create()
	Update()
	Draw()
	Destroy()
}

type World1 struct {
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

	bouncer.Rect.AddTags("bouncer", "solid")

	bouncer.Rect.SetData(bouncer)

	squares = append(squares, bouncer)

	space.AddShape(bouncer.Rect)

	return bouncer

}

func (w *World1) Create() {

	squares = make([]*Bouncer, 0)

	var screenCellWidth = screenWidth / cell
	var screenCellHeight = screenHeight / cell

	// Just so nobody gets confused, yes, this isn't "true" fidelity because while I'm using floats for the speed variables,
	// I'm putting them into ints in the rectangle rather than having extra X and Y variables (just to make it easier to follow).

	space.Clear()
	space.AddShape(resolv.NewRectangle(0, 0, screenWidth, cell))
	space.AddShape(resolv.NewRectangle(0, cell, cell, screenHeight-cell))
	space.AddShape(resolv.NewRectangle(screenWidth-cell, cell, cell, screenHeight-cell))
	space.AddShape(resolv.NewRectangle(cell, screenHeight-cell, screenWidth-(cell*2), cell))

	for i := 0; i < 20; i++ {
		x := rand.Int31n(screenCellWidth - 2)
		y := rand.Int31n(screenCellHeight - 2)
		space.AddShape(resolv.NewRectangle(cell+(x*cell), cell+(y*cell), cell*(1+rand.Int31n(16)), cell*(1+rand.Int31n(16))))
	}

	for _, shape := range *space {
		shape.AddTags("solid")
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
			// This makes the bouncers able to rebound higher if they get a boost from another bouncer below~
			if bouncer.SpeedY < 0 && bouncer.SpeedY > -5 {
				bouncer.SpeedY = -5
			}
			bouncer.BounceFrame = 1
		} else {
			bouncer.Rect.Y += int32(bouncer.SpeedY)
		}

	}

	if keyboard.KeyDown(sdl.K_UP) {
		MakeNewBouncer()
		fmt.Println(len(squares), " bouncers in the world now.")
	}

	if keyboard.KeyPressed(sdl.K_s) { // The ability to trigger solidity
		if !squares[0].Rect.HasTags("solid") {
			space.FilterByTags("bouncer").AddTags("solid")
		} else {
			space.FilterByTags("bouncer").RemoveTags("solid")
		}
	}

	if keyboard.KeyDown(sdl.K_DOWN) {

		bouncers := *space.FilterByTags("bouncer")

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

	for _, shape := range *space {

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

		if b.Rect.HasTags("solid") {
			renderer.SetDrawColor(60, g, 255, 255)
		} else {
			renderer.SetDrawColor(g, g, g, 255)
		}

		renderer.DrawRect(&sdl.Rect{X: b.Rect.X, Y: b.Rect.Y, W: b.Rect.W, H: b.Rect.H})

	}

	if drawHelpText {
		DrawText(32, 16,
			"Bouncer stress test",
			"Press Up to spawn bouncers",
			"Press Down to remove bouncers",
			"Press 'S' to toggle solidity",
			"Press 'R' to restart with a new",
			"layout",
			"Use the number keys to jump to",
			"different worlds",
			strconv.Itoa(len(squares))+" bouncers in the world",
			"Press F1 to turn on or off this text")
	}
}

func (w World1) Destroy() {
	squares = make([]*Bouncer, 0)
	space.Clear()
}
