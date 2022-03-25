package components

import (
	"github.com/Dmitry-dms/moon/pkg/gogl"

	mgl "github.com/go-gl/mathgl/mgl32"
)

type SpriteRenderer struct {
	color  mgl.Vec4 `json:"color"`
	sprite *gogl.Sprite `json:"sprite"`
}

var defaultColor = mgl.Vec4{1, 1, 1, 1}

func DefSpriteRenderer() *SpriteRenderer {
	return &SpriteRenderer{
		color: defaultColor,
	}
}

// func NewSpriteRenderer(color mgl.Vec4, sprite *gogl.Sprite) *SpriteRenderer {
// 	spr := SpriteRenderer{
// 		color:  color,
// 		sprite: sprite,
// 	}
// 	return &spr
// }
type spriteRendererExported struct {
	Color  mgl.Vec4            `json:"color"`
	Sprite *gogl.SpriteExported `json:"sprite,omitempty"`
}

func (r *SpriteRenderer) GetTexture() *gogl.Texture {
	return r.sprite.GetTexture()
}
func (r *SpriteRenderer) GetTextureCoords() []mgl.Vec2 {
	return r.sprite.GetTextureCoords()
}

func (r *SpriteRenderer) Update(dt float32) {

}
func (r *SpriteRenderer) GetColor() mgl.Vec4 {
	return r.color
}

// func (r *SpriteRenderer) Start() {

// }

func (r *SpriteRenderer) SetSprite(sprite *gogl.Sprite) {
	r.sprite = sprite
}

func (r *SpriteRenderer) SetColor(color mgl.Vec4) {
	r.color = color
}
