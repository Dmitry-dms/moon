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
	// cursorX, cursorY float32
	// widgetCounter    int
	// widgets          []widgets.Widget
	mainWidgetSpace    *WidgetSpace
	currentWidgetSpace *WidgetSpace

	virtualHeight float32 // сумма высот всех виджетов (для скроллинга)

	//scrollbar
	// scrlBar       *Scrollbar
	// isScrollShown bool
	// scrlY         float32

	capturedV, capturedH bool
	capturedWin          bool
}

func NewWindow(x, y, w, h float32) *Window {
	tb := NewToolbar(x, y, w, 30)
	wnd := Window{
		toolbar: tb,
		x:       x,
		y:       y,
		w:       w,
		h:       h,

		outerRect: utils.Rect{Min: utils.Vec2{X: x, Y: y}, Max: utils.Vec2{X: x + w, Y: y + h}},

		minW: 200,
		minH: 50,
		// srcX:      x,
		// srcY:      y + tb.h + UiCtx.CurrentStyle.TopMargin,

		// widgets: []widgets.Widget{},
		mainWidgetSpace: newWidgetSpace(x, y+tb.h, w, h-tb.h),
		// scrlBar:         NewScrolBar(utils.NewRect(x+w-10, y, 20, h), utils.NewRect(x+w-10, y, 10, 50), [4]float32{150, 155, 155, 1}),
		buffer: draw.NewBuffer(UiCtx.camera),
	}
	wnd.currentWidgetSpace = wnd.mainWidgetSpace
	// wnd.scrlY = wnd.toolbar.h
	// wnd.mainWidgetSpace.scrlY = wnd.toolbar.h
	return &wnd
}

func (w *Window) AddCommand(cmd draw.Command) {
	w.buffer.AddCommand(cmd)
}

var counter int = 0

const (
	defx, defy, defw, defh = 300, 100, 400, 500
	scrollChange           = 2
)

// var capturedV, capturedH bool

func (c *UiContext) BeginWindow(id string) {
	var wnd *Window
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
	if c.ActiveWindow == wnd && wnd.capturedWin && !c.wantResizeV && !c.wantResizeH {
		newX += c.io.MouseDelta.X
		newY += c.io.MouseDelta.Y
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

	wnd.mainWidgetSpace.cursorX = wnd.mainWidgetSpace.X + UiCtx.CurrentStyle.LeftMargin
	wnd.mainWidgetSpace.cursorY = wnd.mainWidgetSpace.Y + UiCtx.CurrentStyle.TopMargin
	wnd.AddCommand(cmd)

	// Scrollbar behavior
	if wnd.mainWidgetSpace.isVertScrollShown {
		wnd.mainWidgetSpace.vertScrollBar()
		if c.ActiveWindow == wnd && c.io.ScrollY != 0 && c.ActiveWidget == "" {
			wnd.mainWidgetSpace.handleMouseDrag()
			wnd.mainWidgetSpace.handleMouseScroll(float32(c.io.ScrollY))
		}
		scrollCommand := draw.Rounded_rect{
			X:      wnd.mainWidgetSpace.verticalScrollbar.x,
			Y:      wnd.mainWidgetSpace.verticalScrollbar.y,
			W:      wnd.mainWidgetSpace.verticalScrollbar.w,
			H:      wnd.mainWidgetSpace.verticalScrollbar.h,
			Clr:    wnd.mainWidgetSpace.verticalScrollbar.clr,
			Radius: 5,
		}
		sbCmd := draw.Command{
			RRect: &scrollCommand,
			Type:  draw.RoundedRect,
		}
		wnd.AddCommand(sbCmd)

		scrollBtnCommand := draw.Rounded_rect{
			X:      wnd.mainWidgetSpace.verticalScrollbar.bX,
			Y:      wnd.mainWidgetSpace.verticalScrollbar.bY,
			W:      wnd.mainWidgetSpace.verticalScrollbar.bW,
			H:      wnd.mainWidgetSpace.verticalScrollbar.bH,
			Clr:    [4]float32{255, 0, 0, 1},
			Radius: 5,
		}

		sbtnCmd := draw.Command{
			RRect: &scrollBtnCommand,
			Type:  draw.RoundedRect,
		}
		wnd.AddCommand(sbtnCmd)
	}

	c.windowStack.Push(wnd)

	wnd.buffer.InnerWindowSpace = [4]float32{wnd.mainWidgetSpace.X, wnd.mainWidgetSpace.Y, wnd.mainWidgetSpace.W, wnd.mainWidgetSpace.H}
}

var visibleRatio float32
var step float32 = 40
var r, g, b float32 = 231, 158, 162

func (w *Window) addWidget(widg widgets.Widget) bool {
	return w.currentWidgetSpace.addWidget(widg)
}

var (
	whiteColor = [4]float32{255, 255, 255, 1}
)

func (wnd *Window) getWidget(id string, w widgets.WidgetType) widgets.Widget {
	var widg widgets.Widget

	wi, ok := UiCtx.GetWidget(id)
	if ok {
		widg = wi
	} else {
		x := wnd.currentWidgetSpace.cursorX
		y := wnd.currentWidgetSpace.cursorY
		switch w {
		case widgets.ButtonWidget:
			widg = &widgets.Button{
				Id:           id,
				CurrentColor: [4]float32{67, 86, 205, 1},
				IsActive:     false,
				BoundingBox:  [4]float32{x, y, float32(100), float32(100)},
			}
		case widgets.ImageWidget:
			widg = &widgets.Image{
				Id:           id,
				CurrentColor: whiteColor,
				BoundingBox:  [4]float32{x, y, float32(100), float32(100)},
			}
		case widgets.TextWidget:
			widg = &widgets.Text{
				Id:           id,
				CurrentColor: whiteColor,
			}
		case widgets.VerticalSpacingWidget:
			widg = &widgets.VSpace{
				BoundingBox: [4]float32{x, y, float32(100), float32(20)},
				Id:          id,
			}
		}
		wnd.addWidget(widg)

	}
	return widg
}

// func (c *UiContext) ButtonRR(id string, tex *gogl.Texture) bool {

// 	wnd := c.windowStack.GetTop()
// 	var btn *widgets.Button

// 	btn = wnd.getWidget(id, widgets.ButtonWidget).(*widgets.Button)

// 	clr := btn.CurrentColor

// 	inRect := utils.PointInRect(c.io.MousePos, utils.NewRectS(btn.BoundingBox))

// 	if wnd == c.ActiveWindow && inRect {
// 		c.ActiveWidget = btn.GetId()
// 		// fmt.Println(btn.GetId())
// 		if c.io.MouseClicked[0] {
// 			btn.IsActive = !btn.IsActive
// 		}
// 	}

// 	if btn.IsActive {
// 		btn.SetColor([4]float32{150, btn.CurrentColor[1], btn.CurrentColor[2], btn.CurrentColor[3]})
// 	} else {
// 		btn.SetColor([4]float32{80, clr[1], clr[2], clr[3]})
// 		// btn.SetColor(whiteColor)
// 	}

// 	clicked := c.io.MouseClicked[0] && inRect

// 	//DEBUG
// 	scroll := float32(0)
// 	if c.io.ScrollY != 0 && c.ActiveWidget == btn.Id {
// 		scroll += float32(c.io.ScrollY) * 10
// 		btn.AddWidth(scroll)
// 		btn.AddHeight(scroll)
// 	}

// 	rect := draw.Rounded_rect{
// 		X:      wnd.cursorX,
// 		Y:      wnd.cursorY,
// 		W:      btn.Width(),
// 		H:      btn.Height(),
// 		Clr:    clr,
// 		Radius: 5,
// 	}

// 	btn.BoundingBox = [4]float32{rect.X, rect.Y, rect.W, rect.H}

// 	cmd := draw.Command{
// 		RRect: &rect,
// 		Shown: true,
// 	}
// 	if tex != nil {
// 		rect.Texture = tex
// 		cmd.Type = draw.RoundedRectT
// 		rect.Clr = whiteColor
// 	} else {
// 		cmd.Type = draw.RectType
// 	}

// 	wnd.AddCommand(cmd)

// 	wnd.currentWidgetSpace.cursorY += rect.H
// 	return clicked
// }

var scale float32 = 1

func (c *UiContext) Text(id string, msg string, size int) {
	wnd := c.windowStack.GetTop()
	var widg widgets.Widget
	var txt *widgets.Text
	y := wnd.currentWidgetSpace.cursorY
	x := wnd.currentWidgetSpace.cursorX

	// txt = wnd.getWidget(id, widgets.TextWidget).(*widgets.Text)

	widg, ok := UiCtx.GetWidget(id)
	if !ok {
		s := c.font.CalculateTextBounds(msg, size)
		txt = &widgets.Text{
			Id:           id,
			CurrentColor: whiteColor,
			Message:      msg,
			BoundingBox:  [4]float32{x, y, s[0], s[1]},
		}
		wnd.addWidget(txt)
		wnd.currentWidgetSpace.AddVirtualHeight(float32(int(s[1])))
	} else {
		txt = widg.(*widgets.Text)
	}

	txt.Message = msg
	clr := txt.CurrentColor

	// txt.BoundingBox = [4]float32{x, y, s[0], s[1]}

	// DEBUG
	y -= wnd.currentWidgetSpace.scrlY
	// y -= wnd.currentWidgetSpace.verticalScrollbar.bY
	//
	hovered := c.hoverBehavior(wnd, utils.NewRectS(txt.BoundingBox))
	if hovered {
		c.ActiveWidget = txt.GetId()
		txt.CurrentColor = [4]float32{167, 200, 100, 1}
	} else {
		txt.CurrentColor = whiteColor
	}

	//DEBUG
	// scroll := float32(0)
	// if c.io.ScrollY != 0 {
	// 	scroll += float32(c.io.ScrollY) * 10
	// 	scale = scroll
	// }
	// Create command and append it to slice
	{
		rect := draw.Text_command{
			X:    wnd.currentWidgetSpace.cursorX,
			Y:    y,
			Clr:  clr,
			Text: txt.Message,
			Font: *c.font,
			Size: size,
		}

		cmd := draw.Command{
			Text:     &rect,
			Type:     draw.Text,
			WidgetId: txt.Id,
		}
		wnd.AddCommand(cmd)
	}

	txt.BoundingBox = [4]float32{x, y, txt.Width(), txt.Height()}
	wnd.currentWidgetSpace.cursorY += txt.Height()
	// wnd.currentWidgetSpace.AddVirtualHeight(txt.Height())
}

func (c *UiContext) Image(id string, tex *gogl.Texture) bool {
	if tex == nil {
		fmt.Println("error")
		return false
	}

	wnd := c.windowStack.GetTop()
	var img *widgets.Image
	var widg widgets.Widget
	x := wnd.currentWidgetSpace.cursorX
	y := wnd.currentWidgetSpace.cursorY
	// img = wnd.getWidget(id, widgets.ImageWidget).(*widgets.Image)

	widg, ok := UiCtx.GetWidget(id)
	if !ok {
		img = &widgets.Image{
			Id:           id,
			CurrentColor: whiteColor,
			BoundingBox:  [4]float32{x, y, float32(100), float32(100)},
		}
		wnd.addWidget(img)
		wnd.currentWidgetSpace.AddVirtualHeight(img.Height())
	} else {
		img = widg.(*widgets.Image)
	}

	clr := img.CurrentColor

	hovered := c.hoverBehavior(wnd, utils.NewRectS(img.BoundingBox))

	if hovered {
		c.ActiveWidget = img.GetId()
	}
	clicked := c.io.MouseClicked[0] && hovered

	// DEBUG
	y -= wnd.currentWidgetSpace.scrlY
	// y -= wnd.currentWidgetSpace.verticalScrollbar.bY
	//

	// Create command and append it to slice
	{
		rect := draw.Rect_command{
			X:       wnd.currentWidgetSpace.cursorX,
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

	img.BoundingBox = [4]float32{wnd.currentWidgetSpace.cursorX, y, img.Width(), img.Height()}

	wnd.currentWidgetSpace.cursorY += img.Height()
	// wnd.currentWidgetSpace.AddVirtualHeight(img.Height())
	return clicked
}

func (c *UiContext) VSpace(id string) {

	wnd := c.windowStack.GetTop()
	var s *widgets.VSpace
	var widg widgets.Widget
	x := wnd.currentWidgetSpace.cursorX
	y := wnd.currentWidgetSpace.cursorY
	// img = wnd.getWidget(id, widgets.ImageWidget).(*widgets.Image)

	widg, ok := UiCtx.GetWidget(id)
	if !ok {
		s = &widgets.VSpace{
			BoundingBox: [4]float32{x, y, float32(100), float32(20)},
			Id:          id,
		}
		wnd.addWidget(s)
		wnd.currentWidgetSpace.AddVirtualHeight(s.Height())
	} else {
		s = widg.(*widgets.VSpace)
	}

	// s = wnd.getWidget(id, widgets.VerticalSpacingWidget).(*widgets.VSpace)

	wnd.currentWidgetSpace.cursorY += s.Height()
	// wnd.currentWidgetSpace.AddVirtualHeight(s.Height())
}

func (c *UiContext) ButtonT(id string, text string, size int) bool {
	return c.button(id, text, size)
}

func (c *UiContext) Button(id string) bool {
	return c.button(id, "", 0)
}

func (c *UiContext) hoverBehavior(wnd *Window, rect utils.Rect) bool {
	inRect := utils.PointInRect(c.io.MousePos, utils.NewRect(rect.Min.X, rect.Min.Y, rect.Width(), rect.Height()))
	inWindow := RegionHit(c.io.MousePos.X, c.io.MousePos.Y, wnd.x, wnd.y+wnd.toolbar.h, wnd.w, wnd.h-wnd.toolbar.h)
	return c.ActiveWindow == wnd && inRect && inWindow
}

func (c *UiContext) button(id string, text string, size int) bool {

	wnd := c.windowStack.GetTop()
	var btn *widgets.Button

	var clicked bool
	var cmd draw.Command
	var rect draw.Rect_command

	var widg widgets.Widget

	x := wnd.currentWidgetSpace.cursorX
	y := wnd.currentWidgetSpace.cursorY
	widg, ok := UiCtx.GetWidget(id)
	if !ok {
		btn = &widgets.Button{
			Id:           id,
			CurrentColor: [4]float32{67, 86, 205, 1},
			IsActive:     false,
			BoundingBox:  [4]float32{x, y, float32(100), float32(100)},
		}
		wnd.addWidget(btn)
		wnd.currentWidgetSpace.AddVirtualHeight(btn.Height())
	} else {
		btn = widg.(*widgets.Button)
	}

	// btn = wnd.getWidget(id, widgets.ButtonWidget).(*widgets.Button)

	w := btn.Width()
	h := btn.Height()
	clr := btn.CurrentColor

	// DEBUG
	y -= wnd.currentWidgetSpace.scrlY
	// y -= wnd.currentWidgetSpace.verticalScrollbar.bY
	//

	rect = draw.Rect_command{
		X:   x,
		Y:   y,
		W:   w,
		H:   h,
		Clr: clr,
	}

	{
		hovered := c.hoverBehavior(wnd, utils.NewRect(x, y, w, h))
		if hovered {
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

		clicked = c.io.MouseClicked[0] && hovered

		cmd = draw.Command{
			Rect:  &rect,
			Type:  draw.RectType,
			Shown: true,
		}
	}

	wnd.AddCommand(cmd)
	wnd.currentWidgetSpace.cursorY += btn.Height()
	// wnd.currentWidgetSpace.AddVirtualHeight(btn.Height())
	return clicked
}

func (c *UiContext) EndWindow() {

	wnd := c.windowStack.Pop()

	// if wnd.isScrollShown {
	// 	c.srollbar(wnd)
	// } else {
	// 	wnd.scrlY = 0
	// }

	// Toolbar := draw.Toolbar_command{
	// 	X:   wnd.x,
	// 	Y:   wnd.y,
	// 	W:   wnd.w,
	// 	H:   wnd.toolbar.h,
	// 	Clr: wnd.toolbar.clr,
	// }

	// cmdToolbar := draw.Command{

	// 	Type:    draw.ToolbarCmd,
	// 	Toolbar: &Toolbar,
	// 	// window:   cmdw,
	// }
	// wnd.AddCommand(cmdToolbar)

	// cmd := draw.Command{
	// 	Type: draw.WindowCmd,
	// }

	// if wnd.cursorY > wnd.y+wnd.h {
	// 	wnd.isScrollShown = true
	// } else {
	// 	// fmt.Println(wnd.y+wnd.h,wnd.cursorY)
	// 	wnd.isScrollShown = false
	// }

	// wnd.AddCommand(cmd)

	wnd.mainWidgetSpace.checkVerScroll()
	// wnd.currentWidgetSpace.virtualHeight = 0
	c.currentWindow++
}

func RegionHit(mouseX, mouseY, x, y, w, h float32) bool {
	return mouseX >= x && mouseY >= y && mouseX <= x+w && mouseY <= y+h
}
