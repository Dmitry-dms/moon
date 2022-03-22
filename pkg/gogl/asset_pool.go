package gogl

import (
	"sync"
)

func init() {
	s := sync.Once{}
	s.Do(func() {
		AssetPool = newAssetPool()
	})
}

var AssetPool *assetPool

type assetPool struct {
	shaders     map[string]*Shader
	textures    map[string]*Texture
	spritesheet map[string]*Spritesheet
}

func newAssetPool() *assetPool {
	p := assetPool{
		shaders:  make(map[string]*Shader),
		textures: make(map[string]*Texture),
		spritesheet: make(map[string]*Spritesheet),
	}
	return &p
}

func (p *assetPool) GetShader(resourcePath string) (*Shader, error) {
	sh, ok := p.shaders[resourcePath]
	if !ok {
		shader, err := NewShader(resourcePath)
		if err == nil {
			p.shaders[resourcePath] = shader
		}
		return shader, err
	}
	return sh, nil
}
func (p *assetPool) GetTexture(resourcePath string) *Texture {
	sh, ok := p.textures[resourcePath]
	if !ok {
		tex, err := LoadTextureAlpha(resourcePath)
		if err == nil {
			p.textures[resourcePath] = tex
		}
		return tex
	}
	return sh
}

func (p *assetPool) AddSpritesheet(resourceName string, sh *Spritesheet) {
	_, ok := p.spritesheet[resourceName]
	if !ok {
		p.spritesheet[resourceName] = sh
	}
}
func (p *assetPool) GetSpriteSheet(resourceName string) *Spritesheet {
	sh, ok := p.spritesheet[resourceName]
	if !ok {
		return nil
	}
	return sh
}
