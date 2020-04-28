package menupage

import (

	"github.com/ahmedaly113/roadbook/ui/manager/appmanager"
	"github.com/ahmedaly113/roadbook/ui/shared/userevent"
	"github.com/ahmedaly113/roadbook/ui/shared/processevent"
	"github.com/veandco/go-sdl2/sdl"
)

type MenuPage struct {
}

func (is *MenuPage) Init() {

}

// HandleEvent handles SDL events for the menupage
func (is *MenuPage) HandleEvent(event *userevent.UserEvent) bool {

	switch event.EventType {
	case userevent.HomeBump:
		sdl.Delay(100)
		appmanager.GAppManager.SetMode(appmanager.BrightnessPage)
		processevent.OnMenuOpen()
	}

	return false
}

// Render renders the intro state
func (is *MenuPage) Render(mainWindowSurface *sdl.Surface) {

}

// WillShow is called just before this state begins
func (is *MenuPage) WillShow() {
	// call this to move on to the next transition state
	appmanager.GAppManager.WillShowComplete()
}

// WillHide is called just before this state ends
func (is *MenuPage) WillHide() {
}

// DidShow is called just after this statebegins
func (is *MenuPage) DidShow() {
	appmanager.GAppManager.SetEventMode(appmanager.AppManagerEventDriven)
}

// DidHide is called just after this state ends
func (is *MenuPage) DidHide() {
}
