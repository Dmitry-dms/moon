package ui2

import "github.com/Dmitry-dms/moon/pkg/ui2/fonts"

const (
	ImGuiNavInput_COUNT = 20
)

type ImIO struct {
	DisplaySize             Vec2
	DeltaTime               float32
	DisplayFramebufferScale Vec2
	BackendPlatformUserData *ImplGlfw_Data
	BackendRendererUserData *ImGui_ImplOpenGL3_Data

	Name, BackendRendererName string

	SetClipboardTextFn func(string)
	GetClipboardTextFn func() string

	WantCaptureMouse         bool
	WantCaptureKeyboard      bool
	WantTextInput            bool
	WantSetMousePos          bool
	WantSaveIniSettings      bool
	NavActive                bool
	NavVisible               bool
	Framerate                bool
	MetricsRenderVertices    bool
	MetricsRenderIndices     bool
	MetricsRenderWindows     bool
	MetricsActiveWindows     bool
	MetricsActiveAllocations bool
	MouseDelta               Vec2

	MousePos    Vec2
	MouseDown   [5]bool
	MouseWheel  float32
	MouseWheelH float32
	KeyCtrl     bool
	KeyShift    bool
	KeyAlt      bool
	KeySuper    bool
	NavInputs   [ImGuiNavInput_COUNT]float32

	MouseDrawCursor bool

	ConfigFlags  ImGuiConfigFlags
	BackendFlags ImGuiBackendFlags

	//font
	DefaultFont *fonts.Font
}

type ImGuiConfigFlags int
type ImGuiBackendFlags int

const (
	ImGuiBackendFlags_None                 ImGuiBackendFlags = 0
	ImGuiBackendFlags_HasGamepad           ImGuiBackendFlags = 1 << 0
	ImGuiBackendFlags_HasMouseCursors      ImGuiBackendFlags = 1 << 1
	ImGuiBackendFlags_HasSetMousePos       ImGuiBackendFlags = 1 << 2
	ImGuiBackendFlags_RendererHasVtxOffset ImGuiBackendFlags = 1 << 3
)

const (
	ImGuiConfigFlags_None                 ImGuiConfigFlags = 0
	ImGuiConfigFlags_NavEnableKeyboard    ImGuiConfigFlags = 1 << 0
	ImGuiConfigFlags_NavEnableGamepad     ImGuiConfigFlags = 1 << 1
	ImGuiConfigFlags_NavEnableSetMousePos ImGuiConfigFlags = 1 << 2
	ImGuiConfigFlags_NavNoCaptureKeyboard ImGuiConfigFlags = 1 << 3
	ImGuiConfigFlags_NoMouse              ImGuiConfigFlags = 1 << 4
	ImGuiConfigFlags_NoMouseCursorChange  ImGuiConfigFlags = 1 << 5
	// User storage (to allow your backend/engine to
	ImGuiConfigFlags_IsSRGB        ImGuiConfigFlags = 1 << 20
	ImGuiConfigFlags_IsTouchScreen ImGuiConfigFlags = 1 << 21
)



type ImGuiInputEvent struct {
	Type   ImGuiInputEventType
	Source ImGuiInputSource

	mousePos   *ImGuiInputEventMousePos
	mouseWheel *ImGuiInputEventMouseWheel
	mouseBtn   *ImGuiInputEventMouseButton
	text       *ImGuiInputEventText
	appFocused *ImGuiInputEventAppFocused
	key        *ImGuiInputEventKey
}

type ImGuiInputEventType int

const (
	ImGuiInputEventType_None ImGuiInputEventType = 0 + iota
	ImGuiInputEventType_MousePos
	ImGuiInputEventType_MouseWheel
	ImGuiInputEventType_MouseButton
	ImGuiInputEventType_Key
	ImGuiInputEventType_Text
	ImGuiInputEventType_Focus
	ImGuiInputEventType_COUNT
)

type ImGuiInputSource int

const (
	ImGuiInputSource_None ImGuiInputSource = 0 + iota
	ImGuiInputSource_Mouse
	ImGuiInputSource_Keyboard
	ImGuiInputSource_Gamepad
	ImGuiInputSource_Clipboard // Currently only used by InputText()
	ImGuiInputSource_Nav       // Stored in g.ActiveIdSource only
	ImGuiInputSource_COUNT
)

type ImGuiMouseButton int

const (
	ImGuiMouseButton_Left   ImGuiMouseButton = 0
	ImGuiMouseButton_Right                   = 1
	ImGuiMouseButton_Middle                  = 2
	ImGuiMouseButton_COUNT                   = 5
)

type ImGuiInputEventMousePos struct {
	PosX, PosY float32
}
type ImGuiInputEventMouseWheel struct {
	WheelX, WheelY float32
}
type ImGuiInputEventMouseButton struct {
	Button int
	Down   bool
}
type ImGuiInputEventText struct {
	Text string
}
type ImGuiInputEventAppFocused struct {
	Focused bool
}
type ImGuiInputEventKey struct {
	Key         ImGuiKey
	Down        bool
	AnalogValue float32
}

func (i *ImIO) addKeyAnalogEvent(key ImGuiKey, down bool, analogValue float32) {
	if key == ImGuiKey_None {
		return
	}

	g := context

	var e ImGuiInputEvent = ImGuiInputEvent{
		key: &ImGuiInputEventKey{},
	}
	e.Type = ImGuiInputEventType_Key
	e.Source = ImGuiInputSource_Keyboard
	e.key.Key = key
	e.key.AnalogValue = analogValue
	e.key.Down = down

	g.pushEvent(e)
}

func (i *ImIO) AddKeyEvent(key ImGuiKey, down bool) {
	// fmt.Println(key)
	if down {
		i.addKeyAnalogEvent(key, down, 1)
	} else {
		i.addKeyAnalogEvent(key, down, 0)
	}
}

func (i *ImIO) AddMouseButtonEvent(button int, down bool) {
	g := context

	var e ImGuiInputEvent = ImGuiInputEvent{
		mouseBtn: &ImGuiInputEventMouseButton{},
	}
	e.Type = ImGuiInputEventType_MouseButton
	e.Source = ImGuiInputSource_Mouse
	e.mouseBtn.Button = button
	e.mouseBtn.Down = down

	g.pushEvent(e)
}
func (i *ImIO) AddMousePosEvent(x, y float32) {
	g := context

	var e ImGuiInputEvent = ImGuiInputEvent{
		mousePos: &ImGuiInputEventMousePos{},
	}
	e.Type = ImGuiInputEventType_MousePos
	e.Source = ImGuiInputSource_Mouse
	e.mousePos.PosX = x
	e.mousePos.PosY = y

	g.pushEvent(e)
}
func (i *ImIO) AddMouseWheelEvent(wheelX, wheelY float32) {
	g := context

	var e ImGuiInputEvent = ImGuiInputEvent{
		mouseWheel: &ImGuiInputEventMouseWheel{},
	}
	e.Type = ImGuiInputEventType_MouseWheel
	e.Source = ImGuiInputSource_Mouse
	e.mouseWheel.WheelX = wheelX
	e.mouseWheel.WheelY = wheelY

	g.pushEvent(e)
}
