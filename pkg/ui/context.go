package ui

import (
	"github.com/Dmitry-dms/moon/pkg/gogl"
	"github.com/Dmitry-dms/moon/pkg/ui/cache"
	"github.com/Dmitry-dms/moon/pkg/ui/render"
)

type UiContext struct {
	rq       *RenderQueue
	renderer UiRenderer
	camera   *gogl.Camera
	io       *Io

	//Widgets
	windows      []Window
	activeWidget string
	ActiveWindow *Window

	//cache
	idCache *cache.RamCache[Window]
}

func NewContext(renderer UiRenderer, camera *gogl.Camera) *UiContext {
	c := UiContext{
		rq:       NewRenderQueue(),
		renderer: renderer,
		camera:   camera,
		io:       NewIo(),
		windows:  []Window{},
		idCache:  cache.NewRamCache[Window](),
	}
	return &c
}

func (c *UiContext) Io() *Io {
	return c.io
}

func (c *UiContext) NewFrame() {

	c.renderer.NewFrame()
}

func (c *UiContext) EndFrame() {

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
			size := c.camera.GetProjectionSize()

			tolbar := RegionHit(c.io.MousePos[0], c.io.MousePos[1], wnd.x, wnd.y, wnd.w, wnd.toolbar.h)
			w := RegionHit(c.io.MousePos[0], c.io.MousePos[1], wnd.x, wnd.y, wnd.w, wnd.h)

			if tolbar {

			}
			if w {
				w1, _ := c.idCache.Get("debug")
				c.ActiveWindow = w1
			} else {
				c.ActiveWindow = nil
			}

			c.renderer.RoundedRectangleR(wnd.x, size.Y()-wnd.y, wnd.w, wnd.h, 10, render.AllRounded, wnd.clr)
			c.renderer.RoundedRectangleR(wnd.x, size.Y()-wnd.y, wnd.w, wnd.toolbar.h, 10, render.TopRect, wnd.toolbar.clr)
		default:
		}
	}
	c.rq.clearCommands()
	c.renderer.Draw(c.camera)
	c.renderer.End()
	c.windows = []Window{}
}

func (c *UiContext) Button(name string, pushed *bool, clr [4]float32) bool {
	cmd := command{
		rRect: &rounded_rect{
			x:      100,
			y:      100,
			w:      200,
			h:      100,
			radius: 5,
			clr:    clr,
		},
		t: RoundedRect,
	}
	c.rq.AddCommand(cmd)

	return false
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
