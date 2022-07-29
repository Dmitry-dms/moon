package gogl

import (
	"image"
	"image/color"
	"image/draw"
	_ "image/jpeg"
	"image/png"
	_ "image/png"
	"os"

	"github.com/go-gl/gl/v4.2-core/gl"
)

type Texture struct {
	Filepath  string `json:"filepath"`
	TextureId uint32 `json:"texture_id"`
	Width     int32  `json:"texture_width"`
	Height    int32  `json:"texture_height"`
}

type TextureExported struct {
	Filepath  string `json:"filepath"`
	TextureId uint32 `json:"texture_id"`
	Width     int32  `json:"texture_width"`
	Height    int32  `json:"texture_height"`
}

func (t *Texture) GetFilepath() string {
	return t.Filepath
}

func CreateTexture(filepath string, texId uint32, width, height int32) *Texture {
	t := &Texture{
		Filepath:  filepath,
		Width:     width,
		Height:    height,
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

func genBindTexture() uint32 {
	var texId uint32
	gl.GenTextures(1, &texId)
	gl.BindTexture(gl.TEXTURE_2D, texId)
	return texId
}

func (t *Texture) Bind() {
	gl.BindTexture(gl.TEXTURE_2D, t.TextureId)
}
func (t *Texture) BindActive(texture uint32) {
	gl.ActiveTexture(texture)
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

func UploadTextureFromMemory(data image.Image) *Texture {
	// p := data.Pix
	w := data.Bounds().Max.X
	h := data.Bounds().Max.Y

	pixels := make([]byte, w*h*4)
	bIndex := 0

	for y := 0; y < h; y++ {
		for x := 0; x < w; x++ {
			r, g, b, a := data.At(x, y).RGBA()

			if r != 0 {
				r = 255
			}
			if g != 0 {
				g = 255
			}
			if b != 0 {
				b = 255
			}

			pixels[bIndex] = byte(r )
			bIndex++
			pixels[bIndex] = byte(g )
			bIndex++
			pixels[bIndex] = byte(b )
			bIndex++
			if r == 0 && g == 0 && b == 0 {
				// pixels[bIndex] = byte(0)
				pixels[bIndex] = byte(a / 256)
			} else {
				pixels[bIndex] = byte(a )
			}
			bIndex++
		}
	}

	texture := genBindTexture()


	gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_WRAP_S, gl.CLAMP_TO_EDGE)
	gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_WRAP_T, gl.CLAMP_TO_EDGE)
	gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_MIN_FILTER, gl.NEAREST)
	gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_MAG_FILTER, gl.NEAREST)
	gl.TexImage2D(gl.TEXTURE_2D, 0, gl.RGBA, int32(w), int32(h), 0, gl.RGBA, gl.UNSIGNED_BYTE, gl.Ptr(pixels))
	// gl.GenerateMipmap(gl.TEXTURE_2D)

	textureStruct := Texture{
		TextureId: texture,
		Width:     int32(w),
		Height:    int32(h),
	}
	textureStruct.Unbind()
	return &textureStruct
}

func TextureFromPNG(filepath string) (*Texture, error) {
	infile, err := os.Open(filepath)
	if err != nil {
		return nil, err
	}
	defer infile.Close()

	img, err := png.Decode(infile)
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
	textureStruct.Unbind()
	return &textureStruct, nil
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

	w := img.Bounds().Dx()
	h := img.Bounds().Dy()

	pixels := make([]byte, w*h*4)
	i := 0
	for y := h - 1; y >= 0; y-- {
		for x := 0; x < w; x++ {
			c := color.NRGBAModel.Convert(img.At(x, y)).(color.NRGBA)
			pixels[i] = c.R
			pixels[i+1] = c.G
			pixels[i+2] = c.B
			pixels[i+3] = c.A

			i += 4
		}
	}

	texture := genBindTexture()

	gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_WRAP_S, gl.REPEAT)
	gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_WRAP_T, gl.REPEAT)
	gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_MIN_FILTER, gl.NEAREST)
	gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_MAG_FILTER, gl.LINEAR)

	gl.TexImage2D(gl.TEXTURE_2D, 0, gl.RGBA, int32(w), int32(h), 0, gl.RGBA, gl.UNSIGNED_BYTE, gl.Ptr(pixels))
	// gl.GenerateMipmap(gl.TEXTURE_2D)

	textureStruct := Texture{
		Filepath:  filepath,
		TextureId: texture,
		Width:     int32(w),
		Height:    int32(h),
	}
	textureStruct.Unbind()
	return &textureStruct, nil
}

func ImageToBytes(img image.Image) []byte {
	size := img.Bounds().Size()
	w, h := size.X, size.Y

	switch img := img.(type) {
	case *image.Paletted:
		bs := make([]byte, 4*w*h)

		b := img.Bounds()
		x0 := b.Min.X
		y0 := b.Min.Y
		x1 := b.Max.X
		y1 := b.Max.Y

		palette := make([]uint8, len(img.Palette)*4)
		for i, c := range img.Palette {
			rgba := color.RGBAModel.Convert(c).(color.RGBA)
			palette[4*i] = rgba.R
			palette[4*i+1] = rgba.G
			palette[4*i+2] = rgba.B
			palette[4*i+3] = rgba.A
		}
		// Even img is a subimage of another image, Pix starts with 0-th index.
		idx0 := 0
		idx1 := 0
		d := img.Stride - (x1 - x0)
		for j := 0; j < y1-y0; j++ {
			for i := 0; i < x1-x0; i++ {
				p := int(img.Pix[idx0])
				bs[idx1] = palette[4*p]
				bs[idx1+1] = palette[4*p+1]
				bs[idx1+2] = palette[4*p+2]
				bs[idx1+3] = palette[4*p+3]
				idx0++
				idx1 += 4
			}
			idx0 += d
		}
		return bs
	case *image.RGBA:
		if len(img.Pix) == 4*w*h {
			return img.Pix
		}
		return imageToBytesSlow(img)
	default:
		return imageToBytesSlow(img)
	}
}

func imageToBytesSlow(img image.Image) []byte {
	size := img.Bounds().Size()
	w, h := size.X, size.Y
	bs := make([]byte, 4*w*h)

	dstImg := &image.RGBA{
		Pix:    bs,
		Stride: 4 * w,
		Rect:   image.Rect(0, 0, w, h),
	}
	draw.Draw(dstImg, image.Rect(0, 0, w, h), img, img.Bounds().Min, draw.Src)
	return bs
}

func flipImageY(stride, height int, pixels []byte) {
	// Flip image in y-direction. OpenGL's origin is in the lower
	// left corner.
	row := make([]uint8, stride)
	for y := 0; y < height/2; y++ {
		y1 := height - y - 1
		dest := y1 * stride
		src := y * stride
		copy(row, pixels[dest:])
		copy(pixels[dest:], pixels[src:src+len(row)])
		copy(pixels[src:], row)
	}
}
