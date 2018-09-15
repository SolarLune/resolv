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
		space.Resolve(resolver, 4, true, "")
		space.Resolve(resolver, 4, false, "")
	}

}

func (w World4) Draw() {

	// No drawing for this to make sure that GPU isn't slowing down the CPU; this is a pure CPU thrash-test.

}
