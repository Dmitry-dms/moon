package core

import (
	"github.com/Dmitry-dms/moon/internal/platforms"
	"github.com/Dmitry-dms/moon/internal/renderers"
	"github.com/Dmitry-dms/moon/internal/scenes"

	"github.com/go-gl/glfw/v3.3/glfw"
	//imgui "github.com/inkyblackness/imgui-go/v4"
	"github.com/pkg/errors"
)

// type clipboard struct {
// 	platform Platform
// }

// func (board clipboard) Text() (string, error) {
// 	return board.platform.ClipboardText()
// }

// func (board clipboard) SetText(text string) {
// 	board.platform.SetClipboardText(text)
// }

// type Scene interface {
// 	Init()
// 	Start()
// 	Destroy()
// 	Update(dt float32)
// 	Render()
// }

type Core struct {
	width, height int
	glfwWindow    *glfw.Window
	currentScene  scenes.Scene
	platform      platforms.Platform
	renderer      renderers.Renderer
	imGui         *platforms.ImgUi
}

func NewCore(width, height int, glVersion platforms.GLFWClientAPI, scene int) (*Core, error) {

	platform, err := platforms.NewGLFW(glVersion)
	if err != nil {
		return nil, errors.Wrap(err, "Can't initialize GLFW")
	}

	renderer, err := renderers.NewOpenGL42()
	if err != nil {
		return nil, errors.Wrapf(err, "Can't initialize OpenGL %s", glVersion)
	}

	c := Core{
		width:      width,
		height:     height,
		glfwWindow: platform.GetWindow(),
		platform:   platform,
		renderer:   renderer,
		imGui:      platform.ImguiIO,
	}

	switch scene {
	case 0:
		c.currentScene = scenes.NewEditorScene(renderer)
		c.currentScene.Init()
	}

	return &c, nil
}
func (c *Core) Dispose() {
	c.renderer.Dispose()
	//c.imGuiContext.Destroy()
	c.platform.Dispose()
}

func (c *Core) Run() {
	// beginTime := glfw.GetTime()
	// var endTime, dt float64
	//showGoDemoWindow := false

	var dt float32
	for !c.platform.ShouldStop() {
		c.platform.ProcessEvents()

		// Signal start of a new frame
		c.platform.NewFrame(&dt)

		if dt >= 0 {
			c.currentScene.Update(dt)
		}
		c.currentScene.Render()
		c.imGui.Update(c.platform.DisplaySize(), c.platform.FramebufferSize(), dt, c.currentScene)

		// A this point, the application could perform its own rendering...
		// app.RenderScene()
		c.platform.PostRender()

		// endTime = glfw.GetTime()
		// dt = endTime - beginTime
		// beginTime = endTime
	}
}
