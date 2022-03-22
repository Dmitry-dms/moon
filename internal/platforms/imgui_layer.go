package platforms

import (
	"fmt"
	"strings"

	//	"github.com/Dmitry-dms/moon/internal/scenes"
	"github.com/Dmitry-dms/moon/internal/scenes"
	"github.com/go-gl/gl/v4.2-core/gl"
	"github.com/go-gl/glfw/v3.3/glfw"
	imgui "github.com/inkyblackness/imgui-go/v4"
)

type ImgUi struct {
	context *imgui.Context
	io      imgui.IO

	glslVersion          string
	fontTexture            uint32
	shaderProgramId        uint32
	vertexShaderId         uint32
	fragmentShaderId       uint32
	attribLocationTex      int32
	attribLocationProjMtx  int32
	attribLocationPosition int32
	attribLocationUV       int32
	attribLocationColor    int32
	vboHandle              uint32
	elementsHandle         uint32

	clearColor [3]float32
}

func NewImgui() *ImgUi {
	context := imgui.CreateContext(nil)
	io := imgui.CurrentIO()
	err := gl.Init()
	if err != nil {
		panic("failed to initialize OpenGL")
	}
	g := ImgUi{
		context:     context,
		io:          io,
		glslVersion: "#version 420",
		clearColor:  [3]float32{0.0, 0.0, 0.0},
	}
	io.SetBackendFlags(io.GetBackendFlags() | imgui.BackendFlagsRendererHasVtxOffset)
	g.setKeyMapping()
	g.createDeviceObjects()
	return &g
}

func (g *ImgUi) Update(displaySize [2]float32, framebufferSize [2]float32, dt float32, currentScene scenes.Scene) {

	imgui.NewFrame()

	//currentScene.Imgui()

	// if showGoDemoWindow {
	// 	demo.Show(&showGoDemoWindow)
	// }

	// Rendering
	imgui.Render() // This call only creates the draw data list. Actual rendering to framebuffer is done below.

	g.PreRender(g.clearColor)
	g.Render(displaySize, framebufferSize, imgui.RenderedDrawData())
}

// func (g *ImgUi) InitImgui() {
// 	g.setKeyMapping()
// }

var unversionedVertexShader string = `uniform mat4 ProjMtx;

in vec2 Position;
in vec2 UV;
in vec4 Color;

out vec2 Frag_UV;
out vec4 Frag_Color;

void main()
{
    Frag_UV = UV;
    Frag_Color = Color;
    gl_Position = ProjMtx * vec4(Position.xy, 0, 1);
}`

var unversionedFragmentShader string = `uniform sampler2D Texture;

in vec2 Frag_UV;
in vec4 Frag_Color;

out vec4 Out_Color;

void main()
{
    Out_Color = vec4(Frag_Color.rgb, Frag_Color.a * texture(Texture, Frag_UV.st).r);
}`

// PreRender clears the framebuffer.
func (g *ImgUi) PreRender(clearColor [3]float32) {
	gl.ClearColor(clearColor[0], clearColor[1], clearColor[2], 1.0)
	gl.Clear(gl.COLOR_BUFFER_BIT)
}
func (g *ImgUi) createDeviceObjects() {
	//glfw.GetCurrentContext().MakeContextCurrent()
	// Backup GL state
	var lastTexture int32
	var lastArrayBuffer int32
	var lastVertexArray int32
	gl.GetIntegerv(gl.TEXTURE_BINDING_2D, &lastTexture)
	gl.GetIntegerv(gl.ARRAY_BUFFER_BINDING, &lastArrayBuffer)
	gl.GetIntegerv(gl.VERTEX_ARRAY_BINDING, &lastVertexArray)

	vertexShader := g.glslVersion + "\n" + unversionedVertexShader
	fragmentShader := g.glslVersion + "\n" + unversionedFragmentShader

	g.shaderProgramId = gl.CreateProgram()
	g.vertexShaderId = gl.CreateShader(gl.VERTEX_SHADER)
	g.fragmentShaderId = gl.CreateShader(gl.FRAGMENT_SHADER)

	glShaderSource := func(handle uint32, source string) {
		csource, free := gl.Strs(source + "\x00")
		defer free()

		gl.ShaderSource(handle, 1, csource, nil)
	}

	glShaderSource(g.vertexShaderId, vertexShader)
	glShaderSource(g.fragmentShaderId, fragmentShader)

	gl.CompileShader(g.vertexShaderId)
	var s string
	compileError(g.vertexShaderId, &s)
	if s != "" {
		fmt.Println(s)
	}
	gl.CompileShader(g.fragmentShaderId)
	compileError(g.fragmentShaderId, &s)
	if s != "" {
		fmt.Println(s)
	}

	gl.AttachShader(g.shaderProgramId, g.vertexShaderId)
	gl.AttachShader(g.shaderProgramId, g.fragmentShaderId)
	gl.LinkProgram(g.shaderProgramId)

	var status int32
	gl.GetProgramiv(g.shaderProgramId, gl.LINK_STATUS, &status) //logging
	if status == gl.FALSE {
		var logLength int32
		gl.GetProgramiv(g.shaderProgramId, gl.INFO_LOG_LENGTH, &logLength)
		log := strings.Repeat("\x00", int(logLength)+1)
		gl.GetProgramInfoLog(g.shaderProgramId, logLength, nil, gl.Str(log))
		s = log
	}
	if s != "" {
		fmt.Println(s)
	}

	g.attribLocationTex = gl.GetUniformLocation(g.shaderProgramId, gl.Str("Texture"+"\x00"))
	g.attribLocationProjMtx = gl.GetUniformLocation(g.shaderProgramId, gl.Str("ProjMtx"+"\x00"))
	g.attribLocationPosition = gl.GetAttribLocation(g.shaderProgramId, gl.Str("Position"+"\x00"))
	g.attribLocationUV = gl.GetAttribLocation(g.shaderProgramId, gl.Str("UV"+"\x00"))
	g.attribLocationColor = gl.GetAttribLocation(g.shaderProgramId, gl.Str("Color"+"\x00"))

	gl.GenBuffers(1, &g.vboHandle)
	gl.GenBuffers(1, &g.elementsHandle)

	g.createFontsTexture()

	// Restore modified GL state
	gl.BindTexture(gl.TEXTURE_2D, uint32(lastTexture))
	gl.BindBuffer(gl.ARRAY_BUFFER, uint32(lastArrayBuffer))
	gl.BindVertexArray(uint32(lastVertexArray))
}

func compileError(shaderId uint32, info *string) {
	var status int32
	gl.GetShaderiv(shaderId, gl.COMPILE_STATUS, &status) //logging
	if status == gl.FALSE {
		var logLength int32
		gl.GetShaderiv(shaderId, gl.INFO_LOG_LENGTH, &logLength)
		log := strings.Repeat("\x00", int(logLength)+1)
		gl.GetShaderInfoLog(shaderId, logLength, nil, gl.Str(log))
		*info = fmt.Sprintln("failed to compile shader: \n" + log)
	}
}

// Dispose cleans up the resources.
func (g *ImgUi) Dispose() {
	g.invalidateDeviceObjects()
}

func (g *ImgUi) createFontsTexture() {
	// Build texture atlas
	io := imgui.CurrentIO()
	image := io.Fonts().TextureDataAlpha8()

	// Upload texture to graphics system
	var lastTexture int32
	gl.GetIntegerv(gl.TEXTURE_BINDING_2D, &lastTexture)
	gl.GenTextures(1, &g.fontTexture)
	gl.BindTexture(gl.TEXTURE_2D, g.fontTexture)
	gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_MIN_FILTER, gl.LINEAR)
	gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_MAG_FILTER, gl.LINEAR)
	gl.PixelStorei(gl.UNPACK_ROW_LENGTH, 0)
	gl.TexImage2D(gl.TEXTURE_2D, 0, gl.RED, int32(image.Width), int32(image.Height),
		0, gl.RED, gl.UNSIGNED_BYTE, image.Pixels)

	// Store our identifier
	io.Fonts().SetTextureID(imgui.TextureID(g.fontTexture))

	// Restore state
	gl.BindTexture(gl.TEXTURE_2D, uint32(lastTexture))
}

func (g *ImgUi) invalidateDeviceObjects() {
	if g.vboHandle != 0 {
		gl.DeleteBuffers(1, &g.vboHandle)
	}
	g.vboHandle = 0
	if g.elementsHandle != 0 {
		gl.DeleteBuffers(1, &g.elementsHandle)
	}
	g.elementsHandle = 0

	if (g.shaderProgramId != 0) && (g.vertexShaderId != 0) {
		gl.DetachShader(g.shaderProgramId, g.vertexShaderId)
	}
	if g.vertexShaderId != 0 {
		gl.DeleteShader(g.vertexShaderId)
	}
	g.vertexShaderId = 0

	if (g.shaderProgramId != 0) && (g.fragmentShaderId != 0) {
		gl.DetachShader(g.shaderProgramId, g.fragmentShaderId)
	}
	if g.fragmentShaderId != 0 {
		gl.DeleteShader(g.fragmentShaderId)
	}
	g.fragmentShaderId = 0

	if g.shaderProgramId != 0 {
		gl.DeleteProgram(g.shaderProgramId)
	}
	g.shaderProgramId = 0

	if g.fontTexture != 0 {
		gl.DeleteTextures(1, &g.fontTexture)
		imgui.CurrentIO().Fonts().SetTextureID(0)
		g.fontTexture = 0
	}
}

func (renderer *ImgUi) Render(displaySize [2]float32, framebufferSize [2]float32, drawData imgui.DrawData) {
	// Avoid rendering when minimized, scale coordinates for retina displays (screen coordinates != framebuffer coordinates)
	displayWidth, displayHeight := displaySize[0], displaySize[1]
	fbWidth, fbHeight := framebufferSize[0], framebufferSize[1]
	if (fbWidth <= 0) || (fbHeight <= 0) {
		return
	}
	drawData.ScaleClipRects(imgui.Vec2{
		X: fbWidth / displayWidth,
		Y: fbHeight / displayHeight,
	})

	// Backup GL state
	var lastActiveTexture int32
	gl.GetIntegerv(gl.ACTIVE_TEXTURE, &lastActiveTexture)
	gl.ActiveTexture(gl.TEXTURE0)
	var lastProgram int32
	gl.GetIntegerv(gl.CURRENT_PROGRAM, &lastProgram)
	var lastTexture int32
	gl.GetIntegerv(gl.TEXTURE_BINDING_2D, &lastTexture)
	var lastSampler int32
	gl.GetIntegerv(gl.SAMPLER_BINDING, &lastSampler)
	var lastArrayBuffer int32
	gl.GetIntegerv(gl.ARRAY_BUFFER_BINDING, &lastArrayBuffer)
	var lastElementArrayBuffer int32
	gl.GetIntegerv(gl.ELEMENT_ARRAY_BUFFER_BINDING, &lastElementArrayBuffer)
	var lastVertexArray int32
	gl.GetIntegerv(gl.VERTEX_ARRAY_BINDING, &lastVertexArray)
	var lastPolygonMode [2]int32
	gl.GetIntegerv(gl.POLYGON_MODE, &lastPolygonMode[0])
	var lastViewport [4]int32
	gl.GetIntegerv(gl.VIEWPORT, &lastViewport[0])
	var lastScissorBox [4]int32
	gl.GetIntegerv(gl.SCISSOR_BOX, &lastScissorBox[0])
	var lastBlendSrcRgb int32
	gl.GetIntegerv(gl.BLEND_SRC_RGB, &lastBlendSrcRgb)
	var lastBlendDstRgb int32
	gl.GetIntegerv(gl.BLEND_DST_RGB, &lastBlendDstRgb)
	var lastBlendSrcAlpha int32
	gl.GetIntegerv(gl.BLEND_SRC_ALPHA, &lastBlendSrcAlpha)
	var lastBlendDstAlpha int32
	gl.GetIntegerv(gl.BLEND_DST_ALPHA, &lastBlendDstAlpha)
	var lastBlendEquationRgb int32
	gl.GetIntegerv(gl.BLEND_EQUATION_RGB, &lastBlendEquationRgb)
	var lastBlendEquationAlpha int32
	gl.GetIntegerv(gl.BLEND_EQUATION_ALPHA, &lastBlendEquationAlpha)
	lastEnableBlend := gl.IsEnabled(gl.BLEND)
	lastEnableCullFace := gl.IsEnabled(gl.CULL_FACE)
	lastEnableDepthTest := gl.IsEnabled(gl.DEPTH_TEST)
	lastEnableScissorTest := gl.IsEnabled(gl.SCISSOR_TEST)

	// Setup render state: alpha-blending enabled, no face culling, no depth testing, scissor enabled, polygon fill
	gl.Enable(gl.BLEND)
	gl.BlendEquation(gl.FUNC_ADD)
	gl.BlendFunc(gl.SRC_ALPHA, gl.ONE_MINUS_SRC_ALPHA)
	gl.Disable(gl.CULL_FACE)
	gl.Disable(gl.DEPTH_TEST)
	gl.Enable(gl.SCISSOR_TEST)
	gl.PolygonMode(gl.FRONT_AND_BACK, gl.FILL)

	// Setup viewport, orthographic projection matrix
	// Our visible imgui space lies from draw_data->DisplayPos (top left) to draw_data->DisplayPos+data_data->DisplaySize (bottom right).
	// DisplayMin is typically (0,0) for single viewport apps.
	gl.Viewport(0, 0, int32(fbWidth), int32(fbHeight))
	orthoProjection := [4][4]float32{
		{2.0 / displayWidth, 0.0, 0.0, 0.0},
		{0.0, 2.0 / -displayHeight, 0.0, 0.0},
		{0.0, 0.0, -1.0, 0.0},
		{-1.0, 1.0, 0.0, 1.0},
	}
	gl.UseProgram(renderer.shaderProgramId)
	gl.Uniform1i(renderer.attribLocationTex, 0)
	gl.UniformMatrix4fv(renderer.attribLocationProjMtx, 1, false, &orthoProjection[0][0])
	gl.BindSampler(0, 0) // Rely on combined texture/sampler state.

	// Recreate the VAO every time
	// (This is to easily allow multiple GL contexts. VAO are not shared among GL contexts, and
	// we don't track creation/deletion of windows so we don't have an obvious key to use to cache them.)
	var vaoHandle uint32
	gl.GenVertexArrays(1, &vaoHandle)
	gl.BindVertexArray(vaoHandle)
	gl.BindBuffer(gl.ARRAY_BUFFER, renderer.vboHandle)
	gl.EnableVertexAttribArray(uint32(renderer.attribLocationPosition))
	gl.EnableVertexAttribArray(uint32(renderer.attribLocationUV))
	gl.EnableVertexAttribArray(uint32(renderer.attribLocationColor))
	vertexSize, vertexOffsetPos, vertexOffsetUv, vertexOffsetCol := imgui.VertexBufferLayout()
	gl.VertexAttribPointerWithOffset(uint32(renderer.attribLocationPosition), 2, gl.FLOAT, false, int32(vertexSize), uintptr(vertexOffsetPos))
	gl.VertexAttribPointerWithOffset(uint32(renderer.attribLocationUV), 2, gl.FLOAT, false, int32(vertexSize), uintptr(vertexOffsetUv))
	gl.VertexAttribPointerWithOffset(uint32(renderer.attribLocationColor), 4, gl.UNSIGNED_BYTE, true, int32(vertexSize), uintptr(vertexOffsetCol))
	indexSize := imgui.IndexBufferLayout()
	drawType := gl.UNSIGNED_SHORT
	const bytesPerUint32 = 4
	if indexSize == bytesPerUint32 {
		drawType = gl.UNSIGNED_INT
	}

	// Draw
	for _, list := range drawData.CommandLists() {
		vertexBuffer, vertexBufferSize := list.VertexBuffer()
		gl.BindBuffer(gl.ARRAY_BUFFER, renderer.vboHandle)
		gl.BufferData(gl.ARRAY_BUFFER, vertexBufferSize, vertexBuffer, gl.STREAM_DRAW)

		indexBuffer, indexBufferSize := list.IndexBuffer()
		gl.BindBuffer(gl.ELEMENT_ARRAY_BUFFER, renderer.elementsHandle)
		gl.BufferData(gl.ELEMENT_ARRAY_BUFFER, indexBufferSize, indexBuffer, gl.STREAM_DRAW)

		for _, cmd := range list.Commands() {
			if cmd.HasUserCallback() {
				cmd.CallUserCallback(list)
			} else {
				gl.BindTexture(gl.TEXTURE_2D, uint32(cmd.TextureID()))
				clipRect := cmd.ClipRect()
				gl.Scissor(int32(clipRect.X), int32(fbHeight)-int32(clipRect.W), int32(clipRect.Z-clipRect.X), int32(clipRect.W-clipRect.Y))
				gl.DrawElementsBaseVertexWithOffset(gl.TRIANGLES, int32(cmd.ElementCount()), uint32(drawType),
					uintptr(cmd.IndexOffset()*indexSize), int32(cmd.VertexOffset()))
			}
		}
	}
	gl.DeleteVertexArrays(1, &vaoHandle)

	// Restore modified GL state
	gl.UseProgram(uint32(lastProgram))
	gl.BindTexture(gl.TEXTURE_2D, uint32(lastTexture))
	gl.BindSampler(0, uint32(lastSampler))
	gl.ActiveTexture(uint32(lastActiveTexture))
	gl.BindVertexArray(uint32(lastVertexArray))
	gl.BindBuffer(gl.ARRAY_BUFFER, uint32(lastArrayBuffer))
	gl.BindBuffer(gl.ELEMENT_ARRAY_BUFFER, uint32(lastElementArrayBuffer))
	gl.BlendEquationSeparate(uint32(lastBlendEquationRgb), uint32(lastBlendEquationAlpha))
	gl.BlendFuncSeparate(uint32(lastBlendSrcRgb), uint32(lastBlendDstRgb), uint32(lastBlendSrcAlpha), uint32(lastBlendDstAlpha))
	if lastEnableBlend {
		gl.Enable(gl.BLEND)
	} else {
		gl.Disable(gl.BLEND)
	}
	if lastEnableCullFace {
		gl.Enable(gl.CULL_FACE)
	} else {
		gl.Disable(gl.CULL_FACE)
	}
	if lastEnableDepthTest {
		gl.Enable(gl.DEPTH_TEST)
	} else {
		gl.Disable(gl.DEPTH_TEST)
	}
	if lastEnableScissorTest {
		gl.Enable(gl.SCISSOR_TEST)
	} else {
		gl.Disable(gl.SCISSOR_TEST)
	}
	gl.PolygonMode(gl.FRONT_AND_BACK, uint32(lastPolygonMode[0]))
	gl.Viewport(lastViewport[0], lastViewport[1], lastViewport[2], lastViewport[3])
	gl.Scissor(lastScissorBox[0], lastScissorBox[1], lastScissorBox[2], lastScissorBox[3])
}

func (g *ImgUi) CurrentIO() imgui.IO {
	return g.io
}

func (g *ImgUi) setKeyMapping() {
	// Keyboard mapping. ImGui will use those indices to peek into the io.KeysDown[] array.
	g.io.KeyMap(imgui.KeyTab, int(glfw.KeyTab))
	g.io.KeyMap(imgui.KeyLeftArrow, int(glfw.KeyLeft))
	g.io.KeyMap(imgui.KeyRightArrow, int(glfw.KeyRight))
	g.io.KeyMap(imgui.KeyUpArrow, int(glfw.KeyUp))
	g.io.KeyMap(imgui.KeyDownArrow, int(glfw.KeyDown))
	g.io.KeyMap(imgui.KeyPageUp, int(glfw.KeyPageUp))
	g.io.KeyMap(imgui.KeyPageDown, int(glfw.KeyPageDown))
	g.io.KeyMap(imgui.KeyHome, int(glfw.KeyHome))
	g.io.KeyMap(imgui.KeyEnd, int(glfw.KeyEnd))
	g.io.KeyMap(imgui.KeyInsert, int(glfw.KeyInsert))
	g.io.KeyMap(imgui.KeyDelete, int(glfw.KeyDelete))
	g.io.KeyMap(imgui.KeyBackspace, int(glfw.KeyBackspace))
	g.io.KeyMap(imgui.KeySpace, int(glfw.KeySpace))
	g.io.KeyMap(imgui.KeyEnter, int(glfw.KeyEnter))
	g.io.KeyMap(imgui.KeyEscape, int(glfw.KeyEscape))
	g.io.KeyMap(imgui.KeyA, int(glfw.KeyA))
	g.io.KeyMap(imgui.KeyC, int(glfw.KeyC))
	g.io.KeyMap(imgui.KeyV, int(glfw.KeyV))
	g.io.KeyMap(imgui.KeyX, int(glfw.KeyX))
	g.io.KeyMap(imgui.KeyY, int(glfw.KeyY))
	g.io.KeyMap(imgui.KeyZ, int(glfw.KeyZ))
}
