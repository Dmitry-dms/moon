package ui2

type ImVec2 struct {
	X, Y float32
}

type ImVec4 struct {
	X, Y, W, H float32
}

func FloorSigned(vec ImVec2) ImVec2 {
	return ImVec2{floatSigned(vec.X), floatSigned(vec.Y)}
}

func floatSigned(f float32) float32 {
	if f >= 0 || float32(int(f)) == f {
		return float32(int(f))
	} else {
		return float32(int(f) - 1)
	}
}

func LengthSqrVec2(v ImVec2) float32 {
	return v.X*v.X + v.Y*v.Y
}

func (v ImVec2) Sub(x ImVec2) ImVec2 {
	return ImVec2{v.X - x.X, v.Y - x.Y}
}
func (v ImVec2) Add(x ImVec2) ImVec2 {
	return ImVec2{v.X + x.X, v.Y + x.Y}
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
	Type          ImGuiContextHookType
	UserData      T
	Callback      ImGuiContextHookCallback[T]
	HookId, Owner ImGuiID
}
