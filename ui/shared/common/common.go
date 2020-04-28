package common

import (
	"github.com/veandco/go-sdl2/sdl"
)

func RenderMenuBackGround(mainWindowSurface *sdl.Surface, compRegion sdl.Rect, bgColor uint32, borderColor uint32) {
	outerRegionPlus := sdl.Rect{compRegion.X - 9, compRegion.Y - 9, compRegion.W + 18, compRegion.H + 18}
	outerRegion := sdl.Rect{compRegion.X - 7, compRegion.Y - 7, compRegion.W + 14, compRegion.H + 14}
	mainWindowSurface.FillRect(&outerRegionPlus, bgColor)
	mainWindowSurface.FillRect(&outerRegion, borderColor)
	mainWindowSurface.FillRect(&compRegion, bgColor)
} 