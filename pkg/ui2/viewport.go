package ui2

type ImGuiViewport struct {
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
	DrawListsLastFrame [2]int
	DrawLists          [2]*ImDrawList
	DrawDataP          ImDrawData
	DrawDataBuilder    ImDrawDataBuilder

	WorkOffsetMin      Vec2
	WorkOffsetMax      Vec2
	BuildWorkOffsetMin Vec2
	BuildWorkOffsetMax Vec2
}

func NewImGuiViewportP() *ImGuiViewportP {
	l1 := NewImDrawList()
	l2 := NewImDrawList()

	v := ImGuiViewportP{
		DrawLists: [2]*ImDrawList{l1,l2},
	}
	return &v
}
