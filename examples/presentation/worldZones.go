package main

import (
	"github.com/SolarLune/resolv"
	rl "github.com/gen2brain/raylib-go/raylib"
)

type WorldZones struct {
	Space  *resolv.Space
	Player *Square
}

func (w *WorldZones) Create() {

	var cell int32 = 16

	w.Space = resolv.NewSpace()

	w.Space.Add(resolv.NewRectangle(0, 0, screenWidth, cell))
	w.Space.Add(resolv.NewRectangle(0, cell, cell, screenHeight-cell))
	w.Space.Add(resolv.NewRectangle(screenWidth-cell, cell, cell, screenHeight-cell))
	w.Space.Add(resolv.NewRectangle(cell, screenHeight-cell, screenWidth-(cell*2), cell))

	w.Space.AddTags("solid")

	zone := resolv.NewRectangle(screenWidth/2, cell*2, screenWidth/2-(cell*2), screenHeight/2)
	zone.AddTags("zone")
	w.Space.Add(zone)

	zone = resolv.NewRectangle(cell*4, cell*7, cell*2, cell*4)
	zone.AddTags("zone")
	w.Space.Add(zone)

	circle := resolv.NewCircle(cell*7, cell*7, 32)
	circle.AddTags("zone")
	w.Space.Add(circle)

	w.Player = NewSquare(w.Space)
	w.Player.Rect.X = screenWidth / 2
	w.Player.Rect.Y = screenHeight / 2
	w.Space.Add(w.Player.Rect)

}

func (w *WorldZones) Update() {

	var friction float32 = 0.5
	accel := 0.5 + friction
	var maxSpd float32 = 3

	// Note that I'm being lazy and using the squares list / Bouncer struct from World1.go to store player data,
	// rather than making a new set of data here for World2.

	player := w.Player

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

	if rl.IsKeyDown(rl.KeyRight) {
		player.SpeedX += accel
	}
	if rl.IsKeyDown(rl.KeyLeft) {
		player.SpeedX -= accel
	}

	if rl.IsKeyDown(rl.KeyUp) {
		player.SpeedY -= accel
	}
	if rl.IsKeyDown(rl.KeyDown) {
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

	solids := w.Space.FilterByTags("solid")

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

func (w *WorldZones) Draw() {

	player := w.Player

	touching := "You aren't touching a zone."

	for _, other := range *w.Space {

		drawColor := rl.LightGray

		if other.HasTags("zone") {

			drawColor = rl.Yellow
			if other.IsColliding(player.Rect) {
				drawColor = rl.Blue
				touching = "You ARE touching a zone."
			}

		}

		if other == player.Rect {
			drawColor = rl.Green
		}

		switch shape := other.(type) {

		case *resolv.Rectangle:

			rl.DrawRectangleLines(shape.X, shape.Y, shape.W, shape.H, drawColor)

		case *resolv.Circle:

			rl.DrawCircleLines(shape.X, shape.Y, float32(shape.Radius), drawColor)

		}

	}

	if drawHelpText {
		DrawText(32, 16,
			"-Zone collision test-",
			"Use the arrow keys to move the green square.",
			"When touching a zone (which is just a shape),",
			"it will turn blue.",
			touching,
		)
	}

}

func (w *WorldZones) Destroy() {
	w.Player = nil
	w.Space.Clear()
}
