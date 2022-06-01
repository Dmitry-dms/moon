package main

import "github.com/Dmitry-dms/moon/pkg/ui2/fonts"

// "fmt"

// import "github.com/Dmitry-dms/moon/pkg/ui2/fonts"



func main() {

	fonts.NewFont("assets/fonts/Roboto.ttf", 30)

	

	// const asciiArt = ".++8"
	// buf := make([]byte, 0, height*(width+1))
	// for y := 0; y < height; y++ {
	// 	for x := 0; x < width; x++ {
	// 		c := asciiArt[dst.GrayAt(x, y).Y>>6]
	// 		if c != '.' {
	// 			// No-op.
	// 		} else if x == startingDotX-1 {
	// 			c = ']'
	// 		} else if y == startingDotY-1 {
	// 			c = '_'
	// 		}
	// 		buf = append(buf, c)
	// 	}
	// 	buf = append(buf, '\n')
	// }
	// os.Stdout.Write(buf)

}
