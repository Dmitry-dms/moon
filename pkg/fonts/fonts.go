package fonts

import (
	"fmt"
	"github.com/Dmitry-dms/moon/pkg/math"
	"github.com/Dmitry-dms/moon/pkg/ui/utils"
	"image"
	"image/color"
	"image/draw"

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
	Filepath  string
	FontSize  int
	CharMap   map[int]*CharInfo
	CharSlice []*CharInfo
	TextureId uint32
	Face      font.Face
}

func NewFont(filepath string, fontSize int, dpi float32, from, to int) (*Font, *image.RGBA) {
	f := Font{
		Filepath: filepath,
		FontSize: fontSize,
		CharMap:  make(map[int]*CharInfo, 50),
	}
	data := f.generateBitmap(dpi, from, to)
	return &f, data
}

var AtlasWidth = 1024

func (f *Font) GetXHeight() float32 {
	c := f.GetCharacter('X')
	return float32(c.Heigth)
}

type CombinedCharInfo struct {
	Char  CharInfo
	Pos   utils.Vec2
	Width float32
}

// CalculateTextBounds is only optimized when you draw text top to down.
// TODO: Create font interface and add it to mGUI
func (f *Font) CalculateTextBounds(text string, scale float32) (width, height float32, chars []CombinedCharInfo) {
	prevR := rune(-1)

	fontSize := f.GetXHeight() + 2
	height = scale * float32(fontSize)
	tmp := []rune(text)
	chars = make([]CombinedCharInfo, len(tmp))

	var maxDescend, baseline, maxWidth float32
	linesCounter := 1
	baseline = scale * float32(fontSize)
	var dx float32 = 0
	for i, r := range tmp {
		info := f.GetCharacter(r)
		if info.Width == 0 {
			log.Printf("Unknown char = %q", r)
			fmt.Println("unknown char")
			continue
		}
		if prevR >= 0 {
			kern := f.Face.Kern(prevR, r).Ceil()
			dx += float32(kern)
		}
		if r != ' ' {
			dx += float32(info.LeftBearing)
		}
		if r == '\n' {
			linesCounter++

			dx = 0
			height += float32(fontSize)
			baseline += float32(fontSize)
			prevR = rune(-1)
			continue
		}
		xPos := dx
		yPos := baseline
		if info.Descend != 0 {
			d := float32(info.Descend) * scale
			yPos += d
			if d > maxDescend {
				maxDescend = d
			}
		}
		if info.Rune == ' ' {
			chars[i] = CombinedCharInfo{
				Char:  *info,
				Pos:   utils.Vec2{X: xPos, Y: yPos},
				Width: float32(info.Width),
			}
		} else {
			chars[i] = CombinedCharInfo{
				Char:  *info,
				Pos:   utils.Vec2{X: xPos, Y: yPos},
				Width: float32(info.LeftBearing + info.Width + info.RigthBearing),
			}
		}

		//pos[i] = utils.Vec2{X: xPos, Y: yPos}
		dx += float32(info.Width) * scale
		if r != ' ' {
			dx += float32(info.RigthBearing)
		}
		prevR = r
		width = dx
		if linesCounter > 1 {
			if width > maxWidth {
				maxWidth = width
			}
		} else {
			maxWidth = width
		}
	}
	height += maxDescend
	width = maxWidth
	return
}

func (f *Font) generateBitmap(dpi float32, from, to int) *image.RGBA {
	cp := charmap.Windows1251
	var letters []rune
	for i := from; i < to; i++ {
		r := cp.DecodeByte(byte(i))
		letters = append(letters, r)
	}
	f.CharSlice = make([]*CharInfo, len(letters))
	var (
		DPI          = dpi
		startingDotX = 0
		startingDotY = 0
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
			DPI:     float64(DPI),
			Hinting: font.HintingNone,
		})

		if err != nil {
			log.Fatalf("NewFace: %v", err)
		}
	}
	f.Face = face
	defer face.Close()

	dst := image.NewRGBA(image.Rect(0, 0, AtlasWidth, AtlasWidth))
	d := font.Drawer{
		Dst:  dst,
		Src:  image.White,
		Face: face,
		Dot:  fixed.P(startingDotX, startingDotY),
	}
	fontSize := d.Face.Metrics().Height
	f.FontSize = fontSize.Ceil()
	d.Dot = fixed.P(startingDotX, d.Face.Metrics().Ascent.Floor())

	dx := startingDotX
	for i, c := range letters {
		b, a, _ := d.Face.GlyphBounds(c)
		if d.Dot.X.Ceil()+a.Ceil() >= AtlasWidth {
			dx = 0
			d.Dot.X = 0
			d.Dot.Y += d.Face.Metrics().Height
		}
		dr, mask, maskp, _, ok := d.Face.Glyph(d.Dot, c)
		if !ok {
			fmt.Println("error")
			continue
		}
		ch := CharInfo{
			Rune:         c,
			SrcX:         dr.Min.X,
			SrcY:         dr.Max.Y,
			Width:        dr.Dx(),
			Heigth:       dr.Dy(),
			Ascend:       -b.Min.Y.Floor(),
			Descend:      b.Max.Y.Floor(),
			LeftBearing:  b.Min.X.Floor(),
			RigthBearing: a.Floor() - b.Max.X.Floor(),
			TexCoords:    [2]math.Vec2{},
		}
		ch.calcTexCoords(AtlasWidth, AtlasWidth)
		f.CharSlice[i] = &ch
		//printBorder(dst, dr.Min.X, dr.Max.Y, dr.Dx(), dr.Dy(), colornames.Red)
		draw.DrawMask(d.Dst, dr, d.Src, image.Point{}, mask, maskp, draw.Over)
		d.Dot.X += fixed.I(a.Ceil() + 2)
		dx += a.Ceil()
		if c == ' ' {
			ch.Width = f.FontSize / 3
		}
		if c == '\u007f' {
			f.CharMap[CharNotFound] = &ch
		} else {
			f.CharMap[int(c)] = &ch
		}
	}
	return dst
}
func printBorder(m *image.RGBA, x, y, w, h int, clr color.Color) {

	for i := y; i >= y-h; i-- {
		m.Set(x, i, clr)
	}
	for i := x; i <= x+w; i++ {
		m.Set(i, y-h, clr)
	}
	for i := y; i >= y-h; i-- {
		m.Set(x+w, i, clr)
	}
	for i := x + w; x <= i; i-- {
		m.Set(i, y, clr)
	}

}
