package ui2

type ImGuiWindowFlags int

const (
	ImGuiWindowFlags_None                      ImGuiWindowFlags = 0
	ImGuiWindowFlags_NoTitleBar                ImGuiWindowFlags = 1 << iota // Disable title-bar
	ImGuiWindowFlags_NoResize                                               // Disable user resizing with the lower-right grip
	ImGuiWindowFlags_NoMove                                                 // Disable user moving the window
	ImGuiWindowFlags_NoScrollbar                                            // Disable scrollbars (window can still scroll with mouse or programmatically)
	ImGuiWindowFlags_NoScrollWithMouse                                      // Disable user vertically scrolling with mouse wheel. On child window, mouse wheel will be forwarded to the parent unless NoScrollbar is also set.
	ImGuiWindowFlags_NoCollapse                                             // Disable user collapsing window by double-clicking on it. Also referred to as Window Menu Button (e.g. within a docking node).
	ImGuiWindowFlags_AlwaysAutoResize                                       // Resize every window to its content every frame
	ImGuiWindowFlags_NoBackground                                           // Disable drawing background color (WindowBg, etc.) and outside border. Similar as using SetNextWindowBgAlpha(0.0f).
	ImGuiWindowFlags_NoSavedSettings                                        // Never load/save settings in .ini file
	ImGuiWindowFlags_NoMouseInputs                                          // Disable catching mouse, hovering test with pass through.
	ImGuiWindowFlags_MenuBar                                                // Has a menu-bar
	ImGuiWindowFlags_HorizontalScrollbar                                    // Allow horizontal scrollbar to appear (off by default). You may use SetNextWindowContentSize(ImVec2(width,0.0f)); prior to calling Begin() to specify width. Read code in imgui_demo in the "Horizontal Scrolling" section.
	ImGuiWindowFlags_NoFocusOnAppearing                                     // Disable taking focus when transitioning from hidden to visible state
	ImGuiWindowFlags_NoBringToFrontOnFocus                                  // Disable bringing window to front when taking focus (e.g. clicking on it or programmatically giving it focus)
	ImGuiWindowFlags_AlwaysVerticalScrollbar                                // Always show vertical scrollbar (even if ContentSize.y < Size.y)
	ImGuiWindowFlags_AlwaysHorizontalScrollbar                              // Always show horizontal scrollbar (even if ContentSize.x < Size.x)
	ImGuiWindowFlags_AlwaysUseWindowPadding                                 // Ensure child windows without border uses style.WindowPadding (ignored by default for non-bordered child windows, because more convenient)
	ImGuiWindowFlags_NoNavInputs                                            // No gamepad/keyboard navigation within the window
	ImGuiWindowFlags_NoNavFocus                                             // No focusing toward this window with gamepad/keyboard navigation (e.g. skipped by CTRL+TAB)
	ImGuiWindowFlags_UnsavedDocument                                        // Display a dot next to the title. When used in a tab/docking context, tab is selected when clicking the X + closure is not assumed (will wait for user to stop submitting the tab). Otherwise closure is assumed when pressing the X, so if you keep submitting the tab may reappear at end of tab bar.
	ImGuiWindowFlags_NoNav                                      = ImGuiWindowFlags_NoNavInputs | ImGuiWindowFlags_NoNavFocus
	ImGuiWindowFlags_NoDecoration                               = ImGuiWindowFlags_NoTitleBar | ImGuiWindowFlags_NoResize | ImGuiWindowFlags_NoScrollbar | ImGuiWindowFlags_NoCollapse
	ImGuiWindowFlags_NoInputs                                   = ImGuiWindowFlags_NoMouseInputs | ImGuiWindowFlags_NoNavInputs | ImGuiWindowFlags_NoNavFocus

	// [Internal]
	ImGuiWindowFlags_NavFlattened ImGuiWindowFlags = 1 << iota // [BETA] On child window: allow gamepad/keyboard navigation to cross over parent border to this child or between sibling child windows.
	ImGuiWindowFlags_ChildWindow                               // Don't use! For internal use by BeginChild()
	ImGuiWindowFlags_Tooltip                                   // Don't use! For internal use by BeginTooltip()
	ImGuiWindowFlags_Popup                                     // Don't use! For internal use by BeginPopup()
	ImGuiWindowFlags_Modal                                     // Don't use! For internal use by BeginPopupModal()
	ImGuiWindowFlags_ChildMenu                                 // Don't use! For internal use by BeginMenu()
)

type ImGuiWindow struct {
	Name                     string
	ID                       ImGuiID
	Flags                    ImGuiWindowFlags
	Viewport                 *ImGuiViewportP
	Pos                      ImVec2
	Size                     ImVec2
	SizeFull                 ImVec2
	ContentSize              ImVec2
	ContentSizeIdeal         ImVec2
	ContentSizeExplicit      ImVec2
	WindowPadding            ImVec2
	WindowRounding           float32
	WindowBorderSize         float32
	NameBufLen               int
	MoveId                   ImGuiID
	ChildId                  ImGuiID
	Scroll                   ImVec2
	ScrollMax                ImVec2
	ScrollTarget             ImVec2
	ScrollTargetCenterRatio  ImVec2
	ScrollTargetEdgeSnapDist ImVec2
	ScrollbarSizes           ImVec2

	ScrollbarX, ScrollbarY         bool
	Active                         bool
	WasActive                      bool
	WriteAccessed                  bool
	Collapsed                      bool
	WantCollapseToggle             bool
	SkipItems                      bool
	Appearing                      bool
	Hidden                         bool
	IsFallbackWindow               bool
	IsExplicitChild                bool
	HasCloseButton                 bool
	ResizeBorderHeld               byte
	BeginCount                     int
	BeginOrderWithinParent         int
	BeginOrderWithinContext        int
	FocusOrder                     int
	PopupId                        ImGuiID
	AutoFitFramesX, AutoFitFramesY byte
	AutoFitChildAxises             byte
	AutoFitOnlyGrows               bool

	IDStack []ImGuiID

	DC ImGuiWindowTempData
}

type ImGuiWindowTempData struct {
	//Layout
	CursorPos               ImVec2
	CursorPosPrevLine       ImVec2
	CursorStartPos          ImVec2
	CursorMaxPos            ImVec2
	IdealMaxPos             ImVec2
	CurrLineSize            ImVec2
	PrevLineSize            ImVec2
	CurrLineTextBaseOffset  float32
	PrevLineTextBaseOffset  float32
	IsSameLine              bool
	Indent                  float32
	ColumnsOffset           float32
	GroupOffset             float32
	CursorStartPosLossyness ImVec2

	// Keyboard/Gamepad navigation
	NavLayerCurrent          ImGuiNavLayer
	NavLayersActiveMask      int
	NavLayersActiveMaskNext  int
	NavFocusScopeIdCurrent   ImGuiID
	NavHideHighlightOneFrame bool
	NavHasScroll             bool

	// Miscellaneous
	MenuBarAppending          bool
	MenuBarOffset             ImVec2
	MenuColumns               ImGuiMenuColumns
	TreeDepth                 int
	TreeJumpToParentOnPopMask uint32
	ChildWindows              []*ImGuiWindow
	// StateStorage              *ImGuiStorage
	// CurrentColumns            *ImGuiOldColumns
	CurrentTableIdx int
	// LayoutType                ImGuiLayoutType
	// ParentLayoutType          ImGuiLayoutType
	OuterRectClipped  ImRect
	InnerRect         ImRect
	InnerClipRect     ImRect
	WorkRect          ImRect
	ParentWorkRect    ImRect
	ClipRect          ImRect
	ContentRegionRect ImRect

	DrawList                       *ImDrawList
	DrawListInst                   ImDrawList
	ParentWindow                   *ImGuiWindow
	ParentWindowInBeginStack       *ImGuiWindow
	RootWindow                     *ImGuiWindow
	RootWindowPopupTree            *ImGuiWindow
	RootWindowForTitleBarHighlight *ImGuiWindow
	RootWindowForNav               *ImGuiWindow
	NavLastChildNavWindow          *ImGuiWindow
	NavLastIds                     [ImGuiNavLayer_COUNT]ImGuiID
	NavRectRel                     [ImGuiNavLayer_COUNT]ImRect
}

type ImGuiNavLayer int

const (
	ImGuiNavLayer_Main ImGuiNavLayer = iota // Main scrolling layer
	ImGuiNavLayer_Menu                      // Menu layer (access with Alt/ImGuiNavInput_Menu)
	ImGuiNavLayer_COUNT
)

type ImGuiMenuColumns struct {
	TotalWidth     uint32
	NextTotalWidth uint32
	Spacing        uint16
	OffsetIcon     uint16
	OffsetLabel    uint16
	OffsetShortcut uint16
	OffsetMark     uint16
	Widths         [4]uint16
}

type ImGuiLastItemData struct {
	ID          ImGuiID
	InFlags     ImGuiItemFlags
	StatusFlags ImGuiItemStatusFlags
	Rect        ImRect
	NavRect     ImRect
	DisplayRect ImRect
}
type ImGuiWindowStackData struct {
	Window *ImGuiWindow     
	ParentLastItemDataBackup ImGuiLastItemData 
}

type ImGuiItemFlags int
type ImGuiItemStatusFlags int

const (
	ImGuiItemFlags_None                     ImGuiItemFlags = 0
	ImGuiItemFlags_NoTabStop                ImGuiItemFlags = 1 << 0
	ImGuiItemFlags_ButtonRepeat             ImGuiItemFlags = 1 << 1
	ImGuiItemFlags_Disabled                 ImGuiItemFlags = 1 << 2
	ImGuiItemFlags_NoNav                    ImGuiItemFlags = 1 << 3
	ImGuiItemFlags_NoNavDefaultFocus        ImGuiItemFlags = 1 << 4
	ImGuiItemFlags_SelectableDontClosePopup ImGuiItemFlags = 1 << 5
	ImGuiItemFlags_MixedValue               ImGuiItemFlags = 1 << 6
	ImGuiItemFlags_ReadOnly                 ImGuiItemFlags = 1 << 7
	ImGuiItemFlags_Inputable                ImGuiItemFlags = 1 << 8

	ImGuiItemStatusFlags_None             ImGuiItemStatusFlags = 0
	ImGuiItemStatusFlags_HoveredRect      ImGuiItemStatusFlags = 1 << 0
	ImGuiItemStatusFlags_HasDisplayRect   ImGuiItemStatusFlags = 1 << 1
	ImGuiItemStatusFlags_Edited           ImGuiItemStatusFlags = 1 << 2
	ImGuiItemStatusFlags_ToggledSelection ImGuiItemStatusFlags = 1 << 3
	ImGuiItemStatusFlags_ToggledOpen      ImGuiItemStatusFlags = 1 << 4
	ImGuiItemStatusFlags_HasDeactivated   ImGuiItemStatusFlags = 1 << 5
	ImGuiItemStatusFlags_Deactivated      ImGuiItemStatusFlags = 1 << 6
	ImGuiItemStatusFlags_HoveredWindow    ImGuiItemStatusFlags = 1 << 7
	ImGuiItemStatusFlags_FocusedByTabbing ImGuiItemStatusFlags = 1 << 8
)
