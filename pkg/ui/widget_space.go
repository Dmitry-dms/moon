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
}

var defScrollWidth float32 = 10
var scrolconst float32 = 0

func (ws *WidgetSpace) vertScrollBar() {
	vB := ws.verticalScrollbar
	vB.x = ws.X + ws.W - vB.w
	vB.y = ws.Y
	vB.h = ws.H

	// top border
	if ws.scrlY < 0 {
		ws.scrlY = 0
	}
}

func (ws *WidgetSpace) setScrollY(scrollY float32) {
	ws.scrlY = scrollY
}

func (ws *WidgetSpace) handleMouseDrag() {
	vB := ws.verticalScrollbar
	var ratio float32 = (ws.H) / (ws.virtualHeight)

	UiCtx.dragBehavior(utils.NewRect(vB.bX, vB.bY, vB.bW, vB.bH), &ws.captured)
	endReached := ws.scrlY+ws.H <= ws.virtualHeight
	if ws.captured && endReached {
		ws.scrlY += UiCtx.io.MouseDelta.Y

	} else if ws.captured && !endReached {
		if UiCtx.io.MouseDelta.Y < 0 {
			ws.scrlY += UiCtx.io.MouseDelta.Y
		}
	}

	vB.bH = ws.H * ratio
	vB.bX = ws.X + ws.W - vB.w
	vB.bY = ws.Y + ws.scrlY*ratio
}
func (ws *WidgetSpace) handleMouseScroll(scrollY float32) {
	var factor float32 = scrollY * float32(step)
	currentPos := ws.scrlY
	var topBorder float32 = 0
	botBorder := ws.virtualHeight
	var ratio float32 = (ws.H) / (ws.virtualHeight)
	if currentPos <= topBorder {
		if factor > 0 {
			ws.scrlY += factor * ratio
		}
	} else if currentPos+ws.H >= botBorder {
		if factor < 0 {
			ws.scrlY += factor * ratio
		}
	} else {
		ws.scrlY += factor * ratio
	}
}

func (ws *WidgetSpace) checkVerScroll() {
	if ws.cursorY > ws.Y+ws.H {
		ws.isVertScrollShown = true
	} else {
		ws.scrlY = 0
		ws.isVertScrollShown = false
		// ws.verticalScrollbar.bY = 0
	}
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
		virtualHeight: 0,
		verticalScrollbar: NewScrolBar(utils.NewRect(x+w-defScrollWidth, y, defScrollWidth, h),
			utils.NewRect(x+w-defScrollWidth, y, defScrollWidth, 50),
			[4]float32{150, 155, 155, 1}),
	}
	return &vs
}

func (ws *WidgetSpace) addWidget(widg widgets.Widget) bool {
	ws.widgets = append(ws.widgets, widg)
	// ws.virtualHeight += widg.Rectangle()[3]
	return UiCtx.AddWidget(widg.GetId(), widg)
}

func (ws *WidgetSpace) AddVirtualHeight(height float32) {
	ws.virtualHeight += height
}
