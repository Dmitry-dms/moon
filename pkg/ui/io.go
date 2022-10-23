package ui

import (
	"github.com/Dmitry-dms/moon/pkg/ui/utils"
)

// "fmt"

// "fmt"

const (
	MOUSE_BTN_COUNT = 5
)

type Io struct {
	MousePos    utils.Vec2
	DisplaySize *utils.Vec2

	keyPressed [610]bool

	// scrollX, scrollY float64
	WantCaptureMouse bool

	MouseDown     [MOUSE_BTN_COUNT]bool
	MouseReleased [MOUSE_BTN_COUNT]bool

	MouseClicked            [MOUSE_BTN_COUNT]bool
	MouseDoubleClicked      [MOUSE_BTN_COUNT]bool
	MouseClickedTime        [MOUSE_BTN_COUNT]float32
	MouseDownDuration       [MOUSE_BTN_COUNT]float32
	MouseDownDurationPrev   [MOUSE_BTN_COUNT]float32
	MouseDragMaxDistanceSqr [MOUSE_BTN_COUNT]float32 //Squared maximum distance of how much mouse has traveled from the clicking point (used for moving thresholds)
	MouseClickedPos         [MOUSE_BTN_COUNT]utils.Vec2

	MouseClickedCount     [MOUSE_BTN_COUNT]uint8 // == 0 (not clicked), == 1 (same as MouseClicked[]), == 2 (double-clicked), == 3 (triple-clicked)
	MouseClickedLastCount [MOUSE_BTN_COUNT]uint8 //Count successive number of clicks. Stays valid after mouse release. Reset after another click is done.

	IsDragging      bool
	dragDelta       utils.Vec2
	dragStarted     utils.Vec2
	dragStartedMain utils.Vec2

	//refactor
	MousePosPrev utils.Vec2
	MouseDelta   utils.Vec2

	DeltaTime               float32 // = 1.0f/60.0f
	MouseDoubleClickTime    float32 // 0.3
	MouseDoubleClickMaxDist float32 // 6.0

	//callbacks
	SetCursor func(c CursorType)

	//scroll
	ScrollX, ScrollY float64
}

type Key int

func (io *Io) SetDisplaySize(w, h float32) {
	io.DisplaySize.X = w
	io.DisplaySize.Y = h
}

func (io *Io) KeyCallback(key GuiKey, pressed bool) {
	if key == GuiKey_None {
		return
	}
	if pressed {
		io.keyPressed[key] = true
	} else {
		io.keyPressed[key] = false
	}
}

func (io *Io) DragStarted(rect utils.Rect) bool {
	return utils.PointInRect(io.dragStarted, rect)
}

func (io *Io) MousePosCallback(x, y float32) {
	io.MousePos.X = x
	io.MousePos.Y = y

	io.IsDragging = io.MouseDown[0] ||
		io.MouseDown[1] ||
		io.MouseDown[2]
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
			io.dragStarted = utils.Vec2{}
			io.dragStartedMain = utils.Vec2{}
		}
	}
}

func (io *Io) IsKeyPressed(key GuiKey) bool {
	return io.keyPressed[key]
}

const (
	MOUSE_INVALID float32 = -256000.0
)

func (io *Io) IsMousePosValid(pos *utils.Vec2) bool {
	var r utils.Vec2
	if pos != nil {
		r = *pos
	} else {
		r = io.MousePos
	}
	return r.X >= MOUSE_INVALID && r.Y >= MOUSE_INVALID
}

func NewIo() *Io {
	i := Io{
		MousePos:                utils.Vec2{0, 0},
		DisplaySize:             &utils.Vec2{0, 0},
		keyPressed:              [610]bool{},
		WantCaptureMouse:        false,
		MouseDown:               [MOUSE_BTN_COUNT]bool{},
		MouseReleased:           [5]bool{},
		MouseClicked:            [5]bool{},
		MouseDoubleClicked:      [5]bool{},
		MouseClickedTime:        [5]float32{},
		MouseDownDuration:       [5]float32{},
		MouseDownDurationPrev:   [5]float32{},
		MouseDragMaxDistanceSqr: [5]float32{},
		MouseClickedPos:         [5]utils.Vec2{},
		MouseClickedCount:       [5]uint8{},
		MouseClickedLastCount:   [5]uint8{},
		IsDragging:              false,
		dragDelta:               utils.Vec2{},
		dragStarted:             utils.Vec2{},
		dragStartedMain:         utils.Vec2{},
		MousePosPrev:            utils.Vec2{},
		MouseDelta:              utils.Vec2{},
		DeltaTime:               1 / 60,
		MouseDoubleClickTime:    0.3,
		MouseDoubleClickMaxDist: 6,
	}
	return &i
}
