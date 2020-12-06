package main

import (
	"math"

	"github.com/SolarLune/resolv"
	rl "github.com/gen2brain/raylib-go/raylib"
	"github.com/veandco/go-sdl2/sdl"
)

type WorldPlatformer struct {
	Player            *Square
	Space             *resolv.Space
	FloatingPlatform  *resolv.Line
	FloatingPlatformY float64
}

func (w *WorldPlatformer) Create() {

	w.Space = resolv.NewSpace()
	w.Space.Clear()

	w.Player = NewSquare(w.Space)
	w.Player.Rect.X = 64
	w.Player.Rect.Y = 32
	w.Player.Rect.W = 16
	w.Player.Rect.H = 16
	w.Player.SpeedX = 0
	w.Player.SpeedY = 0

	w.Space.Add(w.Player.Rect)

	w.Space.Add(resolv.NewRectangle(0, 0, 16, screenHeight))
	w.Space.Add(resolv.NewRectangle(screenWidth-16, 0, 16, screenHeight))
	w.Space.Add(resolv.NewRectangle(0, 0, screenWidth, 16))
	w.Space.Add(resolv.NewRectangle(0, screenHeight-16, screenWidth, 16))

	c := int32(16)

	w.Space.Add(resolv.NewRectangle(c*4, screenHeight-c*4, c*3, c))

	w.Space.AddTags("solid")

	// A ramp
	line := resolv.NewLine(c*5, screenHeight-c, c*6, screenHeight-c-8)
	line.AddTags("ramp")
	w.Space.Add(line)

	line = resolv.NewLine(c*6, screenHeight-c-8, c*7, screenHeight-c-8)
	line.AddTags("ramp")

	w.Space.Add(line)

	rect := resolv.NewRectangle(c*7, screenHeight-c-8, c*2, 8)
	rect.AddTags("solid")
	w.Space.Add(rect)

	line = resolv.NewLine(c*9, screenHeight-c-8, c*11, screenHeight-c)
	line.AddTags("ramp")
	w.Space.Add(line)

	line = resolv.NewLine(c*13, screenHeight-c*4, c*17, screenHeight-c*6)
	line.AddTags("ramp")
	w.Space.Add(line)

	line = resolv.NewLine(c*6, screenHeight-c*7, c*7, screenHeight-c*7)
	line.AddTags("ramp")
	w.Space.Add(line)

	w.FloatingPlatform = resolv.NewLine(c*8, screenHeight-c*7, c*9, screenHeight-c*6)
	w.FloatingPlatform.AddTags("ramp")
	w.Space.Add(w.FloatingPlatform)
	w.FloatingPlatformY = float64(w.FloatingPlatform.Y)

}

func (w *WorldPlatformer) Update() {

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

	if rl.IsKeyDown(rl.KeyRight) {
		w.Player.SpeedX += accel
	}

	if rl.IsKeyDown(rl.KeyLeft) {
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
	down := w.Space.Resolve(w.Player.Rect, 0, 4)
	onGround := down.Colliding()

	if rl.IsKeyPressed(rl.KeyX) && onGround {
		w.Player.SpeedY = -8
	}

	x := int32(w.Player.SpeedX)
	y := int32(w.Player.SpeedY)

	solids := w.Space.FilterByTags("solid")
	ramps := w.Space.FilterByTags("ramp")

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

func (w *WorldPlatformer) Draw() {

	for _, shape := range *w.Space {

		rect, ok := shape.(*resolv.Rectangle)

		drawColor := rl.LightGray

		if ok {

			if rect == w.Player.Rect {
				drawColor = rl.Green
				rl.DrawLine(rect.X+5, rect.Y+3, rect.X+5, rect.Y+8, drawColor)
				rl.DrawLine(rect.X+8, rect.Y+3, rect.X+8, rect.Y+8, drawColor)
			}

			rl.DrawRectangleLines(rect.X, rect.Y, rect.W, rect.H, drawColor)

		}

		line, ok := shape.(*resolv.Line)

		if ok {

			rl.DrawLine(line.X, line.Y, line.X2, line.Y2, rl.Blue)

		}

	}

	if drawHelpText {
		DrawText(32, 16,
			"-Platformer test-",
			"You are the green square.",
			"Use the arrow keys to move.",
			"Press X to jump.",
			"You can jump through blue ramps / platforms.")
	}

}

func (w *WorldPlatformer) Destroy() {
	w.Space.Clear()
	w.Player = nil
}
