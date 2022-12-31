package main

import (
	"fmt"
	"image/color"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/quartercastle/vector"
	"github.com/solarlune/resolv"
)

type WorldDirectTest struct {
	Game         *Game
	Geometry     []*resolv.ConvexPolygon
	ShowHelpText bool
	LineStartPos vector.Vector
}

func NewWorldDirectTest(game *Game) *WorldDirectTest {

	w := &WorldDirectTest{
		Game:         game,
		ShowHelpText: true,
		LineStartPos: vector.Vector{float64(game.Width) / 2, float64(game.Height) / 2},
	}

	w.Init()

	return w
}

func (world *WorldDirectTest) Init() {

	smallBox := resolv.NewConvexPolygon(
		0, 0, // Position

		-5, -5, // Vertices
		5, -5,
		5, 5,
		-5, 5,
	)

	type boxSetup struct {
		X, Y     float64
		W, H     float64
		Rotation float64 // in degrees
	}

	boxes := []boxSetup{

		{
			X: 150, Y: 150,
		},

		{
			X: 200, Y: 250,
			Rotation: 45,
		},

		{
			X: 220, Y: 250,
		},

		// Big boi
		{
			X: 300, Y: 200,
			W: 20, H: 10,
			Rotation: 10,
		},

		{
			X: 300,
			Y: 50,
		},

		{
			X: 320,
			Y: 60,
		},
	}

	world.Geometry = []*resolv.ConvexPolygon{}

	for _, box := range boxes {
		newBox := smallBox.Clone().(*resolv.ConvexPolygon)
		newBox.SetPosition(box.X, box.Y)
		if box.W > 0 {
			newBox.SetScale(box.W, box.H)
		}
		newBox.SetRotation(resolv.ToRadians(box.Rotation))
		world.Geometry = append(world.Geometry,
			newBox,
		)
	}

}

func (world *WorldDirectTest) Update() {

	// Let's rotate one of them because why not
	world.Geometry[0].Rotate(0.01)

	dx := 0.0
	dy := 0.0

	if ebiten.IsKeyPressed(ebiten.KeyA) {
		dx = -1
	}
	if ebiten.IsKeyPressed(ebiten.KeyD) {
		dx = 1
	}

	if ebiten.IsKeyPressed(ebiten.KeyW) {
		dy = -1
	}
	if ebiten.IsKeyPressed(ebiten.KeyS) {
		dy = 1
	}

	moveSpd := 0.5

	world.LineStartPos[0] += dx * moveSpd
	world.LineStartPos[1] += dy * moveSpd

}

func (world *WorldDirectTest) Draw(screen *ebiten.Image) {

	lineColor := color.RGBA{255, 255, 0, 255}
	mx, my := ebiten.CursorPosition()
	line := resolv.NewLine(world.LineStartPos[0], world.LineStartPos[1], float64(mx), float64(my))

	intersectionPoints := []vector.Vector{}

	for _, box := range world.Geometry {
		if intersection := line.Intersection(0, 0, box); intersection != nil {
			intersectionPoints = append(intersectionPoints, intersection.Points...)
			lineColor = color.RGBA{255, 0, 0, 255}
		}
	}

	l := line.Lines()[0]
	ebitenutil.DrawLine(screen, l.Start[0], l.Start[1], l.End[0], l.End[1], lineColor)
	DrawBigDot(screen, world.LineStartPos[0], world.LineStartPos[1], lineColor)

	for _, o := range world.Geometry {
		DrawPolygon(screen, o, color.White)
	}

	for _, point := range intersectionPoints {
		DrawBigDot(screen, point[0], point[1], color.RGBA{0, 255, 0, 255})
	}

	if world.Game.ShowHelpText {

		world.Game.DrawText(screen, 16, 16,
			"~ Direct Test Demo ~",
			"",
			"This demo tests out direct collision between",
			"objects. A line is cast from the center of the",
			"screen and goes to the mouse position.",
			"The intersection points of all objects that cross the line",
			"is visualized by green dots.",
			"WASD: Move line start.",
			"Mouse position: Move line end.",
			fmt.Sprintf("%d FPS (frames per second)", int(ebiten.CurrentFPS())),
			fmt.Sprintf("%d TPS (ticks per second)", int(ebiten.CurrentTPS())),
			"",
			"F2: Show / Hide help text",
			"F4: Toggle fullscreen",
			"R: Restart world",
			"E: Next world",
			"Q: Previous world",
		)

	}

}
