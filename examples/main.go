package main

import (
	"bytes"
	_ "embed"
	"errors"
	"fmt"
	"image/color"
	"os"
	"runtime/pprof"
	"strconv"
	"time"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
	"github.com/hajimehoshi/ebiten/v2/text/v2"
	"github.com/hajimehoshi/ebiten/v2/vector"
	"github.com/solarlune/resolv"
	"golang.org/x/image/font/gofont/goregular"
)

type Game struct {
	Worlds        []WorldInterface
	CurrentWorld  int
	Width, Height int
	DebugSpace    bool
	ShowHelpText  bool
	Screen        *ebiten.Image
	FontFace      text.Face
	Time          float64
}

func NewGame() *Game {

	ebiten.SetWindowResizingMode(ebiten.WindowResizingModeEnabled)
	ebiten.SetWindowTitle("resolv test")

	g := &Game{
		Width:        640,
		Height:       360,
		ShowHelpText: true,
	}

	g.Worlds = []WorldInterface{
		NewWorldPlatformer(),
		NewWorldBouncer(),
		NewWorldCircle(),
	}

	// g.Worlds = []WorldInterface{
	// 	NewWorldBouncer(g),
	// 	NewWorldPlatformer(g),
	// 	NewWorldLineTest(g),
	// 	// NewWorldMultiShape(g), // MultiShapes are still buggy; gotta fix 'em up
	// 	NewWorldShapeTest(g),
	// 	NewWorldDirectTest(g),
	// }

	// g.FontFace = truetype.NewFace(fontData, opts)

	faceSrc, err := text.NewGoTextFaceSource(bytes.NewReader(goregular.TTF))
	if err != nil {
		panic(err)
	}

	g.FontFace = &text.GoTextFace{
		Source: faceSrc,
		Size:   15,
	}

	// Debug FPS rendering

	go func() {

		for {

			fmt.Println("FPS: ", ebiten.ActualFPS())
			fmt.Println("Ticks: ", ebiten.ActualTPS())
			time.Sleep(time.Second)

		}

	}()

	return g

}

func (g *Game) Update() error {

	var quit error

	if inpututil.IsKeyJustPressed(ebiten.KeyF) {
		fmt.Println("toggle slow-mo")
		if ebiten.TPS() >= 60 {
			ebiten.SetTPS(6)
		} else {
			ebiten.SetTPS(60)
		}
	}

	if inpututil.IsKeyJustPressed(ebiten.KeyP) {
		g.StartProfiling()
	}

	if inpututil.IsKeyJustPressed(ebiten.KeyF2) {
		g.DebugSpace = !g.DebugSpace
	}

	if inpututil.IsKeyJustPressed(ebiten.KeyF1) {
		g.ShowHelpText = !g.ShowHelpText
	}

	if inpututil.IsKeyJustPressed(ebiten.KeyF4) {
		ebiten.SetFullscreen(!ebiten.IsFullscreen())
	}

	if inpututil.IsKeyJustPressed(ebiten.KeyE) {
		g.CurrentWorld++
	}

	if inpututil.IsKeyJustPressed(ebiten.KeyQ) {
		g.CurrentWorld--
	}

	if g.CurrentWorld >= len(g.Worlds) {
		g.CurrentWorld = 0
	} else if g.CurrentWorld < 0 {
		g.CurrentWorld = len(g.Worlds) - 1
	}

	world := g.Worlds[g.CurrentWorld]

	if inpututil.IsKeyJustPressed(ebiten.KeyR) {
		fmt.Println("Restart World")
		world.Init()
	}

	if ebiten.IsKeyPressed(ebiten.KeyEscape) {
		quit = errors.New("quit")
	}

	world.Update()

	g.Time += 1.0 / 60.0

	return quit

}

func (g *Game) Draw(screen *ebiten.Image) {
	g.Screen = screen
	screen.Fill(color.RGBA{20, 20, 40, 255})

	world := g.Worlds[g.CurrentWorld]
	world.Draw(screen)

	if g.ShowHelpText {
		g.DrawText(screen, 0, 0,
			"Press F1 to show and hide help text.",
			"Press F2 to debug draw the Space.",
			"Q and E switch to different Worlds.",
			"FPS: "+strconv.FormatFloat(ebiten.ActualFPS(), 'f', 1, 32),
			"TPS: "+strconv.FormatFloat(ebiten.ActualTPS(), 'f', 1, 32),
		)
	}
	if g.DebugSpace {
		g.DebugDraw(screen, world.Space())
	}
}

func (g *Game) DrawText(screen *ebiten.Image, x, y int, textLines ...string) {
	metrics := g.FontFace.Metrics()
	for _, txt := range textLines {
		w, h := text.Measure(txt, g.FontFace, 16)
		vector.DrawFilledRect(screen, float32(x+2), float32(y), float32(w), float32(h), color.RGBA{0, 0, 0, 192}, false)

		opt := text.DrawOptions{}
		opt.GeoM.Translate(float64(x+2), float64(y+2-int(metrics.VDescent)-4))
		opt.Filter = ebiten.FilterNearest
		// opt.ColorScale.ScaleWithColor(color.RGBA{0, 0, 150, 255})

		text.Draw(screen, txt, g.FontFace, &opt)
		// text.Draw(screen, txt, g.FontFace, x, y, color.RGBA{100, 150, 255, 255})
		y += 16
	}
}

func (g *Game) DebugDraw(screen *ebiten.Image, space *resolv.Space) {

	for y := 0; y < space.Height(); y++ {

		for x := 0; x < space.Width(); x++ {

			cell := space.Cell(x, y)

			cw := float32(space.CellWidth())
			ch := float32(space.CellHeight())
			cx := float32(cell.X) * cw
			cy := float32(cell.Y) * ch

			drawColor := color.RGBA{20, 20, 20, 255}

			if cell.IsOccupied() {
				drawColor = color.RGBA{255, 255, 0, 255}
			}

			vector.StrokeRect(screen, cx, cy, cx+cw, cy, 2, drawColor, false)

			vector.StrokeRect(screen, cx+cw, cy, cx+cw, cy+ch, 2, drawColor, false)

			vector.StrokeRect(screen, cx+cw, cy+ch, cx, cy+ch, 2, drawColor, false)

			vector.StrokeRect(screen, cx, cy+ch, cx, cy, 2, drawColor, false)
		}

	}

}

func (g *Game) StartProfiling() {
	outFile, err := os.Create("./cpu.pprof")
	if err != nil {
		fmt.Println(err.Error())
		return
	}

	fmt.Println("Beginning CPU profiling...")
	pprof.StartCPUProfile(outFile)
	go func() {
		time.Sleep(5 * time.Second)
		pprof.StopCPUProfile()
		fmt.Println("CPU profiling finished.")
	}()
}

func (g *Game) Layout(w, h int) (int, int) {
	return g.Width, g.Height
}

func main() {
	GlobalGame = NewGame()
	ebiten.RunGame(GlobalGame)
}

var GlobalGame *Game
