package main

import (
	"github.com/SolarLune/resolv/resolv"
	"github.com/veandco/go-sdl2/sdl"
	"github.com/veandco/go-sdl2/ttf"
)

type World2 struct{}

func (w World2) Create() {

	var cell int32 = 16

	space = resolv.NewSpace()
	space.AddShape(resolv.NewRectangle(0, 0, screenWidth, cell))
	space.AddShape(resolv.NewRectangle(0, cell, cell, screenHeight-cell))
	space.AddShape(resolv.NewRectangle(screenWidth-cell, cell, cell, screenHeight-cell))
	space.AddShape(resolv.NewRectangle(cell, screenHeight-cell, screenWidth-(cell*2), cell))

	for _, shape := range space {
		shape.SetTag("solid")
	}

	zone := resolv.NewRectangle(screenWidth/2, cell*2, screenWidth/2-(cell*2), screenHeight/2)
	zone.SetTag("zone")
	space.AddShape(zone)

	mainSquare = resolv.NewRectangle(screenWidth/2, screenHeight/2, cell, cell)
	space.AddShape(mainSquare)

}

func (w World2) Update() {

	var friction float32 = 0.5
	accel := 0.5 + friction
	var maxSpd float32 = 3

	if squareSpeedX >= friction {
		squareSpeedX -= friction
	} else if squareSpeedX <= -friction {
		squareSpeedX += friction
	} else {
		squareSpeedX = 0
	}

	if squareSpeedY >= friction {
		squareSpeedY -= friction
	} else if squareSpeedY <= -friction {
		squareSpeedY += friction
	} else {
		squareSpeedY = 0
	}

	if keyboard.KeyDown(sdl.K_RIGHT) {
		squareSpeedX += accel
	}
	if keyboard.KeyDown(sdl.K_LEFT) {
		squareSpeedX -= accel
	}

	if keyboard.KeyDown(sdl.K_UP) {
		squareSpeedY -= accel
	}
	if keyboard.KeyDown(sdl.K_DOWN) {
		squareSpeedY += accel
	}

	if squareSpeedX > maxSpd {
		squareSpeedX = maxSpd
	} else if squareSpeedX < -maxSpd {
		squareSpeedX = -maxSpd
	}

	if squareSpeedY > maxSpd {
		squareSpeedY = maxSpd
	} else if squareSpeedY < -maxSpd {
		squareSpeedY = -maxSpd
	}

	if res := space.Resolve(mainSquare, squareSpeedX, true, "solid"); res.Colliding() {
		mainSquare.X += res.ResolveDistance
		squareSpeedX = 0
	} else {
		mainSquare.X += int32(squareSpeedX)
	}

	if res := space.Resolve(mainSquare, squareSpeedY, false, "solid"); res.Colliding() {
		mainSquare.Y += res.ResolveDistance
		squareSpeedY = 0
	} else {
		mainSquare.Y += int32(squareSpeedY)
	}

}

func (w World2) Draw() {

	for _, shape := range space {

		renderer.SetDrawColor(255, 255, 255, 255)

		rect, ok := shape.(*resolv.Rectangle)

		if ok {

			if rect.GetTag() == "zone" {

				font, _ := ttf.OpenFont("ARCADEPI.TTF", 12)
				defer font.Close()

				surf, _ := font.RenderUTF8Solid("Isn't touching a zone.", sdl.Color{R: 255, G: 255, B: 255, A: 255})

				renderer.SetDrawColor(255, 255, 0, 255)
				if rect.IsColliding(mainSquare) {
					surf, _ = font.RenderUTF8Solid("Touching a zone!", sdl.Color{R: 255, G: 255, B: 255, A: 255})
					renderer.SetDrawColor(255, 0, 0, 255)
				}

				text, _ := renderer.CreateTextureFromSurface(surf)
				defer text.Destroy()

				_, _, w, h, _ := text.Query()

				renderer.Copy(text, &sdl.Rect{X: 0, Y: 0, W: w, H: h}, &sdl.Rect{X: 0, Y: 0, W: w, H: h})

			}

			if rect == mainSquare {
				renderer.SetDrawColor(0, 255, 0, 255)
			}

			renderer.DrawRect(&sdl.Rect{X: rect.X, Y: rect.Y, W: rect.W, H: rect.H})

		}

	}

}
