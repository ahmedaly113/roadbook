package ui

import (
	"fmt"
	"image/color"
	"strings"
	"time"

	"github.com/ahmedaly113/roadbook/model"

	"github.com/ahmedaly113/roadbook/ui/appcontext"
	"github.com/ahmedaly113/roadbook/ui/layout"
	"github.com/ahmedaly113/roadbook/ui/manager/appmanager"

	"github.com/veandco/go-sdl2/img"
	"github.com/veandco/go-sdl2/sdl"
)

type surfacePair struct {
	tulip *sdl.Surface
	notes *sdl.Surface
}

// TODO: something smarter, don't need to hold all tulips and notes
// in memory simultaneously
var surfaceCache map[int]surfacePair

func RenderSpeedCell(state *model.Model, currentSurface *sdl.Surface, speedCell *layout.Layout) {
	fontColor := sdl.Color{255, 255, 255, 255}
	fontBgColor := sdl.Color{0, 0, 0, 255}
	fontBgColorSurface := sdl.MapRGB(appcontext.GContext.PixelFormat, 0, 0, 0)

	speed, err := dinMedium144.RenderUTF8Shaded(strings.Replace(fmt.Sprintf("%.0f", state.Speed), ".", ",", -1), fontColor, fontBgColor)
	if err != nil {
		panic(err)
	}
	currentSurface.FillRect(speedCell.R(), fontBgColorSurface)
	speed.Blit(nil, currentSurface, speedCell.Align(&speed.ClipRect, layout.Middle, layout.Center))


	SPRect := speedCell.R()

	speedLimit, err := dinMedium40.RenderUTF8Shaded(strings.Replace(fmt.Sprintf("%.0f", state.SpeedLimit), ".", ",", -1), fontColor, fontBgColor)
	speedLimitRect :=  speedLimit.ClipRect
	speedLimitX := SPRect.X + SPRect.W/2 - speedLimitRect.W/2
	speedLimitY := SPRect.Y + SPRect.H - 37 - speedLimitRect.H/2
	speedLimit.Blit(nil, currentSurface, &sdl.Rect{speedLimitX, speedLimitY, 0, 0})

	speedLimitBgSurface := assetManager.Surfaces["speedLimitBg"]
	speedLimitBgSurface.Blit(nil, currentSurface, &sdl.Rect{SPRect.X + SPRect.W/2 - 30, SPRect.Y + SPRect.H - 40 - 30, 0, 0})

	
	defer speed.Free()
}

func RenderHeadingCell(state *model.Model, currentSurface *sdl.Surface, headingCell *layout.Layout) {
	heading, err := dinMedium144.RenderUTF8Shaded(strings.Replace(fmt.Sprintf("%d", int(state.Heading)), ".", ",", -1), sdl.Color{0, 0, 0, 255}, sdl.Color{255, 255, 255, 255})
	if err != nil {
		panic(err)
	}

	heading.Blit(nil, currentSurface, headingCell.Align(&heading.ClipRect, layout.Middle, layout.Center))
	defer heading.Free()
}

func RenderDistanceCell(state *model.Model, currentSurface *sdl.Surface, distanceCell *layout.Layout) {
	distanceFormat := "%0.1f"
	if (appcontext.GContext.NumAfterDecimal == 2) {
		distanceFormat = "%0.2f"
	}
	distance, err := dinMedium144.RenderUTF8Shaded(strings.Replace(fmt.Sprintf(distanceFormat, state.Distance), ".", ",", -1), sdl.Color{0, 0, 0, 255}, sdl.Color{255, 255, 255, 255})
	if err != nil {
		panic(err)
	}

	distance.Blit(nil, currentSurface, distanceCell.Align(&distance.ClipRect, layout.Middle, layout.Center))
	defer distance.Free()
}

func getTulipNotesPair(cacheIdx int, waypoint *model.Waypoint) *surfacePair {
	// tulip and notes
	pair, ok := surfaceCache[cacheIdx]
	if !ok {
		pair = surfacePair{}
		buf, err := sdl.RWFromMem(waypoint.Tulip)
		if err != nil {
			panic(err)
		}

		tulip, err := img.LoadRW(buf, false)
		if err != nil {
		}
		buf.Close()

		// notes
		buf, err = sdl.RWFromMem(waypoint.Notes)
		if err != nil {
			panic(err)
		}

		notes, err := img.LoadRW(buf, false)
		if err != nil {
			panic(err)
		}
		buf.Close()

		pair.tulip = tulip
		pair.notes = notes
		if surfaceCache == nil {
			surfaceCache = make(map[int]surfacePair)
		}
		surfaceCache[cacheIdx] = pair
	}
	return &pair
}

// Render draws a single frame based on the current model state. Render
// returns the remaining duration to sleep for FPS regulation.
func Render(state *model.Model) time.Duration {
	t := time.Now()

	renderer.SetDrawColor(255, 255, 255, 255)
	renderer.FillRect(fullScreenRect)

	currentSurface, _ := window.GetSurface()

	l := layout.New(screenWidth, screenHeight)

	instruments := l.Row(instrumentHeight)
	distanceCell := instruments.Col(columnWidth)
	speedCell := instruments.Col(columnWidth)
	headingCell := instruments.Col(columnWidth)

	RenderHeadingCell(state, currentSurface, headingCell)
	RenderDistanceCell(state, currentSurface, distanceCell)
	RenderSpeedCell(state, currentSurface, speedCell)

	renderer.SetDrawColor(255, 255, 255, 255)

	// waypoints 1 - 3
	if len(state.Book) > 0 {
		//TODO: breaks if there aren't Idx + 3 remaining
		for idx, waypoint := range state.Book[state.Idx : state.Idx+3] {
			y := int32(instrumentHeight + idx*waypointHeight)
			bg := &sdl.Rect{0, y, int32(screenWidth), int32(waypointHeight)}
			r, g, b, _ := waypoint.Background.RGBA()

			// box
			renderer.SetDrawColor(uint8(r), uint8(g), uint8(b), 255)
			renderer.FillRect(bg)

			// distance
			textBg := sdl.Color{uint8(r), uint8(g), uint8(b), 255}
			textFg := sdl.Color{0, 0, 0, 255}
			if waypoint.Background == color.Black {
				textFg = sdl.Color{255, 255, 255, 255}
			}

			distanceFormat := "%0.1f"
			if (appcontext.GContext.NumAfterDecimal == 2) {
				distanceFormat = "%0.2f"
			}
			distance, _ := dinMedium144.RenderUTF8Shaded(strings.Replace(fmt.Sprintf(distanceFormat, waypoint.Distance), ".", ",", -1), textFg, textBg)

			distRect := distance.ClipRect
			distance.Blit(nil, currentSurface, &sdl.Rect{(columnWidth - distRect.W) / 2, bg.Y + (waypointHeight-distRect.H)/2, 0, 0})
			distance.Free()

			// tulip and notes
			pair := getTulipNotesPair(state.Idx+idx, &waypoint)

			pair.tulip.Blit(nil, currentSurface, &sdl.Rect{columnWidth, bg.Y + (waypointHeight-TulipHeight)/2, 0, 0})
			pair.notes.Blit(nil, currentSurface, &sdl.Rect{columnWidth * 2, bg.Y + (waypointHeight-NotesHeight)/2, 0, 0})

			// divider
			renderer.SetDrawColor(0, 0, 0, 255)
			bg.H = dividerHeight
			renderer.FillRect(bg)
		}
		//TODO: draw divider even when no waypoints exist
	}

	gc := appcontext.GContext
	gm := appmanager.GAppManager

	mainWindowSurface := gc.MainSurface
	gm.Render(mainWindowSurface)

	// display frame
	renderer.Present()

	// ensure we wait at least frameDuration
	fpsDelay := frameDuration - time.Since(t)
	return fpsDelay
}
