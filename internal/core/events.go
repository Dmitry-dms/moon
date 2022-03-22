package core

func init() {

}

type events struct {
	scrollX, scrollY float32
	xPos, yPos float32
	lastWorldX, lastWorldY float32;
	mouseBtnPressed [9]bool
	isDragging bool
	mouseBtnDown int8
}
