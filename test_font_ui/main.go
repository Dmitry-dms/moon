package main

import (

	"runtime"


	"github.com/Dmitry-dms/moon/pkg/fonts"
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
	// glfw.WindowHint(glfw.OpenGLDebugContext, 1)

	window, err := glfw.CreateWindow(1280, 720, "example test", nil, nil)
	if err != nil {
		panic(err)
	}
	window.MakeContextCurrent()

	glfw.SwapInterval(1)

	gl.Init()

	size := func(w *glfw.Window, width int, height int) {
		gl.Viewport(0, 0, int32(width), int32(height))
	}
	window.SetSizeCallback(size)

	// var sizeTex int32
	// gl.GetIntegerv(gl.MAX_COMBINED_TEXTURE_IMAGE_UNITS, &sizeTex)
	// fmt.Println(sizeTex)

	// font := fonts.NewFont("assets/fonts/rany.otf", 60,true)
	// font := fonts.NewFont("assets/fonts/mono.ttf", 60,true)
	// font := fonts.NewFont("assets/fonts/Roboto.ttf", 60,true)
	font := fonts.NewFont("C:/Windows/Fonts/times.ttf", 40, true)

	batch := fonts.NewTextBatch(font)
	batch.Init()

	// gl.Enable(gl.BLEND)
	// gl.BlendFunc(gl.SRC_ALPHA, gl.ONE_MINUS_SRC_ALPHA)

	for !window.ShouldClose() {
		glfw.PollEvents()
		gl.ClearColor(1, 1, 1, 1)
		gl.Clear(gl.COLOR_BUFFER_BIT)

		// ui2.ImplOpenGL3_NewFrame()
		// ui2.ImplGlfw_NewFrame(window)

		batch.AddText("Привет, мир!\n920043 ~hghguij Progress #$@\n[A-Za-z] {0-9_20-33}", 50, 600, 1, colornames.Black)
		// batch.AddText("My name is Dmitry", 100, 340, 1, colornames.Magenta)

		batch.FlushBatch()

		window.SwapBuffers()
	}
}

