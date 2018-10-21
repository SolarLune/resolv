package main

import (
	"math"

	"github.com/SolarLune/resolv/resolv"
	"github.com/veandco/go-sdl2/sdl"
)

type World6 struct {
	Player            *Bouncer
	FloatingPlatform  *resolv.Line
	FloatingPlatformY float64
}

func (w *World6) Create() {

	space.Clear()

	w.Player = MakeNewBouncer()
	w.Player.Rect.X = 32
	w.Player.Rect.Y = 32
	w.Player.Rect.W = 16
	w.Player.Rect.H = 16
	w.Player.SpeedX = 0
	w.Player.SpeedY = 0

	space.AddShape(w.Player.Rect)

	space.AddShape(resolv.NewRectangle(0, 0, 16, screenHeight))
	space.AddShape(resolv.NewRectangle(screenWidth-16, 0, 16, screenHeight))
	space.AddShape(resolv.NewRectangle(0, 0, screenWidth, 16))
	space.AddShape(resolv.NewRectangle(0, screenHeight-16, screenWidth, 16))

	c := int32(16)

	space.AddShape(resolv.NewRectangle(c*4, screenHeight-c*4, c*3, c))

	for _, shape := range space {
		shape.SetTags("solid")
	}

	// A ramp
	line := resolv.NewLine(c*5, screenHeight-c, c*6, screenHeight-c-8)
	line.SetTags("ramp")
	space.AddShape(line)

	line = resolv.NewLine(c*6, screenHeight-c-8, c*7, screenHeight-c-8)
	line.SetTags("ramp")

	space.AddShape(line)

	rect := resolv.NewRectangle(c*7, screenHeight-c-8, c*2, 8)
	rect.SetTags("solid")
	space.AddShape(rect)

	line = resolv.NewLine(c*9, screenHeight-c-8, c*11, screenHeight-c)
	line.SetTags("ramp")
	space.AddShape(line)

	line = resolv.NewLine(c*13, screenHeight-c*4, c*17, screenHeight-c*6)
	line.SetTags("ramp")
	space.AddShape(line)

	line = resolv.NewLine(c*6, screenHeight-c*7, c*7, screenHeight-c*7)
	line.SetTags("ramp")
	space.AddShape(line)

	w.FloatingPlatform = resolv.NewLine(c*8, screenHeight-c*7, c*9, screenHeight-c*6)
	w.FloatingPlatform.SetTags("ramp")
	space.AddShape(w.FloatingPlatform)
	w.FloatingPlatformY = float64(w.FloatingPlatform.Y)

}

func (w *World6) Update() {

	w.Player.SpeedY += 0.5

	friction := float32(0.5)
	accel := 0.5 + friction

	maxSpd := float32(3)

	w.FloatingPlatformY += math.Sin(float64(sdl.GetTicks()/1000)) * .5

	w.FloatingPlatform.Y = int32(w.FloatingPlatformY)
	w.FloatingPlatform.Y2 = int32(w.FloatingPlatformY) - 16

	if w.Player.SpeedX > friction {
		w.Player.SpeedX -= friction
	} else if w.Player.SpeedX < -friction {
		w.Player.SpeedX += friction
	} else {
		w.Player.SpeedX = 0
	}

	if keyboard.KeyDown(sdl.K_RIGHT) {
		w.Player.SpeedX += accel
	}

	if keyboard.KeyDown(sdl.K_LEFT) {
		w.Player.SpeedX -= accel
	}

	if w.Player.SpeedX > maxSpd {
		w.Player.SpeedX = maxSpd
	}

	if w.Player.SpeedX < -maxSpd {
		w.Player.SpeedX = -maxSpd
	}

	// JUMP

	// Check for a collision downwards by just attempting a resolution downwards and seeing if it collides with something.
	down := space.Resolve(w.Player.Rect, 0, 4)
	onGround := down.Colliding()

	if keyboard.KeyPressed(sdl.K_x) && onGround {
		w.Player.SpeedY = -8
	}

	x := int32(w.Player.SpeedX)
	y := int32(w.Player.SpeedY)

	solids := space.FilterByTags("solid")
	ramps := space.FilterByTags("ramp")

	// X-movement. We only want to collide with solid objects (not ramps) because we want to be able to move up them
	// and don't need to be inhibited on the x-axis when doing so.

	if res := solids.Resolve(w.Player.Rect, x, 0); res.Colliding() {
		x = res.ResolveX
		w.Player.SpeedX = 0
	}

	w.Player.Rect.X += x

	// Y movement. We check for ramp collision first; if we find it, then we just automatically will
	// slide up the ramp because the player is moving into it.

	// We look for ramps a little aggressively downwards because when walking down them, we want to stick to them.
	// If we didn't do this, then you would "bob" when walking down the ramp as the Player moves too quickly out into
	// space for gravity to push back down onto the ramp.
	res := ramps.Resolve(w.Player.Rect, 0, y+4)

	if y < 0 || (res.Teleporting && res.ResolveY < -w.Player.Rect.H/2) {
		res = resolv.Collision{}
	}

	if !res.Colliding() {
		res = solids.Resolve(w.Player.Rect, 0, y)
	}

	if res.Colliding() {
		y = res.ResolveY
		w.Player.SpeedY = 0
	}

	w.Player.Rect.Y += y

}

func (w *World6) Draw() {

	for _, shape := range space {

		rect, ok := shape.(*resolv.Rectangle)

		if ok {

			if rect == w.Player.Rect {
				renderer.SetDrawColor(0, 128, 255, 255)
			} else {
				renderer.SetDrawColor(255, 255, 255, 255)
			}

			renderer.DrawRect(&sdl.Rect{rect.X, rect.Y, rect.W, rect.H})

		}

		line, ok := shape.(*resolv.Line)

		if ok {

			renderer.DrawLine(line.X, line.Y, line.X2, line.Y2)

		}

	}

	if drawHelpText {
		DrawText(0, 0,
			"Platformer test",
			"Use the arrow keys to move",
			"Press X to jump",
			"You can jump through lines or ramps")
	}

}

func (w *World6) Destroy() {
	space.Clear()
	w.Player = nil
}
