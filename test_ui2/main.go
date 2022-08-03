package main

import (
	"fmt"
	"image"

	"os"
	"runtime"

	"github.com/Dmitry-dms/moon/pkg/gogl"
	"github.com/Dmitry-dms/moon/pkg/ui"
	"github.com/Dmitry-dms/moon/pkg/ui/render"
	"github.com/go-gl/gl/v4.2-core/gl"
	"github.com/go-gl/glfw/v3.3/glfw"
)

func init() {
	runtime.LockOSThread()
}


var uiCtx *ui.UiContext
var window *glfw.Window
var Width, Height int = 1280, 720
var steps int = 1

func main() {
	err := glfw.Init()
	if err != nil {
		panic(err)
	}
	glfw.DefaultWindowHints()
	// glfw.WindowHint(glfw.OpenGLDebugContext, 1)

	window, err = glfw.CreateWindow(Width, Height, "example", nil, nil)
	if err != nil {
		panic(err)
	}
	window.MakeContextCurrent()

	glfw.SwapInterval(0)

	size := func(w *glfw.Window, width int, height int) {
		gl.Viewport(0, 0, int32(width), int32(height))
		fmt.Println(width, height)
		Width = width
		Height = height
		// uiCtx.Io().SetDisplaySize(float32(width), float32(height))
	}

	window.SetSizeCallback(size)
	window.SetKeyCallback(onKey)
	window.SetCursorPosCallback(cursorPosCallback)
	window.SetMouseButtonCallback(mouseBtnCallback)
	window.SetScrollCallback(scrollCallback)

	// gogl.InitGLdebug()

	// font := fonts.NewFont("assets/fonts/rany.otf", 60)
	// font := fonts.NewFont("assets/fonts/mono.ttf", 60)
	// font := fonts.NewFont("assets/fonts/Roboto.ttf", 60)

	// ctx.Io.DefaultFont, _ = ui2.LoadFontFromFile("C:/Windows/Fonts/times.ttf", 40)

	// font := fonts.NewFont("C:/Windows/Fonts/times.ttf", 40, true)

	// batch := fonts.NewTextBatch(font)
	// batch.Init()
	gl.Init()

	uiCtx = ui.UiCtx
	

	front := render.NewGlRenderer()
	uiCtx.Initialize(front)
	// uiCtx = ui.NewContext(front, cam)
	uiCtx.Io().SetCursor = setCursor

	gl.Enable(gl.BLEND)
	gl.BlendFunc(gl.SRC_ALPHA, gl.ONE_MINUS_SRC_ALPHA)

	ui.UiCtx.UploadFont("C:/Windows/Fonts/times.ttf", 50)
	// ui.UiCtx.UploadFont("assets/fonts/rany.otf", 50)

	beginTime := float32(glfw.GetTime())
	var endTime float32
	var dt float32
	dt = dt

	var time float32 = 0
	tex, _ = tex.Init("assets/images/mario.png")
	tex2, _ = tex2.Init("assets/images/goomba.png")

	// fb, err := NewFramebuffer(200, 200)
	// if err != nil {
	// 	panic(err)
	// }
	// var p bool
	// gl.Enable(gl.SCISSOR_TEST)
	for !window.ShouldClose() {
		glfw.PollEvents()
		gl.ClearColor(1, 1, 1, 1)
		gl.Clear(gl.COLOR_BUFFER_BIT | gl.DEPTH_BUFFER_BIT | gl.STENCIL_BUFFER_BIT)

		uiCtx.NewFrame([2]float32{float32(Width),float32(Height)})

		firstWindow()

		// secondWindow()

		// fb.Bind()
		uiCtx.EndFrame([2]float32{float32(Width), float32(Height)})
		// fb.Unbind()
		// rend.NewFrame()


		window.SwapBuffers()

		endTime = float32(glfw.GetTime())
		dt = endTime - beginTime
		beginTime = endTime

		// fmt.Println(time / float32(steps))

		time += dt
		steps++
	}
}

func setCursor(c ui.CursorType) {
	defC := glfw.CreateStandardCursor(ui.Cursor(c))
	window.SetCursor(defC)
}

var tex *gogl.Texture
var tex2 *gogl.Texture

func getImageFromFilePath(filePath string) (image.Image, error) {
	f, err := os.Open(filePath)
	if err != nil {
		return nil, err
	}
	defer f.Close()
	image, _, err := image.Decode(f)
	return image, err
}

var ish bool = false

func firstWindow() {
	uiCtx.BeginWindow("first wnd")

	if uiCtx.Button("bfgfhf") {
		fmt.Println("button clicked f 1")
		// 	ish = !ish
		uiCtx.SetScrollY(200)
	}

	// if uiCtx.ButtonRR(tex) {
	// 	fmt.Println("button clicked f 2")
	// }
	// uiCtx.ButtonT("Нажать",24)
	// uiCtx.Text("#t3","xello world!", 40)
	uiCtx.Text("#t3", "xello world!", 40)
	uiCtx.VSpace("#vs1fdgdf")
	uiCtx.VSpace("#vs1")
	uiCtx.Image("#im2", tex)
	uiCtx.VSpace("#vs12")

	uiCtx.Image("#im76", tex2)
	uiCtx.Image("#im4", tex)
	// if uiCtx.Image("#im1",tex) {
	// 	fmt.Println("image clicked s 1")
	// 	ish = !ish
	// }
	if ish {
		// if uiCtx.Image("#im2",tex) {
		// 	fmt.Println("image clicked s 1")
		// 	ish = !ish
		// }
	}
	uiCtx.VSpace("#vs1hfg")
	if ish {

		// uiCtx.Text("#t2","Hello world!:", 30)
		// uiCtx.VSpace("#vs13r")
		// uiCtx.Image("#im2",tex)
		// uiCtx.VSpace("#vs13r")
		// uiCtx.Text("#t2erwer","Hello world!:", 30)
	}
	if uiCtx.Image("#im3treyr", tex) {
		fmt.Println("image clicked s 1")
		ish = !ish
	}
	// uiCtx.VSpace("#vs1")
	// uiCtx.Image("#im76",tex2)
	// 	fmt.Println("image clicked s 1")
	// 	ish = !ish
	// }

	uiCtx.EndWindow()
}

func secondWindow() {
	uiCtx.BeginWindow("second wnd")

	uiCtx.Image("#im4", tex)

	// uiCtx.VSpace("#dfdf")
	// // if uiCtx.ButtonRR(tex) {
	// // 	fmt.Println("button clicked f 2")
	// // }
	// if uiCtx.Button("#v3245g") {
	// 	fmt.Println("button clicked f 2")
	// }
	// uiCtx.VSpace("#vs3354")
	// // if uiCtx.ButtonRR(tex) {
	// // 	fmt.Println("button clicked f 2")
	// // }
	// if uiCtx.Button("#354362") {
	// 	fmt.Println("button clicked f 2")
	// }
	// if uiCtx.Button("#vs243646547") {
	// 	fmt.Println("button clicked f 2")
	// }
	// if uiCtx.Button("#vs234634") {
	// 	fmt.Println("button clicked f 2")
	// }

	uiCtx.EndWindow()
}

func cursorPosCallback(w *glfw.Window, xpos float64, ypos float64) {
	uiCtx.Io().MousePosCallback(float32(xpos), float32(ypos))
}

func mouseBtnCallback(w *glfw.Window, button glfw.MouseButton, action glfw.Action, mods glfw.ModifierKey) {
	uiCtx.Io().MouseBtnCallback(ui.GlfwMouseKey(button), ui.GlfwAction(action))
}

func onKey(window *glfw.Window, key glfw.Key, scancode int, action glfw.Action, mods glfw.ModifierKey) {

	switch action {
	case glfw.Press:
		uiCtx.Io().KeyCallback(ui.GlfwKeyToGuiKey(key), true)
	case glfw.Release:
		uiCtx.Io().KeyCallback(ui.GlfwKeyToGuiKey(key), false)
	}
	if key == glfw.KeyEscape && action == glfw.Press {
		window.SetShouldClose(true)
	}
}

func scrollCallback(w *glfw.Window, xoff float64, yoff float64) {
	uiCtx.Io().ScrollX = xoff
	uiCtx.Io().ScrollY = yoff
}
