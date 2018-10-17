package main

import (
	"github.com/SolarLune/resolv/resolv"
	"github.com/veandco/go-sdl2/sdl"
)

type World5 struct {
	TargetLine *resolv.Line
}

func (w *World5) Create() {

	space.Clear()

	w.TargetLine = resolv.NewLine(screenWidth/2, screenHeight/2, 0, 0)

	space.AddShape(w.TargetLine)

	other := resolv.NewLine(0, 0, 100, 100)
	space.AddShape(other)

	rect := resolv.NewRectangle(160, 16, 32, 32)
	space.AddShape(rect)

}

func (w *World5) Update() {

	x, y, btn := sdl.GetMouseState()

	winW, winH := window.GetSize()

	ratioX := float32(screenWidth) / float32(winW)
	ratioY := float32(screenHeight) / float32(winH)

	mx := int32(float32(x) * ratioX)
	my := int32(float32(y) * ratioY)

	if btn == sdl.Button(sdl.BUTTON_LEFT) {
		w.TargetLine.X = mx
		w.TargetLine.Y = my
	} else {
		w.TargetLine.X2 = mx
		w.TargetLine.Y2 = my
	}

}

func (w *World5) Draw() {

	for _, shape := range space {

		line, ok := shape.(*resolv.Line)

		if ok {

			if line == w.TargetLine {
				if space.IsColliding(line) {
					renderer.SetDrawColor(255, 0, 0, 255)
				} else {
					renderer.SetDrawColor(0, 255, 0, 255)
				}
			} else {
				renderer.SetDrawColor(255, 255, 255, 255)
			}
			renderer.DrawLine(line.X, line.Y, line.X2, line.Y2)

		}

		rect, ok := shape.(*resolv.Rectangle)

		if ok {

			renderer.SetDrawColor(255, 255, 255, 255)
			renderer.DrawRect(&sdl.Rect{rect.X, rect.Y, rect.W, rect.H})

		}

	}

	if drawHelpText {
		DrawText("Click to place the line's start", 0, 0)
		DrawText("Move the mouse to place the end point", 0, 16)
		DrawText("The line turns red when it touches", 0, 32)
		DrawText("something", 0, 48)
		DrawText("Press F1 to hide this text", 0, 64)
	}

}

func (w *World5) Destroy() {
	space.Clear()
	w.TargetLine = nil
}
