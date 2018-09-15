package main

import (
	"fmt"

	"github.com/SolarLune/resolv/resolv"
	"github.com/veandco/go-sdl2/gfx"
	"github.com/veandco/go-sdl2/sdl"
	"github.com/veandco/go-sdl2/ttf"
)

var screenWidth int32 = 320
var screenHeight int32 = 240

var space resolv.Space

var mainSquare *resolv.Rectangle
var squareSpeedX float32
var squareSpeedY float32

var renderer *sdl.Renderer
var avgFramerate int

func main() {

	sdl.Init(sdl.INIT_EVERYTHING)
	defer sdl.Quit()

	ttf.Init()
	defer ttf.Quit()

	var win *sdl.Window

	win, renderer, _ = sdl.CreateWindowAndRenderer(screenWidth, screenHeight, sdl.WINDOW_SHOWN|sdl.WINDOW_RESIZABLE)

	win.SetResizable(true)

	renderer.SetLogicalSize(320, 240)

	fpsMan := &gfx.FPSmanager{}
	gfx.InitFramerate(fpsMan)
	gfx.SetFramerate(fpsMan, 60)

	// Change this to one of the other World structs to change the world and see different tests
	world := World1{}

	world.Create()

	running := true

	var frameCount int
	var framerateDelay uint32

	for running {

		for event := sdl.PollEvent(); event != nil; event = sdl.PollEvent() {
			switch event.(type) {
			case *sdl.QuitEvent:
				running = false
			case *sdl.KeyboardEvent:
				keyboard.ReportEvent(event.(*sdl.KeyboardEvent))
			}
		}

		keyboard.Update()

		if keyboard.KeyPressed(sdl.K_ESCAPE) {
			running = false
		}

		if keyboard.KeyPressed(sdl.K_r) {
			world.Create()
		}

		world.Update()

		renderer.SetDrawColor(20, 30, 40, 255)

		renderer.Clear()

		world.Draw()

		framerateDelay += gfx.FramerateDelay(fpsMan)

		if framerateDelay >= 1000 {
			avgFramerate = frameCount
			framerateDelay -= 1000
			frameCount = 0
			fmt.Println(avgFramerate, " FPS")
		}

		frameCount++

		renderer.Present()

	}

}
