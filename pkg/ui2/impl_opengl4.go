package ui2

import (
	"unsafe"

	"github.com/Dmitry-dms/moon/pkg/gogl"
	"github.com/go-gl/gl/v4.2-core/gl"
)

type ImGui_ImplOpenGL3_Data struct {
	Version                                                           string
	FontTexture                                                       uint
	ShaderHandle, VboHandle, ElementsHandle                           uint32
	AttribLocationTex, AttribLocationProjMtx                          int32
	AttribLocationVtxPos, AttribLocationVtxUV, AttribLocationVtxColor int32
	VertexBufferSize, IndexBufferSize                                 int
	HasClipOrigin                                                     bool
}

func ImplOpenGL3_GetBackendData() *ImGui_ImplOpenGL3_Data {
	return GetCurrentContext().Io.BackendRendererUserData
}

func ImGui_ImplOpenGL3_Init(version string) bool {
	io := GetCurrentContext().Io
	err := gl.Init()
	if err != nil {
		panic(err)
	}
	data := ImGui_ImplOpenGL3_Data{}
	io.BackendRendererUserData = &data

	io.BackendRendererName = "impl_opengl42"
	data.Version = version

	var current_tex int32
	//для проверки
	gl.GetIntegerv(gl.TEXTURE_BINDING_2D, &current_tex)

	return true
}

func CreateDeviceObjects() bool {
	bd := context.Io.BackendRendererUserData

	//backup gl state
	var last_texture, last_array_buffer, last_vertex_array int32
	gl.GetIntegerv(gl.TEXTURE_BINDING_2D, &last_texture)
	gl.GetIntegerv(gl.ARRAY_BUFFER_BINDING, &last_array_buffer)
	gl.GetIntegerv(gl.VERTEX_ARRAY_BINDING, &last_vertex_array)

	programId, err := gogl.CreateProgram("assets/shaders/ui.glsl")
	if err != nil {
		panic(err)
	}
	bd.ShaderHandle = programId

	bd.AttribLocationTex = gl.GetUniformLocation(programId, gogl.Str("Texture"))
	bd.AttribLocationProjMtx = gl.GetUniformLocation(programId, gogl.Str("ProjMtx"))
	bd.AttribLocationVtxPos = gl.GetUniformLocation(programId, gogl.Str("Position"))
	bd.AttribLocationVtxUV = gl.GetUniformLocation(programId, gogl.Str("UV"))
	bd.AttribLocationVtxColor = gl.GetUniformLocation(programId, gogl.Str("Color"))

	//create buffers
	gl.GenBuffers(1, &bd.VboHandle)
	gl.GenBuffers(1, &bd.ElementsHandle)

	//restore modified GL state
	gl.BindTexture(gl.TEXTURE_2D, uint32(last_texture))
	gl.BindBuffer(gl.ARRAY_BUFFER, uint32(last_array_buffer))
	return true
}

func ImplOpenGL3_NewFrame() {
	data := GetCurrentContext().Io.BackendRendererUserData
	if data == nil {
		panic("Opengl data is nil. New frame")
	}

	if data.ShaderHandle == 0 {
		CreateDeviceObjects()
	}
}

func ImplOpenGL3_SetupRenderState(draw_data *ImDrawData, fb_width, fb_height int, vertex_array_object uint32) {
	data := GetCurrentContext().Io.BackendRendererUserData

	//setup render state
	gl.Enable(gl.BLEND)
	gl.BlendEquation(gl.FUNC_ADD)
	gl.BlendFuncSeparate(gl.SRC_ALPHA, gl.ONE_MINUS_SRC_ALPHA, gl.ONE, gl.ONE_MINUS_SRC_ALPHA)
	gl.Disable(gl.CULL_FACE)
	gl.Disable(gl.DEPTH_TEST)
	gl.Disable(gl.STENCIL_TEST)
	gl.Enable(gl.SCISSOR_TEST)
	gl.Disable(gl.PRIMITIVE_RESTART)

	// Setup viewport, orthographic projection matrix
	// Our visible imgui space lies from draw_data->DisplayPos (top left) to
	// draw_data->DisplayPos+data_data->DisplaySize (bottom right). DisplayPos is (0,0)
	// for single viewport apps.
	gl.Viewport(0, 0, int32(fb_width), int32(fb_height))
	L := draw_data.DisplayPos.X
	R := draw_data.DisplayPos.X + draw_data.DisplaySize.X
	T := draw_data.DisplayPos.Y
	B := draw_data.DisplayPos.Y + draw_data.DisplaySize.Y

	orthoProjection := [4][4]float32{
		{2 / (R - L), 0, 0, 0},
		{0, 2 / (T - B), 0, 0},
		{0, 0, -1, 0},
		{(R + L) / (L - R), (T + B) / (B - T), 0, 1},
	}

	gl.UseProgram(data.ShaderHandle)
	gl.Uniform1i(data.AttribLocationTex, 0)
	gl.UniformMatrix4fv(data.AttribLocationProjMtx, 1, false, &orthoProjection[0][0])

	gl.BindSampler(0, 0)

	// Bind vertex/index buffers and setup attributes for ImDrawVert
	gl.BindBuffer(gl.ARRAY_BUFFER, data.VboHandle)
	gl.BindBuffer(gl.ELEMENT_ARRAY_BUFFER, data.ElementsHandle)

	gl.EnableVertexAttribArray(uint32(data.AttribLocationVtxPos))
	gl.EnableVertexAttribArray(uint32(data.AttribLocationVtxUV))
	gl.EnableVertexAttribArray(uint32(data.AttribLocationVtxColor))

	// gogl.SetVertexAttribPointer(uint32(data.AttribLocationVtxPos), 2, gl.FLOAT, 32, 0)
	// gogl.SetVertexAttribPointer(uint32(data.AttribLocationVtxUV), 2, gl.FLOAT, 32, 2)
	// gogl.SetVertexAttribPointer(uint32(data.AttribLocationVtxColor), 4, gl.UNSIGNED_BYTE, 32, 4)

	gl.VertexAttribPointerWithOffset(uint32(data.AttribLocationVtxPos), 2, gl.FLOAT, false, int32(32), uintptr(gl.PtrOffset(0)))
	gl.VertexAttribPointerWithOffset(uint32(data.AttribLocationVtxUV), 2, gl.FLOAT, false, int32(32), uintptr(gl.PtrOffset(2)))
	gl.VertexAttribPointerWithOffset(uint32(data.AttribLocationVtxColor), 4, gl.UNSIGNED_BYTE, true, int32(32), uintptr(gl.PtrOffset(4)))
}

func RenderDrawData(draw_data *ImDrawData) {
	fb_width := draw_data.DisplaySize.X * draw_data.FramebufferScale.X
	fb_height := draw_data.DisplaySize.Y * draw_data.FramebufferScale.Y

	if fb_width <= 0 || fb_height <= 0 {
		return
	}
	data := GetCurrentContext().Io.BackendRendererUserData

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

	//Setup desired GL state
	var vertex_array_object uint32 = 0
	ImplOpenGL3_SetupRenderState(draw_data, int(fb_width), int(fb_height), vertex_array_object)

	clip_off := draw_data.DisplayPos
	clip_scale := draw_data.FramebufferScale

	for _, list := range draw_data.CmdLists {

		vtx_buffer_size := len(list.VtxBuffer) * int(unsafe.Sizeof(ImDrawVert{}))
		idx_buffer_size := len(list.IdxBuffer) * 16 // ImDrawIdx

		if data.VertexBufferSize < vtx_buffer_size {
			data.VertexBufferSize = vtx_buffer_size
			gl.BufferData(gl.ARRAY_BUFFER, data.VertexBufferSize, nil, gl.STREAM_DRAW)
		}
		if data.IndexBufferSize < idx_buffer_size {
			data.IndexBufferSize = idx_buffer_size
			gl.BufferData(gl.ELEMENT_ARRAY_BUFFER, data.IndexBufferSize, nil, gl.STREAM_DRAW)
		}
		gl.BufferSubData(gl.ARRAY_BUFFER, 0, vtx_buffer_size, unsafe.Pointer(&list.VtxBuffer[0]))
		gl.BufferSubData(gl.ELEMENT_ARRAY_BUFFER, 0, idx_buffer_size, unsafe.Pointer(&list.IdxBuffer[0]))

		for _, draw_cmd := range list.CmdBuffer {

			// Project scissor/clipping rectangles into framebuffer space
			clip_min := ImVec2{(draw_cmd.ClipRect.X - clip_off.X) * clip_scale.X, (draw_cmd.ClipRect.Y - clip_off.Y) * clip_scale.Y}
			clip_max := ImVec2{(draw_cmd.ClipRect.Z - clip_off.X) * clip_scale.X, (draw_cmd.ClipRect.W - clip_off.Y) * clip_scale.Y}

			if clip_max.X <= clip_min.X || clip_max.Y <= clip_min.Y {
				continue
			}
			// Apply scissor/clipping rectangle (Y is inverted in OpenGL)
			gl.Scissor(int32(clip_min.X), int32(fb_height-clip_max.Y), int32(clip_max.X-clip_min.X), int32(clip_max.Y-clip_min.Y))
			// Bind texture, Draw
			gl.BindTexture(gl.TEXTURE_2D, uint32(draw_cmd.TextureId))
			// gl.DrawElements(gl.TRIANGLES, int32(draw_cmd.ElemCount), gl.UNSIGNED_INT, )
			gl.DrawElementsBaseVertexWithOffset(gl.TRIANGLES, int32(draw_cmd.ElemCount), gl.UNSIGNED_INT,
				uintptr(64), int32(4))
		}

	}

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

// func ImGui_ImplOpenGL3_CreateFontsTexture() bool {
// 	io := GetCurrentContext().Io
// 	bd := ImplOpenGL3_GetBackendData()
// }
