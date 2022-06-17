package ui2

type ImGuiKey int

const (
	ImGuiKey_NamedKey_BEGIN = 512
	// Keyboard
	ImGuiKey_None ImGuiKey = 0
	ImGuiKey_Tab      ImGuiKey     = ImGuiKey_NamedKey_BEGIN + iota // == ImGuiKey_NamedKey_BEGIN
	ImGuiKey_LeftArrow
	ImGuiKey_RightArrow
	ImGuiKey_UpArrow
	ImGuiKey_DownArrow
	ImGuiKey_PageUp
	ImGuiKey_PageDown
	ImGuiKey_Home
	ImGuiKey_End
	ImGuiKey_Insert
	ImGuiKey_Delete
	ImGuiKey_Backspace
	ImGuiKey_Space
	ImGuiKey_Enter
	ImGuiKey_Escape
	ImGuiKey_LeftCtrl
	ImGuiKey_LeftShift
	ImGuiKey_LeftAlt
	ImGuiKey_LeftSuper
	ImGuiKey_RightCtrl
	ImGuiKey_RightShift
	ImGuiKey_RightAlt
	ImGuiKey_RightSuper
	ImGuiKey_Menu
	ImGuiKey_0
	ImGuiKey_1
	ImGuiKey_2
	ImGuiKey_3
	ImGuiKey_4
	ImGuiKey_5
	ImGuiKey_6
	ImGuiKey_7
	ImGuiKey_8
	ImGuiKey_9
	ImGuiKey_A
	ImGuiKey_B
	ImGuiKey_C
	ImGuiKey_D
	ImGuiKey_E
	ImGuiKey_F
	ImGuiKey_G
	ImGuiKey_H
	ImGuiKey_I
	ImGuiKey_J
	ImGuiKey_K
	ImGuiKey_L
	ImGuiKey_M
	ImGuiKey_N
	ImGuiKey_O
	ImGuiKey_P
	ImGuiKey_Q
	ImGuiKey_R
	ImGuiKey_S
	ImGuiKey_T
	ImGuiKey_U
	ImGuiKey_V
	ImGuiKey_W
	ImGuiKey_X
	ImGuiKey_Y
	ImGuiKey_Z

	ImGuiKey_F1
	ImGuiKey_F2
	ImGuiKey_F3
	ImGuiKey_F4
	ImGuiKey_F5
	ImGuiKey_F6
	ImGuiKey_F7
	ImGuiKey_F8
	ImGuiKey_F9
	ImGuiKey_F10
	ImGuiKey_F11
	ImGuiKey_F12
	ImGuiKey_Apostrophe   // '
	ImGuiKey_Comma        // ,
	ImGuiKey_Minus        // -
	ImGuiKey_Period       // .
	ImGuiKey_Slash        // /
	ImGuiKey_Semicolon    // ;
	ImGuiKey_Equal        // =
	ImGuiKey_LeftBracket  // [
	ImGuiKey_Backslash    // \ (this text inhibit multiline comment caused by backslash)
	ImGuiKey_RightBracket // ]
	ImGuiKey_GraveAccent  // `
	ImGuiKey_CapsLock
	ImGuiKey_ScrollLock
	ImGuiKey_NumLock
	ImGuiKey_PrintScreen
	ImGuiKey_Pause
	ImGuiKey_Keypad0
	ImGuiKey_Keypad1
	ImGuiKey_Keypad2
	ImGuiKey_Keypad3
	ImGuiKey_Keypad4
	ImGuiKey_Keypad5
	ImGuiKey_Keypad6
	ImGuiKey_Keypad7
	ImGuiKey_Keypad8
	ImGuiKey_Keypad9
	ImGuiKey_KeypadDecimal
	ImGuiKey_KeypadDivide
	ImGuiKey_KeypadMultiply
	ImGuiKey_KeypadSubtract
	ImGuiKey_KeypadAdd
	ImGuiKey_KeypadEnter
	ImGuiKey_KeypadEqual

	ImGuiKey_ModCtrl
	ImGuiKey_ModShift
	ImGuiKey_ModAlt
	ImGuiKey_ModSuper

	ImGuiKey_COUNT

    ImGuiKey_NamedKey_END           = ImGuiKey_COUNT
    ImGuiKey_NamedKey_COUNT         = ImGuiKey_NamedKey_END - ImGuiKey_NamedKey_BEGIN
)

const (
	ImGuiKey_KeysData_OFFSET = 0
	ImGuiKey_KeysData_SIZE = ImGuiKey_COUNT
)

type ImGuiKeyData struct {
	Down             bool    // True for if key is down
	DownDuration     float32 // Duration the key has been down (<0.0f: not pressed, 0.0f: just pressed, >0.0f: time held)
	DownDurationPrev float32 // Last frame duration the key has been down
	AnalogValue      float32 // 0.0f..1.0f for gamepad values
}
