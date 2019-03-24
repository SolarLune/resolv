package main

import (
	"fmt"
	"strconv"

	"github.com/SolarLune/resolv/resolv"
	"github.com/veandco/go-sdl2/gfx"
	"github.com/veandco/go-sdl2/sdl"
	"github.com/veandco/go-sdl2/ttf"
)

var screenWidth int32 = 320
var screenHeight int32 = 240
var cell int32 = 4

var space *resolv.Space
var renderer *sdl.Renderer
var window *sdl.Window
var avgFramerate int

var drawHelpText = true

func main() {

	// defer profile.Start(profile.ProfilePath(".")).Stop()

	sdl.Init(sdl.INIT_EVERYTHING)
	defer sdl.Quit()

	ttf.Init()
	defer ttf.Quit()

	window, renderer, _ = sdl.CreateWindowAndRenderer(screenWidth, screenHeight, sdl.WINDOW_SHOWN|sdl.WINDOW_RESIZABLE)

	window.SetResizable(true)

	renderer.SetLogicalSize(screenWidth, screenHeight)

	fpsMan := &gfx.FPSmanager{}
	gfx.InitFramerate(fpsMan)
	gfx.SetFramerate(fpsMan, 60)

	// Change this to one of the other World structs to change the world and see different tests

	space = resolv.NewSpace()

	var world WorldInterface = &World1{}

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

		if keyboard.KeyPressed(sdl.K_F2) {
			gfx.SetFramerate(fpsMan, 10)
		} else if keyboard.KeyPressed(sdl.K_F3) {
			gfx.SetFramerate(fpsMan, 60)
		}

		if keyboard.KeyPressed(sdl.K_1) {
			world.Destroy()
			world = &World1{}
			world.Create()
		}
		if keyboard.KeyPressed(sdl.K_2) {
			world.Destroy()
			world = &World2{}
			world.Create()
		}
		if keyboard.KeyPressed(sdl.K_3) {
			world.Destroy()
			world = &World3{}
			world.Create()
		}
		if keyboard.KeyPressed(sdl.K_4) {
			world.Destroy()
			world = &World4{}
			world.Create()
		}
		if keyboard.KeyPressed(sdl.K_5) {
			world.Destroy()
			world = &World5{}
			world.Create()
		}
		if keyboard.KeyPressed(sdl.K_6) {
			world.Destroy()
			world = &World6{}
			world.Create()
		}
		if keyboard.KeyPressed(sdl.K_7) {
			world.Destroy()
			world = &World7{}
			world.Create()
		}

		if keyboard.KeyPressed(sdl.K_F1) {
			drawHelpText = !drawHelpText
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

		DrawText(screenWidth-32, 0, strconv.Itoa(avgFramerate))

		renderer.Present()

	}

}

func DrawText(x, y int32, textLines ...string) {

	sy := y

	for _, text := range textLines {

		font, _ := ttf.OpenFont("ARCADEPI.TTF", 12)
		defer font.Close()

		var surf *sdl.Surface

		surf, _ = font.RenderUTF8Solid(text, sdl.Color{R: 50, G: 100, B: 255, A: 255})

		textSurface, _ := renderer.CreateTextureFromSurface(surf)
		defer textSurface.Destroy()

		_, _, w, h, _ := textSurface.Query()

		renderer.Copy(textSurface, &sdl.Rect{X: 0, Y: 0, W: w, H: h}, &sdl.Rect{X: x, Y: sy, W: w, H: h})

		sy += 16

	}

}
