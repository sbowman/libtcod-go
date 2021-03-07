package tcod

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/sbowman/tcod/tcod/keys"
)

//
// Misc functions
//
func min(a, b int) int {
	if a < b {
		return a
	} else {
		return b
	}
}

func max(a, b int) int {
	if a < b {
		return b
	} else {
		return a
	}
}

func absf(v float32) float32 {
	if v < 0 {
		return -v
	}
	return v
}

func assertEqual(a interface{}, b interface{}) {
	if a != b {
		panic(fmt.Sprintf("Assertion error: a != b %v %v", a, b))
	}
}

// Inserts byte c in string a given position
// and returns new string
func insert(s string, c int, pos int) string {
	if pos < len(s) {
		return s[0:pos] + string(c) + s[pos:len(s)]
	}
	return s + string(c)
}

func replace(s string, c int, pos int) string {
	// if pos is beyond string length, then extend s
	if len(s) <= pos {
		return s + strings.Repeat(" ", pos-len(s)) + string(c)
	}
	return s[0:pos] + string(c) + s[(pos+1):len(s)]
}

// deletes string char at given position
// and returns new string
func delete(s string, pos int) string {
	if pos < len(s) {
		return s[0:pos] + s[pos+1:len(s)]
	}
	return s
}

func padRight(s string, length int, c int) string {
	if len(s) > length {
		return s[0:length]
	}
	return s + strings.Repeat(string(c), length-len(s))
}

//
// Generic widget interface
//

type IWidget interface {
	IsVisible() bool
	ComputeSize()
	Update(w IWidget, k Key)
	Render(w IWidget)
	GetX() int
	SetX(x int)
	GetY() int
	SetY(y int)
	GetWidth() int
	SetWidth(w int)
	GetHeight() int
	SetHeight(h int)
	GetUserData() interface{}
	SetUserData(data interface{})
	GetTip() string
	SetTip(tip string)
	GetMouseIn() bool
	SetMouseIn(mouseIn bool)
	GetMouseL() bool
	SetMouseL(mouseL bool)
	GetVisible() bool
	SetVisible(visible bool)
	SetGui(*Gui)
	GetGui() *Gui
	SetDefaultBackground(col, colFocus Color)
	SetDefaultForeground(col, colFocus Color)
	GetDefaultBackground() (col, colFocus Color)
	GetDefaultForeground() (col, colFocus Color)
	GetCurrentColors() (fore, back Color)
	onMouseIn()
	onMouseOut()
	onButtonPress()
	onButtonRelease()
	onButtonClick()
	expand(x, y int)
}

//
// GUI info
//

type Gui struct {
	focus         IWidget // focused widget
	keyboardFocus IWidget // keyboard focused widget
	mouse         Mouse
	elapsed       float32
	con           IConsole
	widgetVector  []IWidget
	rbs           *RadioButtonStatic
	tbs           *TextBoxStatic
}

func NewGui(console IConsole) *Gui {
	return &Gui{
		con:          console,
		widgetVector: []IWidget{},
		tbs:          NewTextBoxStatic(),
		rbs:          NewRadioButtonStatic()}
}

func (gui *Gui) Register(w IWidget) {
	w.SetGui(gui)
	gui.widgetVector = append(gui.widgetVector, w)
}

func (gui *Gui) Unregister(w IWidget) {
	if gui.focus == w {
		gui.focus = nil
	}
	if gui.keyboardFocus == w {
		gui.keyboardFocus = nil
	}
	for i, e := range gui.widgetVector {
		if e == w {
			gui.widgetVector = append(gui.widgetVector[0:i], gui.widgetVector[i+1:]...)
		}
	}
}

func (gui *Gui) updateWidgetsIntern(k Key) {
	gui.elapsed = SysGetLastFrameLength()
	for _, w := range gui.widgetVector {
		if w.IsVisible() {
			w.ComputeSize()
			w.Update(w, k)
		}
	}
}

func (gui *Gui) SetConsole(console IConsole) {
	gui.con = console
}

func (gui *Gui) UpdateWidgets(k Key) {
	gui.mouse = MouseGetStatus()
	gui.updateWidgetsIntern(k)
}

func (gui *Gui) RenderWidgets() {
	for _, w := range gui.widgetVector {
		if w.IsVisible() {
			fore, back := gui.con.GetDefaultForeground(), gui.con.GetDefaultBackground()
			w.Render(w)
			gui.con.SetDefaultForeground(fore)
			gui.con.SetDefaultBackground(back)
		}
	}
}

func (gui *Gui) IsFocused(w IWidget) bool {
	return gui.focus == w
}

func (gui *Gui) IsKeyboardFocused(w IWidget) bool {
	return gui.keyboardFocus == w
}

func (gui *Gui) GetFocusedWidget() IWidget {
	return gui.focus
}

func (gui *Gui) GetFocusedKeyboardWidget() IWidget {
	return gui.keyboardFocus
}

func (gui *Gui) UnSelectRadioGroup(group int) {
	gui.rbs.UnSelectGroup(group)
}

func (gui *Gui) SetDefaultRadioGroup(group int) {
	gui.rbs.SetDefaultGroup(group)
}

//
// Widget root class
//

type Widget struct {
	x, y, w, h int
	userData   interface{}
	tip        string
	mouseIn    bool
	mouseL     bool
	visible    bool
	back       Color
	fore       Color
	backFocus  Color
	foreFocus  Color
	gui        *Gui
}

type WidgetCallback func(w IWidget, userData interface{})

func (gui *Gui) newWidget() *Widget {
	result := &Widget{}
	gui.Register(result)
	return result
}

func (gui *Gui) NewWidget() *Widget {
	result := gui.newWidget()
	result.initializeWidget(0, 0, 0, 0)
	return result
}

//

func (gui *Gui) NewWidgetAt(x, y int) *Widget {
	result := gui.newWidget()
	result.initializeWidget(x, y, 0, 0)
	return result
}

func (gui *Gui) NewWidgetDim(x, y, w, h int) *Widget {
	result := gui.newWidget()
	result.initializeWidget(x, y, w, h)
	return result
}

// Multiple dispatch: self
func (self *Widget) initializeWidget(x, y, w, h int) {
	//assertEqual(self, iself)
	self.x = x
	self.y = y
	self.w = w
	self.h = h
	self.tip = ""
	self.mouseIn = false
	self.mouseL = false
	self.visible = true
	self.back = Color{40, 40, 120}
	self.fore = Color{220, 220, 180}
	self.backFocus = Color{70, 70, 130}
	self.foreFocus = Color{255, 255, 255}
}

func (self *Widget) Delete() {
	self.gui.Unregister(self)
}

func (self *Widget) GetGui() *Gui {
	return self.gui
}

func (self *Widget) SetGui(gui *Gui) {
	self.gui = gui
}

func (self *Widget) GetX() int {
	return self.x
}

func (self *Widget) SetX(x int) {
	self.x = x
}

func (self *Widget) GetY() int {
	return self.y
}

func (self *Widget) SetY(y int) {
	self.y = y
}

func (self *Widget) GetWidth() int {
	return self.w
}

func (self *Widget) SetWidth(w int) {
	self.w = w
}

func (self *Widget) GetHeight() int {
	return self.h
}

func (self *Widget) SetHeight(h int) {
	self.h = h
}

func (self *Widget) GetUserData() interface{} {
	return self.userData
}

func (self *Widget) SetUserData(data interface{}) {
	self.userData = data
}

func (self *Widget) GetTip() string {
	return self.tip
}

func (self *Widget) SetTip(tip string) {
	self.tip = tip
}

func (self *Widget) GetMouseIn() bool {
	return self.mouseIn
}

func (self *Widget) SetMouseIn(mouseIn bool) {
	self.mouseIn = mouseIn
}

func (self *Widget) GetMouseL() bool {
	return self.mouseL
}

func (self *Widget) SetMouseL(mouseL bool) {
	self.mouseL = mouseL
}

func (self *Widget) GetVisible() bool {
	return self.visible
}
func (self *Widget) SetVisible(visible bool) {
	self.visible = visible
}

func (self *Widget) SetDefaultBackground(col, colFocus Color) {
	self.back = col
	self.backFocus = colFocus
}

func (self *Widget) SetDefaultForeground(col, colFocus Color) {
	self.fore = col
	self.foreFocus = colFocus
}

func (self *Widget) GetDefaultBackground() (col, colFocus Color) {
	return self.back, self.backFocus
}

func (self *Widget) GetDefaultForeground() (col, colFocus Color) {
	return self.fore, self.foreFocus
}

func (self *Widget) GetCurrentColors() (fore, back Color) {
	return If(self.mouseIn, self.foreFocus, self.fore).(Color),
		If(self.mouseIn, self.backFocus, self.back).(Color)
}

// both self and iself denote the same object
// here we emulate inheritance in go
// function receives self as receiver and as first param in interface
func (self *Widget) Update(iself IWidget, k Key) {
	//assertEqual(self, iself)
	curs := MouseIsCursorVisible()
	g := self.gui

	if curs {
		if g.mouse.Cx >= iself.GetX() && g.mouse.Cx < iself.GetX()+iself.GetWidth() &&
			g.mouse.Cy >= iself.GetY() && g.mouse.Cy < iself.GetY()+iself.GetHeight() {
			if !iself.GetMouseIn() {
				iself.SetMouseIn(true)
				iself.onMouseIn()
			}
			if g.focus != iself {
				g.focus = iself
			}
		} else {
			if iself.GetMouseIn() {
				iself.SetMouseIn(false)
				iself.onMouseOut()
			}
			iself.SetMouseL(false)
			if iself == g.focus {
				g.focus = nil
			}
		}
	}
	if iself.GetMouseIn() || (!curs && iself == g.focus) {
		if g.mouse.LButton && !iself.GetMouseL() {
			iself.SetMouseL(true)
			iself.onButtonPress()
		} else if !g.mouse.LButton && iself.GetMouseL() {
			iself.onButtonRelease()
			g.keyboardFocus = nil
			if iself.GetMouseL() {
				iself.onButtonClick()
			}
			iself.SetMouseL(false)
		} else if g.mouse.LButtonPressed {
			g.keyboardFocus = nil
			iself.onButtonClick()
		}
	}
}

func (self *Widget) Move(x, y int) {
	self.x = x
	self.y = y
}

func (self *Widget) ComputeSize() {
	// abstract
}

func (self *Widget) Render(iself IWidget) {
	// abstract
}

func (self *Widget) IsVisible() bool {
	return self.visible
}

func (self *Widget) onMouseIn() {
	// abstract
}

func (self *Widget) onMouseOut() {
	// abstract
}

func (self *Widget) onButtonPress() {
	// abstract
}

func (self *Widget) onButtonRelease() {
	// abstract
}

func (self *Widget) onButtonClick() {
	// abstract
}

func (self *Widget) expand(x, y int) {
	// abstract
}

//
// Button
//

type Button struct {
	Widget
	pressed  bool
	label    string
	callback WidgetCallback
}

func (gui *Gui) newButton() *Button {
	result := &Button{}
	gui.Register(result)
	return result
}

func (gui *Gui) NewButton(label string, tip string, callback WidgetCallback, userData interface{}) *Button {
	result := gui.newButton()
	result.initializeButton(0, 0, 0, 0, label, tip, callback, userData)
	return result
}

func (gui *Gui) NewButtonDim(x, y, width, height int, label string, tip string, callback WidgetCallback, userData interface{}) *Button {
	result := gui.newButton()
	result.initializeButton(x, y, width, height, label, tip, callback, userData)
	return result
}

func (self *Button) initializeButton(x, y, width, height int, label string, tip string, callback WidgetCallback, userData interface{}) {
	self.Widget.initializeWidget(x, y, width, height)
	self.label = label
	self.tip = tip
	self.userData = userData
	self.callback = callback
	self.x = x
	self.y = y
	self.w = width
	self.h = height
}

func (self *Button) SetLabel(newLabel string) {
	self.label = newLabel
}

func (self *Button) IsPressed() bool {
	return self.pressed
}

func (self *Button) ComputeSize() {
	self.w = len(self.label) + 2
	self.h = 1
}

func (self *Button) Render(iself IWidget) {
	con := self.gui.con
	fore, back := iself.GetCurrentColors()
	con.SetDefaultForeground(fore)
	con.SetDefaultBackground(back)
	con.PrintEx(self.x+self.w/2, self.y, BkgndNone, Center, self.label)
	if self.w > 0 && self.h > 0 {
		con.Rect(self.x, self.y, self.w, self.h, true, BkgndSet)
	}
	if self.label != "" {
		if self.pressed && self.mouseIn {
			//con.PrintCenter(self.x+self.w/2, self.y, BKGND_NONE, "-%s-", self.label)
			con.PrintEx(self.x+self.w/2, self.y, BkgndNone, Center, "%s", self.label)
			//con.PrintLeft(self.x + 1, self.y, BKGND_NONE, self.label)
		} else {
			con.PrintEx(self.x+self.w/2, self.y, BkgndNone, Center, self.label)
			//con.PrintLeft(self.x + 1, self.y, BKGND_NONE, self.label)
		}
	}
}

func (self *Button) onButtonPress() {
	self.pressed = true
}

func (self *Button) onButtonRelease() {
	self.pressed = false
}

func (self *Button) onButtonClick() {
	if self.callback != nil {
		self.callback(self, self.userData)
	}
}

func (self *Button) expand(width, height int) {
	if self.w < width {
		self.w = width
	}
}

//
// Status bar
//
//
type StatusBar struct {
	Widget
}

func (gui *Gui) newStatusBar() *StatusBar {
	result := &StatusBar{}
	gui.Register(result)
	return result
}

func (gui *Gui) NewStatusBar() *StatusBar {
	result := gui.newStatusBar()
	result.initializeStatusBar(0, 0, 0, 0)
	return result
}

func (gui *Gui) NewStatusBarDim(x, y, w, h int) *StatusBar {
	result := gui.newStatusBar()
	result.initializeStatusBar(x, y, w, h)
	return result

}

func (self *StatusBar) initializeStatusBar(x, y, w, h int) {
	self.initializeWidget(x, y, w, h)
}

func (self *StatusBar) Render(iself IWidget) {
	con := self.gui.con
	focus := self.gui.focus
	con.SetDefaultBackground(self.back)
	con.Rect(self.x, self.y, self.w, self.h, true, BkgndSet)
	if focus != nil && focus.GetTip() != "" {
		con.SetDefaultForeground(self.fore)
		con.PrintRectEx(self.x+1, self.y, self.w, self.h, BkgndNone, Left, focus.GetTip())
	}
}

//
//
//
//
// Image
//
//
type ImageWidget struct {
	Widget
	back Color
}

func (gui *Gui) newImageWidget() *ImageWidget {
	result := &ImageWidget{}
	gui.Register(result)
	return result
}

func (gui *Gui) NewImageWidget(x, y, w, h int) *ImageWidget {
	result := gui.newImageWidget()
	result.initializeImageWidget(x, y, w, h, "")
	return result
}

func (gui *Gui) NewImageWidgetWithTip(x, y, w, h int, tip string) *ImageWidget {
	result := gui.newImageWidget()
	result.initializeImageWidget(x, y, w, h, tip)
	return result
}

func (self *ImageWidget) initializeImageWidget(x, y, w, h int, tip string) {
	self.Widget.initializeWidget(x, y, w, h)
	self.tip = tip
	self.back = Black
}

func (self *ImageWidget) Render(iself IWidget) {
	con := self.gui.con
	fore, back := self.GetCurrentColors()
	con.SetDefaultForeground(fore)
	con.SetDefaultBackground(back)
	con.Rect(self.x, self.y, self.w, self.h, true, BkgndSet)

}

func (self *ImageWidget) expand(width, height int) {
	if width > self.w {
		self.w = width
	}
	if height > self.h {
		self.h = height
	}
}

//
//
// Container
//
//
//

type Container struct {
	Widget
	content []IWidget
}

func (gui *Gui) newContainer() *Container {
	result := &Container{}
	gui.Register(result)
	return result
}

func (gui *Gui) NewContainer(x, y, w, h int) *Container {
	result := gui.newContainer()
	result.initializeContainer(x, y, w, h)
	return result
}

func (self *Container) initializeContainer(x, y, w, h int) {
	self.Widget.initializeWidget(x, y, w, h)
	self.content = []IWidget{}
}

func (self *Container) AddWidget(w IWidget) {
	self.content = append(self.content, w)
	self.gui.Unregister(w)
}

func (self *Container) RemoveWidget(w IWidget) {
	for i, e := range self.content {
		if e == w {
			self.content = append(self.content[:i], self.content[i+1:]...)
		}
	}
}

func (self *Container) Render(iself IWidget) {
	for _, w := range self.content {
		if w.IsVisible() {
			w.Render(w)
		}
	}
}

func (self *Container) Clear() {
	self.content = []IWidget{}
}

func (self *Container) Update(iself IWidget, k Key) {
	self.Widget.Update(iself, k)

	for _, w := range self.content {
		if w.IsVisible() {
			w.Update(w, k)
		}
	}
}

//
//
// VBox
//
type VBox struct {
	Container
	padding int
}

func (gui *Gui) newVBox() *VBox {
	result := &VBox{}
	gui.Register(result)
	return result
}

func (gui *Gui) NewVBox(x, y, padding int) *VBox {
	result := gui.newVBox()
	result.initializeVBox(x, y, padding)
	return result

}

func (self *VBox) initializeVBox(x, y, padding int) {
	self.Container.initializeContainer(x, y, 0, 0)
	self.padding = padding
}

func (self *VBox) ComputeSize() {
	cury := self.y
	self.w = 0
	for _, w := range self.content {
		if w.IsVisible() {
			w.SetX(self.x)
			w.SetY(cury)
			w.ComputeSize()
			if w.GetWidth() > self.w {
				self.w = w.GetWidth()
			}
			cury += w.GetHeight() + self.padding
		}
	}
	self.h = cury - self.y

	for _, w := range self.content {
		if w.IsVisible() {
			w.expand(self.w, w.GetHeight())
		}
	}
}

//
//
// HBox
//

type HBox struct {
	VBox
}

func (gui *Gui) newHBox() *HBox {
	result := &HBox{}
	gui.Register(result)
	return result
}

func (gui *Gui) NewHBox(x, y, padding int) *HBox {
	result := gui.newHBox()
	result.initializeHBox(x, y, padding)
	return result
}

func (self *HBox) initializeHBox(x, y, padding int) {
	self.VBox.initializeVBox(x, y, padding)
}

func (self *HBox) ComputeSize() {
	curx := self.x
	self.h = 0
	for _, w := range self.content {
		if w.IsVisible() {
			w.SetY(self.y)
			w.SetX(curx)
			w.ComputeSize()
			if w.GetHeight() > self.h {
				self.h = w.GetHeight()
			}
			curx += w.GetWidth() + self.padding
		}
	}

	self.w = curx - self.x
	for _, w := range self.content {
		if w.IsVisible() {
			w.expand(w.GetWidth(), self.h)
		}
	}
}

//
//
// Toolbar
//

//
// Separator
//
//
type Separator struct {
	Widget
	txt string
}

func (gui *Gui) newSeparator() *Separator {
	result := &Separator{}
	gui.Register(result)
	return result
}

func (gui *Gui) NewSeparator(txt string) *Separator {
	result := gui.newSeparator()
	result.initializeSeparator(txt, "")
	return result
}

func (gui *Gui) NewSeparatorWithTip(txt, tip string) *Separator {
	result := gui.newSeparator()
	result.initializeSeparator(txt, tip)
	return result

}

func (self *Separator) initializeSeparator(txt, tip string) {
	self.Widget.initializeWidget(0, 0, 0, 1)
	self.txt = txt

}

func (self *Separator) ComputeSize() {
	self.w = If(self.txt != "", len(self.txt)+2, 0).(int)
}

func (self *Separator) expand(width, height int) {
	if self.w < width {
		self.w = width
	}
}

func (self *Separator) Render(iself IWidget) {
	con := self.gui.con
	con.SetDefaultBackground(self.back)
	con.SetDefaultForeground(self.fore)
	con.Hline(self.x, self.y, self.w, BkgndSet)
	con.SetChar(self.x-1, self.y, CHAR_TEEE)
	con.SetChar(self.x+self.w, self.y, CHAR_TEEW)
	con.SetDefaultBackground(self.fore)
	con.SetDefaultForeground(self.back)
	con.PrintEx(self.x+self.w/2, self.y, BkgndSet, Center, " %s ", self.txt)
}

type ToolBar struct {
	Container
	name             string
	fixedWidth       int
	shouldPrintFrame bool
}

func (gui *Gui) newToolBar() *ToolBar {
	result := &ToolBar{}
	gui.Register(result)
	return result
}

func (gui *Gui) NewToolBarWithWidth(x, y, w int, name, tip string) *ToolBar {
	result := gui.newToolBar()
	result.initializeToolBar(x, y, w, name, tip)
	return result

}

func (gui *Gui) NewToolBar(x, y int, name, tip string) *ToolBar {
	result := gui.newToolBar()
	result.initializeToolBar(x, y, 0, name, tip)
	return result

}

func (self *ToolBar) initializeToolBar(x, y, w int, name, tip string) {
	self.Container.initializeContainer(x, y, w, 2)
	self.name = name
	self.tip = tip
	if w == 0 {
		self.w = len(name) + 4
		self.fixedWidth = 0
	} else {
		self.w = max(len(name)+4, w)
		self.fixedWidth = max(len(name)+4, w)
	}
	self.shouldPrintFrame = true

}

func (self *ToolBar) SetShouldPrintFrame(value bool) {
	self.shouldPrintFrame = value
}

func (self *ToolBar) GetShouldPrintFrame() bool {
	return self.shouldPrintFrame
}

func (self *ToolBar) Render(iself IWidget) {
	con := self.gui.con
	fore, back := iself.GetCurrentColors()
	con.SetDefaultForeground(fore)
	con.SetDefaultBackground(back)
	if self.shouldPrintFrame {
		con.PrintFrame(self.x, self.y, self.w, self.h, true, BkgndSet, self.name)
	}
	self.Container.Render(iself)
}

func (self *ToolBar) SetName(name string) {
	self.name = name
	self.fixedWidth = max(len(name)+4, self.fixedWidth)

}

func (self *ToolBar) AddSeparator(txt string) {
	self.AddWidget(self.gui.NewSeparator(txt))
}

func (self *ToolBar) AddSeparatorWithTip(txt string, tip string) {
	self.AddWidget(self.gui.NewSeparatorWithTip(txt, tip))
}

func (self *ToolBar) ComputeSize() {
	cury := self.y + 1
	self.w = If(self.name != "", len(self.name)+4, 2).(int)
	for _, w := range self.content {
		if w.IsVisible() {
			w.SetX(self.x + 1)
			w.SetY(cury)
			w.ComputeSize()
			if w.GetWidth()+2 > self.w {
				self.w = w.GetWidth() + 2
			}
			cury += w.GetHeight()
		}
	}
	if self.w < self.fixedWidth {
		self.w = self.fixedWidth
	}
	self.h = cury - self.y + 1
	for _, w := range self.content {
		if w.IsVisible() {
			w.expand(self.w-2, w.GetHeight())
		}
	}
}

//
// ToggleButton
//

type ToggleButton struct {
	Button
	pressed bool
}

func (gui *Gui) newToggleButton() *ToggleButton {
	result := &ToggleButton{}
	gui.Register(result)
	return result
}

func (gui *Gui) NewToggleButton(label, tip string, callback WidgetCallback, userData interface{}) *ToggleButton {
	result := gui.newToggleButton()
	result.initializeToggleButton(0, 0, 0, 0, label, tip, callback, userData)
	return result
}

func (gui *Gui) NewToggleButtonWithTip(x, y, width, height int, label, tip string, callback WidgetCallback, userData interface{}) *ToggleButton {
	result := gui.newToggleButton()
	result.initializeToggleButton(x, y, width, height, label, tip, callback, userData)
	return result
}

func (self *ToggleButton) initializeToggleButton(x, y, width, height int, label string, tip string, callback WidgetCallback, userData interface{}) {
	self.Button.initializeButton(x, y, width, height, label, tip, callback, userData)
}

func (self *ToggleButton) IsPressed() bool {
	return self.pressed
}

func (self *ToggleButton) SetPressed(val bool) {
	self.pressed = val
}

func (self *ToggleButton) onButtonClick() {
	self.pressed = !self.pressed
	if self.callback != nil {
		self.callback(self, self.userData)
	}
}

func (self *ToggleButton) Render(iself IWidget) {
	con := self.gui.con

	fore, back := iself.GetCurrentColors()
	con.SetDefaultBackground(back)
	con.SetDefaultForeground(fore)
	con.Rect(self.x, self.y, self.w, self.h, true, BkgndSet)
	if self.label != "" {
		con.PrintEx(self.x, self.y, BkgndNone, Left, "%c %s",
			If(self.pressed, CHAR_CHECKBOX_SET, CHAR_CHECKBOX_UNSET).(int), self.label)
	} else {
		con.PrintEx(self.x, self.y, BkgndNone, Left, "%c",
			If(self.pressed, CHAR_CHECKBOX_SET, CHAR_CHECKBOX_UNSET).(int), self.label)
	}
}

//
//
// Label
//

type Label struct {
	Widget
	label string
}

func (gui *Gui) newLabel() *Label {
	result := &Label{}
	gui.Register(result)
	return result
}

func (gui *Gui) NewLabel(x, y int, label string) *Label {
	result := gui.newLabel()
	result.initializeLabel(x, y, label, "")
	return result
}

func (gui *Gui) NewLabelWithTip(x, y int, label string, tip string) *Label {
	result := gui.newLabel()
	result.initializeLabel(x, y, label, tip)
	return result
}

func (self *Label) initializeLabel(x, y int, label string, tip string) {
	self.Widget.initializeWidget(x, y, 0, 1)
	self.x = x
	self.y = y
	self.label = label
	self.tip = tip
}

func (self *Label) Render(iself IWidget) {
	con := self.gui.con
	con.SetDefaultBackground(self.back)
	con.SetDefaultForeground(self.fore)
	con.PrintEx(self.x, self.y, BkgndNone, Left, self.label)
}

func (self *Label) ComputeSize() {
	self.w = len(self.label)
}

func (self *Label) SetValue(label string) {
	self.label = label
}

func (self *Label) expand(width, height int) {
	if self.w < width {
		self.w = width
	}
}

//
//
// TextBox
//

type TextBoxCallback func(w IWidget, val string, data interface{})

type TextBoxStatic struct {
	blinkingDelay float32
}

func NewTextBoxStatic() *TextBoxStatic {
	return &TextBoxStatic{
		blinkingDelay: 0.5}
}

type TextBox struct {
	Widget
	label            string
	txt              string
	blink            float32
	pos, offset      int
	boxx, boxw, maxw int
	insert           bool
	callback         TextBoxCallback
	data             interface{}
}

func (gui *Gui) newTextBox() *TextBox {
	result := &TextBox{}
	gui.Register(result)
	return result
}

func (gui *Gui) NewTextBox(x, y, w, maxw int, label, value string) *TextBox {
	result := gui.newTextBox()
	result.initializeTextBox(x, y, w, maxw, label, value, "")
	return result
}

func (gui *Gui) NewTextBoxWithTip(x, y, w, maxw int, label, value, tip string) *TextBox {
	result := gui.newTextBox()
	result.initializeTextBox(x, y, w, maxw, label, value, tip)
	return result
}

func (textbox *TextBox) initializeTextBox(x, y, w, maxw int, label, value, tip string) {
	textbox.Widget.initializeWidget(x, y, w, 0)
	textbox.x = x
	textbox.y = y
	textbox.w = w
	textbox.h = 1
	textbox.maxw = maxw
	textbox.label = label
	if len(value) > maxw {
		textbox.txt = value[0:maxw]
	} else {
		textbox.txt = value
	}
	textbox.tip = tip
	textbox.boxw = w
	if label != "" {
		textbox.boxx = len(label) + 1
		textbox.w += textbox.boxx
	}
}

func (textbox *TextBox) Render(iself IWidget) {
	// save colors
	con := textbox.gui.con
	g := textbox.gui

	con.SetDefaultBackground(textbox.back)
	con.SetDefaultForeground(textbox.fore)
	con.Rect(textbox.x, textbox.y, textbox.w, textbox.h, true, BkgndSet)
	if textbox.label != "" {
		con.PrintEx(textbox.x, textbox.y, BkgndNone, Left, textbox.label)
	}

	con.SetDefaultBackground(If(g.IsKeyboardFocused(textbox), textbox.foreFocus, textbox.fore).(Color))
	con.SetDefaultForeground(If(g.IsKeyboardFocused(textbox), textbox.backFocus, textbox.back).(Color))
	con.Rect(textbox.x+textbox.boxx, textbox.y, textbox.boxw, textbox.h, false, BkgndSet)
	length := len(textbox.txt) - textbox.offset
	if length > textbox.boxw {
		length = textbox.boxw
	}
	if textbox.txt != "" {
		con.PrintEx(textbox.x+textbox.boxx, textbox.y, BkgndNone, Left, padRight(textbox.txt[textbox.offset:], length, ' '))

	}
	if g.IsKeyboardFocused(textbox) && textbox.blink > 0.0 {
		if textbox.insert {
			con.SetCharBackground(textbox.x+textbox.boxx+textbox.pos-textbox.offset, textbox.y, textbox.fore, BkgndSet)
			con.SetCharForeground(textbox.x+textbox.boxx+textbox.pos-textbox.offset, textbox.y, textbox.back)
		} else {
			con.SetCharBackground(textbox.x+textbox.boxx+textbox.pos-textbox.offset, textbox.y, textbox.back, BkgndSet)
			con.SetCharForeground(textbox.x+textbox.boxx+textbox.pos-textbox.offset, textbox.y, textbox.fore)
		}
	}
}

func (textbox *TextBox) Update(iself IWidget, k Key) {
	g := textbox.gui
	tbs := textbox.gui.tbs
	if g.keyboardFocus == IWidget(textbox) {
		textbox.blink -= g.elapsed
		if textbox.blink < -tbs.blinkingDelay {
			textbox.blink += 2 * tbs.blinkingDelay
		}
		if k.VK == keys.Space || k.VK == keys.Char ||
			(k.VK >= keys.Zero && k.VK <= keys.Nine) ||
			(k.VK >= keys.KP0 && k.VK <= keys.KP9) {
			if !textbox.insert || len(textbox.txt) < textbox.maxw {
				if textbox.insert && textbox.pos < len(textbox.txt) {
					textbox.txt = insert(textbox.txt, int(k.C), textbox.pos)
				} else {
					textbox.txt = replace(textbox.txt, int(k.C), textbox.pos)
				}
				if textbox.pos < textbox.maxw {
					textbox.pos++
				}
				if textbox.pos >= textbox.boxw {
					textbox.offset = textbox.pos - textbox.boxw + 1
				}
				if textbox.callback != nil {
					textbox.callback(textbox, textbox.txt, textbox.data)
				}
			}
			textbox.blink = tbs.blinkingDelay
		}
		switch k.VK {
		case keys.Left:
			if textbox.pos > 0 {
				textbox.pos--
			}
			if textbox.pos < textbox.offset {
				textbox.offset = textbox.pos
			}
			textbox.blink = tbs.blinkingDelay
		case keys.Right:
			if textbox.pos < len(textbox.txt) {
				textbox.pos++
			}
			if textbox.pos >= textbox.boxw {
				textbox.offset = textbox.pos - textbox.boxw + 1
			}
			textbox.blink = tbs.blinkingDelay
		case keys.Home:
			textbox.pos, textbox.offset = 0, 0
			textbox.blink = tbs.blinkingDelay
		case keys.Backspace:
			if textbox.pos > 0 {
				textbox.pos--
				textbox.txt = delete(textbox.txt, textbox.pos)
				if textbox.callback != nil {
					textbox.callback(textbox, textbox.txt, textbox.data)
				}
				if textbox.pos < textbox.offset {
					textbox.offset = textbox.pos
				}
			}
			textbox.blink = tbs.blinkingDelay
		case keys.Delete:
			if textbox.pos < len(textbox.txt) {
				textbox.txt = delete(textbox.txt, textbox.pos)
				if textbox.callback != nil {
					textbox.callback(textbox, textbox.txt, textbox.data)
				}
			}
			textbox.blink = tbs.blinkingDelay
		case keys.End:
			textbox.pos = len(textbox.txt)
			if textbox.pos >= textbox.boxw {
				textbox.offset = textbox.pos - textbox.boxw + 1
			}
			textbox.blink = tbs.blinkingDelay
		default:
		}
	}
	textbox.Widget.Update(iself, k)
}

func (textbox *TextBox) SetBlinkingDelay(delay float32) {
	textbox.gui.tbs.blinkingDelay = delay
}

func (textbox *TextBox) GetBlinkingDelay() float32 {
	return textbox.gui.tbs.blinkingDelay
}

func (textbox *TextBox) SetText(txt string) {
	if textbox.maxw < len(txt) {
		textbox.txt = txt[0:textbox.maxw]
	} else {
		textbox.txt = txt
	}
}

func (textbox *TextBox) GetText() string {
	return textbox.txt
}

func (textbox *TextBox) SetCallback(callback TextBoxCallback, data interface{}) {
	textbox.callback = callback
	textbox.data = data
}

func (textbox *TextBox) onButtonClick() {
	g := textbox.gui
	if g.mouse.Cx >= textbox.x+textbox.boxx && g.mouse.Cx < textbox.x+textbox.boxx+textbox.boxw {
		g.keyboardFocus = textbox
	}
}

//
// RadioButton
//
//

type RadioButtonStatic struct {
	defaultGroup int
	groupSelect  [512]*RadioButton
	init         bool
}

func NewRadioButtonStatic() *RadioButtonStatic {
	return &RadioButtonStatic{
		defaultGroup: 0,
		init:         false,
		groupSelect:  [512]*RadioButton{}}
}

func (self *RadioButtonStatic) UnSelectGroup(group int) {
	self.groupSelect[group] = nil
}

func (self *RadioButtonStatic) SetDefaultGroup(group int) {
	self.defaultGroup = group
}

type RadioButton struct {
	Button
	foreSelection, backSelection Color
	useSelectionColor            bool
	group                        int
}

func (gui *Gui) newRadioButton() *RadioButton {
	result := &RadioButton{}
	gui.Register(result)
	return result
}

func (gui *Gui) NewRadioButton(label string, tip string, callback WidgetCallback, userData interface{}) *RadioButton {
	result := gui.newRadioButton()
	result.initializeRadioButton(0, 0, 0, 0, label, tip, callback, userData)
	return result
}

func (gui *Gui) NewRadioButtonWithTip(x, y, width, height int, label string, tip string, callback WidgetCallback, userData interface{}) *RadioButton {
	result := gui.newRadioButton()
	result.initializeRadioButton(x, y, width, height, label, tip, callback, userData)
	return result
}

func (self *RadioButton) initializeRadioButton(x, y, width, height int, label string, tip string, callback WidgetCallback, userData interface{}) {
	self.Button.initializeButton(x, y, width, height, label, tip, callback, userData)
}

func (self *RadioButton) SetGroup(group int) {
	self.group = group
}

func (self *RadioButton) SetUseSelectionColor(use bool) {
	self.useSelectionColor = use
}

func (self *RadioButton) GetUseSelectionColor() bool {
	return self.useSelectionColor
}

func (self *RadioButton) SetSelectionColor(fore, back Color) {
	self.foreSelection, self.backSelection = fore, back
}

func (self *RadioButton) GetSelectionColor() (fore, back Color) {
	return self.foreSelection, self.backSelection
}

func (self *RadioButton) GetCurrentColors() (fore, back Color) {
	fore, back = self.Button.GetCurrentColors()
	if self.useSelectionColor && self.IsSelected() {
		fore, back = self.foreSelection, self.backSelection
	}
	return
}

func (self *RadioButton) Render(iself IWidget) {
	con := self.gui.con
	fore, back := iself.GetCurrentColors()
	con.SetDefaultForeground(fore)
	con.SetDefaultBackground(back)
	self.Button.Render(iself)
	if self.IsSelected() && !self.GetUseSelectionColor() {
		con.PutCharEx(self.x, self.y, '>', fore, back)
	}
}

func (self *RadioButton) IsSelected() bool {
	rbs := self.gui.rbs
	return rbs.groupSelect[self.group] == self
}

func (self *RadioButton) Select() {
	rbs := self.gui.rbs
	rbs.groupSelect[self.group] = self
}

func (self *RadioButton) UnSelect() {
	rbs := self.gui.rbs
	rbs.groupSelect[self.group] = nil
}

func (self *RadioButton) onButtonClick() {
	self.Select()
	self.Button.onButtonClick()
}

//
//
// Slider
//

type SliderCallback func(w IWidget, val float32, data interface{})

type Slider struct {
	TextBox
	min, max     float32
	value        float32
	sensitivity  float32
	onArrows     bool
	drag         bool
	dragx, dragy int
	dragValue    float32
	fmt          string
	callback     SliderCallback
	data         interface{}
}

func (gui *Gui) newSlider() *Slider {
	result := &Slider{}
	gui.Register(result)
	return result
}

func (gui *Gui) NewSlider(x, y, w int, min, max float32, label string, tip string) *Slider {
	result := gui.newSlider()
	result.initializeSlider(x, y, w, min, max, label, tip)
	return result
}

func (self *Slider) initializeSlider(x, y, w int, min, max float32, label string, tip string) {
	self.TextBox.initializeTextBox(x, y, w, 10, label, "", tip)
	self.min = min
	self.max = max
	self.value = (min + max) * 0.5
	self.sensitivity = 1.0
	self.onArrows = false
	self.drag = false
	self.fmt = ""
	self.callback = nil
	self.data = nil
	self.valueToText()
	self.w += 2

}

func (self *Slider) GetCurrentColors() (fore, back Color) {
	fore = If(self.onArrows || self.drag, self.foreFocus, self.fore).(Color)
	back = If(self.onArrows || self.drag, self.backFocus, self.back).(Color)
	return
}

func (self *Slider) Render(iself IWidget) {
	con := self.gui.con
	fore, back := iself.GetCurrentColors()
	con.SetDefaultBackground(back)
	con.SetDefaultForeground(fore)
	self.w -= 2
	self.TextBox.Render(iself)
	self.w += 2
	con.Rect(self.x+self.w-2, self.y, 2, 1, true, BkgndSet)

	con.PutCharEx(self.x+self.w-2, self.y, CHAR_ARROW_W, fore, back)
	con.PutCharEx(self.x+self.w-1, self.y, CHAR_ARROW_E, fore, back)
}

func (self *Slider) Update(iself IWidget, k Key) {
	con := self.gui.con
	mouse := self.gui.mouse
	oldValue := self.value
	self.TextBox.Update(iself, k)
	self.textToValue()

	if mouse.Cx >= self.x+self.w-2 && mouse.Cx < self.x+self.w && mouse.Cy == self.y {
		self.onArrows = true
	} else {
		self.onArrows = false
	}
	if self.drag {
		if self.dragy == -1 {
			self.dragx = mouse.X
			self.dragy = mouse.Y
		} else {
			mdx := (float32(mouse.X-self.dragx) * self.sensitivity) / float32(con.GetWidth()*8)
			mdy := (float32(mouse.Y-self.dragy) * self.sensitivity) / float32(con.GetHeight()*8)
			oldValue := self.value
			if absf(mdy) > absf(mdx) {
				mdx = -mdy
			}
			self.value = self.dragValue + (self.max-self.min)*mdx
			self.value = ClampF(self.min, self.max, self.value)
			if self.value != oldValue {
				self.valueToText()
				self.textToValue()
			}
		}
	}
	if self.value != oldValue && self.callback != nil {
		self.callback(self, self.value, self.data)
	}
}

func (self *Slider) SetMinMax(min, max float32) {
	self.min = min
	self.max = max
}

func (self *Slider) SetCallback(callback SliderCallback, data interface{}) {
	self.callback = callback
	self.data = data
}

func (self *Slider) SetFormat(fmt string) {
	self.fmt = fmt
}

func (self *Slider) SetValue(value float32) {
	self.value = ClampF(self.min, self.max, value)
	self.valueToText()
}

func (self *Slider) SetSensitivity(sensitivity float32) {
	self.sensitivity = sensitivity
}

func (self *Slider) valueToText() {
	self.txt = fmt.Sprintf(If(self.fmt != "", self.fmt, "%.2f").(string), self.value)
}

func (self *Slider) textToValue() {
	f, err := strconv.ParseFloat(self.txt, 32)
	if err != nil {
		self.value = 0
	} else {
		self.value = float32(f)
	}
}

func (self *Slider) onButtonPress() {
	if self.onArrows {
		self.drag = true
		self.dragy = -1
		self.dragValue = self.value
		MouseShowCursor(false)
	}
}

func (self *Slider) onButtonRelease() {
	if self.drag {
		self.drag = false
		MouseMove((self.x+self.w-2)*8, self.y*8)
		MouseShowCursor(true)
	}
}
