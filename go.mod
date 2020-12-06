module solarlune.com/resolv

go 1.14

require (
	github.com/SolarLune/resolv v0.0.0-20190821203317-2f6176d8d107
	github.com/gen2brain/raylib-go v0.0.0-20201123133337-d123299701ae
	github.com/golang/freetype v0.0.0-20170609003504-e2365dfdc4a0
	github.com/hajimehoshi/ebiten v1.12.4
	github.com/kvartborg/vector v0.0.0-20200419093813-2cba0cabb4f0
	github.com/tanema/gween v0.0.0-20200417141625-072eecd4c6ed
	github.com/veandco/go-sdl2 v0.4.4 // indirect
	golang.org/x/image v0.0.0-20200801110659-972c09e46d76
)

replace github.com/SolarLune/resolv v0.0.0-20190821203317-2f6176d8d107 => ./

replace github.com/tanema/gween v0.0.0-20200417141625-072eecd4c6ed => ../gween
