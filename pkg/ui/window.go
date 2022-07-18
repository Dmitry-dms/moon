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

	//shown space
	startYshow float32
	endYshow   float32
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
		// drawList:  []command{},
		outerRect: utils.Rect{Min: utils.Vec2{X: x, Y: y}, Max: utils.Vec2{X: x + w, Y: y + h}},
		// rq:        NewRenderQueue(),
		minW: 200,
		minH: 200,

		startYshow: y + tb.h,
		endYshow:   y + h,
		// srcX:      x,
		// srcY:      y + tb.h + UiCtx.CurrentStyle.TopMargin,

		widgets: []widgets.Widget{},
		scrlBar: NewScrolBar(utils.NewRect(x+w-10, y, 20, h), utils.NewRect(x+w-10, y, 10, 50), [4]float32{150, 155, 155, 1}),
	}
	wnd.scrlY = wnd.toolbar.h
	return &wnd
}

func (w *Window) AddCommand(cmd draw.Command) {
	// w.drawList = append(w.drawList, cmd)
	// w.drawData.cmdLists
	w.buffer.AddCommand(cmd)
}

func (w *Window) ClearDrawList() {
	// w.drawList = []command{}
	// w.buffer.clear()
}

func generateId() string {
	return uuid.NewString()
}

var counter int = 0

const (
	defx, defy, defw, defh = 300, 100, 400, 500
)

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


	

	// fmt.Println(c.Windows)

	wnd := window

	newX := wnd.x
	newY := wnd.y
	newH := wnd.h
	newW := wnd.w

	

	// prior := 0

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
	// wnd.startYshow = newY
	// wnd.endYshow = newH + newY

	cl := [4]float32{r, g, b, 0.8}
	cmdw := draw.Window_command{
		Active: wnd.active,
		Id:     wnd.Id,
		X:      wnd.x,
		Y:      wnd.y,
		H:      wnd.h,
		W:      wnd.w,
		Clr:    cl,
		Toolbar: draw.Toolbar_command{
			H:   30,
			Clr: [4]float32{255, 0, 0, 1},
		},
	}
	cmd := draw.Command{
		// priority: prior,
		Type:   draw.WindowStartCmd,
		Window: &cmdw,
	}


	// rect2 := draw.Rect_command{
	// 	X:   newX,
	// 	Y:   newY,
	// 	W:   newW,
	// 	H:   newH,
	// 	Clr: [4]float32{0, 0, 0, 1},
	// }
	// cmd2 := draw.Command{
	// 	// priority: prior,
	// 	Type:   draw.RectType,
	// 	Rect: &rect2,
	// }
	// wnd.AddCommand(cmd2)
	wnd.cursorX = wnd.x + UiCtx.CurrentStyle.LeftMargin
	wnd.cursorY = wnd.y + wnd.toolbar.h + UiCtx.CurrentStyle.TopMargin
	wnd.AddCommand(cmd)

	c.windowStack.Push(window)
}

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

	wnd.scrlBar.bX = wnd.x + wnd.w - wnd.scrlBar.w + 5
	wnd.scrlBar.bY = wnd.y + wnd.toolbar.h + wnd.scrlY

	if c.ActiveWindow == wnd && c.io.ScrollY != 0 && c.ActiveWidget == "" {
		factor := c.io.ScrollY * 10
		// n := wnd.scrlY - float32(factor)
		//===================
		// wnd.endYshow += float32(factor)
		// wnd.startYshow += float32(factor)

		wnd.scrlY -= float32(factor)

		// if n+wnd.cursorY < wnd.cursorY {

		// } else if n+wnd.cursorY >= wnd.cursorY+wnd.h-wnd.toolbar.h-wnd.scrlBar.bH {

		// } else {

		// }

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

func (w *Window) isWindowEnd(widgHeight float32) (ratio float32, intersect bool) {
	if w.cursorY+widgHeight > w.h+w.y {
		// w.scrlY = 0
		// w.isScrollShown = true
		intersect = true
		g := w.h + w.y - w.cursorY
		ratio = float32(g / widgHeight)
	} else {
		// w.isScrollShown = false
	}
	return
}

func (w *Window) addWidget(widg widgets.Widget) {
	w.widgets = append(w.widgets, widg)
	UiCtx.AddWidget(widg.GetId(), widg)
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
				Height: 20,
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

	if r, ok := wnd.isWindowEnd(btn.Height()); ok && r > 0 {
		rect.H = btn.Height() * r
		// fmt.Println(rect.h, r)

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

	// Create command and append it to slice
	{
		rect := draw.Rect_command{
			X:       wnd.cursorX,
			Y:       wnd.cursorY,
			W:       img.Width(),
			H:       img.Height(),
			Clr:     clr,
			Texture: tex,
		}
		cmd := draw.Command{
			Rect:  &rect,
			Type:     draw.RectTypeT,
			Shown: true,
		}

		if r, ok := wnd.isWindowEnd(img.Height()); ok && r > 0 {
			rect.H = img.Height() * r
			rect.ScaleFactor = r
			// fmt.Println(rect.h, r)
		}
		wnd.AddCommand(cmd)
	}

	img.BoundingBox = [4]float32{wnd.cursorX, wnd.cursorY, img.Width(), img.Height()}
	wnd.widgetCounter++

	wnd.cursorY += img.Height()
	return clicked
}

func (c *UiContext) VSpace() {

	wnd := c.windowStack.GetTop()
	var s *widgets.VSpace

	s = wnd.getWidget(widgets.VerticalSpacingWidget).(*widgets.VSpace)

	// img.BoundingBox = [4]float32{wnd.srcX, wnd.srcY, img.Width(), img.Height()}
	wnd.widgetCounter++

	wnd.cursorY += s.Height
}

func (c *UiContext) Button() bool {

	wnd := c.windowStack.GetTop()
	var btn *widgets.Button
	var clicked, inRect bool
	var cmd draw.Command
	var rect draw.Rect_command

	if wnd.startYshow >= wnd.cursorY {
		// fmt.Println(wnd.startYshow,wnd.cursorY)
		// rect.h -= wnd.startYshow - wnd.cursorY
		// rect.y += wnd.startYshow - wnd.cursorY
		wnd.cursorY -= wnd.scrlY
	}

	btn = wnd.getWidget(widgets.ButtonWidget).(*widgets.Button)

	x := wnd.cursorX
	y := wnd.cursorY
	w := btn.BoundingBox[2]
	h := btn.BoundingBox[3]
	clr := btn.CurrentColor

	// DEBUG
	y += wnd.scrlY
	//

	rect = draw.Rect_command{
		X:   x,
		Y:   y,
		W:   w,
		H:   h,
		Clr: clr,
	}

	// if r, ok := wnd.isWindowEnd(rect.H); ok && r > 0 {
	// 	rect.H = h * r

	// }

	// fmt.Println(wnd.h, wnd.endYshow - wnd.startYshow)
	// fmt.Println( wnd.startYshow,wnd.cursorY)

	{
		inRect = utils.PointInRect(c.io.MousePos, utils.NewRect(x, y, w, rect.H))

		if wnd == c.ActiveWindow && inRect {

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

		clicked = c.io.MouseClicked[0] && inRect

		//DEBUG
		// scroll := float32(0)
		// if c.io.ScrollY != 0 && c.ActiveWidget == btn.Id {
		// 	scroll += float32(c.io.ScrollY) * 10
		// 	btn.AddWidth(scroll)
		// 	btn.AddHeight(scroll)
		// }
		cmd = draw.Command{
			Rect:  &rect,
			Type:     draw.RectType,
			Shown: true,
		}
	}

	if wnd.startYshow != wnd.y {
		// wnd.startYshow += wnd.y - wnd.startYshow
		// fmt.Println(wnd.startYshow,wnd.y)
		// rect.h -= wnd.y - wnd.startYshow
	}

	// if rect.y > wnd.y+wnd.h {

	// } else {

	// }

	wnd.cursorY += btn.Height()
	wnd.AddCommand(cmd)

	wnd.widgetCounter++

	return clicked
}

func (c *UiContext) EndWindow() {

	wnd := c.windowStack.Pop()
	//-----------------------------
	if wnd.isScrollShown {
		c.srollbar(wnd)
	} else {
		// wnd.startYshow = wnd.y + wnd.toolbar.h
		// wnd.endYshow = wnd.y + wnd.h
		wnd.scrlY = 0
	}

	cmd := draw.Command{

		Type:        draw.WindowCmd,
		// window:   cmdw,
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
