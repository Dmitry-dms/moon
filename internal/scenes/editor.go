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

	currentGameWorld *GameWorld

	isRunning bool
}

func NewEditorScene(changeSceneCallback func(scene int)) *EditorScene {
	world := NewGameWorld("first", 20, 20)

	edtrScene := EditorScene{
		showDemoWindow:      true,
		changeSceneCallback: changeSceneCallback,
		camera:              gogl.NewCamera(mgl32.Vec2{0, 0}),
		currentGameWorld:    world,
	}

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
	e.currentGameWorld.loadResources()
	e.currentGameWorld.Init()

}

var indeces = []int32{
	3, 2, 0, 0, 2, 1,
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

var vao, vbo, ebo uint32
var f = true

func (e *EditorScene) Update(dt float32) {
	//fmt.Printf("FPS - %.1f \n", 1/dt)
	e.currentGameWorld.Update(dt)

	// for _, o := range e.gameObjects {
	// 	if o != nil {
	// 		o.Update(dt)
	// 	}
	// }

	// e.renderer.Render(*e.camera)

}

func (e *EditorScene) Render() {
	e.currentGameWorld.Render(e.camera)
	// e.shader.Use()
	// gl.DrawElements(gl.TRIANGLES, 6, gl.UNSIGNED_INT, gl.PtrOffset(0))
	// err := e.shader.CheckShaderForChanges()
	// if err != nil {
	// 	fmt.Println(err)
	// }
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
