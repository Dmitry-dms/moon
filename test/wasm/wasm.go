package main

import (
	"github.com/Dmitry-dms/moon/pkg/ui/draw"
	"github.com/nuberu/webgl"
	"github.com/nuberu/webgl/types"
	"unsafe"
)

type WebGLRenderer struct {
	gl *webgl.RenderingContext
	sp *types.Program
}

func newGl(gl *webgl.RenderingContext) *WebGLRenderer {
	// Create a vertex shader object
	vertShader := gl.CreateVertexShader()
	gl.ShaderSource(vertShader, vert)
	gl.CompileShader(vertShader)

	// Create fragment shader object
	fragShader := gl.CreateFragmentShader()
	gl.ShaderSource(fragShader, frag)
	gl.CompileShader(fragShader)

	// Create a shader program object to store the combined shader program
	shaderProgram := gl.CreateProgram()
	gl.AttachShader(shaderProgram, vertShader)
	gl.AttachShader(shaderProgram, fragShader)
	gl.LinkProgram(shaderProgram)

	return &WebGLRenderer{
		gl: gl,
		sp: shaderProgram,
	}
}

const (
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

func (b *WebGLRenderer) NewFrame() {
	// Create a vertex shader object

}
func (b *WebGLRenderer) Scissor(x int32, y int32, w int32, h int32) {
	b.gl.Scissor(int(x), int(y), int(w), int(h))
}

func (b *WebGLRenderer) Draw(displaySize [2]float32, buffer draw.CmdBuffer) {

	displayWidth := displaySize[0]
	displayHeight := displaySize[1]

	projectionMatrix := b.gl.GetUniformLocation(b.sp, "uProjection")

	// Create vertex buffer
	vertexBuffer := b.gl.CreateBuffer()
	b.gl.BindBuffer(webgl.ARRAY_BUFFER, vertexBuffer)
	indexBuffer := b.gl.CreateBuffer()
	b.gl.BindBuffer(webgl.ELEMENT_ARRAY_BUFFER, indexBuffer)

	b.gl.VertexAttribPointer(0, posSize, webgl.FLOAT, false, vertexSize*4, posOffset)
	b.gl.EnableVertexAttribArray(0)
	b.gl.VertexAttribPointer(1, colorSize, webgl.FLOAT, false, vertexSize*4, colorOffset)
	b.gl.EnableVertexAttribArray(1)
	b.gl.VertexAttribPointer(2, texCoordsSize, webgl.FLOAT, false, vertexSize*4, texCoordsOffset)
	b.gl.EnableVertexAttribArray(2)
	b.gl.VertexAttribPointer(3, texIdSize, webgl.FLOAT, false, vertexSize*4, texIdOffset)
	b.gl.EnableVertexAttribArray(3)

	b.gl.Enable(webgl.SCISSOR_TEST)

	b.gl.UseProgram(b.sp)

	newInd := make([]int, len(buffer.Indeces))
	for i, l := range buffer.Indeces {
		newInd[i] = int(l)
	}

	b.gl.BindBuffer(webgl.ARRAY_BUFFER, vertexBuffer)
	b.gl.BufferData(webgl.ARRAY_BUFFER, buffer.Vertices, webgl.STREAM_DRAW)

	b.gl.BindBuffer(webgl.ARRAY_BUFFER, indexBuffer)
	b.gl.BufferDataI(webgl.ARRAY_BUFFER, newInd, webgl.STREAM_DRAW)

	orthoProjection := [4][4]float32{
		{2.0 / displayWidth, 0.0, 0.0, 0.0},
		{0.0, 2.0 / displayHeight, 0.0, 0.0},
		{0.0, 0.0, -2.0, 0.0},
		{-1.0, -1.0, -1.0, 1.0},
	}

	b.gl.ClearColor(0.5, 0.5, 0.5, 0.9)
	b.gl.ClearDepth(1.0)
	b.gl.Viewport(0, 0, int(displayWidth), int(displayHeight))
	b.gl.DepthFunc(webgl.LEQUAL)
	var viewMatrixBuffer *[16]float32
	viewMatrixBuffer = (*[16]float32)(unsafe.Pointer(&orthoProjection))
	b.gl.UniformMatrix4fv(projectionMatrix, false, []float32((*viewMatrixBuffer)[:]))

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
			//b.gl.ActiveTexture( cmd.TexId)
			//gl.BindTexture(gl.TEXTURE_2D, cmd.TexId)
			//b.shader.UploadTexture("Texture", int32(cmd.TexId))
		}
		b.Scissor(x, y, w, h)
		b.gl.DrawElements(webgl.TRIANGLES, cmd.Elems, webgl.UNSIGNED_INT, int64(cmd.IndexOffset*4))
		//b.gl.DrawElementsBaseVertexWithOffset(gl.TRIANGLES, int32(cmd.Elems), gl.UNSIGNED_INT,
		//	uintptr(cmd.IndexOffset*4), 0)

	}

}
