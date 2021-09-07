package main

// // Multishapes are jank; they need to be fixed.

// import (
// 	"fmt"
// 	"image/color"
// 	"math"
// 	"math/rand"

// 	"github.com/hajimehoshi/ebiten/v2"
// 	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
// 	"github.com/solarlune/resolv"
// )

// type WorldMultiShape struct {
// 	Space              *resolv.Space
// 	Game               *Game
// 	MovingObj          *resolv.Object
// 	ShowHelpText       bool
// 	PreciseCollisionOn bool
// }

// func NewWorldMultiShape(g *Game) *WorldMultiShape {
// 	w := &WorldMultiShape{Game: g, ShowHelpText: true}
// 	w.Init()
// 	return w
// }

// func (world *WorldMultiShape) Init() {

// 	gw := float64(world.Game.Width)
// 	gh := float64(world.Game.Height)
// 	cellSize := 8

// 	world.Space = resolv.NewSpace(int(gw), int(gh), cellSize, cellSize)

// 	world.MovingObj = resolv.NewObject(320, 32, 1, 1)

// 	multiShape := resolv.NewMultiShape()
// 	multiShape.Add(resolv.NewRectangle(0, 0, 16, 64))
// 	multiShape.Add(resolv.NewRectangle(16, 64, 32, 8))
// 	world.MovingObj.Shape = multiShape
// 	world.Space.Add(world.MovingObj)

// 	world.MovingObj.SetBounds(multiShape.Bounds())

// 	for i := 0; i < 50; i++ {
// 		cx := rand.Float64() * float64(world.Game.Width)
// 		cy := rand.Float64() * float64(world.Game.Height)
// 		if !world.Space.Cell(world.Space.WorldToSpace(cx, cy)).Occupied() {
// 			obj := resolv.NewObject(cx, cy, 8, 8)
// 			obj.Shape = resolv.NewRectangle(0, 0, 8, 8)
// 			world.Space.Add(obj)
// 		} else {
// 			i-- // Try again~
// 		}
// 	}

// }

// func (world *WorldMultiShape) Update() {

// 	dx := 0.0
// 	dy := 0.0

// 	if ebiten.IsKeyPressed(ebiten.KeyRight) {
// 		dx++
// 	}

// 	if ebiten.IsKeyPressed(ebiten.KeyLeft) {
// 		dx--
// 	}

// 	if ebiten.IsKeyPressed(ebiten.KeyUp) {
// 		dy--
// 	}

// 	if ebiten.IsKeyPressed(ebiten.KeyDown) {
// 		dy++
// 	}

// 	if check := world.MovingObj.Check(dx, 0); check != nil {

// 		for _, o := range check.Objects {

// 			if contactSet := world.MovingObj.Shape.Intersection(dx, 0, o.Shape); contactSet != nil {
// 				dist := contactSet.MTV.Magnitude()
// 				dx = contactSet.MTV.Unit().Scale(dist + 0.1)[0]
// 				fmt.Println(contactSet)
// 				break
// 			}

// 		}

// 	}

// 	world.MovingObj.X += dx

// 	if check := world.MovingObj.Check(0, dy); check != nil {

// 		for _, o := range check.Objects {

// 			multiShape := world.MovingObj.Shape.(*resolv.MultiShape)
// 			if contactSet := multiShape.Intersection(0, dy, o.Shape); contactSet != nil {
// 				dy = contactSet.MTV.Y()
// 				break

// 			}

// 		}

// 	}

// 	world.MovingObj.Y += dy

// 	world.MovingObj.Update()

// }

// func (world *WorldMultiShape) Draw(screen *ebiten.Image) {

// 	for _, o := range world.Space.Objects {
// 		ebitenutil.DrawRect(screen, o.X, o.Y, o.W, o.H, color.RGBA{60, 60, 60, 255})
// 	}

// 	for _, shape := range world.MovingObj.Shape.(*resolv.MultiShape).Shapes {

// 		if c, ok := shape.(*resolv.Circle); ok {
// 			DrawCircle(screen, c, color.White)
// 		} else if p, ok := shape.(*resolv.ConvexPolygon); ok {
// 			DrawPolygon(screen, p, color.White)
// 		}

// 	}

// 	// o := world.MovingObj
// 	// ebitenutil.DrawRect(screen, o.X, o.Y, o.W, o.H, color.RGBA{0, 255, 64, 255})

// 	if world.Game.Debug {
// 		world.Game.DebugDraw(screen, world.Space)
// 	}

// 	if world.ShowHelpText {

// 		world.Game.DrawText(screen, 16, 16,
// 			"~ Precision Demo ~",
// 			"",
// 			fmt.Sprint("Precise Collision: ", world.PreciseCollisionOn),
// 			"",
// 			"For precise collisions, one can use resolv's built-in shape-checking functions ",
// 			"to see if a true collision happens between objects within the same Cell.",
// 			"Otherwise, objects collide based on overall cellular locations.",
// 			"",
// 			"Space: Turn on or off precise collisions",
// 			"on the moving red rectangle.",
// 			"F1: Toggle Debug View",
// 			"F2: Show / Hide help text",
// 			"R: Restart world",
// 			"E: Next world",
// 			"Q: Previous world",
// 			fmt.Sprintf("%d FPS (frames per second)", int(ebiten.CurrentFPS())),
// 			fmt.Sprintf("%d TPS (ticks per second)", int(ebiten.CurrentTPS())),
// 		)

// 	}

// }

// func (world *WorldMultiShape) DrawCircle(screen *ebiten.Image, circle *resolv.Circle, drawColor color.Color) {

// 	// Actually drawing the circles live is too inefficient, so we will simply draw them to an image and then draw that instead
// 	// when necessary.

// 	if _, exists := circleBuffer[circle]; !exists {
// 		newImg := ebiten.NewImage(int(circle.Radius)*2, int(circle.Radius)*2)

// 		newImg.Set(int(circle.X), int(circle.Y), color.White)

// 		stepCount := float64(32)

// 		// Half image width and height.
// 		hw := circle.Radius
// 		hh := circle.Radius

// 		for i := 0; i < int(stepCount); i++ {

// 			x := (math.Sin(math.Pi*2*float64(i)/stepCount) * (circle.Radius - 2)) + hw
// 			y := (math.Cos(math.Pi*2*float64(i)/stepCount) * (circle.Radius - 2)) + hh

// 			x2 := (math.Sin(math.Pi*2*float64(i+1)/stepCount) * (circle.Radius - 2)) + hw
// 			y2 := (math.Cos(math.Pi*2*float64(i+1)/stepCount) * (circle.Radius - 2)) + hh

// 			ebitenutil.DrawLine(newImg, x, y, x2, y2, color.White)

// 		}
// 		circleBuffer[circle] = newImg
// 	}

// 	drawOpt := &ebiten.DrawImageOptions{}
// 	r, g, b, _ := drawColor.RGBA()
// 	drawOpt.ColorM.Scale(float64(r)/65535, float64(g)/65535, float64(b)/65535, 1)
// 	drawOpt.GeoM.Translate(circle.X-circle.Radius, circle.Y-circle.Radius)
// 	screen.DrawImage(circleBuffer[circle], drawOpt)

// }
