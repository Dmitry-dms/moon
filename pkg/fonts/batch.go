package fonts

import (
	// "log"

	"image/color"
	"log"

	"github.com/Dmitry-dms/moon/pkg/gogl"
	"github.com/go-gl/gl/v4.2-core/gl"
	"github.com/go-gl/mathgl/mgl32"
)

// var vertices = []float32{
// 	0.5, 0.5, 1.0, 0.2, 0.11, 1.0, 0.0,
// 	0.5, -0.5, 1.0, 0.2, 0.11, 1.0, 1.0,
// 	-0.5, -0.5, 1.0, 0.2, 0.11, 0.0, 1.0,
// 	-0.5, 0.5, 1.0, 0.2, 0.11, 0.0, 0.0,
// }
var indeces = []int32{
	0, 1, 3,
	1, 2, 3,
}

const (
	vertexSize = 7
	batchSize  = 100 // 25 quads
)

type TextBatch struct {
	Vertices   []float32
	Size       int
	projection mgl32.Mat4
	Vao, Vbo   uint32
	Shader     *gogl.Shader
	Font       *Font
}

func NewTextBatch(font *Font) *TextBatch {
	fontShader, _ := gogl.NewShader("assets/shaders/fonts.glsl")
	tb := TextBatch{
		Vertices:   make([]float32, batchSize*vertexSize),
		Shader:     fontShader,
		Font:       font,
		projection: mgl32.Ident4(),
	}
	return &tb
}

func generateEbo() {
	elementSize := batchSize * 3
	elementBuffer := make([]int32, elementSize)

	for i := range elementBuffer {
		i := int32(i)
		elementBuffer[i] = indeces[(i%6)] + ((i / 6) * 4)
	}

	// gogl.GenEBO()
	// gogl.BufferData(gl.ELEMENT_ARRAY_BUFFER, elementBuffer, gl.STATIC_DRAW)
	ebo = gogl.GenBindBuffer(gl.ELEMENT_ARRAY_BUFFER)
	gogl.BufferData(gl.ELEMENT_ARRAY_BUFFER, elementBuffer, gl.STATIC_DRAW)
	gl.BindBuffer(gl.ELEMENT_ARRAY_BUFFER, 0)
}

var ebo uint32

func (t *TextBatch) Init() {
	t.projection = mgl32.Ortho(0, 1280, 0, 720, 1, 100)

	t.Vao = gogl.GenBindVAO()
	t.Vbo = gogl.GenBindBuffer(gl.ARRAY_BUFFER)
	// gl.BufferData(gl.ARRAY_BUFFER, 4*vertexSize*batchSize, gl.Ptr(t.Vertices), gl.STATIC_DRAW)
	gogl.BufferData(gl.ARRAY_BUFFER, t.Vertices, gl.STATIC_DRAW)

	generateEbo()

	stride := 7 * 4
	gogl.SetVertexAttribPointer(0, 2, gl.FLOAT, stride, 0)
	gogl.SetVertexAttribPointer(1, 3, gl.FLOAT, stride, 2)
	gogl.SetVertexAttribPointer(2, 2, gl.FLOAT, stride, 5)

	gl.BindBuffer(gl.ARRAY_BUFFER, 0)
	gl.BindVertexArray(0)
}

func (t *TextBatch) FlushBatch() {
	//Clear GPU buffer and upload cpu contents and draw
	gl.BindBuffer(gl.ARRAY_BUFFER, t.Vbo)
	gl.BufferData(gl.ARRAY_BUFFER, 4*vertexSize*batchSize, nil, gl.DYNAMIC_DRAW)
	gl.BufferSubData(gl.ARRAY_BUFFER, 0, 4*vertexSize*batchSize, gl.Ptr(t.Vertices))

	//draw
	t.Shader.Use()
	gl.ActiveTexture(gl.TEXTURE15)
	// gl.BindTexture(gl.TEXTURE_BUFFER, t.Font.TextureId)
	gl.BindTexture(gl.TEXTURE_2D, t.Font.TextureId)
	t.Shader.UploadTexture("uFontTexture", 15)
	t.Shader.UploadMat4("uProjection", t.projection)

	gl.BindVertexArray(t.Vao)

	gl.BindBuffer(gl.ELEMENT_ARRAY_BUFFER, ebo)
	gl.DrawElements(gl.TRIANGLES, int32(t.Size*6), gl.UNSIGNED_INT, nil)
	gl.BindBuffer(gl.ELEMENT_ARRAY_BUFFER, 0)

	gl.BindVertexArray(0)
	gl.BindTexture(gl.TEXTURE_2D, 0)

	t.Shader.Detach()
	//reset batch for use on next draw call
	t.Size = 0
}

func (t *TextBatch) AddText(text string, x, y int, scale float32, rgb color.RGBA) {
	var dx, dy int
	dy = y
	dx = x
	prevR := rune(-1)
	faceHeight := t.Font.Face.Metrics().Height
	// fmt.Println("-------------------------------")
	for _, r := range text {
		info := t.Font.GetCharacter(r)
		if info.Width == 0 {
			log.Printf("Unknown char = %q", r)
			continue
		}
		if prevR >= 0 {
			kern := t.Font.Face.Kern(prevR, r).Ceil()
			dx += kern
			// fmt.Printf("%q %q %d \n", prevR, r, kern)
		}
		if r == '\n' {
			dx = x
			dy -= faceHeight.Ceil()
			prevR = rune(-1)
			continue
		}
		xPos := float32(dx)
		yPos := float32(dy)

		if info.Descend != 0 {
			yPos -= float32(info.Descend) * scale
		}

		t.addCharacter(xPos, yPos, scale, *info, rgb)
		dx += info.Width * int(scale)
		prevR = r
	}
	// fmt.Println("-------------------------------")
}

func (t *TextBatch) addCharacter(x, y float32, scale float32, info CharInfo, rgb color.RGBA) {
	//Если нет места, удаляем и начинаем заного
	if t.Size >= batchSize-4 {
		t.FlushBatch()
	}

	// r := float32(((rgb >> 16) & 0xFF) / 255)
	// g := float32(((rgb >> 8) & 0xFF) / 255)
	// b := float32(((rgb >> 0) & 0xFF) / 255)

	x0 := x
	y0 := y
	x1 := x + scale*float32(info.Width)
	y1 := y + scale*float32(info.Height)

	ux0, uy0 := info.TexCoords[0].X, info.TexCoords[0].Y
	ux1, uy1 := info.TexCoords[1].X, info.TexCoords[1].Y

	index := t.Size * 7
	t.Vertices[index] = x1
	t.Vertices[index+1] = y0
	t.Vertices[index+2] = float32(rgb.R)
	t.Vertices[index+3] = float32(rgb.G)
	t.Vertices[index+4] = float32(rgb.B)
	t.Vertices[index+5] = ux1
	t.Vertices[index+6] = uy0

	index += 7
	t.Vertices[index] = x1
	t.Vertices[index+1] = y1
	t.Vertices[index+2] = float32(rgb.R)
	t.Vertices[index+3] = float32(rgb.G)
	t.Vertices[index+4] = float32(rgb.B)
	t.Vertices[index+5] = ux1
	t.Vertices[index+6] = uy1

	index += 7
	t.Vertices[index] = x0
	t.Vertices[index+1] = y1
	t.Vertices[index+2] = float32(rgb.R)
	t.Vertices[index+3] = float32(rgb.G)
	t.Vertices[index+4] = float32(rgb.B)
	t.Vertices[index+5] = ux0
	t.Vertices[index+6] = uy1

	index += 7
	t.Vertices[index] = x0
	t.Vertices[index+1] = y0
	t.Vertices[index+2] = float32(rgb.R)
	t.Vertices[index+3] = float32(rgb.G)
	t.Vertices[index+4] = float32(rgb.B)
	t.Vertices[index+5] = ux0
	t.Vertices[index+6] = uy0

	t.Size += 4

}
