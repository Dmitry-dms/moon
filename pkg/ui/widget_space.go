package ui

import (
	"github.com/Dmitry-dms/moon/pkg/ui/utils"
	"github.com/Dmitry-dms/moon/pkg/ui/widgets"
)

type WidgetSpace struct {
	X, Y, W, H float32

	//inner widgets
	cursorX, cursorY float32
	widgetCounter    int
	widgets          []widgets.Widget
	virtualHeight    float32

	verticalScrollbar *Scrollbar
	captured          bool

	isVertScrollShown bool
	scrlY             float32

	ratio float32

	rowStack utils.Stack[*widgets.Row]
}

var defScrollWidth float32 = 10
var scrolconst float32 = 0

func (ws *WidgetSpace) vertScrollBar() {
	vB := ws.verticalScrollbar
	vB.x = ws.X + ws.W - vB.w
	vB.y = ws.Y
	vB.h = ws.H
	ws.ratio = (ws.H) / (ws.virtualHeight)
	// top border
	if ws.scrlY < 0 {
		ws.scrlY = 0
	}

	// this is here because we need to update scroll btn position every frame
	ws.handleMouseDrag()
}

func newWidgetSpace(x, y, w, h float32) *WidgetSpace {
	vs := WidgetSpace{
		X:             x,
		Y:             y,
		W:             w,
		H:             h,
		cursorX:       x,
		cursorY:       y,
		widgetCounter: 0,
		widgets:       []widgets.Widget{},
		virtualHeight: 10,
		verticalScrollbar: NewScrolBar(utils.NewRect(x+w-defScrollWidth, y, defScrollWidth, h),
			utils.NewRect(x+w-defScrollWidth, y, defScrollWidth, 50),
			[4]float32{150, 155, 155, 1}),
		rowStack: utils.NewStack[*widgets.Row](),
	}
	return &vs
}

// FIXME: is this right?
func (ws *WidgetSpace) setScrollY(scrollY float32) {
	if scrollY <= ws.virtualHeight {
		ws.scrlY = scrollY * (1 - ws.ratio)
	}
}

func (ws *WidgetSpace) handleMouseDrag() {
	vB := ws.verticalScrollbar

	UiCtx.dragBehavior(utils.NewRect(vB.bX, vB.bY, vB.bW, vB.bH), &ws.captured)
	endReached := ws.scrlY+ws.H <= ws.virtualHeight
	delta := UiCtx.io.MouseDelta.Y
	if ws.captured && endReached {
		// prevent top glitch
		if ws.scrlY <= 0 {
			if delta > 0 {
				ws.scrlY += delta
			}
		} else {
			ws.scrlY += delta
		}
	} else if ws.captured && !endReached {
		if delta < 0 {
			ws.scrlY += delta
		}
	}

	vB.bH = ws.H * ws.ratio
	vB.bX = ws.X + ws.W - vB.w
	vB.bY = ws.Y + ws.scrlY*ws.ratio
}
func (ws *WidgetSpace) handleMouseScroll(scrollY float32) {
	var factor float32 = scrollY * float32(step)
	currentPos := ws.scrlY
	var topBorder float32 = 0
	botBorder := ws.virtualHeight

	if currentPos <= topBorder {
		if factor > 0 {
			ws.scrlY += factor * ws.ratio
		}
	} else if currentPos+ws.H >= botBorder {
		if factor < 0 {
			ws.scrlY += factor * ws.ratio
		}
	} else {
		ws.scrlY += factor * ws.ratio
	}
}

func (ws *WidgetSpace) checkVerScroll() {
	if ws.cursorY > ws.Y+ws.H {
		ws.isVertScrollShown = true
	} else {
		ws.scrlY = 0
		ws.isVertScrollShown = false
	}
}

func (ws *WidgetSpace) addWidget(widg widgets.Widget) bool {
	ws.widgets = append(ws.widgets, widg)
	return UiCtx.AddWidget(widg.WidgetId(), widg)
}

func (ws *WidgetSpace) AddVirtualHeight(height float32) {
	ws.virtualHeight += height
}

func (ws *WidgetSpace) getCurrentRow() (*widgets.Row, bool) {
	if ws.rowStack.Length() == 0 {
		return nil, false
	} else {
		return ws.rowStack.GetTop(), true
	}
}

func (ws *WidgetSpace) getCursorPosition() (x float32, y float32) {
	row, ok := ws.getCurrentRow()
	if !ok {
		x = ws.cursorX
		y = ws.cursorY
	} else {
		x = row.CursorX
		y = ws.cursorY
	}
	return
}
func (ws *WidgetSpace) BeginRow() {

	// ws.rowStack.Push()
}
