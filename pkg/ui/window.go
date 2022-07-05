package ui

import (
	// "fmt"
	"math/rand"

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
	srcX, srcY    float32
	widgetCounter int
	widgets       []widget
}

type widget interface {
	GetColor() [4]float32
	GetId() string
}

type button struct {
	isActive     bool
	currentColor [4]float32
	id           string
}

func (b *button) GetColor() [4]float32 {
	return b.currentColor
}
func (b *button) GetId() string {
	return b.id
}
func (b *button) SetColor(clr [4]float32) {
	b.currentColor = clr
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
		// srcX:      x,
		// srcY:      y + tb.h + UiCtx.CurrentStyle.TopMargin,

		widgets: []widget{},
	}
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

	wnd.srcX = wnd.x + UiCtx.CurrentStyle.LeftMargin
	wnd.srcY = wnd.y + wnd.toolbar.h + UiCtx.CurrentStyle.TopMargin
	wnd.rq.AddCommand(cmd)

	c.windowStack.Push(window)
}

var r, g, b float32 = 231, 158, 162

func (w *Window) addWidget(widg widget) {
	w.widgets = append(w.widgets, widg)
	UiCtx.AddWidget(widg.GetId(), widg)
}

func (c *UiContext) Button() bool {

	wnd := c.windowStack.GetTop()
	var btn *button

	if len(wnd.widgets) == 0 {
		btn = &button{
			id:           generateId(),
			currentColor: [4]float32{67, 86, 205, 1},
			isActive:     false,
		}
		wnd.addWidget(btn)
	} else {
		btn = wnd.widgets[wnd.widgetCounter].(*button)
	}

	x := wnd.srcX
	y := wnd.srcY
	w := float32(100)
	h := float32(100)

	clr := btn.currentColor

	inRect := PointInRect(c.io.MousePos, NewRect(x, y, w, h))

	if wnd == c.ActiveWindow && inRect {
		c.ActiveWidget = btn.GetId()
		if c.io.MouseClicked[0] {
			btn.isActive = !btn.isActive
		}
	} else {
		c.ActiveWidget = ""
	}

	if btn.isActive {
		btn.SetColor([4]float32{150, clr[1], clr[2], clr[3]})
	} else {
		btn.SetColor([4]float32{80, clr[1], clr[2], clr[3]})
	}

	clicked := c.io.MouseClicked[0] && inRect

	rect := rounded_rect{
		x:   x,
		y:   y,
		w:   w,
		h:   h,
		clr: clr,
		radius: 5,
	}
	cmd := command{
		t:    RoundedRect,
		rRect: &rect,
	}
	wnd.rq.AddCommand(cmd)
	wnd.widgetCounter++
	return clicked
}

func (c *UiContext) EndWindow() {

	wnd := c.windowStack.Pop()

	cmd := command{
		priority: 0,
		t:        WindowCmd,
		// window:   cmdw,
	}
	wnd.rq.AddCommand(cmd)
	c.currentWindow++
	wnd.widgetCounter = 0
}

func RegionHit(mouseX, mouseY, x, y, w, h float32) bool {
	return mouseX >= x && mouseY >= y && mouseX <= x+w && mouseY <= y+h
}
