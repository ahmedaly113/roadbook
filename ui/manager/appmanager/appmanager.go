// Package appmanager tracks the current state of the app. States are
// registered in advance, and satisfy the AppMode interface. When states are
// changed, the old and new states are called with appropriate notifcations so
// they can prepare for the change.
//
// Also can execute the frame delay, and poll for SDL events in a number of ways
// (Event, EventWithTimeout, Poll).
package appmanager

import (
	"fmt"

	"github.com/ahmedaly113/roadbook/ui/shared/userevent"
	"github.com/veandco/go-sdl2/sdl"
)

// AppMode is methods for handling app events and state changes
type AppMode interface {
	Init()
	Render(*sdl.Surface)
	HandleEvent(*userevent.UserEvent) bool
	WillShow()
	DidShow()
	WillHide()
	DidHide()
}

// Event modes for eventMode
const (
	AppManagerEventDriven = iota
	AppManagerEventTimeoutDriven
	AppManagerPollDriven
)

// AppManager manages the main app states
type AppManager struct {
	currentModeID int
	nextModeID    int
	modeMap       map[int]AppMode

	FrameDelay   uint32 // ms
	EventTimeout int
	eventMode    int

	prevFrameTime uint32
}

// GAppManager is the global app manager
var GAppManager = New()

// New creates a new initialized AppManager
func New() *AppManager {
	return &AppManager{
		currentModeID: -1,
		modeMap:       make(map[int]AppMode),
		FrameDelay:    1000 / 60,
		EventTimeout:  1000 / 60,
	}
}

// RegisterMode registers a new main app mode
func (g *AppManager) RegisterMode(id int, gm AppMode) {
	g.modeMap[id] = gm
	gm.Init()
}

// SetMode sets the current app mode to the specified value
func (g *AppManager) SetMode(id int) {
	g.nextModeID = id

	if g.currentModeID >= 0 {
		g.modeMap[g.currentModeID].WillHide()
	}
	g.modeMap[id].WillShow()
}

// WillShowComplete tells the AppManager it's time to go to the next stage
// of the SetMode sequence
func (g *AppManager) WillShowComplete() {
	if g.currentModeID >= 0 {
		g.modeMap[g.currentModeID].DidHide()
	}
	g.modeMap[g.nextModeID].DidShow()

	g.currentModeID = g.nextModeID
}

// HandleEvent forwards to the event handler for the current AppMode
func (g *AppManager) HandleEvent(event *userevent.UserEvent) bool {
	return g.modeMap[g.currentModeID].HandleEvent(event)
}

// Render forwards to the renderer for the current AppMode
func (g *AppManager) Render(surface *sdl.Surface) {
	g.modeMap[g.currentModeID].Render(surface)
}

// DelayToNextFrame waits until it's time to do the next event/render loop
func (g *AppManager) DelayToNextFrame() {
	curTime := sdl.GetTicks()

	if g.prevFrameTime == 0 {
		if curTime >= g.FrameDelay {
			g.prevFrameTime = curTime - g.FrameDelay
		}
	}

	diff := curTime - g.prevFrameTime

	if g.FrameDelay > diff {
		frameDelayUnder := g.FrameDelay - diff
		// we have not yet exceeded one frame, so we need to sleep
		//fmt.Printf("Under: %d %d %d %d\n", curTime, g.prevFrameTime, diff, frameDelayUnder)
		sdl.Delay(frameDelayUnder)
	} else {
		//frameDelayOver := diff - g.FrameDelay
		//fmt.Printf("Over: %d %d %d %d\n", curTime, g.prevFrameTime, diff, frameDelayOver)
		// we have exceeded one frame, so no sleep
		// TODO sleep less in the future to make up for it?
	}

	g.prevFrameTime = curTime
}

// SetEventMode sets the event handler mode to polling or event-based
func (g *AppManager) SetEventMode(eventMode int) {
	g.eventMode = eventMode

	// Push an empty event to kick out of WaitEvent() or WaitEventTimeout()
	sdl.PushEvent(&sdl.UserEvent{Type: 0})
}

// GetNextEvent returns the next sdl.Event depending on the eventMode
func (g *AppManager) GetNextEvent() sdl.Event {
	switch g.eventMode {
	case AppManagerEventDriven:
		return sdl.WaitEvent()
	case AppManagerEventTimeoutDriven:
		return sdl.WaitEventTimeout(g.EventTimeout)
	case AppManagerPollDriven:
		return sdl.PollEvent()
	}

	panic(fmt.Sprintf("GetNextEvent: unknown event mode: %d", g.eventMode))
}
