package fonts

import (
	// "fmt"

	"github.com/Dmitry-dms/moon/pkg/math"
)

type CharInfo struct {
	srcX, srcY    int
	width, heigth int
	TexCoords     [2]math.Vec2
}

func (c *CharInfo) calcTexCoords(fontWidth, fontHeight int) {
	x0 := float32(c.srcX) / float32(fontWidth)
	x1 := (float32(c.srcX) + float32(c.width)) / float32(fontWidth)
	y0 := float32(c.srcY) / float32(fontHeight)
	y1 := (float32(c.srcY) + float32(c.heigth)) / float32(fontHeight)

	c.TexCoords[0] = math.Vec2{X: x0, Y: y1}
	c.TexCoords[1] = math.Vec2{X: x1, Y: y0}
}

func (f *Font) GetCharacter(chr rune) CharInfo {
	c, ok := f.CharMap[int(chr)]
	if !ok {
		panic("nothing found")
	}
	// fmt.Println("char width = ", c.width)
	return *c
}
