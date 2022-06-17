package gogl

import (
	"github.com/go-gl/mathgl/mgl32"
)

type Sprite struct {
	texture   *Texture     `json:"texture"`
	texCoords []mgl32.Vec2 `json:"tex_coords"`
	width     int32      `json:"sprite_width"`
	height    int32      `json:"sprite_height"`
}

type SpriteExported struct {
	Texture   *TextureExported `json:"texture,omitempty"`
	TexCoords []mgl32.Vec2     `json:"tex_coords"`
	Width     int32          `json:"sprite_width"`
	Height    int32          `json:"sprite_height"`
}

var defaultTexCoords = []mgl32.Vec2{
	{1, 0},
	{1, 1},
	{0, 1},
	{0, 0},
}

func DefSprite() *Sprite {
	return &Sprite{
		texCoords: defaultTexCoords,
		// texture: &Texture{},
	}
}

// func NewSprite(tex *Texture) *Sprite {
// 	sp := Sprite{
// 		texture: tex,
// 		texCoords: defaultTexCoords,
// 	}
// 	return &sp
// }
func (s *Sprite) SetTexture(tex *Texture) {
	s.texture = tex
}

func (s *Sprite) GetWidth() int32 {
	return s.width
}
func (s *Sprite) GetHeight() int32 {
	return s.height
}
func (s *Sprite) SetHeight(h int32) {
	s.height = h
}
func (s *Sprite) SetWidth(w int32) {
	s.width = w
}

func (s *Sprite) SetTexCoords(texCoords []mgl32.Vec2) {
	s.texCoords = texCoords
}

func (s *Sprite) GetTexture() *Texture {
	return s.texture
}
func (s *Sprite) GetTexId() int {
	if s.GetTexture() != nil {
		return int(s.GetTexture().GetId())
	} else {
		return -1
	}
}

func (s *Sprite) GetTextureCoords() []mgl32.Vec2 {
	return s.texCoords
}
