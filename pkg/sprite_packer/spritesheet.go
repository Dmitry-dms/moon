package sprite_packer

import (
	"golang.org/x/image/colornames"
	"image"
	"image/color"
	"image/draw"
	"math"
)

type SpriteSheet struct {
	Width             int `json:"atlas_width"`
	Height            int `json:"atlas_height"`
	Filename          string
	Group             map[string][]*SpriteInfo `json:"groups"`
	Image             *image.RGBA
	srcX, srcY        int
	currentHeight     int
	wasResized        bool
	prevWidth         int
	prevHeight        int
	currentGroup      string
	currentSeparators []separator
}

func NewSpriteSheet(initWidth int) *SpriteSheet {
	s := SpriteSheet{
		Width:    initWidth,
		Height:   initWidth,
		Filename: "",
		Group:    make(map[string][]*SpriteInfo, 5),
	}
	im := image.NewRGBA(image.Rect(0, 0, initWidth, initWidth))
	s.Image = im
	return &s
}
func (s *SpriteSheet) GetGroup(id string) ([]*SpriteInfo, bool) {
	g, ok := s.Group[id]
	if !ok {
		s.Group[id] = []*SpriteInfo{}
		return s.Group[id], false
	}
	return g, true
}

func (s *SpriteSheet) AddSprite(groupId, spriteId string, data image.Image) *SpriteInfo {
	var spr *SpriteInfo
	s.BeginGroup(groupId, func() []*SpriteInfo {
		d := s.GetData(data)
		info := s.AddToSheet(spriteId, d)
		spr = info
		return []*SpriteInfo{info}
	})
	return spr
}

func (s *SpriteSheet) increaseImage() int {
	s.prevWidth = s.Width
	s.prevHeight = s.Height
	newW := s.Width * 2
	im := image.NewRGBA(image.Rect(0, 0, newW, newW))
	draw.Draw(im, s.Image.Rect, s.Image, image.Point{}, draw.Src)
	s.Image = im
	s.wasResized = true
	return newW
}
func printBorder(m *image.RGBA, x, y, w int) {

	//for i := y; i >= y-h; i-- {
	//	m.Set(x, i, colornames.Red)
	//}
	//for i := x; i <= x+w; i++ {
	//	m.Set(i, y-h, colornames.Red)
	//}
	//for i := y; i >= y-h; i-- {
	//	m.Set(x+w, i, colornames.Red)
	//}
	for i := x + w; x <= i; i-- {
		m.Set(i, y, colornames.Violet)
	}

}

func (s *SpriteSheet) findEmptySpace(srcX, srcY, width int) int {
	if srcY == 0 {
		return 0
	}

	rayStep := int(math.Ceil(float64(width / 10)))
	if rayStep == 0 {
		rayStep = 1
	}
	tmpH := 100000
	for j := srcX; j < srcX+width; j += rayStep {
		inner := 0
		for i := srcY; i < 0; i++ {
			if s.Image.RGBAAt(j, -i).R != 0 {
				//if s.Image.RGBAAt(j, -i) == colornames.Violet {
				inner--
				break
			} else {
				//if c == des {
				//s.Image.Set(j, -i, colornames.Blue)
				//}
				inner++
			}
		}
		if inner < tmpH {
			tmpH = inner
		}
	}
	//if c == des {
	//	fmt.Println(tmpH, -srcY)
	//	fmt.Printf("%q, h = %d, y = %d \n", c, tmpH, -srcY)
	//}

	return tmpH
}

func (s *SpriteSheet) BeginGroup(id string, f func() []*SpriteInfo) {
	s.currentGroup = id
	s.currentSeparators = []separator{}

	gr, _ := s.GetGroup(id)
	updgr := append(gr, f()...)
	s.Group[id] = updgr
	s.cleanSeparators()
	s.currentGroup = ""
	if s.wasResized {
		for _, infos := range s.Group {
			for _, info := range infos {
				if info != nil {
					info.calcTexCoords(s.Width, s.Height)
				}
			}
		}
		s.wasResized = false
	}
	s.currentSeparators = []separator{}
}

type separator struct {
	x, y, w int
}

func (s *SpriteSheet) cleanSeparators() {
	for _, i2 := range s.currentSeparators {
		for i := i2.x + i2.w; i2.x <= i; i-- {
			s.Image.Set(i, i2.y, color.Transparent)
		}
	}
}

func (s *SpriteSheet) AddToSheet(id string, pixels [][]color.Color) *SpriteInfo {
	height := len(pixels)
	width := len(pixels[0])
	//fmt.Printf("%q x = %d, y = %d \n", char, s.srcX, s.srcY)
	//fmt.Println(width+s.srcX, s.Width)
	if width+s.srcX > s.Width {
		s.srcY -= s.currentHeight
		if -s.srcY+height > s.prevHeight {
			s.srcX = 0
		} else {
			s.srcX = s.prevWidth
		}
	}
	if -s.srcY+height > s.Height {
		h := s.increaseImage()
		s.srcX = s.prevWidth
		s.srcY = 0
		s.Width = h
		s.Height = h
	}
	g1 := s.findEmptySpace(s.srcX, s.srcY, width)

	if height > s.currentHeight {
		s.currentHeight = height + 1
	}
	ypos := s.srcY

	if g1 > 0 {
		ypos += g1
	}
	for y := 0; y < len(pixels); y++ {
		for x := 0; x < len(pixels[0]); x++ {
			p := pixels[y][x]
			s.Image.Set(x+s.srcX, y-ypos, p)
		}
	}
	printBorder(s.Image, s.srcX, -ypos+height, width-1)
	s.currentSeparators = append(s.currentSeparators, separator{
		x: s.srcX,
		y: -ypos + height,
		w: width,
	})
	srcInfo := SpriteInfo{
		Id:     id,
		SrcX:   s.srcX,
		SrcY:   -ypos + height,
		Width:  width,
		Height: height,
	}
	srcInfo.calcTexCoords(s.Width, s.Height)
	s.srcX += width + 1
	return &srcInfo
}

func (s *SpriteSheet) GetData(data image.Image) [][]color.Color {
	b := data.Bounds()
	pixels := make([][]color.Color, b.Dy())
	for i := range pixels {
		pixels[i] = make([]color.Color, b.Dx())
	}
	ik := 0 // Needs of using internal counters because [data]
	// is probably subImage that points to the original image
	for i := b.Min.Y; i < b.Max.Y; i++ {
		jk := 0
		for j := b.Min.X; j < b.Max.X; j++ {
			pixels[ik][jk] = data.At(j, i)
			jk++
		}
		ik++
	}
	return pixels
}

type SpriteInfo struct {
	Id         string     `json:"sprite_id"`
	SrcX       int        `json:"sprite_src_x"`
	SrcY       int        `json:"sprite_src_y"`
	Width      int        `json:"sprite_w"`
	Height     int        `json:"sprite_h"`
	TextCoords [4]float32 `json:"sprite_tex_coords"` //uv0 uv0 uv1 uv1
}

func (c *SpriteInfo) calcTexCoords(atlasWidth, atlasHeight int) {
	x0 := float32(c.SrcX) / float32(atlasWidth)
	x1 := (float32(c.SrcX) + float32(c.Width)) / float32(atlasWidth)
	y0 := float32(c.SrcY) / float32(atlasHeight)
	y1 := (float32(c.SrcY) - float32(c.Height)) / float32(atlasHeight)

	c.TextCoords = [4]float32{x0, y0, x1, y1}
}
