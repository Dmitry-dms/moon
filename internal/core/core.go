package core

import (
	"fmt"
	"time"

	"github.com/Dmitry-dms/moon/internal/platforms"
	"github.com/Dmitry-dms/moon/internal/renderers"
	"github.com/go-gl/glfw/v3.3/glfw"
	imgui "github.com/inkyblackness/imgui-go/v4"
	"github.com/pkg/errors"
)

type Scene interface {
	Init()
	Start()
	Destroy()
	Update(dt float32)
	Render()
}

type Core struct {
	width, height int
	glfwWindow    *glfw.Window
	currentScene  Scene
	platform      Platform
	renderer      Renderer
	imGuiContext *imgui.Context
}

func NewCore(width, height int, glVersion platforms.GLFWClientAPI) (*Core,error) {
	context := imgui.CreateContext(nil)

	io := imgui.CurrentIO()

	platform, err := platforms.NewGLFW(io, glVersion)
	if err != nil {
		return nil, errors.Wrap(err, "Can't initialize GLFW")
	}


	renderer, err := renderers.NewOpenGL42(io)
	if err != nil {
		return nil, errors.Wrapf(err, "Can't initialize OpenGL %s", glVersion)
	}

	c := Core{
		width: width,
		height: height,
		glfwWindow: platform.GetWindow(),
		platform: platform,
		renderer: renderer,
		imGuiContext: context,
	}
	return &c, nil
}
func (c *Core) Dispose() {
	c.imGuiContext.Destroy()
	c.renderer.Dispose()
	c.platform.Dispose()
}

func (c *Core) Run() {
	showDemoWindow := true
	//showGoDemoWindow := false
	clearColor := [3]float32{0.0, 0.0, 0.0}

	for !c.platform.ShouldStop() {
		c.platform.ProcessEvents()

		// Signal start of a new frame
		c.platform.NewFrame()
		imgui.NewFrame()

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
		if showDemoWindow {
			// Normally user code doesn't need/want to call this because positions are saved in .ini file anyway.
			// Here we just want to make the demo initial state a bit more friendly!
			const demoX = 650
			const demoY = 20
			imgui.SetNextWindowPosV(imgui.Vec2{X: demoX, Y: demoY}, imgui.ConditionFirstUseEver, imgui.Vec2{})
			imgui.ShowDemoWindow(&showDemoWindow)
		}

		// if showGoDemoWindow {
		// 	demo.Show(&showGoDemoWindow)
		// }

		// Rendering
		imgui.Render() // This call only creates the draw data list. Actual rendering to framebuffer is done below.

		c.renderer.PreRender(clearColor)
		// A this point, the application could perform its own rendering...
		// app.RenderScene()

		c.renderer.Render(c.platform.DisplaySize(),c.platform.FramebufferSize(), imgui.RenderedDrawData())
		c.platform.PostRender()
	}
}
