package draw

import (
	// "fmt"
	"log"
	"math"

	"github.com/Dmitry-dms/moon/pkg/fonts"
	"github.com/Dmitry-dms/moon/pkg/ui/utils"
	"github.com/Dmitry-dms/moon/pkg/ui/widgets"
)

type CmdBuffer struct {
	commands []Command

	displaySize *utils.Vec2

	Inf []Info
	ofs int

	Vertices []float32
	Indeces  []int32

	VertCount int
	lastIndc  int
	lastElems int

	InnerWindowSpace [4]float32
}

type Info struct {
	Elems       int
	IndexOffset int
	TexId       uint32
	Type        string
	ClipRect    [4]float32
}

func NewBuffer(size *utils.Vec2) *CmdBuffer {
	return &CmdBuffer{
		commands:    []Command{},
		Vertices:    []float32{},
		Indeces:     []int32{},
		displaySize: size,
		VertCount:   0,
	}
}

func (c *CmdBuffer) Clear() {
	c.commands = []Command{}
	c.Vertices = []float32{}
	c.Indeces = []int32{}
	c.VertCount = 0
	c.lastIndc = 0
	c.ofs = 0
	c.Inf = []Info{}
	c.lastElems = 0
}

func (c *CmdBuffer) SeparateBuffer(texId uint32, clipRect [4]float32) {
	cmd := Command{
		Type: SeparateBuffer,
		sb: &separate_buff{
			texid:    texId,
			clipRect: clipRect,
		},
	}
	c.AddCommand(cmd)
}

func (c *CmdBuffer) CreateButtonT(btn widgets.TextButton, font fonts.Font) {
	x := btn.Button.BoundingBox[0]
	y := btn.Button.BoundingBox[1]
	// w := btn.Button.BoundingBox
	// y := btn.Button.BoundingBox
	c.CreateRect(x, y, btn.Button.Width(), btn.Button.Height(), 0, StraightCorners, 0, btn.Button.CurrentColor)
	c.CreateText(*btn.Text, font)
}

func (c *CmdBuffer) CreateText(txt widgets.Text, font fonts.Font) {
	tcmd := Text_command{
		X:    txt.BoundingBox[0],
		Y:    txt.BoundingBox[1],
		Clr:  txt.CurrentColor,
		Text: txt.Message,
		Font: font,
		Size: txt.Size,
	}

	cmd := Command{
		Text:     &tcmd,
		Type:     Text,
		WidgetId: txt.Id,
	}
	c.AddCommand(cmd)
	c.SeparateBuffer(font.TextureId, c.InnerWindowSpace)
}

func (c *CmdBuffer) CreateWindow(wnd Window_command) {
	c.CreateRect(wnd.X, wnd.Y, wnd.W, wnd.H, 10, AllRounded, 0, wnd.Clr)
	toolbar := wnd.Toolbar
	c.CreateRect(toolbar.X, toolbar.Y, toolbar.W, toolbar.H, 10, TopRect, 0, toolbar.Clr)
	if wnd.Scrollbar.H != 0 {
		scrl := wnd.Scrollbar
		c.CreateRect(scrl.X, scrl.Y, scrl.W, scrl.H, scrl.Radius, AllRounded, 0, scrl.ScrollClr)
		c.CreateRect(scrl.Xb, scrl.Yb, scrl.Wb, scrl.Hb, scrl.Radius, AllRounded, 0, scrl.BtnClr)
	}

	c.SeparateBuffer(0, [4]float32{wnd.X, wnd.Y, wnd.W, wnd.H})
}

func (c *CmdBuffer) CreateRect(x, y, w, h float32, radius int, shape RoundedRectShape, texId uint32, clr [4]float32) {
	cmd := Command{
		Type: RectType,
		Rect: &Rect_command{
			X:      x,
			Y:      y,
			W:      w,
			H:      h,
			Clr:    clr,
			radius: radius,
			TexId:  texId,
			shape:  shape,
		},
	}
	c.AddCommand(cmd)
}

func (c *CmdBuffer) AddCommand(cmd Command) {
	c.commands = append(c.commands, cmd)

	switch cmd.Type {
	case SeparateBuffer:
		c.Inf = append(c.Inf, Info{
			Elems:       c.VertCount - c.lastElems,
			IndexOffset: c.ofs,
			TexId:       cmd.sb.texid,
			ClipRect:    cmd.sb.clipRect,
		})
		c.ofs += c.VertCount - c.lastElems
		c.lastElems = c.VertCount
	case RectType:
		r := cmd.Rect
		if r.radius == 0 {
			if r.TexId == 0 {
				c.RectangleR(r.X, c.displaySize.Y-r.Y, r.W, r.H, r.Clr)
			} else {
				c.RectangleT(r.X, c.displaySize.Y-r.Y, r.W, r.H, uint32(r.TexId), 0, 1, r.Clr)
				c.SeparateBuffer(r.TexId, c.InnerWindowSpace) // don't forget to slice buffer
			}
		} else {
			if r.TexId == 0 {
				c.RoundedRectangleR(r.X, c.displaySize.Y-r.Y, r.W, r.H, r.radius, r.shape, r.Clr)
			} else {
				// TODO: Add textured rounded rect
			}
		}
	case Text:
		t := cmd.Text
		// size := c.camera.GetProjectionSize()
		// c.Text(t.Text, t.Font, t.X, size.Y()-t.Y, t.Size, t.Clr)
		c.Text(t.Text, t.Font, t.X, c.displaySize.Y-t.Y, t.Size, t.Clr)
		c.SeparateBuffer(t.Font.TextureId, c.InnerWindowSpace) // don't forget to slice buffer
	}

}

func (r *CmdBuffer) render(vert []float32, indeces []int32, vertCount int) {
	r.Vertices = append(r.Vertices, vert...)
	r.Indeces = append(r.Indeces, indeces...)
	r.VertCount += vertCount
}

func (r *CmdBuffer) RectangleR(x, y, w, h float32, clr [4]float32) {

	vert := make([]float32, 9*4)
	ind := make([]int32, 6)

	ind0 := r.lastIndc
	ind1 := ind0 + 1
	ind2 := ind1 + 1
	offset := 0

	fillVertices(vert, &offset, x, y, 0, 0, 0, clr)
	fillVertices(vert, &offset, x, y-h, 0, 0, 0, clr)
	fillVertices(vert, &offset, x+w, y-h, 0, 0, 0, clr)

	ind[0] = int32(ind0)
	ind[1] = int32(ind1)
	ind[2] = int32(ind2)

	last := ind2 + 1

	fillVertices(vert, &offset, x+w, y, 0, 0, 0, clr)
	ind[3] = int32(ind0)
	ind[4] = int32(ind2)
	ind[5] = int32(last)

	r.lastIndc = last + 1
	r.render(vert, ind, 6)
}

func (b *CmdBuffer) Text(text string, font fonts.Font, x, y float32, size int, clr [4]float32) {

	texId := font.TextureId
	inf := font.GetXHeight()

	faceHeight := font.Face.Metrics().Height

	var dx, dy float32
	dx = x
	prevR := rune(-1)
	scale := float32(1) / (float32(font.DefaultFontSize) / float32(size))
	dy = y - scale*inf

	var maxDescend float32
	for _, r := range text {
		info := font.GetCharacter(r)
		if info.Width == 0 {
			log.Printf("Unknown char = %q", r)
			continue
		}
		if prevR >= 0 {
			kern := font.Face.Kern(prevR, r).Ceil()
			dx += float32(kern)
			// fmt.Printf("%q %q %d \n", prevR, r, kern)
		}
		if r == '\n' {
			dx = x
			dy -= float32(faceHeight.Ceil())
			prevR = rune(-1)
			continue
		}
		xPos := float32(dx)
		yPos := float32(dy)

		if info.Descend != 0 {
			d := float32(info.Descend) * scale
			yPos -= d
			if d > maxDescend {
				maxDescend = d
			}
		}
		b.addCharacter(xPos, yPos, scale, uint32(texId), info, clr)
		dx += float32(info.Width) * float32(scale)
		prevR = r
	}
}

func (b *CmdBuffer) addCharacter(x, y float32, scale float32, texId uint32, info fonts.CharInfo, clr [4]float32) {

	vert := make([]float32, 9*4)
	ind := make([]int32, 6)

	x0 := x
	y0 := y
	x1 := x + scale*float32(info.Width)
	y1 := y + scale*float32(info.Heigth)

	ux0, uy0 := info.TexCoords[0].X, info.TexCoords[0].Y
	ux1, uy1 := info.TexCoords[1].X, info.TexCoords[1].Y

	ind0 := b.lastIndc
	ind1 := ind0 + 1
	ind2 := ind1 + 1
	offset := 0

	fillVertices(vert, &offset, x1, y0, ux1, uy0, float32(texId), clr)
	fillVertices(vert, &offset, x1, y1, ux1, uy1, float32(texId), clr)
	fillVertices(vert, &offset, x0, y1, ux0, uy1, float32(texId), clr)

	// h := scale*float32(info.Heigth)
	// w := scale*float32(info.Width)

	// fillVertices(vert, &offset, x0, y0, ux1, uy1, float32(texId), clr)
	// fillVertices(vert, &offset, x0, y0-h, ux1, uy0, float32(texId), clr)
	// fillVertices(vert, &offset, x0+w, y0-h, ux0, uy0, float32(texId), clr)

	ind[0] = int32(ind0)
	ind[1] = int32(ind1)
	ind[2] = int32(ind2)

	last := ind2 + 1

	fillVertices(vert, &offset, x0, y0, ux0, uy0, float32(texId), clr)
	// fillVertices(vert, &offset, x0+w, y0, ux1, uy0, float32(texId), clr)
	ind[3] = int32(ind0)
	ind[4] = int32(ind2)
	ind[5] = int32(last)

	b.lastIndc = last + 1
	b.render(vert, ind, 6)
}

func (r *CmdBuffer) RectangleT(x, y, w, h float32, texId uint32, uv0, uv1 float32, clr [4]float32) {

	vert := make([]float32, 9*4)
	ind := make([]int32, 6)

	ind0 := r.lastIndc
	ind1 := ind0 + 1
	ind2 := ind1 + 1
	offset := 0

	fillVertices(vert, &offset, x, y, uv1, uv1, float32(texId), clr)
	fillVertices(vert, &offset, x, y-h, uv1, uv0, float32(texId), clr)
	fillVertices(vert, &offset, x+w, y-h, uv0, uv0, float32(texId), clr)

	ind[0] = int32(ind0)
	ind[1] = int32(ind1)
	ind[2] = int32(ind2)

	last := ind2 + 1

	fillVertices(vert, &offset, x+w, y, uv0, uv1, float32(texId), clr)
	ind[3] = int32(ind0)
	ind[4] = int32(ind2)
	ind[5] = int32(last)

	r.lastIndc = last + 1
	r.render(vert, ind, 6)
}

func fillVertices(vert []float32, startOffset *int, x, y, uv0, uv1, texId float32, clr [4]float32) {
	offset := *startOffset
	vert[offset] = x
	vert[offset+1] = y

	vert[offset+2] = clr[0] / 255
	vert[offset+3] = clr[1] / 255
	vert[offset+4] = clr[2] / 255
	vert[offset+5] = clr[3]

	vert[offset+6] = uv0
	vert[offset+7] = uv1

	vert[offset+8] = texId

	*startOffset += 9
}

type CircleSector int
type RoundedRectShape int

const (
	TopLeftRect RoundedRectShape = 1 << iota
	TopRigthRect
	BotLeftRect
	BotRightRect

	TopRect = TopLeftRect | TopRigthRect
	BotRect = BotLeftRect | BotRightRect

	AllRounded = TopRect | BotRect
	StraightCorners
)

const (
	BotLeft CircleSector = iota
	BotRight
	TopLeft
	TopRight
)

func (r *CmdBuffer) DrawArc(x, y, radius float32, steps int, sector CircleSector, clr [4]float32) {
	ind0 := r.lastIndc
	ind1 := ind0 + 1
	ind2 := ind1 + 1
	offset := 0
	indOffset := 0

	angle := math.Pi * 2 / float32(steps)

	numV := int(math.Floor(1.57 / float64(angle)))

	ind := make([]int32, 3*(numV+1))    // 3 - triangle
	vert := make([]float32, 9*(3+numV)) //polygon

	var prevX, prevY, lastX, lastY float32

	var ang float32 = angle
	var sX func(x, radius float32) float32
	var sY func(y, radius float32) float32
	// counterTriangles := 0
	switch sector {
	case BotLeft:
		sX = func(x, ang float32) float32 {
			return x - float32(radius)*float32(math.Sin(float64(ang)))
		}
		sY = func(y, ang float32) float32 {
			return y - float32(radius)*float32(math.Cos(float64(ang)))
		}
		prevX = x
		prevY = y - radius
		lastX = x - radius
		lastY = y
	case BotRight:
		sX = func(x, ang float32) float32 {
			return x + float32(radius)*float32(math.Sin(float64(ang)))
		}
		sY = func(y, ang float32) float32 {
			return y - float32(radius)*float32(math.Cos(float64(ang)))
		}
		prevX = x
		prevY = y - radius
		lastX = x + radius
		lastY = y
	case TopLeft:
		sX = func(x, ang float32) float32 {
			return x - float32(radius)*float32(math.Sin(float64(ang)))
		}
		sY = func(y, ang float32) float32 {
			return y + float32(radius)*float32(math.Cos(float64(ang)))
		}
		prevX = x
		prevY = y + radius
		lastX = x - radius
		lastY = y
	case TopRight:
		sX = func(x, ang float32) float32 {
			return x + float32(radius)*float32(math.Sin(float64(ang)))
		}
		sY = func(y, ang float32) float32 {
			return y + float32(radius)*float32(math.Cos(float64(ang)))
		}
		prevX = x
		prevY = y + radius
		lastX = x + radius
		lastY = y
	}

	fillVertices(vert, &offset, x, y, 0, 0, 0, clr)
	fillVertices(vert, &offset, prevX, prevY, 0, 0, 0, clr)
	newx := sX(x, ang)
	newY := sY(y, ang)
	fillVertices(vert, &offset, newx, newY, 0, 0, 0, clr)
	ind[indOffset] = int32(ind0)
	ind[indOffset+1] = int32(ind1)
	ind[indOffset+2] = int32(ind2)
	indOffset += 3
	// ind = append(ind, int32(ind0), int32(ind1), int32(ind2))
	ind1++
	ind2++
	ang += angle

	vertC := 1
	for ang <= 1.57 { // 90 degress ~= 1.57 radians
		newx := sX(x, ang)
		newY := sY(y, ang)

		fillVertices(vert, &offset, newx, newY, 0, 0, 0, clr)

		ind[indOffset] = int32(ind0)
		ind[indOffset+1] = int32(ind1)
		ind[indOffset+2] = int32(ind2)
		indOffset += 3
		ind1++
		ind2++

		ang += angle
		vertC++
		// counterTriangles++
	}
	fillVertices(vert, &offset, lastX, lastY, 0, 0, 0, clr)

	ind[indOffset] = int32(ind0)
	ind[indOffset+1] = int32(ind1)
	ind[indOffset+2] = int32(ind2)
	// indOffset += 3

	r.lastIndc = ind2 + 1

	r.render(vert, ind, (numV+1)*3)
}

var steps = 30

func (r *CmdBuffer) RoundedRectangleR(x, y, w, h float32, radius int, shape RoundedRectShape, clr [4]float32) {

	topLeft := utils.Vec2{x + float32(radius), y - float32(radius)} //origin of arc
	topRight := utils.Vec2{x + w - float32(radius), y - float32(radius)}
	botLeft := utils.Vec2{x + float32(radius), y - h + float32(radius)}
	botRight := utils.Vec2{x + w - float32(radius), y - h + float32(radius)}

	switch shape {
	case TopLeftRect:
		r.DrawArc(topLeft.X, topLeft.Y, float32(radius), steps, TopLeft, clr)
		r.RectangleR(x, y-float32(radius), w, h-float32(radius), clr)               //main rect
		r.RectangleR(x+float32(radius), y, w-float32(radius), float32(radius), clr) //top rect
	case TopRigthRect:
		r.DrawArc(topRight.X, topRight.Y, float32(radius), steps, TopRight, clr)
		r.RectangleR(x, y-float32(radius), w, h-float32(radius), clr) //main
		r.RectangleR(x, y, w-float32(radius), float32(radius), clr)
	case BotLeftRect:
		r.DrawArc(botLeft.X, botLeft.Y, float32(radius), steps, BotLeft, clr)
		r.RectangleR(x, y, w, h-float32(radius), clr) //main
		r.RectangleR(botLeft.X, botLeft.Y, w-float32(radius), float32(radius), clr)
	case BotRightRect:
		r.DrawArc(botRight.X, botRight.Y, float32(radius), steps, BotRight, clr)
		r.RectangleR(x, y, w, h-float32(radius), clr) //main
		r.RectangleR(x, botLeft.Y, w-float32(radius), float32(radius), clr)
	case TopRect:
		r.DrawArc(topLeft.X, topLeft.Y, float32(radius), steps, TopLeft, clr)
		r.DrawArc(topRight.X, topRight.Y, float32(radius), steps, TopRight, clr)
		r.RectangleR(x, y-float32(radius), w, h-float32(radius), clr)
		r.RectangleR(x+float32(radius), y, w-float32(radius)*2, float32(radius), clr)
	case BotRect:
		r.DrawArc(botLeft.X, botLeft.Y, float32(radius), steps, BotLeft, clr)
		r.DrawArc(botRight.X, botRight.Y, float32(radius), steps, BotRight, clr)
		r.RectangleR(x, y, w, h-float32(radius), clr) //main
		r.RectangleR(botLeft.X, botLeft.Y, w-float32(radius)*2, float32(radius), clr)
	case AllRounded:
		r.DrawArc(topLeft.X, topLeft.Y, float32(radius), steps, TopLeft, clr)
		r.DrawArc(topRight.X, topRight.Y, float32(radius), steps, TopRight, clr)
		r.DrawArc(botLeft.X, botLeft.Y, float32(radius), steps, BotLeft, clr)
		r.DrawArc(botRight.X, botRight.Y, float32(radius), steps, BotRight, clr)

		r.RectangleR(topLeft.X, topLeft.Y+float32(radius), w-float32(radius)*2, float32(radius), clr) //top
		r.RectangleR(x, topLeft.Y, w, h-float32(radius)*2, clr)                                       //center
		r.RectangleR(botLeft.X, botLeft.Y, w-float32(radius)*2, float32(radius), clr)                 //bottom
	}

}
