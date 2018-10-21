package main

import (
	"math"

	"github.com/SolarLune/resolv/resolv"
	"github.com/veandco/go-sdl2/sdl"
)

type MultiBodyPlayer struct {
	SpeedX, SpeedY float32
	Body           *resolv.Space
}

type World7 struct {
	Player *MultiBodyPlayer
}

func (w *World7) Create() {

	space.Clear()

	w.Player = &MultiBodyPlayer{}
	w.Player.Body = resolv.NewSpace()

	rect := resolv.NewRectangle(32, 32, 16, 16)
	w.Player.Body.AddShape(rect)

	space.AddShape(w.Player.Body)

	w.Player.SpeedX = 0
	w.Player.SpeedY = 0

	space.AddShape(resolv.NewRectangle(0, 0, 16, screenHeight))
	space.AddShape(resolv.NewRectangle(screenWidth-16, 0, 16, screenHeight))
	space.AddShape(resolv.NewRectangle(0, 0, screenWidth, 16))
	space.AddShape(resolv.NewRectangle(0, screenHeight-16, screenWidth, 16))

	for _, shape := range space {
		shape.SetTags("solid")
	}

	sticky := resolv.NewRectangle(64, 64, 16, 16)
	sticky.SetTags("sticky", "solid")
	space.AddShape(sticky)

	sticky = resolv.NewRectangle(96, 80, 16, 16)
	sticky.SetTags("sticky", "solid")
	space.AddShape(sticky)

	sticky = resolv.NewRectangle(32, 120, 16, 16)
	sticky.SetTags("sticky", "solid")
	space.AddShape(sticky)

	stickyCircle := resolv.NewCircle(140, 96, 8)
	stickyCircle.SetTags("sticky", "solid")
	space.AddShape(stickyCircle)

	line := resolv.NewLine(160, 100, 170, 100)
	line.SetTags("sticky", "solid")
	space.AddShape(line)

	line = resolv.NewLine(180, 140, 190, 160)
	line.SetTags("sticky", "solid")
	space.AddShape(line)

	w.Player.Body.SetTags("player")

}

func (w *World7) Update() {

	friction := float32(0.5)
	accel := 0.5 + friction

	maxSpd := float32(3)

	if w.Player.SpeedX > friction {
		w.Player.SpeedX -= friction
	} else if w.Player.SpeedX < -friction {
		w.Player.SpeedX += friction
	} else {
		w.Player.SpeedX = 0
	}

	if w.Player.SpeedY > friction {
		w.Player.SpeedY -= friction
	} else if w.Player.SpeedY < -friction {
		w.Player.SpeedY += friction
	} else {
		w.Player.SpeedY = 0
	}

	if keyboard.KeyDown(sdl.K_RIGHT) {
		w.Player.SpeedX += accel
	}

	if keyboard.KeyDown(sdl.K_LEFT) {
		w.Player.SpeedX -= accel
	}

	if keyboard.KeyDown(sdl.K_UP) {
		w.Player.SpeedY -= accel
	}

	if keyboard.KeyDown(sdl.K_DOWN) {
		w.Player.SpeedY += accel
	}

	if w.Player.SpeedX > maxSpd {
		w.Player.SpeedX = maxSpd
	}

	if w.Player.SpeedX < -maxSpd {
		w.Player.SpeedX = -maxSpd
	}

	if w.Player.SpeedY > maxSpd {
		w.Player.SpeedY = maxSpd
	}

	if w.Player.SpeedY < -maxSpd {
		w.Player.SpeedY = -maxSpd
	}

	x := int32(w.Player.SpeedX)

	solids := space.FilterByTags("solid")

	var other resolv.Shape

	if res := solids.Resolve(w.Player.Body, x, 0); res.Colliding() {
		x = res.ResolveX
		w.Player.SpeedX = 0
		other = res.ShapeB
	}

	w.Player.Body.Move(int32(x), 0)

	y := int32(w.Player.SpeedY)

	if res := solids.Resolve(w.Player.Body, 0, y); res.Colliding() {
		y = res.ResolveY
		w.Player.SpeedY = 0
		other = res.ShapeB
	}

	w.Player.Body.Move(0, int32(y))

	if other != nil && other.HasTags("sticky") {
		w.Player.Body.AddShape(other)
		w.Player.Body.SetTags("player")
		space.RemoveShape(other)
	}

	if keyboard.KeyPressed(sdl.K_SPACE) {

		if w.Player.Body.Length() > 1 {
			shape := w.Player.Body.Get(w.Player.Body.Length() - 1) // This is annoying, but I don't know of a way around it
			w.Player.Body.RemoveShape(shape)
			shape.SetTags("solid", "sticky")
			space.AddShape(shape)
		}

	}

	if keyboard.KeyPressed(sdl.K_BACKSPACE) {
		w.Player.Body.SetXY(80, 80)
	}

}

func DrawObject(shape resolv.Shape) {

	if shape.HasTags("player") {
		renderer.SetDrawColor(0, 128, 255, 255)
	} else if shape.HasTags("sticky") {
		renderer.SetDrawColor(0, 255, 100, 255)
	} else {
		renderer.SetDrawColor(255, 255, 255, 255)
	}

	rect, ok := shape.(*resolv.Rectangle)

	if ok {

		renderer.DrawRect(&sdl.Rect{rect.X, rect.Y, rect.W, rect.H})

	}

	circle, ok := shape.(*resolv.Circle)

	if ok {

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

	line, ok := shape.(*resolv.Line)

	if ok {

		renderer.DrawLine(line.X, line.Y, line.X2, line.Y2)

	}

	space, ok := shape.(*resolv.Space)

	if ok {

		for _, o := range *space {
			DrawObject(o)
		}

	}

}

func (w *World7) Draw() {

	shapes := space[:]

	for _, shape := range shapes {
		DrawObject(shape)
	}

	if drawHelpText {
		DrawText(0, 0,
			"Compound shape (space) testing",
			"Use the arrow keys to move",
			"Touch green stuff to stick it",
			"to you",
			"Press Space to detach the",
			"last thing you touched",
			"Note that Circle - line collision",
			"is broken, currently (Sorry!)")
	}

}

func (w *World7) Destroy() {
	space.Clear()
	w.Player = nil
}
