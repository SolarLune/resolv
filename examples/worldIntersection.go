package main

import (
	"image/color"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/solarlune/resolv"
)

type IntersectionExample struct {
	Player        *Player
	LevelGeometry []*resolv.Rectangle
	Space         *resolv.Space
}

func (example *IntersectionExample) Create() {

	// Shorthand functions
	// line := resolv.NewLine
	rect := resolv.NewRectangle

	example.Space = resolv.NewSpace()

	solids := []resolv.Shape{

		rect(0, 0, 320, 16),
		rect(0, 240-16, 320, 16),
		rect(0, 16, 16, 240-32),
		rect(320-16, 16, 16, 240-32),
	}

	for _, geom := range solids {
		geom.Tags().Add("solid")
	}

	example.Space.Add(
		rect(80, 80, 32, 128),
		rect(128, 96, 64, 64),
	)

	example.Space.Add(solids...)

	example.Player = NewPlayer(example.Space)

}

func (example *IntersectionExample) Update() {

	player := example.Player

	accel := float64(0.5)
	maxSpd := float64(1.5)
	friction := float64(0.3)

	dx := float64(0)
	dy := float64(0)

	if ebiten.IsKeyPressed(ebiten.KeyRight) {
		dx++
	}
	if ebiten.IsKeyPressed(ebiten.KeyLeft) {
		dx--
	}
	if ebiten.IsKeyPressed(ebiten.KeyUp) {
		dy--
	}
	if ebiten.IsKeyPressed(ebiten.KeyDown) {
		dy++
	}

	player.SpeedX += dx * (accel + friction)
	player.SpeedY += dy * (accel + friction)

	if player.SpeedX > friction {
		player.SpeedX -= friction
	} else if player.SpeedX < -friction {
		player.SpeedX += friction
	} else {
		player.SpeedX = 0
	}

	if player.SpeedY > friction {
		player.SpeedY -= friction
	} else if player.SpeedY < -friction {
		player.SpeedY += friction
	} else {
		player.SpeedY = 0
	}

	if player.SpeedX > maxSpd {
		player.SpeedX = maxSpd
	} else if player.SpeedX < -maxSpd {
		player.SpeedX = -maxSpd
	}

	if player.SpeedY > maxSpd {
		player.SpeedY = maxSpd
	} else if player.SpeedY < -maxSpd {
		player.SpeedY = -maxSpd
	}

	// Physics

	solids := example.Space.FilterByTags("solid")

	player.Rect.X += player.SpeedX

	if movement := player.Rect.Check(solids, player.SpeedX, 0); movement.Colliding() {
		player.Rect.X += movement.Dx
		player.SpeedX = 0
	}

	player.Rect.Y += player.SpeedY

	if movement := player.Rect.Check(solids, 0, player.SpeedY); movement.Colliding() {
		player.Rect.Y += movement.Dy
		player.SpeedY = 0
	}

}

func (example *IntersectionExample) Draw(screen *ebiten.Image) {

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

	if intersection := example.Player.Rect.Check(example.Space, 0, 0); intersection.Colliding() {

		points := intersection.Points
		rect := resolv.NewRectangle(points[0].X, points[0].Y, points[2].X-points[0].X, points[2].Y-points[0].Y)
		drawRect(screen, rect, color.RGBA{255, 127, 0, 255})

	}

}

func (example *IntersectionExample) Destroy() {

}
