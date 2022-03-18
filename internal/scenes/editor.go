package scenes

import (
	"fmt"
	"time"

	"github.com/Dmitry-dms/moon/internal/components"
	"github.com/Dmitry-dms/moon/pkg/gogl"
	"github.com/go-gl/gl/v4.2-core/gl"
	"github.com/go-gl/glfw/v3.3/glfw"
	"github.com/go-gl/mathgl/mgl32"

	//	mgl "github.com/go-gl/mathgl/mgl32"
	imgui "github.com/inkyblackness/imgui-go/v4"
)

type EditorScene struct {
	showDemoWindow      bool
	glfw                *glfw.Window
	shader              *gogl.Shader
	texture             *gogl.Texture
	changeSceneCallback func(scene int)
	camera              *gogl.Camera
	gameObjects         []components.GameObject
	isRunning           bool
}

func NewEditorScene(changeSceneCallback func(scene int)) *EditorScene {
	edtrScene := EditorScene{
		showDemoWindow:      true,
		changeSceneCallback: changeSceneCallback,
		camera:              gogl.NewCamera(mgl32.Vec2{}),
		gameObjects:         make([]components.GameObject, 0),
	}

	return &edtrScene
}

var vao uint32
var testObj *components.GameObject

//var triangleShader *gogl.Shader

var vertices = []float32{
	//pos                //color    //uv
	100, 0, 0, 1, 0, 0, 1, 1, 1, //bottom right 0
	0, 100, 0, 0, 1, 0, 1, 0, 0, //top left 1
	100, 100, 0, 0, 0, 1, 1, 1, 0, //top right 2
	0, 0, 0, 1, 1, 0, 1, 0, 1, //bottom left 3
}
var indices = []int32{
	2, 1, 0,
	0, 1, 3,
}

func (e *EditorScene) Init() {

	fmt.Println("init editor scene")

	testObj = components.NewGameObject("test object", components.Transform{})
	testObj.AddComponent(&components.SpriteRenderer{})
	e.AddGameObjectToScene(testObj)

	shader, err := gogl.NewShader("assets/shaders/default.glsl")
	if err != nil {
		panic(err)
	}
	e.shader = shader

	texture := gogl.LoadTextureAlpha("assets/images/goomba.png")
	e.texture = texture

	gogl.GenBindBuffer(gl.ARRAY_BUFFER) //vbo
	vao = gogl.GenBindVAO()             //vao
	gogl.BufferData(gl.ARRAY_BUFFER, vertices, gl.STATIC_DRAW)

	gogl.GenBindBuffer(gl.ELEMENT_ARRAY_BUFFER) //ebo
	gogl.BufferData(gl.ELEMENT_ARRAY_BUFFER, indices, gl.STATIC_DRAW)

	// // 0 - начало, 3 - размер
	gogl.SetVertexAttribPointer(0, 3, gl.FLOAT, 9, 0)
	gogl.SetVertexAttribPointer(1, 4, gl.FLOAT, 9, 3)
	gogl.SetVertexAttribPointer(2, 2, gl.FLOAT, 9, 7)

	//gl.BindVertexArray(0)

}

// при запуске сцены, запускаем объекты
func (e *EditorScene) Start() {
	for _, g := range e.gameObjects {
		g.Start()
	}
	e.isRunning = true
}

func (e *EditorScene) GetCamera() *gogl.Camera {
	return e.camera
}

func (e *EditorScene) AddGameObjectToScene(obj *components.GameObject) {
	if e.isRunning {
		e.gameObjects = append(e.gameObjects, *obj)
	} else {
		e.gameObjects = append(e.gameObjects, *obj)
		obj.Start()
	}
}
func (e *EditorScene) Destroy() {

}

var x, y float32
var f bool

func (e *EditorScene) Update(dt float32) {

	//e.camera.SetPosition(mgl32.Vec2{float32(-dt*50)})

	//fmt.Printf("%.1f FPS \n", 1.0/dt)
	e.shader.Use()

	e.shader.UploadTexture("uTexture", 0)
	gl.ActiveTexture(gl.TEXTURE0)
	e.texture.Bind()

	e.shader.UploadMat4("uProjection", e.camera.GetProjectionMatrix())
	e.shader.UploadMat4("uView", e.camera.GetViewMatrix())

	gogl.BindVertexArray(vao)

	// gl.EnableVertexAttribArray(0)
	// gl.EnableVertexAttribArray(1)
	gl.DrawElements(gl.TRIANGLES, int32(len(indices)), gl.UNSIGNED_INT, gl.PtrOffset(0))
	// gl.DisableVertexAttribArray(0)
	// gl.DisableVertexAttribArray(1)
	e.texture.Unbind()
	gogl.BindVertexArray(0)
	e.shader.Detach()

	for _, o := range e.gameObjects {
		o.Update(dt)
	}
}

func (e *EditorScene) Render() {
	e.shader.Use()
	//w, h := e.glfw.GetSize()
	// aspRatio := float32(w) / float32(h)
	// projMatrix := mgl.Perspective(mgl.DegToRad(45), aspRatio, 0.1, 100)

	// viewMatrix := mgl.Ident4()

	// e.shader.SetMat4("projection", projMatrix)
	// e.shader.SetMat4("view", viewMatrix)

	//	gogl.BindTexture(e.texture)
	gogl.BindVertexArray(vao)

	//gl.DrawArrays(gl.TRIANGLES, 0, 3)
	gl.DrawElements(gl.TRIANGLES, 6, gl.UNSIGNED_INT, gl.PtrOffset(0))

	err := e.shader.CheckShaderForChanges()
	if err != nil {
		fmt.Println(err)
	}
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
