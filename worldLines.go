package main

import (
	"strconv"

	rl "github.com/gen2brain/raylib-go/raylib"

	"github.com/SolarLune/resolv/resolv"
)

type WorldLines struct {
	TargetLine *resolv.Line
	Space      *resolv.Space
}

func (w *WorldLines) Create() {

	w.Space = resolv.NewSpace()

	w.Space.Clear()

	w.TargetLine = resolv.NewLine(screenWidth/2, screenHeight/2, 0, 0)

	w.Space.Add(w.TargetLine)

	other := resolv.NewLine(0, 0, 100, 100)
	w.Space.Add(other)

	rect := resolv.NewRectangle(160, 16, 32, 32)
	w.Space.Add(rect)

	var lx, ly int32 = 160, 160
	var ls int32 = 16

	line := resolv.NewLine(lx, ly, lx+ls, ly)
	w.Space.Add(line)

	line = resolv.NewLine(lx+ls, ly, lx+ls, ly+ls)
	w.Space.Add(line)

	line = resolv.NewLine(lx+ls, ly+ls, lx, ly+ls)
	w.Space.Add(line)

	line = resolv.NewLine(lx, ly+ls, lx, ly)
	w.Space.Add(line)

}

func (w *WorldLines) Update() {

	// x, y, btn := sdl.GetMouseState()

	x, y := rl.GetMouseX(), rl.GetMouseY()

	winW, winH := rl.GetScreenWidth(), rl.GetScreenHeight()

	ratioX := float32(screenWidth) / float32(winW)
	ratioY := float32(screenHeight) / float32(winH)

	mx := int32(float32(x) * ratioX)
	my := int32(float32(y) * ratioY)

	if rl.IsMouseButtonPressed(rl.MouseLeftButton) {
		w.TargetLine.X = mx
		w.TargetLine.Y = my
	} else {
		w.TargetLine.X2 = mx
		w.TargetLine.Y2 = my
	}

}

func (w *WorldLines) Draw() {

	for _, shape := range *w.Space {

		line, ok := shape.(*resolv.Line)

		if ok {

			drawColor := rl.White

			if line == w.TargetLine {
				if w.Space.IsColliding(line) {
					for i, point := range line.GetIntersectionPoints(w.Space) {
						rl.DrawLine(point.X-5, point.Y-5, point.X+5, point.Y+5, rl.Yellow)
						rl.DrawLine(point.X+5, point.Y-5, point.X-5, point.Y+5, rl.Yellow)
						DrawText(point.X+5, point.Y, "Intersection #"+strconv.Itoa(i+1))
					}
					drawColor = rl.Red
				} else {
					drawColor = rl.Green
				}
			}
			rl.DrawLine(line.X, line.Y, line.X2, line.Y2, drawColor)

		}

		rect, ok := shape.(*resolv.Rectangle)

		if ok {

			rl.DrawRectangleLines(rect.X, rect.Y, rect.W, rect.H, rl.LightGray)

		}

	}

	if drawHelpText {
		DrawText(32, 16, "-Line collision test-",
			"Click to place the line's start.",
			"Move the mouse to place the end point.",
			"The line turns red when it touches",
			"something.",
			"Press F1 to hide this text.")

		DrawText(160, 130, "This square is made out of",
			"individual lines; the inside is",
			"''hollow''.")
	}

}

func (w *WorldLines) Destroy() {
	w.Space.Clear()
	w.TargetLine = nil
}
