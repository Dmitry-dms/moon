package fonts

import (
	"fmt"
	"image"
	"image/color"
	"image/draw"

	"github.com/golang/freetype"
	"github.com/golang/freetype/truetype"
	"golang.org/x/image/math/fixed"
	"golang.org/x/text/encoding/charmap"
)

func tigrFromTTF(options Options, ttfBytes []byte, lowChar int, highChar int) (*image.NRGBA, error) {
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

	image, err := renderTTFSheet(options.FontSize, lowChar, highChar, ctx, font)
	if err != nil {
		return nil, fmt.Errorf("failed to render TTF: %w", err)
	}

	return image, nil
}

func renderTTFSheet(fontSize, lowChar, highChar int, ctx *freetype.Context, font *truetype.Font) (*image.NRGBA, error) {
	bg := image.NewUniform(color.NRGBA{0x00, 0x00, 0xFF, 0x00})

	destHeightPixels := ctx.PointToFixed(float64(fontSize)*1.5).Ceil() + 2
	// fmt.Println("DEST ", destHeightPixels)
	
	dest := image.NewNRGBA(image.Rect(0, 0, destHeightPixels, destHeightPixels))

	actualWidth, err := renderTTFChars(fontSize, lowChar, highChar, ctx, font, dest, bg)
	if err != nil {
		return nil, err
	}

	dest = image.NewNRGBA(image.Rect(0, 0, actualWidth, destHeightPixels))
	// fmt.Println("act ", actualWidth)

	draw.Draw(dest, dest.Bounds(), bg, image.ZP, draw.Src)

	_, err = renderTTFChars(fontSize, lowChar, highChar, ctx, font, dest, bg)
	if err != nil {
		return nil, err
	}

	return dest, nil
}

func renderTTFChars(fontSize, lowChar, highChar int, ctx *freetype.Context, font *truetype.Font, dest draw.Image, bg image.Image) (int, error) {
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
		draw.Draw(dest, image.Rect(xOffset, 0, xOffset+1, dest.Bounds().Dy()), Border, image.ZP, draw.Src)

		xOffset += 1.0
	}

	return xOffset, nil
}

func getPointSizeFromX(font *truetype.Font, fontSize int) (int, error) {
	ctx := freetype.NewContext()
	ctx.SetFont(font)
	ctx.SetDPI(72.0)
	ctx.SetFontSize(float64(fontSize))

	img, err := renderTTFSheet(fontSize, 'X', 'X', ctx, font)
	if err != nil {
		return 0, err
	}

	bounds := contentBounds(img)
	// fmt.Println(bounds)
	actual := float64(bounds.Dy())
	expected := float64(fontSize)
	factor := expected / actual
	return int(expected * factor), nil
}

func contentBounds(img *image.NRGBA) image.Rectangle {
	minNonTransparentRow := img.Bounds().Max.Y
	maxNonTransparentRow := img.Bounds().Min.Y

rows:
	for y := img.Bounds().Min.Y; y < img.Bounds().Max.Y; y++ {
	cols:
		for x := img.Bounds().Min.X; x < img.Bounds().Max.X; x++ {
			color := img.At(x, y)
			if color == Border.C {
				continue cols
			}
			_, _, _, a := color.RGBA()
			if a != 0 {
				// non-transparent
				if y < minNonTransparentRow {
					minNonTransparentRow = y
				}
				if y > maxNonTransparentRow {
					maxNonTransparentRow = y
				}
				continue rows
			}
		}
	}

	return image.Rect(img.Bounds().Min.X, minNonTransparentRow, img.Bounds().Max.X, maxNonTransparentRow)
}

func shrinkToFit(img *image.NRGBA) *image.NRGBA {
	bounds := contentBounds(img)

	if bounds.Min.Y > 0 {
		bounds.Min.Y--
	}
	if bounds.Min.Y > 0 {
		bounds.Min.Y--
	}

	if bounds.Max.Y < img.Bounds().Dy() {
		bounds.Max.Y++
	}
	if bounds.Max.Y < img.Bounds().Dy() {
		bounds.Max.Y++
	}
	return (img.SubImage(bounds)).(*image.NRGBA)
}

func frame(dest draw.Image, border image.Image) {
	minX := dest.Bounds().Min.X
	minY := dest.Bounds().Min.Y
	maxX := minX + dest.Bounds().Dx()
	maxY := minY + dest.Bounds().Dy()

	draw.Draw(dest, image.Rect(minX, minY, maxX, minY+1), border, image.ZP, draw.Src)
	draw.Draw(dest, image.Rect(maxX-1, minY, maxX, maxY), border, image.ZP, draw.Src)
	draw.Draw(dest, image.Rect(minX, maxY-1, maxX, maxY), border, image.ZP, draw.Src)
	draw.Draw(dest, image.Rect(minX, minY, 1, maxY), border, image.ZP, draw.Src)
}