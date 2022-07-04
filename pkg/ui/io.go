package ui

import (
	// "fmt"

	// "fmt"

	"github.com/go-gl/mathgl/mgl32"
)

type Io struct {
	MousePos    mgl32.Vec2
	DisplaySize mgl32.Vec2

	keyPressed [570]bool

	scrollX, scrollY float64
	WantCaptureMouse bool

	mouseBtnPressed [9]bool
	IsDragging      bool
	dragDelta       mgl32.Vec2
	dragStarted     mgl32.Vec2
	dragStartedMain     mgl32.Vec2
}

type Key int

func (io *Io) KeyCallback(key GuiKey, pressed bool) {
	if pressed {
		io.keyPressed[key] = true
	} else {
		io.keyPressed[key] = false
	}
}

var lastMousePos = mgl32.Vec2{}

func (io *Io) MousePosCallback(x, y float32) {
	io.MousePos[0] = x
	io.MousePos[1] = y

	io.IsDragging = io.mouseBtnPressed[0] ||
		io.mouseBtnPressed[1] ||
		io.mouseBtnPressed[2]

	if io.IsDragging {

		// if io.MousePos.X() - io.dragStarted.X()  <=2 && io.MousePos.Y() - io.dragStarted.Y()  <=2 {
		// 	io.dragDelta = mgl32.Vec2{}
		// 	io.dragStarted = io.MousePos
		// }
		io.dragDelta = io.MousePos.Sub(io.dragStarted)
		io.dragStarted = io.MousePos
		// fmt.Println(io.dragDelta)
		// lastMousePos = io.MousePos
	} else {
		io.dragDelta = mgl32.Vec2{0, 0}
		// io.dragStarted = io.MousePos
	}
	// fmt.Println(io.dragStarted, io.MousePos, io.dragDelta)
	// fmt.Println(io.dragDelta)
}

func (io *Io) MouseBtnCallback(mouseBtn MouseKey, action Action) {
	switch action {
	case Press:
		//если кнопок на мышке больше, игнорируем нажатие
		if int(mouseBtn) < len(io.mouseBtnPressed) {
			io.mouseBtnPressed[mouseBtn] = true
			io.dragStarted = io.MousePos
			io.dragStartedMain = io.MousePos
		}
	case Release:
		if int(mouseBtn) < len(io.mouseBtnPressed) {
			io.mouseBtnPressed[mouseBtn] = false
			io.IsDragging = false
			io.dragStarted = mgl32.Vec2{}
			io.dragStartedMain = mgl32.Vec2{}
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
