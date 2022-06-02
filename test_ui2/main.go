package main

import (
	"fmt"
	"runtime"

	"github.com/Dmitry-dms/moon/pkg/gogl"
	"github.com/Dmitry-dms/moon/pkg/ui2"
	"github.com/Dmitry-dms/moon/pkg/ui2/fonts"
	"github.com/go-gl/gl/v4.2-core/gl"
	"github.com/go-gl/glfw/v3.3/glfw"
)

func init() {
	runtime.LockOSThread()
}

var vertices = []float32{
	0.5, 0.5, 1.0, 0.2, 0.11, 1.0, 0.0,
	0.5, -0.5, 1.0, 0.2, 0.11, 1.0, 1.0,
	-0.5, -0.5, 1.0, 0.2, 0.11, 0.0, 1.0,
	-0.5, 0.5, 1.0, 0.2, 0.11, 0.0, 0.0,
}
var indeces = []int32{
	0, 1, 3,
	1, 2, 3,
}

func main() {
	err := glfw.Init()
	if err != nil {
		panic(err)
	}
	glfw.DefaultWindowHints()

	window, err := glfw.CreateWindow(1280, 720, "Dear ImGui GLFW+OpenGL3 example", nil, nil)
	if err != nil {
		panic(err)
	}
	window.MakeContextCurrent()
	//включаем верт. синхронизацию
	glfw.SwapInterval(1)

	ui2.CreateContext()
	ui2.ImGui_ImplGlfw_InitForOpenGL(window, true)
	ui2.ImGui_ImplOpenGL3_Init("42")

	// font, err := ui2.LoadFontFromFile("assets/fonts/mono.ttf", 30)
	// if err != nil {
	// 	panic(err)
	// }
	// ctx.DefaultFont = font

	//font

	fontShader, _ := gogl.NewShader("assets/shaders/fonts.glsl")

	font := fonts.NewFont("assets/fonts/Roboto.ttf", 60)


	texcCoords := font.GetCharacter('&').TexCoords
	vertices[5], vertices[6] = texcCoords[0].X, texcCoords[0].Y
	vertices[12], vertices[13] = texcCoords[1].X, texcCoords[1].Y
	vertices[19], vertices[20] = texcCoords[2].X, texcCoords[2].Y
	vertices[26], vertices[27] = texcCoords[3].X, texcCoords[3].Y

	fmt.Printf("coords = %v", texcCoords)
	initSquare()


	for !window.ShouldClose() {
		glfw.PollEvents()
		gl.ClearColor(1, 1, 1, 1)
		gl.Clear(gl.COLOR_BUFFER_BIT)

		ui2.ImplOpenGL3_NewFrame()
		ui2.ImplGlfw_NewFrame(window)

		fontShader.Use()

		gl.ActiveTexture(gl.TEXTURE0)
		gl.BindTexture(gl.TEXTURE_BUFFER, font.TextureId)
		fontShader.UploadTexture("uFontTexture", 0)


		gl.BindVertexArray(vao)


		gl.DrawElements(gl.TRIANGLES, 6, gl.UNSIGNED_INT, nil)

		window.SwapBuffers()
	}
}

var (
	vbo, vao uint32
)

func initSquare() {
	vao = gogl.GenBindVAO()
	// vbo = gogl.GenBindBuffer(gl.ARRAY_BUFFER)
	// gogl.BufferData(gl.ARRAY_BUFFER, vertices, gl.STATIC_DRAW)

	// ebo := gogl.GenEBO()
	// gl.BindBuffer(gl.ELEMENT_ARRAY_BUFFER, ebo)
	// gogl.BufferData(gl.ELEMENT_ARRAY_BUFFER, indeces, gl.STATIC_DRAW)

	// stride := 7 * 4
	// gogl.SetVertexAttribPointer(0, 2, gl.FLOAT, stride, 0)
	// gogl.SetVertexAttribPointer(1, 3, gl.FLOAT, stride, 2)
	// gogl.SetVertexAttribPointer(2, 2, gl.FLOAT, stride, 5)

	// gl.GenVertexArrays(1, &vao)
	// gl.BindVertexArray(vao)

	gl.GenBuffers(1, &vbo)
	gl.BindBuffer(gl.ARRAY_BUFFER, vbo)
	gl.BufferData(gl.ARRAY_BUFFER, 4*28, gl.Ptr(vertices), gl.STATIC_DRAW)

	var ebo uint32
	gl.GenBuffers(1, &ebo)
	gl.BindBuffer(gl.ELEMENT_ARRAY_BUFFER, ebo)
	gl.BufferData(gl.ELEMENT_ARRAY_BUFFER, 4*6, gl.Ptr(indeces), gl.STATIC_DRAW)

	stride := 7 * 4
	gl.VertexAttribPointer(0, 2, gl.FLOAT, false, int32(stride), gl.PtrOffset(0))
	gl.EnableVertexAttribArray(0)

	gl.VertexAttribPointer(1, 3, gl.FLOAT, false, int32(stride), gl.PtrOffset(2*4))
	gl.EnableVertexAttribArray(1)

	gl.VertexAttribPointer(2, 2, gl.FLOAT, false, int32(stride), gl.PtrOffset(5*4))
	gl.EnableVertexAttribArray(2)

}