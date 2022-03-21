package scenes

import (
	"fmt"
	"time"
	// "unsafe"

	"github.com/Dmitry-dms/moon/internal/components"
	"github.com/Dmitry-dms/moon/pkg/gogl"

	// "github.com/go-gl/gl/v4.2-core/gl"
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

	shader    *gogl.Shader
	isRunning bool
}

func NewEditorScene(changeSceneCallback func(scene int)) *EditorScene {
	world := NewGameWorld("first", 20, 20)

	edtrScene := EditorScene{
		showDemoWindow:      true,
		changeSceneCallback: changeSceneCallback,
		camera:              gogl.NewCamera(mgl32.Vec2{0, 0}),
		//camera:              gogl.NewCamera(mgl32.Vec3{0, 0, 0}, mgl32.Vec3{0, 1, 0}, 0.5, 0.1),
		currentGameWorld: world,
	}

	return &edtrScene
}

func (e *EditorScene) GetCamera() *gogl.Camera {
	return e.camera
}
func generateIndeces() []int32 {
	//6 вершин на 1 квадрат
	elements := make([]int32, 6*maxBatchSize)

	for i := 0; i < maxBatchSize; i++ {
		loadElementIndeces(elements, int32(i))
	}

	return elements
}
func loadElementIndeces(elements []int32, index int32) {
	var offsetArrayIndex int32 = 6 * index
	var offset int32 = 4 * index
	// 3, 2, 0, 0, 2, 1 - вершины квадрата
	elements[offsetArrayIndex] = offset + 3
	elements[offsetArrayIndex+1] = offset + 2
	elements[offsetArrayIndex+2] = offset + 0

	elements[offsetArrayIndex+3] = offset + 0
	elements[offsetArrayIndex+4] = offset + 2
	elements[offsetArrayIndex+5] = offset + 1
}

func (e *EditorScene) Init() {

	fmt.Println("init editor scene")
	e.currentGameWorld.Init()
	

	// s, _ := gogl.NewShader("assets/shaders/default.glsl")
	// e.shader = s
	// vao = gogl.GenBindVAO()

	// vbo = gogl.GenBindBuffer(gl.ARRAY_BUFFER)
	// gl.BufferData(gl.ARRAY_BUFFER, 1000*6*4, nil, gl.DYNAMIC_DRAW)


	// gogl.SetVertexAttribPointer(0, 2, gl.FLOAT, 6, 0)
	// gogl.SetVertexAttribPointer(1, 4, gl.FLOAT, 6, 2)

	// ebo = gogl.GenBindBuffer(gl.ELEMENT_ARRAY_BUFFER)

	// indeces = generateIndeces()
	
	// gogl.BufferData(gl.ELEMENT_ARRAY_BUFFER, indeces, gl.STATIC_DRAW)

	// gl.BindBuffer(gl.ELEMENT_ARRAY_BUFFER, 0)

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
	// fmt.Printf("COMPONENT = %#v\n",obj)
	// if e.isRunning {
	// 	e.renderer.Add(obj)
	// } else {
	// 	e.renderer.Add(obj)
	// 	obj.Start()
	// }

}
func (e *EditorScene) Destroy() {

}

var vao, vbo, ebo uint32
var f = true

func (e *EditorScene) Update(dt float32) {

	e.currentGameWorld.Update(dt)
	// var vertices = []float32{
	// 	-0.5, -0.5, 1, 0, 0, 1,
	// 	-0.5, 0.5, 1, 1, 0, 1,
	// 	0.5, 0.5, 1, 0, 1, 1,
	// 	0.5, -0.5, 0, 1, 1, 1,
	// }
	// gl.BindBuffer(gl.ARRAY_BUFFER, vbo)
	// gl.BufferSubData(gl.ARRAY_BUFFER, 0, len(vertices)*4, gl.Ptr(vertices))
	// gl.BindBuffer(gl.ARRAY_BUFFER, 0)

	// e.shader.Use()
	// e.shader.UploadMat4("uProjection", e.camera.GetProjectionMatrix())
	// e.shader.UploadMat4("uView", e.camera.GetViewMatrix())

	// gl.BindVertexArray(vao)

	//  gl.BindBuffer(gl.ELEMENT_ARRAY_BUFFER, ebo)
	// gl.DrawElements(gl.TRIANGLES, 6, gl.UNSIGNED_INT, nil)
	// gl.BindBuffer(gl.ELEMENT_ARRAY_BUFFER, 0)

	//  gl.BindVertexArray(0)
	// e.shader.Detach()

	// for _, o := range e.gameObjects {
	// 	if o != nil {
	// 		o.Update(dt)
	// 	}
	// }

	// e.renderer.Render(*e.camera)

	//=============================

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
