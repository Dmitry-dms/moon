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

	f.generateAndUploadBitmap()
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

	var (
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
		Size:    float64(f.FontSize),
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
		b, a, _ := d.Face.GlyphBounds(l)
		// fmt.Println(b.Max.X.Ceil())
		// fmt.Printf("char = %q, minX = %d, minY = %d, maxX = %d, maxY = %d, a = %d \n",
		// 	l, b.Min.X.Ceil(), b.Min.Y.Ceil(), b.Max.X.Ceil(), b.Max.Y.Ceil(), a.Ceil())
		width2 += a.Ceil()
		if -b.Min.Y.Ceil() > startingDotY {
			startingDotY = -b.Min.Y.Ceil()
		}
	}
	height = d.Face.Metrics().Height.Ceil()
	dst2 := image.NewGray(image.Rect(0, 0, width2, height))
	d2 := font.Drawer{
		Dst:  dst2,
		Src:  image.White,
		Face: face,
		Dot:  fixed.P(startingDotX, startingDotY),
	}

	for _, l := range letters {

		b, a, _ := d2.Face.GlyphBounds(l)

		d2.DrawString(string(l))
		// fmt.Println(d.Dot.X.Ceil())

		_ = b
		// ch := CharInfo{
		// 	srcX:   d2.Dot.X.Ceil() - a.Ceil(),
		// 	srcY:   d2.Dot.Y.Ceil(),
		// 	width:  a.Ceil(),
		// 	heigth: -b.Min.Y.Ceil() + b.Max.Y.Ceil(),
		// }
		w, h := (b.Max.X - b.Min.X).Ceil(), (b.Max.Y - b.Min.Y).Ceil()
		// if b.Min.X&((1<<6)-1) != 0 {
		// 	w++
		// }
		// if b.Min.Y&((1<<6)-1) != 0 {
		// 	h++
		// }
		sx := d2.Dot.X.Ceil() - a.Ceil()
		sy := d2.Dot.Y.Ceil() - -b.Min.Y.Ceil()
		// w, h := a.Ceil(), -b.Min.Y.Ceil()+b.Max.Y.Ceil()
		ch := CharInfo{
			srcX:   sx,
			srcY:   sy,
			width:  w,
			heigth: h,
		}
		ch.calcTexCoords(width2, height)
		// if l == 'Q' {
		// 	fmt.Printf("char - %q, x - %d, y - %d, w - %d, h - %d \n", l, ch.srcX, ch.srcY, ch.width, ch.heigth)
		// }
		f.CharMap[int(l)] = &ch
	}

	// d.DrawString(string(letters))

	pngFile, _ := os.OpenFile("e.png", os.O_CREATE|os.O_RDWR, 0664)

	defer pngFile.Close()

	encoder := png.Encoder{
		CompressionLevel: png.BestCompression,
	}
	err = encoder.Encode(pngFile, dst2)

	t2 := &gogl.Texture{}

	t2 = gogl.UploadTextureFromMemory(dst2)
	_ = t2
	// _ = t
	f.TextureId = t2.GetId()
}

var siz = 512

func (f *Font) generateAndUploadBitmap() {
	cp := charmap.Windows1252
	letters := []rune{}
	for i := 32; i < 127; i++ {
		r := cp.DecodeByte(byte(i))
		letters = append(letters, r)
	}

	var (
		width        = siz
		height       = siz
		startingDotX = 0
		startingDotY = int(f.FontSize)
	)
	var face font.Face
	{
		ttfBytes, err := ioutil.ReadFile(f.Filepath)
		if err != nil {
			panic(err)
		}

		parsed, err := opentype.Parse(ttfBytes)
		if err != nil {
			log.Fatalf("Parse: %v", err)
		}
		face, err = opentype.NewFace(parsed, &opentype.FaceOptions{
			Size:    float64(f.FontSize),
			DPI:     72,
			Hinting: font.HintingFull,
		})
		if err != nil {
			log.Fatalf("NewFace: %v", err)
		}
	}

	dst := image.NewGray(image.Rect(0, 0, width, height))
	d := font.Drawer{
		Dst:  dst,
		Src:  image.White,
		Face: face,
		Dot:  fixed.P(startingDotX, startingDotY),
	}
	fontSize := d.Face.Metrics().Height
	//check width
	// width2 := 0
	dx := startingDotX
	dy := startingDotY
	maxDesc := 0

	for _, l := range letters {
		b, a, _ := d.Face.GlyphBounds(l)

		// fmt.Printf("char - %q, sX = %d, width = %d \n", l, d.Dot.X.Ceil(), a.Ceil())

		if (siz - dx) <= a.Ceil() {
			dx = 0
			dy += fontSize.Ceil()

			d.Dot = fixed.P(0, dy)
			maxDesc = 0
		}
		dx += a.Ceil()
		
		d.DrawString(string(l))

		w, h := (b.Max.X - b.Min.X).Ceil(), (b.Max.Y - b.Min.Y).Ceil()
		sy := d.Dot.Y.Ceil() - -b.Min.Y.Ceil()
		sx := d.Dot.X.Ceil() - a.Ceil() + b.Min.X.Ceil() //- (a.Ceil()-b.Max.X.Ceil())
		// w, h := a.Ceil(), -b.Min.Y.Ceil()+b.Max.Y.Ceil()
		ch := CharInfo{
			srcX:   sx,
			srcY:   sy,
			width:  w,
			heigth: h,
		}
		// draw.Draw(dst, image.Rect(sx, sy, sx+w, sy+h), Border, image.ZP, draw.Src)
		// width += dx
		if -b.Min.Y.Ceil() > startingDotY {
			startingDotY = -b.Min.Y.Ceil()
		}
		f.CharMap[int(l)] = &ch

		if b.Max.Y.Ceil() > maxDesc {
			maxDesc = b.Max.Y.Ceil()
		}
	}
	dy += maxDesc
	fmt.Println(dy)

	// fmt.Println(width, height)
	dst2 := dst.SubImage(image.Rect(0, 0, siz, siz))
	dst3 := image.NewGray(dst2.Bounds())
	draw.Draw(dst3, dst2.Bounds(), dst2, image.ZP, draw.Src)
	ng := image.NewRGBA(dst2.Bounds())
	// dst2 = ng

	for _, v := range f.CharMap {
		v.calcTexCoords(siz, siz)
		draw.Draw(ng, image.Rect(v.srcX, v.srcY, v.srcX+v.width, v.srcY+v.heigth), Border, image.ZP, draw.Src)
	}

	pngFile, _ := os.OpenFile("e.png", os.O_CREATE|os.O_RDWR, 0664)

	defer pngFile.Close()

	encoder := png.Encoder{
		CompressionLevel: png.BestCompression,
	}
	encoder.Encode(pngFile, dst3)

	t2 := gogl.UploadTextureFromMemory(dst3)
	f.TextureId = t2.GetId()
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
