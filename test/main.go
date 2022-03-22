package main

import (
	"fmt"
	"image"

	"os"
	"time"

	//	stb "neilpa.me/go-stbi"
	"github.com/nicholasblaskey/stbi"
	_ "neilpa.me/go-stbi/png"
)

func main() {
	t := time.Now()
	DecodeTextureStb("assets/images/goomba.png")
	fmt.Printf("%s",time.Now().Sub(t))
}

func changeSlice(sl []int) {
	for i, _ := range sl {
		sl[i] = 1
	}
}

func DecodeTextureStb(filename string) {
	//i, _:= stb.Load(filename)
	_, width, height, _, cleanup, err := stbi.Load(filename,true,0)
	if err != nil {
		fmt.Println(err)
		return
	}
	defer cleanup()

	w := int(width)
	h := int(height)

	// pixels := make([]byte, w*h*4)
	// bIndex := 0
	for y  := 0; y < h; y++ {
		for x := 0; x < w; x++ {
			// r, g, b, a := img.At(x, y).RGBA()
			// pixels[bIndex] = byte(r / 256)
			// bIndex++
			// pixels[bIndex] = byte(g / 256)
			// bIndex++
			// pixels[bIndex] = byte(b / 256)
			// bIndex++
			// pixels[bIndex] = byte(a / 256)
			// bIndex++
		}
	}
}

func DecodeTextureImage(filename string) {
	infile, err := os.Open(filename)
	if err != nil {
		return 
	}
	defer infile.Close()

	img, _, err := image.Decode(infile)
	if err != nil {
		return
	}


	w := img.Bounds().Max.X
	h := img.Bounds().Max.Y

	pixels := make([]byte, w*h*4)
	bIndex := 0
	for y := 0; y < h; y++ {
		for x := 0; x < w; x++ {
			r, g, b, a := img.At(x, y).RGBA()
			pixels[bIndex] = byte(r / 256)
			bIndex++
			pixels[bIndex] = byte(g / 256)
			bIndex++
			pixels[bIndex] = byte(b / 256)
			bIndex++
			pixels[bIndex] = byte(a / 256)
			bIndex++
		}
	}
}
