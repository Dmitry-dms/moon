package ui

import ()

type Io struct {
	mouseX, mouseY      float32
	winWidth, winHeight float32

	keyPressed [350]bool

	scrollX, scrollY float64

	mouseBtnPressed [9]bool
	isDragging      bool
}

type Key int

func (io *Io) KeyCallback(key Key) {

	
}

var (
	KeyTab         = 0
	KeyLeftArrow   = 1
	KeyRightArrow  = 2
	KeyUpArrow     = 3
	KeyDownArrow   = 4
	KeyPageUp      = 5
	KeyPageDown    = 6
	KeyHome        = 7
	KeyEnd         = 8
	KeyInsert      = 9
	KeyDelete      = 10
	KeyBackspace   = 11
	KeySpace       = 12
	KeyEnter       = 13
	KeyEscape      = 14
	KeyKeyPadEnter = 15
	KeyA           = 16 // for text edit CTRL+A: select all
	KeyC           = 17 // for text edit CTRL+C: copy
	KeyV           = 18 // for text edit CTRL+V: paste
	KeyX           = 19 // for text edit CTRL+X: cut
	KeyY           = 20 // for text edit CTRL+Y: redo
	KeyZ           = 21 // for text edit CTRL+Z: undo
)

func NewIo() *Io {
	i := Io{
		mouseX:          0,
		mouseY:          0,
		winWidth:        0,
		winHeight:       0,
		keyPressed:      [350]bool{},
		scrollX:         0,
		scrollY:         0,
		mouseBtnPressed: [9]bool{},
		isDragging:      false,
	}
	return &i
}
