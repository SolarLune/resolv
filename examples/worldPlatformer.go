package main

import (
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
	"github.com/solarlune/resolv"
)

var (
	TagPlayer    = resolv.NewTag("Player")
	TagSolidWall = resolv.NewTag("SolidWall")
	TagPlatform  = resolv.NewTag("Platform")
)

type WorldPlatformer struct {
	space            *resolv.Space
	Player           *PlatformingPlayer
	MovingPlatform   *resolv.ConvexPolygon
	PlatformMovingUp bool
}

func NewWorldPlatformer() *WorldPlatformer {
	// Create the world.
	w := &WorldPlatformer{}
	// Initialize it.
	w.Init()
	return w
}

func (w *WorldPlatformer) Init() {

	// Create the space. It is 640x360 large (the size of the screen), and divided into 16x16 cells.
	// The cell division makes it more efficient to check for shapes.
	w.space = resolv.NewSpace(640, 360, 16, 16)

	solids := resolv.ShapeCollection{
		resolv.NewCircle(128, 128, 32),

		resolv.NewRectangleFromTopLeft(0, 0, 640, 8),
		resolv.NewRectangleFromTopLeft(640-8, 8, 8, 360-16),

		resolv.NewRectangleFromTopLeft(0, 8, 8, 360-32),
		resolv.NewRectangleFromTopLeft(0, 360-8-16, 8, 8),
		resolv.NewRectangleFromTopLeft(0, 360-8-8, 8, 8),

		resolv.NewRectangleFromTopLeft(64, 200, 300, 8),
		resolv.NewRectangleFromTopLeft(64, 280, 300, 8),
		resolv.NewRectangleFromTopLeft(512, 96, 32, 200),

		resolv.NewRectangleFromTopLeft(0, 360-8, 640, 8),
	}

	solids.SetTags(TagSolidWall | TagPlatform)

	w.space.Add(solids...)

	/////

	platforms := resolv.ShapeCollection{
		resolv.NewRectangleFromTopLeft(400, 200, 32, 16),
		resolv.NewRectangleFromTopLeft(400, 240, 32, 16),
		resolv.NewRectangleFromTopLeft(400, 280, 32, 16),
		resolv.NewRectangleFromTopLeft(400, 320, 32, 16),
	}

	platforms.SetTags(TagPlatform)

	w.space.Add(platforms...)

	////

	w.Player = NewPlayer(w.space)

	ramp := resolv.NewConvexPolygon(180, 175,
		[]float64{
			-24, 8,
			8, -8,
			48, -8,
			80, 8,
		},
	)
	ramp.Tags().Set(TagPlatform | TagSolidWall)
	w.space.Add(ramp)

	// Clone and move the Ramp, then place it again
	r := ramp.Clone()
	r.SetPositionVec(resolv.NewVector(240, 344))
	w.space.Add(r)

	w.MovingPlatform = resolv.NewRectangleFromTopLeft(550, 200, 32, 8)
	w.MovingPlatform.Tags().Set(TagPlatform)
	w.space.Add(w.MovingPlatform)

}

func (w *WorldPlatformer) Update() {

	movingPlatformSpeed := 2.0

	if w.PlatformMovingUp {
		w.MovingPlatform.Move(0, -movingPlatformSpeed)
	} else {
		w.MovingPlatform.Move(0, movingPlatformSpeed)
	}

	if w.MovingPlatform.Position().Y <= 200 || w.MovingPlatform.Position().Y > 300 {
		w.PlatformMovingUp = !w.PlatformMovingUp
	}

	w.Player.Update()

}

func (w *WorldPlatformer) Draw(screen *ebiten.Image) {
	CommonDraw(screen, w)
	if GlobalGame.ShowHelpText {
		GlobalGame.DrawText(screen, 0, 128,
			"Platformer Test",
			"Left and right arrow keys to move",
			"X to jump",
			"You can walljump",
			"Orange platforms can be jumped through",
		)
	}
}

// To allow the world's physical state to be drawn using the debug draw function.
func (w *WorldPlatformer) Space() *resolv.Space {
	return w.space
}

type PlatformingPlayer struct {
	Object   *resolv.ConvexPolygon
	Movement resolv.Vector
	Facing   resolv.Vector
	YSpeed   float64
	Space    *resolv.Space

	OnGround    bool
	WallSliding bool
}

func NewPlayer(space *resolv.Space) *PlatformingPlayer {

	player := &PlatformingPlayer{
		Object: resolv.NewRectangle(192, 128, 16, 12),
		Space:  space,
	}
	player.Object.Tags().Set(TagPlayer)
	space.Add(player.Object)
	return player

}

func (p *PlatformingPlayer) Update() {

	moveVec := resolv.Vector{}
	gravity := 0.5
	friction := 0.2
	accel := 0.5 + friction
	maxSpd := 4.0
	jumpSpd := 8.0

	if !p.WallSliding {

		// Only move if you're not on the wall
		p.YSpeed += gravity

		if ebiten.IsKeyPressed(ebiten.KeyLeft) {
			moveVec.X -= accel
		}

		if ebiten.IsKeyPressed(ebiten.KeyRight) {
			moveVec.X += accel
		}

	} else {
		p.YSpeed = 0.2 // Slide down the wall slowly
	}

	// Jumping
	if inpututil.IsKeyJustPressed(ebiten.KeyX) {

		// Jump either if you're on the ground or on a wall
		if p.OnGround || p.WallSliding {
			p.YSpeed = -jumpSpd
		}

		// Jump away from the wall
		if p.WallSliding {
			if p.Facing.X < 0 {
				moveVec.X = maxSpd
			} else {
				moveVec.X = -maxSpd
			}
			p.WallSliding = false
		}

	}

	// Set the facing if the player's not on a wall and attempting to move
	if !p.WallSliding && !moveVec.IsZero() {
		p.Facing = moveVec.Unit()
	}

	// Add in the player's movement, clamping it to the maximum speed and incorporating friction.
	p.Movement = p.Movement.Add(moveVec).ClampMagnitude(maxSpd).SubMagnitude(friction)

	// Filter out shapes that are nearby the player
	nearbyShapes := p.Object.SelectTouchingCells(4).FilterShapes()

	p.OnGround = false

	checkVec := resolv.NewVector(0, p.YSpeed) // Check downwards by the distance of movement speed

	if p.YSpeed >= 0 {
		checkVec.Y += 4 // Add in a bit of extra downwards cast to account for running down ramps if we're not jumping
	}

	// Snap to ground using a shape-based line test.
	p.Object.ShapeLineTest(resolv.ShapeLineTestSettings{
		Vector:      checkVec,
		TestAgainst: nearbyShapes.ByTags(TagPlatform), // Select the shapes near the player object that are platforms
		OnIntersect: func(set resolv.IntersectionSet, index, max int) bool {

			if p.YSpeed >= 0 && set.Intersections[0].Normal.Y < 0 {
				// If we're falling and landing on upward facing line

				p.OnGround = true                 // Then set on ground to true
				p.WallSliding = false             // And wallsliding to false
				p.YSpeed = 0                      // Stop vertical movement
				p.Object.MoveVec(set.MTV.SubY(2)) // Move to contact plus a bit of floating to not be flush with the ground so running up ramps is easier
				return false                      // We can stop iterating past this

			} else if set.Intersections[0].Normal.Y > 0 && p.YSpeed < 0 && set.OtherShape.Tags().Has(TagSolidWall) {
				// Jumping and bonking on downward-facing line and it's solid
				p.YSpeed = 0
				p.Object.MoveVec(set.MTV)
				return false // We can stop iterating past this.
			}

			// No ground or ceiling, so keep looking for collisions
			return true

		},
	})

	// Apply movement - Y speed is separate so that gravity can take effect separate from horizontal movement speed
	p.Object.Move(p.Movement.X, p.Movement.Y+p.YSpeed)

	// Collision test first

	wallslideSet := false

	p.Object.IntersectionTest(resolv.IntersectionTestSettings{

		// Check shapes in surrounding cells that have the "TagSolidWall" tag.
		TestAgainst: nearbyShapes.ByTags(TagSolidWall),

		OnIntersect: func(set resolv.IntersectionSet) bool {

			// If we're not on the ground and attempting to move and the touched nearest contact in the contact set's normal is opposite
			// the direction we're attempting to move (i.e. moving left into a right-facing wall), then it's a wall-slide candidate.
			if !wallslideSet && !p.OnGround && !moveVec.IsZero() && set.Intersections[0].Normal.Dot(moveVec) < 0 {

				// If we just barely clip a wall, that shouldn't count. We want to touch at least half the height of the player's worth of wall to enter wall-slide mode
				if dist := set.Distance(resolv.WorldUp); dist >= p.Object.Bounds().Height() && p.YSpeed >= -2 {

					// Set the wall as the object we're wallsliding against
					p.WallSliding = true

					// Stop movement
					p.Movement.X = 0

					// Set wallslide to true so that we don't have to check any other objects.
					wallslideSet = true

				}

			}

			// Move away from whatever is struck.
			p.Object.MoveVec(set.MTV)

			// Keep iterating in case we're touching something else.
			return true
		},
	})

	// We can also use an IntersectionTest / ShapeLineTest / LineTest function as a bool if we don't need to do anything special on any particular "one"
	if !p.Object.ShapeLineTest(resolv.ShapeLineTestSettings{
		Vector:      p.Facing,
		TestAgainst: nearbyShapes.ByTags(TagSolidWall),
	}) {
		p.WallSliding = false
	}

}
