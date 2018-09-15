package main

import (
	"math"

	"github.com/SolarLune/resolv/resolv"
	"github.com/veandco/go-sdl2/sdl"
	"github.com/veandco/go-sdl2/ttf"
)

var mainCircle *resolv.Circle

var secondCircleZone *resolv.Circle

type World3 struct{}

func (w World3) Create() {

	var cell int32 = 16

	space = resolv.NewSpace()
	space.AddShape(resolv.NewRectangle(0, 0, screenWidth, cell))
	space.AddShape(resolv.NewRectangle(0, cell, cell, screenHeight-cell))
	space.AddShape(resolv.NewRectangle(screenWidth-cell, cell, cell, screenHeight-cell))
	space.AddShape(resolv.NewRectangle(cell, screenHeight-cell, screenWidth-(cell*2), cell))
	space.AddShape(resolv.NewCircle(30, 60, cell))

	for _, shape := range space {
		shape.SetTag("solid")
	}

	zone := resolv.NewRectangle(screenWidth/2, cell*2, screenWidth/2-(cell*2), screenHeight/2)
	zone.SetTag("zone")
	space.AddShape(zone)

	secondCircleZone = resolv.NewCircle(70, 70, 8)
	secondCircleZone.SetTag("zone")
	space.AddShape(secondCircleZone)

	mainCircle = resolv.NewCircle(screenWidth/2, screenHeight/2, 32)
	space.AddShape(mainCircle)

}

func (w World3) Update() {

	var dx int32 = 0
	var dy int32 = 0

	if keyboard.KeyDown(sdl.K_LEFT) {
		dx -= 2
	}
	if keyboard.KeyDown(sdl.K_RIGHT) {
		dx += 2
	}

	res := space.Resolve(mainCircle, float32(dx), true, "solid")

	if res.Colliding() {
		dx = res.ResolveDistance
	}

	mainCircle.X += dx

	if keyboard.KeyDown(sdl.K_UP) {
		dy -= 2
	}
	if keyboard.KeyDown(sdl.K_DOWN) {
		dy += 2
	}

	res = space.Resolve(mainCircle, float32(dy), false, "solid")

	if res.Colliding() {
		dy = res.ResolveDistance
	}

	mainCircle.Y += dy

}

func (world World3) Draw() {

	renderer.Clear()

	touching := false

	for _, shape := range space {

		renderer.SetDrawColor(255, 255, 255, 255)

		rect, ok := shape.(*resolv.Rectangle)

		if ok {

			if rect.GetTag() == "zone" {

				renderer.SetDrawColor(255, 255, 0, 255)
				if rect.IsColliding(mainCircle) {
					touching = true
					renderer.SetDrawColor(255, 0, 0, 255)
				}

			}

			renderer.DrawRect(&sdl.Rect{X: rect.X, Y: rect.Y, W: rect.W, H: rect.H})
		}

		renderer.SetDrawColor(255, 255, 255, 255)

		circle, ok := shape.(*resolv.Circle)

		if ok {

			if circle.GetTag() == "zone" {

				renderer.SetDrawColor(255, 255, 0, 255)
				if circle.IsColliding(mainCircle) {
					touching = true
					renderer.SetDrawColor(255, 0, 0, 255)
				}

			}

			if circle == mainCircle {
				renderer.SetDrawColor(0, 255, 0, 255)
			}

			lineNum := 16

			pi2 := math.Pi * 2
			segRad := pi2 / float64(lineNum)

			for i := 0; i < lineNum; i++ {

				startX := circle.X + int32(math.Cos(segRad*float64(i+1))*float64(circle.Radius))
				startY := circle.Y + int32(math.Sin(segRad*float64(i+1))*float64(circle.Radius))

				endX := circle.X + int32(math.Cos(segRad*float64(i+2))*float64(circle.Radius))
				endY := circle.Y + int32(math.Sin(segRad*float64(i+2))*float64(circle.Radius))

				// For some reason, this doesn't scale correctly visually with SDL2...?
				renderer.DrawLine(startX, startY, endX, endY)

			}

		}

	}

	font, _ := ttf.OpenFont("ARCADEPI.TTF", 12)
	defer font.Close()

	var surf *sdl.Surface

	if touching {
		surf, _ = font.RenderUTF8Solid("Touching a zone!", sdl.Color{R: 255, G: 255, B: 255, A: 255})
	} else {
		surf, _ = font.RenderUTF8Solid("Isn't touching a zone.", sdl.Color{R: 255, G: 255, B: 255, A: 255})
	}

	text, _ := renderer.CreateTextureFromSurface(surf)
	defer text.Destroy()

	_, _, w, h, _ := text.Query()

	renderer.Copy(text, &sdl.Rect{X: 0, Y: 0, W: w, H: h}, &sdl.Rect{X: 0, Y: 0, W: w, H: h})

}
