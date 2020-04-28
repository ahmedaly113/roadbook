package ui

import (
	"github.com/veandco/go-sdl2/img"
	"github.com/veandco/go-sdl2/sdl"
)

// Image functions

func loadImages() {
	buf, err := sdl.RWFromMem(literalTulip)
	if err != nil {
		panic(err)
	}

	literalTulipSurface, err = img.LoadRW(buf, false)
	if err != nil {
		panic(err)
	}
}
