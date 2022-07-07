package render

import (
	"math"

	"github.com/Dmitry-dms/moon/pkg/gogl"
	"github.com/go-gl/gl/v4.2-core/gl"
	"github.com/go-gl/mathgl/mgl32"
)

type GLRender struct {
	Vertices          []float32
	Indeces           []int32
	vaoId, vboId, ebo uint32
	shader            *gogl.Shader
	vertCount         int
	lastIndc          int

	textures []*gogl.Texture
	texSlots []int32
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
		vertCount: 0,
		textures:  make([]*gogl.Texture, 0),
		texSlots:  []int32{0, 1, 2, 3, 4, 5, 6, 7},
	}
	r.vaoId = gogl.GenBindVAO()

	// gl.GenBuffers(1, &r.vboId)
	// gl.GenBuffers(1, &r.ebo)

	//аллоцируем место для vertices
	r.vboId = gogl.GenBindBuffer(gl.ARRAY_BUFFER)
	// gogl.BufferData(gl.ARRAY_BUFFER, r.vertices, gl.DYNAMIC_DRAW)
	// gl.BufferData(gl.ARRAY_BUFFER, 4*1500, nil, gl.DYNAMIC_DRAW)

	r.ebo = gogl.GenBindBuffer(gl.ELEMENT_ARRAY_BUFFER)
	// gl.BindBuffer(gl.ELEMENT_ARRAY_BUFFER, r.ebo)
	// gl.BufferData(gl.ELEMENT_ARRAY_BUFFER, 500*4, nil, gl.DYNAMIC_DRAW)
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
func (r *GLRender) render(vert []float32, indeces []int32, vertCount int) {
	r.Vertices = append(r.Vertices, vert...)
	r.Indeces = append(r.Indeces, indeces...)
	r.vertCount += vertCount
}

func (r *GLRender) addTexture(tex *gogl.Texture) {
	isAdded := false
	for _, v := range r.textures {
		if tex == v {
			isAdded = true
			break
		}
	}
	if !isAdded {
		r.textures = append(r.textures, tex)
	}
}
func (r *GLRender) Rectangle(x, y, w, h float32, clr [4]float32) {
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
func (r *GLRender) RectangleR(x, y, w, h float32, clr [4]float32) {

	vert := make([]float32, 9*4)
	ind := make([]int32, 6)

	ind0 := r.lastIndc
	ind1 := ind0 + 1
	ind2 := ind1 + 1
	offset := 0

	fillVertices(vert, &offset, x, y, 0, 0, 0, clr)
	fillVertices(vert, &offset, x, y-h, 0, 0, 0, clr)
	fillVertices(vert, &offset, x+w, y-h, 0, 0, 0, clr)

	ind[0] = int32(ind0)
	ind[1] = int32(ind1)
	ind[2] = int32(ind2)

	last := ind2 + 1

	fillVertices(vert, &offset, x+w, y, 0, 0, 0, clr)
	ind[3] = int32(ind0)
	ind[4] = int32(ind2)
	ind[5] = int32(last)

	r.lastIndc = last + 1
	r.render(vert, ind, 6)
}

func (r *GLRender) RectangleT(x, y, w, h float32, tex *gogl.Texture, clr [4]float32) {
	founded := false
	texId := 0
	for i := 0; i < len(r.textures); i++ {
		if r.textures[i] == tex {
			texId = i + 1 // 0 - без текстуры
			founded = true
		}
	}
	if !founded {
		r.addTexture(tex)
	}

	vert := make([]float32, 9*4)
	ind := make([]int32, 6)

	ind0 := r.lastIndc
	ind1 := ind0 + 1
	ind2 := ind1 + 1
	offset := 0

	fillVertices(vert, &offset, x, y, 1, 1, float32(texId), clr)
	fillVertices(vert, &offset, x, y-h, 1, 0, float32(texId), clr)
	fillVertices(vert, &offset, x+w, y-h, 0, 0, float32(texId), clr)

	ind[0] = int32(ind0)
	ind[1] = int32(ind1)
	ind[2] = int32(ind2)

	last := ind2 + 1

	fillVertices(vert, &offset, x+w, y, 0, 1, float32(texId), clr)
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

	vert[offset+2] = clr[0] / 255
	vert[offset+3] = clr[1] / 255
	vert[offset+4] = clr[2] / 255
	vert[offset+5] = clr[3]

	vert[offset+6] = uv0
	vert[offset+7] = uv1

	vert[offset+8] = texId

	*startOffset += 9
}

type CircleSector int
type RoundedRectShape int

const (
	TopLeftRect RoundedRectShape = 1 << iota
	TopRigthRect
	BotLeftRect
	BotRightRect

	TopRect = TopLeftRect | TopRigthRect
	BotRect = BotLeftRect | BotRightRect

	AllRounded = TopRect | BotRect
)

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

	vertC := 1
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
		vertC++
		// counterTriangles++
	}
	fillVertices(vert, &offset, lastX, lastY, 0, 0, 0, clr)

	ind[indOffset] = int32(ind0)
	ind[indOffset+1] = int32(ind1)
	ind[indOffset+2] = int32(ind2)
	// indOffset += 3

	r.lastIndc = ind2 + 1

	r.render(vert, ind, (numV+1)*3)
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

// top to down
func (r *GLRender) RoundedRectangleR(x, y, w, h float32, radius int, shape RoundedRectShape, clr [4]float32) {

	topLeft := mgl32.Vec2{x + float32(radius), y - float32(radius)} //origin of arc
	topRight := mgl32.Vec2{x + w - float32(radius), y - float32(radius)}
	botLeft := mgl32.Vec2{x + float32(radius), y - h + float32(radius)}
	botRight := mgl32.Vec2{x + w - float32(radius), y - h + float32(radius)}

	switch shape {
	case TopLeftRect:
		r.DrawArc(topLeft.X(), topLeft.Y(), float32(radius), steps, TopLeft, clr)
		r.RectangleR(x, y-float32(radius), w, h-float32(radius), clr)               //main rect
		r.RectangleR(x+float32(radius), y, w-float32(radius), float32(radius), clr) //top rect
	case TopRigthRect:
		r.DrawArc(topRight.X(), topRight.Y(), float32(radius), steps, TopRight, clr)
		r.RectangleR(x, y-float32(radius), w, h-float32(radius), clr) //main
		r.RectangleR(x, y, w-float32(radius), float32(radius), clr)
	case BotLeftRect:
		r.DrawArc(botLeft.X(), botLeft.Y(), float32(radius), steps, BotLeft, clr)
		r.RectangleR(x, y, w, h-float32(radius), clr) //main
		r.RectangleR(botLeft.X(), botLeft.Y(), w-float32(radius), float32(radius), clr)
	case BotRightRect:
		r.DrawArc(botRight.X(), botRight.Y(), float32(radius), steps, BotRight, clr)
		r.RectangleR(x, y, w, h-float32(radius), clr) //main
		r.RectangleR(x, botLeft.Y(), w-float32(radius), float32(radius), clr)
	case TopRect:
		r.DrawArc(topLeft.X(), topLeft.Y(), float32(radius), steps, TopLeft, clr)
		r.DrawArc(topRight.X(), topRight.Y(), float32(radius), steps, TopRight, clr)
		r.RectangleR(x, y-float32(radius), w, h-float32(radius), clr)
		r.RectangleR(x+float32(radius), y, w-float32(radius)*2, float32(radius), clr)
	case BotRect:
		r.DrawArc(botLeft.X(), botLeft.Y(), float32(radius), steps, BotLeft, clr)
		r.DrawArc(botRight.X(), botRight.Y(), float32(radius), steps, BotRight, clr)
		r.RectangleR(x, y, w, h-float32(radius), clr) //main
		r.RectangleR(botLeft.X(), botLeft.Y(), w-float32(radius)*2, float32(radius), clr)
	case AllRounded:
		r.DrawArc(topLeft.X(), topLeft.Y(), float32(radius), steps, TopLeft, clr)
		r.DrawArc(topRight.X(), topRight.Y(), float32(radius), steps, TopRight, clr)
		r.DrawArc(botLeft.X(), botLeft.Y(), float32(radius), steps, BotLeft, clr)
		r.DrawArc(botRight.X(), botRight.Y(), float32(radius), steps, BotRight, clr)

		r.RectangleR(topLeft.X(), topLeft.Y()+float32(radius), w-float32(radius)*2, float32(radius), clr) //top
		r.RectangleR(x, topLeft.Y(), w, h-float32(radius)*2, clr)                                         //center
		r.RectangleR(botLeft.X(), botLeft.Y(), w-float32(radius)*2, float32(radius), clr)                 //bottom
	}

}

func (r *GLRender) RoundedRectangle(x, y, w, h float32, radius int, clr [4]float32) {

	r.DrawArc(x+float32(radius), y+float32(radius), float32(radius), steps, BotLeft, clr)
	r.DrawArc(x+w-float32(radius), y+float32(radius), float32(radius), steps, BotRight, clr)
	r.DrawArc(x+float32(radius), y+h-float32(radius), float32(radius), steps, TopLeft, clr)
	r.DrawArc(x+w-float32(radius), y+h-float32(radius), float32(radius), steps, TopRight, clr)
	r.Rectangle(x, y+float32(radius), w, h-float32(radius)*2, clr)                                 //center
	r.Rectangle(x+float32(radius), y+h-float32(radius), w-float32(radius)*2, float32(radius), clr) //top
	r.Rectangle(x+float32(radius), y, w-float32(radius)*2, float32(radius), clr)                   //bottom
}

func (r *GLRender) RoundedRectangleT(x, y, w, h float32, radius int, shape RoundedRectShape, tex *gogl.Texture, clr [4]float32) {

	topLeft := mgl32.Vec2{x + float32(radius), y - float32(radius)} //origin of arc
	topRight := mgl32.Vec2{x + w - float32(radius), y - float32(radius)}
	botLeft := mgl32.Vec2{x + float32(radius), y - h + float32(radius)}
	botRight := mgl32.Vec2{x + w - float32(radius), y - h + float32(radius)}

	switch shape {
	case TopLeftRect:
		r.DrawArc(topLeft.X(), topLeft.Y(), float32(radius), steps, TopLeft, clr)
		r.RectangleT(x, y-float32(radius), w, h-float32(radius), tex, clr)               //main rect
		r.RectangleT(x+float32(radius), y, w-float32(radius), float32(radius), tex, clr) //top rect
	case TopRigthRect:
		r.DrawArc(topRight.X(), topRight.Y(), float32(radius), steps, TopRight, clr)
		r.RectangleT(x, y-float32(radius), w, h-float32(radius), tex, clr) //main
		r.RectangleT(x, y, w-float32(radius), float32(radius), tex, clr)
	case BotLeftRect:
		r.DrawArc(botLeft.X(), botLeft.Y(), float32(radius), steps, BotLeft, clr)
		r.RectangleT(x, y, w, h-float32(radius), tex, clr) //main
		r.RectangleT(botLeft.X(), botLeft.Y(), w-float32(radius), float32(radius), tex, clr)
	case BotRightRect:
		r.DrawArc(botRight.X(), botRight.Y(), float32(radius), steps, BotRight, clr)
		r.RectangleT(x, y, w, h-float32(radius), tex, clr) //main
		r.RectangleT(x, botLeft.Y(), w-float32(radius), float32(radius), tex, clr)
	case TopRect:
		r.DrawArc(topLeft.X(), topLeft.Y(), float32(radius), steps, TopLeft, clr)
		r.DrawArc(topRight.X(), topRight.Y(), float32(radius), steps, TopRight, clr)
		r.RectangleT(x, y-float32(radius), w, h-float32(radius), tex, clr)
		r.RectangleT(x+float32(radius), y, w-float32(radius)*2, float32(radius), tex, clr)
	case BotRect:
		r.DrawArc(botLeft.X(), botLeft.Y(), float32(radius), steps, BotLeft, clr)
		r.DrawArc(botRight.X(), botRight.Y(), float32(radius), steps, BotRight, clr)
		r.RectangleT(x, y, w, h-float32(radius), tex, clr) //main
		r.RectangleT(botLeft.X(), botLeft.Y(), w-float32(radius)*2, float32(radius), tex, clr)
	case AllRounded:
		r.DrawArc(topLeft.X(), topLeft.Y(), float32(radius), steps, TopLeft, clr)
		r.DrawArc(topRight.X(), topRight.Y(), float32(radius), steps, TopRight, clr)
		r.DrawArc(botLeft.X(), botLeft.Y(), float32(radius), steps, BotLeft, clr)
		r.DrawArc(botRight.X(), botRight.Y(), float32(radius), steps, BotRight, clr)

		r.RectangleT(topLeft.X(), topLeft.Y(), w-float32(radius)*2, h-float32(radius)*2, tex, clr)        //center
		r.RectangleR(topLeft.X(), topLeft.Y()+float32(radius), w-float32(radius)*2, float32(radius), clr) //top
		r.RectangleR(botLeft.X(), botLeft.Y(), w-float32(radius)*2, float32(radius), clr)                 //bottom
		r.RectangleR(x, topLeft.Y(), float32(radius), h-float32(radius)*2, clr)                           //left
		r.RectangleR(topRight.X(), topRight.Y(), float32(radius), h-float32(radius)*2, clr)               //right
	}

}

func (b *GLRender) Draw(camera *gogl.Camera) {

	// gl.Enable(gl.BLEND)
	// gl.BlendEquation(gl.FUNC_ADD)
	// gl.BlendFunc(gl.SRC_ALPHA, gl.ONE_MINUS_SRC_ALPHA)
	// gl.Disable(gl.CULL_FACE)
	// gl.Disable(gl.DEPTH_TEST)
	// gl.Enable(gl.SCISSOR_TEST)
	// gl.PolygonMode(gl.FRONT_AND_BACK, gl.FILL)

	b.shader.Use()

	gogl.BindVertexArray(b.vaoId)

	// var vaoHandle uint32
	// gl.GenVertexArrays(1, &vaoHandle)
	// gl.BindVertexArray(vaoHandle)

	gl.BindBuffer(gl.ARRAY_BUFFER, b.vboId)
	gl.BufferData(gl.ARRAY_BUFFER, len(b.Vertices)*4, gl.Ptr(b.Vertices), gl.DYNAMIC_DRAW)
	// gl.BindBuffer(gl.ARRAY_BUFFER, 0)

	gl.BindBuffer(gl.ELEMENT_ARRAY_BUFFER, b.ebo)
	gl.BufferData(gl.ELEMENT_ARRAY_BUFFER, len(b.Indeces)*4, gl.Ptr(b.Indeces), gl.DYNAMIC_DRAW)

	b.shader.UploadMat4("uProjection", camera.GetProjectionMatrix())
	b.shader.UploadMat4("uView", camera.GetViewMatrix())

	// for i := 0; i < len(b.textures); i++ {
	// 	b.textures[i].BindActive(gl.TEXTURE0 + uint32(b.texSlots[i]+1))
	// }
	// b.shader.UploadIntArray("uTextures", b.texSlots)

	//

	gl.DrawElements(gl.TRIANGLES, int32(b.vertCount), gl.UNSIGNED_INT, nil)
	// gl.BindBuffer(gl.ELEMENT_ARRAY_BUFFER, 0)

	// gl.BindVertexArray(0)
	// for i := 0; i < len(b.textures); i++ {
	// 	b.textures[i].Unbind()
	// }
	// gl.DeleteVertexArrays(1, &vaoHandle)
	b.shader.Detach()
}

func (r *GLRender) End() {
	r.Vertices = []float32{}
	r.Indeces = []int32{}
	r.lastIndc = 0
	r.vertCount = 0
}
