package fonts

import (
	"github.com/Dmitry-dms/moon/pkg/ui/utils"
	"image"
	// "image/png"
	"io/ioutil"
	// "os"

	"log"

	"golang.org/x/image/font"
	"golang.org/x/text/encoding/charmap"

	// "golang.org/x/image/font/gofont/goitalic"

	"golang.org/x/image/font/opentype"
	"golang.org/x/image/math/fixed"
)

type Font struct {
	Filepath string
	FontSize int

	CharMap map[int]*CharInfo

	TextureId uint32

	Face font.Face
}

func NewFont(filepath string, fontSize int) (*Font, *image.RGBA) {
	f := Font{
		Filepath: filepath,
		FontSize: fontSize,
		CharMap:  make(map[int]*CharInfo, 50),
	}

	data := f.generateBitmap()
	return &f, data
}

var siz = 2048

func (f *Font) GetXHeight() float32 {
	c := f.GetCharacter('X')
	return float32(c.Heigth)
}

func (font *Font) CalculateTextBounds(text string, scale float32) (width, height float32, pos []utils.Vec2) {
	prevR := rune(-1)
	inf := font.GetXHeight()
	faceHeight := font.FontSize
	height = scale * inf
	pos = make([]utils.Vec2, len(text))

	var maxDescend, baseline float32
	baseline = scale * inf
	for i, r := range text {
		info := font.GetCharacter(r)
		if info.Width == 0 {
			log.Printf("Unknown char = %q", r)
			continue
		}
		if prevR >= 0 {
			kern := font.Face.Kern(prevR, r).Ceil()
			width += float32(kern)
		}
		if r != ' ' {
			width += float32(info.LeftBearing)
		}
		if r == '\n' {
			width = 0
			height += float32(faceHeight)
			baseline -= float32(faceHeight)
			prevR = rune(-1)
			continue
		}
		xPos := width
		yPos := baseline
		if info.Descend != 0 {
			d := float32(info.Descend) * scale
			yPos += d
			if d > maxDescend {
				maxDescend = d
			}
		}

		pos[i] = utils.Vec2{X: xPos, Y: yPos}

		width += float32(info.Width) * scale
		if r != ' ' {
			width += float32(info.RigthBearing)
		}
		prevR = r
	}
	height += maxDescend
	return
}

func (f *Font) generateBitmap() *image.RGBA {
	cp := charmap.Windows1251
	var letters []rune
	for i := 32; i < 256; i++ {
		r := cp.DecodeByte(byte(i))
		letters = append(letters, r)
	}

	var (
		DPI          = 157.0
		width        = siz
		height       = siz
		startingDotX = 0
		startingDotY = int(f.FontSize) * 2
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
			Hinting: font.HintingNone,
		})

		if err != nil {
			log.Fatalf("NewFace: %v", err)
		}
	}
	f.Face = face
	defer face.Close()

	dst := image.NewRGBA(image.Rect(0, 0, width, height))
	d := font.Drawer{
		Dst:  dst,
		Src:  image.White,
		Face: face,
		Dot:  fixed.P(startingDotX, startingDotY),
	}
	fontSize := d.Face.Metrics().Height
	f.FontSize = fontSize.Ceil()

	dx := startingDotX
	dy := startingDotY
	maxDesc := 0

	//d.DrawString("the quick brown fox jumps over the lazy dog")
	//d.DrawString("Съешь ещё этих мягких французских булок да выпей чаю")
	sortSlice := make([]*CharInfo, len(letters))

	prevDot := d.Dot
	for i, l := range letters {
		b, a, _ := d.Face.GlyphBounds(l)

		//В случае, если ширина символа выходит за границу полотна
		if (siz - dx) <= a.Ceil() {
			dx = 0
			dy += fontSize.Ceil()
			d.Dot = fixed.P(0, dy)
			maxDesc = 0
		}

		dx += a.Ceil() + 2
		d.Dot = d.Dot.Add(fixed.P(2, 0))
		prevDot = d.Dot
		d.DrawString(string(l))

		w, h := (b.Max.X - b.Min.X).Ceil(), (b.Max.Y - b.Min.Y).Ceil()
		//sy := d.Dot.Y.Ceil() - -b.Min.Y.Ceil()
		//sx := d.Dot.X.Ceil() - a.Ceil() + b.Min.X.Ceil()
		w += 1
		sx := prevDot.X.Ceil() + b.Min.X.Ceil() - 2
		sy := prevDot.Y.Ceil() + b.Max.Y.Ceil() - 1

		ch := CharInfo{
			Rune:         l,
			SrcX:         sx,
			SrcY:         sy,
			Width:        w,
			Heigth:       h,
			Ascend:       -b.Min.Y.Ceil(),
			Descend:      b.Max.Y.Ceil(),
			LeftBearing:  b.Min.X.Ceil(),
			RigthBearing: a.Ceil() - b.Max.X.Ceil(),
		}
		sortSlice[i] = &ch
		//fmt.Printf("char = %q, top = %d  bot = %d , h = %d sy = %d \n  ",
		//	l, ch.Ascend, ch.Descend, ch.Heigth, sy)
		//printBorder(dst, ch.SrcX, ch.SrcY, ch.Width, ch.Heigth, prevDot, a, sy)

		if l == ' ' {
			ch.Width = a.Ceil()
		}
		ch.calcTexCoords(siz, siz)

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
	}

	dy += maxDesc

	return dst

	//fil, err := os.OpenFile("atlas.json", os.O_CREATE|os.O_RDWR, 0664)
	//if err != nil {
	//	panic(err)
	//}
	//enc := json.NewEncoder(fil)
	//enc.Encode(sheet.Group)

}
