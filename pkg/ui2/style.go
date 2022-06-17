package ui2

type ImGuiStyle struct {
	Alpha                    float32
	DisabledAlpha            float32
	WindowPadding            ImVec2
	WindowRounding           float32
	WindowBorderSize         float32
	WindowMinSize            ImVec2
	WindowTitleAlign         ImVec2
	WindowMenuButtonPosition ImGuiDir
	ChildRounding            float32
	ChildBorderSize          float32
	PopupRounding            float32
	PopupBorderSize          float32
	FramePadding             ImVec2
	FrameRounding            float32
	FrameBorderSize          float32
	ItemSpacing              ImVec2
	ItemInnerSpacing         ImVec2
	CellPadding              ImVec2
	TouchExtraPadding        ImVec2

	IndentSpacing              float32
	ColumnsMinSpacing          float32
	ScrollbarSize              float32
	ScrollbarRounding          float32
	GrabMinSize                float32
	GrabRounding               float32
	LogSliderDeadzone          float32
	TabRounding                float32
	TabBorderSize              float32
	TabMinWidthForCloseButton  float32
	ColorButtonPosition        ImGuiDir
	ButtonTextAlign            ImVec2
	SelectableTextAlign        ImVec2
	DisplayWindowPadding       ImVec2
	DisplaySafeAreaPadding     ImVec2
	MouseCursorScale           float32
	AntiAliasedLines           bool
	AntiAliasedLinesUseTex     bool
	AntiAliasedFill            bool
	CurveTessellationTol       float32
	CircleTessellationMaxError float32
	Colors                     [ImGuiCol_COUNT]ImVec4
}

type ImGuiDir int
type ImGuiCol int

const (
	ImGuiDir_None ImGuiDir = -1
	ImGuiDir_Left ImGuiDir = iota
	ImGuiDir_Right
	ImGuiDir_Up
	ImGuiDir_Down
	ImGuiDir_COUNT
)

const (
	ImGuiCol_Text ImGuiCol = iota
    ImGuiCol_TextDisabled
    ImGuiCol_WindowBg          
    ImGuiCol_ChildBg              
    ImGuiCol_PopupBg              
    ImGuiCol_Border
    ImGuiCol_BorderShadow
    ImGuiCol_FrameBg              
    ImGuiCol_FrameBgHovered
    ImGuiCol_FrameBgActive
    ImGuiCol_TitleBg
    ImGuiCol_TitleBgActive
    ImGuiCol_TitleBgCollapsed
    ImGuiCol_MenuBarBg
    ImGuiCol_ScrollbarBg
    ImGuiCol_ScrollbarGrab
    ImGuiCol_ScrollbarGrabHovered
    ImGuiCol_ScrollbarGrabActive
    ImGuiCol_CheckMark
    ImGuiCol_SliderGrab
    ImGuiCol_SliderGrabActive
    ImGuiCol_Button
    ImGuiCol_ButtonHovered
    ImGuiCol_ButtonActive
    ImGuiCol_Header              
    ImGuiCol_HeaderHovered
    ImGuiCol_HeaderActive
    ImGuiCol_Separator
    ImGuiCol_SeparatorHovered
    ImGuiCol_SeparatorActive
    ImGuiCol_ResizeGrip
    ImGuiCol_ResizeGripHovered
    ImGuiCol_ResizeGripActive
    ImGuiCol_Tab
    ImGuiCol_TabHovered
    ImGuiCol_TabActive
    ImGuiCol_TabUnfocused
    ImGuiCol_TabUnfocusedActive
    ImGuiCol_PlotLines
    ImGuiCol_PlotLinesHovered
    ImGuiCol_PlotHistogram
    ImGuiCol_PlotHistogramHovered
    ImGuiCol_TableHeaderBg        
    ImGuiCol_TableBorderStrong    
    ImGuiCol_TableBorderLight     
    ImGuiCol_TableRowBg           
    ImGuiCol_TableRowBgAlt        
    ImGuiCol_TextSelectedBg
    ImGuiCol_DragDropTarget
    ImGuiCol_NavHighlight          
    ImGuiCol_NavWindowingHighlight 
    ImGuiCol_NavWindowingDimBg     
    ImGuiCol_ModalWindowDimBg      
    ImGuiCol_COUNT
)
