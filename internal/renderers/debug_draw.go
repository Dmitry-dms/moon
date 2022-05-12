package renderers

import (
	"fmt"
	// "sync"

	"github.com/Dmitry-dms/moon/pkg/gogl"
	gom "github.com/Dmitry-dms/moon/pkg/math"
	"github.com/go-gl/gl/v4.2-core/gl"
	"github.com/go-gl/mathgl/mgl32"
)

const max_lines = 5000

var DebugDraw *debugDraw

type debugDraw struct {
	lines       []*gogl.Line2D
	vertexArray []float32
	shader      *gogl.Shader
	vao, vbo    uint32
	started     bool
}

func Init() {
	// s := sync.Once{}
	// s.Do(func() {
	DebugDraw = newDebugDraw()
	DebugDraw.Start()
	//})
}

func newDebugDraw() *debugDraw {
	sh, err := gogl.NewShader("assets/shaders/debugLine2D.glsl")
	if err != nil {
		fmt.Println(err)
	}
	lines := make([]*gogl.Line2D, 0)
	db := debugDraw{
		shader: sh,
		lines:  lines,
		//6 floats per vector, 2 vertices per line
		vertexArray: make([]float32, max_lines*6*2),
	}
	return &db
}

func (d *debugDraw) Start() {
	d.vao = gogl.GenBindVAO()
	d.vbo = gogl.GenBindBuffer(gl.ARRAY_BUFFER)
	gogl.BufferData(gl.ARRAY_BUFFER, d.vertexArray, gl.DYNAMIC_DRAW)
	gogl.SetVertexAttribPointer(0, 3, gl.FLOAT, 6, 0)
	gogl.SetVertexAttribPointer(1, 3, gl.FLOAT, 6, 3)

	d.started = true

	//set line width
	gl.LineWidth(2)
}

func (d *debugDraw) BeginFrame() {
	if !d.started {
		fmt.Println("start")
		d.Start()
	}

	for i := 0; i < len(d.lines); i++ {
		if d.lines[i].BeginFrame() < 0 {
			d.lines = append(d.lines[:i], d.lines[i+1:]...)
		}
	}
}

func (d *debugDraw) Draw(camera *gogl.Camera) {
	if len(d.lines) <= 0 {
		return
	}

	index := 0
	for _, l := range d.lines {
		for i := 0; i < 2; i++ {
			var position mgl32.Vec2
			if i == 0 {
				position = l.From()
			} else {
				position = l.To()
			}
			color := l.Color()

			//load pos into vertexarray
			d.vertexArray[index] = position.X()
			d.vertexArray[index+1] = position.Y()
			d.vertexArray[index+2] = -10

			//color
			d.vertexArray[index+3] = color.X()
			d.vertexArray[index+4] = color.Y()
			d.vertexArray[index+5] = color.Z()
			index += 6
		}
	}

	gl.BindBuffer(gl.ARRAY_BUFFER, d.vbo)
	gl.BufferSubData(gl.ARRAY_BUFFER, 0, len(d.lines)*4*6*2, gl.Ptr(d.vertexArray))
	gl.BindBuffer(gl.VERTEX_ARRAY, 0)

	d.shader.Use()
	d.shader.UploadMat4("uProjection", camera.GetProjectionMatrix())
	d.shader.UploadMat4("uView", camera.GetViewMatrix())

	gl.BindVertexArray(d.vao)
	gl.EnableVertexAttribArray(0)
	gl.EnableVertexAttribArray(1)

	//draw batch
	gl.DrawArrays(gl.LINES, 0, int32(len(d.lines)*6*2))

	gl.DisableVertexAttribArray(0)
	gl.DisableVertexAttribArray(1)
	gl.BindVertexArray(0)

	d.shader.Detach()
}

var defaultColor = mgl32.Vec3{0, 1, 0}

func (d *debugDraw) AddLine2d1(from, to mgl32.Vec2) {
	d.AddLine2d(from, to, defaultColor, 1)
}
func (d *debugDraw) AddLine2d(from, to mgl32.Vec2, color mgl32.Vec3, lifetime int) {
	if len(d.lines) >= max_lines {
		return
	}

	line2d := gogl.NewLine2D(from, to, color, lifetime)
	d.lines = append(d.lines, line2d)
}
func (d *debugDraw) AddBox2D(center, dimensions mgl32.Vec2, rotation float32, color mgl32.Vec3, lifetime int) {
	min := center.Sub(dimensions.Mul(0.5))
	max := center.Add(dimensions.Mul(0.5))

	vert := []*mgl32.Vec2{
		{min.X(), min.Y()}, {min.X(), max.Y()},
		{max.X(), max.Y()}, {max.X(), min.Y()},
	}

	if rotation != 0 {
		for _, v := range vert {
			gom.Rotate(v, rotation, center)
		}
	}
	d.AddLine2d(*vert[0], *vert[1], color, lifetime)
	d.AddLine2d(*vert[1], *vert[2], color, lifetime)
	d.AddLine2d(*vert[2], *vert[3], color, lifetime)
	d.AddLine2d(*vert[3], *vert[0], color, lifetime)
}

func (d *debugDraw) AddCircle(center mgl32.Vec2, radius float32, color mgl32.Vec3, lifetime int) {
	points := [20]mgl32.Vec2{} //segments

	increment := 360 / len(points)
	var currentAngle float32 = 0
	for i := 0; i < len(points); i++ {
		tmp := mgl32.Vec2{radius, 0}
		gom.Rotate(&tmp, currentAngle, mgl32.Vec2{})
		points[i] = tmp.Add(center)

		if i > 0 {
			d.AddLine2d(points[i-1], points[i], color, lifetime)
		}
		currentAngle+=float32(increment)
	}
	d.AddLine2d(points[len(points)-1], points[0], color, lifetime)
}
