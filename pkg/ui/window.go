package ui

import (
	// "fmt"
	"fmt"
	"strings"
	"time"

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

	outerRect  utils.Rect
	minW, minH float32

	//render
	buffer *draw.CmdBuffer

	mainWidgetSpace    *WidgetSpace
	currentWidgetSpace *WidgetSpace
	focusedWidgetSpace *WidgetSpace

	//widgSpaceWantFocus bool

	widgSpaces []*WidgetSpace

	capturedV, capturedH  bool
	capturedWin           bool
	capturedInsideWin     bool
	capturedTextSelection bool

	delayedWidgets []func()
}

func genWindowId() string {
	rand.Seed(time.Now().Unix())
	return fmt.Sprint(rand.Intn(100000))
}

func NewWindow(x, y, w, h float32) *Window {
	tb := NewToolbar(x, y, w, 30)
	id := genWindowId()

	wnd := Window{
		Id:              id,
		toolbar:         tb,
		x:               x,
		y:               y,
		w:               w,
		h:               h,
		outerRect:       utils.Rect{Min: utils.Vec2{X: x, Y: y}, Max: utils.Vec2{X: x + w, Y: y + h}},
		minW:            200,
		minH:            50,
		mainWidgetSpace: newWidgetSpace(fmt.Sprintf("main-widg-space-%s", id), x, y+tb.h, w, h-tb.h, Default),
		buffer:          draw.NewBuffer(UiCtx.Io().DisplaySize),
		widgSpaces:      make([]*WidgetSpace, 0),
		delayedWidgets:  []func(){},
	}

	wnd.currentWidgetSpace = wnd.mainWidgetSpace
	return &wnd
}

const (
	defx, defy, defw, defh = 300, 100, 400, 500
	scrollChange           = 2
)

func (c *UiContext) BeginWindow(id string) {
	var wnd *Window
	var cmdw draw.Window_command
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

	// logic
	{
		//Прямоугольник справа
		vResizeRect := utils.NewRect(wnd.x+wnd.w-scrollChange, wnd.y, scrollChange+5, wnd.h)
		hResizeRect := utils.NewRect(wnd.x, wnd.y+wnd.h-scrollChange, wnd.w, scrollChange+5)
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
		if c.ActiveWindow == wnd && wnd.capturedWin && !c.wantResizeV &&
			!c.wantResizeH && !wnd.capturedInsideWin && c.SelectableText == nil {
			newX += c.io.MouseDelta.X
			newY += c.io.MouseDelta.Y
		}
	}

	wnd.x = newX
	wnd.y = newY
	wnd.h = newH
	wnd.w = newW

	wnd.outerRect = utils.NewRect(wnd.x, wnd.y, wnd.w-wnd.mainWidgetSpace.verticalScrollbar.w, wnd.h)

	wnd.mainWidgetSpace.X = newX
	wnd.mainWidgetSpace.Y = newY + wnd.toolbar.h
	wnd.mainWidgetSpace.W = newW
	wnd.mainWidgetSpace.H = newH - wnd.toolbar.h

	{
		cmdw = draw.Window_command{
			Active: wnd.active,
			Id:     wnd.Id,
			X:      wnd.x,
			Y:      wnd.y,
			H:      wnd.h,
			W:      wnd.w,
			Clr:    mainWindowClr2,
			Toolbar: draw.Toolbar_command{
				X:   newX,
				Y:   newY,
				W:   newW,
				H:   wnd.toolbar.h,
				Clr: wnd.toolbar.clr,
			},
		}
	}
	wnd.mainWidgetSpace.cursorX = wnd.mainWidgetSpace.X + UiCtx.CurrentStyle.LeftMargin
	wnd.mainWidgetSpace.cursorY = wnd.mainWidgetSpace.Y + UiCtx.CurrentStyle.TopMargin

	wnd.mainWidgetSpace.ClipRect = [4]float32{wnd.mainWidgetSpace.X, wnd.mainWidgetSpace.Y,
		wnd.mainWidgetSpace.W - wnd.mainWidgetSpace.verticalScrollbar.w, wnd.mainWidgetSpace.H}

	wnd.buffer.CreateWindow(cmdw, draw.NewClip(draw.EmptyClip, [4]float32{wnd.x, wnd.y, wnd.w, wnd.h}))
	wnd.widgetSpaceLogic(wnd.mainWidgetSpace, func() draw.ClipRectCompose {
		cl := [4]float32{wnd.x, wnd.y, wnd.w, wnd.h}
		return draw.NewClip(draw.EmptyClip, cl)
	})
	c.windowStack.Push(wnd)
}

var step float32 = 100

func (wnd *Window) widgetSpaceLogic(ws *WidgetSpace, clip func() draw.ClipRectCompose) {
	c := UiCtx

	if c.hoverBehavior(wnd, utils.NewRectS(ws.ClipRect)) {
		c.ActiveWidgetSpaceId = ws.id
		if ws.flags&Scrollable != 0 {
			c.WantScrollFocusWidgetSpaceId = ws.id
		}
	}
	// Scrollbar behavior
	if ws.flags&Scrollable != 0 {
		ws.vertScrollBar()
		if c.ActiveWindow == wnd && c.HoveredWindow == wnd && c.io.ScrollY != 0 && c.WantScrollFocusWidgetSpaceLastId == ws.id && c.FocusedWidgetSpace == nil {
			ws.handleMouseScroll(float32(c.io.ScrollY))
		}
		if ws.flags&ShowScrollbar != 0 && ws.isVertScrollShown {
			scrlClip := clip()
			scrl := ws.verticalScrollbar

			wnd.buffer.CreateRect(scrl.x, scrl.y, scrl.w, scrl.h, 5, draw.AllRounded, 0, scrl.clr, scrlClip)
			wnd.buffer.CreateRect(scrl.bX, scrl.bY, scrl.bW, scrl.bH, 5, draw.AllRounded, 0, [4]float32{255, 0, 0, 1}, scrlClip)
			wnd.buffer.SeparateBuffer(0, scrlClip)
		}
	}

}

var (
	whiteColor     = [4]float32{255, 255, 255, 1}
	softGreen      = [4]float32{231, 240, 162, 0.8}
	black          = [4]float32{0, 0, 0, 1}
	red            = [4]float32{255, 0, 0, 1}
	transparent    = [4]float32{0, 0, 0, 0}
	mainWindowClr  = [4]float32{231, 158, 162, 0.8}
	mainWindowClr2 = [4]float32{29, 29, 29, 1}
)

func (c *UiContext) getWidget(id string, f func() widgets.Widget) widgets.Widget {
	var widg widgets.Widget
	widg, ok := c.GetWidget(id)
	if !ok {
		widg = f()
		c.AddWidget(id, widg)
	}
	return widg
}

func (c *UiContext) ButtonT(id string, msg string) bool {
	wnd := c.windowStack.Peek()
	var tBtn *widgets.TextButton
	var hovered, clicked bool
	x, y, isRow := wnd.currentWidgetSpace.getCursorPosition()

	tBtn, hovered, clicked = c.textButton(id, wnd, msg, x, y, widgets.Center)
	if hovered {
		tBtn.Button.SetColor(c.CurrentStyle.BtnHoveredColor)
	} else if tBtn.Active() {
		tBtn.Button.SetColor(c.CurrentStyle.BtnActiveColor)
	} else {
		tBtn.Button.SetColor(c.CurrentStyle.BtnColor)
	}

	clip := wnd.endWidget(x, y, isRow, tBtn)
	wnd.buffer.CreateButtonT(x, y, tBtn, *c.font, clip)

	return clicked
}

func (c *UiContext) Slider(id string, i *float32, min, max float32) {
	wnd := c.windowStack.Peek()
	var slider *widgets.Slider
	//var hovered, clicked bool
	x, y, isRow := wnd.currentWidgetSpace.getCursorPosition()

	slider = c.getWidget(id, func() widgets.Widget {
		return widgets.NewSlider(id, x, y, 200, 50, min, max, c.CurrentStyle)
	}).(*widgets.Slider)

	// logic
	{
		slider.HandleMouseDrag(c.io.MouseDelta.X, i, c.dragBehaviorInWindow)
		slider.CalculateNumber(i)
		// In the first launch if number more or less than borders values, we have to make it equal one of them
		if *i > max {
			*i = max
		} else if *i < min {
			*i = min
		}
	}

	clip := draw.NewClip(slider.BoundingBox(), wnd.currentWidgetSpace.ClipRect)

	wnd.endWidget(x, y, isRow, slider)

	wnd.buffer.CreateRect(slider.MainSliderPos()[0], slider.MainSliderPos()[1], slider.MainSliderPos()[2], slider.MainSliderPos()[3], 0,
		draw.StraightCorners, 0, softGreen,
		clip)

	wnd.buffer.CreateRect(slider.BtnSliderPos()[0], slider.BtnSliderPos()[1], slider.BtnSliderPos()[2], slider.BtnSliderPos()[3], 0,
		draw.StraightCorners, 0, black,
		clip)
}

func (c *UiContext) GetTextInfo(x, w float32, msg string) string {

	//splitted := strings.Split(msg, " ")
	r := []rune(msg)
	sb := strings.Builder{}
	sb.Grow(len(r))
	var dw float32 = 0
	for _, l := range r {
		char := c.font.GetCharacter(l)
		if x+dw+float32(char.Advance) < x+w {
			dw += float32(char.Advance)
			sb.Write([]byte(string(l)))
		} else {
			dw = 0
			sb.Write([]byte(string('\n')))
			sb.Write([]byte(string(l)))
		}
	}
	return sb.String()
}

func (c *UiContext) tHelper(id string, x, y, w float32, msg string, key GuiKey, flag widgets.TextFlag) (txt *widgets.Text, hovered bool) {
	wnd := c.windowStack.Peek()
	txt = c.getWidget(id, func() widgets.Widget {
		if flag&widgets.FitContent != 0 {
			msg = c.GetTextInfo(x, w, msg)
		}
		w, h, p := c.font.CalculateTextBounds(msg, c.CurrentStyle.FontScale)
		return widgets.NewText(id, msg, x, y, w, h, p, c.CurrentStyle, flag)
	}).(*widgets.Text)
	if msg != txt.Message && msg != "" {
		if flag&widgets.Editable != 0 {
			if key == GuiKey_Backspace {
				txt.Message = txt.Message[:len(txt.Message)-1]
			} else {
				txt.Message += msg
			}

		} else {
			txt.Message = msg
		}
		if flag&widgets.FitContent != 0 {
			txt.Message = c.GetTextInfo(x, w, msg)
		}
		w, h, p := c.font.CalculateTextBounds(txt.Message, c.CurrentStyle.FontScale)
		txt.Chars = p
		txt.SetWH(w, h)
	}

	hovered = c.hoverBehavior(wnd, utils.NewRectS(txt.BoundingBox()))
	if hovered && txt.Flag&widgets.Selectable != 0 {
		c.SelectableText = txt
	}
	if c.SelectableText == txt {
		c.dragBehavior(wnd.outerRect, &wnd.capturedTextSelection)
		f, w, msg := txt.FindSelectedString(c.io.dragStarted.X-x, c.io.dragDelta.X)
		if w != 0 {
			txt.LastSelectedWidth = w
			txt.LastSelectedX = f
			c.SelectedText = msg
		}
		wnd.buffer.CreateRect(x+txt.LastSelectedX, y, txt.LastSelectedWidth, txt.Height(), 0,
			draw.StraightCorners, 0, softGreen, wnd.DefaultClip())
	}
	return
}
func (c *UiContext) inputTextHelper(id string, x, y float32, msg string, key GuiKey, flag widgets.TextFlag) (txt *widgets.Text, hovered bool) {
	return c.tHelper(id, x, y, 0, msg, key, flag)
}
func (c *UiContext) textHelper(id string, x, y, w float32, msg string, flag widgets.TextFlag) (txt *widgets.Text, hovered bool) {
	return c.tHelper(id, x, y, w, msg, GuiKey_None, flag)
}

func (c *UiContext) getTextInput() (string, GuiKey) {
	k := ""
	key := GuiKey_None
	if c.io.KeyPressedThisFrame {
		key = c.io.PressedKey
		k = c.io.keyToString(key)
	}
	return k, key
}

func (c *UiContext) InputText(id string, size int) {
	wnd := c.windowStack.Peek()
	x, y, isRow := wnd.currentWidgetSpace.getCursorPosition()
	msg, key := c.getTextInput()
	txt, hovered := c.inputTextHelper(id, x, y, msg, key, widgets.Editable)
	y += wnd.currentWidgetSpace.resolveRowAlign(txt.Height())

	if hovered {
		txt.SetBackGroundColor(softGreen)
		txt.CurrentColor = [4]float32{167, 200, 100, 1}
	} else {
		txt.CurrentColor = whiteColor
		txt.SetBackGroundColor(transparent)
	}

	//txt.CurrentColor = [4]float32{255, 255, 255, 1}

	clip := wnd.endWidget(x, y, isRow, txt)

	wnd.buffer.CreateText(x, y, txt, *c.font, clip)
}

func (c *UiContext) TextFitted(id string, w float32, msg string) {
	wnd := c.windowStack.Peek()
	x, y, isRow := wnd.currentWidgetSpace.getCursorPosition()
	txt, hovered := c.textHelper(id, x, y, w, msg, widgets.FitContent)
	y += wnd.currentWidgetSpace.resolveRowAlign(txt.Height())
	wnd.debugDrawS(txt.BoundingBox())
	//wnd.buffer.RoundedBorderRectangle(x, y, txt.Width(), txt.Height(), 30, 15, red, wnd.DefaultClip())
	if hovered {

	} else {
		txt.CurrentColor = whiteColor
		txt.SetBackGroundColor(transparent)
	}

	clip := wnd.endWidget(x, y, isRow, txt)

	wnd.buffer.CreateText(x, y, txt, *c.font, clip)
}

func (c *UiContext) Text(id string, msg string, size int) {
	wnd := c.windowStack.Peek()
	x, y, isRow := wnd.currentWidgetSpace.getCursorPosition()
	txt, hovered := c.textHelper(id, x, y, 0, msg, widgets.Selectable)
	y += wnd.currentWidgetSpace.resolveRowAlign(txt.Height())
	//wnd.buffer.RoundedBorderRectangle(x, y, txt.Width(), txt.Height(), 30, 15, red, wnd.DefaultClip())
	if hovered {
		//txt.SetBackGroundColor(softGreen)
		//txt.CurrentColor = [4]float32{167, 200, 100, 1}
	} else {
		txt.CurrentColor = whiteColor
		txt.SetBackGroundColor(transparent)
	}

	//txt.CurrentColor = [4]float32{255, 255, 255, 1}

	clip := wnd.endWidget(x, y, isRow, txt)

	wnd.buffer.CreateText(x, y, txt, *c.font, clip)
}

func (wnd *Window) DefaultClip() draw.ClipRectCompose {
	return draw.NewClip(wnd.currentWidgetSpace.ClipRect, wnd.mainWidgetSpace.ClipRect)
}

func (c *UiContext) imageHelper(id string, x, y, w, h float32) (img *widgets.Image, hovered, clicked bool) {
	wnd := c.windowStack.Peek()

	img = c.getWidget(id, func() widgets.Widget {
		img = widgets.NewImage(id, x, y, w, h, whiteColor)
		return img
	}).(*widgets.Image)
	{
		hovered := c.hoverBehavior(wnd, utils.NewRectS(img.BoundingBox()))
		if hovered {
			c.setActiveWidget(img.WidgetId())
		}
		clicked = c.io.MouseClicked[0] && hovered
	}
	return
}

func (c *UiContext) Image(id string, w, h float32, tex *gogl.Texture) bool {
	if tex == nil {
		fmt.Println("error")
		return false
	}
	wnd := c.windowStack.Peek()

	x, y, isRow := wnd.currentWidgetSpace.getCursorPosition()
	img, _, clicked := c.imageHelper(id, x, y, w, h)

	clr := img.Color()

	clip := wnd.endWidget(x, y, isRow, img)
	wnd.buffer.CreateTexturedRect(x, y, img.Width(), img.Height(), tex.TextureId, tex.TexCoords, clr, clip)
	//wnd.buffer.CreateRect(x, y, img.Width(), img.Height(), 0, draw.StraightCorners, tex.TextureId, clr, clip)
	return clicked
}

// TODO: measure performance
func (wnd *Window) endWidget(xPos, yPos float32, isRow bool, w widgets.Widget) draw.ClipRectCompose {
	w.UpdatePosition([4]float32{xPos, yPos, w.Width(), w.Height()})
	wnd.addCursor(w.Width(), w.Height())
	if !isRow {
		wnd.currentWidgetSpace.AddVirtualWH(w.Width(), w.Height())
	}

	//wnd.debugDraw(xPos, yPos, w.Width(), w.Height())

	var clip draw.ClipRectCompose
	if wnd.currentWidgetSpace.flags&IgnoreClipping != 0 {
		clip = draw.NewClip(draw.EmptyClip, wnd.currentWidgetSpace.ClipRect)
	} else {
		clip = wnd.DefaultClip()
	}
	return clip
}

func (w *Window) addDelayedWidget(f func()) {
	w.delayedWidgets = append(w.delayedWidgets, f)
}

func (c *UiContext) VSpace(id string) {
	wnd := c.windowStack.Peek()
	var s *widgets.VSpace
	x, y, isRow := wnd.currentWidgetSpace.getCursorPosition()

	s = c.getWidget(id, func() widgets.Widget {
		s := widgets.NewVertSpace(id, [4]float32{x, y, 100, 20})
		return s
	}).(*widgets.VSpace)

	wnd.endWidget(x, y, isRow, s)
}

func (c *UiContext) hoverBehavior(wnd *Window, rect utils.Rect) bool {
	inRect := utils.PointInRect(c.io.MousePos, utils.NewRect(rect.Min.X, rect.Min.Y, rect.Width(), rect.Height()))
	inWindow := RegionHit(c.io.MousePos.X, c.io.MousePos.Y, wnd.x, wnd.y+wnd.toolbar.h, wnd.w, wnd.h-wnd.toolbar.h)
	focusedWidgSpace := false
	// Accept widget actions only from focused widget space
	if c.FocusedWidgetSpace != nil {
		if wnd.currentWidgetSpace != c.FocusedWidgetSpace {
			focusedWidgSpace = true
		}
	}
	return c.ActiveWindow == wnd && inRect && inWindow && !focusedWidgSpace
}

func (c *UiContext) TreeNode(id string, msg string, widgFunc func()) bool {
	wnd := c.windowStack.Peek()
	var tBtn *widgets.TextButton
	var _, hovered bool
	x, y, isRow := wnd.currentWidgetSpace.getCursorPosition()

	tBtn = c.getWidget(id, func() widgets.Widget {
		w, h, p := c.font.CalculateTextBounds(msg, c.CurrentStyle.FontScale)
		return widgets.NewTextButton(id, x, y, w, h, msg, p, widgets.Left, widgets.AllPadding, c.CurrentStyle)
	}).(*widgets.TextButton)

	// logic
	{
		hovered = c.hoverBehavior(wnd, utils.NewRectS(tBtn.BoundingBox()))
		if hovered {
			c.setActiveWidget(tBtn.Id)
			tBtn.Button.SetColor(c.CurrentStyle.BtnHoveredColor)
			if c.io.MouseClicked[0] {
				tBtn.ChangeActive()
			}
		} else if tBtn.Active() {
			tBtn.Button.SetColor(c.CurrentStyle.BtnActiveColor)
		} else {
			tBtn.Button.SetColor(c.CurrentStyle.BtnColor)
		}
		//clicked = c.io.MouseClicked[0] && hovered
	}
	//
	tBtn.SetWidth(wnd.w)
	clip := wnd.endWidget(x, y, isRow, tBtn)
	wnd.buffer.CreateButtonT(x, y, tBtn, *c.font, clip)

	if tBtn.Active() {
		x += 50
		ws := c.subWidgetSpaceHelper(id, x, y+tBtn.Height(), 0, 0, NotScrollable|Resizable, widgFunc)
		wnd.currentWidgetSpace.AddVirtualHeight(ws.H)
		wnd.addCursor(ws.W, ws.H)
	}

	return tBtn.Active()
}

func (c *UiContext) textButton(id string, wnd *Window, msg string, x, y float32, align widgets.TextAlign) (tBtn *widgets.TextButton, hovered, clicked bool) {
	tBtn = c.getWidget(id, func() widgets.Widget {
		w, h, p := c.font.CalculateTextBounds(msg, c.CurrentStyle.FontScale)
		return widgets.NewTextButton(id, x, y, w, h, msg, p, align, widgets.AllPadding, c.CurrentStyle)
	}).(*widgets.TextButton)
	if msg != tBtn.Message {
		tBtn.Message = msg
		w, h, p := c.font.CalculateTextBounds(msg, c.CurrentStyle.FontScale)
		tBtn.Text.Chars = p
		tBtn.SetWH(w, h)
	}
	hovered = c.hoverBehavior(wnd, utils.NewRectS(tBtn.BoundingBox()))
	if hovered {
		c.setActiveWidget(tBtn.Id)
	}
	clicked = c.io.MouseClicked[0] && hovered
	if clicked {
		tBtn.ChangeActive()
	}
	return
}

func (c *UiContext) button(id string, wnd *Window, x, y, w, h float32) (btn *widgets.Button, hovered, clicked bool) {
	btn = c.getWidget(id, func() widgets.Widget {
		return widgets.NewButton(id, x, y, w, h, c.CurrentStyle.BtnColor)
	}).(*widgets.Button)

	hovered = c.hoverBehavior(wnd, utils.NewRectS(btn.BoundingBox()))
	if hovered {
		c.setActiveWidget(btn.WidgetId())
	}
	clicked = c.io.MouseClicked[0] && hovered
	if clicked {
		btn.ChangeActive()
	}

	return
}

func (c *UiContext) Button(id string) bool {
	wnd := c.windowStack.Peek()
	var btn *widgets.Button
	var clicked, hovered bool
	x, y, isRow := wnd.currentWidgetSpace.getCursorPosition()

	btn, hovered, clicked = c.button(id, wnd, x, y, 100, 100)

	if hovered {
		btn.SetColor(c.CurrentStyle.BtnHoveredColor)
	} else if btn.IsActive {
		btn.SetColor(c.CurrentStyle.BtnActiveColor)
	} else {
		btn.SetColor(c.CurrentStyle.BtnColor)
	}

	clip := wnd.endWidget(x, y, isRow, btn)
	wnd.buffer.CreateRect(x, y, btn.Width(), btn.Height(), 0, draw.StraightCorners, 0, btn.Color(), clip)

	return clicked
}

func (wnd *Window) addCursor(width, height float32) {
	row, ok := wnd.currentWidgetSpace.getCurrentRow()
	if !ok {
		wnd.currentWidgetSpace.cursorY += height
	} else {
		if row.RequireColumn {
			row.CursorY += height
			row.UpdateColWidth(width)
			row.AddColHeight(height)
		} else {
			row.CursorX += width
			row.W += width
			row.UpdateHeight(height)
		}
	}
}

func (c *UiContext) setActiveWidget(id string) {
	c.ActiveWidget = id
}

func (c *UiContext) Selection(id string, index *int, data []string, tex *gogl.Texture) {
	wnd := c.windowStack.Peek()
	var s *widgets.Selection

	// Need to use WS because text may not fit into button, so it should be clipped
	c.SubWidgetSpace(id+"---", 0, 0, Resizable|NotScrollable, func() {
		c.Row(id+"row--", func() {
			x, y, isRow := wnd.currentWidgetSpace.getCursorPosition()
			s = c.getWidget(id, func() widgets.Widget {
				return widgets.NewSelection(id, x, y, 300, 40)
			}).(*widgets.Selection)

			wnd.buffer.CreateRect(x, y, s.Width(), s.Height(), 0, draw.StraightCorners, 0, whiteColor, wnd.DefaultClip())
			wnd.buffer.SeparateBuffer(0, wnd.DefaultClip())
			wnd.endWidget(x, y, isRow, s)

			x2, y2, isRow2 := wnd.currentWidgetSpace.getCursorPosition()
			img, _, clicked := c.imageHelper(id+"arrow", x2, y2, s.Height(), s.Height())

			wnd.endWidget(x2-s.Height(), y2, isRow2, img)

			txt, _ := c.textHelper(data[*index]+"--"+id, x, y, 0, data[*index], widgets.Default)
			wnd.buffer.CreateText(x+c.CurrentStyle.Padding, y+(s.Height()-txt.Height())/2, txt,
				*c.font, draw.NewClip(draw.EmptyClip, [4]float32{x, y, s.Width(), s.Height()}))
			wnd.buffer.CreateTexturedRect(x2-s.Height(), y2, img.Width(), img.Height(), tex.TextureId, tex.TexCoords, img.Color(), wnd.DefaultClip())
			if clicked {
				s.Opened = true
				c.setActiveWidget(id)
			}
			if c.ActiveWidget != id {
				s.Opened = false
			}
		})
	})

	c.ContextMenu(id, IgnoreClipping, func() {
		for i, datum := range data {
			x, y, _ := wnd.currentWidgetSpace.getCursorPosition()
			tbt, _, clicked := c.textButton(datum+"_btnT_"+id, wnd, datum, x, y, widgets.Left)
			if clicked {
				*index = i
				c.FocusedWidgetSpace = nil
			}
			tbt.SetWidth(s.Width())
			clip := wnd.endWidget(x, y, false, tbt)
			wnd.buffer.CreateButtonT(x, y, tbt, *c.font, clip)
		}
	})
}

func (c *UiContext) ContextMenu(ownerWidgetId string, flag WidgetSpaceFlag, widgFunc func()) {
	wnd := c.windowStack.Peek()
	var bb [4]float32
	widg, ok := c.GetWidget(ownerWidgetId)
	if !ok {
		return
	}
	bb = widg.BoundingBox()
	id := ownerWidgetId + "-ws-context"
	ws := c.getWidgetSpace(id, 0, 0, wnd, Resizable|FitWidth|flag)
	if c.LastActiveWidget == widg.WidgetId() {
		c.FocusedWidgetSpace = ws
	}
	if c.FocusedWidgetSpace == ws {
		f := func() {
			ws.ClipRect = [4]float32{ws.X, ws.Y, ws.W, ws.H}
			clip := draw.NewClip(draw.EmptyClip, ws.ClipRect)
			wnd.buffer.CreateRect(bb[0], bb[1]+widg.Height(), ws.W, ws.H, 0, draw.StraightCorners, 0, black, clip)
			c.subWidgetSpaceHelper(id, bb[0], bb[1]+widg.Height(), widg.Width(), 0, Resizable|FitWidth|flag, widgFunc)
		}
		wnd.addDelayedWidget(f)
	}
}
func (c *UiContext) Tooltip(id string, widgFunc func()) {
	x, y := c.io.MousePos.X+10, c.io.MousePos.Y+5
	wnd := c.windowStack.Peek()

	ws := c.getWidgetSpace(id, 0, 0, wnd, Resizable|IgnoreClipping)

	wnd.addDelayedWidget(func() {
		wnd.buffer.CreateRect(x, y, ws.W, ws.H, 0, draw.StraightCorners, 0, black, draw.NewClip(draw.EmptyClip, ws.ClipRect))
		c.subWidgetSpaceHelper(id, x, y, 0, 0, Resizable|IgnoreClipping, widgFunc)
	})
}

func (c *UiContext) Column(id string, widgFunc func()) {
	wnd := c.windowStack.Peek()
	var hl *widgets.HybridLayout
	hl, ok := wnd.currentWidgetSpace.getCurrentRow()
	if !ok {
		return
	}
	hl.RequireColumn = true
	hl.CurrentColH, hl.CurrentColW = 0, 0

	widgFunc()

	hl.RequireColumn = false
	hl.CursorY = hl.InitY

	hl.W += hl.CurrentColW
	hl.CursorX += hl.CurrentColW
	hl.UpdateHeight(hl.CurrentColH)
}

func (wnd *Window) debugDrawS(x [4]float32) {
	wnd.buffer.CreateBorderBox(x[0], x[1], x[2], x[3], 2, red)
}
func (wnd *Window) debugDraw(x, y, w, h float32) {
	wnd.buffer.CreateBorderBox(x, y, w, h, 2, red)
}

func (c *UiContext) getWidgetSpace(id string, width, height float32, wnd *Window, flags WidgetSpaceFlag) *WidgetSpace {
	ws, ok := c.widgSpaceCache.Get(id)
	if !ok {
		ws = newWidgetSpace(id, 0, 0, width, height, flags)
		c.widgSpaceCache.Add(id, ws)
	}
	wnd.widgSpaces = append(wnd.widgSpaces, ws)
	return ws
}

func (c *UiContext) subWidgetSpaceHelper(id string, x, y, width, height float32, flags WidgetSpaceFlag, widgFunc func()) *WidgetSpace {
	wnd := c.windowStack.Peek()

	ws := c.getWidgetSpace(id, width, height, wnd, flags)

	var prevWS = wnd.currentWidgetSpace
	wnd.currentWidgetSpace = ws

	ws.X = x
	ws.Y = y
	ws.cursorY = y
	ws.cursorX = x

	outOfWindow := false
	if y < wnd.mainWidgetSpace.Y {
		outOfWindow = true
		// vs-clip-1.png
		if ws.isVertScrollShown {
			ws.ClipRect = [4]float32{x, wnd.mainWidgetSpace.Y, ws.W - ws.verticalScrollbar.w, ws.H - (wnd.mainWidgetSpace.Y - y)}
		} else {
			//ws.ClipRect = [4]float32{x, wnd.mainWidgetSpace.Y, ws.W, ws.H - (wnd.mainWidgetSpace.Y - y)}
			ws.ClipRect = [4]float32{x, wnd.mainWidgetSpace.Y, ws.W, ws.H}
		}
		//wnd.debugDrawS(ws.ClipRect)
	} else {
		if ws.flags&ShowScrollbar != 0 {
			ws.ClipRect = [4]float32{x, y, ws.W - ws.verticalScrollbar.w, ws.H}
		} else {
			ws.ClipRect = [4]float32{x, y, ws.W, ws.H}
		}
	}
	if flags&FitWidth != 0 {
		ws.ClipRect[2] = width
	}

	wnd.widgetSpaceLogic(ws, func() draw.ClipRectCompose {
		cl := [4]float32{ws.X, ws.Y, ws.W, ws.H}
		if outOfWindow {
			cl[1] = wnd.mainWidgetSpace.Y
		}
		return draw.NewClip(cl, wnd.mainWidgetSpace.ClipRect)
	})

	widgFunc()
	ws.checkVerScroll()

	ws.lastVirtualHeight = ws.virtualHeight
	ws.virtualHeight = 0
	ws.lastVirtualWidth = ws.virtualWidth
	ws.virtualWidth = 0

	if ws.flags&Resizable != 0 {
		ws.H = ws.lastVirtualHeight
		ws.W = ws.lastVirtualWidth
	}

	//wnd.buffer.CreateRect(wnd.mainWidgetSpace.X, ws.H+y, wnd.w, 2,
	//	0, draw.StraightCorners, 0, c.CurrentStyle.WidgSpaceDividerColor, wnd.DefaultClip())
	//wnd.debugDraw(x, y, ws.W, ws.H)

	wnd.currentWidgetSpace = prevWS
	return ws
}

func (c *UiContext) SubWidgetSpace(id string, width, height float32, flags WidgetSpaceFlag, widgFunc func()) {
	wnd := c.windowStack.Peek()
	var ws *WidgetSpace

	x, y, _ := wnd.currentWidgetSpace.getCursorPosition()
	ws = c.subWidgetSpaceHelper(id, x, y, width, height, flags, widgFunc)

	wnd.currentWidgetSpace.AddVirtualHeight(ws.H)
	wnd.addCursor(ws.W, ws.H)
}

func (c *UiContext) TabItem(name string, widgFunc func()) {
	wnd := c.windowStack.Peek()
	var tb *widgets.TabBar
	tb, ok := wnd.currentWidgetSpace.getCurrentTabBar()
	x, y, _ := wnd.currentWidgetSpace.getCursorPosition()
	if !ok {
		return
	}
	wspId := name + "-wsp-" + tb.WidgetId()
	_, index := tb.FindTabItem(name, wspId)

	var ws *WidgetSpace
	if index == tb.CurrentTab {
		ws = c.subWidgetSpaceHelper(wspId, x, y, 0, 0, Resizable|HideScrollbar, widgFunc)
	}
	if ws != nil {
		tb.SetHeight(ws.H)
		tb.SetWidth(ws.W)
	}
}
func (c *UiContext) TabBar(id string, widgFunc func()) {
	wnd := c.windowStack.Peek()
	var tab *widgets.TabBar
	x, y, _ := wnd.currentWidgetSpace.getCursorPosition()

	tab = c.getWidget(id, func() widgets.Widget {
		return widgets.NewTabBar(id, x, y, 0, 0)
	}).(*widgets.TabBar)

	var rowHeight, rowWidth float32
	ws := c.subWidgetSpaceHelper(id, x, y, 0, 0, Resizable|NotScrollable, func() {
		//cr := wnd.currentWidgetSpace
		//wnd.buffer.CreateRect(cr.X, cr.Y, cr.W, cr.H, 10, draw.AllRounded, 0, softGreen, wnd.DefaultClip())
		c.Row("rowds", func() {
			row, _ := wnd.currentWidgetSpace.getCurrentRow()
			wnd.buffer.CreateRect(row.X, row.Y, row.Width(), row.Height(), 10, draw.AllRounded, 0, transparent, wnd.DefaultClip())
			for i, item := range tab.Bars {
				x, y, _ := wnd.currentWidgetSpace.getCursorPosition()
				tbtn, hovered, clicked := c.textButton(fmt.Sprint(id, "-", i), wnd, item.Name, x, y, widgets.Center)
				if clicked {
					tab.CurrentTab = i
				}
				if hovered {
					tbtn.SetBackgroundColor(c.CurrentStyle.TabBtnActiveColor)
					if clicked {
						tab.ChangeActive(item)
					}
				} else if item.Active {
					tbtn.SetBackgroundColor(c.CurrentStyle.TabBtnActiveColor)
				} else {
					tbtn.SetBackgroundColor(c.CurrentStyle.TabBtnColor)
				}

				tbtn.Text.CurrentColor = whiteColor
				tbtn.SetHeight(tbtn.Height() - (tbtn.Height() - tbtn.Text.Height()) + c.CurrentStyle.Padding)
				clip := wnd.endWidget(x, y, false, tbtn)
				//wnd.buffer.CreateRect(x, y, tbtn.Width(), tbtn.Height(), 10, draw.TopRect, 0, tbtn.Color(), clip)
				wnd.buffer.CreateRect(x, y, tbtn.Width(), tbtn.Height(), 10, draw.TopRect, 0, tbtn.Color(), clip)
				wnd.buffer.SeparateBuffer(0, clip)
				wnd.buffer.CreateText(tbtn.Text.BoundingBox()[0], tbtn.Text.BoundingBox()[1], tbtn.Text, *c.font, clip)

				row, _ := wnd.currentWidgetSpace.getCurrentRow()
				if i != len(tab.Bars)-1 {
					row.CursorX += 10
					row.W += 10
				}
				if row.Height() > rowHeight {
					rowHeight = row.Height()
				}
				rowWidth = row.Width()
			}
		})

		wnd.buffer.CreateRect(wnd.x, y+rowHeight, wnd.w, 2, 0, draw.StraightCorners, 0, c.CurrentStyle.TabBtnActiveColor, draw.NewClip(draw.EmptyClip, wnd.mainWidgetSpace.ClipRect))
		wnd.buffer.SeparateBuffer(0, draw.NewClip(draw.EmptyClip, wnd.mainWidgetSpace.ClipRect))
		wnd.currentWidgetSpace.cursorY += 5

		tab.BarHeight = rowHeight
		tab.SetWidth(rowWidth)
		wnd.currentWidgetSpace.tabStack.Push(tab)
		widgFunc()
		wnd.currentWidgetSpace.tabStack.Pop()
	})
	ws.W = tab.Width()
	ws.H = tab.Height() + tab.BarHeight + 5

	wnd.buffer.CreateRect(wnd.mainWidgetSpace.X, ws.H+y, wnd.w, 2,
		0, draw.StraightCorners, 0, c.CurrentStyle.WidgSpaceDividerColor, wnd.DefaultClip())

	wnd.addCursor(ws.W, ws.H)
	wnd.currentWidgetSpace.AddVirtualWH(ws.W, ws.H)
}

func (c *UiContext) Row(id string, widgFunc func()) {
	wnd := c.windowStack.Peek()
	var row *widgets.HybridLayout
	x, y, _ := wnd.currentWidgetSpace.getCursorPosition()
	// Need to return cursor back, because internal row cursor shouldn't know anything about outer
	y += wnd.currentWidgetSpace.scrlY

	row = c.getWidget(id, func() widgets.Widget {
		return widgets.NewHLayout(id, x, y, widgets.VerticalAlign, c.CurrentStyle)
	}).(*widgets.HybridLayout)
	row.UpdatePosition([4]float32{x, y, row.Width(), row.Height()})
	wnd.currentWidgetSpace.rowStack.Push(row)

	widgFunc()

	hl := wnd.currentWidgetSpace.rowStack.Pop()
	wnd.addCursor(0, hl.H)
	//wnd.endWidget(x, y, false, row)

	wnd.currentWidgetSpace.AddVirtualWH(hl.W, hl.H)
	hl.LastWidth = hl.W
	hl.LastHeight = hl.H
	hl.H = 0
	hl.W = 0
}

func (c *UiContext) EndWindow() {

	wnd := c.windowStack.Peek()

	for _, f := range wnd.delayedWidgets {
		f()
	}
	wnd.delayedWidgets = []func(){}
	wnd = c.windowStack.Pop()
	wnd.mainWidgetSpace.checkVerScroll()
	var clip = draw.NewClip(draw.EmptyClip, wnd.mainWidgetSpace.ClipRect)
	wnd.buffer.SeparateBuffer(0, clip) // Make sure that we didn't miss anything
	wnd.mainWidgetSpace.AddVirtualHeight(c.CurrentStyle.BotMargin)

	wnd.mainWidgetSpace.lastVirtualHeight = wnd.mainWidgetSpace.virtualHeight
	wnd.mainWidgetSpace.virtualHeight = 0
	wnd.mainWidgetSpace.lastVirtualWidth = wnd.mainWidgetSpace.virtualWidth
	wnd.mainWidgetSpace.virtualWidth = 0
}

func RegionHit(mouseX, mouseY, x, y, w, h float32) bool {
	return mouseX >= x && mouseY >= y && mouseX <= x+w && mouseY <= y+h
}
