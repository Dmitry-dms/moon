package ui2

type ImGuiViewport struct {
	Flags    ImGuiViewportFlags
	Pos      ImVec2
	Size     ImVec2
	WorkPos  ImVec2
	WorkSize ImVec2
}

func (v *ImGuiViewport) GetCenter() ImVec2 {
	return ImVec2{v.Pos.X + v.Size.X*0.5, v.Pos.Y + v.Size.Y*0.5}
}
func (v *ImGuiViewport) GetWorkCenter() ImVec2 {
	return ImVec2{v.WorkPos.X + v.WorkSize.X*0.5, v.WorkPos.Y + v.WorkSize.Y*0.5}
}

type ImGuiViewportP struct {
	*ImGuiViewport
	DrawListsLastFrame [2]int
	DrawLists          [2]*ImDrawList
	DrawDataP          ImDrawData
	DrawDataBuilder    ImDrawDataBuilder

	WorkOffsetMin      ImVec2
	WorkOffsetMax      ImVec2
	BuildWorkOffsetMin ImVec2
	BuildWorkOffsetMax ImVec2
}

func (v *ImGuiViewportP) UpdateWorkRect() {
	v.WorkPos = v.CalcWorkRectPos(v.WorkOffsetMin)
	v.WorkSize = v.CalcWorkRectSize(v.WorkOffsetMin, v.WorkOffsetMax)
}

func (v *ImGuiViewportP) CalcWorkRectPos(off_min ImVec2) ImVec2 {
	return ImVec2{v.Pos.X + off_min.X, v.Pos.Y + off_min.Y}
}
func (v *ImGuiViewportP) CalcWorkRectSize(off_min, off_max ImVec2) ImVec2 {
	return ImVec2{max(0.0, v.Size.X-off_min.X+off_max.X), max(0.0, v.Size.Y-off_min.Y+off_max.Y)}
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
		ImGuiViewport: &ImGuiViewport{},
		DrawLists:     [2]*ImDrawList{l1, l2},
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
