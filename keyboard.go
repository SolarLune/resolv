package main

import (
	"github.com/veandco/go-sdl2/sdl"
)

// Keycodes
const (
	KeyNotHeld = iota
	KeyPressed
	KeyDown
	KeyReleased
)

type Keyboard struct {
	_keyDown   map[sdl.Keycode]bool
	_keyStatus map[sdl.Keycode]int
}

func NewKeyboard() Keyboard {
	kb := Keyboard{}
	kb._keyDown = make(map[sdl.Keycode]bool, 0)
	kb._keyStatus = make(map[sdl.Keycode]int, 0)
	return kb
}

func (kb *Keyboard) ReportEvent(event *sdl.KeyboardEvent) {

	if event.Type == sdl.KEYDOWN {
		kb._keyDown[event.Keysym.Sym] = true
	} else if event.Type == sdl.KEYUP {
		kb._keyDown[event.Keysym.Sym] = false
	}

}

func (kb *Keyboard) Update() {

	for code, value := range kb._keyDown {

		_, exists := kb._keyStatus[code]

		if !exists {
			kb._keyStatus[code] = 0
		}

		status := kb._keyStatus[code]

		if !value {

			if status == KeyPressed || status == KeyDown {
				kb._keyStatus[code] = KeyReleased
			} else {
				kb._keyStatus[code] = KeyNotHeld
			}

		} else {

			if status == KeyNotHeld || status == KeyReleased {
				kb._keyStatus[code] = KeyPressed
			} else {
				kb._keyStatus[code] = KeyDown
			}

		}

	}

}

func (kb *Keyboard) KeyDown(kc sdl.Keycode) bool {
	_, exists := kb._keyStatus[kc]
	if !exists {
		kb._keyStatus[kc] = 0
	}
	return kb._keyStatus[kc] == KeyDown || kb._keyStatus[kc] == KeyPressed
}

func (kb *Keyboard) KeyPressed(kc sdl.Keycode) bool {
	_, exists := kb._keyStatus[kc]
	if !exists {
		kb._keyStatus[kc] = 0
	}
	return kb._keyStatus[kc] == KeyPressed
}
func (kb *Keyboard) KeyReleased(kc sdl.Keycode) bool {
	_, exists := kb._keyStatus[kc]
	if !exists {
		kb._keyStatus[kc] = 0
	}
	return kb._keyStatus[kc] == KeyReleased
}

func (kb *Keyboard) KeyStatus(kc sdl.Keycode) int {
	_, exists := kb._keyStatus[kc]
	if !exists {
		kb._keyStatus[kc] = 0
	}
	return kb._keyStatus[kc]
}

var keyboard = NewKeyboard()
