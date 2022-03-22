package scenes

import (
	"fmt"
	"time"
	// "unsafe"

	"github.com/Dmitry-dms/moon/internal/components"
	"github.com/Dmitry-dms/moon/pkg/gogl"

	//  "github.com/go-gl/gl/v4.2-core/gl"
	"github.com/go-gl/glfw/v3.3/glfw"
	"github.com/go-gl/mathgl/mgl32"

	//	mgl "github.com/go-gl/mathgl/mgl32"
	imgui "github.com/inkyblackness/imgui-go/v4"
)

type EditorScene struct {
	showDemoWindow bool
	glfw           *glfw.Window

	changeSceneCallback func(scene int)
	camera              *gogl.Camera

	activeGameWorld  *GameWorld
	activeGameObject *components.GameObject

	isRunning bool
}

func NewEditorScene(changeSceneCallback func(scene int)) *EditorScene {
	edtrScene := EditorScene{
		showDemoWindow:      true,
		changeSceneCallback: changeSceneCallback,
		camera:              gogl.NewCamera(mgl32.Vec2{0, 0}),
	}
	callback := func(g *components.GameObject) {
		edtrScene.activeGameObject = g
	}
	world := NewGameWorld("first", 20, 20, callback)
	edtrScene.activeGameWorld = world

	return &edtrScene
}

func (e *EditorScene) GetCamera() *gogl.Camera {
	return e.camera
}

func (e *EditorScene) loadResources() {
	gogl.AssetPool.GetShader("assets/shaders/default.glsl")

	gogl.AssetPool.AddSpritesheet("assets/images/spritesheet.png",
		gogl.NewSpritesheet(gogl.AssetPool.GetTexture("assets/images/spritesheet.png"), 16, 16, 26, 0))
}

func (e *EditorScene) Init() {

	fmt.Println("init editor scene")
	e.loadResources()
	e.activeGameWorld.loadResources()
	e.activeGameWorld.Init()

}

// при запуске сцены, запускаем объекты
func (e *EditorScene) Start() {
	// for _, g := range e.gameObjects {
	// 	if g != nil {
	// 		g.Start()
	// 		e.renderer.Add(g)
	// 	}
	// }
	// e.isRunning = true
}

var inc int

func (e *EditorScene) AddGameObjectToScene(obj *components.GameObject) {

}
func (e *EditorScene) Destroy() {

}

func (e *EditorScene) Update(dt float32) {
	//fmt.Printf("FPS - %.1f \n", 1/dt)
	e.activeGameWorld.Update(dt)
}

func (e *EditorScene) Render() {
	e.activeGameWorld.Render(e.camera)
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

	{
		if e.activeGameObject != nil {
			imgui.Begin("Inspector")
			e.activeGameObject.Imgui()
			imgui.End()
		}
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
