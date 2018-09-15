package main

import (
	"math/rand"

	"github.com/SolarLune/resolv/resolv"
	"github.com/veandco/go-sdl2/sdl"
)

type World1 struct{}

func (w World1) Create() {

	var cell int32 = 16
	var screenCellWidth = screenWidth / cell
	var screenCellHeight = screenHeight / cell

	// Just so nobody gets confused, yes, this isn't "true" fidelity because while I'm using floats for the speed variables,
	// I'm putting them into ints in the rectangle rather than having extra X and Y variables (just to make it easier to follow).
	squareSpeedX = (0.5 - rand.Float32()) * 8
	squareSpeedY = (0.5 - rand.Float32()) * 8

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
		shape.SetTag("solid")
	}

	mainSquare = resolv.NewRectangle(40, 64, cell, cell)
	space.AddShape(mainSquare)

}

func (w World1) Update() {

	squareSpeedY += 0.25

	if res := space.Resolve(mainSquare, squareSpeedX, true, "solid"); res.Colliding() {
		mainSquare.X += res.ResolveDistance
		squareSpeedX *= -1
	} else {
		mainSquare.X += int32(squareSpeedX)
	}

	if res := space.Resolve(mainSquare, squareSpeedY, false, "solid"); res.Colliding() {
		mainSquare.Y += res.ResolveDistance
		squareSpeedY *= -1
	} else {
		mainSquare.Y += int32(squareSpeedY)
	}

}

func (w World1) Draw() {

	for _, shape := range space {

		rect, ok := shape.(*resolv.Rectangle)

		renderer.SetDrawColor(255, 255, 255, 255)

		if rect == mainSquare {
			renderer.SetDrawColor(60, 180, 255, 255)
		}

		if ok {
			renderer.DrawRect(&sdl.Rect{X: rect.X, Y: rect.Y, W: rect.W, H: rect.H})
		}

	}

}
