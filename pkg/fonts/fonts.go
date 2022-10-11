package fonts

import (
	"image"
	"image/draw"
	"image/png"
	"os"

	// "image/png"
	"io/ioutil"
	// "os"

	"log"

	"github.com/Dmitry-dms/moon/pkg/gogl"

	"golang.org/x/image/font"
	"golang.org/x/text/encoding/charmap"

	// "golang.org/x/image/font/gofont/goitalic"

	"golang.org/x/image/font/opentype"
	"golang.org/x/image/math/fixed"
)

type Font struct {
	Filepath        string
	DefaultFontSize int32

	CharMap map[int]*CharInfo

	TextureId uint32
	Texture   *gogl.Texture

	Face font.Face
}

const DefaultFontSize = 70

func NewFont(filepath string) *Font {
	f := Font{
		Filepath:        filepath,
		DefaultFontSize: DefaultFontSize,
		CharMap:         make(map[int]*CharInfo, 50),
	}

	f.generateAndUploadBitmap()
	return &f
}

var siz = 2048

func (f *Font) GetXHeight() float32 {
	c := f.GetCharacter('X')
	return float32(c.Heigth)
}

var first = true

func (f *Font) CalculateTextBounds(text string, size int) (width, height float32) {
	//var dx, dy float32
	prevR := rune(-1)

	inf := f.GetXHeight()

	faceHeight := f.Face.Metrics().Height

	scale := 1 / (float32(DefaultFontSize) / float32(size))
	height = scale * inf

	var maxDescend float32
	for _, r := range text {
		info := f.GetCharacter(r)
		if info.Width == 0 {
			log.Printf("Unknown char = %q", r)
			continue
		}
		if prevR >= 0 {
			kern := f.Face.Kern(prevR, r).Ceil()
			width += float32(kern)
			if first {
				//fmt.Printf("%q %.1f \n", r, width)
			}
		}
		if r == '\n' {
			width = 0
			height += float32(faceHeight.Ceil())
			prevR = rune(-1)
			continue
		}

		if info.Descend != 0 {
			d := float32(info.Descend) * scale
			if d > maxDescend {
				maxDescend = d
			}
		}
		width += float32(info.Width) * scale
		prevR = r
	}
	height += maxDescend
	first = false
	return
	//return [2]float32{dx, dy}
}

func (f *Font) generateAndUploadBitmap() {
	cp := charmap.Windows1251
	letters := []rune{}
	for i := 32; i < 256; i++ {
		r := cp.DecodeByte(byte(i))
		letters = append(letters, r)
	}

	var (
		DPI = 144.0
		// DPI          = 256.0
		width        = siz
		height       = siz
		startingDotX = 0
		startingDotY = int(f.DefaultFontSize) + int(DPI)/3
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
			Size:    float64(f.DefaultFontSize),
			DPI:     DPI,
			Hinting: font.HintingFull,
		})

		if err != nil {
			log.Fatalf("NewFace: %v", err)
		}
	}
	f.Face = face
	defer face.Close()

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

		// fmt.Printf("%q %b\n", l, l)

		// fmt.Printf("prev = %q now = %q kern = %d \n", prev, l, face.Kern(prev, l))
		// fmt.Printf("char - %q, sX = %d, width = %d \n", l, d.Dot.X.Ceil(), a.Ceil())

		//В случае, если ширина символа выходит за границу полотна
		if (siz - dx) <= a.Ceil() {
			dx = 0
			dy += fontSize.Ceil()

			d.Dot = fixed.P(0, dy)
			maxDesc = 0
		}
		// Редкий случай, когда символ занимает нижнее пространство другого символа, напр. ij
		if dx > (dx + b.Min.X.Ceil()) {

			d.Dot = fixed.P(dx-b.Min.X.Ceil()*2, dy)
		}

		// special case when 'g' overlaps 'f' in Times New Roman
		if l == 'g' {
			d.Dot = fixed.P(dx+b.Min.X.Ceil()*2, dy)
		}

		// if l == 'i' || l == 'j' {
		// 	fmt.Println(d.Dot.X.Ceil(), dx, b.Min.X.Ceil())
		// }

		dx += a.Ceil()

		// d.Dot = fixed.P(dx, dy)

		d.DrawString(string(l))

		w, h := (b.Max.X - b.Min.X).Ceil(), (b.Max.Y - b.Min.Y).Ceil()
		sy := d.Dot.Y.Ceil() - -b.Min.Y.Ceil()
		sx := d.Dot.X.Ceil() - a.Ceil() + b.Min.X.Ceil()

		ch := CharInfo{
			SrcX:         sx,
			SrcY:         sy,
			Width:        w,
			Heigth:       h,
			Ascend:       -b.Min.Y.Ceil(),
			Descend:      b.Max.Y.Ceil(),
			LeftBearing:  b.Min.X.Ceil(),
			RigthBearing: a.Ceil() - b.Max.X.Ceil(),
		}
		if l == ' ' {
			ch.Width = a.Ceil()
		}
		ch.calcTexCoords(siz, siz)
		// draw.Draw(dst, image.Rect(sx, sy, sx+w, sy+h), Border, image.ZP, draw.Src)
		// width += dx
		if -b.Min.Y.Ceil() > startingDotY {
			startingDotY = -b.Min.Y.Ceil()
		}

		// Если символ не будет найден, вместо него отдаем пустой квадрат
		if l == '\u007f' {
			f.CharMap[CharNotFound] = &ch
		} else {
			f.CharMap[int(l)] = &ch
		}

		if b.Max.Y.Ceil() > maxDesc {
			maxDesc = b.Max.Y.Ceil()
		}
		// prev = l
	}
	dy += maxDesc

	dst2 := dst.SubImage(image.Rect(0, 0, siz, siz))
	dst3 := image.NewRGBA(dst2.Bounds())
	// dst3 := image.NewGray(dst2.Bounds())
	draw.Draw(dst3, dst2.Bounds(), dst2, image.ZP, draw.Src)
	// ng := image.NewRGBA(dst2.Bounds())
	// dst2 = ng

	pngFile, _ := os.OpenFile("e.png", os.O_CREATE|os.O_RDWR, 0664)

	defer pngFile.Close()

	encoder := png.Encoder{
		CompressionLevel: png.BestCompression,
	}
	encoder.Encode(pngFile, dst3)
	// if opengl {
	t2 := gogl.UploadTextureFromMemory(dst3)
	f.TextureId = t2.GetId()
	f.Texture = t2
	// }
}

var opengl = true
