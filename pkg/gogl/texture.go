package gogl

import (
	// "encoding/json"
	"image"
	_ "image/png"
	"os"

	"github.com/go-gl/gl/v4.2-core/gl"
)

type Texture struct {
	Filepath  string  `json:"filepath"`
	TextureId uint32  `json:"texture_id"`
	Width     int32 `json:"texture_width"`
	Height    int32 `json:"texture_height"`
}

type TextureExported struct {
	Filepath  string  `json:"filepath"`
	TextureId uint32  `json:"texture_id"`
	Width     int32 `json:"texture_width"`
	Height    int32 `json:"texture_height"`
}

func (t *Texture) GetFilepath() string {
	return t.Filepath
}

func CreateTexture(filepath string, texId uint32, width, height int32) *Texture {
	t := &Texture{
		Filepath: filepath,
		Width: width,
		Height: height,
		TextureId: texId,
	}
	// t.Init(filepath)
	return t
}

func NewTextureFramebuffer(width, height int32) *Texture {
	texture := genBindTexture()

	gl.TexImage2D(gl.TEXTURE_2D, 0, gl.RGB, width, height, 0, gl.RGB, gl.UNSIGNED_BYTE, nil)
	//gl.GenerateMipmap(gl.TEXTURE_2D)
	gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_MIN_FILTER, gl.LINEAR)
	gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_MAG_FILTER, gl.LINEAR)

	textureStruct := Texture{
		Filepath:  "generated",
		TextureId: texture,
		Width:     width,
		Height:    height,
	}
	return &textureStruct
}

func (t *Texture) Init(filepath string) (*Texture, error) {
	infile, err := os.Open(filepath)
	if err != nil {
		return nil, err
	}
	defer infile.Close()

	img, _, err := image.Decode(infile)
	if err != nil {
		return nil, err
	}

	w := img.Bounds().Max.X
	h := img.Bounds().Max.Y

	pixels := make([]byte, w*h*4)
	bIndex := 0
	for y := 0; y < h; y++ {
		for x := 0; x < w; x++ {
			r, g, b, a := img.At(x, y).RGBA()
			pixels[bIndex] = byte(r / 256)
			bIndex++
			pixels[bIndex] = byte(g / 256)
			bIndex++
			pixels[bIndex] = byte(b / 256)
			bIndex++
			pixels[bIndex] = byte(a / 256)
			bIndex++
		}
	}

	texture := genBindTexture()
	gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_WRAP_S, gl.REPEAT)
	gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_WRAP_T, gl.REPEAT)
	gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_MIN_FILTER, gl.NEAREST)
	gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_MAG_FILTER, gl.NEAREST)
	gl.TexImage2D(gl.TEXTURE_2D, 0, gl.RGBA, int32(w), int32(h), 0, gl.RGBA, gl.UNSIGNED_BYTE, gl.Ptr(pixels))
	//gl.GenerateMipmap(gl.TEXTURE_2D)

	textureStruct := Texture{
		Filepath:  filepath,
		TextureId: texture,
		Width:     int32(w),
		Height:    int32(h),
	}
	return &textureStruct, nil
}

func genBindTexture() uint32 {
	var texId uint32
	gl.GenTextures(1, &texId)
	gl.BindTexture(gl.TEXTURE_2D, texId)
	return texId
}

func (t *Texture) Bind() {
	gl.BindTexture(gl.TEXTURE_2D, t.TextureId)
}
func (t *Texture) Unbind() {
	gl.BindTexture(gl.TEXTURE_2D, 0)
}

func (t *Texture) GetWidth() int32 {
	return t.Width
}
func (t *Texture) GetHeight() int32 {
	return t.Height
}

func (t *Texture) GetId() uint32 {
	return t.TextureId
}
