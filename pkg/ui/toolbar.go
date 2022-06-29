package ui

type Toolbar struct {
	xPos, yPos    float32 // top-left corner
	width, height float32
}

func NewToolbar(x, y, w, h float32) Toolbar {
	tb := Toolbar{
		xPos:   x,
		yPos:   y,
		width:  w,
		height: h,
	}
	return tb
}
