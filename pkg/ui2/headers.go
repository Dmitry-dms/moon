package ui2

type Vec2 struct {
	X, Y float32
}

type Vec4 struct {
	X, Y, Z, W float32
}

//-----------------------------------------------------------------------------
// [SECTION] Generic context hooks
//-----------------------------------------------------------------------------

type ImGuiContextHookType int

const (
	ImGuiContextHookType_NewFramePre ImGuiContextHookType = iota
	ImGuiContextHookType_NewFramePost
	ImGuiContextHookType_EndFramePre
	ImGuiContextHookType_EndFramePost
	ImGuiContextHookType_RenderPre
	ImGuiContextHookType_RenderPost
	ImGuiContextHookType_Shutdown
	ImGuiContextHookType_PendingRemoval_
)
type ImGuiContextHookCallback[T any] func(ctx *UIContext, hook *ImGuiContextHook[T])
type ImGuiID uint32

type ImGuiContextHook[T any] struct {
	Type ImGuiContextHookType
	UserData T
	Callback ImGuiContextHookCallback[T]
	HookId, Owner ImGuiID
}
