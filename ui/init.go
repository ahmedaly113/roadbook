package ui

import (
	"fmt"
	"os"

	// "github.com/ahmedaly113/roadbook/generated"
	"github.com/ahmedaly113/roadbook/generated"
	"github.com/ahmedaly113/roadbook/ui/appcontext"
	"github.com/ahmedaly113/roadbook/ui/manager/appmanager"
	"github.com/ahmedaly113/roadbook/ui/manager/assetmanager"
	"github.com/ahmedaly113/roadbook/ui/shared/userevent"

	"github.com/ahmedaly113/roadbook/ui/components/brightnesspage"
	"github.com/ahmedaly113/roadbook/ui/components/mainpage"
	"github.com/ahmedaly113/roadbook/ui/components/menupage"

	"github.com/veandco/go-sdl2/img"
	"github.com/veandco/go-sdl2/sdl"
	"github.com/veandco/go-sdl2/ttf"
)

var renderer *sdl.Renderer
var window *sdl.Window
var quit chan struct{}
var assetManager *assetmanager.AssetManager

// Done cleans up SDL structures
func Done() {
	for _, f := range cleanupFunctions {
		// ignore errors
		f()
	}
}

func sdlInit() {
	if err := sdl.Init(sdl.INIT_EVERYTHING); err != nil {
		panic(err)
	}
	cleanup(func() error {
		sdl.Quit()
		return nil
	})

	// init the image subsystem
	imgFlags := img.INIT_PNG | img.INIT_JPG

	if img.Init(imgFlags) != nil {
		panic(fmt.Sprintf("Error initializing img: %v\n", img.GetError()))
	}

	// init the TTF subsystem
	if err := ttf.Init(); err != nil {
		panic(fmt.Sprintf("Error initializing ttf: %v\n", err))
	}
}

func createMainWindow() {

	gc := appcontext.GContext
	gc.WindowWidth = screenWidth
	gc.WindowHeight = screenHeight
	gc.NumAfterDecimal = 1

	// WINDOW_SHOWN is windowed
	// WINDOW_FULLSCREEN is fullscreen
	w, err := sdl.CreateWindow("test", sdl.WINDOWPOS_UNDEFINED, sdl.WINDOWPOS_UNDEFINED,
		screenWidth, screenHeight, sdl.WINDOW_SHOWN)
	if err != nil {
		panic(err)
	}
	window = w
	cleanup(window.Destroy)

	// force software renderer
	r, err := sdl.CreateRenderer(window, -1, sdl.RENDERER_SOFTWARE)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed to create renderer: %s\n", err)
		panic(err)
	}
	renderer = r
	cleanup(renderer.Destroy)
	renderer.Clear()

	gc.MainSurface, err = window.GetSurface()
	if err != nil {
		panic(err)
	}

	gc.PixelFormatEnum, err = window.GetPixelFormat()
	if err != nil {
		panic(err)
	}

	gc.PixelFormat, err = sdl.AllocFormat(uint(gc.PixelFormatEnum)) // TODO why the cast? Seems to work?
	if err != nil {
		panic(err)
	}

	gc.MainWindow = window
}

func InitMenuPages() {

	gm := appmanager.GAppManager

	menup := &menupage.MenuPage{}
	brightnessp := &brightnesspage.BrightnessPage{}
	mainp := &mainpage.MainPage{}

	gm.RegisterMode(appmanager.MenuPage, menup)
	gm.RegisterMode(appmanager.BrightnessPage, brightnessp)
	gm.RegisterMode(appmanager.MainPage, mainp)

	gm.SetMode(appmanager.MenuPage)

	gm.SetEventMode(appmanager.AppManagerEventDriven)

}

func getFontObj(fontName string, fontSize int) *ttf.Font {
	buf, err := sdl.RWFromMem(generated.MustAsset(fontName))
	if err != nil {
		panic(err)
	}

	font, err := ttf.OpenFontRW(buf, 1, fontSize)
	if err != nil {
		panic(err)
	}
	return font
}

// Init performs all screen setup steps and returns a channel that will
// close upon receiving a quit gesture in the UI.
func Init() chan struct{} {
	sdlInit()

	// after Init
	loadImages()

	dinMedium144 = getFontObj("assets/fonts/din-medium.ttf", 144)
	dinMedium40 = getFontObj("assets/fonts/din-medium.ttf", 40)

	createMainWindow()
	InitMenuPages()

	assetManager = assetmanager.New()
	assetManager.SetOuterSurface(appcontext.GContext.MainSurface)

	err := assetManager.LoadJSON("json/app/assets.json")
	if err != nil {
		panic(fmt.Sprintf("json/app/assets.json: %v", err))
	}

	// Some UI events may result in an app quit request. Use a channel to communicate
	// quit command to caller.
	quit = make(chan struct{})
	return quit
}

// ProcessEvents is exported because SDL requires it be called from the main thread
func ProcessEvents() {
	gm := appmanager.GAppManager
	for event := sdl.PollEvent(); event != nil; event = sdl.PollEvent() {
		var usrEvt = userevent.UserEvent{}
		usrEvt.InitFromSDLEvent(&event)
		needquit := gm.HandleEvent(&usrEvt)
		if event.GetType() == sdl.QUIT || needquit {
			close(quit)
			return
		}
	}
}

var cleanupFunctions []func() error

func cleanup(f func() error) {
	// prepend for LIFO
	cleanupFunctions = append([]func() error{f}, cleanupFunctions...)
}
