package ui

import (
	// "fmt"
	"fmt"
	"math/rand"
	// "math/rand"

	"github.com/Dmitry-dms/moon/pkg/gogl"
	"github.com/Dmitry-dms/moon/pkg/ui/draw"
	"github.com/Dmitry-dms/moon/pkg/ui/utils"
	"github.com/Dmitry-dms/moon/pkg/ui/widgets"
	"github.com/google/uuid"
	// "github.com/go-gl/mathgl/mgl32"
)

type Window struct {
	toolbar Toolbar
	x, y    float32 // top-left corner
	w, h    float32
	active  bool
	Id      string
	// drawList   []command
	// rq         *RenderQueue
	outerRect  utils.Rect
	minW, minH float32

	//render
	buffer *draw.CmdBuffer

	//inner widgets
	cursorX, cursorY float32
	widgetCounter    int
	widgets          []widgets.Widget
	virtualHeight    float32 // сумма высот всех виджетов (для скроллинга)

	//scrollbar
	isScrollShown bool
	scrlBar       *Scrollbar
	scrlY         float32
}

func NewWindow(x, y, w, h float32) *Window {
	tb := NewToolbar(x, y, w, 30)
	wnd := Window{
		toolbar: tb,
		x:       x,
		y:       y,
		w:       w,
		h:       h,
		buffer:  draw.NewBuffer(UiCtx.camera),

		outerRect: utils.Rect{Min: utils.Vec2{X: x, Y: y}, Max: utils.Vec2{X: x + w, Y: y + h}},

		minW: 200,
		minH: 200,
		// srcX:      x,
		// srcY:      y + tb.h + UiCtx.CurrentStyle.TopMargin,

		widgets: []widgets.Widget{},
		scrlBar: NewScrolBar(utils.NewRect(x+w-10, y, 20, h), utils.NewRect(x+w-10, y, 10, 50), [4]float32{150, 155, 155, 1}),
	}
	wnd.scrlY = wnd.toolbar.h
	return &wnd
}

func (w *Window) AddCommand(cmd draw.Command) {
	w.buffer.AddCommand(cmd)
}

func generateId() string {
	return uuid.NewString()
}

var counter int = 0

const (
	defx, defy, defw, defh = 300, 100, 400, 500
)

var pointy = 0
var initialy = 0

func (c *UiContext) BeginWindow() {
	var window *Window
	if len(c.Windows) <= c.currentWindow {
		r := rand.Intn(500)
		g := rand.Intn(300)
		window = NewWindow(defx+float32(r), defy+float32(g), defw, defh)
		c.Windows = append(c.Windows, window)
		window.Id = generateId()
	} else {
		window = c.Windows[c.currentWindow]
	}

	wnd := window

	newX := wnd.x
	newY := wnd.y
	newH := wnd.h
	newW := wnd.w

	//Прямоугольник справа
	vResizeRect := utils.Rect{Min: utils.Vec2{X: wnd.x + wnd.w - 15, Y: wnd.y}, Max: utils.Vec2{X: wnd.x + wnd.w + 15, Y: wnd.y + wnd.h}}
	hResizeRect := utils.Rect{Min: utils.Vec2{X: wnd.x, Y: wnd.y + wnd.h - 15}, Max: utils.Vec2{X: wnd.x + wnd.w, Y: wnd.y + wnd.h + 15}}
	if utils.PointInRect(c.io.MousePos, hResizeRect) && c.ActiveWindow == wnd {
		c.io.SetCursor(VResizeCursor)
		c.wantResizeH = true
	} else if utils.PointInRect(c.io.MousePos, vResizeRect) && c.ActiveWindow == wnd {
		c.io.SetCursor(HResizeCursor)
		c.wantResizeV = true
	} else {
		c.io.SetCursor(ArrowCursor)
	}

	// Изменение размеров окна
	if c.wantResizeH && c.io.IsDragging && c.ActiveWindow == wnd {
		n := newH
		n += c.io.MouseDelta.Y
		if n > wnd.minH {
			newH = n
			if wnd.scrlY != 0 {
				wnd.scrlY -= c.io.MouseDelta.Y
			}
		}
	} else if c.wantResizeV && c.io.IsDragging && c.ActiveWindow == wnd {
		n := newW
		n += c.io.MouseDelta.X
		if n > wnd.minW {
			newW = n
		}
	}

	// Изменение положения окна
	if c.io.IsDragging && c.ActiveWindow == wnd && utils.PointInRect(c.io.MousePos, wnd.outerRect) && !c.wantResizeV && !c.wantResizeH {
		newX += c.io.MouseDelta.X
		newY += c.io.MouseDelta.Y
	}

	wnd.x = newX
	wnd.y = newY
	wnd.h = newH
	wnd.w = newW
	wnd.outerRect.Min = utils.Vec2{X: wnd.x, Y: wnd.y}
	wnd.outerRect.Max = utils.Vec2{X: wnd.x + wnd.w, Y: wnd.y + wnd.h}

	cl := [4]float32{r, g, b, 0.8}
	cmdw := draw.Window_command{
		Active: wnd.active,
		Id:     wnd.Id,
		X:      wnd.x,
		Y:      wnd.y,
		H:      wnd.h,
		W:      wnd.w,
		Clr:    cl,
	}
	cmd := draw.Command{
		Type:   draw.WindowStartCmd,
		Window: &cmdw,
	}

	wnd.cursorX = wnd.x + UiCtx.CurrentStyle.LeftMargin
	wnd.cursorY = wnd.y + wnd.toolbar.h + UiCtx.CurrentStyle.TopMargin
	wnd.AddCommand(cmd)

	c.windowStack.Push(window)
}

var visibleRatio float32
var step float32 = 40

func (c *UiContext) srollbar(wnd *Window) {

	wnd.scrlBar.x = wnd.x + wnd.w - wnd.scrlBar.w
	wnd.scrlBar.y = wnd.y + wnd.toolbar.h
	wnd.scrlBar.h = wnd.h - wnd.toolbar.h

	scrollCommand := draw.Rounded_rect{
		X:      wnd.scrlBar.x,
		Y:      wnd.scrlBar.y,
		W:      wnd.scrlBar.w,
		H:      wnd.scrlBar.h,
		Clr:    wnd.scrlBar.clr,
		Radius: 5,
	}
	sbCmd := draw.Command{
		RRect: &scrollCommand,
		Type:  draw.ScrollbarCmd,
		Shown: wnd.isScrollShown,
	}
	wnd.AddCommand(sbCmd)

	//проверки
	{

		if wnd.scrlY < 0 {
			wnd.scrlY = 0
		}
		var ratio float32 = (wnd.scrlY + wnd.h - wnd.toolbar.h) / wnd.virtualHeight

		if wnd.scrlY == 0 {
			visibleRatio = 0
		} else {
			visibleRatio = (ratio)*(wnd.h-wnd.toolbar.h) - wnd.scrlBar.bH
		}
	}

	wnd.scrlBar.bX = wnd.x + wnd.w - wnd.scrlBar.w + 5
	wnd.scrlBar.bY = wnd.y + wnd.toolbar.h + visibleRatio

	if wnd.scrlBar.bY >= wnd.y+wnd.h-wnd.scrlBar.bH {
		wnd.scrlBar.bY = wnd.y + wnd.h - wnd.scrlBar.bH
	}

	if c.ActiveWindow == wnd && c.io.ScrollY != 0 && c.ActiveWidget == "" {

		var factor float32 = float32(c.io.ScrollY) * float32(step)

		currentPos := wnd.scrlY + wnd.y + wnd.toolbar.h
		topBorder := wnd.y + wnd.toolbar.h
		botBorder := wnd.y + wnd.h + wnd.scrlY

		if currentPos <= topBorder {
			if factor > 0 {
				wnd.scrlY += step
			}
		} else if currentPos >= botBorder {
			if factor < 0 {
				wnd.scrlY -= step
			}
		} else {
			if wnd.cursorY-wnd.scrlY <= wnd.y+wnd.h {
				if factor < 0 {
					wnd.scrlY -= step
				}
			} else {
				wnd.scrlY += factor
			}
		}
	}

	scrollBtnCommand := draw.Rounded_rect{
		X:      wnd.scrlBar.bX,
		Y:      wnd.scrlBar.bY,
		W:      wnd.scrlBar.bW,
		H:      wnd.scrlBar.bH,
		Clr:    [4]float32{255, 0, 0, 1},
		Radius: 5,
	}

	sbtnCmd := draw.Command{
		RRect: &scrollBtnCommand,
		Type:  draw.ScrollButtonCmd,
		Shown: wnd.isScrollShown,
	}
	wnd.AddCommand(sbtnCmd)
}

var r, g, b float32 = 231, 158, 162

func (w *Window) addWidget(widg widgets.Widget) {
	w.widgets = append(w.widgets, widg)
	UiCtx.AddWidget(widg.GetId(), widg)
	w.virtualHeight += widg.Rectangle()[3]
}

var (
	whiteColor = [4]float32{255, 255, 255, 1}
)

func (wnd *Window) getWidget(w widgets.WidgetType) widgets.Widget {
	var widg widgets.Widget
	if len(wnd.widgets) == 0 || len(wnd.widgets) <= wnd.widgetCounter {
		switch w {
		case widgets.ButtonWidget:
			widg = &widgets.Button{
				Id:           generateId(),
				CurrentColor: [4]float32{67, 86, 205, 1},
				IsActive:     false,
				BoundingBox:  [4]float32{wnd.cursorX, wnd.cursorY, float32(100), float32(100)},
			}
		case widgets.ImageWidget:
			widg = &widgets.Image{
				Id:           generateId(),
				CurrentColor: whiteColor,
				BoundingBox:  [4]float32{wnd.cursorX, wnd.cursorY, float32(100), float32(100)},
			}
		case widgets.VerticalSpacingWidget:
			widg = &widgets.VSpace{
				BoundingBox: [4]float32{wnd.cursorX, wnd.cursorY, float32(100), float32(20)},
			}
		}
		wnd.addWidget(widg)
	} else {
		widg = wnd.widgets[wnd.widgetCounter]
	}
	return widg
}

func (c *UiContext) ButtonRR(tex *gogl.Texture) bool {

	wnd := c.windowStack.GetTop()
	var btn *widgets.Button

	btn = wnd.getWidget(widgets.ButtonWidget).(*widgets.Button)

	clr := btn.CurrentColor

	inRect := utils.PointInRect(c.io.MousePos, utils.NewRectS(btn.BoundingBox))

	if wnd == c.ActiveWindow && inRect {
		c.ActiveWidget = btn.GetId()
		// fmt.Println(btn.GetId())
		if c.io.MouseClicked[0] {
			btn.IsActive = !btn.IsActive
		}
	}

	if btn.IsActive {
		btn.SetColor([4]float32{150, btn.CurrentColor[1], btn.CurrentColor[2], btn.CurrentColor[3]})
	} else {
		btn.SetColor([4]float32{80, clr[1], clr[2], clr[3]})
		// btn.SetColor(whiteColor)
	}

	clicked := c.io.MouseClicked[0] && inRect

	//DEBUG
	scroll := float32(0)
	if c.io.ScrollY != 0 && c.ActiveWidget == btn.Id {
		scroll += float32(c.io.ScrollY) * 10
		btn.AddWidth(scroll)
		btn.AddHeight(scroll)
	}

	rect := draw.Rounded_rect{
		X:      wnd.cursorX,
		Y:      wnd.cursorY,
		W:      btn.Width(),
		H:      btn.Height(),
		Clr:    clr,
		Radius: 5,
	}

	btn.BoundingBox = [4]float32{rect.X, rect.Y, rect.W, rect.H}

	cmd := draw.Command{
		RRect: &rect,
		Shown: true,
	}
	if tex != nil {
		rect.Texture = tex
		cmd.Type = draw.RoundedRectT
		rect.Clr = whiteColor
	} else {
		cmd.Type = draw.RectType
	}

	wnd.AddCommand(cmd)
	wnd.widgetCounter++

	wnd.cursorY += rect.H
	return clicked
}

func (c *UiContext) Image(tex *gogl.Texture) bool {
	if tex == nil {
		fmt.Println("error")
		return false
	}

	wnd := c.windowStack.GetTop()
	var img *widgets.Image

	img = wnd.getWidget(widgets.ImageWidget).(*widgets.Image)

	clr := img.CurrentColor

	inRect := utils.PointInRect(c.io.MousePos, utils.NewRectS(img.BoundingBox))

	if wnd == c.ActiveWindow && inRect {
		c.ActiveWidget = img.GetId()
	}
	clicked := c.io.MouseClicked[0] && inRect

	y := wnd.cursorY

	// DEBUG
	y += wnd.scrlY
	//

	// Create command and append it to slice
	{
		rect := draw.Rect_command{
			X:       wnd.cursorX,
			Y:       y,
			W:       img.Width(),
			H:       img.Height(),
			Clr:     clr,
			Texture: tex,
		}
		cmd := draw.Command{
			Rect:  &rect,
			Type:  draw.RectTypeT,
			Shown: true,
		}
		wnd.AddCommand(cmd)
	}

	img.BoundingBox = [4]float32{wnd.cursorX, y, img.Width(), img.Height()}
	wnd.widgetCounter++

	wnd.cursorY += img.Height()
	return clicked
}

func (c *UiContext) VSpace() {

	wnd := c.windowStack.GetTop()
	var s *widgets.VSpace

	s = wnd.getWidget(widgets.VerticalSpacingWidget).(*widgets.VSpace)

	wnd.widgetCounter++

	wnd.cursorY += s.Height()
}

func (c *UiContext) Button() bool {

	wnd := c.windowStack.GetTop()
	var btn *widgets.Button
	var clicked, inRect bool
	var cmd draw.Command
	var rect draw.Rect_command

	btn = wnd.getWidget(widgets.ButtonWidget).(*widgets.Button)

	x := wnd.cursorX
	y := wnd.cursorY
	w := btn.BoundingBox[2]
	h := btn.BoundingBox[3]
	clr := btn.CurrentColor

	// DEBUG
	y -= wnd.scrlY
	//

	rect = draw.Rect_command{
		X:   x,
		Y:   y,
		W:   w,
		H:   h,
		Clr: clr,
	}

	{
		inRect = utils.PointInRect(c.io.MousePos, utils.NewRect(x, y, w, h))
		inWindow := RegionHit(c.io.MousePos.X, c.io.MousePos.Y, wnd.x, wnd.y+wnd.toolbar.h, wnd.w, wnd.h-wnd.toolbar.h)

		if wnd == c.ActiveWindow && inRect && inWindow {

			c.ActiveWidget = btn.GetId()
			if c.io.MouseClicked[0] {
				btn.IsActive = !btn.IsActive
			}
		}

		if btn.IsActive {
			btn.SetColor([4]float32{150, btn.CurrentColor[1], btn.CurrentColor[2], btn.CurrentColor[3]})
		} else {
			btn.SetColor([4]float32{80, clr[1], clr[2], clr[3]})
		}

		clicked = c.io.MouseClicked[0] && inRect && inWindow

		//DEBUG
		// scroll := float32(0)
		// if c.io.ScrollY != 0 && c.ActiveWidget == btn.Id {
		// 	scroll += float32(c.io.ScrollY) * 10
		// 	btn.AddWidth(scroll)
		// 	btn.AddHeight(scroll)
		// }
		cmd = draw.Command{
			Rect:  &rect,
			Type:  draw.RectType,
			Shown: true,
		}
	}

	wnd.cursorY += btn.Height()
	wnd.AddCommand(cmd)

	wnd.widgetCounter++

	return clicked
}

func (c *UiContext) EndWindow() {

	wnd := c.windowStack.Pop()
	//-----------------------------
	if wnd.isScrollShown {
		// pointy = int(wnd.virtualHeight) + int(wnd.y)
		c.srollbar(wnd)
	} else {
		// wnd.startYshow = wnd.y + wnd.toolbar.h
		// wnd.endYshow = wnd.y + wnd.h
		wnd.scrlY = 0
		pointy = 0
		initialy = 0
	}

	Toolbar := draw.Toolbar_command{
		X:   wnd.x,
		Y:   wnd.y,
		W:   wnd.w,
		H:   wnd.toolbar.h,
		Clr: wnd.toolbar.clr,
	}

	cmdToolbar := draw.Command{

		Type:    draw.ToolbarCmd,
		Toolbar: &Toolbar,
		// window:   cmdw,
	}
	wnd.AddCommand(cmdToolbar)

	cmd := draw.Command{
		Type: draw.WindowCmd,
	}

	if wnd.cursorY > wnd.y+wnd.h {
		wnd.isScrollShown = true
	} else {
		// fmt.Println(wnd.y+wnd.h,wnd.cursorY)
		wnd.isScrollShown = false
	}

	wnd.AddCommand(cmd)
	c.currentWindow++
	wnd.widgetCounter = 0
}

func RegionHit(mouseX, mouseY, x, y, w, h float32) bool {
	return mouseX >= x && mouseY >= y && mouseX <= x+w && mouseY <= y+h
}
