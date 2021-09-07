package main

import (
	"image/color"
	"math"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/solarlune/resolv"
)

var circleBuffer map[resolv.Shape]*ebiten.Image = map[resolv.Shape]*ebiten.Image{}

func DrawPolygon(screen *ebiten.Image, shape *resolv.ConvexPolygon, color color.Color) {

	verts := shape.Transformed()
	for i := 0; i < len(verts); i++ {
		vert := verts[i]
		next := verts[0]

		if i < len(verts)-1 {
			next = verts[i+1]
		}
		ebitenutil.DrawLine(screen, vert.X(), vert.Y(), next.X(), next.Y(), color)

	}

}

func DrawCircle(screen *ebiten.Image, circle *resolv.Circle, drawColor color.Color) {

	// Actually drawing the circles live is too inefficient, so we will simply draw them to an image and then draw that instead
	// when necessary.

	if _, exists := circleBuffer[circle]; !exists {
		newImg := ebiten.NewImage(int(circle.Radius)*2, int(circle.Radius)*2)

		newImg.Set(int(circle.X), int(circle.Y), color.White)

		stepCount := float64(32)

		// Half image width and height.
		hw := circle.Radius
		hh := circle.Radius

		for i := 0; i < int(stepCount); i++ {

			x := (math.Sin(math.Pi*2*float64(i)/stepCount) * (circle.Radius - 2)) + hw
			y := (math.Cos(math.Pi*2*float64(i)/stepCount) * (circle.Radius - 2)) + hh

			x2 := (math.Sin(math.Pi*2*float64(i+1)/stepCount) * (circle.Radius - 2)) + hw
			y2 := (math.Cos(math.Pi*2*float64(i+1)/stepCount) * (circle.Radius - 2)) + hh

			ebitenutil.DrawLine(newImg, x, y, x2, y2, color.White)

		}
		circleBuffer[circle] = newImg
	}

	drawOpt := &ebiten.DrawImageOptions{}
	r, g, b, _ := drawColor.RGBA()
	drawOpt.ColorM.Scale(float64(r)/65535, float64(g)/65535, float64(b)/65535, 1)
	drawOpt.GeoM.Translate(circle.X-circle.Radius, circle.Y-circle.Radius)
	screen.DrawImage(circleBuffer[circle], drawOpt)

}
