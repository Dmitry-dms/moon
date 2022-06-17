package ui2

type ImDrawCmd struct {
	ClipRect                        Vec4
	TextureId                       TextureId
	VtxOffset, IdxOffset, ElemCount uint32
	UserCallback                    ImDrawCallback
}

type ImDrawCallback func(parentList *ImDrawList, cmd ImDrawCmd)

type TextureId int32

type ImDrawIdx int16

type ImDrawVert struct {
	pos, uv ImVec2
	col     uint32
}

type ImDrawCmdHeader struct {
	ClipRect  Vec4
	TextureId TextureId
	VtxOffset uint
}

type ImDrawList struct {
	CmdBuffer []ImDrawCmd
	IdxBuffer []ImDrawIdx
	VtxBuffer []ImDrawVert

	vtxCurrentIdx uint
	ownerName     string

	clipRectStack  []Vec4
	textureIdStack []int32
	path           []ImVec2

	cmdHeader ImDrawCmdHeader
}

func NewImDrawList() *ImDrawList {
	l := ImDrawList{
		CmdBuffer:      make([]ImDrawCmd, 0),
		IdxBuffer:      make([]ImDrawIdx, 0),
		VtxBuffer:      make([]ImDrawVert, 0),
		clipRectStack:  make([]Vec4, 0),
		textureIdStack: make([]int32, 0),
		path:           make([]ImVec2, 0),
	}
	return &l
}

func (l *ImDrawList) AddRect(p_min, p_max ImVec2, col uint32) {

}

type ImDrawData struct {
	Valid                                       bool
	CmdListsCount, TotalIdxCount, TotalVtxCount int
	CmdLists                                    [2]*ImDrawList
	DisplayPos, DisplaySize, FramebufferScale   ImVec2
}

type ImDrawDataBuilder struct {
	Layers []*ImDrawList
}

//TODO check this
func (d *ImDrawData) Clear() {
	n := &ImDrawData{}
	d = n
}
