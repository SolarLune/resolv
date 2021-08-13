package main

import (
	"fmt"
	"image/color"
	"math"

	"github.com/hajimehoshi/ebiten/v2/inpututil"
	"github.com/tanema/gween"
	"github.com/tanema/gween/ease"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/solarlune/resolv"
)

type Player struct {
	Object         *resolv.Object
	SpeedX         float64
	SpeedY         float64
	OnGround       *resolv.Object
	WallSliding    *resolv.Object
	FacingRight    bool
	IgnorePlatform *resolv.Object
}

func NewPlayer(space *resolv.Space) *Player {

	p := &Player{
		Object:      resolv.NewObject(32, 128, 16, 24),
		FacingRight: true,
	}

	space.Add(p.Object)

	return p
}

type WorldPlatformer struct {
	Game   *Game
	Space  *resolv.Space
	Player *Player

	FloatingPlatform      *resolv.Object
	FloatingPlatformTween *gween.Sequence

	ShowHelpText bool
}

func NewWorldPlatformer(game *Game) *WorldPlatformer {

	w := &WorldPlatformer{Game: game, ShowHelpText: true}
	w.Init()
	return w

}

func (world *WorldPlatformer) Init() {

	gw := float64(world.Game.Width)
	gh := float64(world.Game.Height)

	world.Space = resolv.NewSpace(int(gw), int(gh), 16, 16)

	// Construct geometry
	world.Space.Add(
		resolv.NewObject(0, 0, 16, gh),
		resolv.NewObject(gw-16, 0, 16, gh),
		resolv.NewObject(0, 0, gw, 16),
		resolv.NewObject(0, gh-24, gw, 32),
		resolv.NewObject(160, gh-56, 160, 32),
		resolv.NewObject(320, 64, 32, 160),
		resolv.NewObject(64, 128, 16, 160),
		resolv.NewObject(gw-128, 64, 128, 16),
		resolv.NewObject(gw-128, gh-88, 128, 16),
	)

	for _, o := range world.Space.Objects {
		o.AddTags("solid")
	}

	world.Player = NewPlayer(world.Space)

	world.FloatingPlatform = resolv.NewObject(128, gh-32, 128, 8)
	world.FloatingPlatform.AddTags("platform")
	world.FloatingPlatformTween = gween.NewSequence()
	world.FloatingPlatformTween.Add(
		gween.New(float32(world.FloatingPlatform.Y), float32(world.FloatingPlatform.Y-128), 2, ease.Linear),
		gween.New(float32(world.FloatingPlatform.Y-128), float32(world.FloatingPlatform.Y), 2, ease.Linear),
	)

	world.Space.Add(world.FloatingPlatform)

	// Platforms
	platforms := []*resolv.Object{
		resolv.NewObject(352, 64, 48, 8),
		resolv.NewObject(352, 64+64, 48, 8),
		resolv.NewObject(352, 64+128, 48, 8),
		resolv.NewObject(352, 64+192, 48, 8),
	}
	for _, platform := range platforms {
		platform.AddTags("platform")
	}

	world.Space.Add(platforms...)

	// Ramps
	ramps := []*resolv.Object{
		resolv.NewObject(320, gh-56, 16, 16),
		resolv.NewObject(336, gh-40, 16, 16),
	}

	for _, ramp := range ramps {
		ramp.AddTags("ramp")
	}

	world.Space.Add(ramps...)

}

func (world *WorldPlatformer) Update() {

	// fmt.Println(ebiten.GamepadIDs())

	// for i := ebiten.GamepadID; i < 16; i++ {
	// 	fmt.Println(ebiten.GamepadName(i), ebiten.GamepadAxisNum(0), ebiten.GamepadButtonNum(0), ebiten.GamepadSDLID(0))
	// }

	if inpututil.IsKeyJustPressed(ebiten.KeyF2) {
		world.ShowHelpText = !world.ShowHelpText
	}

	// Playform movement needs to be done first to make sure there's no space between the top and the player's bottom
	y, _, seqDone := world.FloatingPlatformTween.Update(1.0 / 60.0)
	world.FloatingPlatform.Y = float64(y)
	if seqDone {
		world.FloatingPlatformTween.Reset()
	}
	world.FloatingPlatform.Update()

	player := world.Player

	friction := 0.5
	accel := 0.5 + friction
	maxSpeed := 4.0
	jumpSpd := 10.0
	gravity := 0.75

	player.SpeedY += gravity

	if player.WallSliding != nil && player.SpeedY > 1 {
		player.SpeedY = 1
	}

	if player.WallSliding == nil {
		if ebiten.IsKeyPressed(ebiten.KeyRight) || ebiten.GamepadAxis(0, 0) > 0.1 {
			player.SpeedX += accel
			player.FacingRight = true
		}

		if ebiten.IsKeyPressed(ebiten.KeyLeft) || ebiten.GamepadAxis(0, 0) < -0.1 {
			player.SpeedX -= accel
			player.FacingRight = false
		}
	}

	if player.SpeedX > friction {
		player.SpeedX -= friction
	} else if player.SpeedX < -friction {
		player.SpeedX += friction
	} else {
		player.SpeedX = 0
	}

	if player.SpeedX > maxSpeed {
		player.SpeedX = maxSpeed
	} else if player.SpeedX < -maxSpeed {
		player.SpeedX = -maxSpeed
	}

	if inpututil.IsKeyJustPressed(ebiten.KeyX) || ebiten.IsGamepadButtonPressed(0, 0) || ebiten.IsGamepadButtonPressed(1, 0) {

		if (ebiten.IsKeyPressed(ebiten.KeyDown) || ebiten.GamepadAxis(0, 1) > 0.1 || ebiten.GamepadAxis(1, 1) > 0.1) && player.OnGround != nil && player.OnGround.HasTags("platform") {

			player.IgnorePlatform = player.OnGround

		} else {

			if player.OnGround != nil {
				player.SpeedY = -jumpSpd
			} else if player.WallSliding != nil {
				// WALLJUMPING
				player.SpeedY = -jumpSpd

				if player.WallSliding.X > player.Object.X {
					player.SpeedX = -4
				} else {
					player.SpeedX = 4
				}

				player.WallSliding = nil

			}

		}

	}

	// Platform check is done early because we need to 1) prioritize standing on platforms if they can move beneath solid objects like
	// in this example, and we need to know if a platform is ignored or not before standing on it

	player.OnGround = nil
	dy := player.SpeedY

	// Lock vertical movement to a maximum of the size of the grid so we don't miss any collisions
	dy = math.Max(math.Min(dy, 16), -16)

	if check := player.Object.Check(0, dy, "solid"); check.Valid() {

		if dy < 0 && check.Slide.Valid && math.Abs(check.Slide.Vector[0]) < 8 {
			player.Object.X += check.Slide.Vector[0] // This allows you to slide around a block if you're just baaarely off
		} else {
			// Move to contact with the surface; ContactY takes into account the size of the Player automatically
			dy = check.Contact.Vector[1]
			player.SpeedY = 0
			player.OnGround = check.Objects()[0]
			player.WallSliding = nil

			if player.OnGround != player.IgnorePlatform {
				player.IgnorePlatform = nil
			}

		}

	}

	// for _, o := range check.Objects() {

	// 	if o == player.Object {
	// 		continue
	// 	}

	// 	fmt.Println(o.X, o.Y, o.Tags())

	// 	if o.HasTags("solid") {

	// 		if player.SpeedY < 0 && check.CanSlide && math.Abs(check.SlideX) < 8 {
	// 			player.Object.X += check.SlideX // This allows you to slide around a block if you're just baaarely off
	// 		} else {
	// 			_, dy = check.ContactDelta() // Move to contact with the surface; ContactY takes into account the size of the Player automatically
	// 			player.SpeedY = 0
	// 			player.OnGround = o
	// 			player.WallSliding = nil

	// 			if player.OnGround != player.IgnorePlatform {
	// 				player.IgnorePlatform = nil
	// 			}

	// 		}

	// 	} else if o.HasTags("platform") {

	// 		if player.SpeedY >= 0 && o.Y >= player.Object.Y+player.Object.H-4 && o != player.IgnorePlatform {
	// 			dy = check.ContactY // Move to contact with the surface; ContactY takes into account the size of the Player automatically
	// 			player.SpeedY = 0
	// 			player.OnGround = o
	// 			player.WallSliding = nil
	// 		}

	// 	}

	// }

	player.Object.Y += dy

	if check := player.Object.Check(player.SpeedX, 0, "solid"); check.Valid() {
		player.Object.X += check.Contact.Vector[0]
		player.SpeedX = 0
		if player.OnGround == nil {
			player.WallSliding = check.ObjectsByTags("solid")[0]
		}
	} else {
		player.Object.X += player.SpeedX
	}

	// If the wall runs out, stop wall sliding
	wallRight := 1.0
	if !player.FacingRight {
		wallRight = -1
	}

	if c := player.Object.Check(wallRight, 0, "solid"); !c.Valid() {
		player.WallSliding = nil
	}

	player.Object.Update() // Update the player's position in the space

}

func (world *WorldPlatformer) Draw(screen *ebiten.Image) {

	for _, o := range world.Space.Objects {
		drawColor := color.RGBA{60, 60, 60, 255}
		if o.HasTags("platform") {
			drawColor = color.RGBA{180, 100, 0, 255}
		} else if o.HasTags("ramp") {
			drawColor = color.RGBA{255, 50, 100, 255}
		}
		ebitenutil.DrawRect(screen, o.X, o.Y, o.W, o.H, drawColor)
	}

	player := world.Player.Object
	ebitenutil.DrawRect(screen, player.X, player.Y, player.W, player.H, color.RGBA{0, 255, 60, 255})

	if world.Game.Debug {
		world.Game.DebugDraw(screen, world.Space)
	}

	if world.ShowHelpText {

		world.Game.DrawText(screen, 16, 16,
			"~ Platformer Demo ~",
			"Move Player: Left, Right Arrow",
			"Jump: X Key",
			"Wallslide: Move into wall in air",
			"Walljump: Jump while wallsliding",
			"Fall through platforms: Down + X",
			"",
			"F1: Toggle Debug View",
			"F2: Show / Hide help text",
			"R: Restart world",
			"E: Next world",
			"Q: Previous world",
			fmt.Sprintf("%d FPS (frames per second)", int(ebiten.CurrentFPS())),
			fmt.Sprintf("%d TPS (ticks per second)", int(ebiten.CurrentTPS())),
		)

	}

}
