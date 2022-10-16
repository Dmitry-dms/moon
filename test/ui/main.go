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

	glfw.SwapInterval(1)

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

	//ui.UiCtx.UploadFont("C:/Windows/Fonts/times.ttf", 70)
	ui.UiCtx.UploadFont("C:/Windows/Fonts/arial.ttf", 14)
	//ui.UiCtx.UploadFont("assets/fonts/rany.otf", 16)
	//ui.UiCtx.UploadFont("assets/fonts/mono.ttf", 16)

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

		uiCtx.NewFrame([2]float32{float32(Width), float32(Height)})

		firstWindow()

		if uiCtx.Io().IsKeyPressed(ui.GuiKey_Space) {
			//fmt.Println(uiCtx.ActiveWidget)
			//fmt.Println(uiCtx.ActiveWidgetSpaceId)
			fmt.Println(uiCtx.ActiveWidget)
			//fmt.Println(uiCtx.ActiveWindow)
		}
		//secondWindow()

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

	//if uiCtx.ButtonT("Нажать", "Press") {
	//	ish = !ish
	//
	//}
	//if ish {
	//	uiCtx.Text("#er", "Wdff213 ello world!", 14)
	//	uiCtx.Text("#fgfgd", "hello world!", 14)
	//}

	uiCtx.TreeNode("tree1", "Configuration", func() {
		uiCtx.Text("text-ttp-1", "Обычная картинка, которая  ничего не делает", 14)
		uiCtx.Text("#t3j", "hello world!", 14)
	})

	//uiCtx.VSpace("#vs1fdgdf")

	//uiCtx.Row("row 13214", func() {
	//	uiCtx.Image("#im4kjdg464", tex)
	//	uiCtx.Column("col fdfd", func() {
	//		uiCtx.Image("#im76", tex2)
	//		uiCtx.Image("#im4", tex)
	//	})
	//
	//	uiCtx.Column("col fdfdвава", func() {
	//		uiCtx.Button("ASsfdffb")
	//		uiCtx.Button("ASsfdffbbb")
	//		uiCtx.Slider("slider-1", &slCounter, 0, 255)
	//	})
	//
	//	uiCtx.Image("#im4kj", tex)
	//})
	uiCtx.SubWidgetSpace("widhsp-1", ui.Default, func() {
		uiCtx.Image("#im4kjdg464tht", tex2)
		uiCtx.Image("#im76erewr", tex)
		uiCtx.Text("#t3df", "world!", 24)
	})
	//uiCtx.VSpace("#vs1")
	uiCtx.Image("#imgj4", tex2)

	if uiCtx.ActiveWidget == "#imgj4" {
		uiCtx.Tooltip("ttp-1", func() {
			uiCtx.Text("text-ttp-1", "Обычная картинка, которая  ничего не делает", 14)
			uiCtx.Text("text-ttp-2", "Hello World", 16)
			uiCtx.Text("text-ttp-3", "Hello World", 16)
		})
	}

	uiCtx.EndWindow()
}

var slCounter float32 = 0

func secondWindow() {
	uiCtx.BeginWindow("second wnd")
	uiCtx.Row("row 1dfdf14", func() {
		uiCtx.Image("#im4", tex)
		uiCtx.Image("#im4", tex)
	})

	cl := fmt.Sprintf("%.0f", slCounter)
	//uiCtx.Text("text-1dff", "The quick brown fox jumps over the lazy dog", 16)
	uiCtx.Text("text-1dff", "Съешь еще этих мягких", 16)
	uiCtx.Text("text-1dfhjyf", cl, 16)
	uiCtx.Slider("slider-1", &slCounter, 0, 255)

	uiCtx.Row("row 13214", func() {
		uiCtx.Image("#im4kjdg464", tex)
		uiCtx.Column("col fdfd", func() {
			uiCtx.Image("#im76", tex2)
			uiCtx.Image("#im4", tex)
		})

		uiCtx.Column("col fdfdвава", func() {
			uiCtx.Button("ASsfdffb")
			uiCtx.Button("ASsfdffbbb")
		})

		uiCtx.Image("#im4kj", tex)
	})

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
