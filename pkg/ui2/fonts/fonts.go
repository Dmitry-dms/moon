package fonts

import (
	"image"
	"image/draw"
	"image/png"
	"os"

	"io/ioutil"
	"log"

	// "github.com/Dmitry-dms/moon/pkg/gogl"

	"github.com/Dmitry-dms/moon/pkg/gogl"
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
		if l == ' ' {
			ch.width = a.Ceil()
		}
		ch.calcTexCoords(siz, siz)
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

	// fmt.Println(width, height)
	dst2 := dst.SubImage(image.Rect(0, 0, siz, siz))
	dst3 := image.NewGray(dst2.Bounds())
	draw.Draw(dst3, dst2.Bounds(), dst2, image.ZP, draw.Src)
	// ng := image.NewRGBA(dst2.Bounds())
	// dst2 = ng

	// for _, v := range f.CharMap {
	// 	v.calcTexCoords(siz, siz)
	// 	// draw.Draw(ng, image.Rect(v.srcX, v.srcY, v.srcX+v.width, v.srcY+v.heigth), Border, image.ZP, draw.Src)
	// }

	pngFile, _ := os.OpenFile("e.png", os.O_CREATE|os.O_RDWR, 0664)

	defer pngFile.Close()

	encoder := png.Encoder{
		CompressionLevel: png.BestCompression,
	}
	encoder.Encode(pngFile, dst3)

	t2 := gogl.UploadTextureFromMemory(dst3)
	f.TextureId = t2.GetId()
}


