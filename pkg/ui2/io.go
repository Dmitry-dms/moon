package ui2

import "github.com/Dmitry-dms/moon/pkg/ui2/fonts"

const (
	ImGuiNavInput_COUNT = 20
)

type ImIO struct {
	DisplaySize             ImVec2
	DeltaTime               float32
	DisplayFramebufferScale ImVec2
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
	AppFocusLost             bool
	MouseDelta               ImVec2

	MousePos     ImVec2
	MousePosPrev ImVec2
	MouseDown    [5]bool
	MouseWheel   float32
	MouseWheelH  float32
	KeyCtrl      bool
	KeyShift     bool
	KeyAlt       bool
	KeySuper     bool
	NavInputs    [ImGuiNavInput_COUNT]float32

	MouseDrawCursor bool

	ConfigFlags  ImGuiConfigFlags
	BackendFlags ImGuiBackendFlags

	//font
	DefaultFont *fonts.Font

	//key
	KeyData [ImGuiKey_KeysData_SIZE]*ImGuiKeyData

	MouseClickedPos                [5]ImVec2  // Position at time of clicking
	MouseClickedTime               [5]float32 // Time of last click (used to figure out double-click)
	MouseClicked                   [5]bool    // Mouse button went from !Down to Down (same as MouseClickedCount[x] != 0)
	MouseDoubleClicked             [5]bool    // Has mouse button been double-clicked? (same as MouseClickedCount[x] == 2)
	MouseClickedCount              [5]uint16  // == 0 (not clicked), == 1 (same as MouseClicked[]), == 2 (double-clicked), == 3 (triple-clicked) etc. when going from !Down to Down
	MouseClickedLastCount          [5]uint16  // Count successive number of clicks. Stays valid after mouse release. Reset after another click is done.
	MouseReleased                  [5]bool    // Mouse button went from Down to !Down
	MouseDownOwned                 [5]bool    // Track if button was clicked inside a dear imgui window or over void blocked by a popup. We don't request mouse capture from the application if click started outside ImGui bounds.
	MouseDownOwnedUnlessPopupClose [5]bool    // Track if button was clicked inside a dear imgui window.
	MouseDownDuration              [5]float32 // Duration the mouse button has been down (0.0f == just clicked)
	MouseDownDurationPrev          [5]float32 // Previous time the mouse button has been down
	MouseDragMaxDistanceSqr        [5]float32 // Squared maximum distance of how much mouse has traveled from the clicking point (used for moving thresholds)
	NavInputsDownDuration          [ImGuiNavInput_COUNT]float32
	NavInputsDownDurationPrev      [ImGuiNavInput_COUNT]float32

	MouseDoubleClickTime    float32
	MouseDoubleClickMaxDist float32
	MouseDragThreshold      float32
	KeyRepeatDelay          float32
	KeyRepeatRate           float32
}

func NewIO() *ImIO {
	io := ImIO{
		DisplaySize:             ImVec2{},
		BackendPlatformUserData: newData(),
		DeltaTime:               1 / 60,
		MouseDoubleClickTime:    0.3,
		MouseDoubleClickMaxDist: 6,
		MouseDragThreshold:      6,
		KeyRepeatDelay:          0.25,
		KeyRepeatRate:           0.05,
	}
	for i, _ := range io.KeyData {
		io.KeyData[i] = &ImGuiKeyData{}
	}
	return &io
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

	MousePos   *ImGuiInputEventMousePos
	MouseWheel *ImGuiInputEventMouseWheel
	MouseBtn   *ImGuiInputEventMouseButton
	Text       *ImGuiInputEventText
	AppFocused *ImGuiInputEventAppFocused
	Key        *ImGuiInputEventKey
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
		Key: &ImGuiInputEventKey{},
	}
	e.Type = ImGuiInputEventType_Key
	e.Source = ImGuiInputSource_Keyboard
	e.Key.Key = key
	e.Key.AnalogValue = analogValue
	e.Key.Down = down

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
		MouseBtn: &ImGuiInputEventMouseButton{},
	}
	e.Type = ImGuiInputEventType_MouseButton
	e.Source = ImGuiInputSource_Mouse
	e.MouseBtn.Button = button
	e.MouseBtn.Down = down

	g.pushEvent(e)
}
func (i *ImIO) AddMousePosEvent(x, y float32) {
	g := context

	var e ImGuiInputEvent = ImGuiInputEvent{
		MousePos: &ImGuiInputEventMousePos{},
	}
	e.Type = ImGuiInputEventType_MousePos
	e.Source = ImGuiInputSource_Mouse
	e.MousePos.PosX = x
	e.MousePos.PosY = y

	g.pushEvent(e)
}
func (i *ImIO) AddMouseWheelEvent(wheelX, wheelY float32) {
	g := context

	var e ImGuiInputEvent = ImGuiInputEvent{
		MouseWheel: &ImGuiInputEventMouseWheel{},
	}
	e.Type = ImGuiInputEventType_MouseWheel
	e.Source = ImGuiInputSource_Mouse
	e.MouseWheel.WheelX = wheelX
	e.MouseWheel.WheelY = wheelY

	g.pushEvent(e)
}

func (i *ImIO) ClearInputsKey() {
	for _, v := range i.KeyData {
		v.Down = false
		v.DownDuration = -1.0
		v.DownDurationPrev = -1.0
	}
	i.KeyCtrl = false
	i.KeyShift = false
	i.KeyAlt = false
	i.KeySuper = false

}
