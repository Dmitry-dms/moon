package ui

import (
	"math/rand"

	"github.com/google/uuid"
	// "github.com/go-gl/mathgl/mgl32"
)

type Window struct {
	toolbar   Toolbar
	x, y      float32 // top-left corner
	w, h      float32
	active    bool
	Id        string
	drawList  []command
	rq        *RenderQueue
	outerRect Rect
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
		counter++
	} else {
		window = c.Windows[c.currentWindow]
	}

	c.windowStack.Push(window)
	cmd := command{
		t: WindowStartCmd,
		winStart: &window_start_command{
			x:  window.x,
			y:  window.y,
			id: window.Id,
		},
	}
	window.rq.AddCommand(cmd)
	// c.rq.AddCommand(cmd)
}

var r, g, b float32 = 231, 158, 162

func (c *UiContext) EndWindow() {

	wnd := c.windowStack.Pop()

	newX := wnd.x
	newY := wnd.y
	newH := wnd.h
	newW := wnd.w

	// tolbar := RegionHit(c.io.MousePos[0], c.io.MousePos[1], wnd.x, wnd.y, wnd.w, wnd.toolbar.h)
	// w := RegionHit(c.io.MousePos[0], c.io.MousePos[1], wnd.x, wnd.y, wnd.w, wnd.h)

	// if w {
	// 	c.HoveredWindow = wnd
	// 	if c.ActiveWindow == nil {
	// 		wnd.active = true
	// 		c.ActiveWindow = wnd
	// 	}
	// }

	//  Проверка на начало перетаскивания за границами окна
	// dragBounds := RegionHit(c.io.dragStartedMain[0], c.io.dragStartedMain[1], wnd.x, wnd.y, wnd.w, wnd.h)

	prior := 0

	if c.io.IsDragging && c.ActiveWindow == wnd && PointInRect(c.io.MousePos, wnd.outerRect) {
		newX += c.io.MouseDelta.X
		newY += c.io.MouseDelta.Y
	}
	// 	c.io.dragStartedMain = c.io.MousePos
	// 	// wnd.active = true
	// 	newX += c.io.dragDelta[0]
	// 	newY += c.io.dragDelta[1]
	// 	prior = 7
	// } else {
	// 	// wnd.active = false
	// 	prior = 1
	// }
	// if tolbar {
	// 	// c.ActiveWindow = wnd
	// } else {
	// 	// c.ActiveWindow = nil
	// }

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
		t:        WindowCmd,
		window:   &cmdw,
	}
	// c.windows = append(c.windows, window)

	// c.rq.AddCommand(cmd)
	wnd.rq.AddCommand(cmd)

	// counter++
	// c.ActiveWindow = nil
	c.currentWindow++
}

func RegionHit(mouseX, mouseY, x, y, w, h float32) bool {
	return mouseX >= x && mouseY >= y && mouseX <= x+w && mouseY <= y+h
}
