package main

import (
	"fmt"
	"image/color"

	"github.com/hajimehoshi/ebiten/inpututil"
	"github.com/kvartborg/vector"

	"github.com/SolarLune/resolv"
	"github.com/hajimehoshi/ebiten"
	"github.com/hajimehoshi/ebiten/ebitenutil"
)

type WorldLineTest struct {
	Game          *Game
	Space         *resolv.Space
	ShowHelpText  bool
	Player        *resolv.Object
	CollidingCell *resolv.Cell
}

func NewWorldLineTest(game *Game) *WorldLineTest {
	w := &WorldLineTest{Game: game, ShowHelpText: true}
	w.Init()
	return w
}

func (world *WorldLineTest) Init() {

	gw := float64(world.Game.Width)
	gh := float64(world.Game.Height)

	cellSize := 8

	world.Space = resolv.NewSpace(int(gw)/cellSize, int(gh)/cellSize, cellSize, cellSize)

	// Construct geometry
	resolv.NewObject(0, 0, 16, gh, world.Space)
	resolv.NewObject(gw-16, 0, 16, gh, world.Space)
	resolv.NewObject(0, 0, gw, 16, world.Space)
	resolv.NewObject(0, gh-24, gw, 32, world.Space)
	resolv.NewObject(0, gh-24, gw, 32, world.Space)

	resolv.NewObject(200, -160, 16, gh, world.Space)

	for _, o := range world.Space.Objects {
		o.AddTag("solid")
	}

	world.Player = resolv.NewObject(160, 160, 16, 16, world.Space)

	world.Player.AddTag("player")

}

func (world *WorldLineTest) Update(screen *ebiten.Image) {

	if inpututil.IsKeyJustPressed(ebiten.KeyF2) {
		world.ShowHelpText = !world.ShowHelpText
	}

	moveVec := vector.Vector{0, 0}
	moveSpd := 2.0

	if ebiten.IsKeyPressed(ebiten.KeyW) {
		moveVec[1] = -moveSpd
	}

	if ebiten.IsKeyPressed(ebiten.KeyS) {
		moveVec[1] += moveSpd
	}

	if ebiten.IsKeyPressed(ebiten.KeyA) {
		moveVec[0] = -moveSpd
	}

	if ebiten.IsKeyPressed(ebiten.KeyD) {
		moveVec[0] += moveSpd
	}

	if col := world.Player.Check(moveVec[0], 0, "solid"); col.Valid() {
		world.Player.X += col.ContactX
	} else {
		world.Player.X += moveVec[0]
	}

	if col := world.Player.Check(0, moveVec[1], "solid"); col.Valid() {
		world.Player.Y += col.ContactY
	} else {
		world.Player.Y += moveVec[1]
	}

	world.Player.Update()

}

func (world *WorldLineTest) Draw(screen *ebiten.Image) {

	screen.Fill(color.RGBA{20, 20, 40, 255})

	for _, o := range world.Space.Objects {
		drawColor := color.RGBA{60, 60, 60, 255}
		if o.HasTags("player") {
			drawColor = color.RGBA{0, 255, 0, 255}
		}
		ebitenutil.DrawRect(screen, o.X, o.Y, o.W, o.H, drawColor)
	}

	mouseX, mouseY := ebiten.CursorPosition()

	mx, my := world.Space.WorldToSpace(float64(mouseX), float64(mouseY))

	cx, cy := world.Player.CellPosition()

	sightLine := world.Space.CellsInLine(cx, cy, mx, my)

	interrupted := false

	for i, cell := range sightLine {

		if i == 0 { // Skip the beginning because that's the player
			continue
		}

		drawColor := color.RGBA{255, 255, 0, 255}

		// if interrupted {
		// 	drawColor = color.RGBA{0, 0, 255, 255}
		// }

		if !interrupted && cell.ContainsTags("solid") {
			drawColor = color.RGBA{255, 0, 0, 255}
			interrupted = true
		}

		ebitenutil.DrawRect(screen,
			float64(cell.X*world.Space.CellWidth),
			float64(cell.Y*world.Space.CellHeight),
			float64(world.Space.CellWidth),
			float64(world.Space.CellHeight),
			drawColor)

	}

	if world.Game.Debug {
		world.Game.DebugDraw(screen, world.Space)
	}

	if world.ShowHelpText {

		world.Game.DrawText(screen, 16, 16,
			"~ Line of sight test ~",
			"WASD keys: Move player",
			"Mouse: Hover over impassible objects",
			"to get the closest wall to the player.",
			"",
			"F1: Toggle Debug View",
			"F2: Show / Hide help text",
			"R: Restart world",
			"E: Next world",
			"Q: Previous world",
			fmt.Sprintf("Mouse X: %d, Mouse Y: %d", mouseX, mouseY),
			fmt.Sprintf("%d FPS (frames per second)", int(ebiten.CurrentFPS())),
			fmt.Sprintf("%d TPS (ticks per second)", int(ebiten.CurrentTPS())),
		)

	}

}
