package ui

import (
	"github.com/Dmitry-dms/moon/pkg/fonts"
	"github.com/Dmitry-dms/moon/pkg/gogl"
	"github.com/Dmitry-dms/moon/pkg/ui/cache"
	"github.com/Dmitry-dms/moon/pkg/ui/draw"
	"github.com/Dmitry-dms/moon/pkg/ui/render"
	"github.com/Dmitry-dms/moon/pkg/ui/utils"
	"github.com/Dmitry-dms/moon/pkg/ui/widgets"
)

var UiCtx *UiContext

func init() {
	UiCtx = NewContext(nil, nil)
}

type UiContext struct {
	// rq *RenderQueue

	renderer UiRenderer
	camera   *gogl.Camera
	io       *Io

	//Widgets
	Windows           []*Window
	sortedWindows     []*Window
	ActiveWidget      string
	ActiveWindow      *Window
	currentWindow     int
	PriorWindow       *Window
	HoveredWindow     *Window
	LastHoveredWindow *Window

	//cache
	idCache      *cache.RamCache[Window]
	windowStack  utils.Stack[*Window]
	widgetsCache *cache.RamCache[widgets.Widget]

	//refactor
	Time float32

	//intent
	wantResizeH, wantResizeV bool

	//style
	CurrentStyle Style

	//fonts
	font *fonts.Font
}

func NewContext(frontRenderer UiRenderer, camera *gogl.Camera) *UiContext {
	c := UiContext{
		// rq:            NewRenderQueue(),
		renderer:      frontRenderer,
		camera:        camera,
		io:            NewIo(),
		Windows:       make([]*Window, 0),
		sortedWindows: make([]*Window, 0),
		idCache:       cache.NewRamCache[Window](),
		widgetsCache:  cache.NewRamCache[widgets.Widget](),
		windowStack:   utils.NewStack[*Window](),
		CurrentStyle:  DefaultStyle,
	}
	

	return &c
}

func (c *UiContext) UploadFont(path string, size int) {
	c.font = fonts.NewFont(path)
}

func (c *UiContext) Initialize(frontRenderer UiRenderer, camera *gogl.Camera) {
	c.renderer = frontRenderer
	c.camera = camera
}

func (c *UiContext) AddWidget(id string, w widgets.Widget) {
	c.widgetsCache.Add(id, &w)
}

func (c *UiContext) GetWidget(id string) *widgets.Widget {
	w, _ := c.widgetsCache.Get(id)
	return w
}

func (c *UiContext) Io() *Io {
	return c.io
}

func (c *UiContext) NewFrame() {
	// c.sortedWindows = c.Windows

	c.UpdateMouseInputs()

	c.renderer.NewFrame()
}
func (c *UiContext) pushWindowFront(w *Window) {
	for i := len(c.sortedWindows) - 1; i >= 0; i-- {
		if c.sortedWindows[i] == w {
			if i == len(c.sortedWindows)-1 {
				return
			}
			c.sortedWindows[i] = c.sortedWindows[len(c.sortedWindows)-1]
			c.sortedWindows[len(c.sortedWindows)-1] = w
			return
		}
	}
}

var initHover = false

func (c *UiContext) findHoveredWindow() {
	var hovered *Window
	if len(c.sortedWindows) == 0 {
		return
	}

	for i := 0; i <= len(c.sortedWindows)-1; i++ {
		window := c.sortedWindows[i]
		bb := window.outerRect

		if !bb.Contains(c.io.MousePos) {
			continue
		}
		if c.io.MouseClicked[0] && c.ActiveWindow != window {
			if !utils.PointInRect(c.io.MousePos, c.ActiveWindow.outerRect) {
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
		c.ActiveWindow = c.sortedWindows[len(c.sortedWindows)-1]
	}
	c.HoveredWindow = hovered

}

func (c *UiContext) UpdateMouseInputs() {

	io := c.Io()

	if io.IsMousePosValid(&io.MousePos) && io.IsMousePosValid(&io.MousePosPrev) {
		io.MouseDelta = io.MousePos.Sub(io.MousePosPrev)
	} else {
		io.MouseDelta = utils.Vec2{0, 0}
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
				var delta utils.Vec2
				if io.IsMousePosValid(&io.MousePos) {
					delta = io.MousePos.Sub(io.MouseClickedPos[i])
				} else {
					delta = utils.Vec2{0, 0}
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
			io.MouseDragMaxDistanceSqr[i] = utils.Max(io.MouseDragMaxDistanceSqr[i], deltaSqrPos)
		}
		// We provide io.MouseDoubleClicked[] as a legacy service
		io.MouseDoubleClicked[i] = (io.MouseClickedCount[i] == 2)

	}

}

func copyWindows(w []*Window) []*Window {
	r := make([]*Window, len(w))
	for i, v := range w {
		r[i] = v
	}
	return r
}

var lastWinL = 0

func (c *UiContext) EndFrame() {

	// Если количество окон не изменилось, в копировании нет нужды
	if lastWinL != len(c.Windows) {
		c.sortedWindows = copyWindows(c.Windows)
		lastWinL = len(c.Windows)
	}

	c.findHoveredWindow()
	if len(c.sortedWindows) == 0 {
		return
	}

	// for _, v := range c.sortedWindows {
	// 	cmds := v.rq.commands
	// 	// size := c.camera.GetProjectionSize()
	// 	// fmt.Println("---------------------------")
	// 	// fmt.Println(int32(v.y))

	// 	// gl.Scissor(int32(v.x), int32(v.y), int32(v.w), int32(v.h))
	// 	for i := 0; i < v.rq.CmdCount; i++ {
	// 		comm := cmds[i]
	// 		switch comm.t {
	// 		case RectType:
	// 			r := comm.rect

	// 			size := c.camera.GetProjectionSize()
	// 			c.renderer.RectangleR(r.x, size.Y()-r.y, r.w, r.h, r.clr)
	// 		case RectTypeT:
	// 			r := comm.rect

	// 			size := c.camera.GetProjectionSize()
	// 			if r.scaleFactor == 0 {
	// 				c.renderer.RectangleT(r.x, size.Y()-r.y, r.w, r.h, r.texture, 0, 1, 0, r.clr)
	// 			} else {
	// 				c.renderer.RectangleT(r.x, size.Y()-r.y, r.w, r.h, r.texture, 0, 1, r.scaleFactor, r.clr)
	// 			}
	// 		case Triangle:
	// 			tr := comm.triangle
	// 			c.renderer.Trinagle(tr.x0, tr.y0, tr.x1, tr.y1, tr.x2, tr.y2, tr.clr)
	// 		case RoundedRectT:
	// 			rr := comm.rRect
	// 			size := c.camera.GetProjectionSize()
	// 			c.renderer.RoundedRectangleT(rr.x, size.Y()-rr.y, rr.w, rr.h, rr.radius, render.AllRounded, rr.texture, 0, 1, rr.clr)
	// 		case RoundedRect:
	// 			rr := comm.rRect
	// 			size := c.camera.GetProjectionSize()

	// 			c.renderer.RoundedRectangleR(rr.x, size.Y()-rr.y, rr.w, rr.h, rr.radius, render.AllRounded, rr.clr)
	// 		case WindowStartCmd:
	// 			wnd := comm.window
	// 			size := c.camera.GetProjectionSize()

	// 			c.renderer.RoundedRectangleR(wnd.x, size.Y()-wnd.y, wnd.w, wnd.h, 10, render.AllRounded, comm.window.clr)
	// 			c.renderer.RoundedRectangleR(wnd.x, size.Y()-wnd.y, wnd.w, wnd.toolbar.h, 10, render.TopRect, comm.window.toolbar.clr)

	// 			// c.renderer.RectangleR(wnd.x, size.Y()-wnd.y, wnd.w, wnd.h,  comm.window.clr)
	// 			// c.renderer.RectangleR(wnd.x, size.Y()-wnd.y, wnd.w, wnd.toolbar.h, comm.window.toolbar.clr)
	// 		default:
	// 		}
	// 	}
	// 	v.rq.clearCommands()
	// }

	for _, v := range c.sortedWindows {
		size := c.camera.GetProjectionSize()
		var x, y, w, h int32
		x = int32(v.x)
		y = int32(v.y)
		w = int32(v.w)
		h = int32(v.h)

		if int32(size.Y())-(int32(v.y)+int32(v.h)) <= 0 {
			y = 0
			h = int32(size[1]) - int32(v.y)
		} else {
			y = int32(size.Y()) - (int32(v.y) + int32(v.h))

		}
		c.renderer.Scissor(x, y, w, h)
		c.renderer.Draw(c.camera, *v.buffer)
		v.buffer.Clear()
	}

	c.renderer.End()

	c.currentWindow = 0

	if !c.io.IsDragging && c.wantResizeH == true {
		c.wantResizeH = false

	} else if !c.io.IsDragging && c.wantResizeV == true {
		c.wantResizeV = false
	}
	// c.sortedWindows = []*Window{}

	c.io.ScrollX = 0
	c.io.ScrollY = 0

	c.ActiveWidget = ""
}

type UiRenderer interface {
	NewFrame()
	Scissor(x, y, w, h int32)
	Rectangle(x, y, w, h float32, clr [4]float32)
	RectangleR(x, y, w, h float32, clr [4]float32)
	Trinagle(x0, y0, x1, y1, x2, y2 float32, clr [4]float32)
	Circle(x, y, radius float32, steps int, clr [4]float32)
	Line(x0, y0, x1, y1 float32, thick int, clr [4]float32)
	RoundedRectangle(x, y, w, h float32, radius int, clr [4]float32)
	RoundedRectangleR(x, y, w, h float32, radius int, shape render.RoundedRectShape, clr [4]float32)
	RectangleT(x, y, w, h float32, tex *gogl.Texture, uv1, uv0, f float32, clr [4]float32)
	RoundedRectangleT(x, y, w, h float32, radius int, shape render.RoundedRectShape, tex *gogl.Texture, uv1, uv0 float32, clr [4]float32)
	Draw(camera *gogl.Camera, buffer draw.CmdBuffer)
	End()
}
