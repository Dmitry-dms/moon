package components

import (
	"encoding/json"
	"fmt"
	"reflect"

	"github.com/Dmitry-dms/moon/pkg/gogl"
	mgl "github.com/go-gl/mathgl/mgl32"
	"github.com/inkyblackness/imgui-go/v4"
)

type GameObject struct {
	name string

	Transform     *Transform
	LastTransform *Transform

	Spr SpriteRenderer

	isDirty bool
	zIndex  int

	uid int
}

func NewGameObject(name string, transform *Transform, zIndex int) *GameObject {
	obj := GameObject{
		name:          name,
		Transform:     transform,
		LastTransform: transform.Copy(),
		zIndex:        zIndex,
		uid:           -1,
	}
	obj.generateId()
	return &obj
}

func (g *GameObject) AddSpriteRenderer(spr *SpriteRenderer) {
	g.Spr = *spr
}

func (g *GameObject) generateId() {
	if g.uid == -1 {
		id := ID_COUNTER
		ID_COUNTER++
		g.uid = id
	}
}
func (g *GameObject) GetUid() int {
	return g.uid
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

func (g *GameObject) SetColor(color mgl.Vec4) {
	if g.Spr.color != color {
		g.isDirty = true
		g.SetColor(color)
	}
}

func (g *GameObject) SetSprite(sprite *gogl.Sprite) {
	g.Spr.SetSprite(sprite)
	g.isDirty = true
}

func (g *GameObject) AddPosition(tr mgl.Vec2) {
	g.LastTransform = g.Transform.Copy()
	g.Transform.position[0] += tr.X()
	g.Transform.position[1] += tr.Y()
	g.isDirty = true
}
func (g *GameObject) SetPosition(tr mgl.Vec2) {
	g.Transform.SetPosition(tr)
	g.isDirty = true
}

func (g *GameObject) IsDirty() bool {
	return g.isDirty
}

func (g *GameObject) SetClean() {
	g.isDirty = false
}

func (g *GameObject) Imgui() {
	colors := [4]float32{g.Spr.color[0], g.Spr.color[1], g.Spr.color[2], g.Spr.color[3]}
	if imgui.ColorPicker4("Color picker", &colors) {
		g.Spr.SetColor(mgl.Vec4{colors[0], colors[1], colors[2], colors[3]})
		g.isDirty = true
	}

}

func reflObj(g *GameObject) {
	v := reflect.ValueOf(*g)

	for i := 0; i < v.NumField(); i++ {
		field := v.Field(i)

		switch field.Kind() {
		case reflect.Int:
			var val int32 = int32(g.zIndex)
			if imgui.DragInt("zIndex", &val) {
				g.zIndex = int(val)
				g.isDirty = true
			}
		case reflect.ValueOf(mgl.Vec2{}).Kind():
			var fl [2]float32
			if imgui.DragFloat2("floats", &fl) {

			}
		}
	}

}

type gameObjExported struct {
	Name      string `json:"name"`
	Transform tranformExported
	Spr       spriteRendererExported
	ZIndex    int `json:"z_index"`
	Uid       int `json:"uid"`
}

func (g *GameObject) MarshalJSON() ([]byte, error) {
	var transfExp tranformExported
	var spriteExported *gogl.SpriteExported
	var sprExported spriteRendererExported
	var texExported *gogl.TextureExported
	transfExp = tranformExported{
		Position: g.Transform.position,
		Scale:    g.Transform.scale,
	}
	if g.Spr.GetTexture() != nil {
		texExported = &gogl.TextureExported{
			Filepath:  g.Spr.GetTexture().GetFilepath(),
			TextureId: g.Spr.GetTexture().GetId(),
			Width:     g.Spr.GetTexture().GetWidth(),
			Height:    g.Spr.GetTexture().GetHeight(),
		}
		spriteExported = &gogl.SpriteExported{
			Texture:   texExported,
			TexCoords: g.Spr.GetTextureCoords(),
			Width:     g.Spr.sprite.GetWidth(),
			Height:    g.Spr.sprite.GetHeight(),
		}
	} else {
		spriteExported = &gogl.SpriteExported{
			Texture:   nil,
			TexCoords: g.Spr.GetTextureCoords(),
		}
	}

	sprExported = spriteRendererExported{
		Color:  g.Spr.GetColor(),
		Sprite: spriteExported,
	}

	rt := gameObjExported{
		Name:      g.name,
		Transform: transfExp,
		Spr:       sprExported,
		ZIndex:    g.zIndex,
		Uid: g.GetUid(),
	}
	return json.Marshal(rt)
}
func (g *GameObject) UnmarshalJSON(data []byte) error {
	var tr gameObjExported
	err := json.Unmarshal(data, &tr)
	if err != nil {
		fmt.Println(err)
		return err
	}
	g.isDirty = true
	g.name = tr.Name
	g.zIndex = tr.ZIndex
	var newTransform Transform
	newTransform.position = tr.Transform.Position
	newTransform.scale = tr.Transform.Scale
	g.Transform = &newTransform
	g.LastTransform = &newTransform
	var spr SpriteRenderer
	spr.color = tr.Spr.Color
	sprite := gogl.DefSprite()

	if tr.Spr.Sprite.Texture != nil {
		tex := gogl.CreateTexture(tr.Spr.Sprite.Texture.Filepath, tr.Spr.Sprite.Texture.TextureId,
			tr.Spr.Sprite.Texture.Width, tr.Spr.Sprite.Texture.Height)
		sprite.SetTexture(tex)
		sprite.SetTexCoords(tr.Spr.Sprite.TexCoords)
		sprite.SetWidth(tr.Spr.Sprite.Width)
		sprite.SetHeight(tr.Spr.Sprite.Height)
	}

	spr.sprite = sprite
	g.Spr = spr
	g.uid = tr.Uid
	return err
}

// func (g *GameObject) Start() {
// 	for _, c := range g.components {
// 		c.Start()
// 	}
// }
