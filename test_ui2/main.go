package main

import (
	"runtime"

	"github.com/Dmitry-dms/moon/pkg/gogl"
	"github.com/Dmitry-dms/moon/pkg/ui"
	"github.com/Dmitry-dms/moon/pkg/ui2"

	// "github.com/Dmitry-dms/moon/pkg/ui2/fonts"
	// "golang.org/x/image/colornames"

	"github.com/go-gl/gl/v4.2-core/gl"
	"github.com/go-gl/glfw/v3.3/glfw"
	"github.com/go-gl/mathgl/mgl32"
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
	ui2.CreateContext()
	ui2.ImGui_ImplGlfw_InitForOpenGL(window, true)
	ui2.ImGui_ImplOpenGL3_Init("42")

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
	rend := ui.NewUIRenderer(100, 0)
	rend.Start()

	gl.Enable(gl.BLEND)
	gl.BlendFunc(gl.SRC_ALPHA, gl.ONE_MINUS_SRC_ALPHA)
	// cam := gogl.NewCamera(mgl32.Vec2{0, 0})
	cam := &gogl.Camera{
		Position: mgl32.Vec2{0, 0},
		Proj: mgl32.Ortho(0, 1280, 0, 720, 0, 100),
		ProjSize: mgl32.Vec2{1280, 720},
	}
	cam.AdjustProjection()


	spr := ui.DefSpriteRenderer()
	tex := &gogl.Texture{}
	tex, err = tex.Init("assets/images/blend1.png")
	if err != nil {
		panic(err)
	}
	sprite := gogl.DefSprite()
	sprite.SetTexture(tex)

	spr.SetSprite(sprite)
	btn := ui.Button{
		UiObject: ui.UiObject{
			Transform: ui.NewTransform(mgl32.Vec2{0, 0}, mgl32.Vec2{200, 200}),
			Name:      "btn1",
			ZIndex:    5,
			Spr:       spr,
		},
	}
	rend.AddUIComponent(&btn)
	

	beginTime := float32(glfw.GetTime())
	var endTime float32
	var dt float32

	for !window.ShouldClose() {
		glfw.PollEvents()
		gl.ClearColor(1, 1, 1, 1)
		gl.Clear(gl.COLOR_BUFFER_BIT)
		rend.Update(dt)

		// ui2.ImplOpenGL3_NewFrame()
		// ui2.ImplGlfw_NewFrame(window)

		// batch.AddText("Привет, мир!\n920043 ~hghguij Progress #$@\n[A-Za-z] {0-9_20-33}", 50, 600, 1, colornames.Black)
		// batch.AddText("My name is Dmitry", 100, 340, 1, colornames.Magenta)

		// batch.FlushBatch()
		rend.Render(cam)
		window.SwapBuffers()
		endTime = float32(glfw.GetTime())
		dt = endTime - beginTime
		beginTime = endTime
	}
}
