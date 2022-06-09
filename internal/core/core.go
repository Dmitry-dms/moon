package core

import (
	"sync"

	"github.com/Dmitry-dms/moon/internal/listeners"
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
		var width int32 = 1280
		var height int32 = 720
		Window = newCore(width, height, platforms.GLFWClientAPIOpenGL42, 0)
		listeners.SetWinHeight(height)
		listeners.SetWinWidth(width)
	})
}

func (c *Core) GetCurrentScene() scenes.Scene {
	return c.currentScene
}

type Core struct {
	width, height *int32
	glfwWindow    *platforms.GLFW
	currentScene  scenes.Scene
	imGui         *platforms.ImgUi
	framebuffer   *renderers.Framebuffer
}

func newCore(width, height int32, glVersion platforms.GLFWClientAPI, scene int) *Core {
	platform, err := platforms.NewGLFW(glVersion, &width, &height)
	if err != nil {
		panic(err)
	}
	renderers.Init()
	c := Core{
		width:      &width,
		height:     &height,
		glfwWindow: platform,
	}
	buffer, err := renderers.NewFramebuffer(200, 200)
	if err != nil {
		panic(err)
	}
	gl.Viewport(0, 0, width, height)
	c.framebuffer = buffer
	c.changeScene(0)
	return &c
}

func (c *Core) changeScene(scene int) {
	switch scene {
	case 0:
		c.currentScene = scenes.NewEditorScene(c.changeScene)
		c.currentScene.Load()
		c.currentScene.Init()
		c.currentScene.Start()
		listeners.SetCamera(c.currentScene.GetCamera())
	case 1:
		// c.currentScene = scenes.NewLevelScene(c.changeScene)
		// c.currentScene.Init()
		// c.currentScene.Start()
	default:
		panic("Unknown scene")
	}
}
func (c *Core) Dispose() {
	c.currentScene.Destroy()
	//c.renderer.Dispose()
	//c.imGuiContext.Destroy()
	c.glfwWindow.Dispose()
}

func (c *Core) Run() {
	beginTime := float32(glfw.GetTime())
	var endTime float32
	var dt float32
	for !c.glfwWindow.ShouldStop() {
		c.glfwWindow.ProcessEvents()

		// renderers.DebugDraw.BeginFrame()

		c.framebuffer.Bind()
		gl.ClearColor(1, 1, 1, 1)
		gl.Clear(gl.COLOR_BUFFER_BIT)

		if dt >= 0 {
			// renderers.DebugDraw.Draw(c.currentScene.GetCamera())
			c.currentScene.Update(dt)
		}
		c.currentScene.Render()
		c.framebuffer.Unbind()

		c.glfwWindow.NewFrame(dt)
		c.glfwWindow.ImguiIO.Update(c.glfwWindow.DisplaySize(), c.glfwWindow.FramebufferSize(), dt, c.currentScene, c.framebuffer.GetTextureId())

		c.glfwWindow.PostRender()

		endTime = float32(glfw.GetTime())
		dt = endTime - beginTime
		beginTime = endTime
	}
}


