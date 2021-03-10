package main

import (
	"image/color"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
	"github.com/solarlune/resolv"
)

type Player struct {
	Rect *resolv.Line
	// Rect           *resolv.Rectangle
	SpeedX, SpeedY float64
	OnGround       bool
}

func NewPlayer(space *resolv.Space) *Player {

	player := &Player{
		Rect: resolv.NewLine(20, 32, 48, 48),
		// Rect: resolv.NewRectangle(32, 32, 16, 16),
	}

	player.Rect.Tags().Add("player")

	space.Add(player.Rect)

	return player
}

type PlatformerExample struct {
	Player        *Player
	LevelGeometry []*resolv.Rectangle
	Space         *resolv.Space
}

func (example *PlatformerExample) Create() {

	// Shorthand functions
	line := resolv.NewLine
	rect := resolv.NewRectangle

	example.Space = resolv.NewSpace()

	platforms := []resolv.Shape{
		line(64, 180, 90, 180),
		line(64, 140, 90, 140),
		line(200, 240-16, 230, 240-33),
	}

	for _, platform := range platforms {
		platform.Tags().Add("platform")
	}

	example.Space.Add(platforms...)

	solids := []resolv.Shape{

		line(120, 180, 150, 180),
		line(120, 140, 150, 140),

		rect(0, 0, 320, 16),
		rect(0, 240-16, 320, 16),
		rect(0, 16, 16, 240-32),
		rect(320-16, 16, 16, 240-32),
		rect(230, 240-32, 128, 240-32),
	}

	for _, geom := range solids {
		geom.Tags().Add("solid")
	}

	example.Space.Add(solids...)

	example.Player = NewPlayer(example.Space)

}

func (example *PlatformerExample) Update() {

	player := example.Player

	player.SpeedY += 0.5

	accel := float64(0.75)
	maxSpd := float64(2)
	friction := float64(0.5)

	dx := float64(0)

	if ebiten.IsKeyPressed(ebiten.KeyRight) {
		dx++
	}
	if ebiten.IsKeyPressed(ebiten.KeyLeft) {
		dx--
	}

	if inpututil.IsKeyJustPressed(ebiten.KeyX) && player.OnGround {
		player.SpeedY = -8
	}

	player.SpeedX += dx * (accel + friction)

	if player.SpeedX > friction {
		player.SpeedX -= friction
	} else if player.SpeedX < -friction {
		player.SpeedX += friction
	} else {
		player.SpeedX = 0
	}

	if player.SpeedX > maxSpd {
		player.SpeedX = maxSpd
	} else if player.SpeedX < -maxSpd {
		player.SpeedX = -maxSpd
	}

	// Physics

	solids := example.Space.FilterByTags("solid")
	platforms := example.Space.FilterByTags("platform")

	player.OnGround = false

	// Horizontal movement application

	player.Rect.Move(player.SpeedX, 0)

	if check := player.Rect.Check(solids, player.SpeedX, 0); check.Colliding() {
		player.Rect.Move(check.Dx, 0)
		player.SpeedX = 0
	}

	// // Vertical movement application

	player.Rect.Move(0, player.SpeedY)

	// Check against either platforms or solids for ground

	check := player.Rect.Check(platforms, 0, player.SpeedY)

	if player.SpeedY < -1 { // You can jump through platforms
		check = nil
	}

	if check == nil || !check.Colliding() {
		check = player.Rect.Check(solids, 0, player.SpeedY)
	}

	if check.Colliding() {
		player.Rect.Move(0, check.Dy)
		player.SpeedY = 0
		player.OnGround = true
	}

	// if intersection := player.Rect.Check(platforms, 0, player.SpeedY); intersection.Colliding() && player.SpeedY >= 0 {
	// 	// First, one-way platforms (ramps are just platforms in that moving horizontally doesn't stop the player, as we simply resolve on the Y axis only.)
	// 	player.Rect.Y += intersection.Points[0].Y
	// 	player.SpeedY = 0
	// 	player.OnGround = true
	// } else if intersection := player.Rect.Check(solids, 0, player.SpeedY); intersection.Colliding() {
	// 	// And then other solids
	// 	player.Rect.Y += intersection.Points[0].Y
	// 	player.SpeedY = 0
	// 	player.OnGround = true
	// }

}

func (example *PlatformerExample) Draw(screen *ebiten.Image) {

	for _, shape := range *example.Space {

		drawColor := color.RGBA{128, 128, 128, 255}
		if shape.Tags().Has("player") {
			drawColor = color.RGBA{0, 255, 0, 255}
		} else if shape.Tags().Has("platform") {
			drawColor = color.RGBA{0, 255, 255, 255}
		}

		switch cast := shape.(type) {

		case *resolv.Rectangle:

			drawRect(screen, cast, drawColor)

		case *resolv.Line:

			ebitenutil.DrawLine(screen, cast.X, cast.Y, cast.X2, cast.Y2, drawColor)

		}

	}

}

func (example *PlatformerExample) Destroy() {

}
