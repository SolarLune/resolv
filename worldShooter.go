package main

import (
	"math/rand"

	"github.com/SolarLune/resolv/resolv"
	rl "github.com/gen2brain/raylib-go/raylib"
)

type WorldShooter struct {
	Space      *resolv.Space
	Player     *resolv.Space
	Bullets    []*resolv.Line
	SpawnTimer int
}

func (w *WorldShooter) Create() {
	w.Space = resolv.NewSpace()
	w.Player = resolv.NewSpace()
	w.Player.Add(
		resolv.NewLine(0, 0, 8, 8),
		resolv.NewLine(8, 8, 0, 16),
		resolv.NewLine(0, 16, 0, 0),
	)
	w.Player.AddTags("player")
	w.Player.Move(screenWidth/2, screenHeight/2)
	w.Space.Add(w.Player)
}

func (w *WorldShooter) Update() {

	// Player

	dx, dy := int32(0), int32(0)

	if rl.IsKeyDown(rl.KeyLeft) {
		dx -= 2
	}
	if rl.IsKeyDown(rl.KeyRight) {
		dx += 2
	}

	if rl.IsKeyDown(rl.KeyUp) {
		dy -= 2
	}
	if rl.IsKeyDown(rl.KeyDown) {
		dy += 2
	}

	if rl.IsKeyPressed(rl.KeyX) {
		// Add bullet
		lx, ly := w.Player.GetXY()
		lx += 8
		ly += 8

		line := resolv.NewLine(lx, ly-2, lx+8, ly-2)
		line.AddTags("bullet")
		w.Space.Add(line)
		w.Bullets = append(w.Bullets, line)

		line2 := resolv.NewLine(lx, ly+2, lx+8, ly+2)
		line2.AddTags("bullet")
		w.Space.Add(line2)
		w.Bullets = append(w.Bullets, line2)
	}

	w.SpawnTimer++

	// Spawn rocks
	if w.SpawnTimer >= 4 {
		w.SpawnTimer = 0
		// Spawn a rock
		r := rand.Int31n(8)
		rock := resolv.NewSpace()
		rock.Add(
			resolv.NewLine(0, 0, 4*r, -2*r),
			resolv.NewLine(4*r, -2*r, 6*r, 3*r),
			resolv.NewLine(6*r, 3*r, 2*r, 4*r),
			resolv.NewLine(2*r, 4*r, -2*r, 2*r),
			resolv.NewLine(-2*r, 2*r, 0, 0),
		)
		rock.Move(screenWidth+16, 0)
		// rock := resolv.NewRectangle(screenWidth, 0, 8+rand.Int31n(16), 8+rand.Int31n(16))
		rock.AddTags("rock")
		rock.Move(0, rand.Int31n(screenHeight-16))
		w.Space.Add(rock)
	}

	bullets := w.Space.FilterByTags("bullet")

	toBeRemoved := []resolv.Shape{}

	for _, shape := range *w.Space {

		if shape.HasTags("bullet") { // Move da bulleys
			shape.Move(4, 0)

			// Remove the bullet if it goes offscreen
			x, _ := shape.GetXY()
			if x > screenWidth {
				toBeRemoved = append(toBeRemoved, shape)
			}

		} else if shape.HasTags("rock") { // Move da rox

			shape.Move(-3, 0)

			// Remove the rock if it goes offscreen
			x, _ := shape.GetXY()
			if x < -64 {
				toBeRemoved = append(toBeRemoved, shape)
			}

			// Remove the rock if it's shot
			for _, bullet := range *bullets.GetCollidingShapes(shape) {
				toBeRemoved = append(toBeRemoved, shape, bullet)
			}

		}

	}

	// Remove the shapes out here so we're not looping through shapes in
	// the Space itself to remove them.
	for _, shape := range toBeRemoved {
		w.Space.Remove(shape)
	}

	w.Player.Move(dx, dy)

}

func (w *WorldShooter) drawShape(shape resolv.Shape) {

	drawColor := rl.LightGray

	if shape.HasTags("player") {
		drawColor = rl.Green
	} else if shape.HasTags("bullet") {
		choices := []rl.Color{rl.White, rl.Yellow, rl.Red}
		drawColor = choices[rand.Intn(len(choices))]
	}

	switch b := shape.(type) {

	case *resolv.Rectangle:
		rl.DrawRectangleLines(b.X, b.Y, b.W, b.H, drawColor)
	case *resolv.Line:
		rl.DrawLine(b.X, b.Y, b.X2, b.Y2, drawColor)
	case *resolv.Space:
		for _, s := range *b {
			w.drawShape(s)
		}
	}

}

func (w *WorldShooter) Draw() {

	for _, shape := range *w.Space {
		w.drawShape(shape)
	}

	if drawHelpText {
		DrawText(32, 16,
			"-Shooter test-",
			"You are the green triangle.",
			"Use the arrow keys to move.",
			"Press X to shoot ~LASER BEAMS~.",
		)
	}

}

func (w *WorldShooter) Destroy() {
	w.Space.Clear()
	w.Player = nil
	w.Bullets = []*resolv.Line{}
}
