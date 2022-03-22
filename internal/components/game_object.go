package components

import (
	"github.com/Dmitry-dms/moon/pkg/gogl"
	mgl "github.com/go-gl/mathgl/mgl32"
	"github.com/inkyblackness/imgui-go/v4"
)

//"fmt"

type GameObject struct {
	name string

	Transform     *Transform
	LastTransform *Transform

	Spr *SpriteRenderer

	isDirty bool

	zIndex int
}

func NewGameObject(name string, transform *Transform, zIndex int) *GameObject {
	//fmt.Println("Creating " + name)
	obj := GameObject{
		name:          name,
		Transform:     transform,
		LastTransform: transform.Copy(),
		zIndex:        zIndex,
	}
	return &obj
}

func (g *GameObject) AddSpriteRenderer(spr *SpriteRenderer) {
	g.Spr = spr
}

func (g *GameObject) Update(dt float32) {
	if g.Transform != g.LastTransform {
		g.Transform.CopyTo(g.LastTransform)
		g.isDirty = true
	}
}

func (g *GameObject) GetZIndex() int {
	return g.zIndex
}

func (g *GameObject) SetColor(color mgl.Vec4) {
	if g.Spr.color != color {
		g.isDirty = true
		g.SetColor(color)
	}
}

func (g *GameObject) SetSprite(sprite *gogl.Sprite) {
	g.Spr.SetSprite(sprite)
	g.isDirty = true
}

func (g *GameObject) AddPosition(tr mgl.Vec2) {
	g.LastTransform = g.Transform.Copy()
	g.Transform.Position[0] += tr.X()
	g.Transform.Position[1] += tr.Y()
	g.isDirty = true
}

func (g *GameObject) IsDirty() bool {
	return g.isDirty
}

func (g *GameObject) SetClean() {
	g.isDirty = false
}

func (g *GameObject) Imgui() {
	colors := [4]float32{g.Spr.color[0],g.Spr.color[1],g.Spr.color[2],g.Spr.color[3]}
	if imgui.ColorPicker4("Color picker", &colors) {
		g.Spr.SetColor(mgl.Vec4{colors[0],colors[1],colors[2],colors[3]})
		g.isDirty = true
	}
}

// func (g *GameObject) Start() {
// 	for _, c := range g.components {
// 		c.Start()
// 	}
// }
