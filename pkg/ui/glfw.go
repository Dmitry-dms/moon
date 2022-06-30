package ui

import (
	"github.com/go-gl/glfw/v3.3/glfw"
	"github.com/go-gl/mathgl/mgl32"
)

type glfwHandler struct {
	GlfwWindow         *glfw.Window
	Time               float32
	MouseWindow        *glfw.Window
	MouseCursorse      [9]*glfw.Cursor
	LastValidMousePos  mgl32.Vec2
	InstalledCallbacks bool

	//callbacks
	// PrevUserCallbackWindowFocus GLFWwindowfocusfun
	// PrevUserCallbackCursorPos   GLFWcursorposfun
	// PrevUserCallbackCursorEnter GLFWcursorenterfun
	PrevUserCallbackMousebutton glfw.MouseButtonCallback
	PrevUserCallbackScroll      glfw.ScrollCallback
	PrevUserCallbackKey         glfw.KeyCallback
	// 	PrevUserCallbackChar        GLFWcharfun
	// 	PrevUserCallbackMonitor     GLFWmonitorfun
}

func NewData() *glfwHandler {
	d := glfwHandler{}
	return &d
}

func GlfwKeyToGuiKey(key glfw.Key) GuiKey {
	switch key {
	case glfw.KeyTab:
		return GuiKey_Tab
	case glfw.KeyLeft:
		return GuiKey_LeftArrow
	case glfw.KeyRight:
		return GuiKey_RightArrow
	case glfw.KeyUp:
		return GuiKey_UpArrow
	case glfw.KeyDown:
		return GuiKey_DownArrow
	case glfw.KeyPageUp:
		return GuiKey_PageUp
	case glfw.KeyPageDown:
		return GuiKey_PageDown
	case glfw.KeyHome:
		return GuiKey_Home
	case glfw.KeyEnd:
		return GuiKey_End
	case glfw.KeyInsert:
		return GuiKey_Insert
	case glfw.KeyDelete:
		return GuiKey_Delete
	case glfw.KeyBackspace:
		return GuiKey_Backspace
	case glfw.KeySpace:
		return GuiKey_Space
	case glfw.KeyEnter:
		return GuiKey_Enter
	case glfw.KeyEscape:
		return GuiKey_Escape
	case glfw.KeyApostrophe:
		return GuiKey_Apostrophe
	case glfw.KeyComma:
		return GuiKey_Comma
	case glfw.KeyMinus:
		return GuiKey_Minus
	case glfw.KeyPeriod:
		return GuiKey_Period
	case glfw.KeySlash:
		return GuiKey_Slash
	case glfw.KeySemicolon:
		return GuiKey_Semicolon
	case glfw.KeyEqual:
		return GuiKey_Equal
	case glfw.KeyLeftBracket:
		return GuiKey_LeftBracket
	case glfw.KeyBackslash:
		return GuiKey_Backslash
	case glfw.KeyRightBracket:
		return GuiKey_RightBracket
	case glfw.KeyGraveAccent:
		return GuiKey_GraveAccent
	case glfw.KeyCapsLock:
		return GuiKey_CapsLock
	case glfw.KeyScrollLock:
		return GuiKey_ScrollLock
	case glfw.KeyNumLock:
		return GuiKey_NumLock
	case glfw.KeyPrintScreen:
		return GuiKey_PrintScreen
	case glfw.KeyPause:
		return GuiKey_Pause
	case glfw.KeyKP0:
		return GuiKey_Keypad0
	case glfw.KeyKP1:
		return GuiKey_Keypad1
	case glfw.KeyKP2:
		return GuiKey_Keypad2
	case glfw.KeyKP3:
		return GuiKey_Keypad3
	case glfw.KeyKP4:
		return GuiKey_Keypad4
	case glfw.KeyKP5:
		return GuiKey_Keypad5
	case glfw.KeyKP6:
		return GuiKey_Keypad6
	case glfw.KeyKP7:
		return GuiKey_Keypad7
	case glfw.KeyKP8:
		return GuiKey_Keypad8
	case glfw.KeyKP9:
		return GuiKey_Keypad9
	case glfw.KeyKPDecimal:
		return GuiKey_KeypadDecimal
	case glfw.KeyKPDivide:
		return GuiKey_KeypadDivide
	case glfw.KeyKPMultiply:
		return GuiKey_KeypadMultiply
	case glfw.KeyKPSubtract:
		return GuiKey_KeypadSubtract
	case glfw.KeyKPAdd:
		return GuiKey_KeypadAdd
	case glfw.KeyKPEnter:
		return GuiKey_KeypadEnter
	case glfw.KeyKPEqual:
		return GuiKey_KeypadEqual
	case glfw.KeyLeftShift:
		return GuiKey_LeftShift
	case glfw.KeyLeftControl:
		return GuiKey_LeftCtrl
	case glfw.KeyLeftAlt:
		return GuiKey_LeftAlt
	case glfw.KeyLeftSuper:
		return GuiKey_LeftSuper
	case glfw.KeyRightShift:
		return GuiKey_RightShift
	case glfw.KeyRightControl:
		return GuiKey_RightCtrl
	case glfw.KeyRightAlt:
		return GuiKey_RightAlt
	case glfw.KeyRightSuper:
		return GuiKey_RightSuper
	case glfw.KeyMenu:
		return GuiKey_Menu
	case glfw.Key0:
		return GuiKey_0
	case glfw.Key1:
		return GuiKey_1
	case glfw.Key2:
		return GuiKey_2
	case glfw.Key3:
		return GuiKey_3
	case glfw.Key4:
		return GuiKey_4
	case glfw.Key5:
		return GuiKey_5
	case glfw.Key6:
		return GuiKey_6
	case glfw.Key7:
		return GuiKey_7
	case glfw.Key8:
		return GuiKey_8
	case glfw.Key9:
		return GuiKey_9
	case glfw.KeyA:
		return GuiKey_A
	case glfw.KeyB:
		return GuiKey_B
	case glfw.KeyC:
		return GuiKey_C
	case glfw.KeyD:
		return GuiKey_D
	case glfw.KeyE:
		return GuiKey_E
	case glfw.KeyF:
		return GuiKey_F
	case glfw.KeyG:
		return GuiKey_G
	case glfw.KeyH:
		return GuiKey_H
	case glfw.KeyI:
		return GuiKey_I
	case glfw.KeyJ:
		return GuiKey_J
	case glfw.KeyK:
		return GuiKey_K
	case glfw.KeyL:
		return GuiKey_L
	case glfw.KeyM:
		return GuiKey_M
	case glfw.KeyN:
		return GuiKey_N
	case glfw.KeyO:
		return GuiKey_O
	case glfw.KeyP:
		return GuiKey_P
	case glfw.KeyQ:
		return GuiKey_Q
	case glfw.KeyR:
		return GuiKey_R
	case glfw.KeyS:
		return GuiKey_S
	case glfw.KeyT:
		return GuiKey_T
	case glfw.KeyU:
		return GuiKey_U
	case glfw.KeyV:
		return GuiKey_V
	case glfw.KeyW:
		return GuiKey_W
	case glfw.KeyX:
		return GuiKey_X
	case glfw.KeyY:
		return GuiKey_Y
	case glfw.KeyZ:
		return GuiKey_Z
	case glfw.KeyF1:
		return GuiKey_F1
	case glfw.KeyF2:
		return GuiKey_F2
	case glfw.KeyF3:
		return GuiKey_F3
	case glfw.KeyF4:
		return GuiKey_F4
	case glfw.KeyF5:
		return GuiKey_F5
	case glfw.KeyF6:
		return GuiKey_F6
	case glfw.KeyF7:
		return GuiKey_F7
	case glfw.KeyF8:
		return GuiKey_F8
	case glfw.KeyF9:
		return GuiKey_F9
	case glfw.KeyF10:
		return GuiKey_F10
	case glfw.KeyF11:
		return GuiKey_F11
	case glfw.KeyF12:
		return GuiKey_F12
	default:
		return GuiKey_None
	}
}

func GlfwMouseKey(btn glfw.MouseButton) MouseKey {
	switch btn {
	case glfw.MouseButtonLeft:
		return MouseBtnLeft
	case glfw.MouseButtonRight:
		return MouseBtnRight
	case glfw.MouseButtonMiddle:
		return MouseBtnMiddle
	default:
		return MouseBtnUnknown
	}
}

func GlfwAction(action glfw.Action) Action {
	switch action {
	case glfw.Press:
		return Press
	case glfw.Release:
		return Release
	case glfw.Repeat:
		return Repeat
	default:
		return UnknownAction
	}
}

type PlatformHandler interface {
	IsKeyPressed(key Key)
	UpdateInputs()
}
