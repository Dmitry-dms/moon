package ui

import (
	// "fmt"
	"fmt"
	// "math/rand"

	"github.com/Dmitry-dms/moon/pkg/gogl"
	"github.com/Dmitry-dms/moon/pkg/ui/widgets"
	"github.com/google/uuid"
	// "github.com/go-gl/mathgl/mgl32"
)

type Window struct {
	toolbar    Toolbar
	x, y       float32 // top-left corner
	w, h       float32
	active     bool
	Id         string
	drawList   []command
	rq         *RenderQueue
	outerRect  Rect
	minW, minH float32

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
		toolbar:   tb,
		x:         x,
		y:         y,
		w:         w,
		h:         h,
		drawList:  []command{},
		outerRect: Rect{Min: Vec2{x, y}, Max: Vec2{x + w, y + h}},
		rq:        NewRenderQueue(),
		minW:      200,
		minH:      200,

		startYshow: y+tb.h,
		endYshow:   y + h,
		// srcX:      x,
		// srcY:      y + tb.h + UiCtx.CurrentStyle.TopMargin,

		widgets: []widgets.Widget{},
		scrlBar: NewScrolBar(NewRect(x+w-10, y, 20, h), NewRect(x+w-10, y, 10, 50), [4]float32{150, 155, 155, 1}),
	}
	wnd.scrlY = wnd.toolbar.h
	return &wnd
}

func (w *Window) AddCommand(cmd command) {
	w.drawList = append(w.drawList, cmd)
}

func (w *Window) ClearDrawList() {
	w.drawList = []command{}
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
		r := 0 //rand.Intn(500)
		g := 0 //rand.Intn(300)
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

	prior := 0

	//Прямоугольник справа
	vResizeRect := Rect{Min: Vec2{wnd.x + wnd.w - 15, wnd.y}, Max: Vec2{wnd.x + wnd.w + 15, wnd.y + wnd.h}}
	hResizeRect := Rect{Min: Vec2{wnd.x, wnd.y + wnd.h - 15}, Max: Vec2{wnd.x + wnd.w, wnd.y + wnd.h + 15}}
	if PointInRect(c.io.MousePos, hResizeRect) && c.ActiveWindow == wnd {
		c.io.SetCursor(VResizeCursor)
		c.wantResizeH = true
	} else if PointInRect(c.io.MousePos, vResizeRect) && c.ActiveWindow == wnd {
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
	if c.io.IsDragging && c.ActiveWindow == wnd && PointInRect(c.io.MousePos, wnd.outerRect) && !c.wantResizeV && !c.wantResizeH {
		newX += c.io.MouseDelta.X
		newY += c.io.MouseDelta.Y
	}

	wnd.x = newX
	wnd.y = newY
	wnd.h = newH
	wnd.w = newW
	wnd.outerRect.Min = Vec2{wnd.x, wnd.y}
	wnd.outerRect.Max = Vec2{wnd.x + wnd.w, wnd.y + wnd.h}
	// wnd.startYshow = newY
	// wnd.endYshow = newH + newY

	cl := [4]float32{r, g, b, 0.8}
	cmdw := window_command{
		active: wnd.active,
		id:     wnd.Id,
		x:      wnd.x,
		y:      wnd.y,
		h:      wnd.h,
		w:      wnd.w,
		clr:    cl,
		toolbar: toolbar_command{
			h:   30,
			clr: [4]float32{255, 0, 0, 1},
		},
	}
	cmd := command{
		priority: prior,
		t:        WindowStartCmd,
		window:   &cmdw,
	}

	wnd.cursorX = wnd.x + UiCtx.CurrentStyle.LeftMargin
	wnd.cursorY = wnd.y + wnd.toolbar.h + UiCtx.CurrentStyle.TopMargin
	wnd.rq.AddCommand(cmd)

	c.windowStack.Push(window)
}

func (c *UiContext) srollbar(wnd *Window) {

	wnd.scrlBar.x = wnd.x + wnd.w - wnd.scrlBar.w
	wnd.scrlBar.y = wnd.y + wnd.toolbar.h
	wnd.scrlBar.h = wnd.h - wnd.toolbar.h

	scrollCommand := rounded_rect{
		x:      wnd.scrlBar.x,
		y:      wnd.scrlBar.y,
		w:      wnd.scrlBar.w,
		h:      wnd.scrlBar.h,
		clr:    wnd.scrlBar.clr,
		radius: 5,
	}
	sbCmd := command{
		rRect: &scrollCommand,
		t:     ScrollbarCmd,
		shown: wnd.isScrollShown,
	}
	wnd.rq.AddCommand(sbCmd)

	wnd.scrlBar.bX = wnd.x + wnd.w - wnd.scrlBar.w + 5
	wnd.scrlBar.bY = wnd.y + wnd.toolbar.h + wnd.scrlY

	if c.ActiveWindow == wnd && c.io.ScrollY != 0 && c.ActiveWidget == "" {
		factor := c.io.ScrollY * 10
		n := wnd.scrlY - float32(factor)
		//===================
		wnd.endYshow += float32(factor)
		wnd.startYshow += float32(factor)



		if n+wnd.cursorY < wnd.cursorY {

		} else if n+wnd.cursorY >= wnd.cursorY+wnd.h-wnd.toolbar.h-wnd.scrlBar.bH {

		} else {

			wnd.scrlY -= float32(factor)
		}

	}

	scrollBtnCommand := rounded_rect{
		x:      wnd.scrlBar.bX,
		y:      wnd.scrlBar.bY,
		w:      wnd.scrlBar.bW,
		h:      wnd.scrlBar.bH,
		clr:    [4]float32{255, 0, 0, 1},
		radius: 5,
	}

	sbtnCmd := command{
		rRect: &scrollBtnCommand,
		t:     ScrollButtonCmd,
		shown: wnd.isScrollShown,
	}
	wnd.rq.AddCommand(sbtnCmd)
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

	inRect := PointInRect(c.io.MousePos, NewRectS(btn.BoundingBox))

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

	rect := rounded_rect{
		x:      wnd.cursorX,
		y:      wnd.cursorY,
		w:      btn.Width(),
		h:      btn.Height(),
		clr:    clr,
		radius: 5,
	}

	btn.BoundingBox = [4]float32{rect.x, rect.y, rect.w, rect.h}

	cmd := command{
		rRect: &rect,
		shown: true,
	}
	if tex != nil {
		rect.texture = tex
		cmd.t = RoundedRectT
		rect.clr = whiteColor
	} else {
		cmd.t = RectType
	}

	if r, ok := wnd.isWindowEnd(btn.Height()); ok && r > 0 {
		rect.h = btn.Height() * r
		fmt.Println(rect.h, r)

	}

	wnd.rq.AddCommand(cmd)
	wnd.widgetCounter++

	wnd.cursorY += rect.h
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

	inRect := PointInRect(c.io.MousePos, NewRectS(img.BoundingBox))

	if wnd == c.ActiveWindow && inRect {
		c.ActiveWidget = img.GetId()
	}
	clicked := c.io.MouseClicked[0] && inRect

	// Create command and append it to slice
	{
		rect := rect_command{
			x:       wnd.cursorX,
			y:       wnd.cursorY,
			w:       img.Width(),
			h:       img.Height(),
			texture: tex,
			clr:     clr,
		}
		cmd := command{
			rect:  &rect,
			t:     RectTypeT,
			shown: true,
		}

		if r, ok := wnd.isWindowEnd(img.Height()); ok && r > 0 {
			rect.h = img.Height() * r
			rect.scaleFactor = r
			// fmt.Println(rect.h, r)
		}
		wnd.rq.AddCommand(cmd)
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
	var cmd command
	var rect rect_command

	if wnd.startYshow >=  wnd.cursorY {
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

	rect = rect_command{
		x:   x,
		y:   y,
		w:   w,
		h:   h,
		clr: clr,
	}

	if r, ok := wnd.isWindowEnd(rect.h); ok && r > 0 {
		rect.h = h * r

	} 

	// fmt.Println(wnd.h, wnd.endYshow - wnd.startYshow)
	// fmt.Println( wnd.startYshow,wnd.cursorY)
	

	{
		inRect = PointInRect(c.io.MousePos, NewRect(x, y, w, rect.h))

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
		cmd = command{
			rect:  &rect,
			t:     RectType,
			shown: true,
		}
	}

	if wnd.startYshow != wnd.y {
		// wnd.startYshow += wnd.y - wnd.startYshow
		// fmt.Println(wnd.startYshow,wnd.y)
		// rect.h -= wnd.y - wnd.startYshow
	}

	
	if rect.y > wnd.y+wnd.h {

	} else {

	}

	wnd.cursorY += btn.Height()
	wnd.rq.AddCommand(cmd)

	wnd.widgetCounter++

	return clicked
}

func (c *UiContext) EndWindow() {

	wnd := c.windowStack.Pop()
	//-----------------------------
	if wnd.isScrollShown {
		c.srollbar(wnd)
	} else {
		wnd.startYshow = wnd.y + wnd.toolbar.h
		wnd.endYshow = wnd.y + wnd.h
		wnd.scrlY = 0
	}

	cmd := command{
		priority: 0,
		t:        WindowCmd,
		// window:   cmdw,
	}

	if wnd.cursorY > wnd.y+wnd.h {
		wnd.isScrollShown = true
	} else {
		// fmt.Println(wnd.y+wnd.h,wnd.cursorY)
		wnd.isScrollShown = false
	}

	wnd.rq.AddCommand(cmd)
	c.currentWindow++
	wnd.widgetCounter = 0
}

func RegionHit(mouseX, mouseY, x, y, w, h float32) bool {
	return mouseX >= x && mouseY >= y && mouseX <= x+w && mouseY <= y+h
}
