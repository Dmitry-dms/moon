package ui2

type ImGuiViewport struct {
	Flags    ImGuiViewportFlags
	Pos      Vec2
	Size     Vec2
	WorkPos  Vec2
	WorkSize Vec2
}

func (v *ImGuiViewport) GetCenter() Vec2 {
	return Vec2{v.Pos.X + v.Size.X*0.5, v.Pos.Y + v.Size.Y*0.5}
}
func (v *ImGuiViewport) GetWorkCenter() Vec2 {
	return Vec2{v.WorkPos.X + v.WorkSize.X*0.5, v.WorkPos.Y + v.WorkSize.Y*0.5}
}

type ImGuiViewportP struct {
	*ImGuiViewport
	DrawListsLastFrame [2]int
	DrawLists          [2]*ImDrawList
	DrawDataP          ImDrawData
	DrawDataBuilder    ImDrawDataBuilder

	WorkOffsetMin      Vec2
	WorkOffsetMax      Vec2
	BuildWorkOffsetMin Vec2
	BuildWorkOffsetMax Vec2
}

func (v *ImGuiViewportP) UpdateWorkRect() {
	v.WorkPos = v.CalcWorkRectPos(v.WorkOffsetMin)
	v.WorkSize = v.CalcWorkRectSize(v.WorkOffsetMin, v.WorkOffsetMax)
}

func (v *ImGuiViewportP) CalcWorkRectPos(off_min Vec2) Vec2 {
	return Vec2{v.Pos.X + off_min.X, v.Pos.Y + off_min.Y}
}
func (v *ImGuiViewportP) CalcWorkRectSize(off_min, off_max Vec2) Vec2 {
	return Vec2{max(0.0, v.Size.X-off_min.X+off_max.X), max(0.0, v.Size.Y-off_min.Y+off_max.Y)}
}

func max[T float32](v1, v2 T) T {
	if v1 > v2 {
		return v1
	} else {
		return v2
	}
}

func NewImGuiViewportP() *ImGuiViewportP {
	l1 := NewImDrawList()
	l2 := NewImDrawList()

	v := ImGuiViewportP{
		DrawLists: [2]*ImDrawList{l1, l2},
	}
	return &v
}

type ImGuiViewportFlags int

const (
	ImGuiViewportFlags_None              ImGuiViewportFlags = 0
	ImGuiViewportFlags_IsPlatformWindow  ImGuiViewportFlags = 1 << 0
	ImGuiViewportFlags_IsPlatformMonitor ImGuiViewportFlags = 1 << 1
	ImGuiViewportFlags_OwnedByApp        ImGuiViewportFlags = 1 << 2
)
