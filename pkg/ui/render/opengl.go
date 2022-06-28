package render

import (

	"math"

	"github.com/Dmitry-dms/moon/pkg/gogl"
	"github.com/go-gl/gl/v4.2-core/gl"
)

type GLRender struct {
	Vertices          []float32
	Indeces           []int32
	vaoId, vboId, ebo uint32
	shader            *gogl.Shader
	triangles         int
	lastIndc          int
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
	s, err := gogl.NewShader("assets/shaders/default.glsl")
	if err != nil {
		panic(err)
	}
	r := GLRender{
		Vertices:  []float32{},
		Indeces:   []int32{},
		vaoId:     0,
		vboId:     0,
		ebo:       0,
		shader:    s,
		triangles: 0,
	}
	r.vaoId = gogl.GenBindVAO()

	//аллоцируем место для vertices
	r.vboId = gogl.GenBindBuffer(gl.ARRAY_BUFFER)
	// gogl.BufferData(gl.ARRAY_BUFFER, r.vertices, gl.DYNAMIC_DRAW)
	gl.BufferData(gl.ARRAY_BUFFER, 4*5000, nil, gl.DYNAMIC_DRAW)

	r.ebo = gogl.GenBindBuffer(gl.ELEMENT_ARRAY_BUFFER)
	gl.BufferData(gl.ELEMENT_ARRAY_BUFFER, 5000*4, nil, gl.DYNAMIC_DRAW)
	// gl.BindBuffer(gl.ELEMENT_ARRAY_BUFFER, 0)

	//включаем layout
	gogl.SetVertexAttribPointer(0, posSize, gl.FLOAT, vertexSize*4, posOffset)
	gogl.SetVertexAttribPointer(1, colorSize, gl.FLOAT, vertexSize*4, colorOffset)
	gogl.SetVertexAttribPointer(2, texCoordsSize, gl.FLOAT, vertexSize*4, texCoordsOffset)
	gogl.SetVertexAttribPointer(3, texIdSize, gl.FLOAT, vertexSize*4, texIdOffset)

	return &r
}

func (r *GLRender) NewFrame() {

}
func (r *GLRender) render(vert []float32, indeces []int32, count int) {
	r.Vertices = append(r.Vertices, vert...)
	r.Indeces = append(r.Indeces, indeces...)
	r.triangles += count
}
func (r *GLRender) Rectangle(x, y, w, h float32, clr [4]float32) {
	// r.Trinagle(x, y, x, y+h, x+w, y+h, clr)
	// r.Trinagle(x+w, y+h, x+w, y, x, y, clr)
	vert := make([]float32, 9*4)
	ind := make([]int32, 6)

	ind0 := r.lastIndc
	ind1 := ind0 + 1
	ind2 := ind1 + 1
	offset := 0

	fillVertices(vert, &offset, x, y, 0, 0, 0, clr)
	fillVertices(vert, &offset, x, y+h, 0, 0, 0, clr)
	fillVertices(vert, &offset, x+w, y+h, 0, 0, 0, clr)

	ind[0] = int32(ind0)
	ind[1] = int32(ind1)
	ind[2] = int32(ind2)

	last := ind2 + 1

	fillVertices(vert, &offset, x+w, y, 0, 0, 0, clr)
	ind[3] = int32(ind0)
	ind[4] = int32(ind2)
	ind[5] = int32(last)

	r.lastIndc = last + 1
	r.render(vert, ind, 4)
}

func fillVertices(vert []float32, startOffset *int, x, y, uv0, uv1, texId float32, clr [4]float32) {
	offset := *startOffset
	vert[offset] = x
	vert[offset+1] = y

	vert[offset+2] = clr[0]
	vert[offset+3] = clr[1]
	vert[offset+4] = clr[2]
	vert[offset+5] = clr[3]

	vert[offset+6] = uv0
	vert[offset+7] = uv1

	vert[offset+8] = texId

	*startOffset += 9
}

type CircleSector int

const (
	BotLeft CircleSector = iota
	BotRight
	TopLeft
	TopRight
)

func (r *GLRender) DrawArc(x, y, radius float32, steps int, sector CircleSector, clr [4]float32) {
	ind0 := r.lastIndc
	ind1 := ind0 + 1
	ind2 := ind1 + 1
	offset := 0
	indOffset := 0

	angle := math.Pi * 2 / float32(steps)

	numV := int(math.Floor(1.57 / float64(angle)))

	ind := make([]int32, 3*(numV+1))    // 3 - triangle
	vert := make([]float32, 9*(3+numV)) //polygon

	var prevX, prevY, lastX, lastY float32

	var ang float32 = angle
	var sX func(x, radius float32) float32
	var sY func(y, radius float32) float32
	// counterTriangles := 0
	switch sector {
	case BotLeft:
		sX = func(x, ang float32) float32 {
			return x - float32(radius)*float32(math.Sin(float64(ang)))
		}
		sY = func(y, ang float32) float32 {
			return y - float32(radius)*float32(math.Cos(float64(ang)))
		}
		prevX = x
		prevY = y - radius
		lastX = x - radius
		lastY = y
	case BotRight:
		sX = func(x, ang float32) float32 {
			return x + float32(radius)*float32(math.Sin(float64(ang)))
		}
		sY = func(y, ang float32) float32 {
			return y - float32(radius)*float32(math.Cos(float64(ang)))
		}
		prevX = x
		prevY = y - radius
		lastX = x + radius
		lastY = y
	case TopLeft:
		sX = func(x, ang float32) float32 {
			return x - float32(radius)*float32(math.Sin(float64(ang)))
		}
		sY = func(y, ang float32) float32 {
			return y + float32(radius)*float32(math.Cos(float64(ang)))
		}
		prevX = x
		prevY = y + radius
		lastX = x - radius
		lastY = y
	case TopRight:
		sX = func(x, ang float32) float32 {
			return x + float32(radius)*float32(math.Sin(float64(ang)))
		}
		sY = func(y, ang float32) float32 {
			return y + float32(radius)*float32(math.Cos(float64(ang)))
		}
		prevX = x
		prevY = y + radius
		lastX = x + radius
		lastY = y
	}

	fillVertices(vert, &offset, x, y, 0, 0, 0, clr)
	fillVertices(vert, &offset, prevX, prevY, 0, 0, 0, clr)
	newx := sX(x, ang)
	newY := sY(y, ang)
	fillVertices(vert, &offset, newx, newY, 0, 0, 0, clr)
	ind[indOffset] = int32(ind0)
	ind[indOffset+1] = int32(ind1)
	ind[indOffset+2] = int32(ind2)
	indOffset += 3
	// ind = append(ind, int32(ind0), int32(ind1), int32(ind2))
	ind1++
	ind2++
	ang += angle

	for ang <= 1.57 { // 90 degress ~= 1.57 radians
		newx := sX(x, ang)
		newY := sY(y, ang)

		fillVertices(vert, &offset, newx, newY, 0, 0, 0, clr)

		ind[indOffset] = int32(ind0)
		ind[indOffset+1] = int32(ind1)
		ind[indOffset+2] = int32(ind2)
		indOffset += 3
		ind1++
		ind2++

		ang += angle
		// counterTriangles++
	}
	fillVertices(vert, &offset, lastX, lastY, 0, 0, 0, clr)

	ind[indOffset] = int32(ind0)
	ind[indOffset+1] = int32(ind1)
	ind[indOffset+2] = int32(ind2)
	// indOffset += 3

	r.lastIndc = ind2 + 1
	r.render(vert, ind, 3*(numV+1))
}

func (r *GLRender) Trinagle(x0, y0, x1, y1, x2, y2 float32, clr [4]float32) {
	vert := make([]float32, 9*3)
	ind := make([]int32, 3)
	offset := 0

	fillVertices(vert, &offset, x0, y0, 0, 0, 0, clr)
	fillVertices(vert, &offset, x1, y1, 0, 0, 0, clr)
	fillVertices(vert, &offset, x2, y2, 0, 0, 0, clr)

	ind0 := r.lastIndc
	ind1 := ind0 + 1
	ind2 := ind1 + 1

	ind[0] = int32(ind0)
	ind[1] = int32(ind1)
	ind[2] = int32(ind2)

	r.lastIndc = ind2 + 1
	r.render(vert, ind, 3)
}

func (r *GLRender) Circle(x, y, radius float32, steps int, clr [4]float32) {
	ind0 := r.lastIndc
	ind1 := ind0 + 1
	ind2 := ind1 + 1
	offset := 0
	indOffset := 0

	angle := math.Pi * 2 / float32(steps)

	numV := int(math.Floor(6.28 / float64(angle)))

	ind := make([]int32, 3*(numV+1))    // 3 - triangle
	vert := make([]float32, 9*(3+numV)) //polygon

	var prevX, prevY float32
	var lastX, lastY float32

	var ang float32 = angle

	prevX = x
	prevY = y + radius

	
	fillVertices(vert, &offset, x, y, 0, 0, 0, clr)
	fillVertices(vert, &offset, prevX, prevY, 0, 0, 0, clr)
	newx := x + radius*float32(math.Sin(float64(ang)))
	newY := y + radius*float32(math.Cos(float64(ang)))
	fillVertices(vert, &offset, newx, newY, 0, 0, 0, clr)
	ind[indOffset] = int32(ind0)
	ind[indOffset+1] = int32(ind1)
	ind[indOffset+2] = int32(ind2)
	indOffset += 3
	ind1++
	ind2++
	ang += angle

	for ang <= 6.28 { // 360 deg ~= 6.28 rad
		newx := x + radius*float32(math.Sin(float64(ang)))
		newY := y + radius*float32(math.Cos(float64(ang)))
		fillVertices(vert, &offset, newx, newY, 0, 0, 0, clr)

		ind[indOffset] = int32(ind0)
		ind[indOffset+1] = int32(ind1)
		ind[indOffset+2] = int32(ind2)
		indOffset += 3
		ind1++
		ind2++

		ang += angle
	}
	lastX = x
	lastY = y + radius
	fillVertices(vert, &offset, lastX, lastY, 0, 0, 0, clr)

	ind[indOffset] = int32(ind0)
	ind[indOffset+1] = int32(ind1)
	ind[indOffset+2] = int32(ind2)
	// indOffset += 3

	r.lastIndc = ind2 + 1
	r.render(vert, ind, 3*(numV+1))
}

func (r *GLRender) Line(x0, y0, x1, y1 float32, thick int, clr [4]float32) {

	r.Trinagle(x0, y0, x1, y1, x1+float32(thick), y1, clr)
	r.Trinagle(x1+float32(thick), y1, x0+float32(thick), y0, x0, y0, clr)
}

var steps = 30

func (r *GLRender) RoundedRectangle(x, y, w, h float32, radius int, clr [4]float32) {

	r.DrawArc(x+float32(radius), y+float32(radius), float32(radius), steps, BotLeft, clr)
	r.DrawArc(x+w-float32(radius), y+float32(radius), float32(radius), steps, BotRight, clr)
	r.DrawArc(x+float32(radius), y+h-float32(radius), float32(radius), steps, TopLeft, clr)
	r.DrawArc(x+w-float32(radius), y+h-float32(radius), float32(radius), steps, TopRight, clr)
	r.Rectangle(x, y+float32(radius), w, h-float32(radius)*2, clr)                                 //center
	r.Rectangle(x+float32(radius), y+h-float32(radius), w-float32(radius)*2, float32(radius), clr) //top
	r.Rectangle(x+float32(radius), y, w-float32(radius)*2, float32(radius), clr)                   //bottom
}

func (b *GLRender) Draw(camera *gogl.Camera) {
	gl.BindBuffer(gl.ARRAY_BUFFER, b.vboId)

	// gl.BufferData(gl.ARRAY_BUFFER, len(b.Vertices)*4, gl.Ptr(b.Vertices), gl.DYNAMIC_DRAW)
	gl.BufferSubData(gl.ARRAY_BUFFER, 0, len(b.Vertices)*4, gl.Ptr(b.Vertices))
	gl.BindBuffer(gl.ARRAY_BUFFER, 0)

	// gl.BindBuffer(gl.ELEMENT_ARRAY_BUFFER, b.ebo)
	// gl.BufferSubData(gl.ELEMENT_ARRAY_BUFFER, 0, len(b.Indeces)*4, gl.Ptr(b.Indeces))
	// gl.BindBuffer(gl.ELEMENT_ARRAY_BUFFER, 0)

	b.shader.Use()

	b.shader.UploadMat4("uProjection", camera.GetProjectionMatrix())
	b.shader.UploadMat4("uView", camera.GetViewMatrix())

	// for i := 0; i < len(b.textures); i++ {
	// 	// gl.ActiveTexture(gl.TEXTURE0 + uint32(i)+1)
	// 	// b.textures[i].Bind()
	// 	b.textures[i].BindActive(gl.TEXTURE0 + uint32(b.texSlots[i]+1))
	// }
	// b.shader.UploadIntArray("uTextures", b.texSlots)

	gogl.BindVertexArray(b.vaoId)
	gl.EnableVertexAttribArray(0)
	gl.EnableVertexAttribArray(1)
	gl.EnableVertexAttribArray(2)
	gl.EnableVertexAttribArray(3)

	// gl.BindBuffer(gl.ELEMENT_ARRAY_BUFFER, b.ebo)
	// gl.DrawArrays(gl.TRIANGLES, 0, int32(b.triangles))

	// gl.DrawElements(gl.TRIANGLES, int32(b.numSprites)*6, gl.UNSIGNED_INT, nil)
	// gl.BindBuffer(gl.ELEMENT_ARRAY_BUFFER, 0)

	gl.BindBuffer(gl.ELEMENT_ARRAY_BUFFER, b.ebo)
	gl.BufferSubData(gl.ELEMENT_ARRAY_BUFFER, 0, len(b.Indeces)*4, gl.Ptr(b.Indeces))
	gl.DrawElements(gl.TRIANGLES, int32(b.triangles), gl.UNSIGNED_INT, nil)
	gl.BindBuffer(gl.ELEMENT_ARRAY_BUFFER, 0)


	gl.DisableVertexAttribArray(0)
	gl.DisableVertexAttribArray(1)
	gl.DisableVertexAttribArray(2)
	gl.DisableVertexAttribArray(3)

	// gl.BindVertexArray(0)
	// for i := 0; i < len(b.textures); i++ {
	// 	b.textures[i].Unbind()
	// }
	b.shader.Detach()
}

func (r *GLRender) End() {
	// r.Vertices = []float32{}
	r.Vertices = make([]float32, 0)
	r.Indeces = []int32{}
	r.lastIndc = 0
}
