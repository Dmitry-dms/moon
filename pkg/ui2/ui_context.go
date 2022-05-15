package ui2

var context *UIContext = nil

type UIContext struct {
	Initialized bool
	Io          *ImIO
	Viewports []*ImGuiViewportP
	InputEventsQueue []ImGuiInputEvent

	MouseCursor ImGuiMouseCursor

}

func GetCurrentContext() *UIContext {
	return context
}

type ImGuiMouseCursor int

const(
	ImGuiMouseCursor_None ImGuiMouseCursor =  iota
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
				DisplaySize: Vec2{},
				BackendPlatformUserData: newData(),
			},
			Viewports: make([]*ImGuiViewportP, 0),
			InputEventsQueue: make([]ImGuiInputEvent, 0),
		}
		context = ctx
	}
	initializeContext()
	return context
}
