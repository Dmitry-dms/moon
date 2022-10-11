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

	widgSpaces []*WidgetSpace

	capturedV, capturedH bool
	capturedWin          bool
	capturedInsideWin    bool
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
		mainWidgetSpace: newWidgetSpace(fmt.Sprintf("main-widg-space-%s", id), x, y+tb.h, w, h-tb.h),
		buffer:          draw.NewBuffer(UiCtx.Io().DisplaySize),
		widgSpaces:      make([]*WidgetSpace, 0),
	}

	wnd.currentWidgetSpace = wnd.mainWidgetSpace
	return &wnd
}

var counter int = 0

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
		vResizeRect := utils.Rect{Min: utils.Vec2{X: wnd.x + wnd.w - scrollChange, Y: wnd.y}, Max: utils.Vec2{X: wnd.x + wnd.w + scrollChange, Y: wnd.y + wnd.h}}
		hResizeRect := utils.Rect{Min: utils.Vec2{X: wnd.x, Y: wnd.y + wnd.h - scrollChange}, Max: utils.Vec2{X: wnd.x + wnd.w, Y: wnd.y + wnd.h + scrollChange}}
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
	wnd.outerRect.Min = utils.Vec2{X: wnd.x, Y: wnd.y}
	wnd.outerRect.Max = utils.Vec2{X: wnd.x + wnd.w - wnd.mainWidgetSpace.verticalScrollbar.w, Y: wnd.y + wnd.h}

	wnd.mainWidgetSpace.X = newX
	wnd.mainWidgetSpace.Y = newY + wnd.toolbar.h
	wnd.mainWidgetSpace.W = newW
	wnd.mainWidgetSpace.H = newH - wnd.toolbar.h

	cl := [4]float32{r, g, b, 0.8}
	{
		cmdw = draw.Window_command{
			Active: wnd.active,
			Id:     wnd.Id,
			X:      wnd.x,
			Y:      wnd.y,
			H:      wnd.h,
			W:      wnd.w,
			Clr:    cl,
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
	if wnd.mainWidgetSpace.isVertScrollShown {
		wnd.mainWidgetSpace.vertScrollBar()
		if c.ActiveWindow == wnd && c.HoveredWindow == wnd && c.io.ScrollY != 0 && c.ActiveWidgetSpaceId == wnd.mainWidgetSpace.id {
			// wnd.mainWidgetSpace.handleMouseDrag()
			wnd.mainWidgetSpace.handleMouseScroll(float32(c.io.ScrollY))
		}
		cmdw.Scrollbar = draw.Scrollbar_command{
			X:         wnd.mainWidgetSpace.verticalScrollbar.x,
			Y:         wnd.mainWidgetSpace.verticalScrollbar.y,
			W:         wnd.mainWidgetSpace.verticalScrollbar.w,
			H:         wnd.mainWidgetSpace.verticalScrollbar.h,
			Xb:        wnd.mainWidgetSpace.verticalScrollbar.bX,
			Yb:        wnd.mainWidgetSpace.verticalScrollbar.bY,
			Wb:        wnd.mainWidgetSpace.verticalScrollbar.bW,
			Hb:        wnd.mainWidgetSpace.verticalScrollbar.bH,
			Radius:    5,
			ScrollClr: wnd.mainWidgetSpace.verticalScrollbar.clr,
			BtnClr:    [4]float32{255, 0, 0, 1},
		}
	}

	wnd.mainWidgetSpace.ClipRect = [4]float32{wnd.mainWidgetSpace.X, wnd.mainWidgetSpace.Y, wnd.mainWidgetSpace.W - wnd.mainWidgetSpace.verticalScrollbar.w, wnd.mainWidgetSpace.H}
	//wnd.buffer.InnerWindowSpace = [4]float32{wnd.mainWidgetSpace.X, wnd.mainWidgetSpace.Y, wnd.mainWidgetSpace.W - wnd.mainWidgetSpace.verticalScrollbar.w, wnd.mainWidgetSpace.H}
	wnd.createWindow(cmdw, draw.ClipRectCompose{
		//MainClipRect: wnd.mainWidgetSpace.ClipRect,
		MainClipRect: [4]float32{wnd.x, wnd.y, wnd.w, wnd.h},
	})
	c.windowStack.Push(wnd)
}

var step float32 = 40
var r, g, b float32 = 231, 158, 162

func (w *Window) createWindow(wnd draw.Window_command, clip draw.ClipRectCompose) {
	w.buffer.CreateWindow(wnd, clip)
}

func (w *Window) addWidget(widg widgets.Widget) bool {
	return w.currentWidgetSpace.addWidget(widg)
}

var (
	whiteColor = [4]float32{255, 255, 255, 1}
	softGreen  = [4]float32{231, 240, 162, 0.8}
	black      = [4]float32{0, 0, 0, 1}
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

var scale float32 = 1

func (c *UiContext) ButtonT(id string, msg string) bool {
	wnd := c.windowStack.Peek()
	var tBtn *widgets.TextButton
	var hovered, clicked bool
	x, y, isRow := wnd.currentWidgetSpace.getCursorPosition()

	tBtn = wnd.getWidget(id, func() widgets.Widget {
		s := c.font.CalculateTextBounds(msg, c.CurrentStyle.TextSize)
		return widgets.NewTextButton(id, x, y, s, msg, widgets.Center, c.CurrentStyle)
	}).(*widgets.TextButton)

	y -= wnd.currentWidgetSpace.scrlY
	// logic
	{
		hovered = c.hoverBehavior(wnd, utils.NewRectS(tBtn.Box()))
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
	wnd.buffer.CreateButtonT(x, y, tBtn, *c.font, draw.ClipRectCompose{MainClipRect: wnd.mainWidgetSpace.ClipRect, ClipRect: wnd.currentWidgetSpace.ClipRect})
	//
	tBtn.UpdatePosition([4]float32{x, y, tBtn.Width(), tBtn.Height()})

	wnd.addCursor(tBtn.Width(), tBtn.Height())
	if !isRow {
		wnd.currentWidgetSpace.AddVirtualHeight(tBtn.Height())
	}
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

	y -= wnd.currentWidgetSpace.scrlY
	// logic
	{
		slider.HandleMouseDrag(c.io.MouseDelta.X, c.dragBehaviorInWindow)
		slider.CalculateNumber(i)

		//hovered = c.hoverBehavior(wnd, utils.NewRectS(tBtn.Box()))
		//if hovered {
		//	c.setActiveWidget(tBtn.Id)
		//	tBtn.Button.SetColor(c.CurrentStyle.BtnHoveredColor)
		//	if c.io.MouseClicked[0] {
		//		tBtn.ChangeActive()
		//	}
		//} else if tBtn.Active() {
		//	tBtn.Button.SetColor(c.CurrentStyle.BtnActiveColor)
		//} else {
		//	tBtn.Button.SetColor(c.CurrentStyle.BtnColor)
		//}
		//clicked = c.io.MouseClicked[0] && hovered
	}
	slider.UpdatePosition([4]float32{x, y, slider.Width(), slider.Height()})

	wnd.buffer.CreateRect(slider.MainSliderPos()[0], slider.MainSliderPos()[1], slider.MainSliderPos()[2], slider.MainSliderPos()[3], 0,
		draw.StraightCorners, 0, softGreen,
		draw.ClipRectCompose{MainClipRect: slider.BoundingBox(), ClipRect: wnd.currentWidgetSpace.ClipRect})

	wnd.buffer.CreateRect(slider.BtnSliderPos()[0], slider.BtnSliderPos()[1], slider.BtnSliderPos()[2], slider.BtnSliderPos()[3], 0,
		draw.StraightCorners, 0, black,
		draw.ClipRectCompose{MainClipRect: slider.BoundingBox(), ClipRect: wnd.currentWidgetSpace.ClipRect})

	wnd.addCursor(slider.Width(), slider.Height())
	if !isRow {
		wnd.currentWidgetSpace.AddVirtualHeight(slider.Height())
	}
}

func (c *UiContext) Text(id string, msg string, size int) {
	wnd := c.windowStack.Peek()
	var txt *widgets.Text
	var hovered bool
	x, y, isRow := wnd.currentWidgetSpace.getCursorPosition()

	txt = wnd.getWidget(id, func() widgets.Widget {
		s := c.font.CalculateTextBounds(msg, size)
		return widgets.NewText(id, msg, x, y, s[0], s[1], size, c.CurrentStyle.TextColor)
	}).(*widgets.Text)

	y -= wnd.currentWidgetSpace.scrlY
	// logic
	{
		hovered = c.hoverBehavior(wnd, utils.NewRectS(txt.BoundingBox))
		if hovered {
			c.setActiveWidget(txt.Id)
			txt.CurrentColor = [4]float32{167, 200, 100, 1}
		} else {
			txt.CurrentColor = whiteColor
		}
	}

	txt.UpdatePosition([4]float32{x, y, txt.Width(), txt.Height()})
	wnd.buffer.CreateText(x, y, txt, *c.font, draw.ClipRectCompose{MainClipRect: wnd.mainWidgetSpace.ClipRect, ClipRect: wnd.currentWidgetSpace.ClipRect})

	wnd.addCursor(txt.Width(), txt.Height())
	if !isRow {
		wnd.currentWidgetSpace.AddVirtualHeight(txt.Height())
	}
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
		img = &widgets.Image{
			Id:           id,
			CurrentColor: whiteColor,
			BoundingBox:  [4]float32{x, y, float32(100), float32(100)},
		}
		return img
	}).(*widgets.Image)

	clr := img.CurrentColor
	y -= wnd.currentWidgetSpace.scrlY
	// logic
	{
		hovered := c.hoverBehavior(wnd, utils.NewRectS(img.BoundingBox))
		if hovered {
			c.setActiveWidget(img.Id)
		}
		clicked = c.io.MouseClicked[0] && hovered
	}
	wnd.buffer.CreateRect(x, y, img.Width(), img.Height(), 0, draw.StraightCorners, tex.TextureId, clr, draw.ClipRectCompose{MainClipRect: wnd.mainWidgetSpace.ClipRect, ClipRect: wnd.currentWidgetSpace.ClipRect})
	img.UpdatePosition([4]float32{x, y, img.Width(), img.Height()})
	// wnd.addYcursor(img.Height())
	wnd.addCursor(img.Width(), img.Height())
	if !isRow {
		wnd.currentWidgetSpace.AddVirtualHeight(img.Height())
	}
	return clicked
}

func (c *UiContext) VSpace(id string) {
	wnd := c.windowStack.Peek()
	var s *widgets.VSpace
	x, y, isRow := wnd.currentWidgetSpace.getCursorPosition()

	s = wnd.getWidget(id, func() widgets.Widget {
		s := &widgets.VSpace{
			BoundingBox: [4]float32{x, y, float32(100), float32(20)},
			Id:          id,
		}
		return s
	}).(*widgets.VSpace)
	s.UpdatePosition([4]float32{x, y, float32(100), float32(20)})
	// wnd.addYcursor(s.Height())
	wnd.addCursor(0, s.Height())
	if !isRow {
		wnd.currentWidgetSpace.AddVirtualHeight(s.Height())
	}
}

func (c *UiContext) hoverBehavior(wnd *Window, rect utils.Rect) bool {
	inRect := utils.PointInRect(c.io.MousePos, utils.NewRect(rect.Min.X, rect.Min.Y, rect.Width(), rect.Height()))
	inWindow := RegionHit(c.io.MousePos.X, c.io.MousePos.Y, wnd.x, wnd.y+wnd.toolbar.h, wnd.w, wnd.h-wnd.toolbar.h)
	return c.ActiveWindow == wnd && inRect && inWindow
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
	// handle scrolling
	y -= wnd.currentWidgetSpace.scrlY
	// logic
	{
		hovered = c.hoverBehavior(wnd, utils.NewRectS(btn.BoundingBox))
		if hovered {
			c.setActiveWidget(btn.Id)
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
	wnd.buffer.CreateRect(x, y, w, h, 0, draw.StraightCorners, 0, btn.CurrentColor, draw.ClipRectCompose{MainClipRect: wnd.mainWidgetSpace.ClipRect, ClipRect: wnd.currentWidgetSpace.ClipRect})
	//
	btn.UpdatePosition([4]float32{x, y, w, h})
	wnd.addCursor(btn.Width(), btn.Height())
	if !isRow {
		wnd.currentWidgetSpace.AddVirtualHeight(btn.Height())
	}
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

func (c *UiContext) SubWidgetSpace(id string, widgFunc func()) {
	wnd := c.windowStack.Peek()
	var ws *WidgetSpace

	x, y, _ := wnd.currentWidgetSpace.getCursorPosition()
	ws, ok := c.widgSpaceCache.Get(id)
	if !ok {
		ws = newWidgetSpace(id, x, y, 100, 200)
		c.widgSpaceCache.Add(id, ws)

		wnd.widgSpaces = append(wnd.widgSpaces, ws)
	}
	var prevWS = wnd.currentWidgetSpace
	wnd.currentWidgetSpace = ws

	y -= prevWS.scrlY
	ws.X = x
	ws.Y = y
	ws.cursorY = y
	ws.cursorX = x

	outOfWindow := false
	if y < wnd.mainWidgetSpace.Y {
		outOfWindow = true
		ws.ClipRect = [4]float32{x, wnd.mainWidgetSpace.Y, ws.W - ws.verticalScrollbar.w, ws.H - (wnd.mainWidgetSpace.Y - y)}
	} else {
		ws.ClipRect = [4]float32{x, y, ws.W - ws.verticalScrollbar.w, ws.H}
	}

	// Scrollbar behavior
	if ws.isVertScrollShown {
		ws.vertScrollBar()
		if c.ActiveWindow == wnd && c.HoveredWindow == wnd && c.io.ScrollY != 0 && c.ActiveWidgetSpaceId == wnd.currentWidgetSpace.id {
			// wnd.mainWidgetSpace.handleMouseDrag()
			ws.handleMouseScroll(float32(c.io.ScrollY))
		}
	}

	var clip = draw.ClipRectCompose{MainClipRect: ws.ClipRect, ClipRect: wnd.mainWidgetSpace.ClipRect}

	if ws.isVertScrollShown {
		cl := [4]float32{ws.X, ws.Y, ws.W, ws.H}
		if outOfWindow {
			cl[1] = wnd.mainWidgetSpace.Y
		}
		var scrlClip = draw.ClipRectCompose{ClipRect: cl, MainClipRect: wnd.mainWidgetSpace.ClipRect}
		scrl := ws.verticalScrollbar

		wnd.buffer.CreateRect(scrl.x, scrl.y, scrl.w, scrl.h, 5, draw.AllRounded, 0, scrl.clr, scrlClip)
		wnd.buffer.CreateRect(scrl.bX, scrl.bY, scrl.bW, scrl.bH, 5, draw.AllRounded, 0, [4]float32{255, 0, 0, 1}, scrlClip)

		wnd.buffer.SeparateBuffer(0, scrlClip)
	}

	widgFunc()

	ws.checkVerScroll()

	//wnd.buffer.SeparateBuffer(0, wnd.buffer.InnerWindowSpace) // Make sure that we didn't miss anything
	//ws.AddVirtualHeight(c.CurrentStyle.BotMargin)
	wnd.buffer.SeparateBuffer(0, clip)
	ws.lastVirtualHeight = ws.virtualHeight
	ws.virtualHeight = 0
	//ws.cursorY = 0
	wnd.currentWidgetSpace = prevWS

	wnd.currentWidgetSpace.AddVirtualHeight(ws.H)
	wnd.addCursor(ws.W, ws.H)
}

func (c *UiContext) Row(id string, widgFunc func()) {
	wnd := c.windowStack.Peek()
	var row *widgets.HybridLayout
	x, y, _ := wnd.currentWidgetSpace.getCursorPosition()
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

// func (c *UiContext) EndRow() {
// 	wnd := c.windowStack.Peek()

// }

func (c *UiContext) EndWindow() {

	wnd := c.windowStack.Pop()

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
