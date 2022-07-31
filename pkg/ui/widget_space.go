package ui

import "github.com/Dmitry-dms/moon/pkg/ui/widgets"

type WidgetSpace struct {
	X, Y, W, H float32

	//inner widgets
	cursorX, cursorY float32
	widgetCounter    int
	widgets          []widgets.Widget
	virtualHeight    float32
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
		virtualHeight: h,
	}
	return &vs
}

func (ws *WidgetSpace) addWidget(widg widgets.Widget) bool{
	ws.widgets = append(ws.widgets, widg)
	ws.virtualHeight += widg.Rectangle()[3]
	return UiCtx.AddWidget(widg.GetId(), widg)
}
