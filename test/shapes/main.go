package main

import (
	"fmt"
	"golang.org/x/image/colornames"
	"image"
	"image/color"
	"image/png"
	"math"
	"os"
)

func main() {

	var startX, startY float32 = 0.0, 0.0
	var endX, endY float32 = 50.0, 50.0
	var supportX, supportY float32 = .0, 50.0

	bezue := func(t float32) (float32, float32) {
		v1 := float32(math.Pow(float64(1-t), 2))
		v2 := 2 * t * (1 - t)
		v3 := float32(math.Pow(float64(t), 2))
		return v1*startX + v2*supportX + v3*endX, v1*startY + v2*supportY + v3*endY
	}

	img := image.NewRGBA(image.Rect(0, 0, 200, 200))
	//drawPoint(img, int(startX), int(startY), colornames.Violet)
	drawPoint(img, int(endX), int(endY), colornames.Violet)
	for t := startX; t < 1.0; t += 0.05 {
		x, y := bezue(t)
		drawPoint(img, int(x), int(y), colornames.Red)
	}

	//CreateImage("bezue.png", img)
	var precise float32 = 20
	fmt.Println(1 / precise)
}
func CreateImage(filename string, img image.Image) {
	pngFile, _ := os.OpenFile(filename, os.O_CREATE|os.O_RDWR, 0664)

	defer pngFile.Close()

	encoder := png.Encoder{
		CompressionLevel: png.BestCompression,
	}
	encoder.Encode(pngFile, img)
}

func drawPoint(img *image.RGBA, x, y int, clr color.Color) {
	img.Set(x, y, clr)
}
