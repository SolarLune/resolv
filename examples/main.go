package main

import (
	_ "embed"
	"errors"
	"fmt"
	"image/color"
	"time"

	"github.com/golang/freetype/truetype"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
	"github.com/hajimehoshi/ebiten/v2/text"
	"github.com/solarlune/resolv"
	"golang.org/x/image/font"
)

//go:embed excel.ttf
var excelFont []byte

type Game struct {
	Worlds        []WorldInterface
	CurrentWorld  int
	Width, Height int
	Debug         bool
	ShowHelpText  bool
	Screen        *ebiten.Image
	FontFace      font.Face
	Time          float64
}

func NewGame() *Game {

	ebiten.SetWindowResizable(true)
	ebiten.SetWindowTitle("resolv test")

	g := &Game{
		Width:        640,
		Height:       360,
		ShowHelpText: true,
	}

	g.Worlds = []WorldInterface{
		NewWorldBouncer(g),
		NewWorldPlatformer(g),
		NewWorldLineTest(g),
		// NewWorldMultiShape(g), // MultiShapes are still buggy; gotta fix 'em up
		NewWorldShapeTest(g),
		NewWorldDirectTest(g),
	}

	fontData, _ := truetype.Parse(excelFont)

	opts := &truetype.Options{
		Size:    10,
		DPI:     72,
		Hinting: font.HintingFull,
	}

	g.FontFace = truetype.NewFace(fontData, opts)

	// Debug FPS rendering

	go func() {

		for {

			fmt.Println("FPS: ", ebiten.CurrentFPS())
			fmt.Println("Ticks: ", ebiten.CurrentTPS())
			time.Sleep(time.Second)

		}

	}()

	return g

}

func (g *Game) Update() error {

	var quit error

	if inpututil.IsKeyJustPressed(ebiten.KeyP) {
		if ebiten.ActualTPS() >= 60 {
			ebiten.SetTPS(6)
		} else {
			ebiten.SetTPS(60)
		}
	}

	if inpututil.IsKeyJustPressed(ebiten.KeyF1) {
		g.Debug = !g.Debug
	}

	if inpututil.IsKeyJustPressed(ebiten.KeyF2) {
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
	g.Worlds[g.CurrentWorld].Draw(screen)
}

func (g *Game) DrawText(screen *ebiten.Image, x, y int, textLines ...string) {
	rectHeight := 10
	for _, txt := range textLines {
		w := float64(font.MeasureString(g.FontFace, txt).Round())
		ebitenutil.DrawRect(screen, float64(x), float64(y-8), w, float64(rectHeight), color.RGBA{0, 0, 0, 192})

		text.Draw(screen, txt, g.FontFace, x+1, y+1, color.RGBA{0, 0, 150, 255})
		text.Draw(screen, txt, g.FontFace, x, y, color.RGBA{100, 150, 255, 255})
		y += rectHeight
	}
}

func (g *Game) DebugDraw(screen *ebiten.Image, space *resolv.Space) {

	for y := 0; y < space.Height(); y++ {

		for x := 0; x < space.Width(); x++ {

			cell := space.Cell(x, y)

			cw := float64(space.CellWidth)
			ch := float64(space.CellHeight)
			cx := float64(cell.X) * cw
			cy := float64(cell.Y) * ch

			drawColor := color.RGBA{20, 20, 20, 255}

			if cell.Occupied() {
				drawColor = color.RGBA{255, 255, 0, 255}
			}

			ebitenutil.DrawLine(screen, cx, cy, cx+cw, cy, drawColor)

			ebitenutil.DrawLine(screen, cx+cw, cy, cx+cw, cy+ch, drawColor)

			ebitenutil.DrawLine(screen, cx+cw, cy+ch, cx, cy+ch, drawColor)

			ebitenutil.DrawLine(screen, cx, cy+ch, cx, cy, drawColor)
		}

	}

}

func (g *Game) Layout(w, h int) (int, int) {
	return g.Width, g.Height
}

func main() {
	ebiten.RunGame(NewGame())
}
