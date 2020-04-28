package mainpage

import (
	"fmt"    
	"os"
    "path/filepath"
    "strings"

    "golang.org/x/text/language"
	"golang.org/x/text/message"

	_ "github.com/ahmedaly113/roadbook/catalog"

	"github.com/ahmedaly113/roadbook/ui/manager/appmanager"
	"github.com/ahmedaly113/roadbook/ui/shared/userevent"
	"github.com/ahmedaly113/roadbook/ui/shared/processevent"

	"github.com/ahmedaly113/roadbook/ui/manager/assetmanager"
	"github.com/ahmedaly113/roadbook/ui/appcontext"
	"github.com/ahmedaly113/roadbook/ui/shared/menu"
	"github.com/ahmedaly113/roadbook/ui/shared/scenegraph"	
	"github.com/ahmedaly113/roadbook/ui/shared/common"
	"github.com/veandco/go-sdl2/sdl"
)

type stateInfo struct {
	menuState     string
	selection     map[string]int
}

type MainPage struct {
	assetManager                *assetmanager.AssetManager
	rootEntity, pauseMenuEntity *scenegraph.Entity

	fontNormalColor, fontHighlightColor sdl.Color
	bgColor                             uint32
	borderColor                         uint32

	brightness int32

	compX, compY, compWidth, compHeight int32
	compRegion sdl.Rect

	menu                                *menu.Menu	

	state stateInfo
}


func (ms *stateInfo) Init() {
	ms.menuState = "1"
	ms.selection = make(map[string]int)

	ms.selection["1"] = -1
	ms.selection["1.1"] = 0
	ms.selection["1.2"] = -1
}

func (ms *stateInfo) selectMenu(selNum int) int{
	if (ms.menuState == "1") {
		switch selNum {
		case 0:
			ms.menuState = "1.1"
		case 1:
			ms.menuState = "1.2"
		case 2:			
			ms.menuState = "1.3"
		case 3:			
			ms.menuState = "1.4"
		case 4:
			return -1
		}
	} else if(ms.menuState == "1.1") {
		switch selNum {
		case 5:
			ms.menuState = "1"
		default:
			return 0
		}
	} else if(ms.menuState == "1.2") {
		switch selNum {
		case 3:
			ms.menuState = "1"
		default:
			return 0
		}
	} else if (ms.menuState == "1.3") {
		if (selNum == len(ms.getMenuStrings()) - 1) {
			ms.menuState = "1"
		} else {
			return 0;
		}
	} else if (ms.menuState == "1.4") {
		if (selNum == len(ms.getMenuStrings()) - 1) {
			ms.menuState = "1"
		} else {
			return 0;
		}
	}
	return 1
}

func (ms *stateInfo) getLanguage() language.Tag {

	switch(ms.selection["1.1"]) {
	case 0:
		return language.English
	case 1:
		return language.French
	case 2:
		return language.German
	case 3:
		return language.Spanish
	case 4:
		return language.Portuguese
	}
	return language.English
}

func (ms *stateInfo) getMenuStrings() []string {
	x := make(map[string][]string)

	x["1"] = append(x["1"], "Language")
	x["1"] = append(x["1"], "Display Prefs")
	x["1"] = append(x["1"], "Routes")
	x["1"] = append(x["1"], "Decimal Point")
	x["1"] = append(x["1"], "Exit")

	x["1.1"] = append(x["1.1"], "English")
	x["1.1"] = append(x["1.1"], "française")
	x["1.1"] = append(x["1.1"], "Deutschland")
	x["1.1"] = append(x["1.1"], "española")
	x["1.1"] = append(x["1.1"], "Português")
	x["1.1"] = append(x["1.1"], "Go Back")

	x["1.2"] = append(x["1.2"], "Prefs1")
	x["1.2"] = append(x["1.2"], "Prefs2")
	x["1.2"] = append(x["1.2"], "Prefs3")
	x["1.2"] = append(x["1.2"], "Go Back")

	var files []string

    root := os.Getenv("ROUTES_DIR") //ROUTE_DIR :"~/routes" 
    err := filepath.Walk(root, func(path string, info os.FileInfo, err error) error {

        if filepath.Ext(path) != ".gpx" {
		    return nil
		}
        files = append(files, strings.Replace(info.Name(), ".gpx", "", -1))
        return nil
    })
    if err != nil {
        panic(err)
    }
    for _, file := range files {
        x["1.3"] = append(x["1.3"], file)
    }
	x["1.3"] = append(x["1.3"], "Go Back")


	x["1.4"] = append(x["1.4"], "0,0")
	x["1.4"] = append(x["1.4"], "0,00")
	x["1.4"] = append(x["1.4"], "Go Back")

	var returnMenus []string
	p := message.NewPrinter(ms.getLanguage())
	for _, menuItemStr := range x[ms.menuState] {
		if (ms.menuState == "1.1" && menuItemStr != "Go Back") {
        	returnMenus = append(returnMenus, menuItemStr)
		} else if (ms.menuState == "1.3" && menuItemStr != "Go Back") {
        	returnMenus = append(returnMenus, menuItemStr)
		} else if (ms.menuState == "1.4" && menuItemStr != "Go Back") {
        	returnMenus = append(returnMenus, menuItemStr)
		} else {
        	returnMenus = append(returnMenus, p.Sprintf(menuItemStr))
		}
    }

	return returnMenus
}

func (ms *stateInfo) getSelected() int {
	return ms.selection[ms.menuState]
}

func (ms *stateInfo) setSelected(toSel int) {
	ms.selection[ms.menuState] = toSel
}

func (ps *MainPage) Init() {
	ps.state.Init()

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

	err := ps.assetManager.LoadJSON("json/mainpage/mainpageAsset.json")
	if err != nil {
		panic(fmt.Sprintf("json/mainpage/mainpageAsset.json: %v", err))
	}
	ps.buildScene()
}

// buildScene constructs the necessary elements for the scene
func (ps *MainPage) buildScene() {
	var err error

	am := ps.assetManager // asset manager

	ps.rootEntity, err = scenegraph.LoadJSON(am, "json/mainpage/mainpage.json", nil)
	ps.rootEntity.X = ps.compX
	ps.rootEntity.Y = ps.compY
	ps.rootEntity.W = ps.compWidth
	ps.rootEntity.H = ps.compHeight

	if err != nil {
		panic(fmt.Sprintf("json/mainpage/mainpage.json: %v", err))
	}

	ps.buildMenu()
}

func (ps *MainPage) buildMenu() {
	ps.rootEntity.EmptyChild()

	am := ps.assetManager // asset manager

	mColor := ps.fontNormalColor
	mHiColor := ps.fontHighlightColor

	menuStrings := ps.state.getMenuStrings()
	menuChecked := ps.state.getSelected()

	menuItems := []menu.Item{}
	for _, menuString := range menuStrings {
		menuItems = append(menuItems, menu.Item{ AssetFontID: "menuFont",Text: menuString,Color: mColor, HiColor: mHiColor})
	}

	ps.menu = menu.New(am, "Menu", menuItems, menuChecked, 60, menu.MenuJustifyCenter, ps.compWidth)

	scenegraph.CenterEntityInParent(ps.menu.RootEntity, ps.rootEntity)
	ps.menu.RootEntity.Y = 200

	ps.rootEntity.AddChild(ps.menu.RootEntity)
}

// HandleEvent handles SDL events for the intro state
func (ps *MainPage) HandleEvent(event *userevent.UserEvent) bool {
	switch event.EventType {
	case userevent.DownBump:
		ps.menu.SelectNext()
	case userevent.UpBump:
		ps.menu.SelectPrev()
	case userevent.ModeBump:
		selectMenuStatus := ps.state.selectMenu(ps.menu.GetSelected())
		if (selectMenuStatus == -1) {
			gm := appmanager.GAppManager
			gm.SetMode(appmanager.MenuPage)
			processevent.OnMenuClose()
		} else if (selectMenuStatus == 0) {
			ps.menu.SetSelectedItemAsChecked()
			if (ps.state.menuState == "1.1") { // Language
				processevent.SetLanguage(ps.state.getMenuStrings()[ps.menu.GetSelected()])
			} else if (ps.state.menuState == "1.4") { // Change Decimal Point
				if (ps.menu.GetSelected() == 0) {
					appcontext.GContext.NumAfterDecimal = 1
				} else {
					appcontext.GContext.NumAfterDecimal = 2
				}
			}
			ps.state.setSelected(ps.menu.GetSelected())
		} else {
			ps.buildMenu()
		}
	}

	return false
}

// Render renders the intro state
func (ps *MainPage) Render(mainWindowSurface *sdl.Surface) {
	common.RenderMenuBackGround(mainWindowSurface, ps.compRegion, ps.bgColor, ps.borderColor)

	ps.rootEntity.Render(mainWindowSurface)
}

// WillShow is called just before this state begins
func (ps *MainPage) WillShow() {
	appmanager.GAppManager.WillShowComplete()
}

// WillHide is called just before this state ends
func (ps *MainPage) WillHide() {
}

// DidShow is called just after this statebegins
func (ps *MainPage) DidShow() {
	appmanager.GAppManager.SetEventMode(appmanager.AppManagerPollDriven)
}

// DidHide is called just after this state ends
func (ps *MainPage) DidHide() {
}
