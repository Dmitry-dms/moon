package ui

import (
	"github.com/Dmitry-dms/moon/pkg/fonts"
	"github.com/Dmitry-dms/moon/pkg/ui/cache"
	"github.com/Dmitry-dms/moon/pkg/ui/draw"

	"github.com/Dmitry-dms/moon/pkg/ui/styles"
	// "github.com/Dmitry-dms/moon/pkg/ui/render"
	"github.com/Dmitry-dms/moon/pkg/ui/utils"
	"github.com/Dmitry-dms/moon/pkg/ui/widgets"
)

var UiCtx *UiContext

func init() {
	UiCtx = NewContext(nil)
}

type UiContext struct {
	// rq *RenderQueue

	renderer UiRenderer
	io       *Io

	//Widgets
	Windows       []*Window
	sortedWindows []*Window
	ActiveWidget  string

	ActiveWindow *Window

	//widg space
	ActiveWidgetSpaceId              string
	WantScrollFocusWidgetSpaceId     string
	WantScrollFocusWidgetSpaceLastId string
	ActiveWidgetSpace                *WidgetSpace
	FocusedWidgetSpace               *WidgetSpace

	PriorWindow       *Window
	HoveredWindow     *Window
	LastHoveredWindow *Window

	displaySize [2]float32

	//cache
	windowCache    *cache.RamCache[*Window]
	windowStack    utils.Stack[*Window]
	widgetsCache   *cache.RamCache[widgets.Widget]
	widgSpaceCache *cache.RamCache[*WidgetSpace]

	//refactor
	Time float32

	//intent
	wantResizeH, wantResizeV bool

	//style
	CurrentStyle       *styles.Style
	prevStyle          *styles.Style
	styleChangeCounter int
	StyleChanged       bool

	//fonts
	font *fonts.Font
}

func NewContext(frontRenderer UiRenderer) *UiContext {
	c := UiContext{
		// rq:            NewRenderQueue(),
		renderer:       frontRenderer,
		io:             NewIo(),
		Windows:        make([]*Window, 0),
		sortedWindows:  make([]*Window, 0),
		windowCache:    cache.NewRamCache[*Window](),
		widgetsCache:   cache.NewRamCache[widgets.Widget](),
		windowStack:    utils.NewStack[*Window](),
		widgSpaceCache: cache.NewRamCache[*WidgetSpace](),
		CurrentStyle:   &styles.DefaultStyle,
	}

	return &c
}

func (c *UiContext) UploadFont(path string, size int) {
	c.font = fonts.NewFont(path, int32(size))
}

func (c *UiContext) Initialize(frontRenderer UiRenderer) {
	c.renderer = frontRenderer
}

func (c *UiContext) dragBehaviorInWindow(rect utils.Rect, captured *bool) {
	if !*captured {
		*captured = utils.PointInRect(c.io.MousePos, rect) && c.io.DragStarted(rect) && c.io.IsDragging
		if c.ActiveWindow != nil {
			c.ActiveWindow.capturedInsideWin = *captured
		}
	}
	if c.io.MouseReleased[0] {
		*captured = false
		if c.ActiveWindow != nil {
			c.ActiveWindow.capturedInsideWin = false
		}
	}
}

func (c *UiContext) dragBehavior(rect utils.Rect, captured *bool) {
	if !*captured {
		*captured = utils.PointInRect(c.io.MousePos, rect) && c.io.DragStarted(rect) && c.io.IsDragging
	}
	if c.io.MouseReleased[0] {
		*captured = false
	}
}

func (c *UiContext) SetScrollY(scrollY float32) {
	wnd := c.windowStack.Peek()
	wnd.currentWidgetSpace.setScrollY(scrollY)
}

func (c *UiContext) AddWidget(id string, w widgets.Widget) bool {
	return c.widgetsCache.Add(id, w)
}

func (c *UiContext) GetWidget(id string) (widgets.Widget, bool) {
	return c.widgetsCache.Get(id)
}

func (c *UiContext) Io() *Io {
	return c.io
}

func (c *UiContext) NewFrame(displaySize [2]float32) {
	c.UpdateMouseInputs()

	c.renderer.NewFrame()

	c.Io().SetDisplaySize(displaySize[0], displaySize[1])
}
func (c *UiContext) pushWindowFront(w *Window) {
	for i := len(c.sortedWindows) - 1; i >= 0; i-- {
		if c.sortedWindows[i] == w {
			if i == len(c.sortedWindows)-1 {
				return
			}
			c.sortedWindows[i] = c.sortedWindows[len(c.sortedWindows)-1]
			c.sortedWindows[len(c.sortedWindows)-1] = w
			return
		}
	}
}

func (c *UiContext) findHoveredWindow() {
	var hovered *Window
	if len(c.sortedWindows) == 0 {
		return
	}

	for i := 0; i <= len(c.sortedWindows)-1; i++ {
		window := c.sortedWindows[i]
		bb := window.outerRect

		if !bb.Contains(c.io.MousePos) {
			continue
		}
		if c.io.MouseClicked[0] && c.ActiveWindow != window {
			if !utils.PointInRect(c.io.MousePos, c.ActiveWindow.outerRect) {
				c.ActiveWindow = window
				c.pushWindowFront(window)
			}
		}

		if c.ActiveWindow == window {
			hovered = window
		} else {
			hovered = c.LastHoveredWindow
		}
		c.LastHoveredWindow = window

	}
	if c.ActiveWindow == nil {
		c.ActiveWindow = c.sortedWindows[len(c.sortedWindows)-1]
	}
	c.HoveredWindow = hovered

}

func (c *UiContext) UpdateMouseInputs() {

	io := c.Io()

	if io.IsMousePosValid(&io.MousePos) && io.IsMousePosValid(&io.MousePosPrev) {
		io.MouseDelta = io.MousePos.Sub(io.MousePosPrev)
	} else {
		io.MouseDelta = utils.Vec2{0, 0}
	}

	io.MousePosPrev = io.MousePos
	for i := 0; i < len(io.MouseDown); i++ {
		io.MouseClicked[i] = io.MouseDown[i] && io.MouseDownDuration[i] < 0
		io.MouseClickedCount[i] = 0
		io.MouseReleased[i] = !io.MouseDown[i] && io.MouseDownDuration[i] >= 0
		io.MouseDownDurationPrev[i] = io.MouseDownDuration[i]
		if io.MouseDown[i] {
			if io.MouseDownDuration[i] < 0 {
				io.MouseDownDuration[i] = 0
			} else {
				io.MouseDownDuration[i] += io.DeltaTime
			}
		} else {
			io.MouseDownDuration[i] = -1
		}

		if io.MouseClicked[i] {
			isRepeatedClick := false
			if c.Time-float32(io.MouseClickedTime[i]) < io.MouseDoubleClickTime {
				var delta utils.Vec2
				if io.IsMousePosValid(&io.MousePos) {
					delta = io.MousePos.Sub(io.MouseClickedPos[i])
				} else {
					delta = utils.Vec2{0, 0}
				}

				if delta.LengthSqr() < io.MouseDoubleClickMaxDist*io.MouseDoubleClickMaxDist {
					isRepeatedClick = true
				}
			}

			if isRepeatedClick {
				io.MouseClickedLastCount[i]++
			} else {
				io.MouseClickedLastCount[i] = 1
			}

			io.MouseClickedTime[i] = c.Time
			io.MouseClickedPos[i] = io.MousePos
			io.MouseClickedCount[i] = io.MouseClickedLastCount[i]
			io.MouseDragMaxDistanceSqr[i] = 0
		} else if io.MouseDown[i] {
			// Maintain the maximum distance we reaching from the initial click position, which is used with dragging threshold
			var deltaSqrPos float32
			if io.IsMousePosValid(&io.MousePos) {
				deltaSqrPos = (io.MousePos.Sub(io.MouseClickedPos[i])).LengthSqr()
			} else {
				deltaSqrPos = 0
			}
			io.MouseDragMaxDistanceSqr[i] = utils.Max(io.MouseDragMaxDistanceSqr[i], deltaSqrPos)
		}
		// We provide io.MouseDoubleClicked[] as a legacy service
		io.MouseDoubleClicked[i] = (io.MouseClickedCount[i] == 2)

	}

}

func copyWindows(w []*Window) []*Window {
	r := make([]*Window, len(w))
	for i, v := range w {
		r[i] = v
	}
	return r
}

var lastWinL = 0

func (c *UiContext) EndFrame(size [2]float32) {

	// Если количество окон не изменилось, в копировании нет нужды
	if lastWinL != len(c.Windows) {
		c.sortedWindows = copyWindows(c.Windows)
		lastWinL = len(c.Windows)
	}

	c.findHoveredWindow()
	if len(c.sortedWindows) == 0 {
		return
	}
	for _, v := range c.sortedWindows {
		c.renderer.Draw(size, *v.buffer)
		v.buffer.Clear()
	}

	// c.renderer.End()

	if !c.io.IsDragging && c.wantResizeH == true {
		c.wantResizeH = false

	} else if !c.io.IsDragging && c.wantResizeV == true {
		c.wantResizeV = false
	}

	c.io.ScrollX = 0
	c.io.ScrollY = 0

	c.ActiveWidget = ""

	if c.ActiveWindow != nil {
		//founded := false
		////Для правильного скроллинга
		//for _, w := range c.ActiveWindow.widgSpaces {
		//	if utils.PointInRect(c.io.MousePos, utils.NewRectS(w.ClipRect)) {
		//		c.ActiveWidgetSpaceId = w.id
		//		if w.flags&Scrollable != 0 {
		//			c.WantScrollFocusWidgetSpaceId = w.id
		//		}
		//		founded = true
		//		break
		//	}
		//}
		//if !founded {
		//	c.ActiveWidgetSpaceId = c.ActiveWindow.mainWidgetSpace.id
		//}
	} else {
		c.ActiveWidgetSpaceId = ""
	}

	// Используется для удаления фокуса, если ЛКМ была кликнута снаружи WidgetSpace
	if c.FocusedWidgetSpace != nil {
		if utils.PointOutsideRect(c.io.MouseClickedPos[0], utils.NewRectS(c.FocusedWidgetSpace.ClipRect)) {
			c.FocusedWidgetSpace = nil
		}
	}
	c.WantScrollFocusWidgetSpaceLastId = c.WantScrollFocusWidgetSpaceId
	c.WantScrollFocusWidgetSpaceId = ""

	c.io.MouseClickedPos[0] = utils.Vec2{}
	c.ActiveWindow.widgSpaces = []*WidgetSpace{}
}

type StyleVar4f uint
type StyleVar1f uint

const (
	ButtonActiveColor StyleVar4f = iota
	ButtonHoveredColor
)

const (
	Margin StyleVar1f = iota
)

func (c *UiContext) PushStyleVar1f(v StyleVar1f, val float32) {
	if c.styleChangeCounter == 0 {
		prev := *c.CurrentStyle
		c.prevStyle = &prev
	}
	switch v {
	case Margin:
		c.CurrentStyle.Margin = val
	}
	c.StyleChanged = true
	c.styleChangeCounter++
}

func (c *UiContext) PushStyleVar4f(v StyleVar4f, m [4]float32) {
	if c.styleChangeCounter == 0 {
		prev := *c.CurrentStyle
		c.prevStyle = &prev
	}
	switch v {
	case ButtonHoveredColor:
		c.CurrentStyle.BtnHoveredColor = m
	case ButtonActiveColor:
		c.CurrentStyle.BtnActiveColor = m
	}
	c.StyleChanged = true
	c.styleChangeCounter++
}
func (c *UiContext) PopStyleVar() {
	c.CurrentStyle = c.prevStyle
	c.styleChangeCounter = 0
}

type UiRenderer interface {
	NewFrame()
	Scissor(x, y, w, h int32)
	Draw(displaySize [2]float32, buffer draw.CmdBuffer)
}
