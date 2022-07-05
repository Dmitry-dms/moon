package ui

import (
// "fmt"

// "fmt"

)

const (
	MOUSE_BTN_COUNT = 5
)

type Io struct {
	MousePos    Vec2
	DisplaySize Vec2

	keyPressed [570]bool

	scrollX, scrollY float64
	WantCaptureMouse bool

	MouseDown     [MOUSE_BTN_COUNT]bool
	MouseReleased [MOUSE_BTN_COUNT]bool

	MouseClicked            [MOUSE_BTN_COUNT]bool
	MouseDoubleClicked      [MOUSE_BTN_COUNT]bool
	MouseClickedTime        [MOUSE_BTN_COUNT]float32
	MouseDownDuration       [MOUSE_BTN_COUNT]float32
	MouseDownDurationPrev   [MOUSE_BTN_COUNT]float32
	MouseDragMaxDistanceSqr [MOUSE_BTN_COUNT]float32 //Squared maximum distance of how much mouse has traveled from the clicking point (used for moving thresholds)
	MouseClickedPos         [MOUSE_BTN_COUNT]Vec2

	MouseClickedCount     [MOUSE_BTN_COUNT]uint8 // == 0 (not clicked), == 1 (same as MouseClicked[]), == 2 (double-clicked), == 3 (triple-clicked)
	MouseClickedLastCount [MOUSE_BTN_COUNT]uint8 //Count successive number of clicks. Stays valid after mouse release. Reset after another click is done.

	IsDragging      bool
	dragDelta       Vec2
	dragStarted     Vec2
	dragStartedMain Vec2

	//refactor
	MousePosPrev Vec2
	MouseDelta   Vec2

	DeltaTime               float32 // = 1.0f/60.0f
	MouseDoubleClickTime    float32 // 0.3
	MouseDoubleClickMaxDist float32 // 6.0

	//callbacks
	SetCursor func(c CursorType)

	
}

type Key int

func (io *Io) KeyCallback(key GuiKey, pressed bool) {
	if pressed {
		io.keyPressed[key] = true
	} else {
		io.keyPressed[key] = false
	}
}

var lastMousePos = Vec2{}

func (io *Io) MousePosCallback(x, y float32) {
	io.MousePos.X = x
	io.MousePos.Y = y

	io.IsDragging = io.MouseDown[0] ||
		io.MouseDown[1] ||
		io.MouseDown[2]

	// if io.IsDragging {

	// if io.MousePos.X() - io.dragStarted.X()  <=2 && io.MousePos.Y() - io.dragStarted.Y()  <=2 {
	// 	io.dragDelta = mgl32.Vec2{}
	// 	io.dragStarted = io.MousePos
	// }
	// io.dragDelta = io.MousePos.Sub(io.dragStarted)
	// io.dragStarted = io.MousePos
	// fmt.Println(io.dragDelta)
	// lastMousePos = io.MousePos
	// } else {
	// 	io.dragDelta = Vec2{0, 0}
	// io.dragStarted = io.MousePos
	// }
	// fmt.Println(io.dragStarted, io.MousePos, io.dragDelta)
	// fmt.Println(io.dragDelta)
}

func (io *Io) MouseBtnCallback(mouseBtn MouseKey, action Action) {
	switch action {
	case Press:
		//если кнопок на мышке больше, игнорируем нажатие
		if int(mouseBtn) < len(io.MouseDown) {
			io.MouseDown[mouseBtn] = true
			io.dragStarted = io.MousePos
			io.dragStartedMain = io.MousePos
		}
	case Release:
		if int(mouseBtn) < len(io.MouseDown) {
			io.MouseDown[mouseBtn] = false
			io.IsDragging = false
			io.dragStarted = Vec2{}
			io.dragStartedMain = Vec2{}
		}
	}
}

func (io *Io) IsKeyPressed(key GuiKey) bool {
	return io.keyPressed[key]
}

const (
	MOUSE_INVALID float32 = -256000.0
)

func (io *Io) IsMousePosValid(pos *Vec2) bool {
	var r Vec2
	if pos != nil {
		r = *pos
	} else {
		r = io.MousePos
	}
	return r.X >= MOUSE_INVALID && r.Y >= MOUSE_INVALID
}

func NewIo() *Io {
	i := Io{
		MousePos:                Vec2{0, 0},
		DisplaySize:             Vec2{0, 0},
		keyPressed:              [570]bool{},
		scrollX:                 0,
		scrollY:                 0,
		WantCaptureMouse:        false,
		MouseDown:               [MOUSE_BTN_COUNT]bool{},
		MouseReleased:           [5]bool{},
		MouseClicked:            [5]bool{},
		MouseDoubleClicked:      [5]bool{},
		MouseClickedTime:        [5]float32{},
		MouseDownDuration:       [5]float32{},
		MouseDownDurationPrev:   [5]float32{},
		MouseDragMaxDistanceSqr: [5]float32{},
		MouseClickedPos:         [5]Vec2{},
		MouseClickedCount:       [5]uint8{},
		MouseClickedLastCount:   [5]uint8{},
		IsDragging:              false,
		dragDelta:               Vec2{},
		dragStarted:             Vec2{},
		dragStartedMain:         Vec2{},
		MousePosPrev:            Vec2{},
		MouseDelta:              Vec2{},
		DeltaTime:               1 / 60,
		MouseDoubleClickTime:    0.3,
		MouseDoubleClickMaxDist: 6,
	}
	return &i
}
