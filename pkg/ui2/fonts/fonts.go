package fonts

import (
	"errors"
	"fmt"
	"image"
	"image/draw"
	"image/png"
	"os"

	// "image/color"
	_ "image/png"
	"io/ioutil"
	"log"

	// "os"

	"github.com/Dmitry-dms/moon/pkg/gogl"
	"github.com/go-gl/gl/v4.2-core/gl"
	"golang.org/x/image/font"
	"golang.org/x/text/encoding/charmap"

	// "golang.org/x/image/font/gofont/goitalic"
	"golang.org/x/image/font/opentype"
	"golang.org/x/image/math/fixed"
)

type Font struct {
	Filepath string
	FontSize int32

	CharMap map[int]*CharInfo

	TextureId uint32
}

func NewFont(filepath string, size int32) *Font {
	f := Font{
		Filepath: filepath,
		FontSize: size,
		CharMap:  make(map[int]*CharInfo, 50),
	}

	f.generatebitmap2()
	return &f
}

func (f *Font) generatebitmap() {
	opt := Options{
		FontSize: int(f.FontSize),
		DPI:      72,
		Measure:  true, //качество лучше, зависит от размеров шрифта
		Codepage: ASCII,
	}

	// ParseTTF(opt, f.Filepath, 32, 127)

	// err := Convert(opt, f.Filepath, "out.png")
	// if err != nil {
	// 	panic(err)
	// }
	lowChar := 32
	highChar := 127
	switch opt.Codepage {
	case ASCII:
		highChar = 127
	case CP1252:
		highChar = 255
	default:
		fmt.Errorf("invalid TIGR codepage: %v", opt.Codepage)
	}
	f.FromTTF(opt, f.Filepath, lowChar, highChar)
}

func (f *Font) generatebitmap2() {
	cp := charmap.Windows1252
	letters := []rune{}
	for i := 32; i < 127; i++ {
		r := cp.DecodeByte(byte(i))
		letters = append(letters, r)
	}

	const (
		width        = 50
		height       = 36
		startingDotX = 0
		startingDotY = 28
	)
	ttfBytes, err := ioutil.ReadFile(f.Filepath)
	if err != nil {
		panic(err)
	}

	parsed, err := opentype.Parse(ttfBytes)
	if err != nil {
		log.Fatalf("Parse: %v", err)
	}
	face, err := opentype.NewFace(parsed, &opentype.FaceOptions{
		Size:    32,
		DPI:     72,
		Hinting: font.HintingNone,
	})
	if err != nil {
		log.Fatalf("NewFace: %v", err)
	}

	dst := image.NewGray(image.Rect(0, 0, width, height))
	d := font.Drawer{
		Dst:  dst,
		Src:  image.White,
		Face: face,
		Dot:  fixed.P(startingDotX, startingDotY),
	}
	//check width
	width2 := 0
	for _, l := range letters {
		d.DrawString(string(l))
		// fmt.Println(d.Dot.X.Ceil())
		_, a, _ := d.Face.GlyphBounds(l)
		// fmt.Println(b.Max.X.Ceil())
		width2 += a.Ceil()
	}
	width2+=1
	fmt.Println(width2, height)
	dst2 := image.NewGray(image.Rect(0, 0, width2, height))
	d2 := font.Drawer{
		Dst:  dst2,
		Src:  image.White,
		Face: face,
		Dot:  fixed.P(startingDotX, startingDotY),
	}
	for _, l := range letters {
		d2.DrawString(string(l))
		// fmt.Println(d.Dot.X.Ceil())
		b, a, _ := d2.Face.GlyphBounds(l)

		ch := CharInfo{
			srcX:   d2.Dot.X.Ceil() - a.Ceil(), // начало
			srcY:   d2.Dot.Y.Ceil(),
			width:  a.Ceil(),
			heigth: b.Max.Y.Ceil() - b.Min.Y.Ceil(),
		}
		ch.calcTexCoords(width2, height)
		// fmt.Printf("char - %q, x - %d, y - %d, w - %d, h - %d \n", l, ch.srcX, ch.srcY, ch.width, ch.heigth)
		f.CharMap[int(l)] = &ch
	}



	// d.DrawString(string(letters))

	pngFile, _ := os.OpenFile("e.png", os.O_CREATE|os.O_RDWR, 0664)

	defer pngFile.Close()

	encoder := png.Encoder{
		CompressionLevel: png.NoCompression,
	}
	err = encoder.Encode(pngFile, dst2)
	// f.uploadTexture((*image.NRGBA)(dst2))
	
	// gogl.TextureFromPNG("assets/images/blend2.png")
	// gogl.TextureFromPNG("assets/images/mario.png")
	// gogl.TextureFromPNG("assets/images/img.png")
	// gogl.TextureFromPNG("assets/images/decorations.png")
	// t, _ := gogl.TextureFromPNG("e.png")
	// t, _ := NewTextureFromFile("assets/images/blend1.png", gl.REPEAT, gl.REPEAT)
	// f.TextureId = t.handle

	t2 := &gogl.Texture{}

	t2 = gogl.UploadTextureFromMemory(dst2)
	_ = t2
	// _ = t
	f.TextureId = t2.GetId()
}

func (f *Font) uploadTexture(img *image.NRGBA) {
	t := &gogl.Texture{}

	t, _ = t.Init("e.png")
	f.TextureId = t.GetId()
	// data := img.Pix

	// var texId uint32
	// gl.GenTextures(1, &texId)
	// gl.BindTexture(gl.TEXTURE_2D, texId)

	// gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_WRAP_S, gl.REPEAT)
	// gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_WRAP_T, gl.REPEAT)
	// gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_MAG_FILTER, gl.LINEAR)
	// gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_MIN_FILTER, gl.LINEAR)

	// gl.TexImage2D(gl.TEXTURE_2D, 0, gl.RGBA, int32(img.Bounds().Dx()),
	// 	int32(img.Bounds().Dy()), 0, gl.RGBA, gl.UNSIGNED_BYTE, gl.Ptr(data))

	// f.TextureId = texId
}

type Texture struct {
	handle  uint32
	target  uint32 // same target as gl.BindTexture(<this param>, ...)
	texUnit uint32 // Texture unit that is currently bound to ex: gl.TEXTURE0
}

var errUnsupportedStride = errors.New("unsupported stride, only 32-bit colors supported")

var errTextureNotBound = errors.New("texture not bound")

func NewTextureFromFile(file string, wrapR, wrapS int32) (*Texture, error) {
	imgFile, err := os.Open(file)
	if err != nil {
		return nil, err
	}
	defer imgFile.Close()

	// Decode detexts the type of image as long as its image/<type> is imported
	img, _, err := image.Decode(imgFile)
	if err != nil {
		return nil, err
	}
	return NewTexture(img, wrapR, wrapS)
}

func NewTexture(img image.Image, wrapR, wrapS int32) (*Texture, error) {
	rgba := image.NewRGBA(img.Bounds())
	draw.Draw(rgba, rgba.Bounds(), img, image.Pt(0, 0), draw.Src)
	if rgba.Stride != rgba.Rect.Size().X*4 { // TODO-cs: why?
		return nil, errUnsupportedStride
	}

	var handle uint32
	gl.GenTextures(1, &handle)

	target := uint32(gl.TEXTURE_2D)
	internalFmt := int32(gl.RGBA)
	format := uint32(gl.RGBA)
	width := int32(rgba.Rect.Size().X)
	height := int32(rgba.Rect.Size().Y)
	pixType := uint32(gl.UNSIGNED_BYTE)
	dataPtr := gl.Ptr(rgba.Pix)

	texture := Texture{
		handle: handle,
		target: target,
	}

	texture.Bind(gl.TEXTURE0)
	defer texture.UnBind()

	// set the texture wrapping/filtering options (applies to current bound texture obj)
	// TODO-cs
	gl.TexParameteri(texture.target, gl.TEXTURE_WRAP_R, wrapR)
	gl.TexParameteri(texture.target, gl.TEXTURE_WRAP_S, wrapS)
	gl.TexParameteri(texture.target, gl.TEXTURE_MIN_FILTER, gl.LINEAR) // minification filter
	gl.TexParameteri(texture.target, gl.TEXTURE_MAG_FILTER, gl.LINEAR) // magnification filter

	gl.TexImage2D(target, 0, internalFmt, width, height, 0, format, pixType, dataPtr)

	gl.GenerateMipmap(texture.handle)

	return &texture, nil
}

func (tex *Texture) Bind(texUnit uint32) {
	gl.ActiveTexture(texUnit)
	gl.BindTexture(tex.target, tex.handle)
	tex.texUnit = texUnit
}

func (tex *Texture) UnBind() {
	tex.texUnit = 0
	gl.BindTexture(tex.target, 0)
}

func (tex *Texture) SetUniform(uniformLoc int32) error {
	if tex.texUnit == 0 {
		return errTextureNotBound
	}
	gl.Uniform1i(uniformLoc, int32(tex.texUnit-gl.TEXTURE0))
	return nil
}

func loadImageFile(file string) (image.Image, error) {
	infile, err := os.Open(file)
	if err != nil {
		return nil, err
	}
	defer infile.Close()

	// Decode automatically figures out the type of immage in the file
	// as long as its image/<type> is imported
	img, _, err := image.Decode(infile)
	return img, err
}
