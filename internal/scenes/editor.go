package scenes

import (
	"fmt"
	"time"

	"github.com/Dmitry-dms/moon/internal/renderers"
	"github.com/Dmitry-dms/moon/pkg/gogl"
	"github.com/go-gl/gl/v4.2-core/gl"
	"github.com/go-gl/glfw/v3.3/glfw"

	//	ma "github.com/go-gl/mathgl/mgl32"
	imgui "github.com/inkyblackness/imgui-go/v4"
)

type EditorScene struct {
	showDemoWindow bool
	glfw           *glfw.Window
	tsh            *gogl.Shader
	texture        gogl.TextureID
}

func NewEditorScene(renderer renderers.Renderer, window *glfw.Window) *EditorScene {
	edtrScene := EditorScene{
		showDemoWindow: true,
		glfw:           window,
	}

	return &edtrScene
}


var vao gogl.BufferID
//var triangleShader *gogl.Shader

var vertices = []float32{
	//pos                //uv coords
	0.5, 0.5, 0.0, 1.0, 1.0,
	0.5, -0.5, 0.0, 1., 0.,
	-0.5, -0.5, 0.0, 0., 0.,
	-0.5, 0.5, 0.0, 0., 1.,
}
var indices = []int32{
	0, 1, 3,
	1, 2, 3,
}

func (e *EditorScene) Init() {

	fmt.Println("init editor scene")
	triangleShader, err := gogl.NewShader("assets/triangle.vert", "assets/quadtexture.frag")
	if err != nil {
		panic(err)
	}
	e.tsh = triangleShader

	texture := gogl.LoadTextureAlpha("assets/images/img.png")
	e.texture = texture

	gogl.GenBindBuffer(gl.ARRAY_BUFFER)//vbo
	vao = gogl.GenBindVAO()//vao
	gogl.BufferData(gl.ARRAY_BUFFER, vertices, gl.STATIC_DRAW)

	gogl.GenBindBuffer(gl.ELEMENT_ARRAY_BUFFER)//ebo
	gogl.BufferData(gl.ELEMENT_ARRAY_BUFFER, indices, gl.STATIC_DRAW)

	// 0 - начало, 3 - размер
	gl.VertexAttribPointer(0, 3, gl.FLOAT, false, 5*4, nil)
	gl.EnableVertexAttribArray(0)
	gl.VertexAttribPointer(1, 2, gl.FLOAT, false, 5*4, gl.PtrOffset(3*4))
	gl.EnableVertexAttribArray(1)

	gl.BindVertexArray(0)

}
func (e *EditorScene) Start() {

}
func (e *EditorScene) Destroy() {

}
func (e *EditorScene) Update(dt float32) {

}

func (e *EditorScene) Render() {
	e.tsh.Use()

	gogl.BindTexture(e.texture)
	gogl.BindVertexArray(vao)

	//gl.DrawArrays(gl.TRIANGLES, 0, 3)
	gl.DrawElements(gl.TRIANGLES, 6, gl.UNSIGNED_INT, gl.PtrOffset(0))

	err := e.tsh.CheckShaderForChanges()
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
