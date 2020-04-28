package layout

import (
	"github.com/veandco/go-sdl2/sdl"
)

type Alignment int

const (
	Top Alignment = iota
	Middle
	Bottom
	Left
	Center
	Right
)

// Layout allocates contiguous full-height or full-width subrectangles
type Layout struct {
	parent *Layout
	width  int
	height int

	// for sub-layout allocation
	remainingWidth  int
	remainingHeight int

	// within parent
	x int
	y int
}

// New creates a root Layout
func New(width, height int) *Layout {
	return &Layout{
		width:           width,
		height:          height,
		remainingWidth:  width,
		remainingHeight: height,
		x:               0,
		y:               0,
	}
}

// Row allocates a full-width rectangle within the remaining width and height
func (l *Layout) Row(height int) *Layout {
	h := height
	if l.remainingHeight < height {
		h = l.remainingHeight // zero is fine
	}
	child := &Layout{
		parent:          l,
		width:           l.remainingWidth,
		height:          h,
		remainingWidth:  l.remainingWidth,
		remainingHeight: h,
		x:               l.width - l.remainingWidth,
		y:               l.height - l.remainingHeight,
	}
	// debit height
	l.remainingHeight -= h
	return child
}

// Col allocates a full-height rectangle within the remaining width and height
func (l *Layout) Col(width int) *Layout {
	w := width
	if l.remainingWidth < width {
		w = l.remainingWidth // zero is fine
	}

	child := &Layout{
		parent:          l,
		width:           w,
		height:          l.remainingHeight,
		remainingWidth:  w,
		remainingHeight: l.remainingHeight,
		x:               l.width - l.remainingWidth,
		y:               l.height - l.remainingHeight,
	}

	// debit width
	l.remainingWidth -= w
	return child
}

// R returns the current layout as an sdl.Rect
func (l *Layout) R() *sdl.Rect {
	return &sdl.Rect{
		W: int32(l.width),
		H: int32(l.height),
		X: int32(l.absX()),
		Y: int32(l.absY()),
	}
}

// Align sub-rectangle within present layout using alignment rules. Align will not
// return an error or panic on logically inconsistent rules; for example, aligning
// Left and Right will return a valid rect, however with unpredictable x alignment.
// Similarly aligning a child rectangle larger than the current layout will yield
// unpredictable results.
func (l *Layout) Align(r *sdl.Rect, alignments ...Alignment) *sdl.Rect {
	var xOff, yOff int
	for _, a := range alignments {
		switch a {
		case Top:
			yOff = 0
		case Middle:
			yOff = (l.height - int(r.H)) / 2
		case Bottom:
			yOff = l.height - int(r.H)
		case Left:
			xOff = 0
		case Center:
			xOff = (l.width - int(r.W)) / 2
		case Right:
			xOff = l.width - int(r.W)
		}
	}

	return &sdl.Rect{
		X: int32(l.absX() + xOff),
		Y: int32(l.absY() + yOff),
		W: r.W,
		H: r.H,
	}
}

// W returns with width
func (l *Layout) W() int {
	return l.width
}

// H returns with height
func (l *Layout) H() int {
	return l.height
}

func (l *Layout) absX() int {
	if l.parent == nil {
		return 0
	}
	return l.x + l.parent.absX()
}

func (l *Layout) absY() int {
	if l.parent == nil {
		return 0
	}
	return l.y + l.parent.absY()
}
