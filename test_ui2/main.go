package main

import (

	"runtime"


	"github.com/Dmitry-dms/moon/pkg/ui2"
	"github.com/Dmitry-dms/moon/pkg/ui2/fonts"
	"github.com/go-gl/gl/v4.2-core/gl"
	"github.com/go-gl/glfw/v3.3/glfw"
	"golang.org/x/image/colornames"
)

func init() {
	runtime.LockOSThread()
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

	glfw.SwapInterval(1)

	ui2.CreateContext()
	ui2.ImGui_ImplGlfw_InitForOpenGL(window, true)
	ui2.ImGui_ImplOpenGL3_Init("42")

	font := fonts.NewFont("assets/fonts/Roboto.ttf", 60)

	batch := fonts.NewTextBatch(font)
	batch.Init()

	gl.Enable(gl.BLEND)
	gl.BlendFunc(gl.SRC_ALPHA, gl.ONE_MINUS_SRC_ALPHA)


	for !window.ShouldClose() {
		glfw.PollEvents()
		gl.ClearColor(1, 1, 1, 1)
		gl.Clear(gl.COLOR_BUFFER_BIT)

		ui2.ImplOpenGL3_NewFrame()
		ui2.ImplGlfw_NewFrame(window)

		batch.AddText("Hello world!", 200, 200, 1, colornames.Magenta)
		batch.AddText("My name is Dmitry", 100, 300, 1, colornames.Black)

		batch.FlushBatch()

		window.SwapBuffers()
	}
}

