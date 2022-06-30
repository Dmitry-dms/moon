package ui2

import (

	"github.com/pkg/errors"
)

type Vec2 struct {
	X, Y float32
}

type Vec4 struct {
	X, Y, Z, W float32
}

var context *UIContext = nil

type UIContext struct {
	Initialized      bool
	Io               *ImIO
	InputEventsQueue []ImGuiInputEvent

	MouseCursor ImGuiMouseCursor

	//new frame

	Time                 float32
	WithinFrameScope     bool
	FrameCount           uint
	WindowsActiveCount   uint
	TooltipOverrideCount float64
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
			InputEventsQueue: make([]ImGuiInputEvent, 0),
		}
		context = ctx
	}
	initializeContext()
	return context
}



func (c *UIContext) NewFrame() {
	// Remove pending delete hooks before frame start.
	// This deferred removal avoid issues of removal while iterating the hook vector

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
	
	// Update main viewport with current platform position.
	
	
}

func removeIndex[T any](s []T, index int) []T {
	return append(s[:index], s[index+1:]...)
}
