package main

import "github.com/hajimehoshi/ebiten"

type WorldInterface interface {
	Init()
	Update(*ebiten.Image)
	Draw(*ebiten.Image)
}
