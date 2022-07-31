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

	isVertScrollShown bool
	scrlY             float32
}

var defScrollWidth float32 = 10

func (ws *WidgetSpace) vertScrollBar() {
	ws.verticalScrollbar.x = ws.X + ws.W - ws.verticalScrollbar.w
	ws.verticalScrollbar.y = ws.Y
	ws.verticalScrollbar.h = ws.H
	{
		if ws.scrlY < 0 {
			ws.scrlY = 0
		}
		var ratio float32 = (ws.scrlY + ws.H - ws.verticalScrollbar.bH) / ws.virtualHeight

		if ws.scrlY == 0 {
			visibleRatio = 0
		} else {
			visibleRatio = (ratio)*(ws.H) - ws.verticalScrollbar.bH
		}
	}
	ws.verticalScrollbar.bX = ws.X + ws.W - ws.verticalScrollbar.w
	ws.verticalScrollbar.bY = ws.Y + visibleRatio // для правильного отоборажения при ресайзинге

	if ws.verticalScrollbar.bY >= ws.Y+ws.H-ws.verticalScrollbar.bH { // конечное положение для кнопки скроллбара
		ws.verticalScrollbar.bY = ws.Y + ws.H - ws.verticalScrollbar.bH
	}
}

func (ws *WidgetSpace) updatePosition(scrollY float32) {
	var factor float32 = scrollY * float32(step)
	currentPos := ws.Y + ws.scrlY
	topBorder := ws.Y
	botBorder := ws.Y + ws.virtualHeight + ws.verticalScrollbar.bH
	if currentPos <= topBorder {
		if factor > 0 {
			ws.scrlY += step
		}
	} else if currentPos+ws.H >= botBorder {
		if factor < 0 {
			ws.scrlY -= step
		} else {
			v := currentPos+ws.H - botBorder
			if v <=step {
				ws.scrlY += v
			}
		}
	} else {
		ws.scrlY += factor
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
	ws.virtualHeight += widg.Rectangle()[3]
	return UiCtx.AddWidget(widg.GetId(), widg)
}
