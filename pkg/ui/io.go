package ui

import "github.com/go-gl/mathgl/mgl32"

type Io struct {
	MousePos    mgl32.Vec2
	DisplaySize mgl32.Vec2

	keyPressed [570]bool

	scrollX, scrollY float64
	WantCaptureMouse bool

	mouseBtnPressed [9]bool
	IsDragging      bool
}

type Key int

func (io *Io) KeyCallback(key GuiKey, pressed bool) {
	if pressed {
		io.keyPressed[key] = true
	} else {
		io.keyPressed[key] = false
	}
}

func (io *Io) MousePosCallback(x, y float32) {
	io.MousePos[0] = x
	io.MousePos[1] = y

	io.IsDragging = io.mouseBtnPressed[0] ||
		io.mouseBtnPressed[1] ||
		io.mouseBtnPressed[2]
}

func (io *Io) MouseBtnCallback(mouseBtn MouseKey, action Action) {
	switch action {
	case Press:
		//если кнопок на мышке больше, игнорируем нажатие
		if int(mouseBtn) < len(io.mouseBtnPressed) {
			io.mouseBtnPressed[mouseBtn] = true
		}
	case Release:
		if int(mouseBtn) < len(io.mouseBtnPressed) {
			io.mouseBtnPressed[mouseBtn] = false
			io.IsDragging = false
		}
	}
}

func (io *Io) IsKeyPressed(key GuiKey) bool {
	return io.keyPressed[key]
}

func NewIo() *Io {
	i := Io{
		MousePos:         [2]float32{},
		DisplaySize:      [2]float32{},
		keyPressed:       [570]bool{},
		scrollX:          0,
		scrollY:          0,
		WantCaptureMouse: false,
		mouseBtnPressed:  [9]bool{},
		IsDragging:       false,
	}
	return &i
}
