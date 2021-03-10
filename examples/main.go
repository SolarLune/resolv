package main

import (
	"errors"
	"fmt"
	"image/color"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
	"github.com/solarlune/resolv"
)

type WorldInterface interface {
	Create()
	Update()
	Draw(screen *ebiten.Image)
	Destroy()
}

type Game struct {
	World WorldInterface
}

func NewGame() *Game {

	ebiten.SetWindowResizable(true)
	ebiten.SetWindowTitle("Resolv Examples")
	world := &PlatformerExample{}
	// world := &IntersectionExample{}
	world.Create()
	return &Game{
		World: world,
	}

}

func (game *Game) Update() error {

	if inpututil.IsKeyJustPressed(ebiten.KeyR) {
		fmt.Println("Restarted Example")
		game.World.Destroy()
		game.World.Create()
	}

	game.World.Update()

	if inpututil.IsKeyJustPressed(ebiten.KeyEscape) {
		return errors.New("quit")
	}

	return nil

}

func (game *Game) Draw(screen *ebiten.Image) {

	screen.Fill(color.RGBA{20, 30, 40, 255})

	game.World.Draw(screen)

}

func (game *Game) Layout(w, h int) (int, int) {
	return 320, 240
}

func main() {

	game := NewGame()

	if err := ebiten.RunGame(game); err != nil {
		panic(err)
	}

}

func drawRect(screen *ebiten.Image, rect *resolv.Rectangle, color color.Color) {

	x, y, w, h := rect.X, rect.Y, rect.W, rect.H

	ebitenutil.DrawLine(screen, x, y, x+w, y, color)
	ebitenutil.DrawLine(screen, x+w, y, x+w, y+h, color)
	ebitenutil.DrawLine(screen, x+w, y+h, x, y+h, color)
	ebitenutil.DrawLine(screen, x, y+h, x, y, color)
}
