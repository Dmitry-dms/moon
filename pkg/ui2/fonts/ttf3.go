package fonts

import (
	"fmt"
	"image"
	"image/color"
	"image/draw"

	// "image/png"
	"io/ioutil"
	// "os"

	// "github.com/go-gl/gl/v4.2-core/gl"
	"github.com/golang/freetype"
	"github.com/golang/freetype/truetype"
	"golang.org/x/image/math/fixed"
	"golang.org/x/text/encoding/charmap"
)

func (f *Font) FromTTF(options Options, path string, lowChar int, highChar int) (*image.NRGBA, error) {
	ttfBytes, err := ioutil.ReadFile(path)
	if err != nil {
		panic(err)
	}
	font, err := freetype.ParseFont(ttfBytes)
	if err != nil {
		return nil, fmt.Errorf("failed to parse TTF: %w", err)
	}

	ctx := freetype.NewContext()
	ctx.SetFont(font)
	ctx.SetDPI(float64(options.DPI))

	if options.Measure {
		options.FontSize, err = getPointSizeFromX(font, options.FontSize)
		if err != nil {
			return nil, fmt.Errorf("failed to measure TTF font: %w", err)
		}
	}
	ctx.SetFontSize(float64(options.FontSize))

	image, err := f.renderTSheet(options.FontSize, lowChar, highChar, ctx, font)
	if err != nil {
		return nil, fmt.Errorf("failed to render TTF: %w", err)
	}

	// Подстраивает высоту
	// image = shrinkToFit(image)
	// f.uploadTexture(image)

	// pngFile, _ := os.OpenFile("target.png", os.O_CREATE|os.O_RDWR, 0664)

	// defer pngFile.Close()

	// encoder := png.Encoder{
	// 	CompressionLevel: png.BestCompression,
	// }
	// err = encoder.Encode(pngFile, image)

	return image, nil
}



func (f *Font) renderTSheet(fontSize, lowChar, highChar int, ctx *freetype.Context, font *truetype.Font) (*image.NRGBA, error) {
	bg := image.NewUniform(color.NRGBA{0x00, 0x00, 0xFF, 0x00})

	destHeightPixels := ctx.PointToFixed(float64(fontSize)*1.4).Ceil() + 1

	dest := image.NewNRGBA(image.Rect(0, 0, destHeightPixels, destHeightPixels))

	actualWidth, actualHeight, err := f.checkAtlas(fontSize, lowChar, highChar, ctx, font, dest, bg)
	if err != nil {
		return nil, err
	}

	dest = image.NewNRGBA(image.Rect(0, 0, actualWidth, actualHeight))

	draw.Draw(dest, dest.Bounds(), bg, image.ZP, draw.Src)

	f.renderTChars(fontSize, lowChar, highChar, ctx, font, dest, bg)
	if err != nil {
		return nil, err
	}

	return dest, nil
}
func (f *Font) checkAtlas(fontSize, lowChar, highChar int, ctx *freetype.Context, font *truetype.Font, dest draw.Image, bg image.Image) (int, int, error) {
	src := image.White
	ctx.SetSrc(src)

	cp := charmap.Windows1252
	scale := ctx.PointToFixed(float64(fontSize))
	baseline := scale.Ceil()

	// Assume no glyph is landscape
	bufferWidth := int(float64(dest.Bounds().Inset(1).Dy()))

	buffer := image.NewNRGBA(image.Rect(0, 0, bufferWidth, bufferWidth))

	ctx.SetDst(buffer)
	ctx.SetClip(buffer.Bounds())

	width := 0

	x := 0
	var y float32 = 0.0

	height := float32(dest.Bounds().Max.Y)
	adj := 1024

	prev := 0
	for c := lowChar; c <= highChar; c++ {
		// fmt.Println(xOffset)
		r := cp.DecodeByte(byte(c))
		index := font.Index(r)

		hMetric := font.HMetric(scale, index)
		yMetric := font.VMetric(scale, index)

		leftSideAdjustment := fixed.I(0)
		if hMetric.LeftSideBearing < 0 {
			leftSideAdjustment = -hMetric.LeftSideBearing
		}

		topSide := fixed.I(0)
		if yMetric.TopSideBearing > 0 {
			topSide = yMetric.TopSideBearing
		}

		// fmt.Println(hMetric.AdvanceWidth.Ceil())
		// Fill with background
		draw.Draw(buffer, buffer.Bounds(), bg, image.Point{}, draw.Src)

		// Draw glyph
		advance, _ := ctx.DrawString(string(r), fixed.Point26_6{X: leftSideAdjustment, Y: fixed.I(baseline)})

		advance.X += leftSideAdjustment
		advance.Y += topSide

		if hMetric.AdvanceWidth.Ceil() < prev && prev != 0 {
			if (prev - hMetric.AdvanceWidth.Ceil()) < 10 {
				x += (prev - hMetric.AdvanceWidth.Ceil())
			}
		}

		if x+advance.X.Ceil() > adj {
			height += float32(advance.Y.Ceil()) * 1.4
			width = x
			x = 0
			y += float32(advance.Y.Ceil()) * 1.4
		}

		ch := &CharInfo{
			srcX:   x,
			srcY:   int(y),
			width:  advance.X.Ceil(),
			heigth: advance.Y.Ceil(),
		}

		f.CharMap[c] = ch
		// draw.Draw(dest, image.Rect(x, y, x+buffer.Bounds().Dx(), y+buffer.Bounds().Dy()), buffer, image.ZP, draw.Src)

		x += hMetric.AdvanceWidth.Ceil()
		// draw.Draw(dest, image.Rect(x, 0, x+1, dest.Bounds().Dy()), Border, image.ZP, draw.Src)
		x += 1
		// draw.Draw(dest, image.Rect(0, dest.Bounds().Dy()+1, dest.Bounds().Dx(), dest.Bounds().Dy()+1), Border, image.ZP, draw.Src)

		// width += x
		prev = advance.X.Ceil()

	}
	height += float32(f.FontSize) * 1.4

	return width, int(height), nil
}

func (f *Font) renderTChars(fontSize, lowChar, highChar int, ctx *freetype.Context, font *truetype.Font, dest draw.Image, bg image.Image) error {
	src := image.White
	ctx.SetSrc(src)

	cp := charmap.Windows1252
	scale := ctx.PointToFixed(float64(fontSize))
	baseline := scale.Ceil()

	// Assume no glyph is landscape
	bufferWidth := int(float64(dest.Bounds().Inset(1).Dy()))

	buffer := image.NewNRGBA(image.Rect(0, 0, bufferWidth, bufferWidth))

	ctx.SetDst(buffer)
	ctx.SetClip(buffer.Bounds())

	// width := 0

	// x := 0
	// y := 0

	// height := dest.Bounds().Max.Y
	// adj := 512

	for c := lowChar; c <= highChar; c++ {
		// fmt.Println(xOffset)
		r := cp.DecodeByte(byte(c))
		index := font.Index(r)
		hMetric := font.HMetric(scale, index)

		leftSideAdjustment := fixed.I(0)
		if hMetric.LeftSideBearing < 0 {
			leftSideAdjustment = -hMetric.LeftSideBearing
		}

		// // Fill with background
		draw.Draw(buffer, buffer.Bounds(), bg, image.Point{}, draw.Src)

		// Draw glyph
		advance, _ := ctx.DrawString(string(r), fixed.Point26_6{X: leftSideAdjustment, Y: fixed.I(baseline)})

		advance.X += leftSideAdjustment

		charInfo := f.CharMap[c]
		charInfo.calcTexCoords(dest.Bounds().Dx(), dest.Bounds().Dy())
		// fmt.Println(charInfo.TexCoords)

		draw.Draw(dest, image.Rect(charInfo.srcX, charInfo.srcY, charInfo.srcX+charInfo.width, charInfo.srcY+charInfo.heigth), buffer, image.ZP, draw.Src)

		draw.Draw(dest, image.Rect(charInfo.srcX+charInfo.width, charInfo.srcY, charInfo.srcX+charInfo.width+1, charInfo.srcY+charInfo.heigth), Border, image.ZP, draw.Src)
		draw.Draw(dest, image.Rect(charInfo.srcX, charInfo.srcY, charInfo.srcX+charInfo.width, charInfo.srcY+1), Border, image.ZP, draw.Src)
		draw.Draw(dest, image.Rect(charInfo.srcX, charInfo.srcY+charInfo.heigth, charInfo.srcX+charInfo.width, charInfo.srcY+charInfo.heigth+1), Border, image.ZP, draw.Src)
		// draw.Draw(dest, image.Rect(100, 10, 900, 150), Border, image.ZP, draw.Src)
		// fmt.Println(charInfo.srcX, charInfo.srcY+charInfo.heigth+1, charInfo.srcX+charInfo.width, charInfo.srcY+charInfo.heigth+1)
		// x += advance.X.Ceil()
		// // width += x

		// if x > adj {
		// 	height += advance.Y.Ceil() + 10
		// 	width = x
		// 	x = 0
		// 	y += advance.Y.Ceil() + 10
		// }

	}

	return nil
}
