package main

import (
	"image/color"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/vector"
	"github.com/solarlune/resolv"
)

type WorldInterface interface {
	Init()
	Update()
	Draw(img *ebiten.Image)
	Space() *resolv.Space
}

func CommonDraw(screen *ebiten.Image, world WorldInterface) {

	world.Space().ForEachShape(func(shape resolv.IShape, index, maxCount int) bool {

		var drawColor color.Color = color.White

		tags := shape.Tags()

		if tags.Has(TagPlatform) && !tags.Has(TagSolidWall) {
			drawColor = color.RGBA{255, 128, 35, 255}
		}
		if tags.Has(TagPlayer) {
			drawColor = color.RGBA{32, 255, 128, 255}
		}
		if tags.Has(TagBouncer) {
			r := uint8(32)
			g := uint8(128)
			bouncer := shape.Data().(*Bouncer)
			r += uint8((255 - float64(r)) * bouncer.ColorChange)
			g += uint8((255 - float64(g)) * bouncer.ColorChange)
			drawColor = color.RGBA{r, g, 255, 255}
		}
		switch o := shape.(type) {
		case *resolv.Circle:
			vector.StrokeCircle(screen, float32(o.Position().X), float32(o.Position().Y), float32(o.Radius()), 2, drawColor, false)
		case *resolv.ConvexPolygon:

			for _, l := range o.Lines() {
				vector.StrokeLine(screen, float32(l.Start.X), float32(l.Start.Y), float32(l.End.X), float32(l.End.Y), 2, drawColor, false)
			}
		}

		return true

	})

}
