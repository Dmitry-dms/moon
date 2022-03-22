package components

import mgl "github.com/go-gl/mathgl/mgl32"

type Transform struct {
	Position mgl.Vec2
	Scale    mgl.Vec2
}

func NewTransform(pos, scale mgl.Vec2) *Transform {
	return &Transform{Position: pos, Scale: scale}
}

func (t *Transform) Copy() *Transform {
	return &Transform{Position: t.Position, Scale: t.Scale}
}

func (t *Transform) CopyTo(to *Transform) {
	to.Position = t.Position
	to.Scale = t.Scale
}
