package ui

type Toolbar struct {
	x, y    float32 // top-left corner
	w, h float32
}

func NewToolbar(x, y, w, h float32) Toolbar {
	tb := Toolbar{
		x:   x,
		y:   y,
		w:  w,
		h: h,
	}
	return tb
}
