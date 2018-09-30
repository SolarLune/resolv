package main

import (
	"fmt"
	"strconv"

	"github.com/SolarLune/resolv/resolv"
	"github.com/veandco/go-sdl2/gfx"
	"github.com/veandco/go-sdl2/sdl"
	"github.com/veandco/go-sdl2/ttf"
)

var screenWidth int32 = 640
var screenHeight int32 = 480

var space resolv.Space
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

	renderer.SetLogicalSize(screenWidth, screenHeight)

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

		if keyboard.KeyPressed(sdl.K_1) {
			gfx.SetFramerate(fpsMan, 10)
		} else if keyboard.KeyPressed(sdl.K_2) {
			gfx.SetFramerate(fpsMan, 60)
		}

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

		DrawText(strconv.Itoa(avgFramerate), screenWidth-32, 0)

		renderer.Present()

	}

}

func DrawText(text string, x, y int32) {

	font, _ := ttf.OpenFont("ARCADEPI.TTF", 12)
	defer font.Close()

	var surf *sdl.Surface

	surf, _ = font.RenderUTF8Solid(text, sdl.Color{R: 255, G: 255, B: 255, A: 255})

	textSurface, _ := renderer.CreateTextureFromSurface(surf)
	defer textSurface.Destroy()

	_, _, w, h, _ := textSurface.Query()

	textSurface.SetAlphaMod(100)
	renderer.Copy(textSurface, &sdl.Rect{X: 0, Y: 0, W: w, H: h}, &sdl.Rect{X: x, Y: y, W: w, H: h})

}
