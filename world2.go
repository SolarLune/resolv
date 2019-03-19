package main

import (
	"github.com/SolarLune/resolv/resolv"
	"github.com/veandco/go-sdl2/sdl"
)

type World2 struct{}

func (w World2) Create() {

	var cell int32 = 16

	space.Clear()
	space.AddShape(resolv.NewRectangle(0, 0, screenWidth, cell))
	space.AddShape(resolv.NewRectangle(0, cell, cell, screenHeight-cell))
	space.AddShape(resolv.NewRectangle(screenWidth-cell, cell, cell, screenHeight-cell))
	space.AddShape(resolv.NewRectangle(cell, screenHeight-cell, screenWidth-(cell*2), cell))

	for _, shape := range *space {
		shape.AddTags("solid")
	}

	zone := resolv.NewRectangle(screenWidth/2, cell*2, screenWidth/2-(cell*2), screenHeight/2)
	zone.AddTags("zone")
	space.AddShape(zone)

	zone = resolv.NewRectangle(cell*4, cell*7, cell*2, cell*4)
	zone.AddTags("zone")
	space.AddShape(zone)

	squares = make([]*Bouncer, 0)
	bouncer := MakeNewBouncer()
	squares = append(squares, bouncer)

}

func (w World2) Update() {

	var friction float32 = 0.5
	accel := 0.5 + friction
	var maxSpd float32 = 3

	// Note that I'm being lazy and using the squares list / Bouncer struct from World1.go to store player data,
	// rather than making a new set of data here for World2.
	player := squares[0]

	if player.SpeedX >= friction {
		player.SpeedX -= friction
	} else if player.SpeedX <= -friction {
		player.SpeedX += friction
	} else {
		player.SpeedX = 0
	}

	if player.SpeedY >= friction {
		player.SpeedY -= friction
	} else if player.SpeedY <= -friction {
		player.SpeedY += friction
	} else {
		player.SpeedY = 0
	}

	if keyboard.KeyDown(sdl.K_RIGHT) {
		player.SpeedX += accel
	}
	if keyboard.KeyDown(sdl.K_LEFT) {
		player.SpeedX -= accel
	}

	if keyboard.KeyDown(sdl.K_UP) {
		player.SpeedY -= accel
	}
	if keyboard.KeyDown(sdl.K_DOWN) {
		player.SpeedY += accel
	}

	if player.SpeedX > maxSpd {
		player.SpeedX = maxSpd
	} else if player.SpeedX < -maxSpd {
		player.SpeedX = -maxSpd
	}

	if player.SpeedY > maxSpd {
		player.SpeedY = maxSpd
	} else if player.SpeedY < -maxSpd {
		player.SpeedY = -maxSpd
	}

	solids := space.FilterByTags("solid")

	if res := solids.Resolve(player.Rect, int32(player.SpeedX), 0); res.Colliding() {
		player.Rect.X += res.ResolveX
		player.SpeedX = 0
	} else {
		player.Rect.X += int32(player.SpeedX)
	}

	if res := solids.Resolve(player.Rect, 0, int32(player.SpeedY)); res.Colliding() {
		player.Rect.Y += res.ResolveY
		player.SpeedY = 0
	} else {
		player.Rect.Y += int32(player.SpeedY)
	}

}

func (w World2) Draw() {

	player := squares[0]

	touching := "You aren't touching a zone"

	for _, shape := range *space {

		renderer.SetDrawColor(255, 255, 255, 255)

		rect, ok := shape.(*resolv.Rectangle)

		if ok {

			if rect.HasTags("zone") {

				renderer.SetDrawColor(255, 255, 0, 255)
				if rect.IsColliding(player.Rect) {
					renderer.SetDrawColor(255, 0, 0, 255)
					touching = "You ARE touching a zone"
				}

			}

			if rect == player.Rect {
				renderer.SetDrawColor(0, 255, 0, 255)
			}

			renderer.DrawRect(&sdl.Rect{X: rect.X, Y: rect.Y, W: rect.W, H: rect.H})

		}

	}

	if drawHelpText {
		DrawText(0, 0,
			"Zone collision test",
			"Use the arrow keys to move",
			touching)
	}

}

func (w World2) Destroy() {
	squares = make([]*Bouncer, 0)
	space.Clear()
}
