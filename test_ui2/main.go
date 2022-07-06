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
	"github.com/go-gl/mathgl/mgl32"
)

func init() {
	runtime.LockOSThread()
}

var cam *gogl.Camera
var uiCtx *ui.UiContext
var window *glfw.Window

func main() {
	err := glfw.Init()
	if err != nil {
		panic(err)
	}
	glfw.DefaultWindowHints()
	// glfw.WindowHint(glfw.OpenGLDebugContext, 1)

	window, err = glfw.CreateWindow(1280, 720, "example", nil, nil)
	if err != nil {
		panic(err)
	}
	window.MakeContextCurrent()

	glfw.SwapInterval(1)

	size := func(w *glfw.Window, width int, height int) {
		cam.UpdateProjection(mgl32.Vec2{float32(width), float32(height)})
		gl.Viewport(0, 0, int32(width), int32(height))
		fmt.Println(width, height)
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

	cam = gogl.NewCamera(mgl32.Vec2{0, 0})

	cam.UpdateProjection(mgl32.Vec2{1280, 720})

	uiCtx = ui.UiCtx

	front := render.NewGlRenderer()
	uiCtx.Initialize(front, cam)
	// uiCtx = ui.NewContext(front, cam)
	uiCtx.Io().SetCursor = setCursor

	gl.Enable(gl.BLEND)
	gl.BlendFunc(gl.SRC_ALPHA, gl.ONE_MINUS_SRC_ALPHA)

	beginTime := float32(glfw.GetTime())
	var endTime float32
	var dt float32
	dt = dt

	tex, _ = tex.Init("assets/images/mario.png")

	// var p bool
	for !window.ShouldClose() {
		glfw.PollEvents()
		gl.ClearColor(1, 1, 1, 1)
		gl.Clear(gl.COLOR_BUFFER_BIT)

		uiCtx.NewFrame()

		firstWindow()
		secondWindow()

		

		if uiCtx.Io().IsKeyPressed(ui.GuiKey_Space) {

			// fmt.Println(uiCtx.Io().MousePos)
			// if uiCtx.ActiveWindow != nil {
			// 	fmt.Println("ACTIVE: ", uiCtx.ActiveWindow.Id)
			// } else {
			// 	fmt.Println("ACTIVE: nil")
			// }
			// if uiCtx.HoveredWindow != nil {
			// 	fmt.Println("HOVERED: ", uiCtx.HoveredWindow.Id)
			// } else {
			// 	fmt.Println("HOVERED: nil")
			// }
			// for _, v := range uiCtx.Windows {

			// 	fmt.Println(v.Id)
			// }
			// fmt.Println("==========")
			// fmt.Println(uiCtx.Io().IsDragging)
			// fmt.Println(uiCtx.Io().MouseDelta)

			fmt.Println(uiCtx.ActiveWidget)
			// fmt.Println(uiCtx.Io().ScrollX, uiCtx.Io().ScrollY)
		}

		uiCtx.EndFrame()

		// rend.NewFrame()

		// rend.Rectangle(500, 500, 100, 100, [4]float32{1, 0, 0, 1})
		// rend.Rectangle(300, 400, 50, 300, [4]float32{0, 1, 0, 1})
		// rend.Rectangle(700, 400, 300, 300, [4]float32{0, 0, 0, 1})
		// rend.Circle(100, 200, 140, 30, [4]float32{1, 0, 0, 1})
		// rend.DrawArc(200,200, 100, 30, render.TopLeft, [4]float32{1, 0, 0, 1})
		// rend.Trinagle(500, 100, 600, 250, 700, 100, [4]float32{1, 0, 0, 1})
		// rend.RoundedRectangle(700, 400, 300, 200, 20, [4]float32{1, 0, 0, 1})
		// rend.Line(800, 10, 900, 50, 7, [4]float32{0, 0, 0, 1})
		// rend.Draw(cam)
		// rend.End()

		// batch.AddText("Привет, мир!\n920043 ~hghguij Progress #$@\n[A-Za-z] {0-9_20-33}", 50, 600, 1, colornames.Black)
		// batch.AddText("My name is Dmitry", 100, 340, 1, colornames.Magenta)

		// batch.FlushBatch()

		window.SwapBuffers()
		endTime = float32(glfw.GetTime())
		dt = endTime - beginTime
		beginTime = endTime
	}
}

func setCursor(c ui.CursorType) {
	defC := glfw.CreateStandardCursor(ui.Cursor(c))
	window.SetCursor(defC)
}

var tex *gogl.Texture

func getImageFromFilePath(filePath string) (image.Image, error) {
	f, err := os.Open(filePath)
	if err != nil {
		return nil, err
	}
	defer f.Close()
	image, _, err := image.Decode(f)
	return image, err
}

func firstWindow() {
	uiCtx.BeginWindow()

	if uiCtx.Button() {
		fmt.Println("button clicked f 1")
	}
	uiCtx.VSpace()
	if uiCtx.ButtonRR(tex) {
		fmt.Println("button clicked f 2")
	}

	uiCtx.VSpace()
	if uiCtx.Button() {
		fmt.Println("button clicked f 3")
	}

	uiCtx.EndWindow()
}

func secondWindow() {
	uiCtx.BeginWindow()

	if uiCtx.Image(tex) {
		fmt.Println("image clicked s 1")
	}

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
