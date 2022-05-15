package ui2

import (
	"github.com/go-gl/glfw/v3.3/glfw"
)

type GlfwClientApi int

const (
	GlfwClientApi_Unknown GlfwClientApi = iota
	GlfwClientApi_OpenGL
	GlfwClientApi_Vulkan
)

type ImplGlfw_Data struct {
	GlfwWindow         *glfw.Window
	ClientApi          GlfwClientApi
	Time               float32
	MouseWindow        *glfw.Window
	MouseCursorse      [9]*glfw.Cursor
	LastValidMousePos  Vec2
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

func newData() *ImplGlfw_Data {
	d := ImplGlfw_Data{}
	return &d
}

func getClipBoardText() string {
	return glfw.GetClipboardString()
}

func setClipboardText(text string) {
	glfw.SetClipboardString(text)
}

func ImplGlfw_KeyToImGuiKey(key glfw.Key) ImGuiKey {
	switch key {
	case glfw.KeyTab:
		return ImGuiKey_Tab
	case glfw.KeyLeft:
		return ImGuiKey_LeftArrow
	case glfw.KeyRight:
		return ImGuiKey_RightArrow
	case glfw.KeyUp:
		return ImGuiKey_UpArrow
	case glfw.KeyDown:
		return ImGuiKey_DownArrow
	case glfw.KeyPageUp:
		return ImGuiKey_PageUp
	case glfw.KeyPageDown:
		return ImGuiKey_PageDown
	case glfw.KeyHome:
		return ImGuiKey_Home
	case glfw.KeyEnd:
		return ImGuiKey_End
	case glfw.KeyInsert:
		return ImGuiKey_Insert
	case glfw.KeyDelete:
		return ImGuiKey_Delete
	case glfw.KeyBackspace:
		return ImGuiKey_Backspace
	case glfw.KeySpace:
		return ImGuiKey_Space
	case glfw.KeyEnter:
		return ImGuiKey_Enter
	case glfw.KeyEscape:
		return ImGuiKey_Escape
	case glfw.KeyApostrophe:
		return ImGuiKey_Apostrophe
	case glfw.KeyComma:
		return ImGuiKey_Comma
	case glfw.KeyMinus:
		return ImGuiKey_Minus
	case glfw.KeyPeriod:
		return ImGuiKey_Period
	case glfw.KeySlash:
		return ImGuiKey_Slash
	case glfw.KeySemicolon:
		return ImGuiKey_Semicolon
	case glfw.KeyEqual:
		return ImGuiKey_Equal
	case glfw.KeyLeftBracket:
		return ImGuiKey_LeftBracket
	case glfw.KeyBackslash:
		return ImGuiKey_Backslash
	case glfw.KeyRightBracket:
		return ImGuiKey_RightBracket
	case glfw.KeyGraveAccent:
		return ImGuiKey_GraveAccent
	case glfw.KeyCapsLock:
		return ImGuiKey_CapsLock
	case glfw.KeyScrollLock:
		return ImGuiKey_ScrollLock
	case glfw.KeyNumLock:
		return ImGuiKey_NumLock
	case glfw.KeyPrintScreen:
		return ImGuiKey_PrintScreen
	case glfw.KeyPause:
		return ImGuiKey_Pause
	case glfw.KeyKP0:
		return ImGuiKey_Keypad0
	case glfw.KeyKP1:
		return ImGuiKey_Keypad1
	case glfw.KeyKP2:
		return ImGuiKey_Keypad2
	case glfw.KeyKP3:
		return ImGuiKey_Keypad3
	case glfw.KeyKP4:
		return ImGuiKey_Keypad4
	case glfw.KeyKP5:
		return ImGuiKey_Keypad5
	case glfw.KeyKP6:
		return ImGuiKey_Keypad6
	case glfw.KeyKP7:
		return ImGuiKey_Keypad7
	case glfw.KeyKP8:
		return ImGuiKey_Keypad8
	case glfw.KeyKP9:
		return ImGuiKey_Keypad9
	case glfw.KeyKPDecimal:
		return ImGuiKey_KeypadDecimal
	case glfw.KeyKPDivide:
		return ImGuiKey_KeypadDivide
	case glfw.KeyKPMultiply:
		return ImGuiKey_KeypadMultiply
	case glfw.KeyKPSubtract:
		return ImGuiKey_KeypadSubtract
	case glfw.KeyKPAdd:
		return ImGuiKey_KeypadAdd
	case glfw.KeyKPEnter:
		return ImGuiKey_KeypadEnter
	case glfw.KeyKPEqual:
		return ImGuiKey_KeypadEqual
	case glfw.KeyLeftShift:
		return ImGuiKey_LeftShift
	case glfw.KeyLeftControl:
		return ImGuiKey_LeftCtrl
	case glfw.KeyLeftAlt:
		return ImGuiKey_LeftAlt
	case glfw.KeyLeftSuper:
		return ImGuiKey_LeftSuper
	case glfw.KeyRightShift:
		return ImGuiKey_RightShift
	case glfw.KeyRightControl:
		return ImGuiKey_RightCtrl
	case glfw.KeyRightAlt:
		return ImGuiKey_RightAlt
	case glfw.KeyRightSuper:
		return ImGuiKey_RightSuper
	case glfw.KeyMenu:
		return ImGuiKey_Menu
	case glfw.Key0:
		return ImGuiKey_0
	case glfw.Key1:
		return ImGuiKey_1
	case glfw.Key2:
		return ImGuiKey_2
	case glfw.Key3:
		return ImGuiKey_3
	case glfw.Key4:
		return ImGuiKey_4
	case glfw.Key5:
		return ImGuiKey_5
	case glfw.Key6:
		return ImGuiKey_6
	case glfw.Key7:
		return ImGuiKey_7
	case glfw.Key8:
		return ImGuiKey_8
	case glfw.Key9:
		return ImGuiKey_9
	case glfw.KeyA:
		return ImGuiKey_A
	case glfw.KeyB:
		return ImGuiKey_B
	case glfw.KeyC:
		return ImGuiKey_C
	case glfw.KeyD:
		return ImGuiKey_D
	case glfw.KeyE:
		return ImGuiKey_E
	case glfw.KeyF:
		return ImGuiKey_F
	case glfw.KeyG:
		return ImGuiKey_G
	case glfw.KeyH:
		return ImGuiKey_H
	case glfw.KeyI:
		return ImGuiKey_I
	case glfw.KeyJ:
		return ImGuiKey_J
	case glfw.KeyK:
		return ImGuiKey_K
	case glfw.KeyL:
		return ImGuiKey_L
	case glfw.KeyM:
		return ImGuiKey_M
	case glfw.KeyN:
		return ImGuiKey_N
	case glfw.KeyO:
		return ImGuiKey_O
	case glfw.KeyP:
		return ImGuiKey_P
	case glfw.KeyQ:
		return ImGuiKey_Q
	case glfw.KeyR:
		return ImGuiKey_R
	case glfw.KeyS:
		return ImGuiKey_S
	case glfw.KeyT:
		return ImGuiKey_T
	case glfw.KeyU:
		return ImGuiKey_U
	case glfw.KeyV:
		return ImGuiKey_V
	case glfw.KeyW:
		return ImGuiKey_W
	case glfw.KeyX:
		return ImGuiKey_X
	case glfw.KeyY:
		return ImGuiKey_Y
	case glfw.KeyZ:
		return ImGuiKey_Z
	case glfw.KeyF1:
		return ImGuiKey_F1
	case glfw.KeyF2:
		return ImGuiKey_F2
	case glfw.KeyF3:
		return ImGuiKey_F3
	case glfw.KeyF4:
		return ImGuiKey_F4
	case glfw.KeyF5:
		return ImGuiKey_F5
	case glfw.KeyF6:
		return ImGuiKey_F6
	case glfw.KeyF7:
		return ImGuiKey_F7
	case glfw.KeyF8:
		return ImGuiKey_F8
	case glfw.KeyF9:
		return ImGuiKey_F9
	case glfw.KeyF10:
		return ImGuiKey_F10
	case glfw.KeyF11:
		return ImGuiKey_F11
	case glfw.KeyF12:
		return ImGuiKey_F12
	default:
		return ImGuiKey_None
	}
}

func ImplGlfw_KeyToModifier(key glfw.Key) glfw.ModifierKey {
	if key == glfw.KeyLeftControl || key == glfw.KeyRightControl {
		return glfw.ModControl
	}
	if key == glfw.KeyLeftShift || key == glfw.KeyRightShift {
		return glfw.ModShift
	}
	if key == glfw.KeyLeftAlt || key == glfw.KeyRightAlt {
		return glfw.ModAlt
	}
	if key == glfw.KeyLeftSuper || key == glfw.KeyRightSuper {
		return glfw.ModSuper
	}
	return 0
}

func ImplGlfw_UpdateKeyModifiers(mods glfw.ModifierKey) {
	io := context.Io
	if mods > 0 {
		io.AddKeyEvent(ImGuiKey_ModCtrl, (mods&glfw.ModControl != 0))
		io.AddKeyEvent(ImGuiKey_ModShift, (mods&glfw.ModShift != 0))
		io.AddKeyEvent(ImGuiKey_ModAlt, (mods&glfw.ModAlt != 0))
		io.AddKeyEvent(ImGuiKey_ModSuper, (mods&glfw.ModSuper != 0))
	}

}

func ImplGlfw_MouseButtonCallback(w *glfw.Window, button glfw.MouseButton, action glfw.Action, mods glfw.ModifierKey) {
	ImplGlfw_UpdateKeyModifiers(mods)

	io := context.Io
	if button >= 0 && button < ImGuiMouseButton_COUNT {
		io.AddMouseButtonEvent(int(button), action == glfw.Press)
	}
}

func ImplGlfw_ScrollCallback(window *glfw.Window, xOffset, yOffset float64) {
	io := context.Io
	io.AddMouseWheelEvent(float32(xOffset), float32(yOffset))
}

func ImplGlfw_KeyCallback(w *glfw.Window, key glfw.Key, scancode int, action glfw.Action, mods glfw.ModifierKey) {
	if action != glfw.Press && action != glfw.Release {
		return
	}

	keycode_to_mod := ImplGlfw_KeyToModifier(key)
	if action == glfw.Press {
		mods = mods | keycode_to_mod
	} else {
		mods = mods &^ keycode_to_mod
	}
	ImplGlfw_UpdateKeyModifiers(mods)

	io := context.Io
	Imkey := ImplGlfw_KeyToImGuiKey(key)
	io.AddKeyEvent(Imkey, action == glfw.Press)
}

func ImplGlfw_InstallCallbacks(window *glfw.Window) {

	io := context.Io
	io.BackendPlatformUserData.PrevUserCallbackKey = window.SetKeyCallback(ImplGlfw_KeyCallback)
	io.BackendPlatformUserData.PrevUserCallbackScroll = window.SetScrollCallback(ImplGlfw_ScrollCallback)
	io.BackendPlatformUserData.PrevUserCallbackMousebutton = window.SetMouseButtonCallback(ImplGlfw_MouseButtonCallback)

	io.BackendPlatformUserData.InstalledCallbacks = true
}

func ImGui_ImplGlfw_Init(window *glfw.Window, install_callbacks bool, client_api GlfwClientApi) bool {
	io := context.Io
	bd := io.BackendPlatformUserData

	io.BackendFlags |= ImGuiBackendFlags_HasMouseCursors
	io.BackendFlags |= ImGuiBackendFlags_HasSetMousePos

	io.ConfigFlags |= ImGuiConfigFlags_NoMouseCursorChange


	io.Name = "imgui_impl_glfw"

	bd.GlfwWindow = window
	bd.Time = 0.0

	io.SetClipboardTextFn = setClipboardText
	io.GetClipboardTextFn = getClipBoardText

	if install_callbacks {
		ImplGlfw_InstallCallbacks(window)
	}
	bd.ClientApi = client_api
	return true
}

func ImGui_ImplGlfw_InitForOpenGL(window *glfw.Window, install_callbacks bool) bool {

	return ImGui_ImplGlfw_Init(window, install_callbacks, GlfwClientApi_OpenGL)
}

func ImplGlfw_NewFrame(window *glfw.Window) {

	io := GetCurrentContext().Io
	bd := io.BackendPlatformUserData

	// Setup display size (every frame to accommodate for window resizing)
	w, h := window.GetSize()
	display_w, display_h := window.GetFramebufferSize()

	io.DisplaySize = Vec2{float32(w), float32(h)}
	if w > 0 && h > 0 {
		io.DisplayFramebufferScale = Vec2{float32(display_w / w), float32(display_h / h)}
	}

	// Setup time step
	current_time := glfw.GetTime()
	if bd.Time > 0 {
		io.DeltaTime = float32(current_time) - bd.Time
	} else {
		io.DeltaTime = 1 / 60
	}
	bd.Time = float32(current_time)

	ImplGlfw_UpdateMouseData(window)
	ImplGlfw_UpdateMouseCursor(window)
}

func ImplGlfw_UpdateMouseCursor(window *glfw.Window) {
	io := GetCurrentContext().Io

	if (io.ConfigFlags&ImGuiConfigFlags_NoMouseCursorChange) != 0 || (window.GetInputMode(glfw.CursorMode) == glfw.CursorDisabled) {
		return
	}

	cursor := GetCurrentContext().GetMouseCursor()

	if cursor == ImGuiMouseCursor_None || io.MouseDrawCursor {
		// Hide OS mouse cursor if imgui is drawing it or if it wants no cursor
		window.SetInputMode(glfw.CursorMode, glfw.CursorHidden)
	} else {
		// Show OS mouse cursor
		// FIXME-PLATFORM: Unfocused windows seems to fail changing the mouse cursor with GLFW 3.2, but 3.3 works here.
		window.SetCursor(glfw.CreateStandardCursor(glfw.ArrowCursor))
		window.SetInputMode(glfw.CursorMode, glfw.CursorNormal)
	}
}

func ImplGlfw_UpdateMouseData(window *glfw.Window) {
	io := GetCurrentContext().Io
	bd := io.BackendPlatformUserData

	isFocused := window.GetAttrib(glfw.Focused) != 0

	if isFocused {
		// (Optional) Set OS mouse position from Dear ImGui if requested (rarely used, only when ImGuiConfigFlags_NavEnableSetMousePos is enabled by user)
		if io.WantSetMousePos {
			window.SetCursorPos(float64(io.MousePos.X), float64(io.MousePos.Y))
		}
		// (Optional) Fallback to provide mouse position when focused (ImGui_ImplGlfw_CursorPosCallback already provides this when hovered or captured)
		if isFocused && bd.MouseWindow == nil {
			mouseX, mouseY := window.GetCursorPos()

			io.AddMousePosEvent(float32(mouseX), float32(mouseY))
			bd.LastValidMousePos = Vec2{float32(mouseX), float32(mouseY)}
		}
	}
}
