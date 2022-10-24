package fonts

import (
	"github.com/Dmitry-dms/moon/pkg/math"
	"github.com/Dmitry-dms/moon/pkg/sprite_packer"
	"github.com/Dmitry-dms/moon/pkg/ui/utils"
	"golang.org/x/image/colornames"
	"image"
	"image/draw"
	"image/png"
	"os"
	"sort"

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
	Filepath string
	FontSize int

	CharMap map[int]*CharInfo

	TextureId uint32
	Texture   *gogl.Texture

	Face font.Face
}

func NewFont(filepath string, fontSize int) *Font {
	f := Font{
		Filepath: filepath,
		FontSize: fontSize,
		CharMap:  make(map[int]*CharInfo, 50),
	}

	f.generateAndUploadBitmap()
	return &f
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

func (f *Font) generateAndUploadBitmap() {
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
	sort.Slice(sortSlice, func(i, j int) bool {
		return sortSlice[i].Heigth > sortSlice[j].Heigth
	})

	dst2 := dst.SubImage(image.Rect(0, 0, siz, siz))
	dst3 := image.NewRGBA(dst2.Bounds())

	draw.Draw(dst3, dst2.Bounds(), dst2, image.ZP, draw.Src)

	pngFile, _ := os.OpenFile("fonts.png", os.O_CREATE|os.O_RDWR, 0664)

	defer pngFile.Close()

	encoder := png.Encoder{
		CompressionLevel: png.BestCompression,
	}
	//encoder.Encode(pngFile, dst3)
	//t2 := gogl.UploadRGBATextureFromMemory(dst3)
	//f.TextureId = t2.GetId()
	//f.Texture = t2
	//initWidth := 400

	sheet := sprite_packer.NewSpriteSheet(128)

	sheet.BeginGroup("font-18", func() []*sprite_packer.SpriteInfo {
		spriteInfo := make([]*sprite_packer.SpriteInfo, len(sortSlice))
		for i, info := range sortSlice {
			if info.Rune == ' ' || info.Rune == '\u00a0' {
				continue
			}
			ret := dst3.SubImage(image.Rect(info.SrcX, info.SrcY, info.SrcX+info.Width, info.SrcY-info.Heigth)).(*image.RGBA)
			pixels := sheet.GetData(ret)
			spriteInfo[i] = sheet.AddToSheet(string(info.Rune), pixels)
		}
		return spriteInfo
	})

	rr, err := sheet.GetGroup("font-18")
	if err != nil {
		panic(err)
	}

	for _, info := range rr {
		if info != nil {
			ll := []rune(info.Id)
			char := f.GetCharacter(ll[0])
			char.TexCoords = [2]math.Vec2{{info.TextCoords[0], info.TextCoords[1]},
				{info.TextCoords[2], info.TextCoords[3]}}
		}
	}

	im := sheet.Image
	t2 := gogl.UploadRGBATextureFromMemory(sheet.Image)
	f.TextureId = t2.GetId()
	f.Texture = t2
	//func(im *image.RGBA) {
	//	m := im
	//	x := 0
	//	y := 0
	//	w := im.Bounds().Dx()
	//	h := im.Bounds().Dy()
	//	for i := 0; i <= h; i++ {
	//		m.Set(x, i, colornames.Red)
	//	}
	//	for i := 0; i <= w; i++ {
	//		m.Set(i, y, colornames.Red)
	//	}
	//	//for i := y; i >= y-h; i-- {
	//	//	m.Set(x+w, i, colornames.Red)
	//	//}
	//	//for i := x + w; x <= i; i-- {
	//	//	m.Set(i, y, colornames.Red)
	//	//}
	//}(im)

	//ret2 := dst3.SubImage(image.Rect(info.SrcX, info.SrcY, info.SrcX+info.Width, info.SrcY-info.Heigth)).(*image.RGBA)
	//nh := image.NewRGBA(image.Rect(0, 0, ret.Bounds().Dx(), ret.Bounds().Dy()))

	//sheet.AddToSheet(pixels2)
	//pixels := sheet.GetData(ret)

	//for y := 0; y < len(pixels); y++ {
	//	for x := 0; x < len(pixels[0]); x++ {
	//		q := pixels[y]
	//		if q == nil {
	//			continue
	//		}
	//		p := pixels[y][x]
	//		if p == nil {
	//			continue
	//		}
	//		original, ok := color.RGBAModel.Convert(p).(color.RGBA)
	//		if ok {
	//			nh.Set(x, y, original)
	//		}
	//	}
	//}
	//srcX := 0
	//srcY := 0
	//nsl := []*CharInfo{}
	//nsl = append(nsl, sortSlice[0])
	//nsl = append(nsl, sortSlice[0], sortSlice[1], sortSlice[2], sortSlice[3], sortSlice[4])
	//for _, info := range sortSlice[0:20] {
	//
	//	//info := sortSlice[10]
	//	re := dst3.SubImage(image.Rect(info.SrcX, info.SrcY, info.SrcX+info.Width, info.SrcY-info.Heigth)).(*image.RGBA)
	//	b := re.Bounds()
	//	var pixels [][]color.Color
	//	//put pixels into two three two dimensional array
	//	for i := b.Min.Y; i < b.Max.Y; i++ {
	//		var y []color.Color
	//		for j := b.Min.X; j < b.Max.X; j++ {
	//			y = append(y, re.RGBAAt(j, i))
	//		}
	//		pixels = append(pixels, y)
	//	}
	//	if srcX+len(pixels) >= initWidth {
	//		srcX = 0
	//		srcY -= fontSize.Ceil()
	//	}
	//	ypos := srcY
	//	xpos := srcX
	//	h := checkForFreeHeight(im, xpos, -ypos, info.Width, info.Heigth)
	//	ypos += h - 1
	//	w := checkForFreeWidth(im, xpos, -ypos, info.Width, info.Heigth)
	//	xpos -= w - 1
	//	//fmt.Println(h)
	//	fmt.Printf("char = %q, h = %d  w = %d\n", info.Rune, h, w)
	//	//if h != 0 {
	//
	//	//}
	//	//if srcY < 0 {
	//	//	srcY = 0
	//	//}
	//	//fmt.Println(srcY)
	//	for y := 0; y < len(pixels); y++ {
	//		for x := 0; x < len(pixels[0]); x++ {
	//			q := pixels[y]
	//			if q == nil {
	//				continue
	//			}
	//			p := pixels[y][x]
	//			if p == nil {
	//				continue
	//			}
	//			original, ok := color.RGBAModel.Convert(p).(color.RGBA)
	//			if ok {
	//				im.Set(x+xpos, y-ypos, original)
	//			}
	//		}
	//	}
	//	//s := im.Bounds().Size()
	//
	//	//}
	//
	//	srcX += info.Width
	//	//srcX += 2
	//}

	encoder.Encode(pngFile, im)
}
func checkForFreeHeight(img *image.RGBA, x, y, w, h int) (height int) {
all:
	for i := y - 1; i > 0; i-- {
		for j := x; j < x+w; j++ {
			r, g, b, _ := img.At(j, i).RGBA()
			//img.Set(j, i, colornames.Red)
			//fmt.Println(r, g, b, a)
			if r == 0 && g == 0 && b == 0 {

			} else {
				break all
			}
		}
		height += 1
	}
	return
}
func checkForFreeWidth(img *image.RGBA, x, y, w, h int) (width int) {
all:
	for i := x; i > 0; i-- {
		for j := y; j < y+h; j++ {
			r, g, b, _ := img.At(i, j).RGBA()
			//img.Set(i, j, colornames.Red)
			if r == 0 && g == 0 && b == 0 {

			} else {
				break all
			}
		}
		width += 1
	}
	return
}

func printBorder(m *image.RGBA, x, y, w, h int, desc fixed.Point26_6, a fixed.Int26_6, sy int) {
	for i := y; i >= y-h; i-- {
		m.Set(x, i, colornames.Red)
	}
	for i := x; i <= x+w; i++ {
		m.Set(i, y-h, colornames.Red)
	}
	for i := y; i >= y-h; i-- {
		m.Set(x+w, i, colornames.Red)
	}
	for i := x + w; x <= i; i-- {
		m.Set(i, y, colornames.Red)
	}

	for i := desc.X.Ceil(); i <= desc.X.Ceil()+a.Ceil(); i++ {
		m.Set(i, desc.Y.Ceil(), colornames.Blue)
	}

	for i := desc.X.Ceil(); i <= desc.X.Ceil()+a.Ceil(); i++ {
		m.Set(i, sy, colornames.Violet)
	}

}
