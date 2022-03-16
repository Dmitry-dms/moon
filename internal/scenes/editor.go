package scenes

import (
	"fmt"
	"time"

	"github.com/Dmitry-dms/moon/internal/renderers"
	imgui "github.com/inkyblackness/imgui-go/v4"
)

type EditorScene struct {
	showDemoWindow bool
}

func NewEditorScene(renderer renderers.Renderer) *EditorScene {

	edtrScene := EditorScene{
		showDemoWindow: true,
	}

	return &edtrScene
}

func (e *EditorScene) Init() {

}
func (e *EditorScene) Start() {

}
func (e *EditorScene) Destroy() {

}
func (e *EditorScene) Update(dt float32) {

}
func (e *EditorScene) Render() {

}
func (e *EditorScene) Imgui() {
	// // 1. Show a simple window.
	// // Tip: if we don't call imgui.Begin()/imgui.End() the widgets automatically appears in a window called "Debug".
	{
		imgui.Begin("Test")
		imgui.Text(fmt.Sprintf("Application average %.3f ms/frame (%.1f FPS)",
			float32(time.Second.Milliseconds())/imgui.CurrentIO().Framerate(), imgui.CurrentIO().Framerate()))
		imgui.End()
	}

	// 3. Show the ImGui demo window. Most of the sample code is in imgui.ShowDemoWindow().
	// Read its code to learn more about Dear ImGui!
	if e.showDemoWindow {
		// Normally user code doesn't need/want to call this because positions are saved in .ini file anyway.
		// Here we just want to make the demo initial state a bit more friendly!
		const demoX = 650
		const demoY = 20
		imgui.SetNextWindowPosV(imgui.Vec2{X: demoX, Y: demoY}, imgui.ConditionFirstUseEver, imgui.Vec2{})
		imgui.ShowDemoWindow(&e.showDemoWindow)
	}
}
