package userevent

import (
	"github.com/veandco/go-sdl2/sdl"
)

const (
	None = iota
	FwdBump
	RevBump
	UpBump
	DownBump
	ModeBump
	HomeBump
)

type UserEvent struct {
	EventType int
}

func New() *UserEvent {
	return &UserEvent{ EventType: None }
}

// init UserEvent from keybard
func (ue *UserEvent) InitFromSDLEvent(event *sdl.Event) {
	ue.EventType = None
	switch event := (*event).(type) {
		case *sdl.KeyboardEvent:
		if (event.GetType() == sdl.KEYDOWN) {
			switch event.Keysym.Sym {
				case sdl.K_a:
					ue.EventType = FwdBump
				case sdl.K_z:
					ue.EventType = RevBump
				case sdl.K_UP:
					ue.EventType = UpBump
				case sdl.K_DOWN:
					ue.EventType = DownBump
				case sdl.K_RIGHT:
					ue.EventType = ModeBump
				case sdl.K_SPACE:
					ue.EventType = HomeBump
			}
		}
	}
}
