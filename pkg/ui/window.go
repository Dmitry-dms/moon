package ui

import (
	// "fmt"
	"fmt"
	"time"

	"math/rand"

	// "math/rand"

	"github.com/Dmitry-dms/moon/pkg/gogl"
	"github.com/Dmitry-dms/moon/pkg/ui/draw"
	"github.com/Dmitry-dms/moon/pkg/ui/utils"
	"github.com/Dmitry-dms/moon/pkg/ui/widgets"
	// "github.com/go-gl/mathgl/mgl32"
)

type Window struct {
	toolbar Toolbar
	x, y    float32 // top-left corner
	w, h    float32
	active  bool
	Id      string

	outerRect  utils.Rect
	minW, minH float32

	//render
	buffer *draw.CmdBuffer

	mainWidgetSpace    *WidgetSpace
	currentWidgetSpace *WidgetSpace
	focusedWidgetSpace *WidgetSpace

	//widgSpaceWantFocus bool

	widgSpaces []*WidgetSpace

	capturedV, capturedH bool
	capturedWin          bool
	capturedInsideWin    bool

	delayedWidgets []func()
}

func genWindowId() string {
	rand.Seed(time.Now().Unix())
	return fmt.Sprint(rand.Intn(100000))
}

func NewWindow(x, y, w, h float32) *Window {
	tb := NewToolbar(x, y, w, 30)
	id := genWindowId()

	wnd := Window{
		Id:              id,
		toolbar:         tb,
		x:               x,
		y:               y,
		w:               w,
		h:               h,
		outerRect:       utils.Rect{Min: utils.Vec2{X: x, Y: y}, Max: utils.Vec2{X: x + w, Y: y + h}},
		minW:            200,
		minH:            50,
		mainWidgetSpace: newWidgetSpace(fmt.Sprintf("main-widg-space-%s", id), x, y+tb.h, w, h-tb.h, Default),
		buffer:          draw.NewBuffer(UiCtx.Io().DisplaySize),
		widgSpaces:      make([]*WidgetSpace, 0),
		delayedWidgets:  []func(){},
	}

	wnd.currentWidgetSpace = wnd.mainWidgetSpace
	return &wnd
}

const (
	defx, defy, defw, defh = 300, 100, 400, 500
	scrollChange           = 2
)

func (c *UiContext) BeginWindow(id string) {
	var wnd *Window
	var cmdw draw.Window_command
	wnd, ok := c.windowCache.Get(id)
	if !ok {
		r := rand.Intn(500)
		g := rand.Intn(300)
		wnd = NewWindow(defx+float32(r), defy+float32(g), defw, defh)
		c.Windows = append(c.Windows, wnd)
		wnd.Id = id
		c.windowCache.Add(id, wnd)
	}
	newX := wnd.x
	newY := wnd.y
	newH := wnd.h
	newW := wnd.w

	// logic
	{
		//Прямоугольник справа
		vResizeRect := utils.NewRect(wnd.x+wnd.w-scrollChange, wnd.y, scrollChange+5, wnd.h)
		hResizeRect := utils.NewRect(wnd.x, wnd.y+wnd.h-scrollChange, wnd.w, scrollChange+5)
		if utils.PointInRect(c.io.MousePos, hResizeRect) && c.ActiveWindow == wnd {
			c.io.SetCursor(VResizeCursor)
			c.wantResizeH = true
		} else if utils.PointInRect(c.io.MousePos, vResizeRect) && c.ActiveWindow == wnd {
			c.io.SetCursor(HResizeCursor)
			c.wantResizeV = true
		} else {
			c.io.SetCursor(ArrowCursor)
		}

		c.dragBehavior(vResizeRect, &wnd.capturedV)
		c.dragBehavior(hResizeRect, &wnd.capturedH)
		// Изменение размеров окна
		if c.wantResizeH && c.ActiveWindow == wnd && wnd.capturedH {
			n := newH
			n += c.io.MouseDelta.Y
			if n > wnd.minH {
				newH = n
				if wnd.mainWidgetSpace.scrlY != 0 {
					wnd.mainWidgetSpace.scrlY -= c.io.MouseDelta.Y
				}
			}
		} else if c.wantResizeV && c.ActiveWindow == wnd && wnd.capturedV {
			n := newW
			n += c.io.MouseDelta.X
			if n > wnd.minW {
				newW = n
			}
		}

		c.dragBehavior(wnd.outerRect, &wnd.capturedWin)
		// Изменение положения окна
		if c.ActiveWindow == wnd && wnd.capturedWin && !c.wantResizeV && !c.wantResizeH && !wnd.capturedInsideWin {
			newX += c.io.MouseDelta.X
			newY += c.io.MouseDelta.Y
		}
	}

	wnd.x = newX
	wnd.y = newY
	wnd.h = newH
	wnd.w = newW
	wnd.outerRect = utils.NewRect(wnd.x, wnd.y, wnd.w-wnd.mainWidgetSpace.verticalScrollbar.w, wnd.h)

	wnd.mainWidgetSpace.X = newX
	wnd.mainWidgetSpace.Y = newY + wnd.toolbar.h
	wnd.mainWidgetSpace.W = newW
	wnd.mainWidgetSpace.H = newH - wnd.toolbar.h

	{
		cmdw = draw.Window_command{
			Active: wnd.active,
			Id:     wnd.Id,
			X:      wnd.x,
			Y:      wnd.y,
			H:      wnd.h,
			W:      wnd.w,
			Clr:    mainWindowClr,
			Toolbar: draw.Toolbar_command{
				X:   newX,
				Y:   newY,
				W:   newW,
				H:   wnd.toolbar.h,
				Clr: wnd.toolbar.clr,
			},
		}
	}
	wnd.mainWidgetSpace.cursorX = wnd.mainWidgetSpace.X + UiCtx.CurrentStyle.LeftMargin
	wnd.mainWidgetSpace.cursorY = wnd.mainWidgetSpace.Y + UiCtx.CurrentStyle.TopMargin

	// Scrollbar behavior
	if wnd.mainWidgetSpace.flags&Scrollable != 0 {
		wnd.mainWidgetSpace.vertScrollBar()
		if c.ActiveWindow == wnd && c.HoveredWindow == wnd && c.io.ScrollY != 0 && c.ActiveWidgetSpaceId == wnd.mainWidgetSpace.id && c.FocusedWidgetSpace == nil {
			wnd.mainWidgetSpace.handleMouseScroll(float32(c.io.ScrollY))
		}
		if wnd.mainWidgetSpace.flags&ShowScrollbar != 0 && wnd.mainWidgetSpace.isVertScrollShown {
			cl := [4]float32{wnd.x, wnd.y, wnd.w, wnd.h}

			scrlClip := draw.NewClip(draw.EmptyClip, cl)
			scrl := wnd.mainWidgetSpace.verticalScrollbar

			wnd.buffer.CreateRect(scrl.x, scrl.y, scrl.w, scrl.h, 5, draw.AllRounded, 0, scrl.clr, scrlClip)
			wnd.buffer.CreateRect(scrl.bX, scrl.bY, scrl.bW, scrl.bH, 5, draw.AllRounded, 0, [4]float32{255, 0, 0, 1}, scrlClip)
			wnd.buffer.SeparateBuffer(0, scrlClip)
		}

	}

	wnd.mainWidgetSpace.ClipRect = [4]float32{wnd.mainWidgetSpace.X, wnd.mainWidgetSpace.Y, wnd.mainWidgetSpace.W - wnd.mainWidgetSpace.verticalScrollbar.w, wnd.mainWidgetSpace.H}

	wnd.createWindow(cmdw, draw.NewClip(draw.EmptyClip, [4]float32{wnd.x, wnd.y, wnd.w, wnd.h}))

	c.windowStack.Push(wnd)
}

var step float32 = 40

func (w *Window) createWindow(wnd draw.Window_command, clip draw.ClipRectCompose) {
	w.buffer.CreateWindow(wnd, clip)
}

func (w *Window) addWidget(widg widgets.Widget) bool {
	return w.currentWidgetSpace.addWidget(widg)
}

var (
	whiteColor    = [4]float32{255, 255, 255, 1}
	softGreen     = [4]float32{231, 240, 162, 0.8}
	black         = [4]float32{0, 0, 0, 1}
	red           = [4]float32{255, 0, 0, 1}
	transparent   = [4]float32{0, 0, 0, 0}
	mainWindowClr = [4]float32{231, 158, 162, 0.8}
)

func (wnd *Window) getWidget(id string, f func() widgets.Widget) widgets.Widget {
	var widg widgets.Widget
	widg, ok := UiCtx.GetWidget(id)
	if !ok {
		widg = f()
		wnd.addWidget(widg)
	}
	return widg
}

func (c *UiContext) ButtonT(id string, msg string) bool {
	wnd := c.windowStack.Peek()
	var tBtn *widgets.TextButton
	var hovered, clicked bool
	x, y, isRow := wnd.currentWidgetSpace.getCursorPosition()

	tBtn = wnd.getWidget(id, func() widgets.Widget {
		w, h := c.font.CalculateTextBounds(msg, c.CurrentStyle.FontScale)
		return widgets.NewTextButton(id, x, y, w, h, msg, widgets.Center, c.CurrentStyle)
	}).(*widgets.TextButton)

	// logic
	{
		hovered = c.hoverBehavior(wnd, utils.NewRectS(tBtn.BoundingBox()))
		if hovered {
			c.setActiveWidget(tBtn.Id)
			tBtn.Button.SetColor(c.CurrentStyle.BtnHoveredColor)
			if c.io.MouseClicked[0] {
				tBtn.ChangeActive()
			}
		} else if tBtn.Active() {
			tBtn.Button.SetColor(c.CurrentStyle.BtnActiveColor)
		} else {
			tBtn.Button.SetColor(c.CurrentStyle.BtnColor)
		}
		clicked = c.io.MouseClicked[0] && hovered
	}
	//

	clip := wnd.endWidget(x, y, isRow, tBtn)
	wnd.buffer.CreateButtonT(x, y, tBtn, *c.font, clip)

	return clicked
}

func (c *UiContext) Slider(id string, i *float32, min, max float32) {
	wnd := c.windowStack.Peek()
	var slider *widgets.Slider
	//var hovered, clicked bool
	x, y, isRow := wnd.currentWidgetSpace.getCursorPosition()

	slider = wnd.getWidget(id, func() widgets.Widget {
		return widgets.NewSlider(id, x, y, 100, 50, min, max, c.CurrentStyle)
	}).(*widgets.Slider)

	// logic
	{
		slider.HandleMouseDrag(c.io.MouseDelta.X, c.dragBehaviorInWindow)
		slider.CalculateNumber(i)
	}

	clip := draw.NewClip(slider.BoundingBox(), wnd.currentWidgetSpace.ClipRect)

	wnd.endWidget(x, y, isRow, slider)

	wnd.buffer.CreateRect(slider.MainSliderPos()[0], slider.MainSliderPos()[1], slider.MainSliderPos()[2], slider.MainSliderPos()[3], 0,
		draw.StraightCorners, 0, softGreen,
		clip)

	wnd.buffer.CreateRect(slider.BtnSliderPos()[0], slider.BtnSliderPos()[1], slider.BtnSliderPos()[2], slider.BtnSliderPos()[3], 0,
		draw.StraightCorners, 0, black,
		clip)
}

func (c *UiContext) Text(id string, msg string, size int) {
	wnd := c.windowStack.Peek()
	var txt *widgets.Text
	var hovered bool
	x, y, isRow := wnd.currentWidgetSpace.getCursorPosition()

	txt = wnd.getWidget(id, func() widgets.Widget {
		w, h := c.font.CalculateTextBounds(msg, c.CurrentStyle.FontScale)
		return widgets.NewText(id, msg, x, y, w, h, c.CurrentStyle)
	}).(*widgets.Text)

	if msg != txt.Message {
		txt.Message = msg
		w, h := c.font.CalculateTextBounds(msg, c.CurrentStyle.FontScale)
		txt.SetWH(w, h)
	}

	// logic
	{
		hovered = c.hoverBehavior(wnd, utils.NewRectS(txt.BoundingBox()))
		if hovered {
			c.setActiveWidget(txt.WidgetId())
			txt.SetBackGroundColor(black)
			txt.CurrentColor = [4]float32{167, 200, 100, 1}
		} else {
			txt.CurrentColor = whiteColor
			txt.SetBackGroundColor(transparent)
		}
		//if c.io.MouseClicked[0] && hovered {
		//	fmt.Println(msg)
		//}

	}
	//txt.CurrentColor = [4]float32{255, 255, 255, 1}

	clip := wnd.endWidget(x, y, isRow, txt)

	wnd.buffer.CreateText(x, y, txt, *c.font, clip)
}

func (wnd *Window) DefaultClip() draw.ClipRectCompose {
	return draw.NewClip(wnd.currentWidgetSpace.ClipRect, wnd.mainWidgetSpace.ClipRect)
}

func (c *UiContext) Image(id string, tex *gogl.Texture) bool {
	if tex == nil {
		fmt.Println("error")
		return false
	}
	wnd := c.windowStack.Peek()
	var img *widgets.Image
	var clicked bool
	x, y, isRow := wnd.currentWidgetSpace.getCursorPosition()

	img = wnd.getWidget(id, func() widgets.Widget {
		img = widgets.NewImage(id, x, y, 100, 100, whiteColor)
		return img
	}).(*widgets.Image)

	clr := img.Color()
	//y -= wnd.currentWidgetSpace.scrlY
	// logic
	{
		hovered := c.hoverBehavior(wnd, utils.NewRectS(img.BoundingBox()))
		if hovered {
			c.setActiveWidget(img.WidgetId())
			//c.Tooltip("This is important tooltip")
		}
		clicked = c.io.MouseClicked[0] && hovered
	}
	clip := wnd.endWidget(x, y, isRow, img)

	wnd.buffer.CreateRect(x, y, img.Width(), img.Height(), 0, draw.StraightCorners, tex.TextureId, clr, clip)

	return clicked
}

// TODO: measure performance
func (wnd *Window) endWidget(xPos, yPos float32, isRow bool, w widgets.Widget) draw.ClipRectCompose {
	w.UpdatePosition([4]float32{xPos, yPos, w.Width(), w.Height()})
	wnd.addCursor(w.Width(), w.Height())
	if !isRow {
		wnd.currentWidgetSpace.AddVirtualWH(w.Width(), w.Height())
	}

	wnd.debugDraw(xPos, yPos, w.Width(), w.Height())

	var clip draw.ClipRectCompose
	if wnd.currentWidgetSpace.flags&IgnoreClipping != 0 {
		clip = draw.NewClip(draw.EmptyClip, wnd.currentWidgetSpace.ClipRect)
	} else {
		clip = wnd.DefaultClip()
	}
	return clip
}

func (w *Window) addDelayedWidget(f func()) {
	w.delayedWidgets = append(w.delayedWidgets, f)
}

func (c *UiContext) VSpace(id string) {
	wnd := c.windowStack.Peek()
	var s *widgets.VSpace
	x, y, isRow := wnd.currentWidgetSpace.getCursorPosition()

	s = wnd.getWidget(id, func() widgets.Widget {
		s := widgets.NewVertSpace(id, [4]float32{x, y, 100, 20})
		return s
	}).(*widgets.VSpace)

	wnd.endWidget(x, y, isRow, s)
}

func (c *UiContext) hoverBehavior(wnd *Window, rect utils.Rect) bool {
	inRect := utils.PointInRect(c.io.MousePos, utils.NewRect(rect.Min.X, rect.Min.Y, rect.Width(), rect.Height()))
	inWindow := RegionHit(c.io.MousePos.X, c.io.MousePos.Y, wnd.x, wnd.y+wnd.toolbar.h, wnd.w, wnd.h-wnd.toolbar.h)
	focusedWidgSpace := false
	if c.FocusedWidgetSpace != nil {
		if wnd.currentWidgetSpace != c.FocusedWidgetSpace {
			focusedWidgSpace = true
		}
	}
	return c.ActiveWindow == wnd && inRect && inWindow && !focusedWidgSpace
}

func (c *UiContext) TreeNode(id string, msg string, widgFunc func()) bool {
	wnd := c.windowStack.Peek()
	var tBtn *widgets.TextButton
	var _, hovered bool
	x, y, isRow := wnd.currentWidgetSpace.getCursorPosition()

	tBtn = wnd.getWidget(id, func() widgets.Widget {
		w, h := c.font.CalculateTextBounds(msg, c.CurrentStyle.FontScale)
		return widgets.NewTextButton(id, x, y, w, h, msg, widgets.Left, c.CurrentStyle)
	}).(*widgets.TextButton)

	// logic
	{
		hovered = c.hoverBehavior(wnd, utils.NewRectS(tBtn.BoundingBox()))
		if hovered {
			c.setActiveWidget(tBtn.Id)
			tBtn.Button.SetColor(c.CurrentStyle.BtnHoveredColor)
			if c.io.MouseClicked[0] {
				tBtn.ChangeActive()
			}
		} else if tBtn.Active() {
			tBtn.Button.SetColor(c.CurrentStyle.BtnActiveColor)
		} else {
			tBtn.Button.SetColor(c.CurrentStyle.BtnColor)
		}
		//clicked = c.io.MouseClicked[0] && hovered
	}
	//
	tBtn.SetWidth(wnd.w)
	clip := wnd.endWidget(x, y, isRow, tBtn)
	wnd.buffer.CreateButtonT(x, y, tBtn, *c.font, clip)

	if tBtn.Active() {
		//wnd.currentWidgetSpace.cursorX += 50
		x += 50
		ws := c.subWidgetSpaceHelper(id, x, y+tBtn.Height(), 0, 0, NotScrollable|Resizable, widgFunc)
		wnd.currentWidgetSpace.AddVirtualHeight(ws.H)
		wnd.addCursor(ws.W, ws.H)
		//c.SubWidgetSpace(id, NotScrollable|Resizable, widgFunc)
		//wnd.currentWidgetSpace.cursorX -= 50
	}

	return tBtn.Active()
}

func (c *UiContext) Button(id string) bool {

	wnd := c.windowStack.Peek()
	var btn *widgets.Button
	var clicked, hovered bool
	x, y, isRow := wnd.currentWidgetSpace.getCursorPosition()

	btn = wnd.getWidget(id, func() widgets.Widget {
		return widgets.NewButton(id, x, y, 100, 100, c.CurrentStyle.BtnColor)
	}).(*widgets.Button)
	w := btn.Width()
	h := btn.Height()

	// logic
	{
		hovered = c.hoverBehavior(wnd, utils.NewRectS(btn.BoundingBox()))
		if hovered {
			c.setActiveWidget(btn.WidgetId())
			if c.io.MouseClicked[0] {
				btn.ChangeActive()
			}
			btn.SetColor(c.CurrentStyle.BtnHoveredColor)
		} else if btn.IsActive {
			btn.SetColor(c.CurrentStyle.BtnActiveColor)
		} else {
			btn.SetColor(c.CurrentStyle.BtnColor)
		}
		clicked = c.io.MouseClicked[0] && hovered
	}
	//
	clip := wnd.endWidget(x, y, isRow, btn)
	wnd.buffer.CreateRect(x, y, w, h, 0, draw.StraightCorners, 0, btn.Color(), clip)

	return clicked
}

func (wnd *Window) addCursor(width, height float32) {
	row, ok := wnd.currentWidgetSpace.getCurrentRow()
	if !ok {
		wnd.currentWidgetSpace.cursorY += height
	} else {
		if row.RequiereColumn {
			row.CursorY += height
			row.UpdateColWidth(width)
			row.AddColHeight(height)
		} else {
			row.CursorX += width
			row.W += width
			row.UpdateHeight(height)
		}
	}
}
func (wnd *Window) addYcursor(x float32) {
	wnd.currentWidgetSpace.cursorY += x
}

func (c *UiContext) setActiveWidget(id string) {
	c.ActiveWidget = id
}
func (c *UiContext) ContextMenu(ownerWidgetId string, widgFunc func()) {
	wnd := c.windowStack.Peek()
	var bb [4]float32
	widg, ok := c.GetWidget(ownerWidgetId)
	if !ok {
		return
	}
	bb = widg.BoundingBox()
	id := ownerWidgetId + "-ws-context"
	ws := c.getWidgetSpace(id, 0, 0, wnd, Resizable)
	if c.ActiveWidget == widg.WidgetId() {
		c.FocusedWidgetSpace = ws
	}
	if c.FocusedWidgetSpace == ws {
		f := func() {
			ws.ClipRect = [4]float32{ws.X, ws.Y, ws.W, ws.H}
			clip := draw.NewClip(ws.ClipRect, draw.EmptyClip)
			wnd.buffer.CreateRect(bb[0], bb[1]+widg.Height(), ws.W, ws.H, 0, draw.StraightCorners, 0, black, clip)
			c.subWidgetSpaceHelper(id, bb[0], bb[1]+widg.Height(), 0, 0, Resizable, widgFunc)
		}
		wnd.addDelayedWidget(f)
	}
}
func (c *UiContext) Tooltip(id string, widgFunc func()) {
	x, y := c.io.MousePos.X+10, c.io.MousePos.Y+5
	wnd := c.windowStack.Peek()

	ws := c.getWidgetSpace(id, 0, 0, wnd, Resizable|IgnoreClipping)

	wnd.addDelayedWidget(func() {
		wnd.buffer.CreateRect(x, y, ws.W, ws.H, 0, draw.StraightCorners, 0, black, draw.NewClip(draw.EmptyClip, ws.ClipRect))
		c.subWidgetSpaceHelper(id, x, y, 0, 0, Resizable|IgnoreClipping, widgFunc)
	})
}

func (c *UiContext) Column(id string, widgFunc func()) {
	wnd := c.windowStack.Peek()
	var hl *widgets.HybridLayout
	hl, ok := wnd.currentWidgetSpace.getCurrentRow()
	if !ok {
		return
	}
	hl.RequiereColumn = true
	hl.CurrentColH, hl.CurrentColW = 0, 0

	widgFunc()

	hl.RequiereColumn = false
	hl.CursorY = hl.InitY

	hl.W += hl.CurrentColW
	hl.CursorX += hl.CurrentColW
	hl.UpdateHeight(hl.CurrentColH)
}

func (wnd *Window) debugDraw(x, y, w, h float32) {
	wnd.buffer.CreateBorderBox(x, y, w, h, 2, red)
}

func (c *UiContext) getWidgetSpace(id string, width, height float32, wnd *Window, flags WidgetSpaceFlag) *WidgetSpace {
	ws, ok := c.widgSpaceCache.Get(id)
	if !ok {
		ws = newWidgetSpace(id, 0, 0, width, height, flags)
		c.widgSpaceCache.Add(id, ws)
		wnd.widgSpaces = append(wnd.widgSpaces, ws)
	}
	return ws
}

func (c *UiContext) subWidgetSpaceHelper(id string, x, y, width, height float32, flags WidgetSpaceFlag, widgFunc func()) *WidgetSpace {
	wnd := c.windowStack.Peek()

	ws := c.getWidgetSpace(id, width, height, wnd, flags)

	var prevWS = wnd.currentWidgetSpace
	wnd.currentWidgetSpace = ws

	ws.X = x
	ws.Y = y
	ws.cursorY = y
	ws.cursorX = x

	outOfWindow := false
	if y < wnd.mainWidgetSpace.Y {
		outOfWindow = true
		ws.ClipRect = [4]float32{x, wnd.mainWidgetSpace.Y, ws.W - ws.verticalScrollbar.w, ws.H - (wnd.mainWidgetSpace.Y - y)}
	} else {
		if ws.flags&ShowScrollbar != 0 {
			ws.ClipRect = [4]float32{x, y, ws.W - ws.verticalScrollbar.w, ws.H}
		} else {
			ws.ClipRect = [4]float32{x, y, ws.W, ws.H}
		}
	}

	if c.hoverBehavior(wnd, utils.NewRectS(ws.ClipRect)) {
		c.ActiveWidgetSpaceId = ws.id
	}

	// Scrollbar behavior
	if ws.flags&Scrollable != 0 {
		ws.vertScrollBar()
		if c.ActiveWindow == wnd && c.HoveredWindow == wnd && c.io.ScrollY != 0 && c.ActiveWidgetSpaceId == wnd.currentWidgetSpace.id && c.FocusedWidgetSpace == nil {
			ws.handleMouseScroll(float32(c.io.ScrollY))
		}
		if ws.flags&ShowScrollbar != 0 && ws.isVertScrollShown {
			cl := [4]float32{ws.X, ws.Y, ws.W, ws.H}
			if outOfWindow {
				cl[1] = wnd.mainWidgetSpace.Y
			}
			scrlClip := draw.NewClip(cl, wnd.mainWidgetSpace.ClipRect)

			scrl := ws.verticalScrollbar

			wnd.buffer.CreateRect(scrl.x, scrl.y, scrl.w, scrl.h, 5, draw.AllRounded, 0, scrl.clr, scrlClip)
			wnd.buffer.CreateRect(scrl.bX, scrl.bY, scrl.bW, scrl.bH, 5, draw.AllRounded, 0, [4]float32{255, 0, 0, 1}, scrlClip)

			wnd.buffer.SeparateBuffer(0, scrlClip)
		}

		//}
	}

	widgFunc()
	ws.checkVerScroll()

	//clip := draw.NewClip(ws.ClipRect, wnd.mainWidgetSpace.ClipRect)
	//
	//wnd.buffer.SeparateBuffer(0, clip)
	ws.lastVirtualHeight = ws.virtualHeight
	ws.virtualHeight = 0
	ws.lastVirtualWidth = ws.virtualWidth
	ws.virtualWidth = 0

	if ws.flags&Resizable != 0 {
		ws.H = ws.lastVirtualHeight
		ws.W = ws.lastVirtualWidth
	}

	wnd.currentWidgetSpace = prevWS
	return ws
}

func (c *UiContext) SubWidgetSpace(id string, width, height float32, flags WidgetSpaceFlag, widgFunc func()) {
	wnd := c.windowStack.Peek()
	var ws *WidgetSpace

	x, y, _ := wnd.currentWidgetSpace.getCursorPosition()
	ws = c.subWidgetSpaceHelper(id, x, y, width, height, flags, widgFunc)

	wnd.currentWidgetSpace.AddVirtualHeight(ws.H)
	wnd.addCursor(ws.W, ws.H)
}

func (c *UiContext) TabItem(name string, widgFunc func()) {
	wnd := c.windowStack.Peek()
	var tb *widgets.TabBar
	tb, ok := wnd.currentWidgetSpace.getCurrentTabBar()
	if !ok {
		return
	}
	wspId := name + "-wsp-" + tb.WidgetId()
	_, index := tb.FindTabItem(name, wspId)
	//ws := c.getWidgetSpace(wspId, 0, 0, wnd, Default)
	if index == tb.CurrentTab {
		c.SubWidgetSpace(wspId, 0, 0, Resizable|HideScrollbar, widgFunc)
	}
	//tb.SetHeight(ws.H)
	//tb.SetWidth(ws.W)
}
func (c *UiContext) TabBar(id string, widgFunc func()) {
	wnd := c.windowStack.Peek()
	var tab *widgets.TabBar
	x, y, _ := wnd.currentWidgetSpace.getCursorPosition()

	tab = wnd.getWidget(id, func() widgets.Widget {
		return widgets.NewTabBar(id, x, y, 0, 0)
	}).(*widgets.TabBar)

	tab.UpdatePosition([4]float32{x, y, wnd.w, wnd.h})

	ws := c.subWidgetSpaceHelper(id, x, y, 0, 0, Default, func() {
		c.Row("rowds", func() {
			for i, item := range tab.Bars {
				if c.ButtonT(fmt.Sprint(id, "-", i), item.Name) {
					tab.CurrentTab = i
				}
				row, _ := wnd.currentWidgetSpace.getCurrentRow()
				row.CursorX += 20
			}
		})

		wnd.currentWidgetSpace.tabStack.Push(tab)
		widgFunc()
		wnd.currentWidgetSpace.tabStack.Pop()
	})

	wnd.addCursor(ws.W, ws.H)
	wnd.currentWidgetSpace.AddVirtualHeight(ws.H)
}

func (c *UiContext) Row(id string, widgFunc func()) {
	wnd := c.windowStack.Peek()
	var row *widgets.HybridLayout
	x, y, _ := wnd.currentWidgetSpace.getCursorPosition()
	// FIXME: Is it really necessary?
	y += wnd.currentWidgetSpace.scrlY

	row = wnd.getWidget(id, func() widgets.Widget {
		return widgets.NewHLayout(id, x, y, c.CurrentStyle)
	}).(*widgets.HybridLayout)

	wnd.currentWidgetSpace.rowStack.Push(row)
	row.UpdatePosition([4]float32{x, y, wnd.w, 0})

	widgFunc()

	hl := wnd.currentWidgetSpace.rowStack.Pop()
	wnd.addCursor(0, hl.H)

	wnd.currentWidgetSpace.AddVirtualHeight(hl.H)
	hl.H = 0
	hl.W = 0
}

func (c *UiContext) EndWindow() {

	wnd := c.windowStack.Peek()

	for _, f := range wnd.delayedWidgets {
		f()
	}
	wnd.delayedWidgets = []func(){}
	wnd = c.windowStack.Pop()
	wnd.mainWidgetSpace.checkVerScroll()
	var clip = draw.ClipRectCompose{MainClipRect: wnd.mainWidgetSpace.ClipRect}
	wnd.buffer.SeparateBuffer(0, clip) // Make sure that we didn't miss anything
	wnd.mainWidgetSpace.AddVirtualHeight(c.CurrentStyle.BotMargin)

	wnd.mainWidgetSpace.lastVirtualHeight = wnd.mainWidgetSpace.virtualHeight
	wnd.mainWidgetSpace.virtualHeight = 0
}

func RegionHit(mouseX, mouseY, x, y, w, h float32) bool {
	return mouseX >= x && mouseY >= y && mouseX <= x+w && mouseY <= y+h
}
