package platforms

import (
	"fmt"
	"math"
	//"math"

	//"github.com/go-gl/gl/v4.2-core/gl"
	"github.com/Dmitry-dms/moon/internal/listeners"
	"github.com/Dmitry-dms/moon/internal/renderers"
	"github.com/go-gl/glfw/v3.3/glfw"
	imgui "github.com/inkyblackness/imgui-go/v4"
	"github.com/pkg/errors"
)

// GLFWClientAPI identifies the render system that shall be initialized.
type GLFWClientAPI string

// This is a list of GLFWClientAPI constants.
const (
	GLFWClientAPIOpenGL3  GLFWClientAPI = "OpenGL3"
	GLFWClientAPIOpenGL42 GLFWClientAPI = "OpenGL42"
)

// GLFW implements a platform based on github.com/go-gl/glfw (v3.2).
type GLFW struct {
	width, height *int32
	ImguiIO       *ImgUi
	renderer      *renderers.OpenGL42
	window *glfw.Window

	time             float64
	mouseJustPressed [3]bool
}

func (g *GLFW) GetWindow() *glfw.Window {
	return g.window
}

// NewGLFW attempts to initialize a GLFW context.
func NewGLFW(clientAPI GLFWClientAPI, width, height *int32) (*GLFW, error) {

	err := glfw.Init()
	if err != nil {
		return nil, fmt.Errorf("failed to initialize glfw: %w", err)
	}
	//configure glfw
	switch clientAPI {
	case GLFWClientAPIOpenGL3:
		glfw.WindowHint(glfw.ContextVersionMajor, 3)
		glfw.WindowHint(glfw.ContextVersionMinor, 2)
		glfw.WindowHint(glfw.OpenGLProfile, glfw.OpenGLCoreProfile)
		glfw.WindowHint(glfw.OpenGLForwardCompatible, 1)
	case GLFWClientAPIOpenGL42:
		glfw.WindowHint(glfw.ContextVersionMajor, 4)
		glfw.WindowHint(glfw.ContextVersionMinor, 2)
		glfw.WindowHint(glfw.OpenGLProfile, glfw.OpenGLCoreProfile)
	default:
		glfw.Terminate()
		return nil, errors.New("Unsupported GL version")
	}
	glfw.DefaultWindowHints()
	glfw.WindowHint(glfw.Visible, glfw.False) //пока окно не скофигугрировано, не показываем его
	glfw.WindowHint(glfw.Resizable, glfw.True)
	//glfw.WindowHint(glfw.Maximized, glfw.True) // открывается сразу в windowed fullscreen


	//Создание окна
	window, err := glfw.CreateWindow(int(*width), int(*height), "GLFW+"+string(clientAPI)+" example", nil, nil)
	if err != nil {
		glfw.Terminate()
		return nil, fmt.Errorf("failed to create window: %w", err)
	}
	//переключение на контекст OpenGL
	window.MakeContextCurrent()
	//включаем верт. синхронизацию
	glfw.SwapInterval(1)
	
	//включаем окно
	window.Show()
	platform := &GLFW{
		window: window,
		width:  width,
		height: height,
	}
	platform.installCallbacks()

	renderer, err := renderers.NewOpenGL42()
	if err != nil {
		panic(err)
	}
	platform.renderer = renderer

	imGui := NewImgui()
	platform.ImguiIO = imGui

	return platform, nil
}

// Dispose cleans up the resources.
func (platform *GLFW) Dispose() {
	// platform.ImguiIO.Dispose()
	platform.window.Destroy()
	glfw.Terminate()
	
}

// ShouldStop returns true if the window is to be closed.
func (platform *GLFW) ShouldStop() bool {
	return platform.window.ShouldClose()
}

// ProcessEvents handles all pending window events.
func (platform *GLFW) ProcessEvents() {
	glfw.PollEvents()
}

// DisplaySize returns the dimension of the display.
func (platform *GLFW) DisplaySize() [2]float32 {
	w, h := platform.window.GetSize()
	return [2]float32{float32(w), float32(h)}
}

// FramebufferSize returns the dimension of the framebuffer.
func (platform *GLFW) FramebufferSize() [2]float32 {
	w, h := platform.window.GetFramebufferSize()
	return [2]float32{float32(w), float32(h)}
}

//NewFrame marks the begin of a render pass. It forwards all current state to imgui IO.
func (platform *GLFW) NewFrame(dt float32) {
	io := platform.ImguiIO.CurrentIO()
	// Setup display size (every frame to accommodate for window resizing)
	displaySize := platform.DisplaySize()
	io.SetDisplaySize(imgui.Vec2{X: displaySize[0], Y: displaySize[1]})
	io.SetDeltaTime(dt)


	// Setup inputs
	if platform.window.GetAttrib(glfw.Focused) != 0 {
		x, y := platform.window.GetCursorPos()
		io.SetMousePosition(imgui.Vec2{X: float32(x), Y: float32(y)})
	} else {
		io.SetMousePosition(imgui.Vec2{X: -math.MaxFloat32, Y: -math.MaxFloat32})
	}

	for i := 0; i < len(platform.mouseJustPressed); i++ {
		down := platform.mouseJustPressed[i] || (platform.window.GetMouseButton(glfwButtonIDByIndex[i]) == glfw.Press)
		io.SetMouseButtonDown(i, down)
		platform.mouseJustPressed[i] = false
	}
}

// PostRender performs a buffer swap.
func (platform *GLFW) PostRender() {
	platform.window.SwapBuffers()
}

func (glfw *GLFW) installCallbacks() {
	// platform.window.SetMouseButtonCallback(platform.mouseButtonChange)
	// platform.window.SetScrollCallback(platform.mouseScrollChange)
	// platform.window.SetKeyCallback(platform.keyChange)
	glfw.window.SetCharCallback(glfw.charChange)

	//==================
	glfw.window.SetCursorPosCallback(listeners.MousePositionCallback)
	glfw.window.SetMouseButtonCallback(glfw.mouseButtonChange)
	glfw.window.SetScrollCallback(glfw.mouseScrollChange)
	glfw.window.SetKeyCallback(glfw.keyChange)

	glfw.window.SetSizeCallback(glfw.sizeCllback)
}
const (

	mouseButtonPrimary   = 0
	mouseButtonSecondary = 1
	mouseButtonTertiary  = 2
	mouseButtonCount     = 3
)
var glfwButtonIndexByID = map[glfw.MouseButton]int{
	glfw.MouseButton1: mouseButtonPrimary,
	glfw.MouseButton2: mouseButtonSecondary,
	glfw.MouseButton3: mouseButtonTertiary,
}

var glfwButtonIDByIndex = map[int]glfw.MouseButton{
	mouseButtonPrimary:   glfw.MouseButton1,
	mouseButtonSecondary: glfw.MouseButton2,
	mouseButtonTertiary:  glfw.MouseButton3,
}

func (platform *GLFW) mouseButtonChange(window *glfw.Window, rawButton glfw.MouseButton, action glfw.Action, mods glfw.ModifierKey) {
	buttonIndex, known := glfwButtonIndexByID[rawButton]

	if known && (action == glfw.Press) {
		platform.mouseJustPressed[buttonIndex] = true
	}
	// if !platform.ImguiIO.CurrentIO().WantCaptureMouse() {
		listeners.MouseButtonCallback(window, rawButton, action, mods)
	// }
}

func (platform *GLFW) sizeCllback(w *glfw.Window, width int, height int) {
	w1 := int32(width)
	h := int32(height)
	platform.width = &w1
	platform.height = &h

	listeners.SizeCllback(w, width, height)
}

func (platform *GLFW) mouseScrollChange(window *glfw.Window, x, y float64) {
	// platform.ImguiIO.CurrentIO().AddMouseWheelDelta(float32(x), float32(y))

	// if !platform.ImguiIO.CurrentIO().WantCaptureMouse() {
	// 	listeners.MouseScrollCallback(window, x,  y)
	// }
}

func (platform *GLFW) keyChange(window *glfw.Window, key glfw.Key, scancode int, action glfw.Action, mods glfw.ModifierKey) {
	// if action == glfw.Press {
		// platform.ImguiIO.CurrentIO().KeyPress(int(key))
	// }
	// if action == glfw.Release {
		// platform.ImguiIO.CurrentIO().KeyRelease(int(key))
	// }

	// Modifiers are not reliable across systems
	// platform.ImguiIO.CurrentIO().KeyCtrl(int(glfw.KeyLeftControl), int(glfw.KeyRightControl))
	// platform.ImguiIO.CurrentIO().KeyShift(int(glfw.KeyLeftShift), int(glfw.KeyRightShift))
	// platform.ImguiIO.CurrentIO().KeyAlt(int(glfw.KeyLeftAlt), int(glfw.KeyRightAlt))
	// platform.ImguiIO.CurrentIO().KeySuper(int(glfw.KeyLeftSuper), int(glfw.KeyRightSuper))

	// if !platform.ImguiIO.CurrentIO().WantCaptureKeyboard() {
		listeners.KeyCallback(window, key, scancode, action, mods)
	// }
}

func (platform *GLFW) charChange(window *glfw.Window, char rune) {
	// platform.ImguiIO.CurrentIO().AddInputCharacters(string(char))
}

// ClipboardText returns the current clipboard text, if available.
func (platform *GLFW) ClipboardText() (string, error) {
	return platform.window.GetClipboardString(), nil
}

// SetClipboardText sets the text as the current clipboard text.
func (platform *GLFW) SetClipboardText(text string) {
	platform.window.SetClipboardString(text)
}
