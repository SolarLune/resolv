package main

import "github.com/hajimehoshi/ebiten/v2"

type WorldInterface interface {
	Init()
	Update()
	Draw(*ebiten.Image)
}
