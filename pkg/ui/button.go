package ui

import (
	"fmt"

	"github.com/Dmitry-dms/moon/internal/listeners"
	"github.com/go-gl/mathgl/mgl32"
)



type Button struct {
	UiObject
	Hot bool
}

func (b *Button) Spr() *SpriteRenderer {
	return b.UiObject.Spr
}

func (b *Button) Transform() *Transform {
	return b.UiObject.Transform
}

func (b *Button) Update(dt float32) {
	pos := b.Transform().position
	scale := b.Transform().scale
	

	if listeners.RegionHit(pos.X(), pos.Y(), scale.X(), scale.Y()) {
		b.UiObject.Spr.SetColor(mgl32.Vec4{1,0,0,1})
		fmt.Println("a ...any")
	} else {
		b.UiObject.Spr.SetColor(mgl32.Vec4{0,1,0,1})
	}
}
