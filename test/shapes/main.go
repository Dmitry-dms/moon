package main

import (
	"fmt"
	"golang.org/x/image/colornames"
	"image"
	"image/color"
	"image/png"
	"math"
	"os"
	"strings"
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
	//str := "Шла Саша по шоссе и сосала сушку"
	//fmt.Println(strings.LastIndex(str, "С"))
	//
	//fmt.Println(strings.LastIndexFunc("go gopher", func(r rune) bool {
	//	if r == 'h' {
	//		return true
	//	} else {
	//		return false
	//	}
	//}))
	//str := "gopher go"
	//fmt.Println(last(str[:7], 'g'))
	//fmt.Println(strings.LastIndex(str[:7], "g"))
	//fmt.Println(strings.Index(str[3:], "r"))
	//fmt.Println(RuneIndex(str, 'r', 3))
	fmt.Println(wrap("the quick brown fox jumps over the lazy dog", 10))
	//fmt.Println(wrap("the quick brown fox", 10))
	//str := "the quick brown fox"
	//
	//fmt.Println(strings.LastIndex(str[:17], "f"))
	//fmt.Println(strings.Index(str[10:], "q"))
}
func RuneIndex(s string, c rune, fromIndex int) int {
	r := []rune(s)
	ind := 0
	for i := fromIndex; i <= len(r)-1; i++ {
		if r[i] != c {
			ind++
		} else {
			break
		}
	}
	return ind + fromIndex
}
func LastRuneIndex(s string, c rune) int {
	r := []rune(s)
	ind := 0
	for i := len(r) - 1; i >= 0; i-- {
		if r[i] != c {
			ind++
		} else {
			ind++
			break
		}
	}
	return len(r) - ind
}
func wrap(msg string, wrapLength int) string {
	str := []rune(msg)
	inputLineLength := len(str)
	offset := 0
	sb := strings.Builder{}
	sb.Grow(len(str))
	for inputLineLength-offset > wrapLength {
		if str[offset] == ' ' {
			offset++
			continue
		}
		spaceToWrapAt := LastRuneIndex(string(str[:wrapLength+offset+1]), ' ')
		fmt.Println(spaceToWrapAt)
		if spaceToWrapAt >= offset {
			sb.Write([]byte(string(str[offset:spaceToWrapAt])))
			sb.Write([]byte(string('\n')))
			offset = spaceToWrapAt + 1
		} else {
			spaceToWrapAt = RuneIndex(string(str), ' ', wrapLength+offset)
			if spaceToWrapAt >= 0 {
				sb.Write([]byte(string(str[offset:spaceToWrapAt])))
				sb.Write([]byte(string('\n')))
				offset = spaceToWrapAt + 1
			} else {
				sb.Write([]byte(string(str[offset:])))
				offset = inputLineLength
			}
		}
	}
	sb.Write([]byte(string(str[offset:])))
	return sb.String()
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
