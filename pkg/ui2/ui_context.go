package ui2

import (
	"github.com/Dmitry-dms/moon/pkg/ui2/fonts"

	"github.com/pkg/errors"
)

var context *UIContext = nil

type UIContext struct {
	Initialized      bool
	Io               *ImIO
	Viewports        []*ImGuiViewportP
	InputEventsQueue []ImGuiInputEvent

	MouseCursor ImGuiMouseCursor

	//new frame
	Hooks                []*ImGuiContextHook[any]
	Time                 float32
	WithinFrameScope     bool
	FrameCount           uint
	WindowsActiveCount   uint
	TooltipOverrideCount float64


}

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
			Initialized: true,
			Io: &ImIO{
				DisplaySize:             Vec2{},
				BackendPlatformUserData: newData(),
			},
			Viewports:        make([]*ImGuiViewportP, 0),
			InputEventsQueue: make([]ImGuiInputEvent, 0),
			Hooks:            make([]*ImGuiContextHook[any], 0),
		}
		context = ctx
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

func (c *UIContext) NewFrame() {
	// Remove pending delete hooks before frame start.
	// This deferred removal avoid issues of removal while iterating the hook vector
	for i, hook := range c.Hooks {
		if hook.Type == ImGuiContextHookType_PendingRemoval_ {
			c.Hooks = removeIndex(c.Hooks, i)
		}
	}
	c.CallContextHooks(ImGuiContextHookType_NewFramePre)

	err := ErrorCheckNewFrameSanityChecks()
	if err != nil {
		panic(err)
	}

	c.Time += c.Io.DeltaTime
	c.WithinFrameScope = true
	c.FrameCount += 1
	c.TooltipOverrideCount = 0
	c.WindowsActiveCount = 0

	c.UpdateViewportsNewFrame()

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
		panic("Viewports len != 1")
	}
	// Update main viewport with current platform position.
	main_viewport := c.Viewports[0]
	main_viewport.Flags = ImGuiViewportFlags_IsPlatformWindow | ImGuiViewportFlags_OwnedByApp
	main_viewport.Pos = Vec2{0, 0}
	main_viewport.Size = c.Io.DisplaySize

	for _, v := range c.Viewports {
		v.WorkOffsetMin = v.BuildWorkOffsetMin
		v.WorkOffsetMax = v.BuildWorkOffsetMax
		v.BuildWorkOffsetMax, v.BuildWorkOffsetMin = Vec2{0, 0}, Vec2{0, 0}
		v.UpdateWorkRect()
	}
}

func removeIndex[T any](s []T, index int) []T {
	return append(s[:index], s[index+1:]...)
}
