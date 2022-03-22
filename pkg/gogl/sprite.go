package gogl

import (
	"github.com/go-gl/mathgl/mgl32"
)

type Sprite struct {
	texture *Texture
	texCoords []mgl32.Vec2
}

var defaultTexCoords = []mgl32.Vec2{
	{1,0},
	{1,1},
	{0,1},
	{0,0},
}

func NewSprite(tex *Texture) *Sprite {
	sp := Sprite{
		texture: tex,
		texCoords: defaultTexCoords,
	}
	return &sp
}

func (s *Sprite) ReplaceTexCoords(texCoords []mgl32.Vec2) {
	s.texCoords = texCoords
}

func (s *Sprite) GetTexture() *Texture {
	return s.texture
}

func (s *Sprite) GetTextureCoords() []mgl32.Vec2 {
	return s.texCoords
}