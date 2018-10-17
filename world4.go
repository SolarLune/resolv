package main

import (
	"math/rand"

	"github.com/SolarLune/resolv/resolv"
)

// This world is specifically for stress-testing the CPU and collision resolution.

type World4 struct{}

var resolvers []*resolv.Rectangle

func (w World4) Create() {

	space.Clear()

	resolvers = make([]*resolv.Rectangle, 0)

	var cell int32 = 16
	space.AddShape(resolv.NewRectangle(0, 0, screenWidth, cell))
	space.AddShape(resolv.NewRectangle(0, cell, cell, screenHeight-cell))
	space.AddShape(resolv.NewRectangle(screenWidth-cell, cell, cell, screenHeight-cell))
	space.AddShape(resolv.NewRectangle(cell, screenHeight-cell, screenWidth-(cell*2), cell))

	for x := 0; x < 8000; x++ {
		rect := resolv.NewRectangle(rand.Int31n(screenWidth), rand.Int31n(screenHeight), 16, 16)
		resolvers = append(resolvers, rect)
	}

	for x := 0; x < 1500000; x++ { // Stuff to check against
		space.AddShape(resolv.NewRectangle(rand.Int31n(screenWidth), rand.Int31n(screenHeight), 16, 16))
	}

}

func (w World4) Update() {

	for _, resolver := range resolvers {
		space.Resolve(resolver, 4, 0)
		space.Resolve(resolver, 0, 4)
	}

}

func (w World4) Draw() {

	if drawHelpText {
		DrawText("There's nothing to do or see here", 0, 0)
		DrawText("because this is a pure CPU thrash-test.", 0, 16)
		DrawText("There are 8000 moving objects, resolving", 0, 32)
		DrawText("in a space full of 1.5 million shapes.", 0, 48)
	}

}

func (w World4) Destroy() {
	space.Clear()
	resolvers = make([]*resolv.Rectangle, 0)
}
