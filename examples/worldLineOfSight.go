package main

import (
	"fmt"
	"image/color"
	"strconv"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/solarlune/resolv"
)

type WorldLineTest struct {
	Game   *Game
	Space  *resolv.Space
	Player *resolv.Object
}

func NewWorldLineTest(game *Game) *WorldLineTest {
	w := &WorldLineTest{Game: game}
	w.Init()
	return w
}

func (world *WorldLineTest) Init() {

	gw := float64(world.Game.Width)
	gh := float64(world.Game.Height)

	cellSize := 8

	world.Space = resolv.NewSpace(int(gw), int(gh), cellSize, cellSize)

	// Construct geometry
	geometry := []*resolv.Object{

		resolv.NewObject(0, 0, 16, gh),
		resolv.NewObject(gw-16, 0, 16, gh),
		resolv.NewObject(0, 0, gw, 16),
		resolv.NewObject(0, gh-24, gw, 32),
		resolv.NewObject(0, gh-24, gw, 32),

		resolv.NewObject(200, -160, 16, gh),
	}

	world.Space.Add(geometry...)

	for _, o := range world.Space.Objects() {
		o.AddTags("solid")
	}

	world.Player = resolv.NewObject(160, 160, 16, 16)
	world.Player.AddTags("player")
	world.Space.Add(world.Player)

}

func (world *WorldLineTest) Update() {

	dx, dy := 0.0, 0.0
	moveSpd := 2.0

	if ebiten.IsKeyPressed(ebiten.KeyW) {
		dy = -moveSpd
	}

	if ebiten.IsKeyPressed(ebiten.KeyS) {
		dy += moveSpd
	}

	if ebiten.IsKeyPressed(ebiten.KeyA) {
		dx = -moveSpd
	}

	if ebiten.IsKeyPressed(ebiten.KeyD) {
		dx += moveSpd
	}

	if col := world.Player.Check(dx, 0, "solid"); col != nil {
		dx = col.ContactWithObject(col.Objects[0]).X
	}

	world.Player.Position.X += dx

	if col := world.Player.Check(0, dy, "solid"); col != nil {
		dy = col.ContactWithObject(col.Objects[0]).Y
	}

	world.Player.Position.Y += dy

	world.Player.Update()

}

func (world *WorldLineTest) Draw(screen *ebiten.Image) {

	for _, o := range world.Space.Objects() {
		drawColor := color.RGBA{60, 60, 60, 255}
		if o.HasTags("player") {
			drawColor = color.RGBA{0, 255, 0, 255}
		}
		ebitenutil.DrawRect(screen, o.Position.X, o.Position.Y, o.Size.X, o.Size.Y, drawColor)
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

		if interrupted {
			break
		}

	}

	if world.Game.Debug {
		world.Game.DebugDraw(screen, world.Space)
	}

	if world.Game.ShowHelpText {

		world.Game.DrawText(screen, 16, 16,
			"~ Line of sight test ~",
			"WASD keys: Move player",
			"Mouse: Hover over impassible objects",
			"to get the closest wall to the player.",
			fmt.Sprintf("Mouse X: %d, Mouse Y: %d", mouseX, mouseY),
			"Clear line of sight: "+strconv.FormatBool(!interrupted),
			"",
			"F1: Toggle Debug View",
			"F2: Show / Hide help text",
			"F4: Toggle fullscreen",
			"R: Restart world",
			"E: Next world",
			"Q: Previous world",
			fmt.Sprintf("%d FPS (frames per second)", int(ebiten.CurrentFPS())),
			fmt.Sprintf("%d TPS (ticks per second)", int(ebiten.CurrentTPS())),
		)

	}

}
