package main

import (
	"github.com/Dmitry-dms/moon/pkg/ui2"
	"github.com/go-gl/gl/v4.2-core/gl"
	"github.com/go-gl/glfw/v3.3/glfw"
)

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


	for !window.ShouldClose(){
		glfw.PollEvents()
		gl.ClearColor(1, 1, 1, 1)
		gl.Clear(gl.COLOR_BUFFER_BIT)

		ui2.ImplOpenGL3_NewFrame()
		ui2.ImplGlfw_NewFrame(window)



		window.SwapBuffers()
	}
}
