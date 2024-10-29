module github.com/solarlune/resolv/examples

go 1.22.0

toolchain go1.23.2

require (
	github.com/golang/freetype v0.0.0-20170609003504-e2365dfdc4a0
	github.com/hajimehoshi/ebiten/v2 v2.8.3
	github.com/solarlune/resolv v0.7.0
	github.com/tanema/gween v0.0.0-20221212145351-621cc8a459d1
	golang.org/x/image v0.20.0
)

require (
	github.com/ebitengine/gomobile v0.0.0-20240911145611-4856209ac325 // indirect
	github.com/ebitengine/hideconsole v1.0.0 // indirect
	github.com/ebitengine/purego v0.8.0 // indirect
	github.com/jezek/xgb v1.1.1 // indirect
	golang.org/x/sync v0.8.0 // indirect
	golang.org/x/sys v0.25.0 // indirect
)

replace github.com/solarlune/resolv => ../
