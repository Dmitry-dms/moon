package fonts

import (
	// "fmt"
	// "fmt"
	// "fmt"
	"image"
	"image/draw"
	"image/png"
	"io/ioutil"
	"os"

	"log"

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

	Face font.Face
}

func NewFont(filepath string, size int32, flag bool) *Font {
	f := Font{
		Filepath: filepath,
		FontSize: size,
		CharMap:  make(map[int]*CharInfo, 50),
	}
	opengl = flag

	f.generateAndUploadBitmap()
	return &f
}

var siz = 2048

func (f *Font) generateAndUploadBitmap() {
	cp := charmap.Windows1251
	letters := []rune{}
	for i := 32; i < 256; i++ {
		r := cp.DecodeByte(byte(i))
		letters = append(letters, r)
	}

	var (
		DPI           = 170.0
		width            = siz
		height           = siz
		startingDotX     = 0
		startingDotY     = int(f.FontSize)+int(DPI)/3
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

		// if l == 'i' || l == 'j' {
		// 	fmt.Println(d.Dot.X.Ceil(), dx, b.Min.X.Ceil())
		// }

		dx += a.Ceil()

		d.DrawString(string(l))

		w, h := (b.Max.X - b.Min.X).Ceil(), (b.Max.Y - b.Min.Y).Ceil()
		sy := d.Dot.Y.Ceil() - -b.Min.Y.Ceil()
		sx := d.Dot.X.Ceil() - a.Ceil() + b.Min.X.Ceil() //- (a.Ceil()-b.Max.X.Ceil())
		// w, h := a.Ceil(), -b.Min.Y.Ceil()+b.Max.Y.Ceil()
		ch := CharInfo{
			srcX:         sx,
			srcY:         sy,
			width:        w,
			heigth:       h,
			Ascend:       -b.Min.Y.Ceil(),
			Descend:      b.Max.Y.Ceil(),
			LeftBearing:  b.Min.X.Ceil(),
			RigthBearing: a.Ceil() - b.Max.X.Ceil(),
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
	if opengl {
		t2 := gogl.UploadTextureFromMemory(dst3)
		f.TextureId = t2.GetId()
	}
}

var opengl = true
