package fonts

import "github.com/Dmitry-dms/moon/pkg/math"

type CharInfo struct {
	srcX, srcY    int
	width, heigth int
	TexCoords     [4]math.Vec2
}

func (c *CharInfo) calcTexCoords(fontWidth, fontHeight int) {
	x0 := float32(c.srcX) / float32(fontWidth)
	x1 := (float32(c.srcX) + float32(c.width)) / float32(fontWidth)
	y0 := float32(c.srcY) / float32(fontHeight)
	y1 := (float32(c.srcY) + float32(c.heigth))  / float32(fontHeight)

	c.TexCoords[0] = math.Vec2{x1, y0}
	c.TexCoords[1] = math.Vec2{x1, y1}
	c.TexCoords[2] = math.Vec2{x0, y1}
	c.TexCoords[3] = math.Vec2{x0, y0}

	// c.TexCoords[0] = math.Vec2{x0,y0}
	// c.TexCoords[1] = math.Vec2{x0,y1}
	// c.TexCoords[2] = math.Vec2{x1,y1}
	// c.TexCoords[3] = math.Vec2{x0,y1}

}

func (f *Font) GetCharacter(chr rune) CharInfo {
	c, ok := f.CharMap[int(chr)]
	if !ok {
		panic("nothing found")
	}
	return *c
}
