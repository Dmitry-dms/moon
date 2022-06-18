package main

import (
	"runtime"



	"github.com/go-gl/gl/v4.2-core/gl"
	"github.com/go-gl/glfw/v3.3/glfw"

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
	glfw.WindowHint(glfw.OpenGLDebugContext, 1)

	window, err := glfw.CreateWindow(1280, 720, "example", nil, nil)
	if err != nil {
		panic(err)
	}
	window.MakeContextCurrent()

	glfw.SwapInterval(1)


	size := func(w *glfw.Window, width int, height int) {
		gl.Viewport(0, 0, int32(width), int32(height))
	}
	window.SetSizeCallback(size)

	// gogl.InitGLdebug()

	// font := fonts.NewFont("assets/fonts/rany.otf", 60)
	// font := fonts.NewFont("assets/fonts/mono.ttf", 60)
	// font := fonts.NewFont("assets/fonts/Roboto.ttf", 60)

	// ctx.Io.DefaultFont, _ = ui2.LoadFontFromFile("C:/Windows/Fonts/times.ttf", 40)

	// font := fonts.NewFont("C:/Windows/Fonts/times.ttf", 40, true)

	// batch := fonts.NewTextBatch(font)
	// batch.Init()


	gl.Enable(gl.BLEND)
	gl.BlendFunc(gl.SRC_ALPHA, gl.ONE_MINUS_SRC_ALPHA)


	

	beginTime := float32(glfw.GetTime())
	var endTime float32
	var dt float32
	dt = dt

	for !window.ShouldClose() {
		glfw.PollEvents()
		gl.ClearColor(1, 1, 1, 1)
		gl.Clear(gl.COLOR_BUFFER_BIT)
	

	

		// batch.AddText("Привет, мир!\n920043 ~hghguij Progress #$@\n[A-Za-z] {0-9_20-33}", 50, 600, 1, colornames.Black)
		// batch.AddText("My name is Dmitry", 100, 340, 1, colornames.Magenta)

		// batch.FlushBatch()

		window.SwapBuffers()
		endTime = float32(glfw.GetTime())
		dt = endTime - beginTime
		beginTime = endTime
	}
}
