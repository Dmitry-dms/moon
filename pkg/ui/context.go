package ui

import (
	// "fmt"

	"github.com/Dmitry-dms/moon/pkg/gogl"
	"github.com/Dmitry-dms/moon/pkg/ui/cache"
	"github.com/Dmitry-dms/moon/pkg/ui/render"
	"github.com/go-gl/mathgl/mgl32"
)

type UiContext struct {
	rq *RenderQueue

	renderer UiRenderer
	camera   *gogl.Camera
	io       *Io

	//Widgets
	windows       []*Window
	activeWidget  string
	ActiveWindow  *Window
	currentWindow int
	PriorWindow   *Window
	HoveredWindow *Window

	//cache
	idCache     *cache.RamCache[Window]
	windowStack stack[*Window]
}

func NewContext(frontRenderer UiRenderer, camera *gogl.Camera) *UiContext {
	c := UiContext{
		rq:          NewRenderQueue(),
		renderer:    frontRenderer,
		camera:      camera,
		io:          NewIo(),
		windows:     make([]*Window, 0),
		idCache:     cache.NewRamCache[Window](),
		windowStack: Stack[*Window](),
	}

	return &c
}

func (c *UiContext) Io() *Io {
	return c.io
}

func (c *UiContext) NewFrame() {



	c.renderer.NewFrame()
}

func (c *UiContext) UpdateMouseInputs() {
	
}

func (c *UiContext) EndFrame() {

	// for _, wnd := range c.windows {
	// 	hovered := RegionHit(c.io.MousePos[0], c.io.MousePos[1], wnd.x, wnd.y, wnd.w, wnd.h)
	// 	if hovered {
	// 		c.ActiveWindow = wnd
	// 	}
	// }

	cmds := c.rq.Commands()


	for i := 0; i < c.rq.CmdCount; i++ {
		v := cmds[i]
		switch v.t {
		case Rect:
			r := v.rect
			c.renderer.Rectangle(r.x, r.y, r.w, r.h, r.clr)
		case Triangle:
			tr := v.triangle
			c.renderer.Trinagle(tr.x0, tr.y0, tr.x1, tr.y1, tr.x2, tr.y2, tr.clr)
		case RoundedRect:
			rr := v.rRect
			c.renderer.RoundedRectangle(rr.x, rr.y, rr.w, rr.h, rr.radius, rr.clr)
		case WindowCmd:
			wnd := v.window

			// fmt.Println(wnd.id, wnd.active)

			size := c.camera.GetProjectionSize()

			c.renderer.RoundedRectangleR(wnd.x, size.Y()-wnd.y, wnd.w, wnd.h, 10, render.AllRounded, v.window.clr)
			c.renderer.RoundedRectangleR(wnd.x, size.Y()-wnd.y, wnd.w, wnd.toolbar.h, 10, render.TopRect, v.window.toolbar.clr)
		default:
		}
	}
	c.rq.clearCommands()

	c.renderer.Draw(c.camera)
	c.renderer.End()
	// c.windows = []*Window{}

	c.currentWindow = 0

	// if !c.checkMousePos() {
	// 	c.HoveredWindow = nil
	// }

	//io
	c.io.dragDelta = mgl32.Vec2{0, 0}

	// lastWindow := len(c.windows)
	// c.PriorWindow = c.windows[lastWindow-1]
	// c.checkWindowPriority()
}

func repl(m []int, x1, x2 int) []int {
	var res []int
	if x1 == 0 {
		f := m[x1 : x2+1]
		r := m[x2+1:]
		r = append(r, f...)
		res = r
	} else {
		f := m[x1 : x2+1]
		ost := []int{}
		ost = append(ost, m[:x1]...)
		ost = append(ost, m[x2+1:]...)
		ost = append(ost, f...)
		res = ost
	}
	return res
}

func (c *UiContext) checkWindowPriority() {

	var lastWindow = NewWindow(0, 0, 0, 0)
	for _, v := range c.windows {
		if lastWindow.x == v.x && lastWindow.y == v.y &&
			lastWindow.h == v.h && lastWindow.w == v.w {
			c.ActiveWindow = v
		}
		lastWindow = v
	}
}

func (c *UiContext) checkMousePos() bool {
	hovered := false
	for _, wnd := range c.windows {
		check := RegionHit(c.io.MousePos[0], c.io.MousePos[1], wnd.x, wnd.y, wnd.w, wnd.h)
		if check != hovered {
			hovered = true
		}
	}
	return hovered
}

type UiRenderer interface {
	NewFrame()
	Rectangle(x, y, w, h float32, clr [4]float32)
	RectangleR(x, y, w, h float32, clr [4]float32)
	Trinagle(x0, y0, x1, y1, x2, y2 float32, clr [4]float32)
	Circle(x, y, radius float32, steps int, clr [4]float32)
	Line(x0, y0, x1, y1 float32, thick int, clr [4]float32)
	RoundedRectangle(x, y, w, h float32, radius int, clr [4]float32)
	RoundedRectangleR(x, y, w, h float32, radius int, shape render.RoundedRectShape, clr [4]float32)
	Draw(camera *gogl.Camera)
	End()
}
