// Package menu renders and controls a text-based menu. There are helper
// functions to control the position by keyboard or mouse, though the actual
// event handling takes place outside this package.
package menu

import (
	"fmt"

	"github.com/ahmedaly113/roadbook/ui/manager/assetmanager"
	"github.com/ahmedaly113/roadbook/ui/shared/scenegraph"
	"github.com/ahmedaly113/roadbook/ui/util"
	"github.com/veandco/go-sdl2/sdl"
)

// Menu holds information about on-screen menus
type Menu struct {
	items             []Item
	selected, clicked, checked int
	spacing           int32 // px
	justification     int
	RootEntity        *scenegraph.Entity
}

// Item describes an individual Menu line item
type Item struct {
	AssetFontID    string
	Text           string
	Color, HiColor sdl.Color
}

// Justification constants for Menu
const (
	MenuJustifyLeft = iota
	MenuJustifyCenter
	MenuJustifyRight
)

// New constructs the menu
func New(am *assetmanager.AssetManager, id string, items []Item, checkedIndex int, spacing int32, justification int, width int32) *Menu {
	menu := &Menu{items: items, spacing: spacing, justification: justification, selected: checkedIndex, checked: checkedIndex}

	menuEntity := scenegraph.NewEntity(nil)

	maxH := int32(0)
	maxW := int32(0)

	idNum := 0

	for _, item := range menu.items {
		var err error
		var surface, surfaceHi, borderSurface, checkSurface, checkSurfaceHi, bgHi *sdl.Surface
		var entity, entityHi, entityBorder, entityCheck, entityCheckHi, entityBgHi *scenegraph.Entity

		if surface, err = am.RenderText(fmt.Sprintf("%s-%d", id, idNum), item.AssetFontID, item.Text, item.Color, item.HiColor); err != nil {
			panic(fmt.Sprintf("Menu render font: %v", err))
		}

		idNum++

		if surfaceHi, err = am.RenderText(fmt.Sprintf("%s-%d", id, idNum), item.AssetFontID, item.Text, item.HiColor, item.Color); err != nil {
			panic(fmt.Sprintf("Intro render font: %v", err))
		}

		idNum++

		checkSurface = am.Surfaces["menuCheckMark"]
		if (checkSurface == nil) {
			panic("menu check asset is not available")
		}

		idNum++

		checkSurfaceHi = am.Surfaces["menuCheckMarkHi"]
		if (checkSurfaceHi == nil) {
			panic("menu check highlight asset is not available")
		}

		idNum++

		if borderSurface, err = util.MakeFillSurfaceAlpha(width, 7, 0, 0, 0, 255); err != nil {
			panic(fmt.Sprintf("Error creating border Surface: %v", err))
		}

		idNum++

		if bgHi, err = util.MakeFillSurfaceAlpha(width, menu.spacing, 0, 0, 0, 255); err != nil {
			panic(fmt.Sprintf("Error creating bgHi Surface: %v", err))
		}

		idNum++

		entity = scenegraph.NewEntity(surface)
		entityHi = scenegraph.NewEntity(surfaceHi)
		entityCheck = scenegraph.NewEntity(checkSurface)
		entityCheckHi = scenegraph.NewEntity(checkSurfaceHi)
		entityBorder = scenegraph.NewEntity(borderSurface)
		entityBgHi = scenegraph.NewEntity(bgHi)

		menuEntity.AddChild(entityBgHi, entity, entityHi, entityCheck, entityCheckHi, entityBorder)

		if entity.W > maxW {
			maxW = entity.W
		}

		entity.Y = maxH
		entityHi.Y = maxH
		entityCheck.Y = maxH
		entityCheckHi.Y = maxH
		entityBorder.Y = maxH + menu.spacing - menu.spacing/10 - menu.spacing/3
		entityBgHi.Y = maxH - menu.spacing/3

		maxH += menu.spacing
	}

	var err error
	var topBorderSurface *sdl.Surface
	var topBorderEntity *scenegraph.Entity

	if topBorderSurface, err = util.MakeFillSurfaceAlpha(width, 7, 0, 0, 0, 255); err != nil {
		panic(fmt.Sprintf("Error creating topBorderSurface Surface: %v", err))
	}

	menuEntity.W = maxW
	menuEntity.H = maxH

	topBorderEntity = scenegraph.NewEntity(topBorderSurface)
	topBorderEntity.Y = - menu.spacing/10 - menu.spacing/3
	menuEntity.AddChild(topBorderEntity)
	scenegraph.CenterEntityInParent(topBorderEntity, menuEntity)

	// position everything now that we have sizes known
	for i := range menu.items {

		eIndex := i*6

		entityBgHi := menuEntity.GetChild(eIndex)
		entity := menuEntity.GetChild(eIndex + 1)
		entityHi := menuEntity.GetChild(eIndex + 2)
		entityCheck := menuEntity.GetChild(eIndex + 3)
		entityCheckHi := menuEntity.GetChild(eIndex + 4)
		entityBorder := menuEntity.GetChild(eIndex + 5)

		switch menu.justification {
		case MenuJustifyLeft:
			entity.X = 0
			entityHi.X = 0
			entityBorder.X = 0
			entityCheck.X = entity.X - entityCheck.W - 10
		case MenuJustifyCenter:
			scenegraph.CenterEntityInParent(entity, menuEntity)
			scenegraph.CenterEntityInParent(entityHi, menuEntity)
			scenegraph.CenterEntityInParent(entityBorder, menuEntity)
			scenegraph.CenterEntityInParent(entityBgHi, menuEntity)
			entityCheck.X = entity.X - entityCheck.W - 10
			entityCheckHi.X = entity.X - entityCheckHi.W - 10
		case MenuJustifyRight:
			scenegraph.RightJustifyEntityInParent(entity, menuEntity)
			scenegraph.RightJustifyEntityInParent(entityHi, menuEntity)
			entityBorder.X = 0
			entityCheck.X = entity.X - entityCheck.W - 10
		}
	}

	menu.RootEntity = menuEntity
	menu.updateVisibility()

	return menu
}

// updateVisibility sets the visibility flags on the appropriate elements
func (m *Menu) updateVisibility() {
	menuEntity := m.RootEntity

	for i := range m.items {
		eIndex := i*6

		entityBgHi := menuEntity.GetChild(eIndex)
		entity := menuEntity.GetChild(eIndex + 1)
		entityHi := menuEntity.GetChild(eIndex + 2)
		entityCheck := menuEntity.GetChild(eIndex + 3)
		entityCheckHi := menuEntity.GetChild(eIndex + 4)

		if i == m.checked {
			if (i == m.selected) {
				entityCheckHi.Visible = true
				entityCheck.Visible = false
			} else {
				entityCheckHi.Visible = false
				entityCheck.Visible = true
			}
		} else {
			entityCheckHi.Visible = false
			entityCheck.Visible = false
		}
		if i == m.selected {
			entity.Visible = false
			entityHi.Visible = true
			entityBgHi.Visible = true
		} else {
			entity.Visible = true
			entityHi.Visible = false
			entityBgHi.Visible = false
		}

	}
}

// SetSelected sets the selected item in the menu
func (m *Menu) SetSelected(i int) {
	m.selected = i
}

// GetSelected returns the selected item in the menu
func (m *Menu) GetSelected() int {
	return m.selected
}

// GetSelected returns the selected item in the menu
func (m *Menu) SetSelectedItemAsChecked() {
	m.checked = m.selected
	m.updateVisibility()
}

// SelectNext selects the next item in the menu
func (m *Menu) SelectNext() {
	m.selected++

	if m.selected >= len(m.items) {
		m.selected = 0
	}

	m.updateVisibility()
}

// SelectPrev selects the previous item in the menu
func (m *Menu) SelectPrev() {
	m.selected--

	if m.selected < 0 {
		m.selected = len(m.items) - 1
	}

	m.updateVisibility()
}

// SelectByMouseY selects based on the Y position given
func (m *Menu) SelectByMouseY(y int32) {
	root := m.RootEntity

	y += root.WorldToEntity.Y

	for i := range m.items {
		eIndex := i * 2

		child := root.GetChild(eIndex)

		if y >= child.Y && y <= child.Y+child.H {
			m.selected = i
			m.updateVisibility()
			break
		}
	}
}

// SelectByMouseClickY selects based on the Y position given
func (m *Menu) SelectByMouseClickY(y int32) {
	root := m.RootEntity

	y += root.WorldToEntity.Y

	m.clicked = -1

	for i := range m.items {
		eIndex := i * 2

		child := root.GetChild(eIndex)

		if y >= child.Y && y <= child.Y+child.H {
			m.clicked = i
		}
	}
}

// GetClicked gets the most recently clicked menu entry
func (m *Menu) GetClicked() int {
	return m.clicked
}
