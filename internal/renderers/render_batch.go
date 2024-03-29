package renderers

import (
	// "fmt"
	"sync"

	"github.com/Dmitry-dms/moon/internal/components"

	"github.com/Dmitry-dms/moon/pkg/gogl"
	"github.com/go-gl/gl/v4.2-core/gl"

)

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

type RenderBatch struct {
	objects    []*components.GameObject
	numSprites int
	hasRoom    bool

	vertices []float32
	indeces  []int32

	vaoId, vboId, ebo uint32
	maxBatchSize      int

	shader   *gogl.Shader
	textures []*gogl.Texture
	texSlots []int32

	zIndex int
}

func NewRenderBatch(maxBatchSize, zIndex int) *RenderBatch {
	s, err := gogl.AssetPool.GetShader("assets/shaders/default.glsl")
	if err != nil {
		panic(err)
	}
	obj := make([]*components.GameObject, 0, maxBatchSize)
	vt := make([]float32, maxBatchSize*4*vertexSize)
	rb := RenderBatch{
		shader:       s,
		maxBatchSize: maxBatchSize,
		vertices:     vt,
		hasRoom:      true,
		objects:      obj,
		textures:     make([]*gogl.Texture, 0),
		texSlots:     []int32{0, 1, 2, 3, 4, 5, 6, 7},
		zIndex:       zIndex,
	}
	rb.indeces = rb.generateIndeces()

	return &rb
}
func (b *RenderBatch) GetZIndex() int {
	return b.zIndex
}

//работа с OpenGL
func (b *RenderBatch) Start() {
	// fmt.Println("START BATCH")
	b.vaoId = gogl.GenBindVAO()

	//аллоцируем место для vertices
	b.vboId = gogl.GenBindBuffer(gl.ARRAY_BUFFER)
	gogl.BufferData(gl.ARRAY_BUFFER, b.vertices, gl.DYNAMIC_DRAW)

	//включаем layout
	gogl.SetVertexAttribPointer(0, posSize, gl.FLOAT, vertexSize, posOffset)
	gogl.SetVertexAttribPointer(1, colorSize, gl.FLOAT, vertexSize, colorOffset)
	gogl.SetVertexAttribPointer(2, texCoordsSize, gl.FLOAT, vertexSize, texCoordsOffset)
	gogl.SetVertexAttribPointer(3, texIdSize, gl.FLOAT, vertexSize, texIdOffset)

	b.ebo = gogl.GenBindBuffer(gl.ELEMENT_ARRAY_BUFFER)
	gogl.BufferData(gl.ELEMENT_ARRAY_BUFFER, b.indeces, gl.STATIC_DRAW)
	gl.BindBuffer(gl.ELEMENT_ARRAY_BUFFER, 0)

}

func (b *RenderBatch) Update(dt float32, wg *sync.WaitGroup) {
	for _, obj := range b.objects {
		obj.Update(dt)
	}
	wg.Done()
}

func (b *RenderBatch) Render(camera *gogl.Camera) {
	rebuffer := false
	for i, obj := range b.objects {
		if obj.IsDirty() {
		
			b.loadVertexProperties(i)
			obj.SetClean()
			rebuffer = true
		}
	}

	if rebuffer {
		gl.BindBuffer(gl.ARRAY_BUFFER, b.vboId)
		gl.BufferSubData(gl.ARRAY_BUFFER, 0, len(b.vertices)*4, gl.Ptr(b.vertices))
		gl.BindBuffer(gl.ARRAY_BUFFER, 0)
	}

	b.shader.Use()


	b.shader.UploadMat4("uProjection", camera.GetProjectionMatrix())
	b.shader.UploadMat4("uView", camera.GetViewMatrix())

	for i := 0; i < len(b.textures); i++ {
		gl.ActiveTexture(gl.TEXTURE0 + uint32(i) + 1)
		b.textures[i].Bind()
	}
	b.shader.UploadIntArray("uTextures", b.texSlots)

	gogl.BindVertexArray(b.vaoId)
	gl.EnableVertexAttribArray(0)
	gl.EnableVertexAttribArray(1)
	gl.EnableVertexAttribArray(2)
	gl.EnableVertexAttribArray(3)

	gl.BindBuffer(gl.ELEMENT_ARRAY_BUFFER, b.ebo)
	gl.DrawElements(gl.TRIANGLES, int32(b.numSprites)*6, gl.UNSIGNED_INT, nil)
	gl.BindBuffer(gl.ELEMENT_ARRAY_BUFFER, 0)

	gl.DisableVertexAttribArray(0)
	gl.DisableVertexAttribArray(1)
	gl.DisableVertexAttribArray(2)
	gl.DisableVertexAttribArray(3)

	gl.BindVertexArray(0)
	for i := 0; i < len(b.textures); i++ {
		b.textures[i].Unbind()
	}
	b.shader.Detach()
}

func (b *RenderBatch) AddGameObject(obj *components.GameObject) {
	index := b.numSprites
	b.objects = append(b.objects, obj)
	b.numSprites++
	tex := obj.Spr.GetTexture()

	if tex != nil {
		isAdded := false
		for _, v := range b.textures {
			if tex == v {
				isAdded = true
				break
			}
		}
		if !isAdded {
			b.textures = append(b.textures, tex)
		}
	}
	b.loadVertexProperties(index)
	if b.numSprites >= b.maxBatchSize {
		b.hasRoom = false
	}
}

func (b *RenderBatch) loadVertexProperties(index int) {
	obj := b.objects[index]


	offset := index * 4 * int(vertexSize)

	spr := obj.Spr
	color := spr.GetColor()
	tex := spr.GetTexture()

	texId := 0
	if tex != nil {
		for i := 0; i < len(b.textures); i++ {
			if b.textures[i] == tex {
				texId = i + 1 // 0 - без текстуры

			}
		}
	}

	var xAdd float32 = 1
	var yAdd float32 = 1

	for i := 0; i < 4; i++ {
		if i == 1 {
			yAdd = 0
		} else if i == 2 {
			xAdd = 0
		} else if i == 3 {
			yAdd = 1
		}
		//load position
		x := obj.Transform.GetPosition().X() + (xAdd * obj.Transform.GetScale().X())
		y := obj.Transform.GetPosition().Y() + (yAdd * obj.Transform.GetScale().Y())
		// fmt.Printf("BATCH - %.1f %.1f \n",x, y)
		b.vertices[offset] = x
		b.vertices[offset+1] = y
		//load color
		b.vertices[offset+2] = color.X()
		b.vertices[offset+3] = color.Y()
		b.vertices[offset+4] = color.Z()
		b.vertices[offset+5] = color.W()

		//load texCoords
		b.vertices[offset+6] = spr.GetTextureCoords()[i].X()
		b.vertices[offset+7] = spr.GetTextureCoords()[i].Y()
		//load texId
		b.vertices[offset+8] = float32(texId)

		offset += vertexSize
	}
}

func (b *RenderBatch) generateIndeces() []int32 {
	//6 вершин на 1 квадрат
	elements := make([]int32, 6*b.maxBatchSize)

	for i := 0; i < b.maxBatchSize; i++ {
		b.loadElementIndeces(elements, int32(i))
	}

	return elements
}
func (b *RenderBatch) loadElementIndeces(elements []int32, index int32) {
	var offsetArrayIndex int32 = 6 * index
	var offset int32 = 4 * index
	// 3, 2, 0, 0, 2, 1 - вершины квадрата
	elements[offsetArrayIndex] = offset + 3
	elements[offsetArrayIndex+1] = offset + 2
	elements[offsetArrayIndex+2] = offset + 0

	elements[offsetArrayIndex+3] = offset + 0
	elements[offsetArrayIndex+4] = offset + 2
	elements[offsetArrayIndex+5] = offset + 1
}

func (b *RenderBatch) HasTextureRoom() bool {
	return len(b.textures) < 8
}
func (b *RenderBatch) HasTexture(texture *gogl.Texture) bool {
	for _, v := range b.textures {
		if texture == v {
			return true
		}
	}
	return false
}
