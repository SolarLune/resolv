package main

import (
	"math"

	"github.com/SolarLune/resolv/resolv"
	"github.com/veandco/go-sdl2/sdl"
)

var mainCircle *resolv.Circle

var secondCircleZone *resolv.Circle

type World3 struct{}

func (w World3) Create() {

	var cell int32 = 16

	space.Clear()
	space.AddShape(resolv.NewRectangle(0, 0, screenWidth, cell))
	space.AddShape(resolv.NewRectangle(0, cell, cell, screenHeight-cell))
	space.AddShape(resolv.NewRectangle(screenWidth-cell, cell, cell, screenHeight-cell))
	space.AddShape(resolv.NewRectangle(cell, screenHeight-cell, screenWidth-(cell*2), cell))
	space.AddShape(resolv.NewCircle(30, 60, cell))

	for _, shape := range *space {
		shape.AddTags("solid")
	}

	zone := resolv.NewRectangle(screenWidth/2, cell*2, screenWidth/2-(cell*2), screenHeight/2)
	zone.AddTags("zone")
	space.AddShape(zone)

	secondCircleZone = resolv.NewCircle(70, 70, 8)
	secondCircleZone.AddTags("zone")
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

	solids := space.FilterByTags("solid")

	res := solids.Resolve(mainCircle, dx, 0)

	if res.Colliding() {
		mainCircle.X += res.ResolveX
	} else {
		mainCircle.X += dx
	}

	if keyboard.KeyDown(sdl.K_UP) {
		dy -= 2
	}
	if keyboard.KeyDown(sdl.K_DOWN) {
		dy += 2
	}

	res = solids.Resolve(mainCircle, 0, dy)

	if res.Colliding() {
		mainCircle.Y += res.ResolveY
	} else {
		mainCircle.Y += dy
	}

}

func (world World3) Draw() {

	renderer.Clear()

	touching := "Not touching a zone"

	for _, shape := range *space {

		renderer.SetDrawColor(255, 255, 255, 255)

		rect, ok := shape.(*resolv.Rectangle)

		if ok {

			if rect.HasTags("zone") {

				renderer.SetDrawColor(255, 255, 0, 255)
				if rect.IsColliding(mainCircle) {
					touching = "Touching a zone"
					renderer.SetDrawColor(255, 0, 0, 255)
				}

			}

			renderer.DrawRect(&sdl.Rect{X: rect.X, Y: rect.Y, W: rect.W, H: rect.H})
		}

		renderer.SetDrawColor(255, 255, 255, 255)

		circle, ok := shape.(*resolv.Circle)

		if ok {

			if circle.HasTags("zone") {

				renderer.SetDrawColor(255, 255, 0, 255)
				if circle.IsColliding(mainCircle) {
					touching = "Touching a zone"
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

	if drawHelpText {
		DrawText(0, 0, "Circle collision testing",
			touching)
	}

}

func (w World3) Destroy() {
	squares = make([]*Bouncer, 0)

	secondCircleZone = nil

	mainCircle = nil

	space.Clear()
}
