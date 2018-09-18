package main

import (
	"math/rand"

	"github.com/SolarLune/resolv/resolv"
)

// This world is specifically for stress-testing the CPU and collision resolution.

type World4 struct{}

var resolvers []*resolv.Rectangle

func (w World4) Create() {

	space = resolv.NewSpace()

	resolvers = make([]*resolv.Rectangle, 0)

	var cell int32 = 16
	space.AddShape(resolv.NewRectangle(0, 0, screenWidth, cell))
	space.AddShape(resolv.NewRectangle(0, cell, cell, screenHeight-cell))
	space.AddShape(resolv.NewRectangle(screenWidth-cell, cell, cell, screenHeight-cell))
	space.AddShape(resolv.NewRectangle(cell, screenHeight-cell, screenWidth-(cell*2), cell))

	for x := 0; x < 400; x++ {
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

	// No drawing for this to make sure that GPU isn't slowing down the CPU; this is a pure CPU thrash-test.

}
