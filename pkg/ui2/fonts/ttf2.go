package fonts

import (
	"fmt"
	"image"
	"image/color"
	"image/draw"
	"io/ioutil"

	"github.com/golang/freetype"
	"github.com/golang/freetype/truetype"
	"golang.org/x/image/math/fixed"
	"golang.org/x/text/encoding/charmap"
)

func ParseTTF(options Options, path string, lowChar int, highChar int) (*image.NRGBA, error) {
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
	ctx.SetFontSize(float64(options.FontSize))

	return renderSheet(options.FontSize, lowChar, highChar, ctx, font)
}

func renderSheet(fontSize, lowChar, highChar int, ctx *freetype.Context, font *truetype.Font) (*image.NRGBA, error) {

	//fake image
	bg := image.NewUniform(color.NRGBA{0x00, 0x00, 0xFF, 0x00})

	destHeightPixels := float32(ctx.PointToFixed(float64(fontSize)*1.4).Ceil() + 1)
	estimatedWidth := 512

	width := 0
	height := destHeightPixels
	lineHeight := destHeightPixels
	x := 0
	y := lineHeight * 1.4

	cp := charmap.Windows1252
	scale := ctx.PointToFixed(float64(fontSize))

	for c := lowChar; c <= highChar; c++ {
		r := cp.DecodeByte(byte(c))
		index := font.Index(r)
		hMetric := font.HMetric(scale, index)

		width = max(x+hMetric.AdvanceWidth.Ceil(), width)

		x += hMetric.AdvanceWidth.Ceil()
		if x > estimatedWidth {
			x = 0
			y += lineHeight * 1.4
			height += destHeightPixels * 1.4
		}
		// width += hMetric.AdvanceWidth.Ceil()
	}
	// fmt.Println(width, height)

	//real image
	dest := image.NewNRGBA(image.Rect(0, 0, width, int(height)))
	draw.Draw(dest, dest.Bounds(), bg, image.ZP, draw.Src)

	//
	src := image.White
	ctx.SetSrc(src)
	baseline := scale.Ceil()

	// Assume no glyph is landscape
	bufferWidth := int(float64(dest.Bounds().Inset(1).Dy()))
	// fmt.Println("buffer width ", bufferWidth)
	buffer := image.NewNRGBA(image.Rect(0, 0, bufferWidth, bufferWidth))
	// fmt.Println("buuffer ", buffer.Bounds())
	ctx.SetDst(buffer)
	ctx.SetClip(buffer.Bounds())

	xOffset := 1 //отступ по горизонтали в пикселях

	x = 0
	y = 0
	for c := lowChar; c <= highChar; c++ {
		// fmt.Println(xOffset)
		r := cp.DecodeByte(byte(c))
		index := font.Index(r)
		hMetric := font.HMetric(scale, index)

		leftSideAdjustment := fixed.I(0)
		if hMetric.LeftSideBearing < 0 {
			leftSideAdjustment = -hMetric.LeftSideBearing
		}

		// Fill with background
		draw.Draw(buffer, buffer.Bounds(), bg, image.Point{}, draw.Src)

		// Draw glyph
		advance, _:= ctx.DrawString(string(r), fixed.Point26_6{X: leftSideAdjustment, Y: fixed.I(baseline)})

		advance.X += leftSideAdjustment

		draw.Draw(dest, image.Rect(xOffset, 1, xOffset+buffer.Bounds().Dx(), dest.Bounds().Dy()-1), buffer, image.ZP, draw.Src)

		// xOffset += advance.X.Ceil()
		x += advance.X.Ceil()

		//отвечает за вертикальную границу между символами
		// draw.Draw(dest, image.Rect(xOffset, 0, xOffset+1, dest.Bounds().Dy()), Border, image.ZP, draw.Src)

		// xOffset += 1.0
	}

	return nil, nil
}
func render(fontSize, lowChar, highChar int, ctx *freetype.Context, font *truetype.Font, dest draw.Image, bg image.Image) (int, error) {
	src := image.White
	ctx.SetSrc(src)

	cp := charmap.Windows1252
	scale := ctx.PointToFixed(float64(fontSize))
	baseline := scale.Ceil()
	// fmt.Println(baseline)

	// Assume no glyph is landscape
	bufferWidth := int(float64(dest.Bounds().Inset(1).Dy()))
	// fmt.Println("buffer width ", bufferWidth)
	buffer := image.NewNRGBA(image.Rect(0, 0, bufferWidth, bufferWidth))
	// fmt.Println("buuffer ", buffer.Bounds())
	ctx.SetDst(buffer)
	ctx.SetClip(buffer.Bounds())

	xOffset := 1 //отступ по горизонтали в пикселях

	for c := lowChar; c <= highChar; c++ {
		// fmt.Println(xOffset)
		r := cp.DecodeByte(byte(c))
		index := font.Index(r)
		hMetric := font.HMetric(scale, index)

		leftSideAdjustment := fixed.I(0)
		if hMetric.LeftSideBearing < 0 {
			leftSideAdjustment = -hMetric.LeftSideBearing
		}

		// Fill with background
		draw.Draw(buffer, buffer.Bounds(), bg, image.Point{}, draw.Src)

		// Draw glyph
		advance, err := ctx.DrawString(string(r), fixed.Point26_6{X: leftSideAdjustment, Y: fixed.I(baseline)})
		if err != nil {
			return 0, err
		}
		advance.X += leftSideAdjustment

		draw.Draw(dest, image.Rect(xOffset, 1, xOffset+buffer.Bounds().Dx(), dest.Bounds().Dy()-1), buffer, image.ZP, draw.Src)

		xOffset += advance.X.Ceil()

		//отвечает за вертикальную границу между символами
		// draw.Draw(dest, image.Rect(xOffset, 0, xOffset+1, dest.Bounds().Dy()), Border, image.ZP, draw.Src)

		// xOffset += 1.0
	}

	return xOffset, nil
}

func max[T int](x1, x2 T) T {
	if x1 > x2 {
		return x1
	} else {
		return x2
	}
}
