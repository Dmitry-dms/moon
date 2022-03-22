package components

import (
	"github.com/Dmitry-dms/moon/pkg/gogl"
	"github.com/go-gl/mathgl/mgl32"
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

func (g *GameObject) SetColor(color mgl32.Vec4) {
	if g.Spr.color != color {
		g.isDirty = true
		g.SetColor(color)
	}
}

func (g *GameObject) SetSprite(sprite *gogl.Sprite) {
	g.Spr.SetSprite(sprite)
	g.isDirty = true
}

func (g *GameObject) AddPosition(tr mgl32.Vec2) {
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

// func (g *GameObject) Start() {
// 	for _, c := range g.components {
// 		c.Start()
// 	}
// }
