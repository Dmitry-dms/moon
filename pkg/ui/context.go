package ui

import (
	// "fmt"

	"github.com/Dmitry-dms/moon/pkg/gogl"
	"github.com/Dmitry-dms/moon/pkg/ui/cache"
	"github.com/Dmitry-dms/moon/pkg/ui/render"
)

type UiContext struct {
	rq *RenderQueue

	renderer UiRenderer
	camera   *gogl.Camera
	io       *Io

	//Widgets
	Windows           []*Window
	activeWidget      string
	ActiveWindow      *Window
	currentWindow     int
	PriorWindow       *Window
	HoveredWindow     *Window
	LastHoveredWindow *Window

	//cache
	idCache     *cache.RamCache[Window]
	windowStack stack[*Window]

	//refactor
	Time float32
}

func NewContext(frontRenderer UiRenderer, camera *gogl.Camera) *UiContext {
	c := UiContext{
		rq:          NewRenderQueue(),
		renderer:    frontRenderer,
		camera:      camera,
		io:          NewIo(),
		Windows:     make([]*Window, 0),
		idCache:     cache.NewRamCache[Window](),
		windowStack: Stack[*Window](),
	}

	return &c
}

func (c *UiContext) Io() *Io {
	return c.io
}

func (c *UiContext) NewFrame() {

	c.UpdateMouseInputs()
	c.findHoveredWindow()
	c.renderer.NewFrame()
}
func (c *UiContext) pushWindowFront(w *Window) {
	for i := len(c.Windows) - 1; i >= 0; i-- {
		if c.Windows[i] == w {
			if i == len(c.Windows)-1 {
				return
			}
			c.Windows[i] = c.Windows[len(c.Windows)-1]
			c.Windows[len(c.Windows)-1] = w
			return
		}
	}
}

var initHover = false

func (c *UiContext) findHoveredWindow() {
	var hovered *Window
	if len(c.Windows) == 0 {
		return
	}

	for i := 0; i <= len(c.Windows)-1; i++ {
		window := c.Windows[i]
		// if i == len(c.Windows) - 1 && c.ActiveWindow == nil {
		// 	window.active = true
		// 	c.ActiveWindow = window
		// }
		// if !window.active {
		// 	continue
		// }
		bb := window.outerRect

		if !bb.Contains(c.io.MousePos) {
			continue
		}

		if c.io.MouseClicked[0] && c.ActiveWindow != window {
			if !PointInRect(c.io.MousePos, c.ActiveWindow.outerRect) {
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
		c.ActiveWindow = c.Windows[len(c.Windows)-1]
	}
	c.HoveredWindow = hovered

}

func (c *UiContext) UpdateMouseInputs() {

	io := c.Io()

	if io.IsMousePosValid(&io.MousePos) && io.IsMousePosValid(&io.MousePosPrev) {
		io.MouseDelta = io.MousePos.Sub(io.MousePosPrev)
	} else {
		io.MouseDelta = Vec2{0, 0}
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
				var delta Vec2
				if io.IsMousePosValid(&io.MousePos) {
					delta = io.MousePos.Sub(io.MouseClickedPos[i])
				} else {
					delta = Vec2{0, 0}
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
			io.MouseDragMaxDistanceSqr[i] = Max(io.MouseDragMaxDistanceSqr[i], deltaSqrPos)
		}
		// We provide io.MouseDoubleClicked[] as a legacy service
		io.MouseDoubleClicked[i] = (io.MouseClickedCount[i] == 2)

	}

}

func (c *UiContext) EndFrame() {

	// for _, wnd := range c.windows {
	// 	hovered := RegionHit(c.io.MousePos[0], c.io.MousePos[1], wnd.x, wnd.y, wnd.w, wnd.h)
	// 	if hovered {
	// 		c.ActiveWindow = wnd
	// 	}
	// }
	for _, v := range c.Windows {
		cmds := v.rq.commands

		for i := 0; i < v.rq.CmdCount; i++ {

			comm := cmds[i]
			switch comm.t {
			case RectType:
				r := comm.rect
				c.renderer.Rectangle(r.x, r.y, r.w, r.h, r.clr)
			case Triangle:
				tr := comm.triangle
				c.renderer.Trinagle(tr.x0, tr.y0, tr.x1, tr.y1, tr.x2, tr.y2, tr.clr)
			case RoundedRect:
				rr := comm.rRect
				c.renderer.RoundedRectangle(rr.x, rr.y, rr.w, rr.h, rr.radius, rr.clr)
			case WindowCmd:
				wnd := comm.window

				// fmt.Println(wnd.id, wnd.active)

				size := c.camera.GetProjectionSize()

				c.renderer.RoundedRectangleR(wnd.x, size.Y()-wnd.y, wnd.w, wnd.h, 10, render.AllRounded, comm.window.clr)
				c.renderer.RoundedRectangleR(wnd.x, size.Y()-wnd.y, wnd.w, wnd.toolbar.h, 10, render.TopRect, comm.window.toolbar.clr)
			default:
			}
		}
		v.rq.clearCommands()
	}

	c.renderer.Draw(c.camera)
	c.renderer.End()
	// c.windows = []*Window{}

	c.currentWindow = 0

	// if !c.checkMousePos() {
	// 	c.HoveredWindow = nil
	// }

	//io
	c.io.dragDelta = Vec2{0, 0}

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
	for _, v := range c.Windows {
		if lastWindow.x == v.x && lastWindow.y == v.y &&
			lastWindow.h == v.h && lastWindow.w == v.w {
			c.ActiveWindow = v
		}
		lastWindow = v
	}
}

func (c *UiContext) checkMousePos() bool {
	hovered := false
	for _, wnd := range c.Windows {
		check := RegionHit(c.io.MousePos.X, c.io.MousePos.Y, wnd.x, wnd.y, wnd.w, wnd.h)
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
