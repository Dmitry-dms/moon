package gogl

import (
	"github.com/go-gl/mathgl/mgl32"
)

type Spritesheet struct {
	texture *Texture
	sprites []*Sprite
}

func NewSpritesheet(tex *Texture, spriteWidth, spriteHeight, numberSprites, spacing int) *Spritesheet{
	sprites := make([]*Sprite, 0)
	sh := Spritesheet{
		texture: tex,
	}
	currentX := 0
	currentY := 0 
	// проход начинается с лево на право верхнего ряда

	for i := 0; i < numberSprites; i++ {
		var topY float32 = (float32(currentY) + float32(spriteHeight)) / float32(tex.GetHeight())
		var rightX float32 = (float32(currentX) + float32(spriteWidth)) / float32(tex.GetWidth())
		var leftX float32 = float32(currentX) / float32(tex.GetWidth())
		var bottomY float32 = float32(currentY) / float32(tex.GetHeight())

		texCoords := []mgl32.Vec2{
			{rightX, bottomY},
			{rightX, topY},
			{leftX, topY},
			{leftX, bottomY},
		}

		sprite := NewSprite(tex)
		sprite.ReplaceTexCoords(texCoords)
		sprites = append(sprites, sprite)

		currentX += spriteWidth + spacing
		if currentX >= tex.GetWidth() {//достигли конца текстуры
			currentX = 0
			currentY += spriteHeight + spacing //переход на ряд выше
		}
	}
	sh.sprites = sprites

	return &sh
}

func (sh *Spritesheet) GetSprite(index int) *Sprite{
	return sh.sprites[index]
}
