package brightnesspage

import (
	"fmt"

	"github.com/ahmedaly113/roadbook/ui/manager/appmanager"
	"github.com/ahmedaly113/roadbook/ui/manager/assetmanager"
	"github.com/ahmedaly113/roadbook/ui/appcontext"
	"github.com/ahmedaly113/roadbook/ui/shared/scenegraph"
	"github.com/ahmedaly113/roadbook/ui/shared/common"
	"github.com/ahmedaly113/roadbook/ui/shared/userevent"
	"github.com/ahmedaly113/roadbook/ui/shared/processevent"

	"github.com/veandco/go-sdl2/sdl"
	"github.com/veandco/go-sdl2/img"
	qrcode "github.com/skip2/go-qrcode"
)

const (
	stateInterludeDuration       = 250  // ms
	stateInterludeDurationLevel1 = 1000 // ms
)

const (
	stateInterlude = iota
	stateAction
)

type stateInfo struct {
	state     int
	startTime uint32
}

type BrightnessPage struct {
	assetManager                *assetmanager.AssetManager
	rootEntity, pauseMenuEntity *scenegraph.Entity

	fontNormalColor, fontHighlightColor sdl.Color
	bgColor                             uint32
	borderColor							uint32

	brightness int32

	compX, compY, compWidth, compHeight int32
	compRegion sdl.Rect

	scrollBallEntity    *scenegraph.Entity
	scrollBarEntity    	*scenegraph.Entity
	qrcodeEntity 		*scenegraph.Entity

	state stateInfo
}

func (ps *BrightnessPage) Init() {
	ps.initBrightness()

	// Create colors
	ps.bgColor = sdl.MapRGB(appcontext.GContext.PixelFormat, 255, 255, 255)
	ps.borderColor = sdl.MapRGB(appcontext.GContext.PixelFormat, 0, 0, 0)
	ps.fontNormalColor = sdl.Color{R: 0, G: 0, B: 0, A: 255}
	ps.fontHighlightColor = sdl.Color{R: 255, G: 255, B: 255, A: 255}

	// set Component Region
	winWidth := appcontext.GContext.WindowWidth
	winHeight := appcontext.GContext.WindowHeight

	ps.compWidth = winWidth / 2
	ps.compHeight = winHeight * 3 / 5
	ps.compX = winWidth / 4
	ps.compY = winHeight / 5
	ps.compRegion = sdl.Rect{ps.compX, ps.compX, ps.compWidth, ps.compHeight}

	ps.assetManager = assetmanager.New()
	ps.assetManager.SetOuterSurface(appcontext.GContext.MainSurface)

	err := ps.assetManager.LoadJSON("json/brightnesspage/brightnesspageAsset.json")
	if err != nil {
		panic(fmt.Sprintf("json/brightnesspage/brightnesspageAsset.json: %v", err))
	}
	ps.buildScene()
}

func (ps *BrightnessPage) initBrightness() {
	ps.brightness = 60
}

// buildScene constructs the necessary elements for the scene
func (ps *BrightnessPage) buildScene() {
	var err error

	am := ps.assetManager // asset manager

	ps.rootEntity, err = scenegraph.LoadJSON(am, "json/brightnesspage/brightnesspage.json", nil)
	ps.rootEntity.X = ps.compX
	ps.rootEntity.Y = ps.compY
	ps.rootEntity.W = ps.compWidth
	ps.rootEntity.H = ps.compHeight
	if err != nil {
		panic(fmt.Sprintf("json/brightnesspage/brightnesspage.json: %v", err))
	}

	ps.scrollBallEntity = ps.rootEntity.SearchByID("scrollBall")
	ps.scrollBarEntity = ps.rootEntity.SearchByID("brightnessScrollBar")
	ps.qrcodeEntity = ps.rootEntity.SearchByID("QRcode")

	ps.setBrightnessUI()
}

// set position of brightness scrollbar
func (ps *BrightnessPage) setBrightnessUI() {
	brightness := ps.brightness
	minX := ps.scrollBarEntity.X
	maxLen := ps.scrollBarEntity.W - ps.scrollBallEntity.W

	x := minX + (maxLen * brightness / 100)
	ps.scrollBallEntity.MoveTo(x, ps.scrollBallEntity.Y)
}

// HandleEvent handles SDL events for the intro state
func (ps *BrightnessPage) HandleEvent(event *userevent.UserEvent) bool {

	switch event.EventType {
	case userevent.DownBump:
		ps.brightness = ps.brightness - 20
		if (ps.brightness < 0){
			ps.brightness = 0
		}
		processevent.SetBrightness(ps.brightness)
		ps.setBrightnessUI()
	case userevent.UpBump:
		ps.brightness = ps.brightness + 20
		if (ps.brightness > 100){
			ps.brightness = 100
		}
		processevent.SetBrightness(ps.brightness)
		ps.setBrightnessUI()
	case userevent.HomeBump:
		sdl.Delay(100)
		appmanager.GAppManager.SetMode(appmanager.MainPage)
	}

	return false
}

// Render renders the intro state
func (ps *BrightnessPage) Render(mainWindowSurface *sdl.Surface) {
	common.RenderMenuBackGround(mainWindowSurface, ps.compRegion, ps.bgColor, ps.borderColor)

	ps.rootEntity.Render(mainWindowSurface)

	var png []byte
	png, err := qrcode.Encode("WIFI:S:Demo;T:WPA;P:password;H:false;", qrcode.Medium, 256)

	buf, err := sdl.RWFromMem(png)
	if err != nil {
		panic(err)
	}

	qrcodeImg, err := img.LoadRW(buf, false)
	qrcodeImg.Blit(nil, mainWindowSurface, &sdl.Rect{ps.qrcodeEntity.X + ps.compRegion.X - 28, ps.qrcodeEntity.Y + ps.compRegion.Y, 0, 0})
}

// WillShow is called just before this state begins
func (ps *BrightnessPage) WillShow() {
	appmanager.GAppManager.WillShowComplete()
}

// WillHide is called just before this state ends
func (ps *BrightnessPage) WillHide() {
}

// DidShow is called just after this statebegins
func (ps *BrightnessPage) DidShow() {
	appmanager.GAppManager.SetEventMode(appmanager.AppManagerPollDriven)
}

// DidHide is called just after this state ends
func (ps *BrightnessPage) DidHide() {
}
