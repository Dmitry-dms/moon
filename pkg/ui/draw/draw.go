package draw

import (
	"math"

	"github.com/Dmitry-dms/moon/pkg/gogl"
	"github.com/Dmitry-dms/moon/pkg/ui/utils"
)

type DrawData struct {
	cmdLists                []CmdBuffer
	DisplayPos, DisplaySize utils.Vec2
}

func drawData() DrawData {
	return DrawData{
		cmdLists:    []CmdBuffer{},
		DisplayPos:  utils.Vec2{},
		DisplaySize: utils.Vec2{},
	}
}

func (d *DrawData) clear() {
	d.cmdLists = []CmdBuffer{}
}

type CmdBuffer struct {
	commands  []Command
	Vertices  []float32
	Indeces   []int32
	VertCount int
	lastIndc  int

	camera *gogl.Camera
}

func NewBuffer(camera *gogl.Camera) *CmdBuffer {
	return &CmdBuffer{
		commands:  []Command{},
		Vertices:  []float32{},
		Indeces:   []int32{},
		VertCount: 0,
		camera:    camera,
	}
}

func (c *CmdBuffer) Clear() {
	c.commands = []Command{}
	c.Vertices = []float32{}
	c.Indeces = []int32{}
	c.VertCount = 0
	c.lastIndc = 0
}

func (c *CmdBuffer) AddCommand(cmd Command) {
	c.commands = append(c.commands, cmd)

	switch cmd.Type {
	case RectType:
		r := cmd.Rect
		size := c.camera.GetProjectionSize()
		c.RectangleR(r.X, size.Y()-r.Y, r.W, r.H, r.Clr)
	case RectTypeT:
		r := cmd.Rect

		size := c.camera.GetProjectionSize()
		if r.ScaleFactor == 0 {
			c.RectangleT(r.X, size.Y()-r.Y, r.W, r.H, r.Texture, 0, 1, 0, r.Clr)
		} else {
			c.RectangleT(r.X, size.Y()-r.Y, r.W, r.H, r.Texture, 0, 1, r.ScaleFactor, r.Clr)
		}
	case RoundedRect:
		rr := cmd.RRect
		size := c.camera.GetProjectionSize()
		c.RoundedRectangleR(rr.X, size.Y()-rr.Y, rr.W, rr.H, rr.Radius, AllRounded, rr.Clr)
	case ScrollbarCmd:
		rr := cmd.RRect
		size := c.camera.GetProjectionSize()
		c.RoundedRectangleR(rr.X, size.Y()-rr.Y, rr.W, rr.H, rr.Radius, AllRounded, rr.Clr)
	case ScrollButtonCmd:
		rr := cmd.RRect
		size := c.camera.GetProjectionSize()
		c.RoundedRectangleR(rr.X, size.Y()-rr.Y, rr.W, rr.H, rr.Radius, AllRounded, rr.Clr)
	case WindowStartCmd:
		wnd := cmd.Window
		size := c.camera.GetProjectionSize()
		c.RoundedRectangleR(wnd.X, size.Y()-wnd.Y, wnd.W, wnd.H, 10, AllRounded, wnd.Clr)
		c.RoundedRectangleR(wnd.X, size.Y()-wnd.Y, wnd.W, wnd.Toolbar.H, 10, TopRect, wnd.Toolbar.Clr)

		// c.RectangleR(wnd.X, size.Y()-wnd.Y, wnd.W, wnd.H, cmd.Window.Clr)
		// c.RectangleR(wnd.X, size.Y()-wnd.Y, wnd.W, wnd.Toolbar.H, cmd.Window.Toolbar.Clr)
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

func (r *CmdBuffer) RectangleT(x, y, w, h float32, tex *gogl.Texture, uv0, uv1, f float32, clr [4]float32) {
	// founded := false
	texId := 0
	// for i := 0; i < len(r.textures); i++ {
	// 	if r.textures[i] == tex {
	// 		texId = i + 1 // 0 - без текстуры
	// 		founded = true
	// 	}
	// }
	// if !founded {
	// 	r.addTexture(tex)
	// }

	vert := make([]float32, 9*4)
	ind := make([]int32, 6)

	ind0 := r.lastIndc
	ind1 := ind0 + 1
	ind2 := ind1 + 1
	offset := 0

	// var v00, v10, v11, v01 float32
	// v00 = uv0
	// v11 = uv1
	// v10 = uv1
	// v01 = uv0

	// factor := float32(1)

	if f != 0 {
		// factor = f
		// uv1 -= 0.4
		// h = f
	} else {
		f = 1
	}

	fillVertices(vert, &offset, x, y, uv1, uv1, float32(texId), clr)
	fillVertices(vert, &offset, x, y-h, uv1, uv0+(1-f), float32(texId), clr)
	fillVertices(vert, &offset, x+w, y-h, uv0, uv0+(1-f), float32(texId), clr)

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
		r.RectangleR(x, topLeft.Y, w, h-float32(radius)*2, clr)                                         //center
		r.RectangleR(botLeft.X, botLeft.Y, w-float32(radius)*2, float32(radius), clr)                 //bottom
	}

}