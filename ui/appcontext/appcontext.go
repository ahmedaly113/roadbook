// Package appcontext holds global, shared information about the system.
package appcontext

import "github.com/veandco/go-sdl2/sdl"

type appContext struct {
	MainWindow      *sdl.Window
	MainSurface     *sdl.Surface
	PixelFormat     *sdl.PixelFormat
	PixelFormatEnum uint32
	NumAfterDecimal	uint32

	WindowWidth, WindowHeight int32
}

// GContext holds the global app state
var GContext = &appContext{}
