package components

import (
	"github.com/Dmitry-dms/moon/pkg/gogl"
	"github.com/go-gl/mathgl/mgl32"
)

func GenerateSpriteObject(sprite *gogl.Sprite, sizeX, sizeY float32) *GameObject {
	gameObj := NewGameObject("Sprite_obj_gen", NewTransform(mgl32.Vec2{0, 0}, mgl32.Vec2{sizeX, sizeY}), 0)
	spr := DefSpriteRenderer()
	spr.SetSprite(sprite)
	gameObj.AddSpriteRenderer(spr)

	return gameObj
}
