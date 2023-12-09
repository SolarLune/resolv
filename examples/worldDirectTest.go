package main

import (
	"fmt"
	"image/color"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/solarlune/resolv"
)

type WorldDirectTest struct {
	Game         *Game
	Geometry     []*resolv.ConvexPolygon
	ShowHelpText bool
	LineStartPos resolv.Vector
}

func NewWorldDirectTest(game *Game) *WorldDirectTest {

	w := &WorldDirectTest{
		Game:         game,
		ShowHelpText: true,
		LineStartPos: resolv.NewVector(float64(game.Width)/2, float64(game.Height)/2),
	}

	w.Init()

	return w
}

func (world *WorldDirectTest) Init() {

	smallBox := resolv.NewConvexPolygonVec(resolv.NewVectorZero(), // Position

		resolv.NewVector(-5, -5), // Vertices
		resolv.NewVector(5, -5),
		resolv.NewVector(5, 5),
		resolv.NewVector(-5, 5),
	)

	type boxSetup struct {
		Position resolv.Vector
		Size     resolv.Vector
		Rotation float64 // in degrees
	}

	boxes := []boxSetup{

		{
			Position: resolv.NewVector(150, 150),
			Size:     resolv.NewVector(4, 4),
		},

		{
			Position: resolv.NewVector(200, 250),
			Rotation: 45,
		},

		{
			Position: resolv.NewVector(220, 250),
		},

		// Big boi
		{
			Position: resolv.NewVector(300, 200),
			Size:     resolv.NewVector(20, 10),
			Rotation: 10,
		},

		{
			Position: resolv.NewVector(300, 50),
		},

		{
			Position: resolv.NewVector(320, 60),
		},
	}

	world.Geometry = []*resolv.ConvexPolygon{}

	for _, box := range boxes {
		newBox := smallBox.Clone().(*resolv.ConvexPolygon)
		newBox.SetPositionVec(box.Position)
		if box.Size.X > 0 {
			newBox.SetScaleVec(box.Size)
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

	if ebiten.IsKeyPressed(ebiten.KeyP) {
		world.Geometry[0].SetRotation(0)
	}

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

	moveSpd := 1.25

	world.LineStartPos.X += dx * moveSpd
	world.LineStartPos.Y += dy * moveSpd

}

func (world *WorldDirectTest) Draw(screen *ebiten.Image) {

	lineColor := color.RGBA{255, 255, 0, 255}
	mx, my := ebiten.CursorPosition()
	line := resolv.NewLine(world.LineStartPos.X, world.LineStartPos.Y, float64(mx), float64(my))

	intersectionPoints := []resolv.Vector{}

	for _, box := range world.Geometry {
		if intersection := line.Intersection(0, 0, box); intersection != nil {
			intersectionPoints = append(intersectionPoints, intersection.Points...)
			lineColor = color.RGBA{255, 0, 0, 255}
		}
	}

	l := line.Lines()[0]
	ebitenutil.DrawLine(screen, l.Start.X, l.Start.Y, l.End.X, l.End.Y, lineColor)
	DrawBigDot(screen, world.LineStartPos, lineColor)

	for _, o := range world.Geometry {
		DrawPolygon(screen, o, color.White)
	}

	for _, point := range intersectionPoints {
		DrawBigDot(screen, point, color.RGBA{0, 255, 0, 255})
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
