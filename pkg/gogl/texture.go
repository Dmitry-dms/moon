package gogl

import (

	_ "image/png"
	"image"
	"os"


	"github.com/go-gl/gl/v4.2-core/gl"

)

type Texture struct {
	filepath  string
	TextureId uint32
	width, height int
}

func LoadTextureAlpha(filepath string) (*Texture, error) {
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
	gl.GenerateMipmap(gl.TEXTURE_2D)

	textureStruct := Texture{
		filepath:  filepath,
		TextureId: texture,
		width: w,
		height: h,
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

func (t *Texture) GetWidth() int {
	return t.width
}
func (t *Texture) GetHeight() int {
	return t.height
}
