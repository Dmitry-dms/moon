package ui

import (
	// "fmt"


	"github.com/Dmitry-dms/moon/internal/listeners"
	"github.com/go-gl/glfw/v3.3/glfw"
	"github.com/go-gl/mathgl/mgl32"
)

type Button struct {
	*UiObject
	Hot bool
}

func (b *Button) Spr() *SpriteRenderer {
	return b.UiObject.Spr
}

func (b *Button) Transform() *Transform {
	return b.UiObject.Transform
}

var dragStarted mgl32.Vec2

func (b *Button) Update(dt float32) {
	pos := b.Transform().position
	scale := b.Transform().scale

	if listeners.RegionHit(pos.X(), pos.Y(), scale.X(), scale.Y()) {
		b.UiObject.Spr.SetColor(mgl32.Vec4{1, 0, 0, 1})
		b.Hot = true
		if listeners.MouseButtonDown(glfw.MouseButton1) && listeners.IsDragging() {
			b.isMovable = true
			
		} else {
			// b.isMovable = false
			dragStarted = mgl32.Vec2{0, 0}
		}

	} else {
		b.UiObject.Spr.SetColor(mgl32.Vec4{0, 1, 0, 1})
		b.Hot = false
		// b.isMovable = false
	}


	if b.isMovable {

		if dragStarted.X() == 0 && dragStarted.Y() == 0 {
			dragStarted = listeners.VecPos()
		}

		before := b.UiObject.Transform.position
		mPos := listeners.VecPos()
		newPos := mgl32.Vec2{before[0] + mPos[0] - dragStarted[0], before[1] + mPos[1] - dragStarted[1]}
		b.UiObject.Transform.SetPosition(newPos)
		dragStarted = mPos

	}

	if !listeners.IsDragging() {
		b.isMovable = false
	}
}
