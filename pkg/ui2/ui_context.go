package ui2

import (
	"log"

	"github.com/Dmitry-dms/moon/pkg/math"
	"github.com/Dmitry-dms/moon/pkg/ui2/fonts"

	"github.com/pkg/errors"
)

var context *UIContext = nil

type UIContext struct {
	Initialized      bool
	Io               *ImIO
	Viewports        []*ImGuiViewportP
	InputEventsQueue []ImGuiInputEvent
	InputEventsTrail []ImGuiInputEvent

	MouseCursor ImGuiMouseCursor

	//new frame
	Hooks                []*ImGuiContextHook[any]
	Time                 float32
	WithinFrameScope     bool
	FrameCount           uint
	WindowsActiveCount   uint
	TooltipOverrideCount float64

	//widget state
	MouseLastValidPos ImVec2

	// Gamepad/keyboard Navigation
	NavDisableMouseHover bool
	NavIdIsAlive         bool
	NavMousePosDirty     bool
	NavDisableHighlight  bool

	PlatformImeData     ImGuiPlatformImeData
	PlatformImeDataPrev ImGuiPlatformImeData

	Windows               []*ImGuiWindow
	WindowsFocusOrder     []*ImGuiWindow
	WindowsTempSortBuffer []*ImGuiWindow

	WithinFrameScopeWithImplicitWindow bool

	NextWindowData ImGuiNextWindowData

	Style ImGuiStyle

	CurrentWindow                  *ImGuiWindow
	HoveredWindow                  *ImGuiWindow
	HoveredWindowUnderMovingWindow *ImGuiWindow
	MovingWindow                   *ImGuiWindow
	WheelingWindow                 *ImGuiWindow
}

type ImGuiNextWindowData struct {
	Flags              ImGuiNextWindowDataFlags
	PosCond            ImGuiCond
	SizeCond           ImGuiCond
	CollapsedCond      ImGuiCond
	PosVal             ImVec2
	PosPivotVal        ImVec2
	SizeVal            ImVec2
	ContentSizeVal     ImVec2
	ScrollVal          ImVec2
	CollapsedVal       bool
	SizeConstraintRect ImRect
	// SizeCallback         ImGuiSizeCallback
	// SizeCallbackUserData int
	BgAlphaVal          float32
	MenuBarOffsetMinVal ImVec2
}
type ImGuiNextWindowDataFlags int

const (
	ImGuiNextWindowDataFlags_None   ImGuiNextWindowDataFlags = 0
	ImGuiNextWindowDataFlags_HasPos ImGuiNextWindowDataFlags = 1 << iota
	ImGuiNextWindowDataFlags_HasSize
	ImGuiNextWindowDataFlags_HasContentSize
	ImGuiNextWindowDataFlags_HasCollapsed
	ImGuiNextWindowDataFlags_HasSizeConstraint
	ImGuiNextWindowDataFlags_HasFocus
	ImGuiNextWindowDataFlags_HasBgAlpha
	ImGuiNextWindowDataFlags_HasScroll
)

func LoadFontFromFile(path string, scale int32) (*fonts.Font, error) {
	f := fonts.NewFont(path, scale, true)
	return f, nil
}

func GetCurrentContext() *UIContext {
	return context
}

type ImGuiMouseCursor int

const (
	ImGuiMouseCursor_None ImGuiMouseCursor = iota
	ImGuiMouseCursor_Arrow
	ImGuiMouseCursor_TextInput
	ImGuiMouseCursor_ResizeAll
	ImGuiMouseCursor_ResizeNS
	ImGuiMouseCursor_ResizeEW
	ImGuiMouseCursor_ResizeNESW
	ImGuiMouseCursor_ResizeNWSE
	ImGuiMouseCursor_Hand
	ImGuiMouseCursor_NotAllowed
	ImGuiMouseCursor_COUNT
)

func initializeContext() {
	viewport := NewImGuiViewportP()
	context.Viewports = append(context.Viewports, viewport)

}

func (i *UIContext) GetMouseCursor() ImGuiMouseCursor {
	return i.MouseCursor
}

func (i *UIContext) pushEvent(e ImGuiInputEvent) {
	i.InputEventsQueue = append(i.InputEventsQueue, e)
}

func CreateContext() *UIContext {
	prevCtx := GetCurrentContext()
	if prevCtx == nil {
		ctx := &UIContext{
			Initialized:           true,
			Io:                    NewIO(),
			Viewports:             make([]*ImGuiViewportP, 1),
			InputEventsQueue:      make([]ImGuiInputEvent, 0),
			InputEventsTrail:      make([]ImGuiInputEvent, 0),
			Hooks:                 make([]*ImGuiContextHook[any], 0),
			Windows:               make([]*ImGuiWindow, 0),
			WindowsFocusOrder:     make([]*ImGuiWindow, 0),
			WindowsTempSortBuffer: make([]*ImGuiWindow, 0),
		}
		context = ctx
	}
	for i, _ := range context.Viewports {
		context.Viewports[i] = NewImGuiViewportP()
	}
	initializeContext()
	return context
}

func (c *UIContext) CallContextHooks(hook_type ImGuiContextHookType) {
	for _, hook := range c.Hooks {
		if hook.Type == hook_type {
			hook.Callback(c, hook)
		}
	}
}

func (g *UIContext) NewFrame() {
	// Remove pending delete hooks before frame start.
	// This deferred removal avoid issues of removal while iterating the hook vector
	for i, hook := range g.Hooks {
		if hook.Type == ImGuiContextHookType_PendingRemoval_ {
			g.Hooks = removeIndex(g.Hooks, i)
		}
	}
	g.CallContextHooks(ImGuiContextHookType_NewFramePre)

	err := ErrorCheckNewFrameSanityChecks()
	if err != nil {
		// panic(err)
	}

	g.Time += g.Io.DeltaTime
	g.WithinFrameScope = true
	g.FrameCount += 1
	g.TooltipOverrideCount = 0
	g.WindowsActiveCount = 0

	g.UpdateViewportsNewFrame()

	// TODO: пропустим момент с текстом
	// Mark rendering data as invalid to prevent user who may have a handle on it to use it.
	for _, v := range g.Viewports {
		v.DrawDataP.Clear()
	}

	//много пропущено
	g.InputEventsTrail = make([]ImGuiInputEvent, 0)
	g.UpdateInputEvents(false)

	g.UpdateKeyboardInputs()

	g.UpdateMouseInputs()

	g.UpdateHoveredWindowAndCaptureFlags() // TODO

	g.MouseCursor = ImGuiMouseCursor_Arrow
	// c.WantCaptureMouseNextFrame = -1
	// c.WantCaptureKeyboardNextFrame = -1
	// c.WantTextInputNextFrame = -1

	// Platform IME data: reset for the frame
	g.PlatformImeDataPrev = g.PlatformImeData
	g.PlatformImeData.WantVisible = false

	for _, window := range g.Windows {
		window.WasActive = window.Active
		window.WasActive = window.Active
		window.BeginCount = 0
		window.Active = false
		window.WriteAccessed = false
	}

	g.WithinFrameScopeWithImplicitWindow = true

	g.SetNextWindowSize(ImVec2{400, 400}, ImGuiCond_FirstUseEver)
}

func (g *UIContext) Begin(name string, p_open *bool, flags ImGuiWindowFlags) bool {
	// style := g.Style

	// Find or create
	window := g.FindWindowByName(name)
	if window == nil {
		window = g.CreateNewWindow(name, flags)
	}

	if (flags & ImGuiWindowFlags_NoInputs) == ImGuiWindowFlags_NoInputs {
		flags |= ImGuiWindowFlags_NoMove | ImGuiWindowFlags_NoResize
	}

	// current_frame := g.FrameCount
	// first_begin_of_the_frame := (window.LastFrameActive != current_frame);

	g.CurrentWindow = window;

	window_stack_data := ImGuiWindowStackData{}
	window_stack_data.Window = window;
    // window_stack_data.ParentLastItemDataBackup = 
	
	
	return false
}

func (g *UIContext) CreateNewWindow(name string, flags ImGuiWindowFlags) *ImGuiWindow {
	window := ImGuiWindow{
		Name:                     name,
		ID:                       0,
		Flags:                    flags,
		Viewport:                 &ImGuiViewportP{},
		Pos:                      ImVec2{},
		Size:                     ImVec2{},
		SizeFull:                 ImVec2{},
		ContentSize:              ImVec2{},
		ContentSizeIdeal:         ImVec2{},
		ContentSizeExplicit:      ImVec2{},
		WindowPadding:            ImVec2{},
		WindowRounding:           0,
		WindowBorderSize:         0,
		NameBufLen:               0,
		MoveId:                   0,
		ChildId:                  0,
		Scroll:                   ImVec2{},
		ScrollMax:                ImVec2{},
		ScrollTarget:             ImVec2{},
		ScrollTargetCenterRatio:  ImVec2{},
		ScrollTargetEdgeSnapDist: ImVec2{},
		ScrollbarSizes:           ImVec2{},
		ScrollbarX:               false,
		ScrollbarY:               false,
		Active:                   false,
		WasActive:                false,
		WriteAccessed:            false,
		Collapsed:                false,
		WantCollapseToggle:       false,
		SkipItems:                false,
		Appearing:                false,
		Hidden:                   false,
		IsFallbackWindow:         false,
		IsExplicitChild:          false,
		HasCloseButton:           false,
		ResizeBorderHeld:         0,
		BeginCount:               0,
		BeginOrderWithinParent:   0,
		BeginOrderWithinContext:  0,
		FocusOrder:               0,
		PopupId:                  0,
		AutoFitFramesX:           0,
		AutoFitFramesY:           0,
		AutoFitChildAxises:       0,
		AutoFitOnlyGrows:         false,
		IDStack:                  []ImGuiID{},
	}

	main_viewport := g.Viewports[0]
	window.Pos = main_viewport.Pos.Add(ImVec2{60, 60})

	window.DC.CursorStartPos = window.Pos
	window.DC.CursorMaxPos = window.Pos
	window.DC.IdealMaxPos = window.Pos

	if (flags & ImGuiWindowFlags_AlwaysAutoResize) != 0 {
		window.AutoFitFramesX = 2
		window.AutoFitFramesY = 2
		window.AutoFitOnlyGrows = false
	} else {
		if window.Size.X <= 0.0 {
			window.AutoFitFramesX = 2
		}
		if window.Size.Y <= 0.0 {
			window.AutoFitFramesY = 2
		}
		window.AutoFitOnlyGrows = (window.AutoFitFramesX > 0) || (window.AutoFitFramesY > 0)
	}
	g.Windows = append(g.Windows, &window)

	return &window
}

func (g *UIContext) FindWindowByName(name string) *ImGuiWindow {
	for _, v := range g.Windows {
		if v.Name == name {
			return v
		}
	}
	return nil
}

func (g *UIContext) SetNextWindowSize(size ImVec2, cond ImGuiCond) {
	g.NextWindowData.Flags |= ImGuiNextWindowDataFlags_HasSize
	g.NextWindowData.SizeVal = size

	if cond != 0 {
		g.NextWindowData.SizeCond = cond
	} else {
		g.NextWindowData.SizeCond = ImGuiCond_Always
	}
}

type ImGuiCond int

const (
	ImGuiCond_None   ImGuiCond = iota
	ImGuiCond_Always ImGuiCond = 1 << iota
	ImGuiCond_Once
	ImGuiCond_FirstUseEver
	ImGuiCond_Appearing
)

type ImGuiPlatformImeData struct {
	WantVisible     bool
	InputPos        ImVec2
	InputLineHeight float32
}

func (g *UIContext) UpdateInputEvents(trickle_fast_inputs bool) {
	io := g.Io

	// var mouse_moved, mouse_wheeled, key_changed bool
	mouse_button_changed := 0x00

	// fmt.Println(len(g.InputEventsQueue))
	for _, e := range g.InputEventsQueue {

		switch e.Type {
		case ImGuiInputEventType_MousePos:
			event_pos := ImVec2{e.MousePos.PosX, e.MousePos.PosY}
			if g.IsMousePosValid(&event_pos) {
				event_pos = ImVec2{event_pos.X, event_pos.Y}
			}
			if io.MousePos.X != event_pos.X || io.MousePos.Y != event_pos.Y {

				io.MousePos = event_pos
				// mouse_moved = true
			}
		case ImGuiInputEventType_MouseButton:
			button := e.MouseBtn.Button
			if io.MouseDown[button] != e.MouseBtn.Down {
				io.MouseDown[button] = e.MouseBtn.Down
				mouse_button_changed |= (1 << button)
			}
		case ImGuiInputEventType_MouseWheel:
			if e.MouseWheel.WheelX != 0.0 || e.MouseWheel.WheelY != 0.0 {
				io.MouseWheelH += e.MouseWheel.WheelX
				io.MouseWheel += e.MouseWheel.WheelY
				// mouse_wheeled = true
			}
		case ImGuiInputEventType_Key:
			key := e.Key.Key

			if key == ImGuiKey_None {
				log.Println("ImGuiKey_None was pressed")
				continue
			}
			keydata_index := (key - ImGuiKey_KeysData_OFFSET)

			keydata := io.KeyData[keydata_index]
			if keydata.Down != e.Key.Down || keydata.AnalogValue != e.Key.AnalogValue {
				keydata.Down = e.Key.Down
				keydata.AnalogValue = e.Key.AnalogValue
				// key_changed = true
			}

			if key == ImGuiKey_ModCtrl || key == ImGuiKey_ModShift || key == ImGuiKey_ModAlt || key == ImGuiKey_ModSuper {
				if key == ImGuiKey_ModCtrl {
					io.KeyCtrl = keydata.Down
				}
				if key == ImGuiKey_ModShift {
					io.KeyShift = keydata.Down
				}
				if key == ImGuiKey_ModAlt {
					io.KeyAlt = keydata.Down
				}
				if key == ImGuiKey_ModSuper {
					io.KeySuper = keydata.Down
				}

			}
		case ImGuiInputEventType_Focus:
			io.AppFocusLost = !e.AppFocused.Focused
		}

	}

	// Record trail (for domain-specific applications wanting to access a precise trail)
	for _, v := range g.InputEventsQueue {
		g.InputEventsTrail = append(g.InputEventsTrail, v)
	}

	if io.AppFocusLost {
		io.ClearInputsKey()
		io.AppFocusLost = false
	}
}

func (g *UIContext) UpdateKeyboardInputs() {
	io := g.Io

	for _, key_data := range io.KeyData {
		key_data.DownDurationPrev = key_data.DownDuration
		if key_data.Down {
			if key_data.DownDuration < 0.0 {
				key_data.DownDuration = 0
			} else {
				key_data.DownDuration += io.DeltaTime
			}

		} else {
			key_data.DownDuration = -1
		}
	}
}

func (g *UIContext) UpdateMouseInputs() {
	io := g.Io

	if g.IsMousePosValid(&io.MousePos) {
		io.MousePos = FloorSigned(io.MousePos)
		g.MouseLastValidPos = FloorSigned(io.MousePos)
	}

	if g.IsMousePosValid(&io.MousePos) && g.IsMousePosValid(&io.MousePosPrev) {
		io.MouseDelta = io.MousePos.Sub(io.MousePosPrev)
	} else {
		io.MouseDelta = ImVec2{0, 0}
	}

	if io.MouseDelta.X != 0.0 || io.MouseDelta.Y != 0.0 {
		g.NavDisableMouseHover = false
	}
	io.MousePosPrev = io.MousePos

	for i := 0; i < len(io.MouseDown); i++ {
		io.MouseClicked[i] = io.MouseDown[i] && io.MouseDownDuration[i] < 0.0
		io.MouseClickedCount[i] = 0 // Will be filled below
		io.MouseReleased[i] = !io.MouseDown[i] && io.MouseDownDuration[i] >= 0.0
		io.MouseDownDurationPrev[i] = io.MouseDownDuration[i]
		if io.MouseDown[i] {
			if io.MouseDownDuration[i] < 0 {
				io.MouseDownDuration[i] = 0
			} else {
				io.MouseDownDuration[i] = io.MouseDownDuration[i] + io.DeltaTime
			}
		} else {
			io.MouseDownDuration[i] = -1
		}

		if io.MouseClicked[i] {
			is_repeated_click := false
			if g.Time-io.MouseClickedTime[i] < io.MouseDoubleClickTime {
				var delta_from_click_pos ImVec2
				if g.IsMousePosValid(&io.MousePos) {
					delta_from_click_pos = io.MousePos.Sub(io.MouseClickedPos[i])
				} else {
					delta_from_click_pos = ImVec2{0, 0}
				}

				if LengthSqrVec2(delta_from_click_pos) < io.MouseDoubleClickMaxDist*io.MouseDoubleClickMaxDist {
					is_repeated_click = true
				}
			}

			if is_repeated_click {
				io.MouseClickedLastCount[i]++
			} else {
				io.MouseClickedLastCount[i] = 0
			}
			io.MouseClickedTime[i] = g.Time
			io.MouseClickedPos[i] = io.MousePos
			io.MouseClickedCount[i] = io.MouseClickedLastCount[i]
			io.MouseDragMaxDistanceSqr[i] = 0.0

		} else if io.MouseDown[i] {
			var delta_sqr_click_pos float32
			if g.IsMousePosValid(&io.MousePos) {
				delta_sqr_click_pos = LengthSqrVec2(io.MousePos.Sub(io.MouseClickedPos[i]))
			} else {
				delta_sqr_click_pos = 0
			}
			io.MouseDragMaxDistanceSqr[i] = math.Max(io.MouseDragMaxDistanceSqr[i], delta_sqr_click_pos)
		}
		// We provide io.MouseDoubleClicked[] as a legacy service
		io.MouseDoubleClicked[i] = (io.MouseClickedCount[i] == 2)

		if io.MouseClicked[i] {
			g.NavDisableMouseHover = false
		}
	}
}

func (g *UIContext) UpdateHoveredWindowAndCaptureFlags() {

}

func (g *UIContext) IsMousePosValid(mouse_pos *ImVec2) bool {
	MOUSE_INVALID := -256000.0
	var p ImVec2
	if mouse_pos != nil {
		p = *mouse_pos
	} else {
		p = g.Io.MousePos
	}
	return p.X >= float32(MOUSE_INVALID) && p.Y >= float32(MOUSE_INVALID)
}

//проверка на ошибки
func ErrorCheckNewFrameSanityChecks() error {
	ctx := GetCurrentContext()

	if !ctx.Initialized {
		return errors.New("Gui is not initialized")
	}
	if ctx.Io.DeltaTime > 0 || ctx.FrameCount == 0 {
		return errors.New("Need a positive DeltaTime!")
	}

	return nil
}
func (c *UIContext) UpdateViewportsNewFrame() {
	if len(c.Viewports) != 1 {
		// panic("Viewports len != 1")
	}
	// Update main viewport with current platform position.
	main_viewport := c.Viewports[0]
	main_viewport.Flags = ImGuiViewportFlags_IsPlatformWindow | ImGuiViewportFlags_OwnedByApp
	main_viewport.Pos = ImVec2{0, 0}
	main_viewport.Size = c.Io.DisplaySize

	for _, v := range c.Viewports {
		v.WorkOffsetMin = v.BuildWorkOffsetMin
		v.WorkOffsetMax = v.BuildWorkOffsetMax
		v.BuildWorkOffsetMax, v.BuildWorkOffsetMin = ImVec2{0, 0}, ImVec2{0, 0}
		v.UpdateWorkRect()
	}
}

func removeIndex[T any](s []T, index int) []T {
	return append(s[:index], s[index+1:]...)
}
