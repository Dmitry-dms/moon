package core

import (
	"sync"

	//"github.com/Dmitry-dms/moon/internal/listeners"
	"github.com/Dmitry-dms/moon/internal/platforms"
	"github.com/Dmitry-dms/moon/internal/renderers"
	"github.com/Dmitry-dms/moon/internal/scenes"

	"github.com/go-gl/gl/v4.2-core/gl"
	"github.com/go-gl/glfw/v3.3/glfw"
	//"github.com/go-gl/glfw/v3.3/glfw"
	//imgui "github.com/inkyblackness/imgui-go/v4"
	//"github.com/pkg/errors"
)

var Window *Core

func init() {
	o := sync.Once{}
	o.Do(func() { //make a singleton
		Window = NewCore(1200, 700, GLFWClientAPIOpenGL42, 0)
	})
}

func (c *Core) GetCurrentScene() scenes.Scene {
	return c.currentScene
}

type Core struct {
	width, height int
	glfwWindow    *GLFW
	currentScene  scenes.Scene
	renderer      renderers.Renderer
	imGui         *platforms.ImgUi
}

func NewCore(width, height int, glVersion GLFWClientAPI, scene int) *Core {

	platform, err := NewGLFW(glVersion, width, height)
	if err != nil {
		panic(err)
		//	return nil, errors.Wrap(err, "Can't initialize GLFW")
	}

	renderer, err := renderers.NewOpenGL42()
	if err != nil {
		panic(err)
		//return nil, errors.Wrapf(err, "Can't initialize OpenGL %s", glVersion)
	}

	c := Core{
		width:      width,
		height:     height,
		glfwWindow: platform,
		renderer:   renderer,
		//imGui:      platform.ImguiIO,
	}

	c.changeScene(0)

	return &c //, nil
}

func (c *Core) changeScene(scene int) {
	switch scene {
	case 0:
		c.currentScene = scenes.NewEditorScene(c.changeScene)
		c.currentScene.Init()
		c.currentScene.Start()
	case 1:
		// c.currentScene = scenes.NewLevelScene(c.changeScene)
		// c.currentScene.Init()
		// c.currentScene.Start()
	default:
		panic("Unknown scene")
	}
}
func (c *Core) Dispose() {
	c.renderer.Dispose()
	//c.imGuiContext.Destroy()
	c.glfwWindow.Dispose()
}

func (c *Core) Run() {
	beginTime := float32(glfw.GetTime())
	var endTime float32
	var dt float32
	for !c.glfwWindow.ShouldStop() {
		c.glfwWindow.ProcessEvents()

		gl.ClearColor(1, 1, 1, 1)
		gl.Clear(gl.COLOR_BUFFER_BIT)

		if dt >= 0 {
			c.currentScene.Update(dt)
		}
		// Signal start of a new frame
		// c.glfwWindow.NewFrame(&dt)

		// c.imGui.Update(c.glfwWindow.DisplaySize(), c.glfwWindow.FramebufferSize(), dt, c.currentScene)

		// // A this point, the application could perform its own rendering...
		// c.currentScene.Render()
		// // app.RenderScene()
		c.glfwWindow.PostRender()

		endTime = float32(glfw.GetTime())
		dt = endTime - beginTime
		beginTime = endTime
	}
}
