package main

import (
	"fmt"
	"image/color"
	"math"

	"github.com/hajimehoshi/ebiten/inpututil"
	"github.com/tanema/gween"
	"github.com/tanema/gween/ease"

	"github.com/SolarLune/resolv"
	"github.com/hajimehoshi/ebiten"
	"github.com/hajimehoshi/ebiten/ebitenutil"
)

type Player struct {
	Object          *resolv.Object
	GroundDetection *resolv.Object
	SpeedX          float64
	SpeedY          float64
	OnGround        bool
	WallSliding     *resolv.Object
	FacingRight     bool
	IgnorePlatform  *resolv.Object
}

func NewPlayer(space *resolv.Space) *Player {
	p := &Player{
		Object:          resolv.NewObject(32, 128, 16, 24, space),
		GroundDetection: resolv.NewObject(32, 128, 16, 1, space),
		FacingRight:     true,
	}
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

	world.Space = resolv.NewSpace(int(gw/16), int(gh/16), 16, 16)

	// Construct geometry
	resolv.NewObject(0, 0, 16, gh, world.Space)
	resolv.NewObject(gw-16, 0, 16, gh, world.Space)
	resolv.NewObject(0, 0, gw, 16, world.Space)
	resolv.NewObject(0, gh-24, gw, 32, world.Space)
	resolv.NewObject(160, gh-56, 160, 32, world.Space)
	resolv.NewObject(320, 64, 32, 160, world.Space)
	resolv.NewObject(64, 128, 16, 160, world.Space)
	resolv.NewObject(gw-128, 64, 128, 16, world.Space)
	resolv.NewObject(gw-128, gh-88, 128, 16, world.Space)

	for _, o := range world.Space.Objects {
		o.AddTag("solid")
	}

	world.Player = NewPlayer(world.Space)
	world.Player.Object.PreciseCollision = true

	world.FloatingPlatform = resolv.NewObject(128, gh-32, 128, 8, world.Space)
	world.FloatingPlatform.AddTag("platform")
	world.FloatingPlatformTween = gween.NewSequence()
	world.FloatingPlatformTween.Add(
		gween.New(float32(world.FloatingPlatform.Y), float32(world.FloatingPlatform.Y-128), 2, ease.Linear),
		gween.New(float32(world.FloatingPlatform.Y-128), float32(world.FloatingPlatform.Y), 2, ease.Linear),
	)

	// Platforms
	resolv.NewObject(352, 64, 48, 8, world.Space).AddTag("platform")
	resolv.NewObject(352, 64+64, 48, 8, world.Space).AddTag("platform")
	resolv.NewObject(352, 64+128, 48, 8, world.Space).AddTag("platform")
	resolv.NewObject(352, 64+192, 48, 8, world.Space).AddTag("platform")

	// Ramps
	resolv.NewObject(320, gh-56, 16, 16, world.Space).AddTag("ramp")
	resolv.NewObject(336, gh-40, 16, 16, world.Space).AddTag("ramp")

}

func (world *WorldPlatformer) Update(screen *ebiten.Image) {

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
		if ebiten.IsKeyPressed(ebiten.KeyRight) {
			player.SpeedX += accel
			player.FacingRight = true
		}

		if ebiten.IsKeyPressed(ebiten.KeyLeft) {
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

	platformCheck := player.Object.Check(0, player.SpeedY+1, "platform")
	var platform *resolv.Object
	if len(platformCheck.ObjectsByTags("platform")) > 0 {
		platform = platformCheck.ObjectsByTags("platform")[0]
		if platform.Y < player.Object.Y+player.Object.H-4 {
			platform = nil
		}
	}

	if inpututil.IsKeyJustPressed(ebiten.KeyX) {

		if ebiten.IsKeyPressed(ebiten.KeyDown) && platform != nil {

			player.IgnorePlatform = platform

		} else {

			if player.OnGround {
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

	if platform != nil && platform != player.IgnorePlatform && player.SpeedY >= 0 {
		player.Object.Y += platformCheck.ContactY
		player.SpeedY = 0
		player.OnGround = true
		player.WallSliding = nil

	} else if check := player.Object.Check(0, player.SpeedY, "ramp"); check.Valid() && player.SpeedY >= 0 {
		ramp := check.Objects()[0]
		dx := math.Abs(player.Object.X-ramp.X-ramp.W) / ramp.W
		dx += 0.1
		if dx > 1 {
			dx = 1
		} else if dx < 0 {
			dx = 0
		}

		player.Object.SetBottom(ramp.Y + ramp.H)
		player.Object.Y -= dx * ramp.H
		player.SpeedY = 0
		player.OnGround = true
		player.WallSliding = nil

	} else if check := player.Object.Check(0, player.SpeedY, "solid"); check.Valid() {
		if player.SpeedY < 0 && check.CanSlide && math.Abs(check.SlideX) < 8 {
			player.Object.X += check.SlideX // This allows you to slide around a block if you're just baaarely off
		} else {
			player.Object.Y += check.ContactY // Move to contact with the surface; ContactY takes into account the size of the Player automatically
			player.SpeedY = 0
			player.OnGround = true
			player.WallSliding = nil
		}
	} else {
		player.Object.Y += player.SpeedY
		player.OnGround = false
	}

	if check := player.Object.Check(player.SpeedX, 0, "solid"); check.Valid() {
		player.Object.X += check.ContactX
		player.SpeedX = 0
		if !player.OnGround {
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

	if player.OnGround {
		player.IgnorePlatform = nil
	}

	player.Object.Update() // Update the player's position in the space

}

func (world *WorldPlatformer) Draw(screen *ebiten.Image) {

	screen.Fill(color.RGBA{20, 20, 40, 255})

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
