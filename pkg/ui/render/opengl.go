////go:build opengl

package render

import (
	"github.com/Dmitry-dms/moon/pkg/gogl"
	"github.com/Dmitry-dms/moon/pkg/ui/draw"
	"github.com/go-gl/gl/v4.2-core/gl"
)

type GLRender struct {
	vaoId, vboId, ebo uint32
	shader            *gogl.Shader
}

const (
	// pos     color       texCoords    texId
	// f,f     f,f,f,f     f,f          f
	posSize       = 2
	colorSize     = 4
	texCoordsSize = 2
	texIdSize     = 1

	vertexSize = posSize + colorSize + texCoordsSize + texIdSize

	posOffset       = 0
	colorOffset     = posOffset + posSize
	texCoordsOffset = colorOffset + colorSize
	texIdOffset     = texCoordsOffset + texCoordsSize
)

func NewGlRenderer() *GLRender {
	s, err := gogl.NewShader("assets/shaders/gui.glsl")
	if err != nil {
		panic(err)
	}
	r := GLRender{

		vaoId:  0,
		vboId:  0,
		ebo:    0,
		shader: s,
	}
	// r.vaoId = gogl.GenBindVAO()

	gl.GenBuffers(1, &r.vboId)
	gl.GenBuffers(1, &r.ebo)

	//включаем layout
	// gogl.SetVertexAttribPointer(0, posSize, gl.FLOAT, vertexSize*4, posOffset)
	// gogl.SetVertexAttribPointer(1, colorSize, gl.FLOAT, vertexSize*4, colorOffset)
	// gogl.SetVertexAttribPointer(2, texCoordsSize, gl.FLOAT, vertexSize*4, texCoordsOffset)
	// gogl.SetVertexAttribPointer(3, texIdSize, gl.FLOAT, vertexSize*4, texIdOffset)

	return &r
}

func (r *GLRender) NewFrame() {

}

// func (r *GLRender) Trinagle(x0, y0, x1, y1, x2, y2 float32, clr [4]float32) {
// 	vert := make([]float32, 9*3)
// 	ind := make([]int32, 3)
// 	offset := 0

// 	fillVertices(vert, &offset, x0, y0, 0, 0, 0, clr)
// 	fillVertices(vert, &offset, x1, y1, 0, 0, 0, clr)
// 	fillVertices(vert, &offset, x2, y2, 0, 0, 0, clr)

// 	ind0 := r.lastIndc
// 	ind1 := ind0 + 1
// 	ind2 := ind1 + 1

// 	ind[0] = int32(ind0)
// 	ind[1] = int32(ind1)
// 	ind[2] = int32(ind2)

// 	r.lastIndc = ind2 + 1
// 	r.render(vert, ind, 3)
// }

// func (r *GLRender) Circle(x, y, radius float32, steps int, clr [4]float32) {
// 	ind0 := r.lastIndc
// 	ind1 := ind0 + 1
// 	ind2 := ind1 + 1
// 	offset := 0
// 	indOffset := 0

// 	angle := math.Pi * 2 / float32(steps)

// 	numV := int(math.Floor(6.28 / float64(angle)))

// 	ind := make([]int32, 3*(numV+1))    // 3 - triangle
// 	vert := make([]float32, 9*(3+numV)) //polygon

// 	var prevX, prevY float32
// 	var lastX, lastY float32

// 	var ang float32 = angle

// 	prevX = x
// 	prevY = y + radius

// 	fillVertices(vert, &offset, x, y, 0, 0, 0, clr)
// 	fillVertices(vert, &offset, prevX, prevY, 0, 0, 0, clr)
// 	newx := x + radius*float32(math.Sin(float64(ang)))
// 	newY := y + radius*float32(math.Cos(float64(ang)))
// 	fillVertices(vert, &offset, newx, newY, 0, 0, 0, clr)
// 	ind[indOffset] = int32(ind0)
// 	ind[indOffset+1] = int32(ind1)
// 	ind[indOffset+2] = int32(ind2)
// 	indOffset += 3
// 	ind1++
// 	ind2++
// 	ang += angle

// 	for ang <= 6.28 { // 360 deg ~= 6.28 rad
// 		newx := x + radius*float32(math.Sin(float64(ang)))
// 		newY := y + radius*float32(math.Cos(float64(ang)))
// 		fillVertices(vert, &offset, newx, newY, 0, 0, 0, clr)

// 		ind[indOffset] = int32(ind0)
// 		ind[indOffset+1] = int32(ind1)
// 		ind[indOffset+2] = int32(ind2)
// 		indOffset += 3
// 		ind1++
// 		ind2++

// 		ang += angle
// 	}
// 	lastX = x
// 	lastY = y + radius
// 	fillVertices(vert, &offset, lastX, lastY, 0, 0, 0, clr)

// 	ind[indOffset] = int32(ind0)
// 	ind[indOffset+1] = int32(ind1)
// 	ind[indOffset+2] = int32(ind2)
// 	// indOffset += 3

// 	r.lastIndc = ind2 + 1
// 	r.render(vert, ind, 3*(numV+1))
// }

// func (r *GLRender) Line(x0, y0, x1, y1 float32, thick int, clr [4]float32) {

// 	r.Trinagle(x0, y0, x1, y1, x1+float32(thick), y1, clr)
// 	r.Trinagle(x1+float32(thick), y1, x0+float32(thick), y0, x0, y0, clr)
// }

var steps = 30

func (b *GLRender) Draw(displaySize [2]float32, buffer draw.CmdBuffer) {

	displayWidth := displaySize[0]
	displayHeight := displaySize[1]

	// Backup GL state
	// var lastActiveTexture int32
	// gl.GetIntegerv(gl.ACTIVE_TEXTURE, &lastActiveTexture)
	// gl.ActiveTexture(gl.TEXTURE0)
	// var lastProgram int32
	// gl.GetIntegerv(gl.CURRENT_PROGRAM, &lastProgram)
	// var lastTexture int32
	// gl.GetIntegerv(gl.TEXTURE_BINDING_2D, &lastTexture)
	// var lastSampler int32
	// gl.GetIntegerv(gl.SAMPLER_BINDING, &lastSampler)
	// var lastArrayBuffer int32
	// gl.GetIntegerv(gl.ARRAY_BUFFER_BINDING, &lastArrayBuffer)
	// var lastElementArrayBuffer int32
	// gl.GetIntegerv(gl.ELEMENT_ARRAY_BUFFER_BINDING, &lastElementArrayBuffer)
	// var lastVertexArray int32
	// gl.GetIntegerv(gl.VERTEX_ARRAY_BINDING, &lastVertexArray)
	// var lastPolygonMode [2]int32
	// gl.GetIntegerv(gl.POLYGON_MODE, &lastPolygonMode[0])
	// var lastViewport [4]int32
	// gl.GetIntegerv(gl.VIEWPORT, &lastViewport[0])
	// var lastScissorBox [4]int32
	// gl.GetIntegerv(gl.SCISSOR_BOX, &lastScissorBox[0])
	// var lastBlendSrcRgb int32
	// gl.GetIntegerv(gl.BLEND_SRC_RGB, &lastBlendSrcRgb)
	// var lastBlendDstRgb int32
	// gl.GetIntegerv(gl.BLEND_DST_RGB, &lastBlendDstRgb)
	// var lastBlendSrcAlpha int32
	// gl.GetIntegerv(gl.BLEND_SRC_ALPHA, &lastBlendSrcAlpha)
	// var lastBlendDstAlpha int32
	// gl.GetIntegerv(gl.BLEND_DST_ALPHA, &lastBlendDstAlpha)
	// var lastBlendEquationRgb int32
	// gl.GetIntegerv(gl.BLEND_EQUATION_RGB, &lastBlendEquationRgb)
	// var lastBlendEquationAlpha int32
	// gl.GetIntegerv(gl.BLEND_EQUATION_ALPHA, &lastBlendEquationAlpha)
	// lastEnableBlend := gl.IsEnabled(gl.BLEND)
	// lastEnableCullFace := gl.IsEnabled(gl.CULL_FACE)
	// lastEnableDepthTest := gl.IsEnabled(gl.DEPTH_TEST)
	// lastEnableScissorTest := gl.IsEnabled(gl.SCISSOR_TEST)
	b.shader.Use()
	vaoId := gogl.GenBindVAO()
	gl.BindBuffer(gl.ARRAY_BUFFER, b.vboId)
	//включаем layout
	gogl.SetVertexAttribPointer(0, posSize, gl.FLOAT, vertexSize*4, posOffset)
	gogl.SetVertexAttribPointer(1, colorSize, gl.FLOAT, vertexSize*4, colorOffset)
	gogl.SetVertexAttribPointer(2, texCoordsSize, gl.FLOAT, vertexSize*4, texCoordsOffset)
	gogl.SetVertexAttribPointer(3, texIdSize, gl.FLOAT, vertexSize*4, texIdOffset)

	//аллоцируем место для vertices
	// vboId := gogl.GenBindBuffer(gl.ARRAY_BUFFER)
	// gogl.BufferData(gl.ARRAY_BUFFER, r.vertices, gl.DYNAMIC_DRAW)
	// gl.BufferData(gl.ARRAY_BUFFER, 4*1500, nil, gl.DYNAMIC_DRAW)

	// ebo := gogl.GenBindBuffer(gl.ELEMENT_ARRAY_BUFFER)
	// gl.BindBuffer(gl.ELEMENT_ARRAY_BUFFER, r.ebo)
	// gl.BufferData(gl.ELEMENT_ARRAY_BUFFER, 500*4, nil, gl.DYNAMIC_DRAW)
	// gl.BindBuffer(gl.ELEMENT_ARRAY_BUFFER, 0)

	// gl.Enable(gl.BLEND)
	// gl.BlendEquation(gl.FUNC_ADD)
	// gl.BlendFunc(gl.SRC_ALPHA, gl.ONE_MINUS_SRC_ALPHA)
	// gl.Disable(gl.CULL_FACE)
	// gl.Disable(gl.DEPTH_TEST)
	gl.Enable(gl.SCISSOR_TEST)
	// gl.PolygonMode(gl.FRONT_AND_BACK, gl.FILL)

	gl.BindBuffer(gl.ARRAY_BUFFER, b.vboId)
	gl.BufferData(gl.ARRAY_BUFFER, len(buffer.Vertices)*4, gl.Ptr(buffer.Vertices), gl.STREAM_DRAW)

	gl.BindBuffer(gl.ELEMENT_ARRAY_BUFFER, b.ebo)
	gl.BufferData(gl.ELEMENT_ARRAY_BUFFER, len(buffer.Indeces)*4, gl.Ptr(buffer.Indeces), gl.STREAM_DRAW)

	orthoProjection := [4][4]float32{
		{2.0 / displayWidth, 0.0, 0.0, 0.0},
		{0.0, 2.0 / displayHeight, 0.0, 0.0},
		{0.0, 0.0, -2.0, 0.0},
		{-1.0, -1.0, -1.0, 1.0},
	}

	b.shader.UploadMatslice("uProjection", orthoProjection)

	for _, cmd := range buffer.Inf {
		//mainRect := cmd.Clip.MainClipRect
		clipRect := cmd.ClipRect

		x := int32(clipRect[0])
		y := int32(clipRect[1])
		w := int32(clipRect[2])
		h := int32(clipRect[3])

		y = int32(displayHeight) - (y + h)
		// fmt.Printf("type = %s, elems = %d, ofs = %d, texId = %d \n", cmd.Type, cmd.Elems, cmd.IndexOffset, cmd.TexId)
		if cmd.TexId != 0 {
			gl.ActiveTexture(gl.TEXTURE0 + cmd.TexId)
			gl.BindTexture(gl.TEXTURE_2D, cmd.TexId)
			b.shader.UploadTexture("Texture", int32(cmd.TexId))
		}
		gl.Scissor(x, y, w, h)
		gl.DrawElementsBaseVertexWithOffset(gl.TRIANGLES, int32(cmd.Elems), gl.UNSIGNED_INT,
			uintptr(cmd.IndexOffset*4), 0)

	}

	b.shader.Detach()
	gl.DeleteVertexArrays(1, &vaoId)
	// Restore modified GL state
	// gl.UseProgram(uint32(lastProgram))
	// gl.BindTexture(gl.TEXTURE_2D, uint32(lastTexture))
	// gl.BindSampler(0, uint32(lastSampler))
	// gl.ActiveTexture(uint32(lastActiveTexture))
	// gl.BindVertexArray(uint32(lastVertexArray))
	// gl.BindBuffer(gl.ARRAY_BUFFER, uint32(lastArrayBuffer))
	// gl.BindBuffer(gl.ELEMENT_ARRAY_BUFFER, uint32(lastElementArrayBuffer))
	// gl.BlendEquationSeparate(uint32(lastBlendEquationRgb), uint32(lastBlendEquationAlpha))
	// gl.BlendFuncSeparate(uint32(lastBlendSrcRgb), uint32(lastBlendDstRgb), uint32(lastBlendSrcAlpha), uint32(lastBlendDstAlpha))
	// if lastEnableBlend {
	// 	gl.Enable(gl.BLEND)
	// } else {
	// 	gl.Disable(gl.BLEND)
	// }
	// if lastEnableCullFace {
	// 	gl.Enable(gl.CULL_FACE)
	// } else {
	// 	gl.Disable(gl.CULL_FACE)
	// }
	// if lastEnableDepthTest {
	// 	gl.Enable(gl.DEPTH_TEST)
	// } else {
	// 	gl.Disable(gl.DEPTH_TEST)
	// }
	// if lastEnableScissorTest {
	// 	gl.Enable(gl.SCISSOR_TEST)
	// } else {
	gl.Disable(gl.SCISSOR_TEST)
	// }
	// gl.PolygonMode(gl.FRONT_AND_BACK, uint32(lastPolygonMode[0]))
}
