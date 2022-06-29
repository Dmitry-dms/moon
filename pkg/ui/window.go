package ui

type Window struct {
	toolbar       Toolbar
	xPos, yPos    float32 // top-left corner
	width, height float32
}

func NewWindow(x, y, w, h float32) Window {
	tb := NewToolbar(x, y, w, 30)
	wnd := Window{
		toolbar: tb,
		xPos:    x,
		yPos:    y - tb.height,
		width:   w,
		height:  h,
	}
	return wnd
}

func clrConverter(clr float32) float32 {
	return clr * 100 / 255
}

func (c *UiContext) AddWindow() {
	window := NewWindow(300, 100, 400, 500)
	var r float32 = 231
	var g float32 = 158
	var b float32 = 162

	cl := [4]float32{r, g, b, 0.8}
	cmdw := window_command{
		x:   window.xPos,
		y:   window.yPos,
		h:   window.height,
		w:   window.width,
		clr: cl,
		toolbar: toolbar_command{
			h:   30,
			clr: [4]float32{255, 0, 0, 1},
		},
		// clr: [4]float32{1, 0, 0, 0.8},
	}
	cmd := command{
		t:      WindowCmd,
		window: &cmdw,
	}
	c.windows = append(c.windows, window)
	c.rq.AddCommand(cmd)
}
